# 34_04_MySQL系列之十二——MySQL日志管理

---

## 笔记

---

### 日志

* 错误日志
* 查询日志, 一般查询日志
* 慢查询日志, 未必是真正的慢, 可能因为隔离级别获取资源慢
* 二进制日志, 记录格式二进制, 只记录`DML,DCL,DDL`语句. 记录引起或可能引起数据库变化的操作.
	* 实现`mysql`复制.
		* 从服务器不停的从主服务器上读取二进制日志, 保存在本地文件中(**中继日志**), 一个本地的`sql`和`io`线程, 读取日志事件并在本地执行一次.
	* `mysql`即时点恢复(恢复到一个指定的时间点上).
* 中继日志, 复制二进制日志.
* 事务日志. 保证事务的`acid`的重要组件, 将随机io转换为顺序io.

```sql
mysql> SHOW GLOBAL VARIABLES LIKE '%log%';
+-----------------------------------------+--------------------------------------+
| Variable_name                           | Value                                |
+-----------------------------------------+--------------------------------------+
| back_log                                | 80                                   |
| binlog_cache_size                       | 32768                                |
| binlog_checksum                         | CRC32                                |
| binlog_direct_non_transactional_updates | OFF                                  |
| binlog_error_action                     | IGNORE_ERROR                         |
| binlog_format                           | STATEMENT                            |
| binlog_gtid_simple_recovery             | OFF                                  |
| binlog_max_flush_queue_time             | 0                                    |
| binlog_order_commits                    | ON                                   |
| binlog_row_image                        | FULL                                 |
| binlog_rows_query_log_events            | OFF                                  |
| binlog_stmt_cache_size                  | 32768                                |
| binlogging_impossible_mode              | IGNORE_ERROR                         |
| expire_logs_days                        | 0                                    |
| general_log                             | OFF                                  |
| general_log_file                        | /var/lib/mysql/90af5cf2869e.log      |
| innodb_api_enable_binlog                | OFF                                  |
| innodb_flush_log_at_timeout             | 1                                    |
| innodb_flush_log_at_trx_commit          | 1                                    |
| innodb_locks_unsafe_for_binlog          | OFF                                  |
| innodb_log_buffer_size                  | 8388608                              |
| innodb_log_compressed_pages             | ON                                   |
| innodb_log_file_size                    | 50331648                             |
| innodb_log_files_in_group               | 2                                    |
| innodb_log_group_home_dir               | ./                                   |
| innodb_mirrored_log_groups              | 1                                    |
| innodb_online_alter_log_max_size        | 134217728                            |
| innodb_undo_logs                        | 128                                  |
| log_bin                                 | OFF                                  |
| log_bin_basename                        |                                      |
| log_bin_index                           |                                      |
| log_bin_trust_function_creators         | OFF                                  |
| log_bin_use_v1_row_events               | OFF                                  |
| log_error                               |                                      |
| log_output                              | FILE                                 |
| log_queries_not_using_indexes           | OFF                                  |
| log_slave_updates                       | OFF                                  |
| log_slow_admin_statements               | OFF                                  |
| log_slow_slave_statements               | OFF                                  |
| log_throttle_queries_not_using_indexes  | 0                                    |
| log_warnings                            | 1                                    |
| max_binlog_cache_size                   | 18446744073709547520                 |
| max_binlog_size                         | 1073741824                           |
| max_binlog_stmt_cache_size              | 18446744073709547520                 |
| max_relay_log_size                      | 0                                    |
| relay_log                               |                                      |
| relay_log_basename                      |                                      |
| relay_log_index                         |                                      |
| relay_log_info_file                     | relay-log.info                       |
| relay_log_info_repository               | FILE                                 |
| relay_log_purge                         | ON                                   |
| relay_log_recovery                      | OFF                                  |
| relay_log_space_limit                   | 0                                    |
| simplified_binlog_gtid_recovery         | OFF                                  |
| slow_query_log                          | OFF                                  |
| slow_query_log_file                     | /var/lib/mysql/90af5cf2869e-slow.log |
| sql_log_bin                             | ON                                   |
| sql_log_off                             | OFF                                  |
| sync_binlog                             | 0                                    |
| sync_relay_log                          | 10000                                |
| sync_relay_log_info                     | 10000                                |
+-----------------------------------------+--------------------------------------+
61 rows in set (0.02 sec)
```

* `binlog`: 二进制日志
	* `log_bin`: 是否记录二进制日志文件.
* `general_log`: 一般查询日志
* `innodb_log`: `innodb`的事务日志
* `log_slow_queries`: 慢查询日志
* `relay_log`: 中继日志

`mysql`在启动的时候, 默认情况下没有启用任何日志. 借助于配置文件启用日志功能.

开启户关闭日志, 可以动态操作.修改日志记录文件, 必须重启操作.

### 错误日志

* 服务器(`mysqld`)启动和关闭过程中的信息
* 服务器运行过程中的错误信息
* 事件(`event`)调度器运行一个事件时产生的信息
* 在从服务器上启动从服务器进程时产生的信息

### 一般查询日志

* `general_log`: 是否启用查询日志
* `general_log_file`: 启用后记录在哪个文件中
* `log_output`: 定义一般查询日志和慢查询日志的保存方式. 可以`TABLE,FILE`一起, 即记录在`TABLE`又记录在`FILE`
	* `TABLE`: 表
	* `FILE`: 文件
	* `NONE`: 不记录

如果`log_output`为`none`, 就算`general_log`开启了也不行.

### 慢查询日志

执行超出了`long_query_time`时间. 这里的语句执行时长为实际的执行时间, 而非在CPU上执行时长, 因此, 负载较重的的服务器上更容易产生慢查询. 但为使`秒钟`, 也支持毫秒级的解析度.

```sql
mysql> SHOW GLOBAL VARIABLES LIKE 'long_query_time';
+-----------------+-----------+
| Variable_name   | Value     |
+-----------------+-----------+
| long_query_time | 10.000000 |
+-----------------+-----------+
1 row in set (0.00 sec)
```

* `slow_query_log`: 是否启用慢查询日志
* `slow_query_log_file`: 慢查询日志文件的名称

### 二进制日志

* 索引文件. `mysql-bin.index`, 文本文件. 说明开始文件和结束文件是谁.
* 二进制日志文件. 

无论任何存储引擎, 只要在服务器级别任何可能引起数据库发生变化的操作都会记录下来.
记录格式为二进制.

`mysqlbinlog`命令查询日志.

`mysql`服务器每重启一次都需要滚动一次.

记录操作本身, 还要记录数据.

即时点恢复就是将记录的语句重放.

```sql
CREATE DATABASE mydb; //记录语句
INSERT INTO ... CURRENT_TIME(); //记录语句会出现问题.
```

#### 二进制日志的格式

* 基于语句`statement`(2次操作结果一样, 记录语句简单)
* 基于行`row`, 记录这一行的改变.
* 混合方式`mixed`. `mysql`自动选取是基于语句还是行的.

服务器日志每增长到一定量, 就滚动一次(生成新的).

#### 二进制日志事件

* 产生的时间`starttime`.
* 相对位置`position`, 相对当前文件所处的位置(当前事件的开始位置, 上个事件的结束位置).

查找定位事件可以根据**时间**和**位置**来查找.

`mysql-bin.00001`滚动后创建一个新的`mysql-bin.00002`以此类推. 老文件不动, 创建新的文件.

`mysql`重启也会创建一个新的日志文件.

#### 查看当前正在使用的二进制格式文件

`SHOW MASTER STATUS`. 查看当前正在使用的二进制格式文件

`SHOW BINLOG EVENTS IN '二进制日志文件名' [FROM pos]` 查看事件记录.

* `Pos`: 事件起始位置.
* `Eng_log_pos`: 事件结束位置.
* `Server_id`: 服务器`id`号, 表示由哪个服务器产生的.
* `Event_type`: 事件类型.
	* `Query`: 查询
	* `Intvar`: 内部变量
	* `Xid`: 事务id


```sql
mysql>SHOW MASTER STATUS;

+------------------+--------------------+------------------------+----------------------------+-------------------------------------------------------------------------------------------+

| File             | Position           | Binlog_Do_DB           | Binlog_Ignore_DB           | Executed_Gtid_Set                                                                         |

+------------------+--------------------+------------------------+----------------------------+-------------------------------------------------------------------------------------------+

| mysql-bin.001268 | 4987518            |                        |                            | d8a0a856-8f4f-11e6-b95e-7cd30ab90a64:1-30865946,
e685bcb0-8f4f-11e6-b95f-7cd30ab79b6e:1-4 |

+------------------+--------------------+------------------------+----------------------------+-------------------------------------------------------------------------------------------+

mysql>SHOW BINLOG EVENTS IN 'mysql-bin.001268';

+--------------------+---------------+----------------------+---------------------+-----------------------+-------------------------------------------------------------------------------------------+

| Log_name           | Pos           | Event_type           | Server_id           | End_log_pos           | Info                                                                                      |

+--------------------+---------------+----------------------+---------------------+-----------------------+-------------------------------------------------------------------------------------------+

| mysql-bin.001268   | 4             | Format_desc          |          3467990554 | 120                   | Server ver: 5.6.16-log, Binlog ver: 4                                                     |

| mysql-bin.001268   | 120           | Previous_gtids       |          3467990554 | 231                   | d8a0a856-8f4f-11e6-b95e-7cd30ab90a64:1-30863277,
e685bcb0-8f4f-11e6-b95f-7cd30ab79b6e:1-4 |

| mysql-bin.001268   | 231           | Gtid                 |          3467990554 | 279                   | SET @@SESSION.GTID_NEXT= 'd8a0a856-8f4f-11e6-b95e-7cd30ab90a64:30863278'                  |

| mysql-bin.001268   | 279           | Query                |          3467990554 | 353                   | BEGIN                                                                                     |

| mysql-bin.001268   | 353           | Table_map            |          3467990554 | 428                   | table_id: 9486 (cattle.cluster_membership)                                                |

| mysql-bin.001268   | 428           | Update_rows_v1       |          3467990554 | 722                   | table_id: 9486 flags: STMT_END_F                                                          |

| mysql-bin.001268   | 722           | Xid                  |          3467990554 | 753                   | COMMIT /* xid=783379588 */                                                                |

| mysql-bin.001268   | 753           | Gtid                 |          3467990554 | 801                   | SET @@SESSION.GTID_NEXT= 'd8a0a856-8f4f-11e6-b95e-7cd30ab90a64:30863279'                  |

| mysql-bin.001268   | 801           | Query                |          3467990554 | 875                   | BEGIN                                                                                     |

| mysql-bin.001268   | 875           | Table_map            |          3467990554 | 950                   | table_id: 9486 (cattle.cluster_membership)                                                |

| mysql-bin.001268   | 950           | Update_rows_v1       |          3467990554 | 1244                  | table_id: 9486 flags: STMT_END_F                                                          |

| mysql-bin.001268   | 1244          | Xid                  |          3467990554 | 1275                  | COMMIT /* xid=783379629 */                                                                |

| mysql-bin.001268   | 1275          | Gtid                 |          3467990554 | 1323                  | SET @@SESSION.GTID_NEXT= 'd8a0a856-8f4f-11e6-b95e-7cd30ab90a64:30863280'                  |

| mysql-bin.001268   | 1323          | Query                |          3467990554 | 1391                  | BEGIN                                                                                     |

| mysql-bin.001268   | 1391          | Table_map            |          3467990554 | 1453                  | table_id: 71 (mysql.ha_health_check)                                                      |

| mysql-bin.001268   | 1453          | Update_rows_v1       |          3467990554 | 1509                  | table_id: 71 flags: STMT_END_F                                               
```

#### mysqlbinlog

因为日志是二进制日子, 所以不能用`cat`命令, 只能用`mysqlbinlog`命令.

查询二进制日志的命令.

* `--start-date-time xxxx-xx-xx HH:MM:SS`
* `--stop-datetime`
* `--start-position`
* `--stop-position`

```shell
mysqlbinlog --start-position=107 --stop-position=358 mysql-bin.001268
```

可以把`mysqlbinlog`重定向到文件中, 因为里面都是`sql`语句, 数据库可以直接执行.

#### 手动滚动日志.

`FLUSH LOGS`, 只会滚动二进制日志. 其他日志,错误日志关闭子在打开. 从服务器上, 滚动中继日志, 其他日志,错误日志关闭在打开.

#### 删除二进制日志

备份后删除二进制日志.

`PURGE BINARY LOGS TO 'log_name'`.

这个`log_name`**之前**的所有日志都删除. 该`log_name`文件不删除.

#### `SHOW BINARY LOGS`

查看当前所有二进制日志文件.

## 整理知识点

---