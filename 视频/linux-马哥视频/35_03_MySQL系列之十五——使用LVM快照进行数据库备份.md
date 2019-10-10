# 35_03_MySQL系列之十五——使用LVM快照进行数据库备份

---

## 笔记

---

### mysqldump 

逻辑备份, 还原过程涉及写操作, 会记录到二进制日志.

假设逻辑数据`3G`, 还原过程产生的二进制日志可能会高于`3G`.

还原过程生成的二进制日志没用, 还会产生大量`IO`.

**还原过程应该把二进制日志关掉, 还原之后在开启二进制日志**

`sql_log_bin`对当前会话改为`off`.

`SET sql_lob_bin=0`, 在当前会话临时性禁用二进制日志.

逻辑备份:

1. 浮点数据丢失精度.
2. 备份出来的数据更占用存储空间; 压缩后可大大节省空间.
3. 不适合对大数据库做完全备份. 速度慢, 消耗资源.
4. 只能对`MyISAM`做温备, 需要锁定表.
5. 对`InnoDB`做热备, 就算执行了`--lock-all-tables`锁表操作. 如果这时候有其他事物在进行, 会需要等待很久. 因为写锁优先级比读锁高, 或者别人在执行大事物. 而且也有可能因为事物日志存在, `mysql`后台进程自我调度把事物日志数据同步到数据文件中, 所有可能后台还在写(前端写不了).需要要通过命令`SHOW ENGINE INNODB STATUS`观察缓冲没有数据后在执行`mysqldump`.

```sql
mysql> SHOW ENGINE INNODB STATUS \G
*************************** 1. row ***************************
  Type: InnoDB
  Name:
Status:
=====================================
2017-08-23 13:19:07 7f44a4d10700 INNODB MONITOR OUTPUT
=====================================
Per second averages calculated from the last 29 seconds
-----------------
BACKGROUND THREAD
-----------------
srv_master_thread loops: 0 srv_active, 0 srv_shutdown, 4412 srv_idle
srv_master_thread log flush and writes: 4412
----------
SEMAPHORES
----------
OS WAIT ARRAY INFO: reservation count 2
OS WAIT ARRAY INFO: signal count 2
Mutex spin waits 0, rounds 0, OS waits 0
RW-shared spins 2, rounds 60, OS waits 2
RW-excl spins 0, rounds 0, OS waits 0
Spin rounds per wait: 0.00 mutex, 30.00 RW-shared, 0.00 RW-excl
------------
TRANSACTIONS
------------
Trx id counter 1797
Purge done for trx's n:o < 0 undo n:o < 0 state: running but idle
History list length 0
LIST OF TRANSACTIONS FOR EACH SESSION:
---TRANSACTION 0, not started
MySQL thread id 13, OS thread handle 0x7f44a4d10700, query id 41 localhost root init
SHOW ENGINE INNODB STATUS
--------
FILE I/O
--------
I/O thread 0 state: waiting for completed aio requests (insert buffer thread)
I/O thread 1 state: waiting for completed aio requests (log thread)
I/O thread 2 state: waiting for completed aio requests (read thread)
I/O thread 3 state: waiting for completed aio requests (read thread)
I/O thread 4 state: waiting for completed aio requests (read thread)
I/O thread 5 state: waiting for completed aio requests (read thread)
I/O thread 6 state: waiting for completed aio requests (write thread)
I/O thread 7 state: waiting for completed aio requests (write thread)
I/O thread 8 state: waiting for completed aio requests (write thread)
I/O thread 9 state: waiting for completed aio requests (write thread)
Pending normal aio reads: 0 [0, 0, 0, 0] , aio writes: 0 [0, 0, 0, 0] ,
 ibuf aio reads: 0, log i/o's: 0, sync i/o's: 0
Pending flushes (fsync) log: 0; buffer pool: 0
171 OS file reads, 5 OS file writes, 5 OS fsyncs
0.00 reads/s, 0 avg bytes/read, 0.00 writes/s, 0.00 fsyncs/s
-------------------------------------
INSERT BUFFER AND ADAPTIVE HASH INDEX
-------------------------------------
Ibuf: size 1, free list len 0, seg size 2, 0 merges
merged operations:
 insert 0, delete mark 0, delete 0
discarded operations:
 insert 0, delete mark 0, delete 0
Hash table size 276707, node heap has 0 buffer(s)
0.00 hash searches/s, 0.00 non-hash searches/s
---
LOG
---
Log sequence number 1625997
Log flushed up to   1625997
Pages flushed up to 1625997
Last checkpoint at  1625997
0 pending log writes, 0 pending chkp writes
8 log i/o's done, 0.00 log i/o's/second
----------------------
BUFFER POOL AND MEMORY
----------------------
Total memory allocated 137363456; in additional pool allocated 0
Dictionary memory allocated 52179
Buffer pool size   8192
Free buffers       8035
Database pages     157
Old database pages 0
Modified db pages  0
Pending reads 0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 0, not young 0
0.00 youngs/s, 0.00 non-youngs/s
Pages read 157, created 0, written 1
0.00 reads/s, 0.00 creates/s, 0.00 writes/s
No buffer pool page gets since the last printout
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 157, unzip_LRU len: 0
I/O sum[0]:cur[0], unzip sum[0]:cur[0]
--------------
ROW OPERATIONS
--------------
0 queries inside InnoDB, 0 queries in queue
0 read views open inside InnoDB
Main thread process no. 1260, id 139932761548544, state: sleeping
Number of rows inserted 0, updated 0, deleted 0, read 0
0.00 inserts/s, 0.00 updates/s, 0.00 deletes/s, 0.00 reads/s
----------------------------
END OF INNODB MONITOR OUTPUT
============================

1 row in set (0.00 sec)
```

备份出来的数据保存为文本, 会比原版数据文件大. 因为比如整数在导出来后会被保存为字符型. 也会导致浮点数据丢失精度.

#### `InnoDB`热备

`MVCC`, 如果事物隔离级别是`REPEATABLE-READ`, `--single-transaction`会启动一个大事物直到我们完成备份.

#### 恢复时临时禁用二进制日志

1. 找一个有权限恢复数据的`mysql`用户
2. `SET sql_lob_bin=0`关闭`sql_log_bin`.
3. 在**当前会话(不能在外部使用命令导入)**执行`source`或者`\.`导入备份文件.

### `SELECT * INTO OUTFILE '/path/file' FROM tb_name`

保存为纯文本数据(不是`sql`语句),各字段用制表符`tab`分开. 倒出的路径必须对用户`mysql`有写权限.

恢复需要使用命令.

`LOAD DATE INFILE '/path/file' INTO TABLE tb_name`. `tb_name`这张表需要事先创建好.

比`mysqldump`更节约空间, 恢复更快. 一般只做**单张表备份**场景使用.

恢复过程**写入二进制日志**, 因为既不是`DML`也不是`DDL`语句, 执行了一遍基于行的复制, 所以不是记录`sql`语句形式.

#### 从日志中的一个点恢复.

```shell
需要先 FLUSH LOGS;. 要不然从日志倒出的文件会报错bin log正在使用.
mysqlbinlog --start-position-xxx /path/binlog,0000xx > xxx.sql
```

### 几乎热备: LVM

#### 备份步骤

`snapshot`: 快照的访问路径, 就是那一刻的数据.

* 对所有表执行**读锁**`FLUSH TABLES WITH READ LOCK;`
* 滚动日志`FLUSH LOGS`;
* 备份当前二进制日志事件点 
	```
	mysql -uxx -pxx -e 'SHOW MASTER STATUS \G' > /backup/master-`date +%F`.info
	```
* 快照 `lvcreate -L # -s -p r -n LV_NAME /path/to/source_lv`
	```
	lvcreate -L 50M -s -p r(只读) -n mydata-snap /dev/xxx(mysql数据文件存放的卷)
	
	/dev/xxx 挂载在 /mydata 下, mysql 数据存放在 /mydata下
	```
* 释放锁 `UNLOCK TABLES`
* 挂载快照卷,备份

	```
	mount /dev/xxx/mydata-snap(快照) /mnt -o ro
	cd /mnt
	里面的mysql数据目录, 除了二进制日志所有其他文件都备份.
	如果只是备份单个数据库, 且开启了innodb_file_per_table, 则备份单独文件夹即可. 否则需要备份主文件夹下的 ibdata1. ib_logfile1 和 ib_logfile0 是事务日志.
	
	mkdir /backup/full-backup-`date +%F`
	cp -a ./* /backup/full-backup-xxx
	删除二进制日志
	rm /backup/full-backup-xxx/mysql-bin.* -f
	```
* 删除快照卷

	```
	umount /mnt
	lvremove --force /dev/xxx/mydata-snap
	cd /backup/full-backup-xxx
	```

* 怎量备份二进制日志

	```
	在原来的数据文件内, 备份二进制日志.
	在 mysql -e 'SHOW MASTER STATUS \G' 中产生的事件点, 我们只需要备份这个事件点以后的数据即可.
	mysqlbinlog --start-position=事件点 binlog > xxx
	mysqlbinlog binlog(后续的日志,假如备份时产生新的日志)
	或者通过匹配时间一次从多个二进制日志倒出
	mysqlbinlog --start-datetime='xxx-xxx-xx' mysql-bin.xxx mysql-bin.xx > xxx
	```

#### 前提

* 数据文件要在逻辑卷上.
* 此逻辑卷所在卷组必须有足够的空间使用快照卷.
	* 快照卷不一定要和原卷一样大, 快照卷主要目的是备份. 备份结束之前所有**变化的**数据需要容纳在快照卷上. 比如备份需要3小时, 数据库每1小时产生50MB数据, 那么创建一个大于150MB的快照卷就足够了. 只要备份完成后删除就可.
* 事务日志一定要和数据文件在一个卷上. 因为创建的快照, 假如事务日志和数据文件不在一个卷上, 那么快照的时间点可能不一致. 二进制还是分开存放在其他位置.
* 仍然需要二进制日志做即时点还原.

#### 正在执行事务时是否可以添加锁

会话A: 添加数据到一张表中, 事务不提交.

会话B: 另外一个会话锁表, 可以锁住.

会话A: 在添加语句添加不正常.

#### 快照还原

1. 把复制出来的文件复制到`mysql`数据目录下. `cp -a`保证属主属组是`mysql`.
2. 关闭当前会话二进制日志, 通过`source`导入增量.
3. 开启当前会话二进制日志


## 整理知识点

---