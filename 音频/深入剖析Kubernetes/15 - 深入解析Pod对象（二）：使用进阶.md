# 15 | 深入解析Pod对象（二）：使用进阶

## 笔记

### Projected Volume

投射数据卷

在`Kubernetes`中, 几种特殊的`Volume`, 存在的意义不是为了存放容器里的数据, 也不是用来进行容器和宿主机之间的数据交换. 这些特殊的`Volume`的作用, 是为容器提供预先定义好的数据. 从容器角度来看, 这些`Volume`里的信息就是仿佛**被Kubernetes"投射"(Project)进入容器当中的**.

`Projected Volume`有四种:

* `Secret`
* `ConfigMap`
* `Download API`
* `ServiceAccountToken`

### Secret

把`Pod`想要访问的加密数据, 存放到`Etcd`中. 可以通过在`Pod`的容器挂载`Volume`的方式, 访问到这些`Secret`里保存的信息.

可以通过`kubectl create secret`指令创建, 也可以通过`YAML`创建. 这些数据必须是经过`Base64`转码的, 以免出现明文密码的安全隐患.

创建完`Pod`之后, 保存在`Etcd`里的用户名和密码信息, 已经以文件形式出现在了容器的`Volume`目录里. 文件的名字就是, 指定的`Key`.

一旦其对应的`Etcd`里的数据被更新, 这些`Volume`里的文件内容, 同样也会被更新. **更新会有一定的延时, 需要编写好重试和超时的逻辑**.

### ConfigMap

与`Secret`的区别在于, `ConfigMap`保存的是不需要加密的, 应用所需的配置信息.

可以通过`kubectl create configmap`命令创建, 也可以通过`YAML`创建.

### Downward API

让`Pod`里的容器能够直接获取到这个`Pod API`对象本身的信息.

```

apiVersion: v1
kind: Pod
metadata:
  name: test-downwardapi-volume
  labels:
    zone: us-est-coast
    cluster: test-cluster1
    rack: rack-22
spec:
  containers:
    - name: client-container
      image: k8s.gcr.io/busybox
      command: ["sh", "-c"]
      args:
      - while true; do
          if [[ -e /etc/podinfo/labels ]]; then
            echo -en '\n\n'; cat /etc/podinfo/labels; fi;
          sleep 5;
        done;
      volumeMounts:
        - name: podinfo
          mountPath: /etc/podinfo
          readOnly: false
  volumes:
    - name: podinfo
      projected:
        sources:
        - downwardAPI:
            items:
              - path: "labels"
                fieldRef:
                  fieldPath: metadata.labels
```

`Volume`的数据来源变成了`Downward API`, 声明了要暴露`Pod`的`metadata.labels`信息给容器.

通过这样的声明方式, 当前`Pod`的`Labels`字段的值就会被`Kubernetes`自动挂载称为容器里的`/tc/podinfo/labels`文件.

`Downward API`支持的字段:

```

1. 使用fieldRef可以声明使用:
spec.nodeName - 宿主机名字
status.hostIP - 宿主机IP
metadata.name - Pod的名字
metadata.namespace - Pod的Namespace
status.podIP - Pod的IP
spec.serviceAccountName - Pod的Service Account的名字
metadata.uid - Pod的UID
metadata.labels['<KEY>'] - 指定<KEY>的Label值
metadata.annotations['<KEY>'] - 指定<KEY>的Annotation值
metadata.labels - Pod的所有Label
metadata.annotations - Pod的所有Annotation

2. 使用resourceFieldRef可以声明使用:
容器的CPU limit
容器的CPU request
容器的memory limit
容器的memory request
```

`Downward API`能够获取到的信息, **一定是 Pod 里的容器进程启动之前就能够确定下来的信息**, 如果想要获取`Pod`容器运行后才会出现的信息, 不如容器进程的`PID`, 可以通过`sidecar`容器方式(`shareProcessNamespace=true`).

#### 环境变量 

处理`Secret, ConfigMap, Downward API`这三种定义的信息, 还可以通过环境变量的福安上出现在容器里. 但是, 通过环境变量获取这些信息的方式, 不具备自动更新的能力.

##### Service Account Token

加入在一个`pod`里安装一个`Kubernetes`的`Client`, 这样是否可以直接操作这个`Kubernetes`的`API`了呢.

需要先解决**API Server**的授权问题.

`Service Account`, 是`Kubernetes`系统内置的一种"服务账户", 它是`Kubernetes`进行权限分配的对象. 比如`Service Account A`, 可以只被允许对`Kubernetes API`进行`GET`操作, 而`Service Account B`, 则可以用`Kubernetes API`的所有操作的权限.

`Service Account`的授权信息和文件, 保存在一个特殊的`Secret`对象里**ServiceAccountToken**. 任何运行在`Kubernetes`集群上的应用, 都必须使用这个`ServiceAccountToken`里保存的授权信息, 也就是`token`, 才可以合法地访问`API Server`.

**为了方便使用, Kubernetes 已经提供了一个默认"服务账户"(default Service Account)**, 任何一个运行在 Kubernetes 里的 Pod, 都可以直接使用这个默认的 Service Account, 而无需显示地声明挂载它.

每一个`Pod`, 都已经自动声明一个类型是`Secret`, 名为`default-token-xxx`的`Volume`, 然后 **自动挂载每个容器的一个固定目录上**.

```shell
kubectl describe pod zhongjie-frontend-oa-74b7777ddc-tgscm -n=hzlh
...
Volumes:
  shared-data:
    Type:    EmptyDir (a temporary directory that shares a pod's lifetime)
    Medium:
  default-token-5xg84:
    Type:        Secret (a volume populated by a Secret)
    SecretName:  default-token-5xg84
    Optional:    false
QoS Class:       BestEffort
Node-Selectors:  type=frontend
Tolerations:     node.kubernetes.io/not-ready:NoExecute for 300s
                 node.kubernetes.io/unreachable:NoExecute for 300s
Events:          <none>
``` 

这个过程对于用户是**完全透明的**, 这样一旦`Pod`创建完成, 容器里的应用就可以直接从这个默认`ServiceAccountToken`的挂载目录访问到授权信息和文件. 这个容器的骷颅在`/var/run/secrets/kubernetes.io/serviceaccount`, 这个`Secret`类型的`Volume`里面的内容如下所示:

```shell
$ ls /var/run/secrets/kubernetes.io/serviceaccount 
ca.crt namespace  token
```

程序只要加载这些授权文件就可以访问并操作`Kubernetes API`了, 如果是官方的`Clinet`包, 可以自动加载这个目录下的文件.

**这种把 Kubernetes 客户端以容器的方式运行在集群里, 然后使用 default Service Account 自动授权的方式, 被称作"InClusterConfig"**.

可以设置默认不为`Pod`里的容器自动挂载这个`Volume`.

#### 容器健康检查和恢复机制

可以为`Pod`里的容器定义个健康检查"探针"(Probe). 这样, kubelet 就会根据这个`Probe`的返回值决定这个容器的状态, 而不是直接以容器进行是否运行(从`Docker`)作为依据.

```shell

apiVersion: v1
kind: Pod
metadata:
  labels:
    test: liveness
  name: test-liveness-exec
spec:
  containers:
  - name: liveness
    image: busybox
    args:
    - /bin/sh
    - -c
    - touch /tmp/healthy; sleep 30; rm -rf /tmp/healthy; sleep 600
    livenessProbe:
      exec:
        command:
        - cat
        - /tmp/healthy
      initialDelaySeconds: 5
      periodSeconds: 5
```

`livenessProbe`(健康检查). 它的类型是`exec`. 意味着, 容器启动后, 在容器里面执行依据我们指定的命令. 返回`0`就是健康的. 

* 容器启动后 5s 后开始执行
* 没 5s 执行一次

因为 30秒后文件会被删除, Kubernetes里的**Pod 恢复机制(restartPolicy, 默认值是 always)**会`Restart`这个容器(实际是重新创建了容器).

Pod 的恢复过程, 永远都是发生在当前节点上, 除非这个绑定发生了变化(pod.spec.node 字段被修改), 而不会跑到别的节点上去. 除非使用`Deployment`这样的"控制器"来管理 Pod.

restartPolicy:

* Always：在任何情况下, 只要容器不在运行状态, 就自动重启容器；
* OnFailure: 只在容器 异常时才自动重启容器；
* Never: 从来不重启容器.

如果关心容器退出后的上下文环境, 比如容器退出后的日志, 文件和目录, 就需要将`restartPolicy`设置为`Never`. 因为一旦容器被自动重新创建, 这些内容就有可能丢失掉.

**两个基本设计原理**

1. 只要 Pod 的 restartPolicy 指定的策略允许重启异常的容器(比如：Always), 那么这个 Pod 就会保持 Running 状态, 并进行容器重启. 否则, Pod 就会进入 Failed 状态.
2. 对于包含多个容器的 Pod, **只有它里面所有的容器都进入异常状态后, Pod 才会进入 Failed 状态**. 在此之前, Pod 都是 Running 状态. 此时, Pod 的 READY 字段会显示正常容器的个数.

假如一个 Pod 里只有一个容器，然后这个容器异常退出了. 那么, 只有当 `restartPolicy=Never`时, 这个 Pod 才会进入 Failed 状态. 而其他情况下, 由于 Kubernetes 都可以重启这个容器, 所以 Pod 的状态保持 Running 不变.

而如果这个 Pod 有多个容器, 仅有一个容器异常退出, 它就始终保持 Running 状态, 哪怕即使 restartPolicy=Never. 只有当所有容器也异常退出之后, 这个 Pod 才会进入 Failed 状态.

`readinessProbe`检查结果的成功与否, 决定的这个 Pod 是不是能被通过 Service 的方式访问到, 而并不影响 Pod 的声明周期.

#### PodPreset

运维人员可以定一个`PodPreset`对象. 这个对象可以对后续匹配的`Pod`里的字段进行追加.

**PodPreset 里定义的内容, 只会在 Pod API 对象被创建之前追加在这个对象本身上, 而不会影响任何 Pod 的控制器的定义.**

`Deployment`对象本身不会被`PodPreset`改变, 被修改的只是这个`Deployment`创建出来的所有`Pod`.

如果`PodPreset`有多个, `Kubernetes`会进行合并, 如果修改有冲突, 这些冲突的字段不会被修改.

## 扩展