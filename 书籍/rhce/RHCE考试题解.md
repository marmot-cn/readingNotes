# RHCE考试题解

## 目录

* [1. 设定 SELinux](#1) **简单**
* [2. 配置防火墙](#2) **中等**
* [3. 自定义用户环境](#3) **简单**
* [4. 配置端口转发](#4)
* [5. 配置链路聚合](#5)
* [6. 配置 IPv6 地址](#6)
* [7. 配置本地邮件服务](#7)
* [8. 通过 SMB 共享目录](#8)
* [9. 配置多用户 SMB 挂载](#9)
* [10. 配置 NFS 服务](#10)
* [11. 挂载一个 NFS 共享](#11)
* [12. 实现一个 web 服务器](#12) **中等**
* [13. 配置安全 web 服务](#13)
* [14. 配置虚拟主机](#14)
* [15. 配置 web 内容的访问](#15)
* [16. 实现动态 web 内容](#16)
* [17. 创建一个脚本](#17) **简单**
* [18. 创建一个添加用户的脚本](#18) **中等**
* [19. 配置 ISCSI 服务端](#19)
* [20. 配置 ISCSI 的客户端](#20)
* [21. 配置一个数据库](#21) **简单**
* [22. 数据库查询](#22) **简单**

### <a name="1">1. 设定 SELinux</a>

在`system1`和`system2`上要求`SELinux`的工作模式为`enforcing`.

* 系统中期后依然生效

#### 答题步骤

在`system1`和`system2`上依次执行

```shell
# 确认现在模式
getenforce

# 设置
setenforce 1

# 再次确认
getenforce

# 编辑/etc/sysconfig/selinux
...
SELINUX=enforcing
...
```

#### 难点

### <a name="2">2. 配置防火墙</a>

在`system1`和`system2`上设定防火墙系统:

* 允许`group8.example.com`域的客户对`system1`和`system2`进行`ssh`访问
* 禁止`my133t.org`域的客户对`system1`和`system2`进行`ssh`访问
* 备注: `my133t.org`是在`172.13.8.0/24`网络

#### 答题步骤

在`system1`和`system2`上执行

```shell
systemctl enable firewalld
firewall-cmd --permanent --add-service=ssh
firewall-cmd --permanent --add-rich-rule 'rule family="ipv4" source address="172.13.8.0/24" service name="ssh" reject'
firewall-cmd --reload

#验证结果
firewall-cmd --list-all 
```

#### 难点

`firewall-cmd`的语法

```
[family="ipv4|ipv6"]
　　　　　　[source |destination] address="address[/mask]" [invert="True|yes"]
　　　　　　[[service name="service name" ]| [port port="number_or_range" protocol="tcp|udp"] | [protocol value="协议名"] ]
　　　　　　[ icmp-block name="icmptype name" ]
　　　　　　[masquerade]
　　　　　　[forward-port port="number_or_range" protocol="tcp|udp" to-port="number_or_range" to-addr="address"]
　　　　　　[log [prefix=prefix text] [level=log level] limit value=rate/duration]
　　　　　　[audit]
　　　　　　[accept | reject [type="reject type"] | drop]
```

### <a name="3">3. 自定义用户环境</a>

在`system1`和`system2`上创建自动以命令为`qstat`, 要求:

* 此自定义命令将执行以下命令: `/bin/pas -Ao pid,tt,user,fname,rsz`
* 此命令对系统中的所有用户有效

#### 答题步骤

在`system1`和`system2`上执行

```shell
/etc/bashrc 文件

alias astat='/bin/pas -Ao pid,tt,user,fname,rsz'
```

验证

```shell
bash
which astat
```

#### 难点

### <a name="4">4. 配置端口转发</a>

#### 答题步骤

#### 难点

### <a name="5">5. 配置链路聚合</a>

#### 答题步骤

#### 难点

### <a name="6">6. 配置 IPv6 地址</a>

#### 答题步骤

#### 难点

### <a name="7">7. 配置本地邮件服务</a>

#### 答题步骤

#### 难点

### <a name="8">8. 通过 SMB 共享目录</a>

#### 答题步骤

#### 难点

### <a name="9">9. 配置多用户 SMB 挂载</a>

#### 答题步骤

#### 难点

### <a name="10">10. 配置 NFS 服务</a>

#### 答题步骤

#### 难点

### <a name="11">11. 挂载一个 NFS 共享</a>

#### 答题步骤

#### 难点

### <a name="12">12. 实现一个 web 服务器</a>

在`system1`配置上一个站点`http://system1.group8.example.com/`, 然后执行下述步骤:

* 从`http://system1.group8.example.com/pub/system1.html`下载文件， 并且将文件重名为`index.html`不要修改此文件内容
* 将文件`index.html`拷贝到您的`web`服务器的`DocumentRoot`目录下
* 来自于`group8.example.com`域的客户端可以访问此`web`服务
* 来自于`my133t.org`域的客户端拒绝访问此`web`服务

#### 答题步骤

##### 12.1 安装软件包

```shell
yum install httpd -y
```

##### 12.2 创建配置文件

```shell
vim /etc/httpd/conf.d/httpd-vhosts.conf
<VirtualHost *:80>
	DocumentRoot "/var/www/html"
	ServerName system1.group8.example.com
	
	<Directory "/var/www/html"> 
		<RequireAll>
			Require all granted
			Require not host .my133t.org 
		</RequireAll>
	</Directory>
</VirtualHost>
```

##### 12.3 复制首页

```shell
wget -O /var/www/html/index.html http://system1.group8.example.com/pub/system1.html
```

##### 12.4 设定服务开机启动并且马上启动服务

```shell
systemctl enable httpd
systemctl start httpd
```

##### 12.5 设定防火墙

```shell
firewall-cmd --permanent --add-service=http
firewall-cmd --reload
```

##### 12.6 验证

```shell
curl system1.group8.example.com
```

#### 难点

* `httpd`配置文件格式, 包括禁用来自xx域
	* `Require not host`
* 设定返货强

### <a name="13">13. 配置安全 web 服务</a>

#### 答题步骤

#### 难点

### <a name="14">14. 配置虚拟主机</a>

#### 答题步骤

#### 难点

### <a name="15">15. 配置 web 内容的访问</a>

#### 答题步骤

#### 难点

### <a name="16">16. 实现动态 web 内容</a>

#### 答题步骤

#### 难点

### <a name="17">17. 创建一个脚本</a>

在`system1`上创建一个名为`/root/foo.sh`的脚本, 让其提供下列特性:

* 运行`/root/foo.sh redhat`, 输出为`fedora`
* 运行`/root/foo.sh fedora`, 输出为`redhat`
* 当没有任何参数或者参数不是`redhat`或者`fedora`时, 其错误输出产生一下的信息:

```
/root/foo.sh redhat | fedora
```

#### 答题步骤

```
#!/bin/bash
case $1 in
  redhat)
  echo "fedora"
  ;;
  fedora)
  echo "redhat"
  ;;
  *)
  echo "/root/foo.sh redhat|fedora"
  ;;
esac
```

赋予脚本执行权限

```
chmod 755 /root/foo.sh
```

#### 难点

* `shell`脚本的`case`用法

### <a name="18">18. 创建一个添加用户的脚本</a>

在`system1`上创建一个脚本, 名为`/root/batchusers`, 此脚本能实现为系统`system1`创建本地用户, 并且这些用户的用户名来自一个包含用户名列表的文件, 同时满足下列要求:

* 此脚本要求提供一个参数, 此参数就是包含用户列表的文件.
* 如果没有提供参数, 此脚本应该给出下面的提示信息`Usage: /root/batchusers userfile`然后退出并返回相应的值
* 如果提供一个不存在文件名, 此脚本应该给出下面的提示信息`Input file not found`然后退出并返回影响的值
* 创建的用户登录`shell`为`/bin/false`
* 此脚本不需要为用户设置密码
* 用户名列表`http://xxxx/userlist`

#### 答题步骤

```shell
/root/batchusers

#!/bin/bash
if [ $# -eq 1 ];then
	if [ -f "$1" ];then
		while read username; do
			useradd -s /bin/false $username &> /dev/null
		done < $1
	else
		echo "Input file not found"
		exit 1
	fi
else
	echo "Usage: /root/batchusers userfile"
	exit 2
fi
	
```

#### 难点

* `done < $1`, `while`循环从文件中读取内容
* 设置登录`shell`, `useradd -s xxx`

### <a name="19">19. 配置 ISCSI 服务端</a>

#### 答题步骤

#### 难点

### <a name="20">20. 配置 ISCSI 的客户端</a>

#### 答题步骤

#### 难点

### <a name="21">21. 配置一个数据库</a>

在`system1`上创建一个`Maraia DB`数据库, 名为`Contacts`, 要求:

* 数据库应该包含来自数据库复制的内容, 复制文件的`URL`为`http://xxxx.mdb`, 数据库只能被`localhost`访问
* 除了`root`用户, 此数据库只能被用户`Mary`查询, 此用户密码为`redhat`
* `root`用户的数据库密码为`redhat`, 同步不先允许空密码登录.

#### 答题步骤

##### 21.1 创建数据库

```shell
yum install mariadb* -y

systemctl enable mariadb
systemctl start mariadb

# 设置数据库
mysql_secure_installation
# 创建root密码为 redhat
```

##### 21.2 导入数据库, 建立用户并授权

```shell
aget -O /root/users.mdb http://xxxx.mdb
mysql -uroot -predhat
mysql> create database Contacts;
mysql> use Contactsl
mysql> source /root/users.mdb;
# 创建 Maray, 只能被localhost访问, 只有查询权限
mysql> grant select on Contacts.* to Maray@localhost identified by 'redhat';
mysql> quit;
```

#### 难点

### <a name="22">22. 数据库查询</a>

在`system1`上使用数据库`Contacts`, 并使用响应的`SQL`查询以回答下列问题:

* 密码是`fadora`的人的名字是什么？
* 有多少人姓名是`John`, 同时居住在`Santa Clara`

#### 答题步骤

##### 22.1 密码是`fadora`的人的名字是什么

```shell
use Contacts;

SELECT u_name.firstname FROM u_name, u_passwd WHERE u_name.userid = u_passwd.uid AND u_passwd.password = 'fadora'
```

##### 22.2 有多少人姓名是`John`, 同时居住在`Santa Clara`

```shell
use Contacts;

SELECT COUNT(*) FROM u_name, u_loc WHERE u_name.userid = u_loc.uid AND u_name.firstname = 'John' AND u_loc.location = 'Santa Clara';
```

#### 难点