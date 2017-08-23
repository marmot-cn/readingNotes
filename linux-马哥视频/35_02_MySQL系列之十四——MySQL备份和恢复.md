# 35_02_MySQL系列之十四——MySQL备份和恢复

---

## 笔记

---

### MySQL的备份和还原

备份: 副本.

`RAID1`, `RAID0`: 保证硬件损坏而不会业务中止, 逻辑上不提供保证.

备份和`RAID`是两个不同层次的概念.

#### 备份类型

如果通过`cp`命令直接复制, 复制数据的时间点不一致. 数据库停了, 在复制.

根据备份时服务器是否在线:

* 热备份: 在线备份, 读写操作可继续进行不受影响.
* 温备份: 能读不能写, 仅可以执行读操作.
* 冷备份: 关机备份, 离线备份. 读写操作均不能进行.

备份方式:

* 物理备份: 直接复制数据文件.
* 逻辑备份: 将数据导出至文本文件中.

根据备份内容:

* 完全备份: 备份整个数据集.
* 增量备份: 仅备份上次完全备份或增量备份以后变化的数据.
* 差异备份: 仅备份上次完全备份以来变化的数据.

差异备份和增量备份区别:

* 增量备份: 上次完全, 周一的增量, 周二的增量. 恢复需要: 上次完全,周一,周二.
* 差异备份: 上次完全, 周一和上次完全的差异, 周二和上次完全的差异. 恢复需要: 上次完全,周二差异即可.比增量备份多占据空间,但是还原起来快.

备份策略:

* 完全+增量.
* 完全+差异.

备份时间选取:

* 数据变化的频率.
* 能够忍受的还原时长.

#### 备份什么

* 数据文件.
* 配置文件.
* 备份日志
	* 二进日志, 做即时点还原.
	* 事务日志, 因为数据文件中可能存在运行了一半的事务.
* 任何和`mysql`相关的配置文件.

一个月一次完全, 每天一次增量. 当一天过了一半,挂了需要还原. 这时候增量备份数据文件只能还原到昨天的数据. 就需要二进制日志做即时点还原.

#### 还原测试

备份后需要进行还原测试. 演练还原. 经常进行消防预演.

要有备份还原策略, 以及紧急状态还原策略.

#### 热备份

备份需要时间, 可能会有数据时间前后不一致. 所以技术复杂度最高.

* `MyISAM`热备几乎是不可能的.借助逻辑卷快照可以使用几乎热备.否则只能使用温备,以共享方式锁定`MyISAM`所有表(别人不能改).
* `InnoDB`可以. 
	* `xtrabackup`物理备份工具.
	* `mysqldump`逻辑备份工具.

`MySQL`做一个**从服务器**, 需要备份停下从服务器, 然后备份, 备份完成后在开启从服务器.或者把从服务器的从进程停了(就变成只读了)可以进行温备份.

* 物理备份: 因为文件系统不一致, 移植性不强. 速度快. 不需要借助`mysql`, 直接在文件系统层面完成.
* 逻辑备份: 借助`mysql`服务器进程将数据从表中读出,另存为文件. 速度慢. 丢失浮点数, 因为要将数据保存至文本字符. 方便使用文本工具直接对其处理. 可移植能力强. 跨`MySQL`服务器版本.

在原服务器做备份, 影响服务器性能.

#### 备份工具

mysql自带工具:

* `mysqldunmp`: 逻辑备份工具. 对MyISAM(温备), InnoDB(热备).
* `mysqlhotcopy`: **温备份**,实际是冷备, 需要锁表. 物理备份工具.

文件系统工具:

* `cp`: 冷备. 借助逻辑卷可以实现几乎热备, 创建好快照后在`cp`.
* `逻辑卷快照功能`, 实现几乎热备.
	1. `mysql> FLUSH TABLES;`
	2. `mysql> LOCK TABLES;`
	3. 创建快照.
	4. 释放锁.
	5. 复制数据.

逻辑卷快照功能: 对`MyISAM`引擎来说锁表可以实现, 但是对`InnoDB`来说锁表后背后有可能有写操作:

1. 事务可能在日志中, 正在从日志往数据文件中同步.
2. 已经提交的事务在内存中, 正在往日志中同步.

所以需要监控确保`InnoDB`的缓冲区的内容都已经同步完成.

第三方工具:

* `ibbackup`(`innodb`): 商业工具.
* `xtrabackup`: 开源工具.

##### mysqldump

小规模数据库使用可以, 大规模使用效率太低.

mysqldump(完全备份,逻辑方式) + 二进制文件 

* 完全(mysqldump备份整个数据库) + 增量(每天备份二进制信息)
* 完全 + 差异

###### 选项

`mysqldump DB_NAME [tb1] [tb2]`

* `db_name [tb1] [tb2]`: 备份指定数据库,或者库中的某张表. 如果指定一个数据库,则不备份创建数据库的语句. 还原需要手动创建数据库.
* `--all-database`: 备份所有库. 会自动创建`CREATE DATABASE`命令.
* `--databases DB_NAME,DB_NAME`: 备份指定库. 会自动创建`CREATE DATABASE`命令.
* `--master-data={0|1|2}`:
	* `0`: 不记录二进制日志文件及事件位置.
	* `1`: 以`CHANGE MASTER TO`的方式记录位置, 可用于恢复后直接启动**从服务器**的.
	* `2`: 以`CHANGE MASTER TO`的方式记录位置, 但默认被注释.
* `--lock-all-tables`: 备份前先锁定**所有**表, 温备.
* `--flush-logs`: 备份前锁完表, 执行日志滚动. 这样可以把之前的二进制日志一起备份了. 增量备份, 每天备份二进制日志即可.
* `--single-transaction`: 启动热备, 如果指定库中表类型均为`innodb`. 不用在启动`--lock-all-tables`选项. 可能备份时间较长.
* `--events`: 备份事件.
* `--routines`: 备份存储过程和函数.
* `--triggers`: 备份触发器.

#### 备份策略

每周完全+每日增量.

完全备份: mysqldump

增量备份: 备份二进制日志文件(flush logs), 或者从时间点之间备份.

##### 手动模拟示例

完全备份:

我自己服务器上数据库全是`innodb`引擎,可以执行热备., `mysqldump -uroot -p --single-transaction --flush-logs --all-databases --master-data=2 > /root/alldatabases-`date +%F-%H-%S`.sql`

```shell
这里的命令 --lock-all-tables 是执行的温备.

[root@iZ944l0t308Z ansible]# mysqldump -uroot -p --lock-all-tables --flush-logs --all-databases --master-data=2 > /root/alldatabases-`date +%F-%H-%S`.sql
Enter password:

[root@iZ944l0t308Z ~]# less alldatabases-2017-08-22-12-50.sql
...
-- CHANGE MASTER TO MASTER_LOG_FILE='mysql-bin.000004', MASTER_LOG_POS=120;
```

连入mysql, 如果要删除旧的bin-log, 最好**先备份**(这里没有演示备份)在删除:

```shell
mysql> show binary logs;
+------------------+-----------+
| Log_name         | File_size |
+------------------+-----------+
| mysql-bin.000001 |       386 |
| mysql-bin.000002 |       167 |
| mysql-bin.000003 |       167 |
| mysql-bin.000004 |       120 |
+------------------+-----------+
4 rows in set (0.00 sec)
mysql> PURGE BINARY LOGS TO 'mysql-bin.000004';
Query OK, 0 rows affected (0.02 sec)

mysql> show binary logs;
+------------------+-----------+
| Log_name         | File_size |
+------------------+-----------+
| mysql-bin.000004 |       120 |
+------------------+-----------+
1 row in set (0.00 sec)
```

我们对数据做一些修改,测试增量备份.

```shell
mysql> use test;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from test;
+----+-----+
| id | num |
+----+-----+
|  1 |   2 |
|  2 |   2 |
|  3 |   4 |
|  4 |   5 |
+----+-----+
4 rows in set (0.00 sec)
mysql> update test set num=6 where id=4;
Query OK, 1 row affected (0.01 sec)
Rows matched: 1  Changed: 1  Warnings: 0

mysql> select * from test;
+----+-----+
| id | num |
+----+-----+
|  1 |   2 |
|  2 |   2 |
|  3 |   4 |
|  4 |   6 |
+----+-----+
4 rows in set (0.00 sec)
```

增量备份:

```shell
滚动日志
[root@iZ944l0t308Z ~]# mysql -uroot -p
Enter password:
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 12
Server version: 5.6.37-log MySQL Community Server (GPL)

Copyright (c) 2000, 2017, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> flush logs;
Query OK, 0 rows affected (0.03 sec)

其中 mysql-bin.000004 就是我们的增量.
[root@iZ944l0t308Z ~]# ls /var/lib/mysql
auto.cnf  ib_logfile0  ib_logfile1  ibdata1  mysql  mysql-bin.000004  mysql-bin.000005  mysql-bin.index  mysql.sock  performance_schema  test

[root@iZ944l0t308Z ~]# mysqlbinlog /var/lib/mysql/mysql-bin.000004 > /root/mon-incremental.sql
```

在新增一些数据:

```shell
mysql> use test;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from test;
+----+-----+
| id | num |
+----+-----+
|  1 |   2 |
|  2 |   2 |
|  3 |   4 |
|  4 |   6 |
+----+-----+
4 rows in set (0.00 sec)

mysql> insert into test values(5,7);
Query OK, 1 row affected (0.01 sec)
```

恢复:

```shell
复制最新的日志
[root@iZ944l0t308Z lib]# cp mysql/mysql-bin.000005 /root/

清空mysql数据目录
[root@iZ944l0t308Z mysql]# ls
[root@iZ944l0t308Z mysql]# pwd
/var/lib/mysql

关闭mysql服务
[root@iZ944l0t308Z mysql]# systemctl stop mysql.service

现在也启动不起来mysql, 需要初始化数据库

还原所有的数据库
mysql -uroot -p < alldatabases-2017-08-22-12-50.sql

还原第一次增量备份
mysql -uroot -p < mon-incremental.sql

还原最后一刻数据, 从二进制日志导入
mysqlbinlog mysql-bin.000005 | mysql -uroot -p
```

###### 示例: `--master-data`

通过被注释的事件位置, 增量备份可以从该事件位置以后来备份.

```shell
root@90af5cf2869e:~# mysqldump -uroot -p --master-data=2 test > ./test-`date +%F-%H-%S`.sql
Enter password:


可以看见被注释的 CHANGE MASTER TO MASTER_LOG_FILE='mysql-bin.000001', MASTER_LOG_POS=339;

[root@iZ944l0t308Z ansible]# cat test-2017-08-22-12-19.sql
-- MySQL dump 10.13  Distrib 5.6.37, for Linux (x86_64)
--
-- Host: localhost    Database: test
-- ------------------------------------------------------
-- Server version	5.6.37-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Position to start replication or point-in-time recovery from
--

-- CHANGE MASTER TO MASTER_LOG_FILE='mysql-bin.000001', MASTER_LOG_POS=339;

--
-- Table structure for table `test`
--

DROP TABLE IF EXISTS `test`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `test` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `num` smallint(5) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `test`
--

LOCK TABLES `test` WRITE;
/*!40000 ALTER TABLE `test` DISABLE KEYS */;
INSERT INTO `test` VALUES (1,2),(2,2),(3,4),(4,5);
/*!40000 ALTER TABLE `test` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-08-22 12:21:22
```

###### 示例: 备份单个数据库

**如果有其他人正在往数据库写数据, 会导致时间点的不一致.**

```shell
将整个数据库备份成批量插入的sql语句. 没有创建数据库的语句.

root@90af5cf2869e:~# mysqldump -uroot -p test > ./test.sql
Enter password:
root@90af5cf2869e:~# cat saas_product.sql
root@90af5cf2869e:~# mysqldump -uroot -p test > ./test.sql
Enter password:
root@90af5cf2869e:~# cat test.sql
-- MySQL dump 10.13  Distrib 5.6.37, for Linux (x86_64)
--
-- Host: localhost    Database: test
-- ------------------------------------------------------
-- Server version	5.6.37

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `test`
--

DROP TABLE IF EXISTS `test`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `test` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `num` smallint(5) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `test`
--

LOCK TABLES `test` WRITE;
/*!40000 ALTER TABLE `test` DISABLE KEYS */;
INSERT INTO `test` VALUES (1,5),(2,20),(3,0);
/*!40000 ALTER TABLE `test` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-08-21  6:44:00

还原
1. 创建数据库
mysql> create database test_backup;
Query OK, 1 row affected (0.01 sec)

2. 还原数据
root@90af5cf2869e:~# mysql -uroot -p test_backup < test.sql
Enter password:
```

###### 示例: 锁表备份单个数据库

```shell
刷新并以读锁锁住表, 可以读数据, 但是不能写数据
mysql> FLUSH TABLES WITH READ LOCK;
Query OK, 0 rows affected (0.04 sec);

FLUSH LOGS; 滚动日志.

SHOW BINARY LOGS; 最新的日志就是锁表后新插入的数据.

另外一个终端, 查询可以, 不能写
root@90af5cf2869e:/var/lib# mysql -uroot -p123456
Warning: Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 412
Server version: 5.6.37 MySQL Community Server (GPL)

Copyright (c) 2000, 2017, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> use test;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> show tables;
+----------------+
| Tables_in_test |
+----------------+
| test           |
+----------------+
1 row in set (0.01 sec)

mysql> select * from test;
+----+-----+
| id | num |
+----+-----+
|  1 |   5 |
|  2 |  20 |
|  3 |   0 |
+----+-----+
3 rows in set (0.00 sec)
mysql> insert into test values (4,10);
... 卡主,等待释放锁

备份数据库 mysqldump

mysql> UNLOCK TABLES;
Query OK, 0 rows affected (0.00 sec)
```

###### 示例: 备份全部库

```shell
[root@iZ944l0t308Z ansible]# mysqldump -uroot -p --lock-all-tables --flush-logs --all-databases --master-data=2 > all.sql
Enter password:

...
日志已经滚动为最新

mysql> show binary logs;
+------------------+-----------+
| Log_name         | File_size |
+------------------+-----------+
| mysql-bin.000001 |       386 |
| mysql-bin.000002 |       120 |
+------------------+-----------+
...
```


## 整理知识点

---