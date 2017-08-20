# lastb 

---

`lastb`命令**记录失败**的登录尝试. 必须拥有`root`权限才能运行`lastb`命令. `lastb`会解析`/var/log/btmp`的信息.

```shell
[root@localhost ansible]# lastb
ansible  ssh:notty    192.168.0.139    Fri Aug 18 00:42 - 00:42  (00:00)
ansible  ssh:notty    192.168.0.139    Fri Aug 18 00:42 - 00:42  (00:00)
ansible  ssh:notty    192.168.0.139    Fri Aug 18 00:42 - 00:42  (00:00)
ansible  ssh:notty    192.168.0.139    Fri Aug 18 00:42 - 00:42  (00:00)
root     ssh:notty    192.168.0.139    Fri Aug 18 00:41 - 00:41  (00:00)
root     pts/0                         Thu Aug 17 20:06 - 20:06  (00:00)
root     ssh:notty    192.168.0.139    Thu Aug 17 20:05 - 20:05  (00:00)
root     ssh:notty    192.168.0.139    Thu Aug 17 20:04 - 20:04  (00:00)
root     ssh:notty    192.168.0.139    Thu Aug 17 18:21 - 18:21  (00:00)
root     ssh:notty    192.168.0.139    Thu Aug 17 15:25 - 15:25  (00:00)

btmp begins Thu Aug 17 15:25:29 2017
```

## 语法

`lastb [-adRx][-f <记录文件>][-n <显示列数>][帐号名称...][终端机编号...]`

## 选项

* `-a`: 把从何处登入系统的主机名称或ip地址，显示在最后一行.
* `-d`：将IP地址转换成主机名称.
* `-f <记录文件>`：指定记录文件. 默认是显示`/var/log`目录下的`wtmp`文件的记录, 但`/var/log`目录下的`btmp`能显示的内容更丰富, 可以显示远程登录, 例如ssh登录 ,包括失败的登录请求.
* `-n <显示列数>或-<显示列数>`：设置列出名单的显示列数.
* `-R`：不显示登入系统的主机名称或IP地址
* `-x`：显示系统关机，重新开机，以及执行等级的改变等信息.
* ` -F`: 显示完整登入登出时间日期

## 示例

````shell
[root@localhost ansible]# lastb
ansible  ssh:notty    192.168.0.139    Fri Aug 18 00:42 - 00:42  (00:00)
ansible  ssh:notty    192.168.0.139    Fri Aug 18 00:42 - 00:42  (00:00)
ansible  ssh:notty    192.168.0.139    Fri Aug 18 00:42 - 00:42  (00:00)
ansible  ssh:notty    192.168.0.139    Fri Aug 18 00:42 - 00:42  (00:00)
root     ssh:notty    192.168.0.139    Fri Aug 18 00:41 - 00:41  (00:00)
root     pts/0                         Thu Aug 17 20:06 - 20:06  (00:00)
root     ssh:notty    192.168.0.139    Thu Aug 17 20:05 - 20:05  (00:00)
root     ssh:notty    192.168.0.139    Thu Aug 17 20:04 - 20:04  (00:00)
root     ssh:notty    192.168.0.139    Thu Aug 17 18:21 - 18:21  (00:00)
root     ssh:notty    192.168.0.139    Thu Aug 17 15:25 - 15:25  (00:00)

btmp begins Thu Aug 17 15:25:29 2017
```
