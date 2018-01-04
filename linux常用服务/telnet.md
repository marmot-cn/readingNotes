# telnet

---

## 概述

`telnet`命令用于登录远程主机, 对远程主机进行管理. `telnet`因为采用**明文**传送报文, 安全性不好, 很多`Linux`服务器都不开放`telnet`服务, 而改用更安全的`ssh`方式了.

我这里只是测试, 现在服务器配置都是采用`ssh`的方式了.

## 语法

`telnet`(选项)(参数)

## 选项

* `-8`：允许使用8位字符资料, 包括输入与输出；
* `-a`：尝试自动登入远端系统；
* `-b<主机别名>`：使用别名指定远端主机名称；
* `-c`：不读取用户专属目录里的.telnetrc文件；
* `-d`：启动排错模式；
* `-e<脱离字符>`：设置脱离字符；
* `-E`：滤除脱离字符；
* `-f`：此参数的效果和指定"-F"参数相同；
* `-F`：使用Kerberos V5认证时, 加上此参数可把本地主机的认证数据上传到远端主机；
* `-k<域名>`：使用Kerberos认证时, 加上此参数让远端主机采用指定的领域名, 而非该主机的域名；
* `-K`：不自动登入远端主机；
* `-l<用户名称>`：指定要登入远端主机的用户名称；
* `-L`：允许输出8位字符资料；
* `-n<记录文件>`：指定文件记录相关信息；
* `-r`：使用类似rlogin指令的用户界面；
* `-S<服务类型>`：设置telnet连线所需的ip TOS信息；
* `-x`：假设主机有支持数据加密的功能,就使用它；
* `-X<认证形态>`：关闭指定的认证形态.

## 参数

* 远程主机: 指定要登录进行管理的远程主机.
* 端口: 指定`TELNET`协议使用的端口号.

## 示例

### 安装

```
服务端
yum -y install telnet-server.x86_64

客户端
yum -y install telnet.x86_64

yum -y install xinetd.x86_64

设置开机启动:
systemctl enable xinetd.service
systemctl enable telnet.socket

开启service:
systemctl start telnet.socket
systemctl start xinetd

可见23端口监听
[root@demo ansible]# netstat -tnlp
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
tcp        0      0 0.0.0.0:17456           0.0.0.0:*               LISTEN      1221/sshd
tcp        0      0 127.0.0.1:25            0.0.0.0:*               LISTEN      1985/master
tcp6       0      0 :::17456                :::*                    LISTEN      1221/sshd
tcp6       0      0 :::1521                 :::*                    LISTEN      12746/docker-proxy
tcp6       0      0 :::23                   :::*                    LISTEN      1/systemd
tcp6       0      0 ::1:25                  :::*                    LISTEN      1985/master
```

### 配置文件

`/etc/securetty`文件允许你规定"root"用户可以从哪个`TTY`设备登录.


### 客户端登录

```
telnet 192.168.0.239
Trying 192.168.0.239...
Connected to 192.168.0.239.
Escape character is '^]'.
Password:
Login incorrect

demo login: ansible
Password:
Last login: Thu Jan  4 12:22:05 from 192.168.0.201
[ansible@demo ~]$ ls
[ansible@demo ~]$
```

### 允许`root`登录的配置方案

修改`/etc/securetty`, 因为该文件是固定"root"用户可以从哪个`TTY`设备登录.

我们加入

```
pts/0  
pts/1 
```

因为`pts/0`是伪终端, 被我当前的用户登录占用.

```
[root@demo ansible]# who
ansible  pts/0        Jan  4 12:22 (192.168.0.201)
```

添加后我们在用`root`登录.

```
telnet 192.168.0.239
Trying 192.168.0.239...
Connected to 192.168.0.239.
Escape character is '^]'.
Password:
Login incorrect

demo login: root
Password:
Last failed login: Thu Jan  4 14:03:34 CST 2018 from ::ffff:192.168.0.201 on pts/1
There were 5 failed login attempts since the last successful login.
Last login: Thu Jan  4 13:59:06 on pts/1
[root@demo ~]# exit
logout
Connection closed by foreign host.
```

可见登录到`pts/1`上了. 但是为了方便最好不要让`root`直接登录.