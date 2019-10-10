#16_04_Linux系统裁减之六 为系统添加ssh服务

---

###笔记

---

####telnet: 远程登录协议

应用层协议,基于`TCP`,`23`(默认)端口.

`C/S`架构.

* S: telnet服务器
* C: telnet客户端

早期远程登录都是使用`telnet`服务,但是无论其命令还是认证过程都是明文发布,不安全.后来使用`ssh`

####SSH: Secure SHell

远程登录协议.引用层协议,基于`TCP`,`22`(默认)端口.

通信过程及认证过程是加密的.

* 主机认证
* 用户认证过程加密
* 数据传输过程加密

ssh v1, v2. 协议版本.

v1 已经不再安全了,能够被破解了,无法避免中间人攻击(man-in-middle).

**SSH 支持2种认证**

认证过程:

* 基于口令认证
* 基于密钥认证

**协议**

协议只是一种`规范`.

实现:服务端,客户端.

服务端和客户端使用`规范`来实现.

**openSSH**

Linux: openSSH

c/s

* 服务器端: sshd, 配置文件`/etc/ssh/sshd_config`
* 客户端, ssh, 配置文件`/etc/ssh_config`

`ssh-keygen`: 密钥生成器. 为某个用户生成一对密钥.

`ssh-copy-id`: 将公钥传输至远程服务器的.

`scp`: 跨主机安全复制工具,两台主机之间复制数据,复制的过程是加密的.

**ssh**

* `ssh USERNAME@HOST`,指定用户ssh登录
* `ssh -l USERNAME HOST`,指定用户ssh登录
* `ssh USERNAME@HOST 'COMMAND'`, 以这个用户的身份在远程主机执行`COMMAND`命令,并退出.命令发送和结果回送都是加密的.

**scp(和`cp`命令一样)**

* `scp SRC DEST`,如果有多个源,目标必须是目录.
 * `-r`: 递归
 * `-a`

**ssh-keygen**

`ssh-keygen -t rsa`,(`-t`指定加密方式,一般使用`rsa`),生成密钥后保存在当前用户的家目录下.`-f /path/to/KEY_FILE`指定秘钥文件.`-P ''`指定加密私钥的密码串.

* `~/.ssh/id_rsa`: 私钥
* `~/.ssh/id_rsa.pub`: 共钥,公钥需要追加保存到远程主机对应用户的家目录下的`~/.ssh/authorized_keys`这个文件当中或`~/.ssh/authorized_keys2`文件当中.

可能有多个秘钥保存在`~/.ssh/authorized_keys`中,一定要`追加`到这个文件中,`不能复制`到这个文件中.

**ssh-copy-id**

`ssh-copy-id`

* `-i` 指定`key`的位置,一般是`~/.ssh/id_rsa.pub`

####示例

**scp:远程文件复制到本地**

		SCP USERNAME@HOST:/path/to/somefile /path/to/local
		以USERNAME这个用户身份去远程主机找somefile 并复制到本地
		
**scp:本地文件到远程主机,放在某个目录下**

		SCP /path/to/local USERNAME@HOST:/path/to/somewhere

**ssh 没有指用户,则默认用当前主机用户**

		[chloroplast@k8s-master-test ~]$ ssh 120.25.161.1 -p 17456
		chloroplast@120.25.161.1's password:

**服务器传输做主机认证的公钥**

		[chloroplast@k8s-master-test ~]$ ssh chloroplast@120.25.161.1 -p 17456
		The authenticity of host '[120.25.161.1]:17456 ([120.25.161.1]:17456)' can't be established.
		ECDSA key fingerprint is 58:a2:7c:11:7b:af:97:80:ba:d1:a1:9e:1e:b3:f2:96.
		Are you sure you want to continue connecting (yes/no)?
		
主机密钥放在`~/.ssh/known_hosts`,保存在登录用户的家目录下.

		[chloroplast@k8s-master-test ~]$ cat .ssh/known_hosts
		[120.25.161.1]:17456 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBHRnA7fqDMzEu8V/0GMRKS8jI5CKb0VEaLiTa5Jj745NBa2E7manPuIIlMZ2CWybAtoFYP+qWe86ymQuIZFwL/Q=

**ssh-keygen**

		[chloroplast@k8s-master-test ~]$ ssh-keygen -t rsa
		Generating public/private rsa key pair.
		Enter file in which to save the key (/home/chloroplast/.ssh/id_rsa):
		Enter passphrase (empty for no passphrase):
		Enter same passphrase again:
		Your identification has been saved in /home/chloroplast/.ssh/id_rsa.
		Your public key has been saved in /home/chloroplast/.ssh/id_rsa.pub.
		The key fingerprint is:
		a4:f6:21:84:53:fc:03:34:f7:fb:ec:ca:17:56:e6:03 chloroplast@k8s-master-test
		The key's randomart image is:
		+--[ RSA 2048]----+
		|     o+ .        |
		|     ooo .       |
		|    o .o. .      |
		|     o oo  .E o  |
		|      + S..  =   |
		|     . o . oo o  |
		|        .  .o. . |
		|         . ..    |
		|          oo.    |
		+-----------------+

		[chloroplast@k8s-master-test ~]$ ls .ssh/
		id_rsa  id_rsa.pub  known_hosts

**ssh-copy-id**

		[chloroplast@k8s-master-test ~]$ ssh-copy-id -p 17456 -i ~/.ssh/id_rsa.pub chloroplast@120.25.161.1
		/usr/bin/ssh-copy-id: INFO: attempting to log in with the new key(s), to filter out any that are already installed
		/usr/bin/ssh-copy-id: INFO: 1 key(s) remain to be installed -- if you are prompted now it is to install the new keys
		chloroplast@120.25.161.1's p assword:
		
		Number of key(s) added: 1
		
		Now try logging into the machine, with:   "ssh -p '17456' 'chloroplast@120.25.161.1'"
		and check to make sure that only the key(s) you wanted were added.
		
		[chloroplast@k8s-master-test ~]$ ssh -p 17456 chloroplast@120.25.161.1
		Last login: Mon Jun 27 10:47:29 2016
		
		Welcome to aliyun Elastic Compute Service!
		
		[chloroplast@k8s-minon2-test ~]$ exit

####dropbear: 嵌入式系统专用的ssh服务器端和客户端工具

服务器端: 

* `dropbear`
* `dropbearkey`,秘钥生成工具.服务器端为主机认证时生成秘钥
	* `-t`: `rsa|dsa`
	* `-f`: 指定文件存放位置 `path/to/KEY_FILE`
	* `-s SIZE`: 指定长度

客户端: `dbclinet`

在宿主机上编译`dropbear`默认使用`nsswitch`实现名称解析.

* `/etc/nsswitch.conf`
* `/lib/libnss_files*`
* `/usr/lib/libnss3.so`
* `/usr/lib/libnss_files*`

dropbear会在用户登录时检查其默认shell是否当前系统的安全shell

* `/etc/shells`

主机秘钥默认位置:

* `/etc/dropbear/`
	* RSA: dropbear_rsa_host_key,长度可变,只要是8的整数倍,默认为1024
	* DSS: dropbear_dss_host_key,长度固定,默认为1024
	

###整理知识点

---

####SSH协议交互过程 


####cp -a

相当于`cp --dpR`,保持文件的连接(`d`),保持原文件的属性(`p`)并作递归处理(`R`).

**cp -d**
		
		创建一个文件a
		[ansible@k8s-minon2-test ~]$ touch a
		[ansible@rancher-agent-2 ~]$ echo 1 > a
		建立一个快捷方式,软连接
		[ansible@k8s-minon2-test ~]$ cp -s a a_slink
		创建一个硬连接				
		[ansible@k8s-minon2-test ~]$ cp -l a a_hlink
		
		[ansible@k8s-minon2-test ~]$ ll -s
		总用量 0
		0 -rw-rw-r-- 2 ansible ansible 0 6月  27 10:39 a
		0 -rw-rw-r-- 2 ansible ansible 0 6月  27 10:39 a_hlink
		0 lrwxrwxrwx 1 ansible ansible 1 6月  27 10:40 a_slink -> a
		
将`a_slink`复制成为`a_slink2`.
		
		[ansible@k8s-minon2-test ~]$ cp a_slink a_slink2
		原本要复制的是连接文件，却将连接文件连接的实际文件复制过来了
		[ansible@k8s-minon2-test ~]$ ll -s a_slink2
		0 -rw-rw-r-- 1 ansible ansible 0 6月  27 10:41 a_slink2
		[ansible@rancher-agent-2 ~]$ cat a_slink2
		1
		
若要复制连接文件而不是它指向的源文件,就要使用`-d`参数.

		[ansible@k8s-minon2-test ~]$ cp -d a_slink a_slink3
		[ansible@k8s-minon2-test ~]$ ll -s a_slink3				
		0 lrwxrwxrwx 1 ansible ansible 1 6月  27 10:42 a_slink3 -> a
		
####netstat

netstat 命令用于显示各种网络相关信息,如网络连接,路由表,接口状态(Interface Statistics),masquerade 连接,多播成员 (Multicast Memberships) 等等.

netstat的输出结果可以分为两个部分:

* `Active Internet connections`,称为有源TCP连接,其中"Recv-Q"和"Send-Q"指%0A的是接收队列和发送队列.这些数字一般都应该是0.如果不是则表示软件包正在队列中堆积.这种情况只能在非常少的情况见到.
* `Active UNIX domain sockets`,称为有源Unix域套接口(和网络套接字一样,但是只能用于本机通信,性能可以提高一倍).
	* `Proto`: 显示连接使用协议
	* `RefCnt`: 表示连接到本套接口上的进程号
	* `Types`: 显示套接口的类型
	* `State`: 显示套接口当前的状态
	* `Path`: 表示连接到套接口的其他进程使用的路径名 

**常见参数**
		
		常用
		netstat -tnl

* `-a(all)`: 显示所有选项,默认不显示LISTEN相关
* `-t(tcp)`: 仅显示tcp相关选项
* `-u(udp)`: 仅显示udp相关选项
* `-n`: 拒绝显示别名,能显示数字的全部转化成数字.
* `-l`: 仅列出在`Listen(监听)`的服务状态
* `-p`: 显示建立相关链接的程序名
* `-r`: 显示路由信息,路由表
* `-e`: 显示扩展信息,例如uid等
* `-s`: 俺各个协议进行统计
* `-c`: 每隔一个固定时间,执行该`netstat`命令

####nsswitch

`/etc/nsswitch.conf`(name service switch configuration,名字服务切换配置)规定通过哪些途径以及按照什么顺序通过这些途径来查找特定类型的信息.还可以指定某个方法奏效或失效时系统将采取什么动作.

nsswitch.conf中的每一行配置都指明了如何搜索信息,每行配置的格式如下:

		info: method[[action]] [method[[action]]...]
		
**nsswitch.conf的工作原理**

当需要提供`nsswitch.conf`文件所描述的信息的时候,系统将检查含有适当info字段的配置行.它按照从左向右的顺序开始执行配置行中指定的方法.在默认情况下,如果找到期望的信息,系统将停止搜索.如果没有指定action,那么当某个方法未能返回结果时,系统就会尝试下一个动作.有可能搜索结束都没有找到想要的信息.

**信息(info)**

nsswitch.conf文件通常控制着用户(在passwd中),口令(在shadow中),主机IP和组信息(在group中)的搜索.下面的列表描述了nsswitch.conf文件控制搜索的大多数信息(Info项)的类型:

* `automount`: 自动挂载(/etc/auto.master和/etc/auto.misc)
* `bootparams`: 无盘引导选项和其他引导选项
* `ethers`: MAC地址
* `group`: 用户所在组(/etc/group),getgrent()函数使用该文件
* `hosts`: 主机名和主机号(/etc/hosts),gethostbyname()以及类似的函数使用该文件
* `networks`: 网络名及网络号(/etc/networks),getnetent()函数使用该文件
* `passwd`: 用户口令(/etc/passwd),getpwent()函数使用该文件
* `protocols`: 网络协议(/etc/protocols),getprotoent()函数使用该文件
* `publickey`: NIS+及NFS所使用的secure_rpc的公开密钥
* `rpc`: 远程过程调用名及调用号(/etc/rpc),getrpcbyname()及类似函数使用该文件
* `services`: 网络服务(/etc/services),getservent()函数使用该文件
* `shadow`: 映射口令信息(/etc/shadow),getspnam()函数使用该文件
* `aiases`: 邮件别名,sendmail()函数使用该文件

**方法(method)**

下面列出了nsswich.conf文件控制搜索信息类型的方法,对于每一种信息类型,都可以指定下面的一种或多种方法:

* `files`: 搜索本地文件,如`/etc/passwd`和`/etc/hosts`
* `nis`: 搜索NIS数据库,nis还有一个别名,即yp
* `dns`: 查询DNS(只查询主机)
* `compat`: passwd、group和shadow文件中的±语法

**搜索顺序(从左至右)**

两个或者更多方法所提供的信息可能会重叠.举例来说,files和nis可能都提供同一个用户的口令信息.如果出现信息重叠现象,就需要考虑将哪一种方法作为权威方法(优先考虑),并将该方法放在方法列表中`靠左`的位置上.

默认nsswitch.conf文件列出的方法并没有动作项,并假设没有信息重叠(正常情况).在这种情况下,搜索顺序无关紧要:当一种方法失败之后,系统就会尝试下一种方法,只是时间上受到一点损失.如果在方法之间设置了动作,或者重叠的项的内容不同,那么搜索顺序就变得重要起来.

例如下面两行nsswitch.conf文件配置行:

		passwd files nis
		host nis files dns
		
第一行让系统在`/etc/passwd`文件中搜索口令信息,如果失败的话,就使用NIS来查找信息.如果正在查找的用户同时出现在这两个地方,就会使用本地文件中的信息.第二行先使用NIS搜索;如果失败的话,就搜索/etc/hosts文件;如果再次失败的话,核对DNS以找出主机信息.

**动作项([action])**

在每个方法后面都可以选择跟一个动作项,用来指定如果由于某种原因该方法成功抑或失败需要做些什么.动作项的格式如下:

		[[!]STATUS =action]
		
其中,开头和末尾的方括号属于格式的一部分,并不是用来指出括号中的内容是可选的.`STATUS`(按照约定使用大写字母,但本身并不区分大小写)是待测试的状态,`action`是如果`STATUS`匹配前面的方法所返回的状态将要执行的动作.开头的感叹号(!)是可选的,其作用是将状态取反.

`STATUS`取值如下: 

* `NOTFOUND`: 方法已经执行,但是并没有找到待搜索的值.默认的动作是`continue`.
* `SUCCESS`: 方法已经执行,并且已经找到待搜索的值,没有返回错误.默认动作是`return`.
* `UNAVAIL`: 方法失败,原因是永久不可用.举例来说,所需的文件不可访问或者所需的服务器可能停机.默认的动作是`continue`.
* `TRYAGAIN`: 方法失败,原因是临时不可用.举例来说,某个文件被锁定,或者某台服务器超载.默认动作是`continue`.

`action`:

* `return`: 返回到调用例程,带有返回值,或者不带返回值.
* `continue`: 继续执行下一个方法.任何返回这都会被下一个方法找到的值覆盖.

`示例`:

		hosts    dns [!UNAVAIL=return] files
		
		如果DNS方法没有返回UNAVAIL(!UNAVAIL),也就是说DNS返回SUCCESS,NOTFOUND或者TRYAGAIN,那么系统就会执行与该STATUS相关的动作(return).其结果就是,只有在DNS服务器不可用的情况下才会使用后面的方法(files).
		
**compat方法 passwd,group和shadow文件中的"±"**

可以在`/etc/passwd`,`/etc/group`和`/etc/shadow`文件中放入一些特殊的代码,(如果在nsswitch.conf文件中指定compat方法的话)让系统将本地文件和NIS映射表中的项进行合并和修改.

在这些文件中,如果在行首出现加号'＋',就表示添加NIS信息;如果出现减号'－',就表示删除信息.举例来说,要想使用passwd文件中的这些代码,可以在nsswitch.conf文件中指定passwd: compat.然后系统就会按照顺序搜寻passwd文件,当它遇到以+或者-开头的行时,就会添加或者删除适当的NIS项.

虽然可以在passwd文件的末尾放置加号,在nsswitch.conf文件中指定passwd: compat,以搜索本地的passwd文件,然后再搜寻NIS映射表,但是更高效的一种方法是在nsswitch.conf文件中添加passwd: file nis而不修改passwd文件.


