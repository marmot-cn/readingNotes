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
* [9.  升级系统内核](#9)
* [10. 绑定验证服务](#10)
* [11. 配置 autofs](#11)
* [12. 配置 NTP](#12)
* [13. 创建一个归档](#13) **简单**
* [14. 配置一个用户账户](#14) **简单**
* [15. 创建一个 swap 分区](#15)
* [16. 查找文件](#16) **中等**
* [17. 查找一个字符串](#17) **简单**
* [18. 创建一个逻辑卷](#18)

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

#### 难点

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