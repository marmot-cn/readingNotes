# last

---

**成功**登录, 和`lastb`区别是,`lastb`是记录失败的登录.

用来列出目前与过去登录系统的用户相关信息.

它会读取位于`/var/log`目录下名称为`wtmp`的文件(二进制文件), 并把该给文件的内容记录的登录系统的用户名单全部显示出来.

默认是显示`wtmp`的记录,`btmp`能显示的更详细, 可以显示远程登录, 例如`ssh`登录.

## 语法

`last (选项) (参数)`

## 选项

* `-a`: 把从何处登入系统的主机名称或ip地址，显示在最后一行.
* `-d`：将IP地址转换成主机名称.
* `-f <记录文件>`：指定记录文件. 默认是显示`/var/log`目录下的`wtmp`文件的记录, 但`/var/log`目录下的`btmp`能显示的内容更丰富, 可以显示远程登录, 例如ssh登录 ,包括失败的登录请求.
* `-n <显示列数>或-<显示列数>`：设置列出名单的显示列数.
* `-R`：不显示登入系统的主机名称或IP地址
* `-x`：显示系统关机，重新开机，以及执行等级的改变等信息.
* ` -F`: 显示完整登入登出时间日期

## 命令输出字段介绍

* 第一列: 用户名
* 第二列: 终端位置. 
	* `pts/0`(伪终端).意味着从诸如`SSH`或`telnet`的远程连接的用户.
	* `tty(teletypewriter)`.意味着直接连接到计算机或者本地连接的用户.
	* 如果是重启或者启动操作, 这里会显示`system boot`.
* 第三列: 登录ip或者内核. 如果你看见`:0.0`或者什么都没有,这意味着用户通过本地终端连接.系统操作, 开机,关机,重启.内核版本会显示在状态中.
* 第四列信息: 开始时间,其中的日期格式为date +"%a %b %d".
* 第五列信息：结束时间(still login in 还未退出 `down` 直到正常关机 `crash`直到强制关机)
* 第六列信息: 持续时间

## 示例

```shell
[root@localhost ansible]# last
ansible  pts/0        192.168.0.103    Sun Aug 20 16:13   still logged in
root     tty1                          Sun Aug 20 16:13   still logged in
reboot   system boot  3.10.0-514.el7.x Sun Aug 20 16:12 - 16:59  (00:46)
ansible  pts/1        192.168.0.139    Fri Aug 18 01:17 - crash (2+14:55)
ansible  pts/0        192.168.0.139    Fri Aug 18 00:45 - 01:18  (00:32)
chloropl pts/0        192.168.0.139    Thu Aug 17 20:05 - 20:28  (00:23)
chloropl pts/0        192.168.0.139    Thu Aug 17 20:02 - 20:02  (00:00)
chloropl pts/0        192.168.0.139    Thu Aug 17 19:50 - 19:51  (00:00)
chloropl pts/1        192.168.0.139    Thu Aug 17 18:26 - 19:50  (01:23)
chloropl pts/1        192.168.0.139    Thu Aug 17 18:21 - 18:21  (00:00)
root     pts/0        192.168.0.139    Thu Aug 17 17:38 - 19:06  (01:28)
root     pts/0        192.168.0.139    Thu Aug 17 17:02 - 17:08  (00:05)
root     pts/1        192.168.0.139    Thu Aug 17 16:39 - 16:47  (00:08)
root     pts/1        192.168.0.139    Thu Aug 17 16:35 - 16:36  (00:00)
root     pts/1        192.168.0.139    Thu Aug 17 16:28 - 16:35  (00:06)
root     pts/1        192.168.0.139    Thu Aug 17 15:25 - 15:30  (00:05)
root     pts/0        192.168.0.139    Thu Aug 17 14:56 - 16:49  (01:52)
root     pts/0        192.168.0.139    Thu Aug 17 14:50 - 14:56  (00:05)
root     tty1                          Thu Aug 17 14:49 - crash (3+01:23)
reboot   system boot  3.10.0-514.el7.x Thu Aug 17 00:19 - 16:59 (3+16:39)
root     pts/0        192.168.0.139    Thu Aug 17 00:17 - down   (00:02)
root     pts/0        192.168.0.139    Thu Aug 17 00:15 - 00:17  (00:02)
root     tty1                          Thu Aug 17 00:11 - 00:19  (00:07)
root     pts/0        192.168.0.103    Thu Aug 17 00:06 - 00:11  (00:04)
root     tty1                          Thu Aug 17 00:04 - 00:11  (00:07)
reboot   system boot  3.10.0-514.el7.x Thu Aug 17 00:04 - 00:19  (00:15)

[root@localhost ansible]# last -f /var/log/btmp
ansible  ssh:notty    192.168.0.139    Fri Aug 18 00:42    gone - no logout
ansible  ssh:notty    192.168.0.139    Fri Aug 18 00:42 - 00:42  (00:00)
ansible  ssh:notty    192.168.0.139    Fri Aug 18 00:42 - 00:42  (00:00)
ansible  ssh:notty    192.168.0.139    Fri Aug 18 00:42 - 00:42  (00:00)
root     ssh:notty    192.168.0.139    Fri Aug 18 00:41 - 00:42  (00:00)
root     pts/0                         Thu Aug 17 20:06    gone - no logout
root     ssh:notty    192.168.0.139    Thu Aug 17 20:05 - 00:41  (04:36)
root     ssh:notty    192.168.0.139    Thu Aug 17 20:04 - 20:05  (00:00)
root     ssh:notty    192.168.0.139    Thu Aug 17 18:21 - 20:04  (01:43)
root     ssh:notty    192.168.0.139    Thu Aug 17 15:25 - 18:21  (02:56)

btmp begins Thu Aug 17 15:25:29 2017
```

### 打印特定的用户名

`last 用户名`

```shell
[root@localhost ansible]# last ansible
ansible  pts/0        192.168.0.103    Sun Aug 20 16:13   still logged in
ansible  pts/1        192.168.0.139    Fri Aug 18 01:17 - crash (2+14:55)
ansible  pts/0        192.168.0.139    Fri Aug 18 00:45 - 01:18  (00:32)

wtmp begins Thu Aug 17 00:04:11 2017
```

### 打印特定`pts`

```shell
[root@localhost ansible]# last tty1
root     tty1                          Sun Aug 20 16:13   still logged in
root     tty1                          Thu Aug 17 14:49 - crash (3+01:23)
root     tty1                          Thu Aug 17 00:11 - 00:19  (00:07)
root     tty1                          Thu Aug 17 00:04 - 00:11  (00:07)
```
