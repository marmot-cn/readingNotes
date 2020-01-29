# 10 | Kubernetes一键部署利器：kubeadm

## 笔记

### Kubernetes 简单的部署方法

**kubeadm**

```
# 创建一个Master节点
$ kubeadm init

# 将一个Node节点加入到当前集群中
$ kubeadm join <Master节点的IP和端口>
```

### kubeadm 的工作原理

`Kubernetes`的架构和它的组件. 在部署时, 它的每一个组件都是一个需要被执行的, 单独的二进制文件. 

**为什么不用容器部署Kunernetes呢?**

如何容器化`kubelet`.

`kubelet`是`Kubernetes`项目用来操作`Docker`等容器运行时的核心组件. 除了跟容器运行时打交道, `kubelet`在配置容器网络, 管理容器数据卷时, 都需要直接操作宿主机.

**容器部署kubernetes的问题**

* 配置网络可以通过不开启`Network Namespace`(Docker 的`host network`模式)的方式, 直接共享宿主机的网络栈. **可以解决**
* 配置挂载文件系统, 除非使用`setns()`系统调用在宿主机的`Mount Namespace`中执行这些挂载操作. **不好解决**

`kubeadm`选择了一种妥协的方案.

**`kubelet`直接运行在宿主机上, 然后使用容器部署其他的`Kubernetes`组件**.

### kubeadm init 工作流程

#### 1.Preflight Checks

首先要做一系列检查工作, 以确定这台机器可以用来部署`Kubernetes`, 称之为`Preflight Checks`.

* Linux 内核的版本3.10以上
* Linux Cgroups 模块是否可用
* hostname 是否标准, 在 Kubernetes 项目里, 机器的名字以及一切存储在 Etcd 中的 API 对象, 都必须使用标准的 DNS 命名
* 安装的`kubeadm`和`kubelet`的版本是否匹配
* 是不是已经安装了`Kubernetes`的二进制文件
* `Kubernetes`的工作端口`10250/10251/10252`端口是不是已经被占用
* `ip, mount`等`Linux`指令是否存在
* `Docker`是否已经安装

#### 2. 生成证书

检查之后, `kubeadm`生成`Kubernetes`对外提供服务所需的各种证书和对应目录.

`Kubernetes`对外提供服务时, 除非专门开启**不安全模式**, 否则都要通过`HTTPS`才能访问`kube-apiserver`. 需要为`Kubernetes`集群配置好证书文件.

证书文件存储在`Master`节点的`/etc/kubernetes/pki`目录下. 在这个目录下最主要的证书文件是`ca.crt`和对应的私钥`ca,key`.

使用`kubelet`获取容器日志等`streaming`操作时, 需要通过`kube-apiserver`向`kubelet`发起请求, 这个连接也必须是安全的. `kubeadm`为这一步生成的是:

* `apiserver-kubelet-client.crt`
* `apiserver-kubelet-key`

#### 3. 生成配置文件

`kubeadm`为其他组件生成访问`kube-apiserver`所需的配置文件. 文件路径是`/etc/kubernetes/xxx.conf`

```shell
ls /etc/kubernetes/
admin.conf  controller-manager.conf  kubelet.conf  scheduler.conf
```

这些文件记录的是, 当前`Master`节点的服务器地址, 监听端口, 证书目录等信息. 对应的客户端(`scheduler, kubelet`等), 可以直接加载相应的文件, 使用里面的信息与`kube-apiserver`建立安全连接.

#### 4. 为`Master`组件生成`Pod`配置文件

`Master`的组件`kube-apiserver, kube-controller-manager, kube-scheduler`, 它们都会被使用`Pod`的方式部署起来.

`Kubeadm`还会再生成一个`Etcd`的`Pod YAML`文件, 通过同样的`Static Pod`的方式启动`Etcd`.

在`Kubernetes`中, 有一种特殊的容器启动方法叫做"`Static Pod`". 它允许你把要部署的`Pod`的`YAML`文件放在一个指定的目录里. 这样, 当这台其上的`kubelet`启动时. 它会自动检查这个目录, 加载所有的`Pod YAML`文件, 然后在这台机器上启动它们.

`Master`组件的 YAML 文件会被生成在`/etc/kubernetes/manifests`路径下.

`Matser`组件的`Pod YAML`文件如下所示:

```
$ ls /etc/kubernetes/manifests/
etcd.yaml  kube-apiserver.yaml  kube-controller-manager.yaml  kube-scheduler.yaml
```

#### 5. 健康检查

当`Master`容器启动后, `kubeadm`会通过检查`localhost:6443/healthz`这个`Master`组件的健康检查`URL`, 等待`Master`组件完全运行起来.

#### 6. 生成 token

`kubeadm`会为集群生成一个`bootstrap token`. 在后面, 只要持有这个`token`，任何一个安装了`kuebelet`和`kubeadm`的节点, 都可以通过`kubeadm join`加入到这个集群当中.

在`token`生成之后, `kubeadm`会将`ca.crt`等`Master`节点的重要信息,通过`ConfigMap`(cluster-info)的方式保存在`Etcd`当中, 供后续部署`Node`节点使用.

#### 7. 安装插件

`kube-proxy`和`DNS`这两个插件是必须安装的. 提供整个集群的服务发现和`DNS`功能.

### kubeadm join 的工作流程

任何一台机器想要成为`Kubernetes`集群中的一个节点, 就必须在集群的`kube-apiserver`上注册. 可是, 要想跟`apiserver`打交道, 这台机器就必须要获取相应的证书文件(CA文件).

`kubeadm`至少需要发起一次**不安全模式的**的访问到`kube-apiserver`, 从而拿到保存在`ConfigMap`中的`cluster-info`(它保存了APIServer的授权信息). `bootstrap token`, 在这个过程中做安全验证.

有了`cluster-info`里的`kube-apiserver`的地址, 端口, 证书, `kubelet`就可以**安全模式**连接到`apiserver`上, 这样一个新的节点就部署完成了.

### 配置 kubeadm 的部署参数

```
$ kubeadm init --config kubeadm.yaml
```

通过部署参数配置文件启动.

## 扩展