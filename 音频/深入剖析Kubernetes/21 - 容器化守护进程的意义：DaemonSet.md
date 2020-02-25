# 21 | 容器化守护进程的意义：DaemonSet

## 笔记

**StatefulSet**就是对现有典型运维业务的容器化抽象. 相较于传统部署, `k8s`可以更高的管理升级, 版本的事务.

修改`StstefulSet`的`Pod`模板就会自动触发"滚动更新".

```
$ kubectl patch statefulset mysql --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/image", "value":"mysql:5.7.23"}]'
statefulset.apps/mysql patched
```

`StatefulSet Controller`会按照与`Pod`编号相反的顺讯, **从最后一个`Pod`开始,逐一更新每个`Pod`**. 

还可以实现金丝雀更新, 应用的多个实例中被指定的一部分不会被更新到最新的版本. ` spec.updateStrategy.rollingUpdate`的`partition`字段. 只有序号**大于等于**`partition`的`Pod`会被更新到这个版本. 序号小于的`Pod`即使删除了，重启了也不会被影响, 还是保持原来的版本.

### DaemonSet

`Daemon Pod`有三个特征:

1. 这个`Pod`运行在`k8s`集群里的**每一个节点(Node)**上
2. 每个节点上**只有一个**这样的`Pod`实例
3. 当有新的节点加入`k8s`集群, 该`Pod`会自动地在新节点上被创建出来; 当旧节点被删除后, 它上面的`Pod`也相应地会被回收掉.

作用:

1. 各种**网络**插件的`agent`组件, 必须运行在每一个节点上, 用来处理这个节点上的容器网络.
2. 各种**存储**插件的`agent`组件, 必须运行在每一个节点上, 用来在这个节点上挂载远程存储目录, 操作容器的`Volume`目录.
3. 各种**监控**组件和**日志**组件, 必须运行在每一个节点上, 负责这个节点上的监控信息和日志搜集.

**`DaemonSet`开始运行的时机, 很多时候比整个k8s集群出现的时机要早**.

### DaemonSet 的API定义

```
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd-elasticsearch
  namespace: kube-system
  labels:
    k8s-app: fluentd-logging
spec:
  selector:
    matchLabels:
      name: fluentd-elasticsearch
  template:
    metadata:
      labels:
        name: fluentd-elasticsearch
    spec:
      tolerations:
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
      containers:
      - name: fluentd-elasticsearch
        image: k8s.gcr.io/fluentd-elasticsearch:1.20
        resources:
          limits:
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 200Mi
        volumeMounts:
        - name: varlog
          mountPath: /var/log
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
      terminationGracePeriodSeconds: 30
      volumes:
      - name: varlog
        hostPath:
          path: /var/log
      - name: varlibdockercontainers
        hostPath:
          path: /var/lib/docker/containers
```
  
### DaemonSet 如何保证每个 Node 上有且只有一个被管理的 Pod 

`DaemonSet Controller`, 首先从`Etcd`里获取所有的`Node`列表, 然后遍历所有`Node`, 当前这个`Node`是否携带了`name=fluentd-elasticsearch`(上面的例子)的标签的`Pod`在运行.

* 没有, 创建一个
* 有, 但是数量 > 1, 删除多余的
* 有, 数量=1, 正常

#### 如何在指定 Node 上创建新 Pod 呢?

使用`nodeSelector`(要被废弃了).

```
nodeSelector:
    name: <Node名字>
``` 

新版用`nodeAffinity`来替换了.

```

apiVersion: v1
kind: Pod
metadata:
  name: with-node-affinity
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: metadata.name
            operator: In
            values:
            - node-geektime
```

`nodeAffinity`的含义是:

1. `requiredDuringSchedulingIgnoredDuringExecution`, 必须在每次调度时候予以考虑
2. 只允许运行在`metadata.name`是`node-geektime`的节点上

**DaemonSet Controller 会在创建 Pod 的时候, 自动在这个 Pod 的 API 对象里, 加上这样一个 nodeAffinity 定义**. 需要绑定的节点名字, 正式当前正在遍历的这个`Node`. 它并不需要修改用户提交的`YAML`文件里的`Pod`模板, 而是在向`k8s`发起请求你之前, 直接修改根据模板生成的`Pod`对象.

`DaemonSet`还会给这个`Pod`自动加上一个字段`tolerations`, 意味着这个`Pod`, 会"容忍"(`Tolerations`)某些`Node`的"污点"(`Taint`).

```

apiVersion: v1
kind: Pod
metadata:
  name: with-toleration
spec:
  tolerations:
  - key: node.kubernetes.io/unschedulable
    operator: Exists
    effect: NoSchedule
```

"容忍"所有被标记为`unschedulable`"污点"的`Node`; "容忍"的效果是允许调度.

正常情况下, 被标记了`unschedulable`"污点"的`Node`, 是不会有任何`Pod`被调度上去的(effect: NoSchedule). 可是, `DaemonSet`自动地给被管理的`Pod`加上了这个特殊的 `Toleration`, 就使得这些`Pod`可以忽略这个限制, 继而保证每个节点上都会被调度一个 `Pod`. 当然, 如果这个节点有故障的话, 这个`Pod`可能会启动失败, 而`DaemonSet`则会始终尝试下去, 直到`Pod`启动成功.

#### 具体应用示例 - 网络插件

```
...
template:
    metadata:
      labels:
        name: network-plugin-agent
    spec:
      tolerations:
      - key: node.kubernetes.io/network-unavailable
        operator: Exists
        effect: NoSchedule
```

容忍`node.kubernetes.io/network-unavailable`污点的`Toleration`.

在`k8s`项目中, 当一个节点的网络插件尚未安装时, 这个节点就会被自动加上名为`node.kubernetes.io/network-unavailable`的"污点". **而通过这样一个 `Toleration`, 调度器在调度这个`Pod`的时候, 就会忽略当前节点上的"污点", 从而成功地将网络插件的`Agent`组件调度到这台机器上启动起来.**

这种机制, 正式我们先部署`k8s`集群, 在部署网络插件的根本原因.


#### 具体应用示例 - 给`mastet`添加`Pod`

```
tolerations:
- key: node-role.kubernetes.io/master
  effect: NoSchedule
 ```
 
 正常`Master`节点默认携带了一个叫作`node-role.kubernetes.io/master`的"污点". 不允许用户在上`Master`节点部署`Pod`. 可以通过让这个`Pod`容忍这个污点, 而部署在`Master`节点

### DaemonSet 如何维护版本

`Deployment`是通过"一个版本对应一个 ReplicaSet 对象". `DaemonSet`是直接操作的`Pod`, 通过"API对象".

`k8s v1.7`后添加了一个`API`对象, **ControllerRevision**, 专门用来记录某种`Conntroller`对象的版本.

ControllerRevision:

* `data`: 保存了该版本对应的完整的`DaemonSet`的`API`对象.
* `Annotation`: 保存了创建这个对象所使用的`kubectl`命令.

```
$ kubectl rollout undo daemonset fluentd-elasticsearch --to-revision=1 -n kube-system
daemonset.extensions/fluentd-elasticsearch rolled back
```

对现有的`DaemonSet`做一次`PATCH`操作(等价于执行一次 kubectl apply -f "旧的 DaemonSet 对象"), 从而把这个`DaemonSet`"更新"到一个旧版本. 

`Revision`不会回退, 会增加成`3`. 因为, 一个新的`ControllerRevision`被创建出来.

### ControllerRevision

**StatefulSet 也是直接控制`Pod`对象, 也是通过ControllerRevision来进行版本管理**.

## 扩展