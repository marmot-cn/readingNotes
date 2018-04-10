# mysql

---

## mysql 最低需求

(根据阿里云,官方手册没有找见)发现最低需求**1核1G**也可以. 考虑到实际使用暂时先给**2核4G**.

## mysql mgr 额外需求

* InnoDB Storage Engine.
* Primary Keys. 每张表都必须有一个主键.
* IPv4 Network.
* 网络, 因为`mgr`被设计为一个集群环境, 所以对彼此节点最好相邻, 且带宽延迟低.

