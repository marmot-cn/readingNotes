# 42_03_MySQL主从复制——MySQL-5.6基于GTID及多线程的复制

---

## 笔记

### 数据库复制过滤

* 主服务器
	* `binlog-do-db`: 白名单, 可以写多次. 
		* 仅将制定数据库的修改操作记入二进制日志, 但是主库恢复的时候会导致主段二进制日志不完整.
	* `binlog-ignore-db`: 黑名单.
	* 如果主服务器不使用`binlog-do-db`和`binlog-ignore-db`那么就会将所有的数据都记录到二进制日志当中.
* 从服务器
	* `replicate-do-db`: 从中继日志读取的白名单.
	* `replicate-ignore-db`: 从中继日志读取的黑名单.
	* `replicate-do-table`: 表级别
	* `replicate-ignore-table`: 表级别
	* `replicate-wild-do-table=magedudb.tb%(tb开头所有的表都复制)`: 通配符
	* `replicate-wild-ignore-table`: 通配符

`wild`支持`%`和`_`.

这些变量**只读**, 不允许在服务器运行时修改.

应该在从端设置, 保证铸锻二进制日志完成.

### GTID

由服务器的`UUID` + 事务`ID`号 = GTID.

`5.6`中每个事务的事件首部都会有`GTID`.

可以让主从两端自动发现去哪复制.

### 多线程复制

每个数据库仅能使用一个线程.

复制涉及到多个数据库是多线程复制才有意义.

从服务器`slave-parallel-workers=#`表示启用的线程数. 这个数字要小于等于库数.
`0`: 表示禁用多线程功能.

### mysql5.6 配置启用`gtid`参数

* `log_slave_updates = ON`: 当我们从服务器上, 从中继日志中读取事件在本地应用时, 是否把事件的写操作应用在二进制日志上.
* `gtid_mode=ON`: 是否启用`gitid`功能.
* `enforce_gtid_consistency=ON`: 是否强制`gitid`一致性.
* `report-port`: `gtid`模式下, 每个从服务器连入时候都需要告知端口号.
* `report-host`: `gtid`模式下, 每个从服务器连入时候都需要告知主机号.
* `sync-master-info=1`: 启用之可确保无信息丢失；任何一个事务提交后, 将二进制日志的文件名及事件位置记录到文件中 
* `master-info-repository=TABLE`: #定义master-info(主服务器) 记录在table中
	* 文件: `master.info`
	* TABLE: `mysql.salve_master_info`
* `relay_log_info_repository=TABLE`: #定义relay-log-info(从服务器) 记录在table中
* `binlog-checksum`: 主服务器端在启动时, 是否校验二进制日志校验码.
* `master-verify-checksum`: 复制有关校验功能.
* `slave-sql-verify-checksum`: 复制有关校验功能.
* `binlog-rows-query-log-events`: 启用之后可用于在二进制日志记录事件相关的信息, 可降低故障排除的复杂度(但是日志会变大)

## 整理知识点

---