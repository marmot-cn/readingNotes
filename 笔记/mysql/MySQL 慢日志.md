# MySQL 慢日志

## 概述

### 参数解释

* `slow_query_log = on`: 开启慢查询日志
* `slow_query_log_file = filename`: 指定慢日志文件存放位置，可以为空，系统会给一个缺省的文件host_name-slow.log
* `long_query_time = #`: 设定慢查询的阀值，超出次设定值的SQL即被记录到慢查询日志，缺省值为10s
* `log-slow-admin-statements`: 慢管理语句例如OPTIMIZE TABLE、ANALYZE TABLE和ALTER TABLE记入慢查询日志(我们不开启)
* `log_queries_not_using_indexes`: 记录没有使用索引的查询语句(我们不开启)

### 配置文件

```
[mysqld]
slow_query_log_file = /var/lib/mysql/slow.log
slow_query_log = on
long_query_time = 2 
```

### 示例

检查配置

```
mysql> show variables like '%slow%';
+---------------------------+-------------------------+
| Variable_name             | Value                   |
+---------------------------+-------------------------+
| log_slow_admin_statements | OFF                     |
| log_slow_slave_statements | OFF                     |
| slow_launch_time          | 2                       |
| slow_query_log            | ON                      |
| slow_query_log_file       | /var/lib/mysql/slow.log |
+---------------------------+-------------------------+
5 rows in set (0.01 sec)
```

执行一个慢查询

```
mysql> select sleep(3);
+----------+
| sleep(3) |
+----------+
|        0 |
+----------+
1 row in set (3.01 sec)
```

查看慢日志

```
root@cec9de744a90:/# cat /var/lib/mysql/slow.log
mysqld, Version: 5.7.21-log (MySQL Community Server (GPL)). started with:
Tcp port: 0  Unix socket: /var/run/mysqld/mysqld.sock
Time                 Id Command    Argument
# Time: 2018-10-26T06:58:46.377680Z
# User@Host: root[root] @ localhost []  Id:     2
# Query_time: 3.005214  Lock_time: 0.000000 Rows_sent: 1  Rows_examined: 0
SET timestamp=1540537126;
select sleep(3);
```

使用`mysqldumpslow mysql-slow.log`工具查看

```
root@cec9de744a90:/# mysqldumpslow /var/lib/mysql/slow.log

Reading mysql slow query log from /var/lib/mysql/slow.log
Count: 1  Time=0.00s (0s)  Lock=0.00s (0s)  Rows=0.0 (0), 0users@0hosts
  mysqld, Version: N.N.N-log (MySQL Community Server (GPL)). started with:
  # Time: N-N-26T06:N:N.377680Z
  # User@Host: root[root] @ localhost []  Id:     N
  # Query_time: N.N  Lock_time: N.N Rows_sent: N  Rows_examined: N
  SET timestamp=N;
  select sleep(N)
```

### 日志

日志收集, 通过`fluntd`收集日志, 具体配置可见[fluentd](https://github.com/chloroplast1983/docker-image-fluentd)