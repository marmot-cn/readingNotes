# RHCE考试题解

## 目录

* [1. 设定 SELinux](#1) **简单**
* [2. 配置防火墙](#2) **中等**
* [3. 自定义用户环境](#3) **简单**
* [4. 配置端口转发](#4) **简单**
* [5. 配置链路聚合](#5) **中等**
* [6. 配置 IPv6 地址](#6) **简单**
* [7. 配置本地邮件服务](#7) **中等**
* [8. 通过 SMB 共享目录](#8) **困难**
* [9. 配置多用户 SMB 挂载](#9) **困难**
* [10. 配置 NFS 服务](#10) **困难**
* [11. 挂载一个 NFS 共享](#11) **中等**
* [12. 实现一个 web 服务器](#12) **中等**
* [13. 配置安全 web 服务](#13) **简单**
* [14. 配置虚拟主机](#14) **简单**
* [15. 配置 web 内容的访问](#15) **简单**
* [16. 实现动态 web 内容](#16) **中等**
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

* 命令对系统中的所有用户有效, 需要修改文件`/etc/bashrc`文件
* 了解`alias`命令

### <a name="4">4. 配置端口转发</a>

在系统`system1`设定端口转发

* 在`172.24.8.0/24`网络中的系统, 访问`system1`的本地端口`5423`将被转发到`80`
* 此设置必须永久有效

#### 答题步骤

在`system1`上执行

```shell
firewall-cmd --permanent --add-rich-rule 'rule family="ipv4" source address="172.24.8.0/24" forward-port port="5423" protocol="tcp" to-port="80"'
```

验证

```shell
firewall-cmd --list-all
```

#### 难点

`firewall-cmd`的语法

### <a name="5">5. 配置链路聚合</a>

在`system2`和`system1`之间按以下要求设定一个链路:

* 此链路使用接口`eth1`和`eth2`
* 此链路在一个接口失效时仍然能工作
* 此链路在`system1`使用下面的地址`172.16.3.40/255.255.255.0`
* 此链路在`system1`使用下面的地址`172.16.3.45/255.255.255.0`
* 此链路在系统重启之后依然保持正常状态

#### 答题步骤

在`system1`和`system2`上执行

```shell
nmcli connection add con-name team0 type team ifname team0 config '{"runner":{"name":"activebackup"}}'

#system1
nmcli connection modify team0 ipv4.address "172.16.3.40/24" connection.autoconnect yes ipv4.method manual

#system2
nmcli connection modify team0 ipv4.address "172.16.3.45/24" connection.autoconnect yes ipv4.method manual

nmcli connection add con-name team0-salve1 ifname eth1 type team-slave master team0
nmcli connection add con-name team0-salve2 ifname eth2 type team-slave master team0
nmcli connection up team0
```

验证:

```shell
ifconfig team0
....

teamdctl team0 state

在system2上 ping 172.16.3.40
```

如果配置错误

```shell
nmcli connection delete team0-slave1
nmcli connection delete team0-slave2
nmcli connection delete team0
nmcli connection reload
```

#### 难点

* `nmcli`的命令使用
* 配置主备`activebackup`

### <a name="6">6. 配置 IPv6 地址</a>

设定接口`eth0`使用下列`IPv6`地址:

* `system1`上的地址应该是`2003:ac18::305/64`
* `system2`上的地址应该是`2003:ac18::30a/64`
* 两个系统必须能与网络`2003:ac18/64`内的系统通信
* 地址必须在重启后依然生效
* 两个系统必须保持当前的`IPv4`地址并能通信

#### 答题步骤

在`system1`和`system2`上执行

```shell
# system1
nmcli connection modify eth0 ipv6.addresses "2003:ac18::305/64" ipv6.method manual connection.autoconnect yes

# system2
nmcli connection modify eth0 ipv6.addresses "2003:ac18::30a/64" ipv6.method manual connection.autoconnect yes

nmcli connection reload
nmcli connection down eth0 && nmcli connection up eth0
```

验证

```shell
ping6 -c1 2003:ac18::30a
```

#### 难点

* `nmcli`配置`IPV6`的命令

### <a name="7">7. 配置本地邮件服务</a>

在系统`system1`和`system2`上配置邮件服务, 要求:

* 这些系统不接受外部发来的邮件
* 在这些系统上本地发送的任何邮件都会自动路由到`mail.group8.exmaple.com`
* 从这些系统上发送的优先见识来自于`server.group8.example.com`

测试: 可以发送邮件到本地用户`dave`来测试配置, 系统`server.group8.example.com`已经把此用户的邮件转到`URL`:`http://server.group8.example.com/pub/received_mail/8`.

#### 答题步骤

在`system1`和`system2`上执行

```shell
postconf -e inet_interfaces=loopback-only
postconf -e mydestindation=
postconf -e local_transport=error:err
postconf -e relayhost=[mail.group8.exmaple.com]
postconf -e myorigin=server.group8.exmaple.com
systemctl enable postfix
systemctl restart postfix
```

验证

```shell
echo "hello" | mail -s testmail dave
curl http://server.group8.example.com/pub/received_mail/8
....
```

#### 难点

* `postconf`的配置
	* `inet_interfaces`
	* `mydestindation`
	* `local_transport`禁止本地分发邮件到本地用户邮箱 
	* `relayhost`
	* `myorigin`

### <a name="8">8. 通过 SMB 共享目录</a>

在`system1`上配置`SMB`服务, 要求:

* `SMB`服务器必须是`STFF`工作组的一个成员
* 共享`/common`目录, 共享名必须为`common`
* 只有`group8.example.com`域内的客户端可以访问`common`共享
* `common`必须是可以浏览的
* 用户`andy`必须能够读取共享中的内容, 验证密码是`redhat`

#### 答题步骤

在`system1`上执行

##### 1. 安装软件包

```shell
yum install samba samba-client -y
systemctl enable smb.service nmb.service
```

##### 2. 设定防火墙

```shell
firewall-cmd --permanent --add-service=samba
firewall-cmd --reload
```

##### 3. 编辑配置文件

```shell
/etc/samba/smb.conf

...
[global]
	...
	workgroup=STAFF
	...

[common]
	...
	path=/common
	hosts allow=172.24.8.0/24
	borwseable=yes
```

##### 4. 创建目录并设定`selinux`上下文
```shell
mkdir /common
semanage fcontext -a -t "samba_share_t" '/common(/.*)?'
restorecon -Rv /common/
```

##### 5. 创建 samba 用户

```shell
smbpasswd -a andy
...输入密码 redhat
```

##### 6. 启动服务

```shell
systemctl restart smb nmb
```

##### 7. 验证

在`system2`上执行

```shell
# 安装软件包
yum install samba-client -y

# 验证服务是否可以访问
smbclient -L //172.24.8.11/ -U andy
...

smbclient -L //172.24.8.11/ -U andy
```

#### 难点

* `/etc/samba/smb.conf`配置文件修改规则
* `selinux`上下文

### <a name="9">9. 配置多用户 SMB 挂载</a>

在`system1`上通过`SMB`共享目录`/devops`, 并满足下列要求:

* 共享名为`devops`
* 共享目录`devops`只能`group8.example.com`域中的客户端使用
* 共享目录`devops`必须可以被浏览
* 用户`silene`必须能以读的方式访问此共享, 访问密码是`redhat`
* 用户`akira`必须能以读写的方式访问此共享, 访问密码是`redhat`
* 此共享永久挂载在`system2.group8.exmaple.com`上的`/mnt/dev`目录, 并使用用户`silene`违认证任何用户, 可以通过用户`akira`来临时获取写的权限

#### 答题步骤

##### 1. 修改配置文件 (system1)

```shell
# /etc/samba/smb.conf

[devops]
	path=/devops
	hosts allow=172.24.8.
	browseable=yes
	writable=no
	write list=akira
```

##### 2. 设定目录 (system1)

```shell
mkdir /devops
semanage fcontext -a -t 'samba_share_t' '/devops(/.*)?'
restorecon -Rv /devops
setfacl -m u:akira:rwx /dev/ops
```

##### 3. 建立用户 (system1)

```shell
smbpasswd -a silene
...输入密码 redhat

smbpasswd -a akira
...输入密码 redhat
```

##### 4. 重启服务 (system1)

```shell
systemctl restart smb nmb
```

##### 5. 安装软件 (system2)

```shell
yum install cifs-utils -y
```

##### 6. 建立挂载点 (system2)

```shell
mkdir /mnt/dev
```

##### 7. 编辑 fstab 增加持久挂载 (system2)

```shell
#查看共享
smbclient -L //172.24.8.11/ -U silene

#编辑/etc/fstab
//172.24.8.11/devops /mnt/dev cifs defaults,multiuser,username=silene,password=redhat,sec=ntlmssp 0 0

#测试
mount -a
df -h

#测试silene
su - silene
cd /mnt/dev
cifscreds add 172.24.8.11
...输入密码
touch testfile #不能创建成功, 可以读不能写
exir

#测试akira
su = akira
cd /mnt/dev
cifscreds add 172.24.8.11
...输入密码
ls
touch testfile #可以创建成功
ls
exit
```

挂载参数需要添加"multiuser,sec=ntlmssp", 客户机上的普通用户可以通过`cifscreds`命令提交新的身份凭据.

#### 难点

* `smb`服务端配置
* `smb`客户端挂载

### <a name="10">10. 配置 NFS 服务</a>

在`system1`配置`NFS`服务, 要求:

* 以只读方式共享目录`/public`, 同时只能被`group8.example.com`域中的系统访问
* 以读写的方式共享目录`/protected`, 同时只能被`group8.example.com`域中的系统访问
* 访问`/protocted`需要通过`Kerberos`安全加密, 您可以使用下面`URL`提供密钥`http://server.group8.example.com/pub/keytabs/system1.keytab`
* 目录`/protocted`应该包含名为`project`拥有认为`andres`的子目录
* 用户`andres`能以读写访问访问`/protocted/project`

#### 答题步骤

##### 1. 安装软件

```shell
yum install nfs-utils -y
```

##### 2. 服务器开机启动

```shell
systemctl enable nfs-server nfs-secure-server
```

##### 3. 设定防火墙规则

```shell
firewall-cmd --permanent --add-service=nfs
firewall-cmd --permanent --add-service=rpc-bind
firewall-cmd --permanent --add-service=mountd
firewall-cmd --reload
```

##### 4. 设定目录及 SELinux 安全上下文

```shell
mkdir -p /public /protected/project
chown andres /protected/project
semanage fcontext -a -t'public_contect_t' '/protected(/.*)?'
semanage fcontext -a -t'public_contect_rw_t' '/protected/project(/.*)?'
restorecon -Rv /protected
```

##### 5. 下载 kerberos 证书

```shell
wget -O /etc/krb5.keytab htt[://server.group8.example.com/pub/keytabs/system1.keytab
```

##### 6. 编辑 nfs 资源导出配置文件

```shell
# /etc/exports

/public *.group8.example.com(ro,sec=sys,sync)
/protected *.group8.example.com(rw,sec=krb5p,sync)
```

##### 7. 修改 nfs 启动参数, 重启服务

```shell
# /etc/sysconfig/nfs
...
RPCNFSDARGS="-V 4.2"
...

systemctl restart nfs-server.service nfs-secure-server.service
```

##### 8. 刷新并验证导出资源

```shell
exports -ra
exportfs 
```

#### 难点

* 设定防火墙
* 设定`SELinux`安全上下文
* `nfs`先关配置文件
	* `/etc/exports`
	* `/etc/sysconfig/nfs`

### <a name="11">11. 挂载一个 NFS 共享</a>

在`system2`上挂载一个来自`system1.group8.example.com`的`NFS`共享, 要求:

* `/public`挂载在下面的目录上`/mnt/nfsmount`
* `/protected`挂载在下面的目录上`/mnt/nfssecure`并使用安全的方式, 密钥下载的`URL`如下:`http://server.group8.example.com/pub/keytabs/system2.keytab`
* 用户`andres`能够在`/mnt/nfssecure/project`上创建文件
* 这些文件系统在系统启动时自动挂载

#### 答题步骤

##### 1.建立挂载点

```shell
mkdir /mnt/nfsmount /mnt/nfssecure
```

##### 2. 下载 kerberos 证书

```shell
wget -O /etc/krb5.keytab htt[://server.group8.example.com/pub/keytabs/system2.keytab
```

##### 3. 修改 fstab 添加持久化挂载选项

```
# /etc/fstab
system1:/public /mnt/nfsmount nfs defaults,sec=sys 0 0
system1:/protected /mnt/nfssecure nfs defaults,sec=krb5p,v4.2 0 0
```

###### 4. 设定开启及启动

```shell
systemctl enable nfs-secure
systemctl start nfs-secure.service
```

###### 5. 测试

```shell
mount -a
df -h

su = andres
kinit
...输入密码
cd /mnt/nfssecure/project/
touch testfile
ls -l #确认创建成功
```

#### 难点

* 客户端挂载
* kerberos验证

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
* 设定防火墙

### <a name="13">13. 配置安全 web 服务</a>

为站点`http://system1.group8.example.com`配置`TLS`加密:

* 一个以签名证书从`http://证书链接`获取
* 此证书的密钥从`http://证书密钥链接`获取
* 此证书的签名授权信息从`http://证书签名授权信息链接`获取

#### 答题步骤

##### 1. 安装`ssl`模块

```shell
yum install mod_ssl -y
```

##### 2. 修改配置文件

```shell
# /etc/httpd/conf.d/httpd-vhosts.conf

# 增加虚拟主机配置
...
<VirtualHost *:443>
	DocumentRoot "/var/www/html" 
	ServerName 	system1.group8.example.com
	<Directory "/var/www/html"> 
		<RequireAll>
			Require all granted
			Require not host .my133t.org 
		</RequireAll>
	</Directory>
	
	SSLEngine on
	SSLProtocol all ‐SSLv2 ‐SSLv3
	SSLCertificateFile /etc/pki/tls/certs/system1.crt 	SSLCertificateKeyFile /etc/pki/tls/private/system1.key 	SSLCACertificateFile /etc/pki/tls/certs/ssl‐ca.crt
</VirtualHost>
...
```

##### 3. 下载证书

```shell
wget -O /etc/pki/certs/system1.crt http://证书链接
wget -O /etc/pki/tls/private/system1.key http://证书密钥链接
wget -O /etc/pki/tls/certs/ssl‐ca.crt http://证书签名授权信息链接
```

##### 4. 添加防火墙规则

```shell
firewall‐cmd ‐‐permanent ‐‐add‐service=https
firewall‐cmd ‐‐reload
```

##### 5. 重启服务

```shell
systemctl restart httpd
```

##### 6. 验证

```shell
# -k 忽略证书不受信问题
curl -k https://system1.group8.example.com
```

#### 难点

* `apache`配置文件添加`ssl`证书

### <a name="14">14. 配置虚拟主机</a>

在`system1`上扩展您的`web`服务器, 为站点`http://www8.group8.example.com`创建一个虚拟主机, 然后执行下述步骤:

* 设置`DocumentRoot`为`/var/www/virtual`
* 从`http://server.group8.exmaple.com/pub/wwww8.html`下载文件重名为`index.html`(不对文件做修改)
* 将文件`index.html`放到虚拟主机的`DocumentRoot`目录下
* 确保`andy`用户能够在`/var/www/virtual`目录下创建文件

源站点`http://system1.group8.example.com`必须仍然能够服务.

#### 答题步骤

##### 1. 创建网站目录并下载首页

```shell
mkdir /var/www/virtual
wget -O /var/www/virtual/index.html http://server.group8.exmaple.com/pub/wwww8.html
```

##### 2. 设定网站目录权限

```shell
setfacl -m u:andy:rwx /var/www/virtual
```

##### 3. 建立虚拟主机

```shell
# 修改配置文件 /etc/httpd/conf.d/httpd-vhosts.conf

...
<VirtualHost *:80>
	DocumentRoot "/var/www/virtual" 
	ServerName www8.group8.example.com
	<Directory "/var/www/virtual"> 
		<RequireAll>
		Require all granted 
		</RequireAll>
	</Directory> 
</VirtualHost>
...
```

##### 4. 重启服务

```shell
systemctl restart httpd
```

##### 5. 测试

在`system2`上

```shell
curl http://www8.group8.example.com
```

#### 难点

* 添加虚拟主机配置

### <a name="15">15. 配置 web 内容的访问</a>

在`system1`上的`web`服务器的`DocumentRoot`目录下, 创建一个名为`private`的目录, 要求：

* 从`http://server.group8.example.com/pub/private.html`下载一个文件副本到这个目录, 并且命名为`index.html`
* 不要对这个文件的内容做修改
* 从`system1`上, 任何人都可以浏览`private`的内容, 但是从其他系统不能访问这个目录的内容

#### 答题步骤

##### 1. 建立目录

```shell
mkdir /var/www/html/private
mkdir /var/www/virtual/private
```

##### 2. 下载页面

```shell
wget ‐O /var/www/html/private/index.html http://server.group8.example.com/pub/private.html
wget ‐O /var/www/virtual/private/index.html http://server.group8.example.com/pub/private.html
```

##### 3. 修改虚拟主机的配置

```shell
# /etc/httpd/conf.d/httpd‐vhosts.conf
# www8.group8.example.com 的配置文件添加如下信息

# 在 system1.group8.example.com
...
<Directory "/var/www/html/private"> 
	Require all denied
	Require local 
</Directory>
...

# www8.group8.example.com 
...
<Directory "/var/www/private/private"> 
	Require all denied
	Require local 
</Directory>
...
```

##### 4. 重启服务

```shell
systemctl restart httpd
```

##### 5. 测试

```shell
#在system1, 可以访问正常
curl http://system1.group8.example.com/private/

#在system2, 访问出现403
curl http://system1.group8.example.com/private/
```

#### 难点

* 配置信息
	* Require all denied
	* Require local 

### <a name="16">16. 实现动态 web 内容</a>

实现动态`WEB`内容

在`system1`上配置提供动态`web`内容, 要求:

* 动态内容由名为`wsgi.group8.example.com`的虚拟主机提供
* 虚拟主机侦听在端口`8909`
* 从`http://server.group8.example.com/pub/webinfo.wsgi`下载一个脚本, 放在适当位置, 不要修改
* 客户端访问`http://wsgi.group8.example.com:8909`时, 应该接收到动态生成的`web`页面
* 此`http://wsgi.group8.example.com:8909`, 必须能被`group8.example.com`域内的所有系统访问

#### 答题步骤

在`system1`上执行

##### 1. 建立虚拟主机

```shell
# /etc/httpd/conf.d/httpd‐vhosts.conf
...
Listen 8909
<VirtualHost *:8909>
	ServerName wsgi.group8.example.com
	WSGIScriptAlias / /var/www/html/webinfo.wsgi 
</VirtualHost>
```

##### 2. 下载页面

```shell
wget -O /var/www/html/webinfo.wsgi http://server.group8.example.com/pub/webinfo.wsgi
```

##### 3. 安装模块包

```shell
yum install mod_wsgi -y
```

##### 4. 添加防火墙规则

```shell
firewall-cmd --permanent ‐‐add‐rich‐rule 'rule family="ipv4" port port=8909 protocol=tcp accept'
firewall-cmd --reload
```

##### 5. 设定 SELinux

```shell
# SELinux 放行端口
semanage port -a -t http_port_t -p tcp 8909
```

##### 6. 重启服务

```shell
systemctl restart httpd
```

##### 7. 测试

```shell
curl http://wsgi.group8.example.com:8909
...
```

#### 难点

* `semanage port`针对端口做修改

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