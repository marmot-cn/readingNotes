# 38_01_Linux集群系列之七——高可用集群原理详解

---

## 笔记

---

### 资源粘性

资源更倾向于运行于哪个节点.

定义资源和节点之间的倾向性.

### 资源约束

Constraint.

定义资源和资源之间的倾向性.

* 排列约束: coloation, 资源间的互斥性, 定义资源是否能够在一起(能否运行在同一节点上).
	* score:
		* 正值: 可以在一起.
		* 负值: 不能再一起. 
* 位置约束: location, 资源对节点的依赖程度的, 定义资源运行在哪个节点, 给一个整数值(score).
	* 正值: 倾向于此.
	* 负值: 倾向于逃离此.
* 顺序约束: order, 定义资源和资源之间的启动或关闭次序(解决资源之间的依赖关系).

### 资源隔离

1. 节点级别: STONITH
2. 资源级别: 
	* FC SAN switch 可以实现在存储资源级别拒绝某节点的访问.
	* ... 

### 集群事务信息

一般集群事务信息通过`UDP`传递.

### CRM

Cluster Resource Manager

* DC
* LRM: 本地资源管理器. 接收由`TE`传递过来的指挥, 在某一个节点上采取相应动作.

#### DC

决策节点, 自动选出的.

Designated Coordinator. 

* PE: Policy Engine
* TE: Transaction Engine

### LSB

Linux Standard Base

* start
* stop
* restart
* statue

符合`LSB`标准库的脚本.

### 集群文件系统

保证任何一个节点写, 其他节点也能看见.

Cluster Filesystem:

* GFS
* OCFS2

## 整理知识点

---