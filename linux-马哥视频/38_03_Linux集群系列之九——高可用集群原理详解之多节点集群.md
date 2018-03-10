# 38_03_Linux集群系列之九——高可用集群原理详解之多节点集群

---

## 笔记

---

### watchdog

在某一个节点上用来协调两个进程, 每一个进程在启动后要通过一个unix套接字(管道), 不停的像watchdog写数据, 要是其他进程发现对应的进程没有在写入数据了, 就认为进程挂掉了, 会试图重启进程.

### quorum 法定票数

without_quorum_policy

* freeze
* stop
* ignore

## 整理知识点

---