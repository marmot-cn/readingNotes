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

## 整理知识点

---