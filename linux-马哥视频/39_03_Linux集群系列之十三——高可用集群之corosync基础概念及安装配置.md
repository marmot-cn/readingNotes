# 39_03_Linux集群系列之十三——高可用集群之corosync基础概念及安装配置

---

## 笔记

---

### High Availability

```
A = MTBF / (MTBF + MTTR)
```

* `MTBF` = mean time between failures. 平均无故障时间.
* `MTTR` = mean time to repair. 平均修复时间.
* `A` = probability system will provide service at a random time(ranging from 0 to 1)

提高可用性:

* 提高`MTBF`
* 降低`MTTR`

硬件故障引起原因:

* Design failure (设计故障)
* Random failure (随机故障)
* Infant Mortality (用法不正确, 导致服务器寿命提前终止)
* Wear out (用坏了)

### Heartbeart 和 corosync
 
非常早的实现高可用的`messaging layer`

RHEL 6.x RHCS: corosync

RHEL 5.x RHCS: openais, cman, rgmanager 

`corosync`是一个实现了`Messaging layer`层的集群引擎. 用于完成集群事务信息传递.

#### openAIS

应用接口规范(AIS)是用来定义应用程序接口(API)的开放性规范的集合, 这些应用程序作为中间件为应用服务提供一种开放, 高移植性的程序接口.

OpenAIS提供一种集群模式, 这个模式包括集群框架, 集群成员管理, 通信方式, 集群检测等.

是一个接口, 一大堆API.

`openais`的分支`Wilson`. `Wilson `把`Openais`核心架构组件独立出来放在`Corosync`(`Corosync`是一个集群管理引擎).

`Wilson`是`openais`的一个版本的名称.

#### corosync --> pacemaker

* SUSE Linux Enterprise Server: Hawk, WebGUI
* LCMC: Linux Cluster Management Console
* RHCS: Conga(luci/ricci)
	* 浏览器 -> luci -> ssh不基于密码,基于密钥 -> ricci

* `corosync`支持节点多, 因为是基于组播通信.
* `keepalived`支持2个节点, 基于`VRRP`协议.

#### corosync 安装前提

* 时间同步
* 主机名
* SSH 互信

## 整理知识点

---