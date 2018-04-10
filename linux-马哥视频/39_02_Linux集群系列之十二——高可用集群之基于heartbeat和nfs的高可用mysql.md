# 39_02_Linux集群系列之十二——高可用集群之基于heartbeat和nfs的高可用mysql

---

## 笔记

---

### hb v2, CRM来实现`MySQL`高可用集群

NFS: Mysql 数据

### nfs

`mysql`用户和`mysql`组要对文件有读写权限. 

### drbd

分布式磁盘块设备, 主机级别的`raid`, 在内核层通过网络同步.

主节点可以挂载, 读写.

备节点**不能**挂载. 所以也不能读写.

支持双主模型`Dual master`. 但是必须使用集群文件系统. `ocfs2, gfs`

## 整理知识点

---