# 11 | 从0到1：搭建一个完整的Kubernetes集群

## 笔记

### 准备工作

### 安装 kubeadm 和 Docker

### 部署 Kubernetes 的 Master 节点

```
horizontal-pod-autoscaler-use-rest-clients: "true"
```

将来部署的 `kube-controller-manager` 能够使用自定义资源（Custom Metrics）进行自动水平扩展.

第一次使用 Kubernetes 集群所需要的配置命令.

```
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

`Kubernetes`集群默认需要加密方式访问. 这几条命令, 就是将刚刚部署生成的`Kubernetes`集群的安全配置文件, 保存到当前用户的`.kube`目录下, `kubectl`默认会使用这个目录下的授权信息访问`Kubernetes`集群.

`kube-system`是`Kubernetes`项目预留的系统`Pod`的工作空间.

### 部署网络插件

`Kubernetes`支持容器网络插件, 使用的是一个名叫`CNI`的通用接口, 它也是当前容器网络的事实标准.

### 部署 Kubernetes 的 Worker 节点

`Kubernetes`的`Worker`节点跟`Master`节点几乎是相同的, 运行着的都是一个`kubelet`组件. 区别在于, `kubelet`启动后, `Master`节点上还会自动运行`kube-apiserver、kube-scheduler、kube-controller-manger`这三个系统`Pod`.

### 通过 Taint/Toleration 调整 Master 执行 Pod 的策略

默认情况下`Master`节点是不允许运行用户`Pod`的. 而`Kubernetes`做到这一点, 依靠的是`Kubernetes`的`Taint/Toleration`机制.

**原理**

一旦某个节点被加上了一个`Taint`，即被“打上了污点”, 那么所有`Pod`就都不能在这个节点上运行, 因为`Kubernetes`的`Pod`都有“洁癖”. 除非, 有个别的`Pod`声明自己能“容忍”这个“污点”, 即声明了`Toleration`, 它才可以在这个节点上运行.

### 部署 Dashboard 可视化插件

### 部署容器存储插件

很多时候我们需要用数据卷（Volume）把外面宿主机上的目录或者文件挂载进容器的`Mount Namespace`中, 从而达到容器和宿主机共享这些目录或者文件的目的. 容器里的应用, 也就可以在这些数据卷中新建和写入文件.

如果你在某一台机器上启动的一个容器，显然无法看到其他机器上的容器在它们的数据卷里写入的文件. **这是容器最典型的特征之一: 无状态**.

容器的持久化存储, 就是用来保存容器存储状态的重要手段.

**持久化**, 存储插件会在容器里挂载一个基于网络或者其他机制的远程数据卷，使得在容器里创建的文件，实际上是保存在远程存储服务器上，或者以分布式的方式保存在多个节点上，而与当前宿主机没有任何绑定关系。这样，无论你在其他哪个宿主机上启动新的容器，都可以请求挂载指定的持久化存储卷，从而访问到数据卷里保存的内容.



## 扩展