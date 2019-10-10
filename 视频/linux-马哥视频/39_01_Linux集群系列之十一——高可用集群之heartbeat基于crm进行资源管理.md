# 39_01_Linux集群系列之十一——高可用集群之heartbeat基于crm进行资源管理

---

## 笔记

---

### CRM

为那些非`ha-aware`, 本身不具备高可用管理功能(机制)的应用程序提供高可用的基础平台的. 通用性平台.

由`CRM`代为管理, 代为监控状态.

`CRM`需要提供基础管理框架.

#### crmd 管理API

* GUI
* CLI

### Resource Type

38-3笔记里面提及

* primitive(natived): 基本资源
* group:
* clone:
	* STONITH
	* Cluster Filesystem
		* dlm: 分布式锁管理器 Distributed Lock Manager
* master/slave: 主从类型资源
	* drbd 分布式复制块设备

### heartbeat v1配置文件

* authkeys: 集群事务信息在内部通知之间进行加密, 防止别的节点冒充成员.
* ha.cf(核心配置文件)
	* node指令定义各节点
	* 定义两者之间互相传递心跳信息的机制(集群事务信息的传递机制)
		* 广播 bcast(broadcast)
		* 多播 mcast(multicast)
		* 单播 ucast(unicast)
	* haresource: 资源配置信息, 自带资源管理器读取该配置文件

### 配置高可用需要准备的条件

1. 时间同步.
2. SSH双机互信.
3. 主机名称要与`uname -n`保持一致, 并且要通过`/etc/hosts`解析. 不能依赖于`DNS`解析.

### 示例, 配置组播地址

在`225.0.100.19`地址进行组播.

```
mcast eth0 225.0.100.19 694 1 0
```

### CIB

`C`luster `I`nformation `B`ase.

集群信息库.`XML`格式.

### 资源粘性

资源是否倾向于留在当前节点

* 正数: 乐意
* 负数: 离开

最终还要结合位置约束来结合在一起判断资源是否离开.

### send_arp 脚本

现在有两个节点. A 和 B

VIP 可能某一时间在A节点, 在与VIP通信之前需要解析到A节点的mac地址, 万一A节点故障, VIP被B节点抢走. 但是路由器缓存的还是A节点mac地址, 要通过**ARP欺骗的方式**进行更新, B节点自问自答.

使用脚本`send_arp`, 来完成.


## 整理知识点

---