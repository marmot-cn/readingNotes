# binlog，redo log，undo log区别

---

## 概述

`binlog`是MySQL Server层记录的日志.

`redo log`是InnoDB存储引擎层的日志.

选择`binlog`日志作为replication我想主要原因是MySQL的特点就是支持多存储引擎, 为了兼容绝大部分引擎来支持复制这个特性，那么自然要采用MySQL Server自己记录的日志.

## 什么是binlog

Mysql的`binlog`日志作用是用来记录mysql内部**增删改查等对mysql数据库有更新的内容的记录(对数据库的改动)**, 对数据库的查询select或show等不会被binlog日志记录; 主要用于数据库的**主从复制以及增量恢复**.

## 事务日志(redo log)

减少提交事务时的开销. 因为日志中已经记录了事务, 就无须在每个事务提交时把缓冲池的脏块数据刷新(`flush`)到磁盘中. 

事务修改的数据和索引通常会映射到表空间的随机位置, 所以刷新这些变更需要很多随机`I/O`.

通常会初始化2个或更多的`ib_logfile`存储`redo log`, 由参数 `innodb_log_files_in_group`确定个数, 命名从`ib_logfile0`开始, 依次写满 `ib_logfile`并顺序重用(in a circular fashion). 如果最后1个`ib_logfile`被写满, 而第一个`ib_logfile`中所有记录的事务对数据的变更已经被持久化到磁盘中, 将清空并重用之. **不会覆盖还没应用到数据文件的日志记录**

`InnoDB`使用一个后台线程智能地刷新这些变更到数据文件. 可以**批量组合写入**

### 日志的大小

#### 日志太小

如果日志太小, `InnoDB`将必须做更多的检查点, 导致更多的日志写. 在日志没有空间继续写入前, 必须等待变更被应用到数据文件.

#### 日志太大

在崩溃恢复`InnoDB`可能不得不做大量的工作.

## undo log

记录了数据修改的前镜像. 存放于ibdata中.

用于在实例故障恢复时, 借助`undo log`将尚未`commit`的事务, 回滚到事务开始前的状态.