# 29_03_samba之一

---

### 笔记

---

* FTP
	* 应用层
* NFS 
	* 应用层
	* 数据传输基于RPC

#### CIFS/SMB

通用互联网文件系统,`Common Internet File System`

`SMB`: 服务信息块

基于RPC

在`Windows`上还要依赖`NETBIOS`协议(网络基本输入输出系统).

`NETBIOS`: 利用广播模式, 在`Windows`同一个工作组之内, 同一个网络之内, 实现主机名称解析的协议. 第一打开会较慢, 需要发送大量广播包去解析主机名称. 效率低后来使用`WINNS`.

`WINNS`, 类似`DNS`服务, 类似`NETBIOS`但是是单拨. 向`WINNS`服务器请求当前网路的主机信息.

`Windows`主机贡献后可以通过`UNC`路径访问:

		UNC:\\IP\\Shared_path(共享名)
		
		同一台主机的共享名不能重复.
		
#### samba

在`Linux`上模拟实现`NetBIOS`和`CIFS/SMB`. 可以和`Windows`交互.

`Linux`打开实现上述协议的服务, 打开服务以后就可以和`Windows`交互了.

`CIFS`弥合里两种不同操作系统的不同文件系统.

`samba`会启动两个进程:

* `nmbd`: 提供 `NetBIOS`
* `smbd`: 提供 文件共享

**Windows主机共享文件监听端口**

* NETBIOS: `udp:137`,`udp:138`,`tcp:139`
* 共享文件: `tcp:445`

使用`iptables`在`linux`主机开通服务, 需要放行这些端口.

**samba名字由来**

早期就是模拟`SMB`,但是`SMB`是一类协议的名称, 后来改名叫`samba`.

**samba如今提供的功能**

* Winbindd + LDAP(windows实现该开源协议的软件是Openldap), 让Linux主机加入到Windows的`AD`域.

**samba如何验证客户端访问**

`NFS`不验证用户.

`samba`验证的账户密码是`samba server`端的账户密码.账户是`samba server`的系统用户. 为了保证系统安全. 账户都是系统账户, 但是密码是`samba`密码. 单独一个文件存储系统用户访问`samba`服务的密码, 是加密存储的.

可以单独服务器验证系统账户.

`samba`的(用户认证)安全级别:

* `user`(默认): 必须提供账户和密码
* `share`: 允许匿名访问
* `server`: 用户名和密码放在服务器认证
* `domain`: 用户名和密码放在`ad域`认证

**客户端访问服务端的权限**

* 共享权限
* 文件系统权限

是上述两个权限的交集.

		账户权限在文件上只有读权限, 共享权限是读写. 但是交集只有读, 所以权限就为读.

**samba配置文件**

配置文件(/etc/samba/)

* `smb.conf`: 主配置文件
* `lmhosts`: `windows`中对`hosts`文件的一个补充.

`smb.conf`:

* `global`: 全局定义
* `homes(share)`: 共享定义
* `printers`: 打印机定义

参数:

* `Network-Related Options`:
	* `netbios name`(winwos访问网路邻居需要主机提供netbios名称): 当前主机的`netbios`名称. 如果没启用就是当前主机名的第一段(www.xxx.com 显示 www)
	* `workgroup`: 工作组
	* `server string`: 主机的说明信息(描述信息)
* `Logging Options`: 日志定义
	* `log file = /var/log/samba/log.%m` %m(客户端主机自己的名称), 每个客户端使用自己独立的日志文件.
	* `max log size = 50` 日志文件最大为50kb, 查过50日志会滚动
* `Standalone Server Options`: 独立守护进程
	* `security = user` 安全级别
	* `passdb backend = tdbsam` 用户的账户和密码存放的格式, tdb格式存储的sam(windows存储账户密码的文件)文件
* `Browser Control Options`: 浏览控制信息
* `Name Resolution`: 名称解析相关信息
* `Printing Options`: 
	* `load printers = yes` 是否加载打印机
	* `cups options = raw` unix打印服务的选项(哪种打印机驱动)
* `File System Options`: 文件系统选项
* `Share Definitions`: 共享定义
	* `homes` 定义一个用户是否可以访问自己的家目录
		* `comment`: 注释, 说明信息
		* `browseable`: 能否可以被浏览到
		* `writable`: 是否具有写权限
	* `printers` 共享打印机
		* `guest ok = no` 是否允许来宾账户访问(匿名访问)

**一般共享一个目录用到的参数**

`[共享名称]`: 

* `comment = ` 注释
* `path = ` 资源路径
* `browseable = ` 是否可以被浏览
* `guest ok = ` 来宾账户, 匿名访问
* `public = ` 是否被所有用户读
* `read only = ` 是否只读
* `writable = ` 是否可写(不允许 read only = yes 和 writable = yes 同时出现)
* `wriet list = user1, user2... 或 @group(该组用户可写) +group(该组用户可写)` 具有写权限的用户列表
* `valid users = ` 白名单, 限定哪些用户可以访问
* `invalid users = ` 黑名单(如果和白名单同时使用, 白名单生效)

```shell
share 目录下的 test 共享, 共享名叫 tools

安装samba
[root@iZ944l0t308Z ~]# yum install samba -y
...
[root@iZ944l0t308Z ~]# ls /etc/samba/
lmhosts  smb.conf  smb.conf.example

[root@iZ944l0t308Z ~]# mkdir /shared/test -pv
mkdir: created directory '/shared/test'

编辑samba配置文件

vim /etc/samba/smb.conf
...
[tools]
	comment = Share testing
	path = /shared/test/
	public = yes
	writable = yes
...

使用 testparm 测试配置文件是否正确
Load smb config files from /etc/samba/smb.conf
Processing section "[homes]"
Processing section "[printers]"
Processing section "[print$]"
Processing section "[tools]"
Loaded services file OK.
Server role: ROLE_STANDALONE

Press enter to see a dump of your service definitions

...
一些参数被自动换成等价参数
[tools]
	comment = Share testing
	path = /shared/test/
	public = yes
	writable = yes
...

启动 smba
[root@iZ944l0t308Z ~]# service smb start
Redirecting to /bin/systemctl start  smb.service

启动 nmb
[root@iZ944l0t308Z ~]# service nmb start
Redirecting to /bin/systemctl start  nmb.service

可见 tcp 445, 139 端口启动. udp 137 138 
[root@iZ944l0t308Z ~]# netstat -tunlp
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
tcp        0      0 0.0.0.0:445             0.0.0.0:*               LISTEN      5046/smbd
tcp        0      0 0.0.0.0:139             0.0.0.0:*               LISTEN      5046/smbd
tcp        0      0 0.0.0.0:111             0.0.0.0:*               LISTEN      1/systemd
tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      831/sshd
tcp        0      0 127.0.0.1:631           0.0.0.0:*               LISTEN      4662/cupsd
tcp6       0      0 :::445                  :::*                    LISTEN      5046/smbd
tcp6       0      0 :::139                  :::*                    LISTEN      5046/smbd
tcp6       0      0 :::111                  :::*                    LISTEN      3785/rpcbind
udp        0      0 0.0.0.0:992             0.0.0.0:*                           3785/rpcbind
udp        0      0 0.0.0.0:111             0.0.0.0:*                           3785/rpcbind
udp        0      0 120.25.87.35:123        0.0.0.0:*                           500/ntpd
udp        0      0 10.170.148.109:123      0.0.0.0:*                           500/ntpd
udp        0      0 127.0.0.1:123           0.0.0.0:*                           500/ntpd
udp        0      0 0.0.0.0:123             0.0.0.0:*                           500/ntpd
udp        0      0 10.170.151.255:137      0.0.0.0:*                           7483/nmbd
udp        0      0 10.170.148.109:137      0.0.0.0:*                           7483/nmbd
udp        0      0 120.25.87.255:137       0.0.0.0:*                           7483/nmbd
udp        0      0 120.25.87.35:137        0.0.0.0:*                           7483/nmbd
udp        0      0 0.0.0.0:137             0.0.0.0:*                           7483/nmbd
udp        0      0 10.170.151.255:138      0.0.0.0:*                           7483/nmbd
udp        0      0 10.170.148.109:138      0.0.0.0:*                           7483/nmbd
udp        0      0 120.25.87.255:138       0.0.0.0:*                           7483/nmbd
udp        0      0 120.25.87.35:138        0.0.0.0:*                           7483/nmbd
udp        0      0 0.0.0.0:138             0.0.0.0:*                           7483/nmbd
udp6       0      0 :::992                  :::*                                3785/rpcbind
udp6       0      0 :::111                  :::*                                3785/rpcbind
udp6       0      0 :::123                  :::*                                500/ntpd

创建用户账户密码
[root@iZ944l0t308Z ~]# smbpasswd -a chloroplast
New SMB password:
Retype new SMB password:
Added user chloroplast.

```

`Windows`链接`samba`, 我使用`win10`链接有问题, 最后换了阿里云的`win 2008 server 标准版`测试可以访问.

### 整理知识点

---