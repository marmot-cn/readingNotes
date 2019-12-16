# RHCSA考试题解

## 目录

* [1.修改系统密码, 完成网络配置](#1) **难**
* [2. 配置Selinux](#2) **简单**
* [3. 配置 YUM 仓库](#3) **简单**
* [4. 调整逻辑卷容量](#4) **难**
* [5. 创建用户和用户组](#5) **简单**
* [6. 配置文件`/var/tmp/fstab`的权限](#6) **简单**
* [7. 建立计划任务](#7) **简单**
* [8. 创建一个共享目录](#8) **简单**
* [9. 升级系统内核](#9) **简单**
* [10. 绑定验证服务](#10) **中等**
* [11. 配置 autofs](#11) **中等**
* [12. 配置 NTP](#12) **简单**
* [13. 创建一个归档](#13) **简单**
* [14. 配置一个用户账户](#14) **简单**
* [15. 创建一个 swap 分区](#15) **难**
* [16. 查找文件](#16) **中等**
* [17. 查找一个字符串](#17) **简单**
* [18. 创建一个逻辑卷](#18) **难**

### <a name="1">1. 重置系统密码, 完成网络配置</a>

* 修改`root`账户密码
* 配置主机名
* 配置`IP`地址
* 配置子网掩码
* 配置网关
* 配置DNS

#### 答题步骤

##### 1.1 重启系统

##### 1.2 同时按下`Ctrl + Alt +Del`

##### 1.3 倒计时读秒结束前, 按下任意键

##### 1.4 选中默认的第一个启动条目, 并按下键盘字母`e`键, 进入**编辑**

##### 1.5 找到第一个以`linux16`开头的行, 定位到行尾, 附加**<空格>rd.break**

##### 1.6 按下`Ctrl +x`使用这些更改启动系统, 进入临时内核`shell`界面

##### 1.7 因为实际系统的`root`文件系统会在`/sysroot`中以只读方式挂载. 以读写形式重新挂载`/sysroot`

```shell
mount -o remount,rw / /sysroot
```

#### 1.8 切换 chroot 存放位置, 其中`/sysroot`被视为文件系统树的`root`

```shell
chroot /sysroot
```

#### 1.9 设置新的 root 密码

```shell
passwd root
```

#### 1.10 下次重启对文件进行`Selinux`上下文重新打标记

```shell
touch /.autorelabel
```

使selinux生效. 有这个文件存在, 系统统在重启时就会对整个文件系统进行relabeling.

#### 1.11 退出真实系统根目录

```shell
exit
```

#### 1.12 网络配置

查看本机的网卡设备

```shell
nmcli device
```

查看本机的连接配置

```shell
nmcli connection
```

修改名为`eth0`的连接配置信息

```shell
nmcli connection modify eth0 ipv4.address xxx.xxx.xxx.xx/xx ipv4.gateway xxx.xxx.xx.xx ipv4.dns xxx.xxx.xxx.xxx ipv4.method manual connection.autoconnect yes connection.interface-name eth0
```

重新加载连接配置信息

```shell
nmcli connection reload
```

重启链接

```shell
nmcli connection down eth0
nmcli connection up eth0
```

`ping`网关来测试网络连通性

```shell
ping xxxx.xxx.xxx.xxx
```

#### 1.13 设定主机名

```shell
hostnamectl set-hostname xxx
```

#### 难点

`1.5`,`1.6`,`1.7`,`1.8`和`1.10`为难点 **修改密码**

`1.12` 修改网络配置

`1.13` 设定主机名

### <a name="2">2. 配置Selinux</a>

* `SeLinux`的工作模式为`enforcing`
* 系统重启后依然生效

#### 答题步骤

##### 2.1 临时开启

检查现有配置

```shell
getenforce
```

修改配置

```shell
setenforce 1
```

再次检查

```shell
getenforce
```

##### 2.2 永久设置

修改`/etc/sysconfig/selinux`

找见`SELINUX=xxx`修改为`SELINUX=enforcing`

#### 难点

* `setenforce`命令
* `/etc/sysconfig/selinux`配置文件

### <a name="3">3. 配置 YUM 仓库</a>

* 修改仓库地址

#### 答题步骤

##### 3.1 创建仓库配置文件

在`/etc/yum.repo.d/`下创建一个`xxx.repo`

```
[仓库名字]
name=仓库名字
baseurl=仓库地址
enabled=1
gpgcheck=0
```

##### 3.2 检查

执行命令检查配置是否成功

```shell
yum repolist
```

#### 难点

* 仓库配置文件的创建
	* name
	* baseurl **仓库地址**
	* enabled=1
	* gpgcheck=0

### <a name="4">4. 调整逻辑卷容量</a>

* 调整本地逻辑卷`lv0`的容量
	* 调整文件系统大小
	* 调整后确保文件系统已存在的内容不能被破坏
	* 调整后的容量可能出现误差
	* 调整后, 保证其挂载目录不改变, 文件系统完成

#### 答题步骤

##### 4.1 检查文件系统大小

```shell
df -hT
```

* `h` 方便阅读方式显示
* `T` 显示文件系统类型

##### 4.2 扩容逻辑卷

```shell
lvextend -L 290M /dev/vg0/lv0
```

检查是否扩容成功

```shell
lvs
```

#### 4.3 刷新文件系统容量

```shell
resize2fs /dev/vg0/lv0
```

再次检查

```shell
df -hT
```

#### 难点

* 扩容逻辑卷命令 `lvextend`
* 刷新文件系统容量 `resize2fs`

### <a name="5">5. 创建用户和用户组</a>

* 创建一个名为`adminuser`的用户组, 组`id`为`40000`
* 创建一个名为`natasha`的用户组, 并将`adminuser`作为其附属组
* 创建一个名为`harry`的用户组, 并将`adminuser`作为其附属组
* 创建一个名为`sarah`的用户组, 其不属于`adminuser`组, 其在系统中没有任何可交互的`shell`
* `natasha`,`harry`和`sarah`三个用户的密码均设置为`redhat`

#### 答题步骤

```shell
groupadd -g 40000 adminuser
useradd -G adminuser natasha
useradd -G adminuser harry
useradd -s /sbin/nologin sarah
echo redhat | passwd --stdin natasha
echo redhat | passwd --stdin harry
echo redhat | passwd --stdin sarah
```

#### 难点

* 创建组命令`groupadd`
* 创建用户命令`useradd`
* 创建用户命令, 且不能登录`useradd -s /sbin/nologin`

### <a name="6">6. 配置文件`/var/tmp/fstab`的权限</a>

* 复制文件`/etc/fstab`到`/var/tmp`下, 并配置权限
	* 文件所属人为`root`
	* 文件的所属组为`root`
	* 文件对任何人均没有执行权限
	* 用户`natasha`对该文件有读和写的权限
	* 用户`harry`对该文件既不能读也不能写
	* 所有其他用户对该文件都有读的权限

#### 答题步骤

```
cp /etc/fstab /var/tmp
chonw root:root /var/tmp
chmod a-x /var/tmp

# 通过acl设定 natasha 和 harry 权限
setfacl -m u:natasha:rw,u:harry:--- /var/tmp/fstab

# 检查
getface /var/tmp/fstab

# 如果其他用户没有对该文件有读的权限, 可以设置. 答案里面没有该设置
chmod o+r /var/tmp
```

#### 难点

* `setface` 命令

### <a name="7">7. 建立计划任务</a>

对`natasha`用户建立计划任务, 要求在本地时间的每天`14:23`执行以下命令

```shell
/bin/echo "rchsa"
```

#### 答题步骤

```shell
crontab -e -u natasha
23 14 * * * /bin/echo  "rhcsa"
```

#### 难点

`crontab`中的`*`依次是:

* 分钟
* 小时
* 每天
* 每月
* 每星期

### <a name="8">8. 创建一个共享目录</a>

 在`/home`目录下创建名为`admins`的子目录, 并按以下需求设置权限:
 
 * `/home/admins`目录的所属组为`adminuser`
 * 该目录对`adminuser`组的成员可读可执行可写, 但对其他用户没有任何权限, 但`root`不受限制
 * 在`/home/admins`目录下创建的所属组自动被设置为`adminuser`
 
#### 答题步骤

```shell
mkdir -p /home/admins
chown :adminuser /home/admins 或者 chgrp adminuser /home/admins
chmod grwx,o=--- /home/admins
chmod g+s /home/admins
```

#### 难点

* `g+s`: 该权限只对目录有效. 目录被设置该位后, 任何用户在此目录下创建的文件都具有和该目录所属的组相同的组

### <a name="9">9. 升级系统内核</a>

从`xxx`下找到需要升级的内核

* 当系统重新启动之后, 升级的内核要作为默认的内核
* 原来的内核要保留, 并且仍然可以正常启动

#### 答题步骤

##### 9.1 安装内核文件

1. 通过火狐浏览器打开链接, 手动下载`kernel-xxx`内核文件, 然后使用`rpm -ivh *.rpm`命令安装
2. 通过`curl --silent xxx | grep kenel`过滤出内核`rpm`包地址, 使用`yum install -y xxxx.rpm`远程安装

##### 9.2 验证

`grub2-editenv list`

#### 难点

* `curl --silent xxx | grep kenel`远程过滤
* 手动安装`rpm -ivh *.rpm`(也可以使用`rpm -uvh *.rpm`)
* 自动安装`yum install -y 远程地址.rpm`
* 查看默认启动项`grub2-editenv list`
* 修改默认启动项`grub2-set-default`

### <a name="10">10. 绑定外部验证服务</a>

系统`server.group8.example.com`提供了一个`LDAP`验证服务. 系统新药按照以下要求绑定到这个服务上:

* 验证服务器的基本`DN`是: `dc=group8, dc=example, dc=com`
* 账户信息和验证信息都是由`LDAP`提供
* 链接需要使用证书加密, 证书可以在下面的链接下`证书下载链接/cacert.crt`

正确配置完成后, 用户`thales`可以登录系统, 登录密码是`redhat`

#### 答题步骤

##### 10.1 下载相关组件

```shell
yum install authconfig-gtk sssd -y
```

##### 10.2 配置相关组件

```
authconfig-gtk &
```

配置:

* DN
* Server
* 证书

##### 10.3 测试

```shell
su - thales
```

或者

```shell
# 相当于 cat /etc/passwd|grep thales
getent passwd thales
```

如果无法登录, 先完成`ntp`时间同步, 在重启`sssd`, 并确认`sssd`启动正常.

#### 难点

* `authconfig`配置以及相关组件安装
* `sssd`一个守护进程，该进程可以用来访问多种验证服务器, 如LDAP, Kerberos等, 并提供授权. SSSD是介于本地用户和数据存储之间的进程, 本地客户端首先连接SSSD, 再由SSSD联系外部资源提供者(一台远程服务器).

### <a name="11">11. 配置autofs</a>

按照下述要求配置`autofs`用来自动挂载`DLAP`用户的主目录.

* `server.group8.example.com`通过`NFS`输出了`/rhome`目录到您的系统. 这个文件系统包含了用户`thales`的主目录, 并且已经预先配置好了.
* `thales`用户的主目录是`server.group8.example.com:/rhome/thales`
* `thales`用户的主目录应该挂载到本地的`/home/ldap/thales`
* 用户对其主目录必须是读写的
* `thales`的登录密码是`redhat`

#### 答题步骤

##### 11.1 安装服务

```shell
yum install autofs -y
```

##### 11.2 修改配置文件

```shell
vim /etc/auto.master
# 当系统访问以/home/ldap路径开头的资源时, 读取/etc/.ldap配置文件进行自动挂载
/home/ldap	/etc/auto.ldap
```

```shell
vim /etc/auto.ldap

# 配置含义: 当访问 /home/ldap/* (* 代表任意路径)的资源时, 自动挂载到server.group8.example.com:/rhome/对应路径的资源
*	-rw,sync,soft  server.group8.example.com:/rhome/&
```

##### 11.3 开机启动服务

```shell
systemctl enable autofs
systemctl restart autofs
```

##### 11.4 验证

```shell
ssh thales@localhost
pwd
/home/ldap/thales
```

#### 难点

* 理解配置文件

### <a name="12">12. 配置NTP</a>

配置系统时间与服务器`server.group8.example.com`同步, 要求系统重启后依然生效.

#### 答题步骤

##### 12.1 修改`/etc/chrony.con`配置文件

注释掉原来的同步时间服务器

```shell
# server 0.rhel.pool.ntp.org iburst
# server 1.rhel.pool.ntp.org iburst
# ...
```

添加新的同步时间服务器

```
server server.group8.example.com iburst
```

##### 12.2 配置服务

```shell
stsremctl enable chronyd
systemctl restart chronyd
# 手动同步时间
chronyc

chronyc>waitsync
```

#### 难点

* 修改配置文件`/etc/chrony.con`
* 手动同步时间`chronyc`

### <a name="13">13. 创建一个归档</a>

创建一个名为`/root/sysconfig.tar.bz2`的归档文件, 其中包含了`/etc/sysconfig`目录中的内容.

#### 答题步骤

```shell
tar -jcf /root/sysconfig.tar.bz2 /etc/susconfig/
```

#### 难点

* `tar -jcf` 使用 `bz2`
* `tar -zxf` 使用 `gzip`

### <a name="14">14. 配置一个用户账户</a>

创建一个名为`jay`的用户.

* 用户`id`为`3456`
* 密码为`gleunge`

#### 答题步骤

```shell
useradd -u 3456 jay
echo gleunge | passwd --stdin jay 
```

#### 难点

* `useradd -u`, 添加用户并设定用户`uid`
* `passwd --stdin`, 从`stdin`输入密码

### <a name="15">15. 添加一个 swap 分区</a>

添加一个新的`swap`分区

* `swap`分区容量为`512 MiB`
* 系统启动时, `swap`分区应该可以自动挂载
* 不要移除或者修改其他已经存在于您的系统中的`swap`分区

#### 答题步骤

##### 15.1 格式化磁盘

```shell
fdisk /dev/xxx
n
e #使用扩展分区
回车
回车
# 创建完成扩展分区

n
回车
+512M

t
分区号
82 #82代表swap
w #保存

partprobe
# 如果生成设备文件失败可以使用 partx -a /dev/sda 再次生成
```

##### 15.2 格式化分区

```shell
# 格式化分区, 生成 UUID, 再次查看UUID 使用 blkid /dev/xxx
mkswap /dev/xxx

# 开启swap
swapon /dev/xxx

# 开机启动 /etc/fstab

UUID(也可以使用/dev/xxx) swap swap defaults 0 0
```

#### 难点

* `fdisk`
	* 创建逻辑分区
	* 格式化一个确认大小的磁盘
* `mkswap` 格式化分区
* `/etc/fstab` 开机自动挂载 
* `partprobe` 用于重读分区表, 当出现删除文件后, 出现仍然占用空间. 可以`partprobe`在不重启的情况下重读分区.
* `partx` 告诉内核当前磁盘的分区情况
	* `-a`: 增加指定的分区或读磁盘新增的分区
* `blkid` 对查询设备上所采用文件系统类型进行查询, 可以查看UUID

### <a name="16">16. 查找文件</a>

把系统上拥有者为`jay`用户的所有文件拷贝到`/root/findfiles`目录中

#### 答题步骤

```shell
mkdir -p /root/findfiles
find / -user jay -exec cp -a {} /root/finedfiles/ \l
```

#### 难点

* `-user` 按文件属主来查找
* `-exec`  find命令对匹配的文件执行该参数所给出的shell命令, 相应命令的形式为`'command' { } \;`，注意{ }和\；之间的空格
	* 花括号代表前面find查找出来的文件名.

### <a name="17">17. 查找一个字符串</a>

把`/user/share/dict/words`文件中所包含`seismic`字符串的行找到, 并将这些按照原始文件中的顺序存放到`/root/wordlist`中, `/root/wordlist`文件中不能包含换行.

#### 答题步骤

```shell
grep seismic /usr/share/dict/words > /root/wordlist
```

#### 难点

* `grep`命令

### <a name="18">18. 创建一个逻辑卷</a>

创建一个新的逻辑卷

* 创建一个名为`datastore`的卷组, 卷组`PE`尺寸为`16 MiB`
* 逻辑卷的名字为`database`, 所属卷组为`datastore`, 该逻辑卷由`50`个`PE`组成
* 将新建的逻辑卷格式化为`xfs`文件系统, 要求系统启动时, 该逻辑卷能被自动挂碍到`/mnt/database`

#### 答题步骤

##### 18.1 格式化磁盘

```shell
fdisk /dev/sda
n #这里自动使用上个扩展分区
回车
回车 #用完所有控件

t #修改分区类型标记
8e #修改类型为 LVM
w #保存分区设定

partprobe 
```

##### 18.2 创建逻辑卷

```shell
pvcreate /dev/sda#

# 创建卷组
vgcreate -s 16M datastore /dev/sda#
vgdisplay

# 创建逻辑卷, -l 50 代表使用50个PE的容量 
lvcreate -n database -l 50 datastore
lvs

# 格式化操作系统
mkfs.xfs /dev/datastore/database
mkdir -p /mnt/database

# 自动化挂载
# 查看 UUID
blkid /dev/datastore/database
# 修改/etc/fstab 自动化挂载

UUID="xxx" /mnt/dayabase xfs defaults 0 0

# 挂载
mount -a

# 确认已挂载上
df -h
```

#### 难点

* `vgcreate` 创建卷组
	* `-s` 卷组上的物理卷`PE`大小
* `lvcreate` 创建逻辑卷
* `mks.xfs` 格式化操作系统