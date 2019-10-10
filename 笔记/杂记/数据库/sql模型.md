# sql模型

---

通过定义某些规定, 限制用户行为, 并定义对应的处理机制.

## 常见的模型

### `ANSI`

宽松模式, 对插入数据进行校验, 如果不符合定义类型或长度, 对数据类型调整或截断保存, 报`warning`警告.

### `TRADITIONAL`

严格模式, 当向mysql数据库插入数据时, 进行数据的严格校验, 保证错误数据不能插入, 报`error`错误. 用于事物时,会进行事物的回滚.

### `STRICT_TRANS_TABLES`

严格模式, 进行数据的严格校验, 不允许向一个支持事物的表中插入非法数据, 报error错误.

### `STRICT_ALL_TABLES`

未设置的情况下, 所有的非法数值都允许, 返回警告信息. 设置以后只要违反数据规则, 都不允许填入, 并返回错误.

### `ANSI QUOTES`

双引号和反引号作用相同, 只能用来引用字段名称/表名等, 单引号只能引用在字符串. mysql中默认3者可以随意引用.

### `IGNORE_SPACE`

在内建函数中忽略多余空格.

### `no_engine_substitution`

如果我把引擎指定成一个并不存在的引擎, 这个时候mysql可以有两种行为供选择 

1. 直接报错,当`sql_mode`中包涵`no_engine_subtitution`时, 如果`create table`时指定的`engine`项不被支持, 这个时候`mysql`会支持报错.
2. 把表的存储引擎替换成innodb, 在`sql_mode`中不包涵`no_engine_subtitution`且`create table`中`engine`子句指定的存储引擎不被支持时, `mysql`会把表的引擎改为`innodb`

#### 没有`no_engine_substitution`

```sql
mysql> SHOW VARIABLES LIKE 'sql_mode';
+---------------+--------------------------------------------+
| Variable_name | Value                                      |
+---------------+--------------------------------------------+
| sql_mode      | STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION |
+---------------+--------------------------------------------+
1 row in set (0.00 sec)
mysql> SET SESSION sql_mode=STRICT_TRANS_TABLES;
Query OK, 0 rows affected (0.01 sec)

mysql> SHOW VARIABLES LIKE 'sql_mode';
+---------------+---------------------+
| Variable_name | Value               |
+---------------+---------------------+
| sql_mode      | STRICT_TRANS_TABLES |
+---------------+---------------------+
1 row in set (0.00 sec)

创建一张不存在的表
mysql> use test;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> create table t(x int) engine=federated;
Query OK, 0 rows affected, 2 warnings (0.03 sec)

mysql> show create table t;
+-------+------------------------------------------------------------------------------------+
| Table | Create Table                                                                       |
+-------+------------------------------------------------------------------------------------+
| t     | CREATE TABLE `t` (
  `x` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 |
+-------+------------------------------------------------------------------------------------+
1 row in set (0.00 sec)
```

### 有`no_engine_substitution`

```sql
mysql> SHOW VARIABLES LIKE 'sql_mode';
+---------------+--------------------------------------------+
| Variable_name | Value                                      |
+---------------+--------------------------------------------+
| sql_mode      | STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION |
+---------------+--------------------------------------------+
1 row in set (0.00 sec)

mysql> use test;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> create table t(x int) engine=federated;
ERROR 1286 (42000): Unknown storage engine 'federated'
```

## 示例

```sql

mysql> SHOW GLOBAL VARIABLES LIKE 'sql_mode';
+---------------+--------------------------------------------+
| Variable_name | Value                                      |
+---------------+--------------------------------------------+
| sql_mode      | STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION |
+---------------+--------------------------------------------+
1 row in set (0.00 sec)

mysql>
```