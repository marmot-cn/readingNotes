# 35_04_MySQL系列之十六——使用xtrabackup进行数据库备份

---

## 笔记

---

### 二进制日志相关的几个选项

* `innodb_support_xa={TRUE|FALSE}` 分布式`xa`事务.
* `sync_binlog=#` 设定多久同步一次二进制日志至磁盘文件中. 一般建议设置为`1`.
	* `0`表示不同步.
	* 任何正数值都表示对二进制每多少次写操作之后同步一次.
	* 当`autocommit`的值为`1`时, 每条语句的执行都会引起二进制日志同步, 否则, 每个事务的提交会引起二进制日志同步. 

在`MySQL`中系统默认的设置是`sync_binlog=0`, 也就是不做任何强制性的磁盘刷新指令, 这时候的性能是最好的, 但是风险也是最大的. 因为一旦系统Crash, 在binlog_cache中的所有binlog信息都会被丢失. 而当设置为“1”的时候, 是最安全但是性能损耗最大的设置. 因为当设置为1的时候, 即使系统Crash, 也最多丢失binlog_cache中未完成的一个事务, 对实际数据没有任何实质性影响.

从以往经验和相关测试来看, 对于高并发事务的系统来说, “sync_binlog”设置为`0`和设置为1的系统写入性能差距可能高达5倍甚至更多.

### percona

`ibbackup`: InnoDB 在线物理备份.

* 全量
* 增量

`lvm`方式是**几乎**热备, 并不是完全热备, 施加锁的时间可能会长一些. `mylvmbackup`工具 自动的快照备份在逻辑卷上的数据.

`xtrabackup`开源的. 可以对`xtradb`(`innodb`引擎的升级版)和`innodb`存储引擎备份.

* 全量
* 增量

#### xtrabackup + 二进制日志备份

xtrabackup + 备份二进制日志方式. 需要手动定期备份二进制日志. 支持对`innodb`增量备份, 但是对`myisam`不支持增量. `mysql`库下面的大多数表都是`myisam`引擎.

##### 完全备份

`innobackupex  --host=name --port=PORT --user=DBUSER(需要用备份权限的用户) --pasword=DBUSERPASS /path/to/BACKUP-DIR/`

如果要使用一个最小权限的用户进行备份, 则可基于如下命令创建此类用户.

`PROCESS`这个权限是新版新增的.

```sql
mysql> CREATE USER 'bkpuser'@'localhost' IDENTIFIED BY 'password';
Query OK, 0 rows affected (0.00 sec)
mysql> GRANT RELOAD, LOCK TABLES, PROCESS, REPLICATION CLIENT ON *.* TO 'bkpuser'@'localhost';
Query OK, 0 rows affected (0.00 sec)
mysql> SHOW GRANTS FOR 'bkpuser'@'localhost';
+--------------------------------------------------------------------------------------------------------------------------------------------------+
| Grants for bkpuser@localhost                                                                                                                     |
+--------------------------------------------------------------------------------------------------------------------------------------------------+
| GRANT RELOAD, LOCK TABLES, REPLICATION CLIENT ON *.* TO 'bkpuser'@'localhost' IDENTIFIED BY PASSWORD '*2470C0C06DEE42FD1618BB99005ADCA2EC9D1E19' |
+--------------------------------------------------------------------------------------------------------------------------------------------------+
1 row in set (0.00 sec)
mysql> FLUSH PRIVILEGES;
Query OK, 0 rows affected (0.00 sec)
```

使用工具备份

```shell
root@iZ94xwu3is8Z ~]# innobackupex --user=bkpuser --password=password ./backup
170823 23:47:18 innobackupex: Starting the backup operation

IMPORTANT: Please check that the backup run completes successfully.
           At the end of a successful backup run innobackupex
           prints "completed OK!".

perl: warning: Setting locale failed.
perl: warning: Please check that your locale settings:
	LANGUAGE = (unset),
	LC_ALL = (unset),
	LC_CTYPE = "UTF-8",
	LANG = "en_US.UTF-8"
    are supported and installed on your system.
perl: warning: Falling back to the standard locale ("C").
Can't locate Digest/MD5.pm in @INC (@INC contains: /usr/local/lib64/perl5 /usr/local/share/perl5 /usr/lib64/perl5/vendor_perl /usr/share/perl5/vendor_perl /usr/lib64/perl5 /usr/share/perl5 .) at - line 693.
BEGIN failed--compilation aborted at - line 693.
170823 23:47:18 Connecting to MySQL server host: localhost, user: bkpuser, password: set, port: 0, socket: /var/lib/mysql/mysql.sock
Using server version 5.6.37-log
innobackupex version 2.3.6 based on MySQL server 5.6.24 Linux (x86_64) (revision id: )
xtrabackup: uses posix_fadvise().
xtrabackup: cd to /var/lib/mysql
xtrabackup: open files limit requested 0, set to 65535
xtrabackup: using the following InnoDB configuration:
xtrabackup:   innodb_data_home_dir = ./
xtrabackup:   innodb_data_file_path = ibdata1:12M:autoextend
xtrabackup:   innodb_log_group_home_dir = ./
xtrabackup:   innodb_log_files_in_group = 2
xtrabackup:   innodb_log_file_size = 50331648
170823 23:47:18 >> log scanned up to (1649486)
xtrabackup: Generating a list of tablespaces
170823 23:47:19 [01] Copying ./ibdata1 to /root/backup/2017-08-23_23-47-18/ibdata1
170823 23:47:19 [01]        ...done
170823 23:47:19 [01] Copying ./mysql/slave_relay_log_info.ibd to /root/backup/2017-08-23_23-47-18/mysql/slave_relay_log_info.ibd
170823 23:47:19 [01]        ...done
170823 23:47:19 [01] Copying ./mysql/slave_worker_info.ibd to /root/backup/2017-08-23_23-47-18/mysql/slave_worker_info.ibd
170823 23:47:19 [01]        ...done
170823 23:47:19 [01] Copying ./mysql/innodb_index_stats.ibd to /root/backup/2017-08-23_23-47-18/mysql/innodb_index_stats.ibd
170823 23:47:19 [01]        ...done
...
170823 23:47:20 Finished backing up non-InnoDB tables and files
170823 23:47:20 [00] Writing xtrabackup_binlog_info
170823 23:47:20 [00]        ...done
170823 23:47:20 Executing FLUSH NO_WRITE_TO_BINLOG ENGINE LOGS...
xtrabackup: The latest check point (for incremental): '1649486'
xtrabackup: Stopping log copying thread.
.170823 23:47:20 >> log scanned up to (1649486)

170823 23:47:21 Executing UNLOCK TABLES
170823 23:47:21 All tables unlocked
170823 23:47:21 Backup created in directory '/root/backup/2017-08-23_23-47-18'
MySQL binlog position: filename 'mysql-bin.000003', position '3542'
170823 23:47:21 [00] Writing backup-my.cnf
170823 23:47:21 [00]        ...done
170823 23:47:21 [00] Writing xtrabackup_info
170823 23:47:21 [00]        ...done
xtrabackup: Transaction log of lsn (1649486) to (1649486) was copied.
170823 23:47:21 completed OK!

[root@iZ94xwu3is8Z 2017-08-23_23-47-18]# ls
backup-my.cnf  ibdata1  mysql  performance_schema  test  xtrabackup_binlog_info  xtrabackup_checkpoints  xtrabackup_info  xtrabackup_logfile

[root@iZ94xwu3is8Z 2017-08-23_23-47-18]# file xtrabackup_binlog_info
xtrabackup_binlog_info: ASCII text
[root@iZ94xwu3is8Z 2017-08-23_23-47-18]# cat xtrabackup_binlog_info
mysql-bin.000003	3542
[root@iZ94xwu3is8Z 2017-08-23_23-47-18]# file xtrabackup_logfile
xtrabackup_logfile: data
[root@iZ94xwu3is8Z 2017-08-23_23-47-18]# file xtrabackup_checkpoints
xtrabackup_checkpoints: ASCII text
[root@iZ94xwu3is8Z 2017-08-23_23-47-18]# cat xtrabackup_checkpoints
backup_type = full-backuped
from_lsn = 0
to_lsn = 1649486
last_lsn = 1649486
compact = 0
recover_binlog_info = 0
```

* `ibdata1`表空间文件.
* `backup-my.cnf`备份的配置文件.
* `xtrabackup_binlog_info`文本文件, 存储备份那一刻二进制文件名以及位置信息.
* `xtrabackup_binlog_pos_innodb`二进制日志文件及用于`InnoDB`或`XtraDB`表的二进制日志文件的当前`position`.
* `xtrabackup_logfile`数据文件, 
* `xtrabackup_checkpoints`备份类型(如完全或增量),备份状态(如是否已经为prepared状态)和LSN(日志序列号)范围信息.
	* `backup_type`备份类型: `full-backuped`完全备份
	* `from_lsn`从哪个逻辑号码开始的备份.(`lsn`: innodb每一个数据块(存储在磁盘上)都有一个日志序列号, innodb会在内部维护每个数据库的日志序列号, 当数据发生了改变这个序列号会增加, 可以根据这个号码做增量备份). 因为这次是完全备份, 所以从`0`开始.
	* `to_lsn`最后一次的日志序列号.
	* `last_lsn`, 增量备份从该号码开始. 靠追踪`lsn`来完成增量备份.

##### 恢复前准备工作

`innobackup`命令的`--apply-log`选项可用于实现下面功能.

* 将尚未提交的事务回滚
* 将已经提交的事务, 从事务日志同步到数据文件

```shell
[root@iZ94xwu3is8Z ~]# innobackupex --apply-log backup/2017-08-23_23-47-18
170824 13:24:45 innobackupex: Starting the apply-log operation

IMPORTANT: Please check that the apply-log run completes successfully.
           At the end of a successful apply-log run innobackupex
           prints "completed OK!".

innobackupex version 2.3.6 based on MySQL server 5.6.24 Linux (x86_64) (revision id: )
xtrabackup: cd to /root/backup/2017-08-23_23-47-18
xtrabackup: This target seems to be not prepared yet.
xtrabackup: xtrabackup_logfile detected: size=2097152, start_lsn=(1649486)
xtrabackup: using the following InnoDB configuration for recovery:
xtrabackup:   innodb_data_home_dir = ./
xtrabackup:   innodb_data_file_path = ibdata1:12M:autoextend
xtrabackup:   innodb_log_group_home_dir = ./
xtrabackup:   innodb_log_files_in_group = 1
xtrabackup:   innodb_log_file_size = 2097152
xtrabackup: using the following InnoDB configuration for recovery:
xtrabackup:   innodb_data_home_dir = ./
xtrabackup:   innodb_data_file_path = ibdata1:12M:autoextend
xtrabackup:   innodb_log_group_home_dir = ./
xtrabackup:   innodb_log_files_in_group = 1
xtrabackup:   innodb_log_file_size = 2097152
xtrabackup: Starting InnoDB instance for recovery.
xtrabackup: Using 104857600 bytes for buffer pool (set by --use-memory parameter)
InnoDB: Using atomics to ref count buffer pool pages
InnoDB: The InnoDB memory heap is disabled
InnoDB: Mutexes and rw_locks use GCC atomic builtins
InnoDB: Memory barrier is not used
InnoDB: Compressed tables use zlib 1.2.7
InnoDB: Using CPU crc32 instructions
InnoDB: Initializing buffer pool, size = 100.0M
InnoDB: Completed initialization of buffer pool
InnoDB: Highest supported file format is Barracuda.
InnoDB: The log sequence numbers 1625987 and 1625987 in ibdata files do not match the log sequence number 1649486 in the ib_logfiles!
InnoDB: Database was not shutdown normally!
InnoDB: Starting crash recovery.
InnoDB: Reading tablespace information from the .ibd files...
InnoDB: Restoring possible half-written data pages
InnoDB: from the doublewrite buffer...
InnoDB: 128 rollback segment(s) are active.
InnoDB: Waiting for purge to start
InnoDB: 5.6.24 started; log sequence number 1649486
xtrabackup: Last MySQL binlog file position 2996, file name mysql-bin.000003

xtrabackup: starting shutdown with innodb_fast_shutdown = 1
InnoDB: FTS optimize thread exiting.
InnoDB: Starting shutdown...
InnoDB: Shutdown completed; log sequence number 1649496
xtrabackup: using the following InnoDB configuration for recovery:
xtrabackup:   innodb_data_home_dir = ./
xtrabackup:   innodb_data_file_path = ibdata1:12M:autoextend
xtrabackup:   innodb_log_group_home_dir = ./
xtrabackup:   innodb_log_files_in_group = 2
xtrabackup:   innodb_log_file_size = 50331648
InnoDB: Using atomics to ref count buffer pool pages
InnoDB: The InnoDB memory heap is disabled
InnoDB: Mutexes and rw_locks use GCC atomic builtins
InnoDB: Memory barrier is not used
InnoDB: Compressed tables use zlib 1.2.7
InnoDB: Using CPU crc32 instructions
InnoDB: Initializing buffer pool, size = 100.0M
InnoDB: Completed initialization of buffer pool
InnoDB: Setting log file ./ib_logfile101 size to 48 MB
InnoDB: Setting log file ./ib_logfile1 size to 48 MB
InnoDB: Renaming log file ./ib_logfile101 to ./ib_logfile0
InnoDB: New log files created, LSN=1649496
InnoDB: Highest supported file format is Barracuda.
InnoDB: 128 rollback segment(s) are active.
InnoDB: Waiting for purge to start
InnoDB: 5.6.24 started; log sequence number 1649676
xtrabackup: starting shutdown with innodb_fast_shutdown = 1
InnoDB: FTS optimize thread exiting.
InnoDB: Starting shutdown...
InnoDB: Shutdown completed; log sequence number 1649686
170824 13:24:51 completed OK!
```

##### 备份后数据又发生变化, 我们需要即时点还原

我们插入几条新的数据(备份后).

```shell
[root@iZ94xwu3is8Z ~]# mysql -uroot -p
Enter password:
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 37
Server version: 5.6.37-log MySQL Community Server (GPL)

Copyright (c) 2000, 2017, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> use test;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from test;
+----+-----+
| id | num |
+----+-----+
|  1 |   5 |
|  2 |  20 |
|  3 |   0 |
|  4 |  10 |
+----+-----+
4 rows in set (0.00 sec)

mysql> insert into test values(5,11);
Query OK, 1 row affected (0.01 sec)

mysql> insert into test values(6,12);
Query OK, 1 row affected (0.00 sec)

mysql> select * from test;
+----+-----+
| id | num |
+----+-----+
|  1 |   5 |
|  2 |  20 |
|  3 |   0 |
|  4 |  10 |
|  5 |  11 |
|  6 |  12 |
+----+-----+
5 rows in set (0.00 sec)
```

##### 备份二进制日志文件

```shell
[root@iZ94xwu3is8Z ~]# cd /var/lib/mysql
[root@iZ94xwu3is8Z mysql]# ls
auto.cnf  ib_logfile0  ib_logfile1  ibdata1  mysql  mysql-bin.000001  mysql-bin.000002  mysql-bin.000003  mysql-bin.index  mysql.sock  performance_schema  test

日志滚动(因为备份正在写的二进制日志, 导出数据会出错)
...
mysql> FLUSH LOGS;
Query OK, 0 rows affected (0.04 sec)
mysql> SHOW MASTER STATUS;
+------------------+----------+--------------+------------------+-------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000004 |      120 |              |                  |                   |
+------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)
...

备份二进制日志
[root@iZ94xwu3is8Z mysql]# cp mysql-bin.000003 /root/
```

##### 模拟mysql数据库损毁

```shell
[root@iZ94xwu3is8Z mysql]# systemctl stop mysql.service

删除mysql数据目录
[root@iZ94xwu3is8Z mysql]# rm -rf ./*
```

##### 从一个完全备份中恢复数据

* 恢复时不需要启动数据库.
* 恢复后文件的属主和属组不一定正确, 还需要修改.

`innobackupex --copy-back /path/to/BACKUP-DIR`

```shell
[root@iZ94xwu3is8Z mysql]# innobackupex --copy-back /root/backup/2017-08-23_23-47-18/
170824 13:42:56 innobackupex: Starting the copy-back operation

IMPORTANT: Please check that the copy-back run completes successfully.
           At the end of a successful copy-back run innobackupex
           prints "completed OK!".

innobackupex version 2.3.6 based on MySQL server 5.6.24 Linux (x86_64) (revision id: )
170824 13:42:56 [01] Copying ib_logfile0 to /var/lib/mysql/ib_logfile0
170824 13:42:57 [01]        ...done
170824 13:42:57 [01] Copying ib_logfile1 to /var/lib/mysql/ib_logfile1
170824 13:42:58 [01]        ...done
170824 13:42:58 [01] Copying ibdata1 to /var/lib/mysql/ibdata1
170824 13:42:58 [01]        ...done
170824 13:42:58 [01] Copying ./xtrabackup_binlog_pos_innodb to /var/lib/mysql/xtrabackup_binlog_pos_innodb
170824 13:42:58 [01]        ...done
170824 13:42:58 [01] Copying ./performance_schema/events_waits_summary_by_thread_by_event_name.frm to /var/lib/mysql/performance_schema/events_waits_summary_by_thread_by_event_name.frm
170824 13:42:58 [01]        ...done
...
170824 13:42:59 [01] Copying ./test/test.frm to /var/lib/mysql/test/test.frm
170824 13:42:59 [01]        ...done
170824 13:42:59 completed OK!

[root@iZ94xwu3is8Z mysql]# pwd
/var/lib/mysql
[root@iZ94xwu3is8Z mysql]# ls
ib_logfile0  ib_logfile1  ibdata1  mysql  performance_schema  test  xtrabackup_binlog_pos_innodb  xtrabackup_info

[root@iZ94xwu3is8Z mysql]# cd ..
[root@iZ94xwu3is8Z lib]# pwd
/var/lib
[root@iZ94xwu3is8Z lib]# chown -R mysql:mysql mysql

启动mysql
[root@iZ94xwu3is8Z lib]# systemctl start mysql.service
[root@iZ94xwu3is8Z lib]# mysql -uroot -p
Enter password:
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 3
Server version: 5.6.37-log MySQL Community Server (GPL)

Copyright (c) 2000, 2017, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| test               |
+--------------------+
4 rows in set (0.00 sec)

mysql> use test;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from test;
+----+-----+
| id | num |
+----+-----+
|  1 |   5 |
|  2 |  20 |
|  3 |   0 |
|  4 |  10 |
+----+-----+
4 rows in set (0.01 sec)

mysql> show master status;
+------------------+----------+--------------+------------------+-------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000001 |      120 |              |                  |                   |
+------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)

即时点还原(还原增量)
1. 查看还原点位置
[root@iZ94xwu3is8Z mysql]# cat xtrabackup_info
uuid = 599468fd-881a-11e7-8d33-00163e008577
name =
tool_name = innobackupex
tool_command = --user=bkpuser --password=... ./backup
tool_version = 2.3.6
ibbackup_version = 2.3.6
server_version = 5.6.37-log
start_time = 2017-08-23 23:47:18
end_time = 2017-08-23 23:47:21
lock_time = 0
binlog_pos = filename 'mysql-bin.000003', position '3542'
innodb_from_lsn = 0
innodb_to_lsn = 1649486
partial = N
incremental = N
format = file
compact = N
compressed = N
encrypted = N
[root@iZ94xwu3is8Z mysql]# cat /root/backup/2017-08-23_23-47-18/xtrabackup_binlog_info
mysql-bin.000003	3542
[root@iZ94xwu3is8Z mysql]# mysqlbinlog /root/mysql-bin.000003  --start-position=3542 > /tmp/test.sql

登录mysql, 禁止二进制日志, 还原, 然后在开启.

SET sql_log_bin=0;
SOURCE /tmp/test.sql
SET sql_log_bin=1;
```

#### xtrabackup + 增量备份

##### 完全备份

```shell
[root@iZ94xwu3is8Z ~]# innobackupex --user=bkpuser --password=password ./backup
170824 14:36:17 innobackupex: Starting the backup operation
...
[root@iZ94xwu3is8Z ~]# ls backup/2017-08-24_14-36-17/
backup-my.cnf  ibdata1  mysql  performance_schema  test  xtrabackup_binlog_info  xtrabackup_checkpoints  xtrabackup_info  xtrabackup_logfile
```

##### 增量备份

`innobackupex --incremental /backup(增量备份文件存放的目录) --incremental-basedir=BASEDIR`

其中, `BASEDIR`指的是完全备份所在的目录, 此命令执行结束后, `innobackupex`命令会在`/backup`目录中创建一个新的以时间命名的沐浴露以存放所有的增量备份数据. 另外, 在执行过增量备份之后在进行第一次进行增量备份时, 其`--incremental-basedir`应该指向上一次的增量备份所在的目录. 

```shel
在插入两条数据, 然后开始做增量备份

mysql> use test;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from test;
+----+-----+
| id | num |
+----+-----+
|  1 |   5 |
|  2 |  20 |
|  3 |   0 |
|  4 |  10 |
|  5 |  11 |
|  6 |  12 |
+----+-----+
6 rows in set (0.01 sec)

mysql> insert into test values(7,13),(8,14);
Query OK, 2 rows affected (0.01 sec)
Records: 2  Duplicates: 0  Warnings: 0

mysql> select * from test;
+----+-----+
| id | num |
+----+-----+
|  1 |   5 |
|  2 |  20 |
|  3 |   0 |
|  4 |  10 |
|  5 |  11 |
|  6 |  12 |
|  7 |  13 |
|  8 |  14 |
+----+-----+
8 rows in set (0.00 sec)

[root@iZ94xwu3is8Z 2017-08-24_14-36-17]# innobackupex --user=bkpuser --password=password --incremental /backup --incremental-basedir=/root/backup/2017-08-24_14-36-17
[root@iZ94xwu3is8Z 2017-08-24_14-36-17]# ls /backup/2017-08-24_14-42-38
backup-my.cnf  ibdata1.delta  ibdata1.meta  mysql  performance_schema  test  xtrabackup_binlog_info  xtrabackup_checkpoints  xtrabackup_info  xtrabackup_logfile
```

完全备份的`last_lsn`和后续增量备份的`from_lsn`是连续的.

```shell
完全备份的lsn
[root@iZ94xwu3is8Z ~]# cat /root/backup/2017-08-24_14-36-17/xtrabackup_checkpoints
backup_type = full-backuped
from_lsn = 0
to_lsn = 1653068
last_lsn = 1653068
compact = 0
recover_binlog_info = 0

增量备份的lsn
[root@iZ94xwu3is8Z ~]# cat /backup/2017-08-24_14-42-38/xtrabackup_checkpoints
backup_type = incremental
from_lsn = 1653068
to_lsn = 1654717
last_lsn = 1654717
compact = 0
recover_binlog_info = 0
```

##### 执行第二次增量备份

`basedir`是**上一次增量备份的路径**

##### 模拟数据库崩溃

```shell
[root@iZ94xwu3is8Z mysql]# systemctl stop mysql.service
[root@iZ94xwu3is8Z mysql]# pwd
/var/lib/mysql
[root@iZ94xwu3is8Z mysql]# rm -rf ./*
```

##### 恢复

增量备份的准备(prepare)和完全备份有一些不同:

1. 需要在每个备份(包括完全和各个增量备份)上, 将已经提交的事务进行"重放", “重放"之后, 所有的备份数据将合并到完全备份上.
2. 基于所有的备份将未提交的事务进行"回滚"

操作为: `innobackupex --apply-log --redo-only BASE-DIR`

接着执行: `innobackupex --apply-log --redo-only BASE-DIR --incremental-dir=INCREMENTAL-DIR-1`. 如果有多个增量备份, 重复执行

只执行`redo`操作是因为, 可能第一次完全备份没有执行过的, 在第一次增量中已经执行了. 如果执行`undo`操作数据就会发生回滚了.

```shell
完全备份
[root@iZ94xwu3is8Z ~]# innobackupex --apply-log --redo-only /root/backup/2017-08-24_14-36-17/
[root@iZ94xwu3is8Z ~]# innobackupex --apply-log --redo-only /root/backup/2017-08-24_14-36-17/ --incremental-dir=/backup/2017-08-24_14-42-38/

所有的数据都合并到完全备份上
[root@iZ94xwu3is8Z mysql]# innobackupex --copy-back /root/backup/2017-08-24_14-36-17/
170824 15:02:19 innobackupex: Starting the copy-back operation

IMPORTANT: Please check that the copy-back run completes successfully.
           At the end of a successful copy-back run innobackupex
           prints "completed OK!".

innobackupex version 2.3.6 based on MySQL server 5.6.24 Linux (x86_64) (revision id: )
170824 15:02:19 [01] Copying ib_logfile0 to /var/lib/mysql/ib_logfile0
170824 15:02:19 [01]        ...done
170824 15:02:20 [01] Copying ib_logfile1 to /var/lib/mysql/ib_logfile1
170824 15:02:21 [01]        ...done
170824 15:02:21 [01] Copying ibdata1 to /var/lib/mysql/ibdata1
170824 15:02:21 [01]        ...done
170824 15:02:21 [01] Copying ./xtrabackup_binlog_pos_innodb to /var/lib/mysql/xtrabackup_binlog_pos_innodb
170824 15:02:21 [01]        ...done
170824 15:02:21 [01] Copying ./performance_schema/events_waits_summary_by_thread_by_event_name.frm to /var/lib/mysql/performance_schema/events_waits_summary_by_thread_by_event_name.frm
...

[root@iZ94xwu3is8Z ~]# chown -R mysql:mysql /var/lib/mysql
[root@iZ94xwu3is8Z ~]# systemctl start mysql.service
[root@iZ94xwu3is8Z mysql]# mysql -uroot -p
Enter password:
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 2
Server version: 5.6.37-log MySQL Community Server (GPL)

Copyright (c) 2000, 2017, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> use test;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from test;
+----+-----+
| id | num |
+----+-----+
|  1 |   5 |
|  2 |  20 |
|  3 |   0 |
|  4 |  10 |
|  5 |  11 |
|  6 |  12 |
|  7 |  13 |
|  8 |  14 |
+----+-----+
8 rows in set (0.00 sec)
```

#### 导入导出单张表

必须启用了`innodb_file_per_table`(导出和导入都需要)选项和`innodb_expand_import`(导入的数据库需要启动)选项.

##### 导出表

导出表是在备份的`prepare`阶段进行的, 因此, 一旦完全备份完成, 就可以在`prepare`过程中通过`--export`选项将其表导出了.

`innobackupex --apply-log --export /path/to/backup`

此命令会为每个`innodb`表的表空间创建一个以`.exp`结尾的文件, 这些以`.exp`结尾的文件则可以用于导入至其他服务器.

##### 导入表

要在`mysql`服务器上导入来自于其他服务器的某`innodb`表, 需要现在当前服务器上创建一个跟原表结构一致的表, 而后才能实现将表导入.

创建表`xxx`.

将其表空间删除. `ALTER TABLE xxx.xxx DISCARD TABLESPACE;`(把刚才创建的表空间删除);

将来自于"导出"表的服务器的`mytable`表的`mytable.ibd`和`mytable.exp`文件复制到当前服务器的数据目录, 然后使用如下命令将其"导入":

`ALTER TABLE xxx.xxx IMPORT TABLESPACE;`;

### Xtrabackup的"流"及"备份压缩"功能

将备份的数据通过`STDOUT`传输给`tar`程序进行归档, 而不是默认的直接保存至其备份目录中.

```
innobackupex --stream=tar /backup | gzip > /backup/`date +%F_%H-%M-%S`.tar.gz
```

也可以使用如下命令将数据备份至其他服务器:

```
innobackupex --stream=tar /backup | ssh user@xxx "cat - > /backup/`date +%F_%H-%M-%S`.tar.gz"
```

在执行本地备份时, 还可以使用`--parallel`选项对多个文件进行并行复制. 此选项用于指定在复制时启动的线程数目. 也需要启用`innodb_file_per_table`选项或共享的表空间通过`innodb_data_file_path`选项存储在多个`ibdata`文件中.

`innobackupex --parallel /path/to/backup`

同时, `innobackupex`备份的数据也可以存储至远程主机, 这可以使用`--remote-host`(**确认新版已经移除**)选项来实现.

`innobackupex --remote-host=xxx /path/IN/REMOTE/to/backup`

```shell
[root@iZ94xwu3is8Z mysql]# innobackupex --user=bkpuser --password=password --remote-host=root@120.25.87.35 /backup
```

## 整理知识点

### xtrabackup_binlog_info 和 xtrabackup_binlog_pos_innodb 区别

---