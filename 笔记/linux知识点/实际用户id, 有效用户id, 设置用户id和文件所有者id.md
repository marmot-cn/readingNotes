# 实际用户id, 有效用户id, 设置用户id和文件所有者id

---

### 实际用户ID

这个ID就是我们登录`unix`系统时的身份ID.

一个进程的`real user ID`是指运行此进程的用户ID.

### 有效用户ID

定义了操作者的权限. 有效用户ID是进程的属性, 决定了该进程对文件的访问权限.

一个进程的`effective user ID`是指**此进程目前实际有效用户ID**,`effective user ID`主要用来校验权限使用. 如打开文件, 创建文件, 修改文件, kill别的进程, 等等.

如果一个进程是以`root`身份来运行的, 那么上面这两个`ID`可以用`setuid/seteuid`随便修改.

但是如果一个进程是以普通用户身份来运行的, 那么上面这两个`ID`一般来说是相同的, 并且也不能随便修改.

只有一种情况例外: 此进程的可执行文件的权限标记中, 设置了**设置用户ID**位. 那么它的有效用户ID就是可执行文件的拥有者.

下文的`s`就是设置用户id位.

### 设置用户id

作用是我们如何去修改有效用户ID.

```shell
chmod +s /path/to/file
```

使用这个命令后, 在执行这个文件, 那么生成的进程`effective user ID`就变成了这个可执行文件的`ownser user ID`(属主用户ID), 而`real user ID`仍然是启动这个程序时所用的用户`ID`.

```shell
无执行权限的a, 赋予设置用户id位后权限位变为S
[root@localhost ~]# touch a
[root@localhost ~]# ll -s a
0 -rw-r--r--. 1 root root 0 Oct  4 15:45 a
[root@localhost ~]# chmod +s a
[root@localhost ~]# ll -s a
0 -rwSr-Sr--. 1 root root 0 Oct  4 15:45 a

有执行权限的b, 赋予设置用户id位权限变为S
[root@localhost ~]# touch b
[root@localhost ~]# chmod +x b
[root@localhost ~]# ll -s b
0 -rwxr-xr-x. 1 root root 0 Oct  4 15:46 b
[root@localhost ~]# chmod +s b
[root@localhost ~]# ll -s b
0 -rwsr-sr-x. 1 root root 0 Oct  4 15:46 b
```

### 保存设置用户ID (SUID)

是有效用户ID的副本, 既然有效用户ID是副本, 它的作用肯定是为了以后**恢复有效用户ID**用的. 

## ruid, euid suid, fuid

* `ruid`: 实际用户id: `real userid`
* `euid`: 有效用户id: `effective userid`
* `suid`: 保存用户id: `saved userid`
* `fuid`: 文件系统用户id
	* 在Linux中使用, 且只用于对文件系统的访问权限控制, 在没有明确设定的情况下与EUID相同. 设立FSUID是为了允许程序(如NFS服务器)在不需获取向给定UID账户发送信号的情况下以给定UID的权限来限定自己的文件系统权限.

还有一个设置用户id位, `set user id bit`, 就是`rwx`之外的那个`s`标志位.

```shell
[root@localhost ~]# ps -ax -o ruid -o euid -o suid -o fuid -o pid -o fname
 RUID  EUID  SUID  FUID   PID COMMAND
    0     0     0     0     1 systemd
    0     0     0     0     2 kthreadd
...
```

## 示例

```shell
可见 passwd 有s位, 即设置用户id位
[test@localhost root]$ ll -s /usr/bin/passwd
28 -rwsr-xr-x. 1 root root 27832 Jun 10  2014 /usr/bin/passwd

创建用户 test 执行 passwd 命令
[root@localhost ~]# useradd test
[root@localhost ~]# su test
[test@localhost root]$ passwd
Changing password for user test.
Changing password for test.
(current) UNIX password:
...

在另外一个终端可见, 该进程是以root用户的ID来运行的
[root@localhost ~]# ps -ef | grep "passwd"
root     12803 12731  0 19:02 pts/0    00:00:00 passwd

```