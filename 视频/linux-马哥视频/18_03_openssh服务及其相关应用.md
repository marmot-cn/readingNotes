#18_03_OpenSSH服务及其相关应用

---

###笔记

---

###OpenSSH:

telnet, TCP/`23`,远程登录.

* 认证明文
* 数据传输明文

因为以上2个缺点,`ssh`诞生了.

ssh, Secure SHell, TCP/`22`

OpenSSH(ssh 开源版本的实现).

ssh v1(设计有缺陷,不能解决中间人的问题). v2. ssh既是软件又是协议.

OpenSSH 支持 ssh v1, v2.

####客户端

**Linux**

ssh

**Windows**

putty, SecureCRT, SSHSecureShellClient, Xmanager

####服务器端

* sshd

OpenSSH:

* ssh: 客户端
* sshd: 服务端

####认证

telnet,用户输入密码都是明文.

ssh如何实现加密数据传输:

sshd: 主机秘钥

**基于口令认证**

用户名 密码

**基于密钥认证**

client 用户的`公钥`放到服务器对应用户的家目录里面.

client 传输数据用户私钥加密后传到服务器端,服务器就能使用用户的公钥进行解密.

基于秘钥更安全.

####ssh rpm

ssh rpm 通过一些包组成,包括客户端,服务端和一些通用组件.

`netstat -tnl`: 如果当中有`22`(默认)则代码ssh服务启动.

####配置文件

* `/etc/ssh/ssh_config`: 客户端配置文件(红帽系统上的位置)
* `/etc/ssh/sshd_config`: 服务端配置文件(红帽系统上的位置)

####服务器端

服务器端运行的就是`sshd`服务.

`pstree`

	 ├─sshd─┬─sshd───sshd───bash
     │      └─sshd───sshd───bash───pstree
     
当前主机可能有多个`sshd`进程,`sshd`服务起来以后,`sshd`启动主进程,主进程监听一个用户请求,每一个请求进来以后分配一个子进程.登陆的用户多了,`sshd`会启动多个.

####`sshd_config`

* `Protocol 2`, 协议版本,2代表v2.
* `Port xxx`, 端口号
* `AddressFamily any`, `any: ipv4,ipv6`都可以
* `ListenAddress 0.0.0.0`: 指定监听地址(假如有2个IP地址,只希望一个地址提供服务)
* `ListenAddress ::`
* `KeyRegenerationInterval 1h`: 秘钥重新生成的时间间隔
* `ServerKeyBits`: 服务器端秘钥长度
* `SyslogFacility AUTHPRIV`: 日志来源:认证,授权相关的
* `LogLevel INFO`: 日志级别
* `LoginGraceTime 2m`: 登陆宽容器,输入了账号,等待输入密码的时间(因为等待需要服务器的进程一直监控)
* `PermitRootLogin no`: 不允许管理员直接登陆
* `StrictModes yes`: 是否使用严格限定模式
* `MaxAuthTries 6`: 最多允许尝试登陆几次
* `RSAAuthentication yes`: 是否支持`RSA`认证
* `PubkeyAuthentication yes`: 是否基于秘钥认证
* `RhostsRSAAuthentication no`: 如何实现主机认证(我们是不是信任对方主机以后就不需要登陆)
* `PasswordAuthentication yes`: 是否需要基于口令认证
* `Banner none`: 文件路径信息,用户登陆之前显示该信息.
* `PrintMotd yes`: 用户登陆时是否显示`/etc/motd`信息.
* `PrintLastLog yes`: 显示上次通过哪个主机什么时候登陆过
* `Subsystem`: 子系统(Subsystem sftp /usr/libexec/openssh/sftp-server)

		ssh ansible@120.25.161.1 -p17456
		ansible@120.25.161.1's password:
		//PrintLastLog yes
		Last login: Mon Sep 12 13:02:54 2016 from 1.86.241.202
		//PrintLastLog yes
		Welcome to aliyun Elastic Compute Service!

改变配置需要重启服务

####FTP 明文文件传输协议

* `SFTP`: ftp基于`ssh`来实现 
* `FTPS`: 基于`ssl`来实现

####ssh使用

ssh 如果不指用户名,则默认指的是当前主机用户名作为登录名.

* ssh -l USERNAME REMOTE_HOST
* ssh USERNAME@REMOTE_HOST
* ssh `-X -Y`,`Y`比`X`更安全.登陆到远程主机,启动窗口命令.本地必须是窗口的.`enabl X11 forwarding`.

#####实现基于秘钥的认证

一台主机为客户端(基于某个用户实现):

1. 生成一堆秘钥(公钥和私钥)
	* `ssh-keygen` 
		* `-t` {rsa|dsa}
		* `-f` /path/to/keyfile
		* `-N` 'password'
2. 将公钥传输至服务器端某用户的家目录下的`.ssh/authorized_keys`文件中
	* 文件传输工具(ssh-copy-id, scp) 
		* `ssh-copy-id -i /path/to/pubkey USERNAME@REMOTE_HOST`
3. 测试登陆

####scp

基于ssh的远程复制命令,能够实现在主机之间传输数据,基于ssh传输的,传输过程是加密的.

`scp [options] SRC DEST`,和`cp`命令使用一样

* `-a`,复制目录和`cp -a` 一样,相当于 `-rp`,递归,保留文件属性.
* `-r`,复制目录,递归

REMOTE_MACHINE

		USERNAME@HOSTNAME:/path/to/somefile,如果USERNAME省略则以当前主机用户名为登陆名
		
**远程到本地**
		
scp 远程主机:远程主机文件(xxx@xxx:/xx/xx) 当前主机目录

**本地到远程**

scp 当前主机文件 远程主机:远程主机目录(xxx@xxx:/xx/xx)

####sftp

如果已经在ssh配置过秘钥登陆,则不需要登陆.

####dropbear

嵌入式系统专用的`ssh`软件`dropbear`.

###总结

* 秘钥应该经常换且足够复杂
* 使用非默认端口
* 限制登陆客户地址
* 禁止管理直接登陆
* 仅允许有限制用户登陆
* 使用基于秘钥的认证
* 禁止使用版本1.

###整理知识点

---

####sftp



####netstat

`netstat -tnl`

* `-t`: 代表`tcp`.
* `-u`: 代表`udp`.
* `-n`: 代表用数字代表端口名称(`22`->`ssh`).
* `-l`: 代表监听,服务已经启动,等待客户端访问.
	* 不加`l`,表示出于当前正在连接的状态
* `-p`: 那个进程(程序)监听了该服务 
	
`udp`的协议是不会`listen`的,因为不需要三次握手.是没有`tcp`状态的.