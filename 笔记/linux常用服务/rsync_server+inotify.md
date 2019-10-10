# Rsync + Inotify

---

只有主节点可以写.

## rsync

### 安装

```
yum install rsync -y
```

### 编辑配置文件`/etc/rsyncd.conf`


```
vim /etc/rsyncd.conf 

#设置进行数据传输时所使用的账户名称或ID号, 默认使用nobody, 这里使用www-data
uid = 33 
#设置进行数据传输时所使用的组名称或GID号, 默认使用nobody, 这里使用www-data
gid = 33 

#设置user chroot为yes后，rsync会首先进行chroot设置，将根映射到path参数路径下，对客户 
#端而言，系统的根就是path参数所指定的路径。但这样做需要root权限，并且在同步符号 
#连接资料时仅会同步名称，而内容将不会同步。 
use chroot = no 

#开启Rsync数据传输日志功能 
transfer logging = yes

#是否允许客户端上传数据
read only = false

#模块，Rsync通过模块定义同步的目录, 模块以[name]的形式定义, 这与Samba定义共 
#享目录是一样的效果. 在Rsync中也可以定义多个模块 
[www-data] 
#comment定义注释说明字串 
www-data = Web content 

#忽略一些IO错误 
ignore errors 

#同步目录的真实路径通过path指定 
path = /data/www-data

#exclude可以指定例外的目录，即将common目录下的某个目录设置为不同步数据 
#exclude = test/ 

#设置允许连接服务器的账户，账户可以是系统中不存在的用户 
auth users = tom,jerry 

#设置允许哪些主机可以同步数据，可以是单个IP，也可以是网段，多个IP与网段之间使用空 
#格分隔 
hosts allow=192.168.0.205 192.168.0.203

#设置拒绝所有（除hosts allow定义的主机外） 
hosts deny=* 
#客户端请求显示模块列表时，本模块名称是否显示，默认为true 
list= false

auth users = rsync
secrets file = /etc/rsyncd.passwd
```

```
[root@localhost www-data]# cat /etc/rsyncd.passwd
rsync:J,wt["52>D4%Akq<
```

### 启动

```
rsync --daemon
```

### 使用

我们在`192.168.0.204`机器上安装`rsync`并启动.

在`192.168.0.203`机器上修改文件, 并进行同步.

```
rsync -avz --delete /data/www-data/ rsync@192.168.0.204::www-data
```

`192.168.0.203`上的任何文件修改都会自动同步到`204`上.

### 一些注意点

如果使用`auth users`和`secret file`, 必须提前在服务器设置用户名和密码.

```
[root@localhost www-data]# useradd rsync
[root@localhost www-data]# passwd rsync
```

传输

```
rsync -avz rsync@192.168.0.204::www-data /data/www-data
```

如果没有设置这两项, 则可以不用用户名和密码即可传输

```
rsync -avz 192.168.0.204::www-data /data/www-data
```

为了安全, 需要设置密码文件权限未`600`

## 配合`inotify`实现自动同步

### 安装`inotify`

官网`https://sourceforge.net/projects/inotify-tools/?source=typ_redirect`

```
[root@localhost ~]# tar -zxvf inotify-tools-3.13.tar.gz
[root@localhost ~]# cd inotify-tools-3.13
[root@localhost inotify-tools-3.13]# ./configure
checking for a BSD-compatible install... /usr/bin/install -c
checking whether build environment is sane... yes
checking for gawk... gawk
checking whether make sets $(MAKE)... yes
checking whether make sets $(MAKE)... (cached) yes
checking for gcc... no
checking for cc... no
checking for cl.exe... no
configure: error: no acceptable C compiler found in $PATH
See `config.log' for more details.

安装 gcc

[root@localhost inotify-tools-3.13]# yum install gcc -y
[root@localhost inotify-tools-3.13]# make
[root@localhost inotify-tools-3.13]# make install
```

### 编写`shell`脚本

```shell
#!/bin/bash
src=/data/www-data/
des=rsync@192.168.0.204::www-data
/usr/local/bin/inotifywait -mrq --timefmt '%d/%m/%y %H:%M' --format '%T %w%f' \
-e modify,delete,create,attrib ${src} \
| while read x
    do
        /usr/bin/rsync -avz --delete --progress $src $des --password-file=/root/rsyncpass &&
        echo "$x was rsynced" >> /var/log/rsync.log
    done
```

这里需要注意, 同步目录`src=/data/www-data/`最后需要有一个`/`, 否则在同步的时候会把该文件夹自动再次同步一份过去.

### 编写服务脚本

```
[www-data@localhost www-data]$ cat /etc/systemd/system/sync.service
[Unit]
Description = SyncService
After = network.target

[Service]
PIDFile = /run/syncservice/syncservice.pid
User = root
Group = root
WorkingDirectory = /opt
ExecStartPre = /bin/mkdir /run/syncservice
ExecStartPre = /bin/chown -R root:root /run/syncservice
ExecStart = /bin/bash /root/rsync.sh
ExecReload = /bin/kill -s HUP $MAINPID
ExecStop = /bin/kill -s TERM $MAINPID
ExecStopPost = /bin/rm -rf /run/syncservice
PrivateTmp = true

[Install]
WantedBy = multi-user.target

[www-data@localhost www-data]$ chmod 755 /etc/systemd/system/sync.service  
[www-data@localhost www-data]$ systemctl daemon-reload
[root@localhost system]# systemctl start sync.service
[root@localhost system]# systemctl status sync.service
● sync.service - SyncService
   Loaded: loaded (/etc/systemd/system/sync.service; disabled; vendor preset: disabled)
   Active: active (running) since Mon 2018-08-06 17:37:43 CST; 3s ago
  Process: 3177 ExecStartPre=/bin/chown -R root:root /run/syncservice (code=exited, status=0/SUCCESS)
  Process: 3175 ExecStartPre=/bin/mkdir /run/syncservice (code=exited, status=0/SUCCESS)
 Main PID: 3181 (bash)
   Memory: 384.0K
   CGroup: /system.slice/sync.service
           ├─3181 /bin/bash /root/rsync.sh
           ├─3182 /usr/local/bin/inotifywait -mrq --timefmt %d/%m/%y %H:%M --format %T %w%f -e modify,delete,create,attrib /data/www-data/
           └─3183 /bin/bash /root/rsync.sh

Aug 06 17:37:43 localhost.localdomain systemd[1]: Starting SyncService...
Aug 06 17:37:43 localhost.localdomain systemd[1]: Started SyncService.
```

需要设置开机启动

```
systemctl enable sync.service
```