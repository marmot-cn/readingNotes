# 41_04_Linux集群系列之二十二——keepalived详解

---

## 笔记

### keepalived

提供`HA`的一个底层工具.

最早设计为`ipvs`提供`HA`功能. `ipvs`是在内核中提供, `keepalived`添加了能够将`vip`在节点之间流转的功能. `vip`流转基于`vrrp`协议.

#### vrrp

`vrrp`协议, 将多个物理设备虚拟为一个物理设备.



## 整理知识点

---