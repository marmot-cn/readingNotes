#Mysql InnoDB非索引行检索使用行锁还是表锁?

---

在`深入浅出Mysql`一书中(20.mysql锁机制)中提及:

		InnoDB 只有通过`索引条件检索数据`,InnoDB才使用`行锁`,否则,InnoDB将使用`表锁`!
		

其实表示有**错误的**,我们在不同隔离级别下面测试.


**RR**

我们直接引用`深入浅出Mysql`一文中的例子(略作修改,给name字段添加PK,innodb如果没有主键索引会自动寻找唯一,如果没有唯一会自动创建一个隐含字段):

第一步:
	
		创建表table_no_index表,id没有索引,name为PK	
		mysql> create table tab_no_index(id int,name varchar(10)) engine=innodb;		Query OK, 0 rows affected (0.15 sec)		mysql> insert into tab_no_index values(1,'1'),(2,'2'),(3,'3'),		(4,'4');		Query OK, 4 rows affected (0.00 sec)		Records: 4  Duplicates: 0  Warnings: 0		mysql> alter table tab_no_index add primary key (name);		Query OK, 4 rows affected (6.01 sec)第二步:
* session1:
		mysql> set autocommit=0;		Query OK, 0 rows affected (0.00 sec)		mysql> select * from tab_no_index where id = 1 ;		+------+------+		| id   | name |		+------+------+		| 1    | 1    |		+------+------+		1 row in set (0.00 sec)
* session2:
		mysql> set autocommit=0;		Query OK, 0 rows affected (0.00 sec)		mysql> select * from tab_no_index where id = 2 ;		+------+------+		| id   | name |		+------+------+		| 2    | 2    |		+------+------+		1 row in set (0.00 sec)
第三步:
* session1:
	`给一行加了排他锁`
		mysql> select * from tab_no_index where id = 1 for update;		+------+------+		| id   | name |		+------+------+		| 1    | 1    |		+------+------+		1 row in set (0.00 sec)
* session2:
	`求其他行的排他锁时,却出现了锁等待.原因就是在没有索引的情况下,InnoDB只能使用表锁`
		mysql> select * from tab_no_index where id = 2 for update;				等待		
`注意这里的结论是错误的,我们后面在做结论`

第四步:
* session2:
	
	`插入数据失败`
			mysql> insert into tab_no_index (id,name) VALUES (6,6);
		ERROR 1205 (HY000): Lock wait timeout exceeded; try restarting transaction		**RC**
第一步:
	
		创建表table_no_index表,id没有索引,name为PK	
		mysql> create table tab_no_index(id int,name varchar(10)) engine=innodb;		Query OK, 0 rows affected (0.15 sec)		mysql> insert into tab_no_index values(1,'1'),(2,'2'),(3,'3'),		(4,'4');		Query OK, 4 rows affected (0.00 sec)		Records: 4  Duplicates: 0  Warnings: 0		mysql> alter table tab_no_index add primary key (name);		Query OK, 4 rows affected (6.01 sec)
第二步:
* session1:
		mysql> set tx_isolation='read-committed';
		Query OK, 0 rows affected (0.00 sec)

		mysql> select @@tx_isolation;
		+----------------+
		| @@tx_isolation |
		+----------------+
		| READ-COMMITTED |
		+----------------+
		1 row in set (0.00 sec)
		
		mysql> set autocommit=0;		Query OK, 0 rows affected (0.00 sec)
		
		mysql> select * from tab_no_index where id = 1 ;
		+------+------+
		| id   | name |
		+------+------+
		|    1 | 1    |
		+------+------+
		1 row in set (0.04 sec)
		* session2:
		mysql> set tx_isolation='read-committed';
		Query OK, 0 rows affected (0.00 sec)

		mysql> select @@tx_isolation;
		+----------------+
		| @@tx_isolation |
		+----------------+
		| READ-COMMITTED |
		+----------------+
		1 row in set (0.00 sec)
		
		mysql> set autocommit=0;		Query OK, 0 rows affected (0.00 sec)
				mysql> select * from tab_no_index where id = 2 ;		+------+------+		| id   | name |		+------+------+		| 2    | 2    |		+------+------+		1 row in set (0.00 sec)
第三步:
* session1:
	`给一行加了排他锁`
		mysql> select * from tab_no_index where id = 1 for update;		+------+------+		| id   | name |		+------+------+		| 1    | 1    |		+------+------+		1 row in set (0.00 sec)
* session2:
	`和RR相同,等待锁`
		mysql> select * from tab_no_index where id = 2 for update;				等待第四步:
* session2:
	
	`和RR不同,插入数据成功`
			mysql> insert into tab_no_index (id,name) VALUES (6,6);
		Query OK, 1 row affected (0.00 sec)
		

###结论

示例在`MySQL加锁处理分析`笔记中,`id列上没有索引,RC隔离级别`和`id列上没有索引,RR隔离级别`;

**RC**
聚簇索引上所有的记录，都被加上了X锁。无论记录是否满足条件，全部被加上X锁。既不是加表锁，也不是在满足条件的记录上加行锁。

**RR**
`聚簇索引上的所有记录`,都被加上了X锁.其次,聚簇索引每条记录间的间隙(GAP),也同时被加上了GAP锁.

**innodb不会主动升级表锁**

rc隔离级别下,有区别,记录仍旧可以插入.rr下,功能上无区别.但是`innodb不会主动升级表锁`.

