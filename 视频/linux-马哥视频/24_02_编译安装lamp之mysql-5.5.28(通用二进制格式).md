# 24_02_编译安装LAMP之MySQL-5.5.28(通用二进制格式)

---

### 笔记

---

#### httpd 2.4新特性

* MPM 可以运行时装载
	* `--enable-mpms-shared=all --with-mpm=event`
* Event MPM
* 异步读写
* 在每模块及每目录上指定日志级别
* 每请求配置,单独文件进行授权: <If>, <ElseIf>, <Else>
* 增强的表达式分析器
* 毫秒级的KeepAlive Timeout,原来是秒级的,时间解析度更加精确
* 基于域名的虚拟主机不再需要NameVirtualHost指令
* 降低了内存占用
* 支持在配置文件中使用自定义变量
* 对于基于IP的访问控制 
		
		2.2中
		
		Order allow,deny
		allow from all
		
		2.4中不再支持此方法
		
		2.4中使用 Require user
		
		Require not .... 就是对所有指令取反 (Require not ip IPADDR)
		
		Require ip IPADDR
			IP
			NETWORK/NETMASK
			NETWORK/LENGTH
			NET
			
			172.16.0.0/255.255.0.0 = 172.16.0.0/16 = 172.16
			
		Require host HOSTNAME
			HOSTNAME
			DOMAIN
			
			www.magedu.com
			magedu.com
			.magedu.com
			.com
		
		Require user USERNAME
		Require group GRPNAME
		
		Require all granted 允许所有主机访问
		Require all deny 拒绝所有主机访问
		
**新增加模块**

mod_proxy 是核心模块

* mod_proxy_fcgi
* mod_proxy_scgi
* mod_proxy_express
* mod_remoteip
* mod_session
* mod_ratelimit
* mod_request

#### 安装mysql

* 把通用二进制包解压到`/usr/local`内
* `ln -sv mysql-5.5.28-linux2.6-i686 mysql`,通过软连接形式链接到mysql(官方要求必须在/usr/local 下 且名字必须为 mysql),用软链的好处是保留了源文件(源文件有版本号等信息)
* 创建mysql用户,mysql组

		groupadd -r -g 306 mysql (这里设定306 只是为了举例.因为3306大于500了就不是系统组了)
		useradd -g 306 -r -u 306 mysql (-r 代表系统用户,不能登录,只是为了运行)
* `chown -R mysql.mysql /usr/local/mysql/*`,修改所属用户,所属组
* mysql初始化db
		
		scripts/mysql_install_db --user=mysql --datadir=/mydata/data

* 复制mysql service脚本

		cp supprot-files/mysql.server /etc/init.d/mysql.d
		如果已经有了执行权限,添加到服务列表中
		chkconfig --add mysqld
		chkconfig --list msqld
			
* 修改mysql 配置文件
	
		根据主机内存大小选择不同的配置文件
		cp supprot-files/my-large.cnf /etc/my.cnf

* 启动mysql

		service mysql start

		但是 客户端在 /usr/local/mysql/bin/mysql 
		
		vim /etc/prodile.d/mysql.sh
		
		export PATH=$PATH:/usr/local/mysql/bin
		
* 输出mysql的库文件(其他程序会使用到mysql的头文件,库文件)

		vim /etc/ld.so.conf.d/mysql.conf
		
		/usr/local/mysql/lib
		
		因为不会立即生效,重新读取库文件,建立库文件缓存
		ldconfig -v
		
		/etc/ld.so.cache 是库缓存文件

* 输出mysql头文件

		ln -sv /usr/local/mysql/include /usr/include/mysql
		

#### mysql配置文件

**配置文件格式**
		
集中式配置文件,可以为多个程序提供配置		
		
		片段式
		[xx]
		xxxx
		
		[mysql] 客户端配置文件
		
		[mysqld] 服务端配置文件
		
		[client] 对所有客户端程序均生效

**配置文件路径**

先查询`/etc/my.cnf`,再查询`/etc/mysql/my.cnf`,在查询`$BASEDIR /my.cnf`,在查询`~/.my.cnf(用户家目录下)`

`$BASEDIR`, mysql进程(实例)运行的目录,一般就是安装目录

在四处查找配置文件.

如果配置冲突,则以最后的为准.后面的会覆盖前面的.

**mysql默认配置文件**
 
mysql默认存在 small,medium,large,huge 等配置文件. 这个是根据内存大小划分的.

**配置参数**

		[mysqld]
		prot	= 3306 (不在同一台主机)
		socket	= /tmp/mysql.sock (客户端和服务端在同一台机器上使用socket)
		...
		thread_cache_size = 8 线程缓存大小
		thread_concurrency = 8 线程并发量(每一个线程占一个cpu)

#### mysql 变量

mysql 服务器维护了两类变量.

**服务器变量**

使用参数可以定义,可以改变服务器工作状态. 定义Mysql服务器运行特性.

* datadir 在什么路径
* 是否启动日志
* 日志放在什么位置
* ....

`SHOW GLOBAL VARIABLES [LIKE 'STRING']`

**状态变量**

保存了mysql服务器运行的统计数据.

`SHOW GLOBAL STATUS [LIKE 'STRING']`

#### mysql 通配符

* `_`: 任意单个字符
* `%` 任意长度的任意字符

#### 编译安装php

		tar xf php5.4.13.tar.hz2
		cd php5.4.13
		
		./configure ...
		
* `--prefix=/usr/local/php` : 安装目录
* `--with-mysql=/usr/local/msql`: mysql的安装目录
* `--with-openssl`: 支持 openssl 功能
* `--with-mysqli=/usr/local/mysql/bin/mysql_config(二进制程序)`: mysql 另外一种和php交互接口
* `--enable-mbstring`: 支持多字节文字(中文,1个字节表述不了的文字)
* `--with-freetype-dir`: 引用特定字体
* `--with-jpeg-dir`: 支持jpeg 
* `--with-png-dir`: 支持png
* `--with-zlib`: 常用通用格式压缩裤
* `--with-libxml-dir=/usr`: xml 库文件路径
* `--enable-xml`: 支持xml(扩展标记语言)
* `--enable-sockets`: php 支持基于套接字的通信
* `--with-apxs2=/usr/local/apacge/bin/apxs`: 可以让php支持让apache编译成模块,把这项改位`--enable-fpm`就会支持fastcgi
* `--with-mcrypt`: 支持额外的加密库
* `--with-config-file-path=/etc`: php 配置文件路径(php.ini)
* `--with-config-file-scan-dir=/etc/php.d`: php 额外配置文件路径(*.ini)
* `--with-bz2`: 压缩裤
* `--enable-maintainer-zts`: apache 如果使用`event`,`worker`,以线程方式工作必须编译成这个形式

### 整理知识点

---

#### partprobe

将磁盘分区表变化信息通知内核,请求操作系统重新加载分区表.

使用fdisk工具只是将分区信息写到磁盘.

而使用partprobe则可以使kernel重新读取分区信息,从而避免重启系统.

		partprobe /dev/sda
		
#### `cat /proc/cpuinfo`

可以查看 cpu 有几核

其他文件是其他配置参数

#### `/etc/profile(文件)`和`/etc/profile.d(目录)` 的区别和用法

1. 两个文件都是设置环境变量文件的,/etc/profile是永久性的环境变量,是全局变量,/etc/profile.d/设置所有用户生效
2. /etc/profile.d/比/etc/profile好维护,不想要什么变量直接删除/etc/profile.d/下对应的shell脚本即可,不用像/etc/profile需要改动此文件

`/etc/profile.d/**.sh` 需要使用自己的环境变量,在其下创建一个`sh脚本`即可.

#### 添加`man 手册`

		vim /etc/my.confg 添加 MANPATH
		
		我的测试服务器是 centos 7 (路径为 vi /etc/man_db.conf)
		
#### `etc/ld.so.conf`

#### `rpm -ivh`

**参数**

* -a: 查询所有套件
* -b<完成阶段><套件档>+或-t <完成阶段><套件档>+：设置包装套件的完成阶段，并指定套件档的文件名称
* -c: 只列出组态配置文件,本参数需配合"-l"参数使用
* -d: 只列出文本文件,本参数需配合"-l"参数使用
* -e<套件档>或--erase<套件档>: 删除指定的套件;
* -f<文件>+: 查询拥有指定文件的套件;
* -h或--hash：套件安装时列出标记; 
* -i<套件档>或--install<套件档>: 安装指定的套件档
* -l: 显示套件的文件列表 
* -p<套件档>+: 查询指定的RPM套件档 
* -q: 使用询问模式,当遇到任何问题时,rpm指令会先询问用户
* -R: 显示套件的关联性信息
* -s: 显示文件状态,本参数需配合"-l"参数使用 
* -U<套件档>或--upgrade<套件档>: 升级指定的套件档
* -v: 显示指令执行过程
* -vv: 详细显示指令执行过程,便于排错

**安装rpm包**

		rpm -ivh your-package.rpm

**安装.src.rpm软件包**

方法一:

		rpm -i your-package.src.rpm 
		cd /usr/src/redhat/SPECS 
		rpmbuild -bp your-package.specs #一个和你的软件包同名的specs文件 
		cd /usr/src/redhat/BUILD/your-package/ #一个和你的软件包同名的目录 
		./configure #这一步和编译普通的源码软件一样，可以加上参数 
		make 
		make install

方法二:

		rpm -i your-package.src.rpm 
		cd /usr/src/redhat/SPECS 
		rpmbuild -bb your-package.specs #一个和你的软件包同名的specs文件

		这时在/usr/src/redhat/RPM/i386/ (根据具体包的不同,也可能是i686,noarch等等)在这个目录下,有一个新的rpm包,这个是编译好的二进制文件.执行rpm -i new-package.rpm即可安装完成.

		rpm -i new-package.rpm即可安装完成

**卸载rpm软件包**

		rpm -e 软件包名(包名可以包含版本号等信息,但是不可以有后缀.rpm)

**显示已经安装的包**

		rpm -aq
		
**显示一个包的文件全名**

		[ansible@rancher-agent-1 ~]$ rpm -q docker
		docker-1.10.3-59.el7.centos.x86_64
		
**显示一个包的文件被安装到哪**

		[ansible@rancher-agent-1 ~]$ rpm -ql docker
		
**某个程序是哪个软件包安装的,或者哪个软件包包含这个程序**

		rpm -qf `which 程序名` #返回软件包的全名 
		rpm -qif `which 程序名` #返回软件包的有关信息 
		rpm -qlf `which 程序名` #返回软件包的文件列表
