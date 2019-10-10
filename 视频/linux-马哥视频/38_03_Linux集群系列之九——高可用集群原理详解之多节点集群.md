# 38_03_Linux集群系列之九——高可用集群原理详解之多节点集群

---

## 笔记

---

### watchdog

在某一个节点上用来协调两个进程, 每一个进程在启动后要通过一个unix套接字(管道), 不停的像watchdog写数据, 要是其他进程发现对应的进程没有在写入数据了, 就认为进程挂掉了, 会试图重启进程.

### quorum 法定票数

without_quorum_policy

当一些节点认为自己不再是集群中的节点时候(不具备法定票数)采取的动作:

* freeze: 冻结, 不在接受新的请求, 当时当前链接进来的请求继续提供服务.
* stop: 停机.
* ignore: 忽略.

高可用集群必须要使用隔离设备.

### 故障转移域 failover domain

当一个节点`fail`时, 只能转移的节点范围叫做故障转移域.

### 节点运行模型

一共`N`个节点, 运行`m`个服务. **N-m**

`N`个节点运行`N`个服务. **N-n**

### RHCS

红帽集群套件.

RedHat Cluster Suite 红帽集群套件.

### Messaging Layer

让多个节点组合成高可用集群服务.

* heartbeat
	* v1
	* v2
	* v3
		* heartbear
		* pacemaker 集群资源管理器
		* cluster-glue
* corosync(需要和`pacemaker`组合起来使用)
* cman: cluster-manager (红帽5.0核心的 Messaging Layer)
* keepalived

### CRM

CRM 是附着在 Messaging Layer 之上运行的.

* heartbeat v1 (同时提供 Messaging Layer 也提供 CRM)
	* haresources
* heartbeat v2
	* 兼容 v1 的 haresources
	* crm
* heartbeat v3: 资源管理器crm发展为独立的项目: pacemaker(也可以使用corosync作为Messaging Layer)
* 以`CMAN`所提供的`rgmanager`(资源组管理器)

### 资源类型

* Primitive: 主资源, 在某一时刻只能运行在某一节点资源.
* clone: 把主资源`clone`N份, 同时运行在多个节点上的资源.
* group: 把多个`Primitive`资源归集在一个组里.
* master/slave: 独特的`clone`类资源, 只能运行在两个节点上(两个节点具有主从关系). 一个主节点, 一个从节点. 

### RA Resource Agent

是脚本. 接收`LRM`传递过来的控制指令, 控制资源. 主要提供资源管理功能.

`LRM`: 本地资源管理器.

#### RA Classes 的类别

* Legacy heartbeat v1 RA(专用于heartbeat v1)
* LSB(遵循Linux Shell 编程风格的, 位于`/etc/rc.d/init.d`下的脚本都是高可用资源管理脚本)
* OCF(Open Cluster Framework 开放式集群框架)
* STONITH(专门用于管理硬件STONITM设备的)

### STONITH 设备

#### Power Distribution Units 

电源分布单元(PDU), 电源交换机. 可以接上网线, 可以接收由其他节点所发来的控指令, 可以暂停某个接口的资源.

#### Uninterruptible Power Supplies(UPS)

不间断电源.

#### Blade Power Control Devices

刀片服务器的电源控制设备.

#### Lights-out Devices

轻型管理设备. 在服务器上的小管理模块, 可以实现远程管理.

#### Testing Devices

测试性的设备, `ssh`基于密钥认证. 通过`ssh`来关闭. 但是网络本来就不通了, 所以`ssh`可能也通不了.

### STONITH的实现

#### stonithd

为了实现`STONITH`功能, 在每一个节点都要启动一个进程`stonithd`, 用于监控当前主机所运行的`stonith`资源, 让多个节点的`STONITH`进程互相通信的一种机制.

#### STONITH Plug-ins

所有支持`fencing`的设备, 都通过`STONITH Plug-in`来管理设备.

`STONITH Plug-in`是实现`fencing`设备的接口.

### 资源隔离级别

* 节点级别
	* STONITH
* 资源级别
	* FC SAN Switch

## 整理知识点

---

高可用服务每个资源都不能开机自动启动, 需要由`CRM`来决定在哪个节点启动.