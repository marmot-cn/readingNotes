#Kubernetes

---

###Kubernetes核心概念

####Pods

在Kubernetes系统中,调度的最小颗粒不是单纯的容器,而是抽象成一个`Pod`,`Pod`是一个可以被创建,销毁,调度,管理的最小的`部署单元`.

一个`pod`是由若干个`Docker`容器构成的容器组(pod意味豆荚,里面容纳了多个豆子).

`pod里的容器是共享网络和存储的`.


####Replication Controllers

Replication Controller实现复制多个Pod副本,往往一个应用需要多个Pod来支撑,并且可以保证其复制的副本数,即使副本所调度分配的宿主机出现异常,通过Replication Controller可以保证在其它主宿机启用同等数量的Pod.

**设计理念**

为每个pod"`外挂`"一个控制器进程,从而避免了健康检查组建称为性能瓶颈;即时这个控制器进程失效,容器依然可以正常运行,pod和容器无需知道这个控制器,也不会把这个控制器作为依赖.

####Services

Services是Kubernetes最外围的单元,通过虚拟一个访问IP及服务端口,可以访问我们定义好的Pod资源,目前的版本是通过`iptables`的`nat`转发来实现,转发的目标端口为`Kube_proxy`生成的随机端口.

对于每个`service`,`kube-proxy`都会在宿主机上随机监听一个端口与这个`service`对应起来,它会在宿主机上建立起`iptables`规则。

####Labels

Labels是用于区分`Pod`、`Service`、`Replication Controller`的`key/value`键值对,仅使用在`Pod`、`Service`、`Replication Controller`之间的`关系识别`,但对这些单元本身进行操作时得使用`name`标签.

####Proxy

Proxy不但解决了同一主宿机相同服务端口冲突的问题,还提供了Service转发服务端口对外提供服务的能力,Proxy后端使用了随机,轮循负载均衡算法.

###master运行组件

####kube-apiserver

作为kubernetes系统的入口,封装了核心对象的增删改查操作,以RESTFul接口方式提供给外部客户和内部组件调用.它维护的REST对象将持久化到`etcd`.

####kube-scheduler

负责集群的资源调度,为新建的pod分配机器.这部分工作分出来变成一个组件,意味着可以很方便地替换成其他的调度器.

####kube-controller-manager

负责执行各种控制器,目前有两类:

* `endpoint-controller`
* `replication-controller`

**`endpoint-controller`**

定期关联service和pod(关联信息由endpoint对象维护),保证service到pod的映射总是最新的.

**`replication-controller`**

定期关联replicationController和pod,保证replicationController定义的复制数量与实际运行pod的数量总是一致的.

###minion运行组件

####kube-proxy

负责为`pod`提供代理.它会定期从`etcd`获取所有的`service`,并根据service信息创建代理.当某个客户`pod`要访问其他`pod`时,访问请求会经过本机`proxy`做转发.

####kubelet

负责管控`docker`容器,如启动/停止,监控运行状态等.它会定期从`etcd`获取分配到本机的pod,并根据pod信息启动或停止相应的容器.同时,它也会接收apiserver的HTTP请求，汇报pod的运行状态.
