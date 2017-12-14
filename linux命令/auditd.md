# auditd

## 简介

`auditd`是`Linux`审计系统中用户空间的一个组件, 负责**将审计记录写入磁盘**.

## 相关工具

安装完成后会有如下相关工具.

* `audictl`: 即时控制审计守护进程的行为的工具, 添加规则等.
* `/etc/audit/audit.rules`: 记录审计规则的文件.
* `aureport`: 查看和生成审计报告工具.
* `ausearch`: 查找审计事件的工具.
* `auditspd`: 转发事件通知给其他应用程序, 而不是写入到审计日志文件中.
* `autrace`: 跟踪进程的命令.
* `/etc/audit/auditd.conf`: `auditd`工具的配置文件.

```shell
[root@localhost ansible]# auditctl -l
No rules
```

## 如何使用

### 文件和目录访问审计

#### 文件审计

```shell
sudo auditctl -w /etc/passwd -p rwxa
```

选项:

* `-w path`: 只要要监控的路径, 上述监控`/etc/passwd`文件
* `-p`: 指定触发审计的文件/目录的访问权限.
* `rwxa`: 指定要触发条件, `r`读取权限, `w`写入权限, `x`执行权限, `a`属性(attr)改变

#### 目录审计

```shell
sudo auditctl -w /data/

[root@localhost data]# auditctl -l
-w /etc/passwd -p rwxa
-w /data -p rwxa
```

### 查看审计日志

添加规则后, 我们可以查看`auditd`的日志. 使用`ausearch`工具可以查看`auditd`日志.

```shell
sudo ausearch -f /etc/passwd
```

* `-f`意味着基于文件名搜索.


为了演示我添加一个而用户, 然后查看审计日志.

```shell
[root@localhost data]# id tom
id: tom: no such user
[root@localhost data]# useradd tom
[root@localhost data]# ausearch -f /etc/passwd
----
time->Mon Nov 27 19:20:19 2017
type=PATH msg=audit(1511781619.598:6034): item=0 name="/etc/passwd" inode=18757551 dev=fd:00 mode=0100644 ouid=0 ogid=0 rdev=00:00 objtype=NORMAL
type=CWD msg=audit(1511781619.598:6034):  cwd="/data"
type=SYSCALL msg=audit(1511781619.598:6034): arch=c000003e syscall=2 success=yes exit=3 a0=7ffa773dd432 a1=80000 a2=1b6 a3=24 items=1 ppid=16304 pid=16329 auid=1001 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=pts3 ses=255 comm="id" exe="/usr/bin/id" key=(null)
----
time->Mon Nov 27 19:20:23 2017
type=PATH msg=audit(1511781623.358:6035): item=0 name="/etc/passwd" inode=18757551 dev=fd:00 mode=0100644 ouid=0 ogid=0 rdev=00:00 objtype=NORMAL
type=CWD msg=audit(1511781623.358:6035):  cwd="/data"
type=SYSCALL msg=audit(1511781623.358:6035): arch=c000003e syscall=2 success=yes exit=4 a0=7fd9d8f81432 a1=80000 a2=1b6 a3=24 items=1 ppid=16304 pid=16330 auid=1001 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=pts3 ses=255 comm="useradd" exe="/usr/sbin/useradd" key=(null)
----
time->Mon Nov 27 19:20:23 2017
type=PATH msg=audit(1511781623.359:6036): item=0 name="/etc/passwd" inode=18757551 dev=fd:00 mode=0100644 ouid=0 ogid=0 rdev=00:00 objtype=NORMAL
type=CWD msg=audit(1511781623.359:6036):  cwd="/data"
type=SYSCALL msg=audit(1511781623.359:6036): arch=c000003e syscall=2 success=yes exit=5 a0=7fd9e18adce0 a1=20902 a2=0 a3=0 items=1 ppid=16304 pid=16330 auid=1001 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=pts3 ses=255 comm="useradd" exe="/usr/sbin/useradd" key=(null)
----
time->Mon Nov 27 19:20:23 2017
type=PATH msg=audit(1511781623.359:6037): item=0 name="/etc/passwd" inode=18757551 dev=fd:00 mode=0100644 ouid=0 ogid=0 rdev=00:00 objtype=NORMAL
type=CWD msg=audit(1511781623.359:6037):  cwd="/data"
type=SYSCALL msg=audit(1511781623.359:6037): arch=c000003e syscall=2 success=yes exit=8 a0=7fd9d8f81432 a1=80000 a2=1b6 a3=24 items=1 ppid=16304 pid=16330 auid=1001 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=pts3 ses=255 comm="useradd" exe="/usr/sbin/useradd" key=(null)
----
time->Mon Nov 27 19:20:23 2017
type=PATH msg=audit(1511781623.359:6039): item=0 name="/etc/passwd" inode=18757551 dev=fd:00 mode=0100644 ouid=0 ogid=0 rdev=00:00 objtype=NORMAL
type=CWD msg=audit(1511781623.359:6039):  cwd="/data"
type=SYSCALL msg=audit(1511781623.359:6039): arch=c000003e syscall=2 success=yes exit=10 a0=7fd9d8f81432 a1=80000 a2=1b6 a3=24 items=1 ppid=16304 pid=16330 auid=1001 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=pts3 ses=255 comm="useradd" exe="/usr/sbin/useradd" key=(null)
----
time->Mon Nov 27 19:20:23 2017
type=CONFIG_CHANGE msg=audit(1511781623.364:6041): auid=1001 ses=255 op="updated_rules" path="/etc/passwd" key=(null) list=4 res=1
----
time->Mon Nov 27 19:20:23 2017
type=PATH msg=audit(1511781623.364:6042): item=4 name="/etc/passwd" inode=19347849 dev=fd:00 mode=0100644 ouid=0 ogid=0 rdev=00:00 objtype=CREATE
type=PATH msg=audit(1511781623.364:6042): item=3 name="/etc/passwd" inode=18757551 dev=fd:00 mode=0100644 ouid=0 ogid=0 rdev=00:00 objtype=DELETE
type=PATH msg=audit(1511781623.364:6042): item=2 name="/etc/passwd+" inode=19347849 dev=fd:00 mode=0100644 ouid=0 ogid=0 rdev=00:00 objtype=DELETE
type=PATH msg=audit(1511781623.364:6042): item=1 name="/etc/" inode=16777281 dev=fd:00 mode=040755 ouid=0 ogid=0 rdev=00:00 objtype=PARENT
type=PATH msg=audit(1511781623.364:6042): item=0 name="/etc/" inode=16777281 dev=fd:00 mode=040755 ouid=0 ogid=0 rdev=00:00 objtype=PARENT
type=CWD msg=audit(1511781623.364:6042):  cwd="/data"
type=SYSCALL msg=audit(1511781623.364:6042): arch=c000003e syscall=82 success=yes exit=0 a0=7ffe6c0751a0 a1=7fd9e18adce0 a2=7ffe6c075110 a3=9 items=5 ppid=16304 pid=16330 auid=1001 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=pts3 ses=255 comm="useradd" exe="/usr/sbin/useradd" key=(null)
```

输出结果:

* `time`: 审计时间.
* `name`: 审计对象.
* `cwd`: 当前路径.
* `syscall`: 相关的系统调用.
* `auid`: 审计用户ID.
* `uid`和`gid`: 访问文件的用户ID和用户组ID.
* `comm`: 用户访问文件的命令.
* `exe`: 上面命令的可执行文件路径.

### 查看审计报告

`aureport`是使用系统审计日志生成简要报告的工具.

生成审计报告,我们可以使用`aureport`工具.不带参数运行的话,可以生成审计活动的概述.

```shell
sudo aureport

Summary Report
======================
Range of time in logs: 08/17/17 00:04:12.819 - 11/27/17 19:25:23.698
Selected time for report: 08/17/17 00:04:12 - 11/27/17 19:25:23.698
Number of changes in configuration: 4682
Number of changes to accounts, groups, or roles: 59
Number of logins: 89
Number of failed logins: 19
Number of authentications: 228
Number of failed authentications: 142
Number of users: 6
Number of terminals: 15
Number of host names: 7
Number of executables: 26
Number of commands: 22
Number of files: 6
Number of AVC's: 6
Number of MAC events: 24
Number of failed syscalls: 2
Number of anomaly events: 171
Number of responses to anomaly events: 0
Number of crypto events: 1307
Number of integrity events: 0
Number of virt events: 35
Number of keys: 0
Number of process IDs: 2983
Number of events: 22263
```

#### 选项

* `-m`: 查看账户修改相关事件.
* `-au`: 授权报告.

### auditd 配置文件

可以把我们的规则添加到配置文件`/etc/audit/audit.rules`中.

```
-w /etc/passwd -p rwxa
-w /production/
```

重启`auditd`守护程序.

```
systemctl restart auditd.service
```

