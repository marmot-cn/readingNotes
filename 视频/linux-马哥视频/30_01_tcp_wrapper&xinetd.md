# 30_01_tcp_wrapper&xinetd

---

### 笔记

---

#### tcp wrapper

网络资源控制器. 工作在`tcp`层访问控制工具.

通常只能对基于`tcp`协议的去做控制.

`tcp wrapper`对应`libwrap.so`库文件. 不是一个服务, 是一个库文件. 工作在内核和应用程序之间. `内核 -> 库 -> 应用程序`. 应用程序提供了一个调用`libswarp.so`库的调用接口, 就可以受`tcp wrapper`控制.

检查二进制文件依赖于什么库, 命令`ldd`.

```shell
可见 libwrap.so 库
[root@iZ944l0t308Z ~]# ldd `which sshd`
	linux-vdso.so.1 =>  (0x00007fff85dfe000)
	libfipscheck.so.1 => /lib64/libfipscheck.so.1 (0x00007fbcfef70000)
	libwrap.so.0 => /lib64/libwrap.so.0 (0x00007fbcfed65000)
	libaudit.so.1 => /lib64/libaudit.so.1 (0x00007fbcfeb3c000)
	libpam.so.0 => /lib64/libpam.so.0 (0x00007fbcfe92d000)
	libselinux.so.1 => /lib64/libselinux.so.1 (0x00007fbcfe706000)
	libcrypto.so.10 => /lib64/libcrypto.so.10 (0x00007fbcfe320000)
	libdl.so.2 => /lib64/libdl.so.2 (0x00007fbcfe11c000)
	libldap-2.4.so.2 => /lib64/libldap-2.4.so.2 (0x00007fbcfdeca000)
	liblber-2.4.so.2 => /lib64/liblber-2.4.so.2 (0x00007fbcfdcba000)
	libutil.so.1 => /lib64/libutil.so.1 (0x00007fbcfdab7000)
	libz.so.1 => /lib64/libz.so.1 (0x00007fbcfd8a1000)
	libnsl.so.1 => /lib64/libnsl.so.1 (0x00007fbcfd687000)
	libcrypt.so.1 => /lib64/libcrypt.so.1 (0x00007fbcfd450000)
	libresolv.so.2 => /lib64/libresolv.so.2 (0x00007fbcfd236000)
	libgssapi_krb5.so.2 => /lib64/libgssapi_krb5.so.2 (0x00007fbcfcfe7000)
	libkrb5.so.3 => /lib64/libkrb5.so.3 (0x00007fbcfcd00000)
	libk5crypto.so.3 => /lib64/libk5crypto.so.3 (0x00007fbcfcace000)
	libcom_err.so.2 => /lib64/libcom_err.so.2 (0x00007fbcfc8c9000)
	libc.so.6 => /lib64/libc.so.6 (0x00007fbcfc508000)
	libcap-ng.so.0 => /lib64/libcap-ng.so.0 (0x00007fbcfc303000)
	libpcre.so.1 => /lib64/libpcre.so.1 (0x00007fbcfc0a1000)
	/lib64/ld-linux-x86-64.so.2 (0x00007fbcff421000)
	libsasl2.so.3 => /lib64/libsasl2.so.3 (0x00007fbcfbe84000)
	libssl3.so => /lib64/libssl3.so (0x00007fbcfbc46000)
	libsmime3.so => /lib64/libsmime3.so (0x00007fbcfba1e000)
	libnss3.so => /lib64/libnss3.so (0x00007fbcfb6f9000)
	libnssutil3.so => /lib64/libnssutil3.so (0x00007fbcfb4cd000)
	libplds4.so => /lib64/libplds4.so (0x00007fbcfb2c8000)
	libplc4.so => /lib64/libplc4.so (0x00007fbcfb0c3000)
	libnspr4.so => /lib64/libnspr4.so (0x00007fbcfae85000)
	libpthread.so.0 => /lib64/libpthread.so.0 (0x00007fbcfac68000)
	libfreebl3.so => /lib64/libfreebl3.so (0x00007fbcfa9eb000)
	libkrb5support.so.0 => /lib64/libkrb5support.so.0 (0x00007fbcfa7db000)
	libkeyutils.so.1 => /lib64/libkeyutils.so.1 (0x00007fbcfa5d7000)
	librt.so.1 => /lib64/librt.so.1 (0x00007fbcfa3ce000)
	
xinetd 也包含 libwrap.so
[root@iZ944l0t308Z ~]# ldd `which xinetd`
	linux-vdso.so.1 =>  (0x00007fffa6ffe000)
	libselinux.so.1 => /lib64/libselinux.so.1 (0x00007f2b8d1ce000)
	libwrap.so.0 => /lib64/libwrap.so.0 (0x00007f2b8cfc3000)
	libnsl.so.1 => /lib64/libnsl.so.1 (0x00007f2b8cda9000)
	libm.so.6 => /lib64/libm.so.6 (0x00007f2b8caa7000)
	libcrypt.so.1 => /lib64/libcrypt.so.1 (0x00007f2b8c870000)
	libc.so.6 => /lib64/libc.so.6 (0x00007f2b8c4ae000)
	libpcre.so.1 => /lib64/libpcre.so.1 (0x00007f2b8c24d000)
	libdl.so.2 => /lib64/libdl.so.2 (0x00007f2b8c049000)
	/lib64/ld-linux-x86-64.so.2 (0x00007f2b8d62b000)
	libfreebl3.so => /lib64/libfreebl3.so (0x00007f2b8bdcb000)
	libpthread.so.0 => /lib64/libpthread.so.0 (0x00007f2b8bbaf000)
```

`tcp warpper`可以理解为工作在用户请求和服务所监听的套接字之间的一种检查过滤机制.

`xinetd`接受`tcp warpper`, 其下所有控制的进程都受到`tcp warpper`控制.

##### 实现对一个服务的访问控制

* `/etc/hosts.allow`: 能够被访问的列表
* `/etc/hosts.deny`: 拒绝被访问的列表

首先检查`/etc/hosts.allow`是否明确允许, 否则继续检查`/etc/hosts.deny`，还不存在则按照默认规则:**允许**.

##### 文件的语法格式

`daemon_list: client_list [:options]`

`daemon_list`: 可执行程序的二进制文件名.

```shell
这个客户端列表在访问sshd服务 允许(allow文件)/拒绝(deny文件)
sshd: 192.168.0.

vsftpd,sshd,in.telnetd:

ALL(本机上所有接收tcp wrapper控制的服务)

daemon@host
假设本机有两个地址
1. 172.16.100.1
2. 192.168.0.186

来在于网段 1.0.0.0/255.0.0.0 访问在 192.168.0.186 地址的 vsftpd, 接收tcp wrapper控制
vsftpd@192.168.0.186: 1. (1.0.0.0/255.0.0.0)
```

`client_list`:

* `IP` 单个ip
* `network address`
	* `network/mask`: 不能使用位数格式, 1.0.0.0/255.0.0.0(可以) 1.0.0.0/8(不可以)
	* 172.16.
* `HOSTNAME` 主机名
	* `fqdn`
	* .a.org (*.a.org 都可以访问)
* `MACRO` 宏
	* `ALL` 本机上所有接收tcp wrapper控制的服务
	* `LOCAL` 本地来宾, 和本地网卡在同一个网段的主机
	* `KNOWN` fqdn 可以被正常解析的
	* `UNKNOWN` fqdn 无法解析的
	* `PARANOID` 主机名正向解析和反向解析不匹配(www.xxx.com正向解析地址是1.1.1.1, 而1.1.1.1反向解析的地址是www.xxx.com)
	* `EXCEPT` 除了,不包含
	
`options`:

* `ALLOW`文件中可以`DENY`, `DENY`文件中可以`ALLOW`
* `spawn` 日志
	
```shell
/etc/hosts.allow
in.telnetd: 172.16. :DENY

当172.16.网段主机访问记录日志
/etc/hosts.allow
in.telnetd: 172.16. :spawn echo "somebody enteted, `date`" >> /var/log/tcpwrapper.log
```

##### 宏

`man 5 hosts_access`获取

* `%c` 客户端信息: client information(user@host)
* `%s` 主机上的哪一个服务: service info(service@host)
* `%h` 客户端主机名: client hostname
* `%p` 明确说明哪个服务: server PID
	
##### 示例

**示例1**

sshd 服务, 仅允许 172.16.0.0 网段访问

```shell
/etc/hosts.allow
sshd: 172.16.

/etc/hosts.deny
sshd: ALL
```

**示例2**

telnet 服务, 不允许 172.16.0.0 网段访问, 但是允许 172.16.100.200 访问, 其他客户端不做控制,

```shell
方法一:
/etc/hosts.allow
in.telnetd: 172.16.100.200

/etc/hosts.deny
in.telnetd: 172.16.

方法二:
/etc/hosts.deny
in.telnetd: 172.16. EXCEPT 172.16.100.200

方法三:
/etc/hosts.allow
in.telnetd: ALL EXCEPT 172.16. EXCEPT 172.16.100.200

/etc/hosts.deny
in.telnetd: ALL
```

**示例3**

主机有两块网卡:

1. 10.170.148.109
2. 120.25.87.35

禁止所有主机从`10.170.148.109`使用`sshd`服务.

```shell
[root@iZ944l0t308Z ~]# cat /etc/hosts.deny
...
sshd@10.170.148.109:ALL
...

重启 xinetd
service xinetd restart

另外一台主机:

一开始可以访问
[root@iZ94ebqp9jtZ ~]# ssh root@10.170.148.109
root@10.170.148.109's password:
Last login: Mon Jul 24 14:26:50 2017 from 10.116.138.44

Welcome to aliyun Elastic Compute Service!

禁止后:
[root@iZ94ebqp9jtZ ~]# ssh root@10.170.148.109
sh_exchange_identification: read: Connection reset by peer

但是另外一个ip可以登录
[root@iZ94ebqp9jtZ ~]# ssh root@120.25.87.35
The authenticity of host '120.25.87.35 (120.25.87.35)' can't be established.
ECDSA key fingerprint is ab:e5:65:fd:39:c2:e4:57:be:61:a0:00:ed:c0:16:df.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '120.25.87.35' (ECDSA) to the list of known hosts.
root@120.25.87.35's password:
Last failed login: Mon Jul 24 14:30:31 CST 2017 from pool-74-102-79-239.nwrknj.fios.verizon.net on pts/1
There were 2 failed login attempts since the last successful login.
Last login: Mon Jul 24 14:28:20 2017 from 10.116.138.44

Welcome to aliyun Elastic Compute Service!

[root@iZ944l0t308Z ~]#
```

### 整理知识点

---