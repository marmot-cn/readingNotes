# 浅谈RR隔离到底会不会出现幻读

---

最近拜读了相关数据库书籍,描述了RR隔离级别`会导致幻读`.但是在拜读了何登成的`MySQL 加锁处理分析`一文后,描述了Mysql InnoDB的RR级别不会导致幻读(用间隙锁解决了幻读).

这篇文章简单描述下Mysql InnDB如何用**间隙锁**解决了幻读.

**为什么有的描述说RR有幻读,有的没有**

标准SQL规范中定义的RR隔离级别只是解决了`不可重复读`和`脏读`,但是允许`幻读`.但是`mysql innodb引擎的实现,跟标准有所不同`也就是用间隙锁解决了幻读.


### 幻读和不可重复读的区别

**不可重复读**

不可重复读的重点是`修改`.同样的条件,你读取过的数据,再次读取出来发现值不一样了.

1. 在事务1中,Mary读取了自己的工资为1000,操作并没有完成.
2. 事务2中,这时财务人员修改了Mary的工资为2000,并提交了事务. 
3. 在事务1中，Mary 再次读取自己的工资时，工资变为了2000 

在一个事务中前后两次读取的结果并不致,导致了不可重复读.

**幻读**

幻读的重点在于新增或者删除
同样的条件, 第1次和第2次读出来的记录数不一样

例子:
目前工资为1000的员工有10人.

1. 读取所有工资为1000的员工,共读取10条记录 
2. 这时另一个事务向employee表插入了一条员工记录，工资也为1000 
3. 事务1再次读取所有工资为1000的员工,共读取到了11条记录,这就产生了幻像读.

# RC和RR

首先我们用RC(有幻读)和RR(Mysql InnoDB无幻读)来举个列子.

**首先创建一张表**

		CREATE TABLE IF NOT EXISTS `test` (
  		  `id` int(10) NOT NULL,
  		  `name` varchar(20) NOT NULL,
  		  PRIMARY KEY (`name`),
  		  KEY `id` (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
		
注意我们的`name`是`主键索引`,`id`是`非唯一索引`.

**插入一些数据**

		INSERT INTO `test` (`id`, `name`) VALUES
		(2, 'a'),
		(6, 'b'),
		(10, 'c'),
		(10, 'd'),
		(11, 'e'),
		(15, 'f');
		

**现在我们看下我们的数据**

		SELECT * FROM test;
		
		+----+------+
		| id | name |
		+----+------+
		|  2 | a    |
		|  6 | b    |
		| 10 | c    |
		| 10 | d    |
		| 11 | e    |
		| 15 | f    |
		+----+------+
		6 rows in set (0.00 sec)

**修改隔离级别**

		SET [SESSION | GLOBAL] TRANSACTION ISOLATION LEVEL {READ UNCOMMITTED | READ COMMITTED | REPEATABLE READ | SERIALIZABLE}
		
`默认的行为(不带session和global)是为下一个(未开始)事务设置隔离级别`.如果你使用GLOBAL关键字,语句在全局对从那点开始创建的所有新连接(除了不存在的连接)设置默认事务级别.你需要SUPER权限来做这个.使用SESSION关键字为将来在当前连接上执行的事务设置默认事务级别. 任何客户端都能自由改变会话隔离级别(甚至在事务的中间),或者为下一个事务设置隔离级别.

**RC测试不可重复读**

* 第一步:

	`我们首先开启2个session,并修改隔离级别`:

	session1:

		mysql> select @@tx_isolation;
		+-----------------+
		| @@tx_isolation  |
		+-----------------+
		| REPEATABLE-READ |
		+-----------------+
		1 row in set (0.00 sec)
		
		mysql> set tx_isolation='read-committed';
		Query OK, 0 rows affected (0.00 sec)

		mysql> select @@tx_isolation;
		+----------------+
		| @@tx_isolation |
		+----------------+
		| READ-COMMITTED |
		+----------------+
		1 row in set (0.00 sec)
		
		mysql> start transaction;
		Query OK, 0 rows affected (0.00 sec)
		
	session2:
		
		mysql> select @@tx_isolation;
		+-----------------+
		| @@tx_isolation  |
		+-----------------+
		| REPEATABLE-READ |
		+-----------------+
		1 row in set (0.00 sec)
		
		mysql> set tx_isolation='read-committed';
		Query OK, 0 rows affected (0.00 sec)

		mysql> select @@tx_isolation;
		+----------------+
		| @@tx_isolation |
		+----------------+
		| READ-COMMITTED |
		+----------------+
		1 row in set (0.00 sec)
		
		mysql> start transaction;
		Query OK, 0 rows affected (0.00 sec)

* 第二步:

	`seesion2`修改数据并提交
	
	session1:
		
		mysql> SELECT * FROM test WHERE id=10;
		+----+------+
		| id | name |
		+----+------+
		| 10 | c    |
		| 10 | d    |
		+----+------+
		2 rows in set (0.00 sec)
		
	session2:
		
		mysql> UPDATE test SET id=11 WHERE name='c';
		Query OK, 1 row affected (0.34 sec)

		mysql> commit;
		Query OK, 0 rows affected (0.05 sec)

* 第三步:

	`session1`发现了**不可重复读**,因为`session2`提交了修改数据.
	
	session1:
	
		mysql> SELECT * FROM test WHERE id=10;
		+----+------+
		| id | name |
		+----+------+
		| 10 | d    |
		+----+------+
		1 row in set (0.00 sec)

**RC测试幻读**

* 第一步:

	`我们首先开启2个session,并修改隔离级别`:

	session1:

		mysql> select @@tx_isolation;
		+-----------------+
		| @@tx_isolation  |
		+-----------------+
		| REPEATABLE-READ |
		+-----------------+
		1 row in set (0.00 sec)
		
		mysql> set tx_isolation='read-committed';
		Query OK, 0 rows affected (0.00 sec)

		mysql> select @@tx_isolation;
		+----------------+
		| @@tx_isolation |
		+----------------+
		| READ-COMMITTED |
		+----------------+
		1 row in set (0.00 sec)
		
		mysql> start transaction;
		Query OK, 0 rows affected (0.00 sec)
		
	session2:
		
		mysql> select @@tx_isolation;
		+-----------------+
		| @@tx_isolation  |
		+-----------------+
		| REPEATABLE-READ |
		+-----------------+
		1 row in set (0.00 sec)
		
		mysql> set tx_isolation='read-committed';
		Query OK, 0 rows affected (0.00 sec)

		mysql> select @@tx_isolation;
		+----------------+
		| @@tx_isolation |
		+----------------+
		| READ-COMMITTED |
		+----------------+
		1 row in set (0.00 sec)
		
		mysql> start transaction;
		Query OK, 0 rows affected (0.00 sec)
		
* 第二步:

	`session1`使用当前读,`seesion2`插入数据并提交
	
	session1:
		
		mysql> SELECT * FROM test WHERE id=10 for update;
		+----+------+
		| id | name |
		+----+------+
		| 10 | c    |
		| 10 | d    |
		+----+------+
		2 rows in set (0.00 sec)
		
	session2:
		
		mysql> INSERT INTO `test` (`id`, `name`) VALUES (10,'g');
		Query OK, 1 row affected (0.34 sec)

		mysql> commit;
		Query OK, 0 rows affected (0.05 sec)

* 第三步:

	`session1`发现了**幻读**,因为`session2`提交了插入数据.
	
	session1:
	
		mysql> SELECT * FROM test WHERE id=10 for update;
		+----+------+
		| id | name |
		+----+------+
		| 10 | c    |
		| 10 | d    |
		| 10 | g    |
		+----+------+
		3 rows in set (6.97 sec)


**RR测试不可重复读**

* 第一步:

	`我们首先开启2个session,开启事务`:

	session1:

		mysql> select @@tx_isolation;
		+-----------------+
		| @@tx_isolation  |
		+-----------------+
		| REPEATABLE-READ |
		+-----------------+
		1 row in set (0.00 sec)
		
		mysql> start transaction;
		Query OK, 0 rows affected (0.00 sec)
		
	session2:
		
		mysql> select @@tx_isolation;
		+-----------------+
		| @@tx_isolation  |
		+-----------------+
		| REPEATABLE-READ |
		+-----------------+
		1 row in set (0.00 sec)
		
		mysql> start transaction;
		Query OK, 0 rows affected (0.00 sec)
		
* 第二步:

	`seesion2`修改数据并提交
	
	session1:
		
		mysql> SELECT * FROM test WHERE id=10;
		+----+------+
		| id | name |
		+----+------+
		| 10 | c    |
		| 10 | d    |
		+----+------+
		2 rows in set (0.00 sec)
		
	session2:
		
		mysql> UPDATE test SET id=11 WHERE name='c';
		Query OK, 1 row affected (0.34 sec)

		mysql> commit;
		Query OK, 0 rows affected (0.05 sec)
		
* 第三步:

	`session1没有发现`**不可重复读**.
	
	session1:
	
		SELECT * FROM test WHERE id=10;
		+----+------+
		| id | name |
		+----+------+
		| 10 | c    |
		| 10 | d    |
		+----+------+
		2 rows in set (0.00 sec)
		
* 第四步:
	
	`session1退出事务,发现数据变化`
	
	session1:
	
		mysql> rollback;
		Query OK, 0 rows affected (0.00 sec)
		
		mysql> SELECT * FROM test WHERE id=10;
		+----+------+
		| id | name |
		+----+------+
		| 10 | d    |
		+----+------+
		1 row in set (0.00 sec)

**RR测试幻读**

* 第一步:

	`我们首先开启2个session,开启事务`:

	session1:

		mysql> select @@tx_isolation;
		+-----------------+
		| @@tx_isolation  |
		+-----------------+
		| REPEATABLE-READ |
		+-----------------+
		1 row in set (0.00 sec)
		
		mysql> start transaction;
		Query OK, 0 rows affected (0.00 sec)
		
	session2:
		
		mysql> select @@tx_isolation;
		+-----------------+
		| @@tx_isolation  |
		+-----------------+
		| REPEATABLE-READ |
		+-----------------+
		1 row in set (0.00 sec)
		
		mysql> start transaction;
		Query OK, 0 rows affected (0.00 sec)

* 第二步:
		
	`session1`使用当前读,`seesion2`插入数据并提交
	
	session1:
		
		mysql> SELECT * FROM test WHERE id=10 for update;
		+----+------+
		| id | name |
		+----+------+
		| 10 | c    |
		| 10 | d    |
		+----+------+
		2 rows in set (0.00 sec)
		
	session2:
		
		INSERT INTO `test` (`id`, `name`) VALUES (10,'g');
		ERROR 1205 (HY000): Lock wait timeout exceeded; try restarting transaction
		插入阻塞,因为间隙锁.
		INSERT INTO `test` (`id`, `name`) VALUES (8,'g');
		ERROR 1205 (HY000): Lock wait timeout exceeded; try restarting transaction
		插入阻塞,因为间隙锁.
		
		INSERT INTO `test` (`id`, `name`) VALUES (5,'g');
		Query OK, 1 row affected (0.05 sec)
		插入成功

* 第三步:

	`session1`没有发现**幻读**
		I
	session1:
	
		mysql> SELECT * FROM test WHERE id=10 for update;
		+----+------+
		| id | name |
		+----+------+
		| 10 | c    |
		| 10 | d    |
		+----+------+
		2 rows in set (0.00 sec)

* 结论:

	`B+树索引的有序性,满足条件的项一定是连续存放的`.
	
	1. 记录[6,b]之前，不会插入id=10的记录
	2. [6,b]与[10,c]间可以插入`[10, aa]`
	3. [10,c]与[10,d]间，可以插入[10,cc]
	4. [10,d]与[11,e]间可以插入满足条件的[10,dd],[10,z]
	5. 而[11,e]之后也不会插入满足条件的记录

	为了保证`[6,b]与[10,c]间`,`[10,c]与[10,d]间`和`[10,d]与[11,e]`不会插入新的满足条件的记录,MySQL选择了用GAP锁,将这三个GAP给锁起来.
		