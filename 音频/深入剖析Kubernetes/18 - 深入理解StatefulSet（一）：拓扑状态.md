# 18 | 深入理解StatefulSet（一）：拓扑状态

## 笔记

**Deployment**做了一个假设:

一个应用的所有`Pod`, 是完全一样的. 所以, 它们互相之间**没有顺序**, 也**无所谓运行在哪台宿主机上**. 需要的时候, 通过`Pod`模板创建新的, 不需要的时候, 可以"杀掉"任意一个.

**没有考虑到的情况**:

* 依赖关系: 注册, 主备
* 数据存储类应用, 有多个实例: 在本地磁盘上保存一份数据. 这些实例一旦被杀掉, 即便重建出来, **实例与数据库之间的对应关系也已经丢失**.

**这种示例之间有不对等关系, 以及示例对外部数据有依赖关系的应用, 被称为"有状态应用"(Statsful Application)**

容器, 比较易于用来封装"无状态应用"(Stateless Application), 如`Web`服务.

`Kubernetes`在`Deployment`的基础上, 扩展了对"有状态应用"的支持. 就是**StatefulSet**

`StatefulSet`的设计, 把真实世界里的应用状态, 抽象为两种情况:

* **拓扑状态**. 多个示例之间不是完全对等的关系, 这些应用示例, 必须按照某些顺序启动. 比如删掉后重新创建的`Pod`必须和原来`Pod`的网络标识一样, 这样原先的访问者才能使用同样的方法, 访问到这个新`Pod`.
* **存储状态**. 多个示例分别绑定了不同的存储数据. 

**StatefulSet的核心功能, 就是通过某种方式记录这些状态, 然后在 Pod 被重新创建时, 能够为新 Pod 恢复这些状态**

### Headless Service

`Service`是`Kubernetes`项目中用来将一组`Pod`暴露给外界访问的一种机制. 如, 一个`Deployment`有3个`Pod`, 那么就可以定义一个`Service`. 然后, 用户只要能访问到这个`Service`, 它就能访问到某个具体的`Pod`.

#### Service 如何被访问的

##### 1. `Service`的`VIP`(虚拟IP)

访问`10.0.23.1`这个`Service`的`IP`地址时, `10.0.23.1`其实就是一个`VIP`, 它会把请求转发到该`Service`所代理的某一个`Pod`上.

##### 2. `Service`的`DNS`

访问`my-svc.my-namespace.svc.cluster.local`这条`DNS`记录, 就可以访问到名叫`my-svc`的`Service`所代理的某一个`Pod`.

**DNS的两种处理方式**

1. `Normal Service`, 域名解析到`Service`的`VIP`
2. `HeadlessService`, 域名解析到`Service`代理的某一个`Pod`的`IP`地址.

**HeadlessService 不需要分配一个 VIP, 而是可以直接以 DNS 记录的方式解析出被代理 Pod 的 IP 地址**

`Headless Service`对应的`YAML`文件

```
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  ports:
  - port: 80
    name: web
  clusterIP: None
  selector:
    app: nginx
```

`clusterIP`字段的值是: `None`. 这个`Service`, 没有一个`VIP`作为"头". **这也就是 Headless 的含义**. 这个`Service`被创建后并不会被分配一个`VIP`, **而是会以 DNS 记录的方式暴露出它所代理的 Pod**.

这个`Headless Service`所代理的`Pod`, 通过`Label Serector`机制选择出来的, 即携带了`app=nginx`标签的`Pod`, 会被这个`Service`代理起来.

以这样的方式创建了一个Headless Service之后, **它所代理的所有`Pod`的`IP`地址, 都会被绑定一个这样格式的 DNS 记录**, 如下所示:

```
<pod-name>.<svc-name>.<namespace>.svc.cluster.local
```

这个`DNS`记录, 正式`Kubernetes`项目为`Pod`分配的唯一的"可解析身份"(`Resolvable Identity`).

只要你知道了一个`Pod`的名字, 以及它对应的`Service`的名字, 就可以非常确定地通过这条`DNS`记录访问到`Pod`的`IP`地址.

`StatefulSet`的`YAML`配置文件:

```
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  serviceName: "nginx"
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.9.1
        ports:
        - containerPort: 80
          name: web
```

多了一个`serviceName: "nginx"`字段. 这个字段的作用就是告诉`StatefulSet`控制器, 在执行控制循环(`Control Loop`)的时候, 请使用`nginx`这个`Headless Service`来保证`Pod`的"可解析身份".

`StatefulSet`给它所管理的所有`Pod`的名字进行了编号, 编号规则是`-`. 这些编号都是从`0`开始累加, 与`StatefulSet`的每个`Pod`实例一一对应, 绝不重复. 这些`Pod`的创建也是严格按照编号顺序进行的. 比如, 在 web-0 进入到 Running 状态、并且细分状态（Conditions）成为 Ready 之前, web-1 会一直处于 Pending 状态.

这两个容器的`hostname`也是与`Pod`名字是一致的, 都被分配了对应的编号.

```
$ kubectl exec web-0 -- sh -c 'hostname'
web-0
$ kubectl exec web-1 -- sh -c 'hostname'
web-1
```

可以通过`web-0.nginx`和`web-1.nginx`分配访问`Pod`.

如果删除后在重新创建, `Kubernetes`会按照原先编号的顺序, 创建出两个新的`Pod`, 分配了与原来相同的"网络身份".

**StatefulSet 保证了 Pod 网络标识的稳定性**, 如`web-0.nginx`是主节点, `web-1`是从节点. 这个关系绝对不会发生变化.

通过这种方法, **Kubernetes 就成功地将 Pod 的拓扑状态(比如: 哪个节点先启动, 那个节点后启动), 按照 Pod 的"名字 + 编号"的方式固定了下来**. `Kubernetes`还未每一个`Pod`提供了一个固定并且唯一的访问入口. 

**注意**, 尽管`web-0.nginx`这条记录本身不会变, 但是它解析到的`Pod`的`IP`地址, 并不是固定的. 这就意味着, 对于"有状态应用"实例的访问, 必须使用`DNS`记录或者`hostname`的方式, 而绝不应该直接访问这些`Pod`的`IP`地址.

## 扩展