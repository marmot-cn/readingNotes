# 33_02_MySQL系列之六——MySQL管理表和索引

---

## 笔记

---

### `SQL`语句

* 数据库
* 表
* 索引
* 视图
* DML语句

#### 数据库

##### 创建数据库

`CREATE DATABASE|SCHEMA [IF NOT EXISTS] db_name [CHARACTER SET=] [COLLATE=]`

* `CHARACTER SET` 字符集
* `COLLATE` 排序规则

##### 修改数据库

`ALTER DATABASE|SCHEMA db_name [CHARACTER SET=] [COLLATE=]`

* 添加,删除,修改字段
* 添加.删除.修改索引
* 改表名
* 修改表属性

修改表引擎 = 创建一张新表, 把旧表删掉.

##### 删除数据库

`DROP DATABASE|SCHEMA [IF NOT EXISTS] db_name`

#### 表

##### 创建表

1. 直接定义一张空表
2. 从其他表中查询出数据, 并以之创建新表
3. 以其他表为模板创建一个空表

###### 直接定义一张空表

`CREATE TABLE [IF NOT EXISTS] tb_name (col_name col_defination,constraint)`

创建表以后也可以单独创建索引. `CREATE INDEX`

查看表的状态`SHOW TABLE STATUS LIKE 'tb_name';`

删除表 `DROP TABLE tb_name;`

查看表的索引: `SHOW INDEXES FROM tb_name;`

* Table: 哪张表的索引
* Non_unique: 是不是非唯一键
* Key_name: 键名称
* Seq_in_index: 该字段在索引中的位置，单列索引该值为1，复合索引为每个字段在索引定义中的顺序
* Column_name: 定义索引的列字段
* Collation: 排序规则,有值'A'(升序)或NULL(无分类)
* Cardinality: 索引中唯一值的数目的估计值.通过运行ANALYZE TABLE或myisamchk -a可以更新。基数根据被存储为整数的统计数据来计数,所以即使对于小型表,该值也没有必要是精确的.基数越大,当进行联合时,mysql使用该索引的机会就越大.
* Sub_part: 如果列只是被部分地编入索引,则为被编入索引的字符的数目.如果整列被编入索引,则为NULL.
* Packed: 指示关键字如何被压缩.如果没有被压缩,则为NULL.
* Null: 如果列含有NULL,则含有YES. 如果没有,则该列含有NO.
* Index_type: 索引类型(BTREE, FULLTEXT, HASH, RTREE)
* Comment:
* Index_comment:

```sql
mysql> SHOW INDEXES FROM pcore_user;
+------------+------------+-----------+--------------+-------------+-----------+-------------+----------+--------+------+------------+---------+---------------+
| Table      | Non_unique | Key_name  | Seq_in_index | Column_name | Collation | Cardinality | Sub_part | Packed | Null | Index_type | Comment | Index_comment |
+------------+------------+-----------+--------------+-------------+-----------+-------------+----------+--------+------+------------+---------+---------------+
| pcore_user |          0 | PRIMARY   |            1 | user_id     | A         |           0 |     NULL | NULL   |      | BTREE      |         |               |
| pcore_user |          0 | user_name |            1 | user_name   | A         |           0 |     NULL | NULL   |      | BTREE      |         |               |
| pcore_user |          0 | cellphone |            1 | cellphone   | A         |           0 |     NULL | NULL   | YES  | BTREE      |         |               |
+------------+------------+-----------+--------------+-------------+-----------+-------------+----------+--------+------+------------+---------+---------------+
3 rows in set (0.00 sec)
```

创建索引:

* PRIMARY KEY
* UNIQUE KEY

单或多字段:

* PRIMARY KEY (col1,...)
* UNIQUE KEY (col1,...)
* INDEX (col1,...)


```sql
CREATE TABLE tb1 (id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY, Name CHAR(20) NOT NULL, Age TINYINT NOT NULL) ENGINE [=] engine_name

CREATE TABLE tb1 (id INT UNSIGNED NOT NULL AUTO_INCREMENT,  Name CHAR(20) NOT NULL, Age TINYINT NOT NULL, PRIMAEY KEY(id), UNIQUE KEY(name), INDEX(age))
```

###### 从其他表中查询出数据, 并以之创建新表

`CREATE TABLE 创建表名称 SELECT * FROM 来源表名称 WHERE 查询数据条件`

字段的属性格式会不存在.

###### 以其他表为模板创建一个空表

`CREATE TABLE 创建表名称 LIKE 来源表名称`

表结构定义完全一样.

##### 表选项

* ENGINE: 表的存储引擎, 如果不指定从数据库直接继承
* AUTO_INCREMENT
* AVG_ROW_LENGTH: 平均每行包括的字节数
* CHARACTER SET: 字符集
* CHECKSUM: 是否启用校验和, 每次更新都会重新计算
* DELAT_KEY_WRITE
* COLLATE: 排序
* COMMENT: 注释
* DATA DIRECTORY: 数据目录
* MAX_ROWS: 最多允许存储多少行
* ROW_FORMAT
* TABLESPACE
* ...

##### 索引类型

* BTREE
* HASH

##### 唯一键,主键和索引的区别

键也称作约束可用作索引. 是特殊的索引, B+Tree索引结构:

* unique 数据不能相同
* 主键 数据不能相同, 不能为空

##### 删除表

InnoDB支持外键.

`DROP TABLE`

删除数据和表, 如果表里有外键需要先解决依赖关系.

外键会消耗系统资源.

#### 创建索引

`CREATE INDEX index_name ON tb_name (col,...)`

索引指定的`col`上还可以指定长度和排序.

`col_name[(length)][ASC | DESC]`

* `length` 索引的长度, 从最左侧开始只索引多少个字符, 比较越长资源消耗越大, 太短重复率越高. 索引使用的时候会全部载入内存. 只有字符串类型的字段才能指定索引长度.
* `ASC | DESC` 指定升序或降序的索引值存储

#### 删除索引.

`DROP INDEX index_name ON tb_name`

## 整理知识点

---