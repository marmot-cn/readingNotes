# 33_01_MySQL系列之五——MySQL数据类型及sql模型

---

## 笔记

---

### 存储引擎

也被称为表类型.

* MyISAM
	* 不支持事务
	* 支持表锁
	* `.frm` 表结构定义文件
	* `.MYD` 表数据文件
	* `.MYI` 索引
* InnoDB
	* 支持事务
	* 支持行锁 
	* `.frm` 表结构定义文件
	* `.ibd` 表空间(数据和索引)

`mysql`库中的表是`myisam`引擎.

#### 命令

`SHOW ENGINES` 查看存储引擎.

`SHOW TABLE STATUS [LIKE...]` 查看表的状态(包括存储引擎).

#### 服务端

* 客户端: mysql,mysqladmin,mysqldump,mysqlimport,mysqlcheck
* 服务器: mysqld, mysqld_safe, mysqld_multi(多实例, 在一个物理机运行多个实例, 但是绑定多个不同的端口)

默认启动是`mysqld_safe`.

`my.cnf`中的`[client]`是对所有客户端生效. 其中大部分参数, 都可以在命令行中直接指定.

```
使用mysqld --help --verbose 查看mysql支持的选项.

root@90af5cf2869e:/# mysqld --help --verbose
..
```

#### DBA

* 开发DBA:
	* 数据库设计
	* SQL语句
	* 存储过程
	* 存储函数
	* 触发器
* 管理DBA:
	* 安装
	* 升级
	* 备份
	* 恢复
	* 用户管理
	* 权限管理
	* 监控
	* 性能分析
	* 基准测试

#### 数据类型

* 数值型
	* 精确数值
		* int
		* bit
	* 近似数值
		* float(g,f) g:最大多少位 f:浮点多少位 1.36 g=3(整体为3位) f=2(小数点后2位)
		* double
		* decimal
		* real(实数)
* 字符型
	* 定长
		* char(不区分大小写)
		* binary(区分大小写)
	* 变长
		* varchar(不区分大小写)
		* varbinary(区分大小写)
	* text(不区分大小写)
	* blob(区分大小写)
* 内置型(可以理解为字符型)
	* enum 枚举，性别:M,F 只能填写M或F. 包含1-65,535个字符串(不同变化).
	* set 集合,随意组合字符, 定义了两个字符M,F 可以填写M,F,MF,FM(FF,MM一般不可以). 1-64个字符串. 在表中存储的是集合中的索引下标.
* 日期时间型
	* date : 3字节 1000-01-01 to 9999-12-31
	* time : 3字节 -838:59:59 to 838:59:58
	* datetime: 8字节 date 和 time 组合起来
	* timestamp: 4字节 1970-01-01 00:00:00 to 2038-01-18 22:14:07
	* year
		* 4位: year(4) 1个字节 1901-2155
		* 2位: year(2) 1个字节 00-99
* `bool`型, `MySQL`没有布尔型, 可以理解为数字型.

定义了数据类型会做出以下限制:

1. 表示了哪个种类的数据
2. 最大存储空间确定
3. 定义了是变长还是定长的
4. 定义了比较(compare)和排序(sort)
5. 是否能够索引

#### `int`

* tinyint 1个字节
* smallint 2个字节
* mediumint 3个字节
* int	4个字节
* bigint 8个字节

##### `AUTO_INCREMENT`

自动增长. 没增加一个自动加N.

满足两个条件:

* 必须使用整数类型
* 只能包含正数, 不能为负数, 无符号
* 不能为空
* 一定要索引
	* 主键
	* 唯一

函数`LAST_INSERT_ID()`. 显示出上一次自动生成自增序列(`AUTO_INCREMENT`)的id.

```sql
SELECT LAST_INSERT_ID();
```

#### 字符型

* char(M) 存储`M`个字符, `M`最大为`255`
* varchar(M) 最多存储`M`个字符, 但是占据`M+1(1个字节)`(如果在255需要一个结束符, 超过255需要2个结束符)个空间, 这额外的一个是结束符, 因为是变长的, 不知道在哪里结束. 最多为 65535个字符.
* tinytext 最大255个字符, 额外占据一个字节.
* text 最大65545个字符, 额外占据2个字节.
* mediumtext 最大16777215个字符. 额外占据3个字节.
* longtext 最大4294967295个字符. 额外占据4个字节.

没`2^8`次方字符都额外多占据一个字节.

`char(255)`和`tinytext`的区别, 虽然最大都是255个字符: char(255)可以索引, tinytext不能索引.

##### 字符型常用修饰符

* `NOT NULL` 不允许为空
* `NULL` 可以为空
* `DEFAULT` 默认值
* `BINARY` 大小写
* `CHARACTER SET` 字符集
* `COLLATION` 排序规则

如果字段没有设定排序规则, 从表继承, 表没有从数据库继承, 数据库没有从服务器继承.

```shell
显示当前服务器所支持的所有字符集
SHOW CHARACTER SET; 

显示当前服务器各个字符集下的排序规则, 同一个字符集下面有多种不同的排序规则.
SHOW COLLATION;
```

#### 域属性, 修饰符

修饰符定义域的限制.

#### blob 无法建立索引

大数据段不是存在行和列, 而是存储在其他地方, 行中字段存储的是指针, 所以指针无法创建索引.

#### 函数和存储过程

所有函数使用`SELECT`执行.

所有存储过程使用`CALL`执行.

### SQL模型

常见的SQL模型:

* ANSI QUOTES: 双引号相当于反引号(表名,字段名称), 单引号用来表示字符串.
* IGNORE_SPACE: 内建函数中忽略多余的空白字符.
* STRICT_ALL_TABLES: 如果没有设置, 所有非法数据都允许. 如果设置, 不允许非法数据, 并返回错误.
* STRICT_TRANS_TABLES: 像支持事务的表中提供非法数据不允许, 并且返回错误.
* TRADITIONAL:

```sql
SHOW GLOBAL VARIABLES LIKE 'sql_mode';
```

### mysql服务器变量

按作用域分为:

* 全局变量, 和用户没有关系, 服务器启动就生效. `SHOW GLOBAL VARIABLES`.
* 会话变量, 只对当前会话生效. `SHOW [SESSION] VARIABLES`.

按生效时间划分:

* 动态调整, 可以即时调整.
* 静态, 写到配置文件中或通过参数传递给`mysqld`

动态调整参数的生效方式:

* 全局: 对当前会话无效, 只对新建立会话有效.
* 会话: 即时生效, 只对当前会话有效.

显示变量:

`SELECT @@global.sql_mode`

* `@@` 服务器变量或内置变量.
* `@` 用户自定义变量.

设定变量:

`SET GLOBAL|SESSION 变量名='value'`

## 整理知识点

---