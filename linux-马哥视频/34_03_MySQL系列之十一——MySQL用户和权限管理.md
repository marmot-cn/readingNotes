# 34_03_MySQL系列之十一——MySQL用户和权限管理

---

## 笔记

---

### 用户

用于认证.

用户是虚拟用户, 和操作系统无关.

密码默认使用`password()`加密函数.

用户: `用户名@主机`, 用户必须通过对应的主机才能登陆系统.

用户是保存在`mysql`库的多个表中实现的.

* `user`
	* 用户账户
	* 全局权限
	* 非权限字段
* `db`
	* 数据库级别权限定义
* `host`
	* 已经废弃
* `tables_priv`
	* 表级别权限
* `columns_priv`
	* 列级别权限
* `procs_priv`
	* 存储过程和存储函数相关权限
* `proxies_priv` 代理用户权限, 将mysql访问给某个代理服务器来验证的时候.

`mysql`启动时会读取这六张表, 并在内存中生成授权表. 任何`sql`语句都可能需要检查权限, 为了加速所以放在内存处理.`FLUSH PRIVILEGES;`重读授权表, 并生效权限. 有些时候需要重新建立会话, 以读取最新的权限.

#### 用户账户

用户账户由**用户名**@**主机**共同构成.

主机:

* 主机名: www.xxx.com, mysql
* ip: 172.16.10.177
* 网络地址: 172.16.0.0/255.255.0.0
* 通配符:
	* `%`: 表任何字符串. 172.16.%.%(172.16.0.0/255.255.0.0), %.xxx.com(域内所有主机)
	* `_`: 任意一个字符.

为了验证用户来源主机是否合法, 一般mysql主机需要反解或者解析用户主机名. `--skip-name-resolve`略过名称解析, 可以提升用户链接速度. 

#### 权限级别

* 全局级别
* 库
* 表: `DELETE`,`ALTER`,`SELECT`,`INSERT`,`UPDATE`,`TRIGGER`,`DROP`,`INDEX`,`CREATE`
* 列: `SELECT`,`INSERT`,`UPDATE`
* 存储过程和存储函数

不同级别权限支撑的命令不一致. 其他可以参见`mysql`手册.

* CREATE(Create_priv),应用场景
	* databases 库
	* tables 表
	* indexes 索引
....

场景是`server administration`, 是服务器级别的管理权限.

* `ALTER_ROUTINE`: 修改存储过程和函数
* `CREATE_ROUTINE`: 创建存储过程和函数
* `EXECUTE`: 执行存储过程和函数
* `CREATE_TABLESPACE`: 创建表空间(就是`innodb`老版是默认一个共享空间,可以通过参数修改(后来是默认的)一个表一个空间, 空间是一个独立的文件存储数据)
* `CREATE_USER`: 创建用户
* `PROCESS`: 查看进程列表
* `PROXY`: 代理服务,代理用户的创建
* `RELOAD`: 重载授权表(`flush privileges`)
* `REPLICATION_CLINET`: 确定用户是否可以确定复制从服务器和主服务器的位置
* `REPLICATION_SLAVE`: 确定用户是否可以读取用于维护复制数据库环境的二进制日志文件。此用户位于主系统中，有利于主机和客户机之间的通信
* `SHOW_DATABASES`
* `SHUTDOWN`: 关闭服务
* `SUPER`: 服务器级别管理权限, 确定用户是否可以执行某些强大的管理功能，例如通过KILL命令删除用户进程，使用SET GLOBAL修改全局MySQL变量，执行关于复制和日志的各种命令, `CHANGE MAGETER TO`(设定主从中的主服务器).
* `ALL [PRIVILEGES]`
* `USAGE`, 没有任何权限, 仅仅可以连接到`mysql`服务器.
* `FILE`, 确定用户是否可以执行`SELECT INTO OUTFILE`(表中的数据备份到文件)和`LOAD DATA INFILE`(从文件中导入数据到表)命令.
* ...

##### 库级别权限

```shell
mysql> select * from db \G
*************************** 1. row ***************************
                 Host: %
                   Db: test
                 User:
          Select_priv: Y
          Insert_priv: Y
          Update_priv: Y
          Delete_priv: Y
          Create_priv: Y
            Drop_priv: Y
           Grant_priv: N
      References_priv: Y
           Index_priv: Y
           Alter_priv: Y
Create_tmp_table_priv: Y
     Lock_tables_priv: Y
     Create_view_priv: Y
       Show_view_priv: Y
  Create_routine_priv: Y
   Alter_routine_priv: N
         Execute_priv: N
           Event_priv: Y
         Trigger_priv: Y
*************************** 2. row ***************************
                 Host: %
                   Db: test\_%
                 User:
          Select_priv: Y
          Insert_priv: Y
          Update_priv: Y
          Delete_priv: Y
          Create_priv: Y
            Drop_priv: Y
           Grant_priv: N
      References_priv: Y
           Index_priv: Y
           Alter_priv: Y
Create_tmp_table_priv: Y
     Lock_tables_priv: Y
     Create_view_priv: Y
       Show_view_priv: Y
  Create_routine_priv: Y
   Alter_routine_priv: N
         Execute_priv: N
           Event_priv: Y
         Trigger_priv: Y
2 rows in set (0.00 sec)
```

只授权了`test`库和`test_%`库能被`%`(所有)主机访问.

```shell
添加一个用户对库的权限
mysql> GRANT ALL ON marmot.* TO 'marmot'@'10.132.27.109' IDENTIFIED BY '123456';
Query OK, 0 rows affected (0.01 sec)

再次查看
mysql> select * from db \G
*************************** 1. row ***************************
                 Host: %
                   Db: test
                 User:
          Select_priv: Y
          Insert_priv: Y
          Update_priv: Y
          Delete_priv: Y
          Create_priv: Y
            Drop_priv: Y
           Grant_priv: N
      References_priv: Y
           Index_priv: Y
           Alter_priv: Y
Create_tmp_table_priv: Y
     Lock_tables_priv: Y
     Create_view_priv: Y
       Show_view_priv: Y
  Create_routine_priv: Y
   Alter_routine_priv: N
         Execute_priv: N
           Event_priv: Y
         Trigger_priv: Y
*************************** 2. row ***************************
                 Host: %
                   Db: test\_%
                 User:
          Select_priv: Y
          Insert_priv: Y
          Update_priv: Y
          Delete_priv: Y
          Create_priv: Y
            Drop_priv: Y
           Grant_priv: N
      References_priv: Y
           Index_priv: Y
           Alter_priv: Y
Create_tmp_table_priv: Y
     Lock_tables_priv: Y
     Create_view_priv: Y
       Show_view_priv: Y
  Create_routine_priv: Y
   Alter_routine_priv: N
         Execute_priv: N
           Event_priv: Y
         Trigger_priv: Y
*************************** 3. row ***************************
                 Host: 10.132.27.109
                   Db: marmot
                 User: marmot
          Select_priv: Y
          Insert_priv: Y
          Update_priv: Y
          Delete_priv: Y
          Create_priv: Y
            Drop_priv: Y
           Grant_priv: N
      References_priv: Y
           Index_priv: Y
           Alter_priv: Y
Create_tmp_table_priv: Y
     Lock_tables_priv: Y
     Create_view_priv: Y
       Show_view_priv: Y
  Create_routine_priv: Y
   Alter_routine_priv: Y
         Execute_priv: Y
           Event_priv: Y
         Trigger_priv: Y
3 rows in set (0.00 sec)
```
##### 表级别权限

```shell
mysql>  GRANT ALL ON marmot.pcore_user TO 'marmot'@'10.132.27.109' IDENTIFIED BY '123456';
mysql> select * from tables_priv \G
*************************** 1. row ***************************
       Host: 10.132.27.109
         Db: marmot
       User: marmot
 Table_name: pcore_user
    Grantor: root@localhost
  Timestamp: 0000-00-00 00:00:00
 Table_priv: Select,Insert,Update,Delete,Create,Drop,References,Index,Alter,Create View,Show view,Trigger
Column_priv:
1 row in set (0.00 sec)
```

### 临时表

存储在内存中.

* 速度快
* 不需要永久保存
* 空间有限

### 触发器

执行`INSET`,`DELETE`,`UPDATE`可以自动触发一些我们自定义的操作.

### 创建用户

#### CREATE

`CREATE USER user_name%host [IDENTIFIED BY password]`

这样的用户只有:

* `USAGE`
* `SHOW_DATABASES`

权限.

自动触发`FLUSH PRIVILEGES;`.

#### GRANT

`SHOW GRANTES FOR 'user@host'` 查看某个用户的授权信息.

```
mysql> SHOW GRANTS FOR 'marmot'@'10.132.27.109';
+-------------------------------------------------------------------------------------------------------------------+
| Grants for marmot@10.132.27.109                                                                                   |
+-------------------------------------------------------------------------------------------------------------------+
| GRANT USAGE ON *.* TO 'marmot'@'10.132.27.109' IDENTIFIED BY PASSWORD '*6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9' |
| GRANT ALL PRIVILEGES ON `marmot`.* TO 'marmot'@'10.132.27.109'                                                    |
| GRANT ALL PRIVILEGES ON `marmot`.`pcore_user` TO 'marmot'@'10.132.27.109'                                         |
+-------------------------------------------------------------------------------------------------------------------+
3 rows in set (0.00 sec)
```

```sql
GRANT
    priv_type [(column_list)]
      [, priv_type [(column_list)]] ...
    ON [object_type] priv_level
    TO user_specification [, user_specification] ...
    [REQUIRE {NONE | tls_option [[AND] tls_option] ...}]
    [WITH {GRANT OPTION | resource_option} ...]

GRANT PROXY ON user_specification
    TO user_specification [, user_specification] ...
    [WITH GRANT OPTION]

object_type: {
    TABLE
  | FUNCTION
  | PROCEDURE
}

priv_level: {
    *
  | *.*
  | db_name.*
  | db_name.tbl_name
  | tbl_name
  | db_name.routine_name
}

user_specification:
    user [ auth_option ]

auth_option: {
    IDENTIFIED BY 'auth_string'
  | IDENTIFIED BY PASSWORD 'hash_string'
  | IDENTIFIED WITH auth_plugin
  | IDENTIFIED WITH auth_plugin AS 'hash_string'
}

tls_option: {
    SSL
  | X509
  | CIPHER 'cipher'
  | ISSUER 'issuer'
  | SUBJECT 'subject'
}

resource_option: {
  | MAX_QUERIES_PER_HOUR count
  | MAX_UPDATES_PER_HOUR count
  | MAX_CONNECTIONS_PER_HOUR count
  | MAX_USER_CONNECTIONS count
}
```

可以限定必须使用`ssl`(`REQUIRE ssl`). **必须**使用`ssl`链接, 如果没有限制则是可用或可不用.

限定用户链接进来的资源请求:

* `GRANT OPTION`: 被授权的用户可以将这些权限赋予给别的用户(危险), 可以逐级授权.
* `MAX_QUERIES_PER_HOUR count`: 每小时最多查询请求次数,`count`修改为`0`表示不限定.
* `MAX_UPDATES_PER_HOUR count`: 每小时最多更新次数请求,`count`修改为`0`表示不限定.
* `MAX_CONNECTIONS_PER_HOUT count`: 每小时只允许发起多少个新的链接请求,`count`修改为`0`表示不限定.
* `MAX_USER_CONNECTIONS count`: 某一个用户账户最多同时连进来多少次,`count`修改为`0`表示不限定.

##### 示例

在`db`的`abc`存储函数上有执行权限.
`GRANT EXECUTE ON FUNCTION db.abc TO 'username'@'%';`

只授权修改某张表中某个字段的权限.
`GRANT UPDATE(字段) ON db.table TO 'username'@'%';`

给予`SUPER`权限.
`GRANT SUPER ON *.* TO 'username'@'%';`


#### 主动插入`user`表

`INERT INTO mysql.user`

需要手动`FLUSH PRIVILEGES;`

### 删除用户

`DROP USER 'username'@'%';`

### 重命名用户

`REANME USER 'oldname'@'oldhost' TO 'newname'@'newhost';`

### 取消授权

`REVOKE 权限 ON *.* FROM  'newname'@'newhost';`

### 找回管理员密码

1. 关掉数据库服务.
2. 手动启动`mysql`(命令行启动, 修改服务脚本).
		
		编辑脚本 /etc/init.d/mysqld
		
		在start相中 mysqld_safe 后面添加参数 --skip-grant-tables(跳过授权表) --skip-networking(跳过网络, 防止有人通过网络尝试连接进去)
		
3. 启动后可以直接通过`mysql`命令进入.
4. 修改管理员密码.
	* `SET PASSWORD FOR 'root'@'host'=PASSWORD('xxx');`, 因为跳过授权表,这样**不能修改**.
	* `UPDATE user SET Password=PASSWORD('xxx') WHERE User='root';`
5. 退出, 还原第`2`步编辑的脚本.

## 整理知识点

---