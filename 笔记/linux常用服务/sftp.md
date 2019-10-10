# sftp 配置过程

## 配置

### 创建用户

* 修改家目录
* 设定不能登录

```
sudo usermod -s /usr/sbin/nologin dev

# 这里的ftp目录是我们设定好的
sudo usermod -d /ftp/dev dev
```

### 修改对应目录权限

设定家目录属主属组都为`root`

```
sudo chown root:root /ftp/dev

[ansible@ftp-1 ~]$ ll -s /ftp/
total 0
0 drwxr-xr-x 7 root root 86 Aug 30 15:04 dev
```

设定家目录里面的文件夹权限和用户匹配

```
[ansible@ftp-1 ~]$ sudo chown -R dev:dev /ftp/dev/
[ansible@ftp-1 ~]$ ll -s /ftp/dev/
total 0
0 drwxr-xr-x 3 dev dev 20 Aug 30 15:04 homestay
0 drwxr-xr-x 6 dev dev 70 Aug 29 18:09 huizhong
0 drwxr-xr-x 6 dev dev 70 Aug 29 18:09 pingxiang
0 drwxr-xr-x 4 dev dev 38 Aug 29 18:09 qixinyun
0 drwxr-xr-x 6 dev dev 70 Aug 29 18:09 renhang
```

### 修改`sshd`配置文件

```
[ansible@ftp-1 ~]$ sudo cat /etc/ssh/sshd_config
...
Subsystem       sftp    internal-sftp
Match Group sftp
X11Forwarding no
AllowTcpForwarding no
ChrootDirectory /ftp/%u
ForceCommand internal-sftp
...
```

注意上述配置文件内容需要放在该段配置下面

```
...
UseDNS no
AddressFamily inet
PermitRootLogin no
SyslogFacility AUTHPRIV
PasswordAuthentication yes
AllowGroups sshallow
```

## 参数说明

### `Match`

```
Match
             引入一个条件块. 块的结尾标志是另一个 Match 指令或者文件结尾. 
             如果 Match 行上指定的条件都满足, 那么随后的指令将覆盖全局配置中的指令. 
             Match 的值是一个或多个"条件-模式"对. 可用的"条件"是：User, Group, Host, Address . 
             只有下列指令可以在 Match 块中使用：AllowTcpForwarding, Banner,
             ForceCommand, GatewayPorts, GSSApiAuthentication,
             KbdInteractiveAuthentication, KerberosAuthentication,
             PasswordAuthentication, PermitOpen, PermitRootLogin,
             RhostsRSAAuthentication, RSAAuthentication, X11DisplayOffset,
             X11Forwarding, X11UseLocalHost
```

`Match Group sftp`这里意思是匹配用户组`sftp`

### `X11Forwarding`

是否允许进行 X11 转发. 默认值是"no", 设为"yes"表示允许. 如果允许X11转发并且sshd(8)代理的显示区被配置为在含有通配符的地址(X11UseLocalhost)上监听.  那么将可能有额外的信息被泄漏. 由于使用X11转发的可能带来的风险, 此指令默认值为"no". 
需要注意的是, 禁止X11转发并不能禁止用户转发X11通信, 因为用户可以安装他们自己的转发器. 如果启用了 UseLogin , 那么X11转发将被自动禁止. 

`X11Forwarding no`禁止转发

### `AllowTcpForwarding`

```
是否允许TCP转发, 默认值为"yes". 禁止TCP转发并不能增强安全性, 除非禁止了用户对shell的访问, 因为用户可以安装他们自己的转发器. 
```

`AllowTcpForwarding no` 禁止`TCP`转发

**关闭此功能**会禁止端口转发, 即不能使用:

```
ssh -D [bind_address:]port
```

### `ChrootDirectory`

设定属于用户组`sftp`的用户访问的根文件夹.

`ChrootDirectory /ftp/%u`, 这里设定为`/ftp/下的用户名`文件夹

### `ForceCommand`

强制执行这里指定的命令而忽略客户端提供的任何命令. 这个命令将使用用户的登录shell执行(shell -c). 这可以应用于 shell 、命令、子系统的完成, 通常用于 Match 块中. 
这个命令最初是在客户端通过 SSH_ORIGINAL_COMMAND 环境变量来支持的. 