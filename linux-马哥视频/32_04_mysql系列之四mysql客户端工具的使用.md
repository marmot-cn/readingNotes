# 32_04_mysql系列之四mysql客户端工具的使用

---

## 笔记

---

### `mysql`客户端

工作模式:

* 交互模式`mysql>`
* 批处理模式(脚本模式)`mysql < xxx.sql(sql语句的文件)`

交互式命令:

`mysql>`两类命令

* 客户端命令(`\?`可以查看命令列表)
	* `source`or(`\.`): 可以读取一个`sql`文件

			mysql> \.	xxx.sql
			mysql> source xxx.sql
	* `delimiter`or(`\d`): 定义语句结束符(默认为`;`)
	* `status`or(`\s`)
	* `clear`or(`\c`) 提前终止语句执行

			mysql> show databas \c
			mysql>
	* `connect`or(`\r`) 重新连接服务器.

			mysql> connect
			Connection id:    2
			Current database: *** NONE ***		
	* `ego`or(`\G`) 无论语句结束符是什么, 直接将此语句送至服务器端执行, 而且结果以竖排方式显示.
	* `go`or(`\g`) 无论语句结束符是什么, 直接将此语句送至服务器端执行.
	* `print`or(`\p`) 显示当前执行的命令.
	* `quit`or(`\q`) 退出。
	* `system`or(`\!`) 执行`shell`命令.

			mysql> \! ls
	* `warnings`or(`\W`) 语句执行结束后显示警告信息.
	* `nowarning`or(`\w`) 语句执行结束后不显示警告信息.
	* `rehash`or(`\#`) 对新建的对象支持补全功能.
			
* 服务器语句: 有语句结束符, 默认为`;`.

#### mysql中的promots(提示符)

* `mysql>` 准备好输入新的语句.
* `->` 还没输入完, 等待新的语句.
* `'>` 等待还没输入完成的单引号(单引号没有成对出现)
* `">` 等到还没输入完成的双引号(双引号没有成对出现)
* \`> 等待还没输入完成的` 缺少反引号的后一半.
* `/*>`缺少多行注释的后一半.

#### mysql的其他选项

mysql:

* `--compress` 压缩所有语句从客户端到服务端.
* `--ssl-ca` `ca`证书文件.
* `--ssl-capath` 如果有多个`ca`, ca证书的目录, 和上面的选项出现一个.
* `--ssl-cert` 自己的证书.
* `--ssl-cipher` 加密算法列表.
* `--ssl-key` 自己的密钥文件.
* `--ssl-verify-server-cert` 是否验证服务器端证书.
		
#### mysql名称补全

需要载入到内存才能补全, 连接到数据库需要载入会非常慢.

禁用补全: 

* `-A`
* `--no-auto-rehash`
* `--disable-auto-rehash`
		
对新建的功能支持补全:

* `\#`
* `rehash`		
		
#### mysql的输出格式

* `--html`or(`-H`)以`html`格式输出.
* `--xml`or(`-X`)以`xml`格式输出.

```shell
root@90af5cf2869e:/# mysql -uroot -p123456 --html
Warning: Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 3
Server version: 5.6.37 MySQL Community Server (GPL)

Copyright (c) 2000, 2017, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> use marmot;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> show tables;
<TABLE BORDER=1><TR><TH>Tables_in_marmot</TH></TR><TR><TD>pcore_user</TD></TR></TABLE>1 row in set (0.00 sec)

mysql>
```
		
#### 交互式客户端命令示例

```shell

mysql> \?

For information about MySQL products and services, visit:
   http://www.mysql.com/
For developer information, including the MySQL Reference Manual, visit:
   http://dev.mysql.com/
To buy MySQL Enterprise support, training, or other products, visit:
   https://shop.mysql.com/

List of all MySQL commands:
Note that all text commands must be first on line and end with ';'
?         (\?) Synonym for `help'.
clear     (\c) Clear the current input statement.
connect   (\r) Reconnect to the server. Optional arguments are db and host.
delimiter (\d) Set statement delimiter.
edit      (\e) Edit command with $EDITOR.
ego       (\G) Send command to mysql server, display result vertically.
exit      (\q) Exit mysql. Same as quit.
go        (\g) Send command to mysql server.
help      (\h) Display this help.
nopager   (\n) Disable pager, print to stdout.
notee     (\t) Don't write into outfile.
pager     (\P) Set PAGER [to_pager]. Print the query results via PAGER.
print     (\p) Print current command.
prompt    (\R) Change your mysql prompt.
quit      (\q) Quit mysql.
rehash    (\#) Rebuild completion hash.
source    (\.) Execute an SQL script file. Takes a file name as an argument.
status    (\s) Get status information from the server.
system    (\!) Execute a system shell command.
tee       (\T) Set outfile [to_outfile]. Append everything into given outfile.
use       (\u) Use another database. Takes database name as argument.
charset   (\C) Switch to another charset. Might be needed for processing binlog with multi-byte charsets.
warnings  (\W) Show warnings after every statement.
nowarning (\w) Don't show warnings after every statement.

For server side help, type 'help contents'
```	

#### mysql服务器端命令

`help keyword` 可以显示`keyword`的帮助信息.

```shell
mysql> help SELECT;
Name: 'SELECT'
Description:
Syntax:
SELECT
    [ALL | DISTINCT | DISTINCTROW ]
      [HIGH_PRIORITY]
      [STRAIGHT_JOIN]
      [SQL_SMALL_RESULT] [SQL_BIG_RESULT] [SQL_BUFFER_RESULT]
      [SQL_CACHE | SQL_NO_CACHE] [SQL_CALC_FOUND_ROWS]
    select_expr [, select_expr ...]
    [FROM table_references
      [PARTITION partition_list]
    [WHERE where_condition]
    [GROUP BY {col_name | expr | position}
      [ASC | DESC], ... [WITH ROLLUP]]
    [HAVING where_condition]
    [ORDER BY {col_name | expr | position}
      [ASC | DESC], ...]
    [LIMIT {[offset,] row_count | row_count OFFSET offset}]
    [PROCEDURE procedure_name(argument_list)]
    [INTO OUTFILE 'file_name'
        [CHARACTER SET charset_name]
        export_options
      | INTO DUMPFILE 'file_name'
      | INTO var_name [, var_name]]
    [FOR UPDATE | LOCK IN SHARE MODE]]
```

#### mysqladmin 客户端命令

用来执行管理命令.

`mysqladmin [options] command [arg] [command [arg]]`

`mysqladmin -u USERNAME -h HOSTNAME password` 设定用户密码的命令, 其中`password`就是一个子命令.

* `create DATABASE` 创建数据库

```sql
root@90af5cf2869e:/# mysqladmin -uroot -p123456 create hellodb;
Warning: Using a password on the command line interface can be insecure.
```

* `drop DATABASE ` 删除数据库
* `password` 修改命令
* `ping` 检查对方数据库是否在线

		root@90af5cf2869e:/# mysqladmin -uroot -p ping
		Enter password:
		mysqld is alive
* `processlist` 进程列表.

		root@90af5cf2869e:/# mysqladmin -uroot -p123456 processlist
		Warning: Using a password on the command line interface can be insecure.
		+----+------+-----------+----+---------+------+-------+------------------+
		| Id | User | Host      | db | Command | Time | State | Info             |
		+----+------+-----------+----+---------+------+-------+------------------+
		| 11 | root | localhost |    | Query   | 0    | init  | show processlist |
		+----+------+-----------+----+---------+------+-------+------------------+
* `status` 状态

		root@90af5cf2869e:/# mysqladmin -uroot -p123456 status
		Warning: Using a password on the command line interface can be insecure.
		Uptime(开机时间): 31415  Threads: 1  Questions: 42  Slow queries: 0  Opens: 72  Flush tables: 1  Open tables: 65  Queries per second avg: 0.001
		
	* `--sleep` 睡眠多久显示一次.
	* `--count` 一共显示多少次.
* `extended-status` 显示状态变量(服务器工作状态的统计信息).
* `variables` 服务器变量(定义服务器的工作属性).
* `flush-privileges` 刷新授权表.
* `flush-tables` 关闭所有当前打开的表的文件句柄(处理), 关闭所有已打开的表.
* `flush-threads` 重置线程池(线程缓存), 把空闲线程清楚.
* `flush-status` 重置大多数服务器状态变量.
* `flush-logs` 二进制和中继日志滚动. 对错误日志只是关闭在打开.
* `flush-hosts` 刷新主机, `mysql`在单位时间内限制一个用户的失败链接次数. 重置计数器. 清楚主机内部信息, DNS的解析缓存和由于太多链接错误导致用户拒绝登录.
* `kill` 杀死一个客户端线程. 可以用逗号分隔线程id(一次多个线程)
* `reload` 等同于`flush-privileges`.
* `refresh` 等同于`flush-logs`和`flush-hosts`同时执行.
* `shutdown` 关闭`mysql`服务器进程.
* `start-slave` 启动复制, 启动从服务器复制线程.
	* `SQL thread` 线程.
	* `IO thread` 线程.
* `stop-slave` 关闭复制, 关闭从服务器复制线程.
		
## 整理知识点

---