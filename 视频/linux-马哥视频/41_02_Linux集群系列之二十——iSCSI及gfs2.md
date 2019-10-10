# 41_02_Linux集群系列之二十——iSCSI及gfs2

---

## 笔记

### chap认证

双向认证, 客户端和服务端都需要认证.

### iscsi-initiator-utils

不支持`discovery`认证.

如果使用基于用户的认证, 必须首先开放基于`IP`的认证.

### 集群文件系统

`mkfs.gfs2` 

* `-j` 指定日志区域个数, 有几个就能够被几个节点所挂载.
* `-J` 指定日志区域大小(单位是#MB), 默认是`128MB`
* `-p {lock_dlm|lock_nolck }` 锁协议的名称.
	* `lock_dlm` 分布式文件锁(多个节点挂载必须使用锁)
	* `lock_nolck` 不使用锁(如果只被一个节点挂载可以不使用锁) 
* `-t name` 锁表的名称. 唯一的表示某一个文件系统锁的持有情况.
	* 格式为`clustername:locktablename`.
		* `clustername` 为当前节点所在的集群的民此功能.
		* `locktablename` 必须要在当前集群内部唯一.

为了性能考虑, 一个集群文件日志区域个数最好不要超过`16`个.

## 整理知识点

---