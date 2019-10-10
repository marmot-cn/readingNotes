# 29_04_samba之二

---

### 笔记

---

#### smbclient 访问共享

```shell
安装 smb 客户端
[root@iZ94ebqp9jtZ ~]# yum install -y samba-client
```

**参数**

* `-L NetBios_Name` 指定netbios主机名, 跟`ip`地址也可以
* `-U Username` 指定用户名
* `-P` 输入密码

**示例**

```shell
我测试的访问另外一台linux的samba共享目录, 没有链接windwos服务器

[root@iZ94ebqp9jtZ ~]# smbclient -L 120.25.87.35 -U chloroplast
Enter chloroplast's password:
Domain=[SAMBA] OS=[Windows 6.1] Server=[Samba 4.4.4]

	Sharename       Type      Comment
	---------       ----      -------
	print$          Disk      Printer Drivers
	tools           Disk      Share testing
	IPC$            IPC       IPC Service (Samba 4.4.4)
	chloroplast     Disk      Home Directories
Domain=[SAMBA] OS=[Windows 6.1] Server=[Samba 4.4.4]

	Server               Comment
	---------            -------

	Workgroup            Master
	---------            -------
	
直接访问 samba 服务器
[root@iZ94ebqp9jtZ ~]# smbclient //120.25.87.35 -U chloroplast

\\120.25.87.35: Not enough '\' characters in service
Usage: smbclient [-?EgBVNkPeC] [-?|--help] [--usage] [-R|--name-resolve NAME-RESOLVE-ORDER] [-M|--message HOST] [-I|--ip-address IP] [-E|--stderr] [-L|--list HOST] [-m|--max-protocol LEVEL] [-T|--tar <c|x>IXFqgbNan]
        [-D|--directory DIR] [-c|--command STRING] [-b|--send-buffer BYTES] [-t|--timeout SECONDS] [-p|--port PORT] [-g|--grepable] [-B|--browse] [-d|--debuglevel DEBUGLEVEL] [-s|--configfile CONFIGFILE]
        [-l|--log-basename LOGFILEBASE] [-V|--version] [--option=name=value] [-O|--socket-options SOCKETOPTIONS] [-n|--netbiosname NETBIOSNAME] [-W|--workgroup WORKGROUP] [-i|--scope SCOPE] [-U|--user USERNAME] [-N|--no-pass]
        [-k|--kerberos] [-A|--authentication-file FILE] [-S|--signing on|off|required] [-P|--machine-pass] [-e|--encrypt] [-C|--use-ccache] [--pw-nt-hash] service <password>
[root@iZ94ebqp9jtZ ~]# smbclient //120.25.87.35/tools -U chloroplast
Enter chloroplast's password:
Domain=[SAMBA] OS=[Windows 6.1] Server=[Samba 4.4.4]
smb: \> ls
  .                                   D        0  Sat Jul 22 15:41:03 2017
  ..                                  D        0  Sat Jul 22 14:12:24 2017
  222.txt                             N        3  Sat Jul 22 15:41:10 2017
  111                                 N        0  Sat Jul 22 15:40:53 2017

		20510332 blocks of size 1024. 17599396 blocks available
smb: \> get 222.txt
getting file \222.txt of size 3 as 222.txt (0.1 KiloBytes/sec) (average 0.1 KiloBytes/sec)
smb: \> exit
[root@iZ94ebqp9jtZ ~]# ls
222.txt
[root@iZ94ebqp9jtZ ~]# cat 222.txt
33

把共享目录挂载到mnt下
[root@iZ94ebqp9jtZ ~]# mount -t cifs //120.25.87.35/tools /mnt -o username=chloroplast
mount: //120.25.87.35/tools is write-protected, mounting read-only
mount: cannot mount //120.25.87.35/tools read-only

挂载提示不能挂载上去, 后查询需要下载 cifs-utils
[root@iZ94ebqp9jtZ mnt]# yum install cifs-utils -y
...

可以挂载上了, 可以向本地文件系统一样使用samba
[root@iZ94ebqp9jtZ mnt]# mount -t cifs //120.25.87.35/tools /mnt -o username=chloroplast
Password for chloroplast@//120.25.87.35/tools:  **********
[root@iZ94ebqp9jtZ mnt]# ls /mnt/
111  222.txt
可以写入文件
[root@iZ94ebqp9jtZ mnt]# touch 333

实现开机自动挂载, 写入 fstab. 我们不能把密码直接写入到 /etc/fstab 文件中去, 所以我们把账户和密码写到/etc/samba/crde.passwd(只有root有访问权限), 使用 credentials=/etc/samba/crde.passwd
[root@iZ94ebqp9jtZ mnt]# vim /etc/fstab
...
//120.25.87.35/tools /mnt	cifs	credentials=/etc/samba/crde.passwd 0 0

[root@iZ94ebqp9jtZ mnt]# cat /etc/samba/crde.passwd
username=chloroplast
password=19831030dK

设定 og 没有任何权限
[root@iZ94ebqp9jtZ mnt]# chmod og=--- /etc/samba/crde.passwd
卸载刚才挂的samba 
[root@iZ94ebqp9jtZ mnt]# umount /mnt
[root@iZ94ebqp9jtZ ~]# mount -a
[root@iZ94ebqp9jtZ ~]# df -h
Filesystem            Size  Used Avail Use% Mounted on
/dev/xvda1             40G  1.6G   36G   5% /
devtmpfs              487M     0  487M   0% /dev
tmpfs                 496M     0  496M   0% /dev/shm
tmpfs                 496M  292K  496M   1% /run
tmpfs                 496M     0  496M   0% /sys/fs/cgroup
tmpfs                 100M     0  100M   0% /run/user/0
//120.25.87.35/tools   20G  1.8G   17G  10% /mnt
[root@iZ94ebqp9jtZ ~]# ls /mnt/
111  222.txt  333
```

#### 示例

新建一个共享, 共享名为tools, 开放给组mygrp中的所有用户具有读写权限, 其他用户只有读权限.

首先需要保证`mygrp`对目录具有读写权限.

```shell
[tools]
	comment = 
	path = 
	public = yes
	write list = @mygrp 
```

#### samba 基于ip的访问

iptables 开放 tcp 139,445 端口. udp 137,138 端口.

samba配置文件 `hosts allow = 172.16. 127.` (允许172.16. 和 127. 网段的主机访问)白名单

#### samba-swat

基于 WEB GUI 的图形配置界面

#### 守护进程

* standalone(独立的, 自己管理所有事情)
* transient(瞬时进程): 自己平时不启动, 用`xinetd`(超级守护进程,可以同时为多个进程提供服务)代理, 

`xinetd`为本身平时不经常访问的进程代为监听端口.

`samba-swat`监听在`901/tcp`, 由`xinetd`代为监听. 当请求过来, `xinetd`临时的将`swat`启动进来, 并将请求转交给它.

`xinted`是独立守护进程由运行级别概念. 瞬时守护进程没有运行级别概念.

**开启瞬时守护进程方法**

`chkconfig xxx on` 瞬时守护进程不能设定级别. 重启`xinetd`然后可以代理监听.

每一个瞬时守护进程都有一个独立的配置文件: `/etc/xint.d/`目录中, 名字和服务名字相同的配置文件.

配置参数:

* `disable` = no (是否禁用, yes 代表禁用, 直接修改该参数然后重启服务也可以打开服务)
* `port` = xx 监听在哪个端口
* `socket_type` = stream(tcp) dgram(udp)
* `wait` = no 是不是可以两个以上的用户同时访问, yes 就是等待. no 就是可以同时访问, 启动多个进程.
* `only_from` = 127.0.0.1 白名单(默认是本机, 也可以设置为一个网段)
* `user` = root 以root用户身份运行服务
* `server` = /usr/sbin/swat 执行的程序
* `log_on_failure` += USERID 一旦产生错误就在原有记录上记录用户的uid(所以是`+=`)

### 整理知识点

---