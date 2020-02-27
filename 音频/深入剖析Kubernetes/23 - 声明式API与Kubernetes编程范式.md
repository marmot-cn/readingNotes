# 23 | 声明式API与Kubernetes编程范式

## 笔记

**命令式命令行操作**

```
$ docker service create --name nginx --replicas 2  nginx
$ docker service update --image nginx:1.7.9 nginx
```

**命令式配置文件操作**

```
kubectl create ...
kubectl replace ...
```

和"命令式命令行操作"没什么本质上的区别. 只不过, 把`Docker`命令行里的参数, 写在了配置文件里.

### 什么是"声明式API"

`kubectl apply`命令, 推荐代替`kubectl create`.

`kubectl replace`的执行过程, 是使用新的`YAML`文件中的`API`对象, 替换原有的`API`对象.

而`kubectl apply`, 则是执行了一个对原有`API`对象的`PATCH`操作.

`kubectl set image`和`kubectl edit`也是对**已有`API`对象的修改**.

#### 以`Istio`项目为例，来为你讲解一下声明式 API 在实际使用时的重要意义

**Istio**, 是一个基于`Kubernetes`项目的微服务治理框架.

![](./img/23_01.jpg)

`Istio`架构的核心所在. `lstio`最根本的组件, 是运行在每一个应用`Pod`里的`Envoy`容器. `Envoy`是一个高性能`C++`网络代理.

`Istio`项目, 把这个代理服务以`sidecar`容器的方式, 运行在了每一个被治理的应用`Pod`中. 因为`Pod`里的所有容器都共享同一个`Network Namespace`. 所以, **`Envoy`容器就能通过配置`Pod`里的`iptables`规则, 把整个`Pod`的进出流量接管下来**.

`Istio`的控制层(`Control Plane`)里的`Pilot`组件, 就能够通过调用每个`Envoy`容器里的`API`, 对这个`Envoy`代理进行配置, 从而实现微服务治理.

如`Pilot`可以通过调节流量控制把`90%`的流量给左边的`Pod`, `10%`的流量给右边的`Pod`. 从而实现了灰度发布. 针对用户和应用都是完全"无感"的.

`Istio`项目明明需要在每个`Pod`里安装一个`Envoy`容器, 怎么做到无感. 使用的是`Kubernetes`中的`Dynamic Admission Control`.

`Kubernetes`中, 当一个`Pod`后者任何一个`API`对象被提交给`APIServer`之后, **有一些"初始化"性质的工作需要在它们被`Kubernetes`项目正式处理之前进行**. 如, 自动为所有`Pod`加某些标签(Labels).

这个"初始化"操作的实现, 借助于"Admission"的功能. 是`Kubernetes`项目里一组被称为`Admission Controller`的代码, 可以选择性地被编译进`APIServer`中, 在`API`对象创建之后会被立刻调用到.

如果想自己添加规则到`Admission Controller`, 要求重新编译并重启`APIServer`. 

`Kuberntes`提供一种**"热插拔"**式的`Admission`机制, 就是`Dynamic Admission Control`, 也叫做: `Initializer`.

如, 下面例子

```
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app: myapp
spec:
  containers:
  - name: myapp-container
    image: busybox
    command: ['sh', '-c', 'echo Hello Kubernetes! && sleep 3600']
```

`Istio`需要在这个`Pod YAML`被提交给`Kubernetes`之后, 在它对应的`API`对象里自动加上`Envoy`容器的配置, 使之变成下面的样子:

```
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app: myapp
spec:
  containers:
  - name: myapp-container
    image: busybox
    command: ['sh', '-c', 'echo Hello Kubernetes! && sleep 3600']
  - name: envoy
    image: lyft/envoy:845747b88f102c0fd262ab234308e9e22f693a1
    command: ["/usr/local/bin/envoy"]
    ...
```

多了一个`envoy`的容器, 它就是`Istio`要使用的`Envoy`代理.

`Istio`, 编写一个用来为`Pod`"自动注入"`Envoy`容器的`Initializer`.

### Istio 如何实现 Envoy 的 Initializer

#### 1. Istio 会将这个 Envoy 容器本身的定义, 以`ConfigMap`的方式保存在`Kubernetes`当中.

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: envoy-initializer
data:
  config: |
    containers:
      - name: envoy
        image: lyft/envoy:845747db88f102c0fd262ab234308e9e22f693a1
        command: ["/usr/local/bin/envoy"]
        args:
          - "--concurrency 4"
          - "--config-path /etc/envoy/envoy.json"
          - "--mode serve"
        ports:
          - containerPort: 80
            protocol: TCP
        resources:
          limits:
            cpu: "1000m"
            memory: "512Mi"
          requests:
            cpu: "100m"
            memory: "64Mi"
        volumeMounts:
          - name: envoy-conf
            mountPath: /etc/envoy
    volumes:
      - name: envoy-conf
        configMap:
          name: envoy
```

这个`ConfigMap`的`data`部分, 正是一个`Pod`对象的一部分定义. 

`Initalizer`要做的工作, 就是把这部分`Envoy`相关的字段, 自动添加到用户提交的`Pod`的`API`对象里. 因为用户提交的`Pod`本身就有`containers`和`volumes`, 所以`Kubernetes`在处理这样的更新请求时, 就必须使用类似于`git merge`这样的操作, 才能将这两部分内容合并在一起.

`Inttializer`更新用户的`Pod`对象的时候, 必须使用`PATCH API`来完成. 正式声明`API`最主要的能力.

### 2. Istio 将一个编写好的 Initializer, 作为一个 Pod 部署在 Kubernetes 中

```
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: envoy-initializer
  name: envoy-initializer
spec:
  containers:
    - name: envoy-initializer
      image: envoy-initializer:0.0.1
      imagePullPolicy: Always
```

`envoy-initializer:0.0.1`镜像, 是一个实现编写好的**"自定义控制器"(`Custom Controller`)**.

`Kubernets`的控制器, 实际上就是一个"死循环". 它不断地获取"实际状态", 然后与"期望状态"做对比, 并以此为依据决定下一步的操作.

`Initializer`的控制器, 不断获取到的"实际状态", 就是用户新创建的`Pod`. 而它的"期望状态", 则是: 这个`Pod`里被添加了`Envoy`容器的定义.

控制器逻辑的伪代码:

```go
for {
  // 获取新创建的Pod
  pod := client.GetLatestPod()
  // Diff一下，检查是否已经初始化过
  if !isInitialized(pod) {
    // 没有？那就来初始化一下
    doSomething(pod)
  }
}
```

`Istio`要往`Pod`里合并的字段, 正是之前保存在`envoy-initializer`这个`ConfigMap`里的数据.

#### 3. `Initializer`控制器的工作逻辑

首先会从`APIServer`中拿到这个`ConfigMap`.

```
func doSomething(pod) {
  cm := client.Get(ConfigMap, "envoy-initializer")
}
```

然后, 把这个`ConfigMap`里存储的`containers`和`volumes`字段, 直接添加进一个空的`Pod`对象里.

```
func doSomething(pod) {
  cm := client.Get(ConfigMap, "envoy-initializer")
  
  newPod := Pod{}
  newPod.Spec.Containers = cm.Containers
  newPod.Spec.Volumes = cm.Volumes
}
```

`Kuberntes`的`API`库, 提供了一个直接使用新旧两个`Pod`对象, 生成一个`TwoWayMergePatch`:

```
func doSomething(pod) {
  cm := client.Get(ConfigMap, "envoy-initializer")

  newPod := Pod{}
  newPod.Spec.Containers = cm.Containers
  newPod.Spec.Volumes = cm.Volumes

  // 生成patch数据
  patchBytes := strategicpatch.CreateTwoWayMergePatch(pod, newPod)

  // 发起PATCH请求，修改这个pod对象
  client.Patch(pod.Name, patchBytes)
}
```

合并后, `Initializer`的代码就可以使用这个`ptach`的数据, 调用`Kubernets`的`Client`, 发起一个`PATCH`请求.

这样, 用户提交的`Pod`对象里, 就会被自动加上`Envoy`容器相关的字段.

`Kubernetes`可以指定什么样的资源进行这个`Initialize`操作.

```
apiVersion: admissionregistration.k8s.io/v1alpha1
kind: InitializerConfiguration
metadata:
  name: envoy-config
initializers:
  // 这个名字必须至少包括两个 "."
  - name: envoy.initializer.kubernetes.io
    rules:
      - apiGroups:
          - "" // 前面说过， ""就是core API Group的意思
        apiVersions:
          - v1
        resources:
          - pods
```

这个配置:

* 队友所有的`Pod`进行这个`Initialize`操作.
* 负责这个操作的`Initializer`, 名叫: `envoy-initializer`.

这个`InitializerConfigueration`被创建, `Kubernetes`会把这个`Initializer`的名字, **加在所有新创建的`Pod`的`Metadata`上**.

```
apiVersion: v1
kind: Pod
metadata:
  initializers:
    pending:
      - name: envoy.initializer.kubernetes.io
  name: myapp-pod
  labels:
    app: myapp
...
```

每一个新创建的`Pod`, 都会自动携带了`metadata.initializers.pending`的`Metadata`信息.



**这个`Metadata`, 正是接下来`Initializer`的控制器判断这个`Pod`有没有执行过自己所负责的初始化操作的重要依据**.

当自己开发的`Initializer`完成了要做的操作后, 一定要清除掉`metadata.initializers.pending`标志.


需要清除掉下面这个标签
```
metadata.initializers.pending[0].name=" envoy.initializer.kubernetes.io
```
否则这些Pod会一直出在 uninitialized 状态.

还可以通过具体的`Pod`的`Annotation`里添加一个字段, 从而声明要使用某个`Initializer`.

```
apiVersion: v1
kind: Pod
metadata
  annotations:
    "initializer.kubernetes.io/envoy": "true"
    ...
```

### 总结

声明式`API`

* 只需要提交一个定义好的`API`对象来"声明", 我所期望的状态是什么样子.
* 声明式`API`"允许有多个`API`写端, 以`PATCH`的方式对`API`对象进行修改, 而无需关心本地原始`YAML`文件的内容.

标准的"Kubernetes 编程范式" - **如何使用控制器模式, 同`Kubernetes`里`API`对象的"增、删、改、查"进行协作, 进而完成用户业务逻辑的编写过程**.

## 扩展