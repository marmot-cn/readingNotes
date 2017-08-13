# 32_03_mysql系列之三MySQL数据库基础及编译安装

---

## 笔记

---

### MySQL 特性

* Speed: 高性能
	* 完全多线程(单进程?)
	* 查询缓存
* Reliability: 稳定性
* Scalability: 伸缩性, 数据分区
* Ease of use: 易用
* 多用户支持
* 开源

### MySQL 产品

* MySQL Server(mysqld:服务器, mysql:客户端)
* MySQL Cluster: mysql集群套件, 将多个物理服务器组合起来高可用,负载均衡.
* MySQL Proxy: mysql代理服务器, 实现语句路由. 实现读写分离.
* MySQL Administrator: mysql图形化管理工具.
* MySQL Query Browser: mysql查询浏览器.
* MySQL Workbench: mysql数据库设计工具.
* MySQL Migration Toolkit: mysql移植工具箱(access->mysql, sqlserver->mysql, mysql不同版本,不同平台(win-linux)).
* MySQL Embedded Server: 嵌入式.
* MySQL Drivers and Connextors: mysql驱动和连接器.

`MySQL`社区办和商业版.

#### percona

www.percona.com

高性能`mysql`博客组.

### 安装 MySQL

1. 专用软件包管理器包
	* deb
	* rpm(要根据OS选取不同的系列)
		* RHEL红包, 包括Oracle Linux, CentOS
		* SUSE
2. 通用二进制格式包
	* gcc: x86(32位), x64(64位)
	* icc 
3. 源代码包(自己编译安装)
	* 5.1: `make`编译
	* 5.5+: 使用`cmake(跨平台编译器)`编译 

#### MySQL 官方 rpm 包

* `MySQL-client`(我们需要装): 客户端, 以及其他客户端组件.
* `MySQL-debuginfo`: 调试功能, 针对客户端和服务端.
* `MySQL-devel`(我们需要装, 如果有需要): 开发组件, 开发时用到的头文件和库文件. 如果编译一个软件, 依赖于`MySQL`, 就必须把该包装上.
* `MySQL-embedded`: 嵌入式环境专用
* `MySQL-ndb-management`: Mysql Cluster server `name db`专用
* `MySQL-server`(我们需要装): 服务器端
* `MySQL-shared`(我们需要装): 共享库
* `MySQL-shared-compat`(我们需要装): 兼容老版本客户端工具, 兼容库. 是`shared`补充.
* `MySQL-test`: 测试组件
* `MySQL-VERSION-PLATFORM.src.rpm`: 源码格式的`rpm`包

#### MySQL安装后的目录结构

* `bin`: 二进制程序
* `data`: 数据目录
* `include`: 头文件
* `lib`: 库文件
* `man`: 手册
* `mysql-test`: 测试组件
* `scripts`: `MySQL`初始化脚本.
* `share`: 对于各种不同语言出现错误信息的各种语言版本.
* `sql-bench`: 基准性能测试(测试针对当前主机的响应,负载能力).
* `support-files`: 启动服务脚本, 配置文件样例等等.

#### MySQL完成初始化配置

中心`my.cnf`配置文件.

```
[mysql] 针对客户端生效

[mysqld] 针对服务器端生效

...
```

##### 找寻配置文件路径以及顺序

1. `/etc/my.cnf`
2. `/etc/mysql/my.cnf`
3. `$MYSQL_HOME/my.cnf`
	* `$MYSQL_HOME`:`MySQL`每个实例(监听不同的端口)自己的家目录
4. `/path/to/file when defaults-extra-file=/path/to/file is specified`启动时手动指定
5. `~/.my.cnf`: `mysql`用户家目录下.

每个文件都会定位读取一次, 如果多个文件都找见了, 需要合并起来. 如果信息重复了, 最后读取的配置将会是最终生效的结果. 上面配置文件的生效顺序, 优先级由低到高.

#### 安装完成后的操作

mysql 用户名包含两部分:

* 用户名
* 授权的主机

##### 默认安装完成后会生成5个用户

* 3个`root`
	* `root@127.0.0.1`
	* `root@localhost`
	* `root@hostname(当前主机名)`
* 2个匿名, 匿名用户不安全, 要删除匿名用户
	* `@localhost`
	* `@hostname`

##### 删除匿名用户

```
musql> DROP USER ''@localhost(主机名)
```

##### 给`root`用户设置密码

`root`用户默认密码为空.

`MySQL`用户密码修改:

1. `mysqladmin -u USERNAME -h HOSTNAME password 'NEW_PASS' -p(此前用户有密码需要指定)`

	```shell
	mysqladmin -u root password 'new-password' -p(指定老密码, 指定老密码才能改新密码)
	mysqladmin -u root -h this_host_name password 'new-password' -p
	```
2. `mysql> SET PASSWORD FOR 'USERNAME'@'HOST'=PASSWORD('new_pass');`
3. `mysql> UPDATE mysql.user SET PASSWORD=PASSWORD('new_pass') WHERE condition`(需要 flush privileges 授权表)

##### 修改客户端配置文件

```shell
家目录下 .my.cnf 定义

[clinet]
user='root'
password='xxxxx'
host='localhsot'
```

本地直接调用`mysql`命令可以直接加载该配置文件.

### MySQL安装:

源码安装`MySQL`.

`cmake`.

`MySQL`数据应该放在独立的目录下(`lv`下,数据的属主数组属于`mysql`用户).

#### 使用`cmake`编译`mysql5.5`

`cmake`指定编译选项的方式不同于`make`, 其实现方式对比如下:

```
./configure				cmake .
./configure --help		cmake . -LH or ccmake .
```

指定安装文件的安装路径时常用的选项:

* `DCMAKE_INSTALL_PREFIC=/usr/local/mysql`
* `DMYSQL_DATADIR=/data/mysql`
* `DSYSCONFIGDIR=/etc`

默认编译的存储引擎包括: `csv`,`myisam`,`myisammrg`和`heap`. 若要安装其他存储引擎, 可以使用类似如下编译选项:

* `-DWITH_INNODBASE_STORAGE_ENGINE=1`(innodb)
* `-DWITH_ARCHIVE_StorAGE_ENGINE=1`
* `-DWITH_BLOACKHOLE_STORAGE_ENFINE=1`(类似/dev/null, 中继复制有用)
* `-DWITH_FEDERATED_STORAGE_ENGINE=1`

`5.5`默认没有包含`innodb`, `5.6`默认包含了.

若要明确指定不编译某存储引擎, 可以使用类似如下的选项:

* `-DWITHOUT_<ENGINE>_STORAGE_ENGINE=1`
	* `-DWITHOUT_EXAMPLE_STORAGE_ENGINE=1`
	* `-DWITHOUT_FEDERATED_STORAGE_ENGINE=1`
	* `-DWITHOUT_PARTION_STORAGE_ENGINE=1`

如若要编译进其他功能, 比如`ssl`等, 则可使用类似如下选项来实现编译时使用某库或者不使用某库.

* `-DWITH_READLINE=1` 能够使用`load in file`批量导入mysql数据.
* `-DWITH_SSL=system` mysql支持基于ssl会话.
* `-DWITH_ZLIB=system` 压缩库.
* `-DWITH_LIBWRAP=0` 针对是否可以使用`tcpwrapper`来进行访问控制.

其他常用的选项:

* `-DMYSQL_TCP_PORT=3306`
* `-DMYSQL_UNIX_ADDR=/tmp/mysql.sock` (unix 套接字路径)
* `-DENABLED_LOCAL_INFILE=1` load in file
* `-DEXTRA_CHARSETS=all`
* `-DDEFAULT_CHARSET=utf8`
* `-DDEFAULT_COLLATION=utf8_general_ci`
* `-DWITH_DEBUG=0`
* `-DEANBLE_PROFILING=1` profiling 性能分析.

#### 编译安装

```shell
groupadd -r mysql (-r 表示建立系统账户)
useradd -g mysql -r -d /data/mydata mysql (-d 指定mysql用户的家目录, 替换系统的默认/home/<用户名>, 如果该某不存在可以使用 -m 自动创建) -s /sbin/noligin mysql
tar xf mysql-5.5.25a.tar.gz
cd mysql-5.5.25a
cmake . -DCMAKE_INSTALL_PREFIC=/usr/local/mysql
		-DMYSQL_DATADIR=/data/mysql
		-DSYSCONFIGDIR=/etc
		-DWITH_INNODBASE_STORAGE_ENGINE=1
		-DWITH_ARCHIVE_StorAGE_ENGINE=1
		-DWITH_BLOACKHOLE_STORAGE_ENFINE=1
		-DWITH_READLINE=1
		-DWITH_SSL=system
		-DWITH_ZLIB=system
		-DWITH_LIBWRAP=0
		-DMYSQL_UNIX_ADDR=/tmp/mysql.sock
		-DDEFAULT_CHARSET=utf8
		-DDEFAULT_COLLATION=utf8_general_ci
		
make 
make install
```

#### 访问`mysql`

客户端`mysql`访问`mysqld`

* unix 在同一台主机上
	* mysql -> mysql.sock --> mysqld, 通过mysql.sock完成进程间通信
* windows 在同一台主机上
	* mysql -> memory(共享内存) or pipe(管道) -> mysqld 
* 客户端和服务器端不在同一台主机上
	* 基于`TCP/IP协议`

```
mysql 默认用户管理员 默认密码为空, 默认本机, 使用 mysql.sock 通信

mysql -u root -h xxxx(xxx是本机ip), 使用 tcp/ip 协议通信
```

#### `mysql`客户端程序

* mysqlimport
* mysqldump
* mysql
* mysqlcheck

配置文件中`[client]`都会生效.

链接到服务端支持的选项:

* -u(--user) USERNAME
* -h(--host) HOST
* -p(--password) password
* --protocol 指定协议
	* `tcp`(all)
	* `socket`(unix only), 默认使用
	* `pipe`(windows)
	* `memory`(windows)
* --port 指定端口(如果使用mysql.sock, 指定端口没有意义)
* `-D(--database)`设定默认库

#### `mysql`非客户端工具

* `myisamck` 检查myisam表
* `myisampack` 压缩myisam表

#### 数据库数据目录文件格式

Myasiam 引擎:

* `frm`: 表结构定义文件
* `MYD`: 表数据文件
* `MYI`: 索引文件

InnoDB 引擎:

* 所有表共享一个表空间文件
* 建议: 每表一个独立的表空间文件
	* `.frm`: 表结构定义文件
	* `.idb`: 表空间(数据和索引)

```
SHOW VARIABLES LIKE '%innodb%'
...
innodb_file_per_table 如果是 ON 就是每一个表一个数据空间文件, 是OFF的话每所有表放在一个数据空间文件(我自己本机的mysql默认这个选项也是打开的).
...

如果要永久生效:

vim my.cnf

[mysqld]
innodb_file_per_table = 1
```

## 整理知识点

---

### `tcpwrapper`

### 字符集(`charset`)

计算机存储`01`代码. 比如中文: 哪个编码表示哪个汉字, 就是字符集.

### 排序规则(`collation`)

同一种字符集的字符也有不同的排序规则. 要和字符集语言对应.

### `partprobe`

在硬盘分区发生改变时, 更新`Linux`内核中读取的硬盘分区数据表.