# 09 | 从容器到容器云：谈谈Kubernetes的本质

## 笔记

### 容器

一个"容器", 是一个由:

* Linux Namespace
* Linux Cgroups
* rootfs 

三种技术构建出来的进程的隔离环境.

### 一个正在运行的容器

* 一组联合挂载在`/var/lib/docker/aufs/mnt`上的`rootfs`, **容器镜像**(Container Image). 容器的**静态**视图.
* 一个由 `Namespace+Cgroups` 构成的隔离环境. **容器运行时**(Container Runtime). 容器的**动态**视图.

### Kubernetes 项目要解决什么问题

![](./img/09_01.png)

* Master
	* API 服务, kube-apiserver
	* 调度, kube-scheduler
	* 容器编排, kube-controller-manager
	* 集群的持久化数据, 有 kube-apiserver 处理后保存在 Etcd 中
* 计算节点
	* kubelet
		* 负责通容器运行时(如 Docker)打交道, 依赖于 CRI(Container Runtime Interface)的远程调用接口, 定义了容器运行时的各项核心操作(如, 启动一个容器需要的参数).
			* 只要容器运行时能够运行标准的容器镜像, 就可以通过 `CRI` 介入到 k8s 中.
				* Docker 通过 `OCI` 这个容器运行时规范通底层的 Linux 操作系统进行交互. 把`CRI`请求翻译成对 Linux 操作系统的调用(操作 Linux Namespace 和 Cgroups)
		* 通过`gRPC`协议通`Device Plugin`进行交互, 管理 GPU 等宿主机物理设备的主要组件.
		* 调用网络插件`CNI(Container Networking Interface)`为容器配置网络.
		* 调用存储插件`CSI(Container Storage Interface)`为容器配置持久化存储.

### Borg 对于 k8s 的指导作用体现

k8s 没有把`Docker`作为整个结构的核心, 而**仅仅把它作为最底层的一个容器运行时表现**.

### 微服务落地的先决条件

容器的本质是**进程**.

那些原先拥挤在同一个虚拟机里的各个应用、组件、守护进程，都可以被分别做成镜像, 然后运行在一个个专属的容器中. 它们之间互不干涉, 拥有各自的资源配额, 可以被调度在整个集群里的任何一台机器上.

### k8s 设计思想

从更宏观的角度, 以统一的方式来定义任务之间的各种关系, 并且为将来支持更多种类的关系留有余地.

#### Pod

`Pod` 里的容器共享同一个`Network Namespace`, 同一组数据卷, 从而达到高效交换信息的目的.

#### Service

`Service` 服务声明的 IP 地址等信息是"终生不变". 这个`Service`服务的主要作用, 就是作为`Pod`的代理入口(Portal), 从而代替`Pod`对外暴露一个固定的网络地址.

`Service` 后端真正代理的`Pod`的`IP`地址, 端口等信息的自动更新, 是`k8s`的职责.

#### 核心功能全景图

![](./img/09_02.png)

#### Secret

`Secret`是一个保存在 Etcd 里的键值对数据. 把 Credential 信息以 Secret 的方式存在`Etcd`里, K8s 就会在你指定的 Pod 启动时, 自动把 Secret 里的数据以 Volume 的方式挂载到容器里.

#### Job

描述一次性运行的`Pod`

#### DaemonSet

每个宿主机必须且只能运行一个副本的守护进程服务

#### CronJob

描述定时任务

### 声明式API

* 通过一个编排对象, 比如 Pod、Job、CronJob 等，来描述你试图管理的应用.
* 再为它定义一些"服务对象", 比如 Service、Secret、Horizontal Pod Autoscaler(自动水平扩展器)等. 这些对象, 会负责具体的平台级功能。

## 扩展

### 调度

是把一个容器, 按照某种规则, 放置在某个最佳节点上运行起来.

### 编排

是按照用户的意愿和整个系统的规则, 完全自动化地处理好容器之间的各种关系.