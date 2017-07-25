# 服务器添加ftp功能

---

因为需求不能`ssh`登录, 需要`ftp`登录到服务器上.

我们需要给服务器添加`ftp`服务, 并添加对应账户.

**下载配置`ftp`服务器**

```shell
下载 vsftpd 软件
sudo yum install vsftpd
...
Y
...

sudo vi /etc/vsftpd/vsftpd.conf

设置指定的用户执行chroot
chroot_local_user=YES
chroot_list_enable=YES
chroot_list_file=/etc/vsftpd/chroot_list

添加读取用户配置目录
原来没有,需要自己添加该配置项
user_config_dir=/etc/vsftpd/userconf

设置 discuz 不能切换根目录
sudo vi  /etc/vsftpd/chroot_list
...
discuz
...

mkdir -p /etc/vsftpd/userconf

sudo vi /etc/vsftpd/userconf/discuz
...
local_root=/data/html/DiscuzX/upload/
...

service vsftpd restart
```

**添加`ftp`使用账户**

```shell
添加用户 xxx 不能让其登录

或可以用此命令设置用户不能登录
usermod -s /usr/sbin/nologin xxx
```
