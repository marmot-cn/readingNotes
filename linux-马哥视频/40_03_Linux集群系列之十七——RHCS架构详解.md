# 40_03_Linux集群系列之十七——RHCS架构详解

---

## 笔记

---

### rhcs 高可用集群

集群套件:

* 负载均衡功能: LVS
* 高可用: HA
* 集群文件系统: GFS
* 集群逻辑卷: cLVM

```
资源代理(RA)
 |
资源管理器(CRM), 实现资源配置, 资源删除, 资源添加, 资源属性清理, 资源迁移...还能够完成基于资源代理(RA)在相应节点上启动, 停止或者重新配置资源.
 |
集群基础架构层(信息层, Cluster Infrastructure), 各节点传递心跳信息, 完成相关配置, 检查法定票数
  |
集群节点

```

RA:

* Lsb(/etc/init.d/)
* ocf(开放集群框架)

CRM:

* Heartbeat v1: HA resources
* Heartbeat v2: CRM
* Heartbeat v3: pacemaker

Cluster Infrastructure:

* heartbeat
* corosync
* keepalived

红帽早期(红帽4)使用的是`cman`(Cluster Manager), 独立组件.

RCHS: rgmanager 资源组管理器

* ra
	* internal
	* script: lsb
		* /etc/rc.d/init.d/*

dlm: Distributed Lock Manager 分布式锁管理器. 运行在各节点上的守护进程. 锁管理器彼此之间需要通信, 基于`TCP/IP`. 

gfs(Global File System): 要几个节点能挂载, 就要创建几个日志区域(日志文件).

ocfs2: 集群文件系统, Oracle Cluster File System.

Google File System. google 分布式文件系统. 简写也是`gfs`.

### cLVM

Cluster LVM.

借助于`HA`的功能, 将某节点对`lvm`的操作通知给其他节点. 把`lvm`做成分布式机制.

`/etc/lvm/lvm.conf`

`locking_type = 1`, 默认为`1`. 表示基于本机使用.

* `3` 内建的集群锁.

设置为`3`后, 并且借助于`cLVM`即可.

### CCS

Cluster Configuration System. 集群配置系统

`cman`中的`ccsd`(在每一个高可用节点运行守护进程), 集群配置文件管理. 在任何一个节点修改完成后, 都会通过底层信息同步层, 同步到其他节点.

### CMAN

红帽4, `CMAN`是工作在内核空间中.

红帽5, `CMAN`是工作在用户空间中.

* `dlm_controld`: 分布式锁控制器.
* `lock_dlmd`: 锁管理器.

### Failover Domain

**服务**故障转移域: 当一个节点发生故障后, 运行在A节点发生故障, 这些节点上的服务和资源所能转移到的目标节点.

故障转移域和服务(`Service`)相关. 

### PXE, COBBLER

早期一般使用人工配置`pxe`+`dhcp`+`tftp`配合`kickstart`.

现在可以使用开源工具, 如`cobbler`.

### 配置文件管理, 软件分发

`puppet`

### n-m 模型

`N`个节点运行`M`个模型(M<=N).

## 整理知识点

---