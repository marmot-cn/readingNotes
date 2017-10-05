# Docker以及Docker网络 

---

我在`unutun`系统下进行的测试

## 解压image, 分析构成

把镜像保存为一个`tar`包.

```shell
root@iZ94xwu3is8Z: docker image save busybox:glibc -o busybox.tar
root@iZ94xwu3is8Z: mkdir busybox
root@iZ94xwu3is8Z: tar -xvf busybox.tar -C busybox/
44425705f75f9f3139c5e7dc94b94f7fc80787fba82eca8cc8fd5aee600a1f1f/
44425705f75f9f3139c5e7dc94b94f7fc80787fba82eca8cc8fd5aee600a1f1f/VERSION
44425705f75f9f3139c5e7dc94b94f7fc80787fba82eca8cc8fd5aee600a1f1f/json
44425705f75f9f3139c5e7dc94b94f7fc80787fba82eca8cc8fd5aee600a1f1f/layer.tar
8e24a6e63b274e4e7518dbbd5f8e8cdd5ebd3f10e960048870c2375de6a1127d.json
manifest.json
tar: manifest.json: implausibly old time stamp 1970-01-01 08:00:00
repositories
tar: repositories: implausibly old time stamp 1970-01-01 08:00:00
root@iZ94xwu3is8Z: cd busybox
[root@localhost busybox]# ls
44425705f75f9f3139c5e7dc94b94f7fc80787fba82eca8cc8fd5aee600a1f1f  8e24a6e63b274e4e7518dbbd5f8e8cdd5ebd3f10e960048870c2375de6a1127d.json  manifest.json  repositories
```

`manifest.json` 配置文件

```
root@iZ94xwu3is8Z:~# cd busybox/
root@iZ94xwu3is8Z:~/busybox# cat manifest.json
[{"Config":"8e24a6e63b274e4e7518dbbd5f8e8cdd5ebd3f10e960048870c2375de6a1127d.json","RepoTags":["busybox:glibc"],"Layers":["44425705f75f9f3139c5e7dc94b94f7fc80787fba82eca8cc8fd5aee600a1f1f/layer.tar"]}]
```

进入镜像的层

```
root@iZ94xwu3is8Z:~/busybox# cd 44425705f75f9f3139c5e7dc94b94f7fc80787fba82eca8cc8fd5aee600a1f1f/
root@iZ94xwu3is8Z:~/busybox/44425705f75f9f3139c5e7dc94b94f7fc80787fba82eca8cc8fd5aee600a1f1f# ls
VERSION  json  layer.tar

解开 laver.tar
root@iZ94xwu3is8Z:~/busybox/44425705f75f9f3139c5e7dc94b94f7fc80787fba82eca8cc8fd5aee600a1f1f# tar -xvf layer.tar
bin/
bin/[
bin/[[
bin/acpid
bin/add-shell
...
...

解压后可见是类似于一个linux的根目录,少了一些不用的精简掉. docker 本就是把一个linux根目录整体打包.
root@iZ94xwu3is8Z:~/busybox/44425705f75f9f3139c5e7dc94b94f7fc80787fba82eca8cc8fd5aee600a1f1f# ls -l
total 4556
-rw-r--r-- 1 root   root          3 Sep 13 18:14 VERSION
drwxr-xr-x 2 root   root      12288 Aug 23 07:28 bin
drwxr-xr-x 2 sys    sys        4096 Aug 23 07:28 dev
drwxr-xr-x 2 root   root       4096 Aug 23 07:28 etc
drwxr-xr-x 2 nobody nogroup    4096 Aug 23 07:28 home
-rw-r--r-- 1 root   root       1174 Sep 13 18:14 json
-rw-r--r-- 1 root   root    4610048 Sep 13 18:14 layer.tar
drwxr-xr-x 2 root   root       4096 Aug 23 07:28 lib
lrwxrwxrwx 1 root   root          3 Aug 23 07:28 lib64 -> lib
drwxr-xr-x 2 root   root       4096 Aug 23 07:28 root
drwxrwxrwt 2 root   root       4096 Aug 23 07:28 tmp
drwxr-xr-x 3 root   root       4096 Aug 23 07:28 usr
drwxr-xr-x 4 root   root       4096 Aug 23 07:28 var

移动出来留作它用
root@iZ94xwu3is8Z: mv 44425705f75f9f3139c5e7dc94b94f7fc80787fba82eca8cc8fd5aee600a1f1f /var/lib/xocker/image/busybox
```

## 分层文件系统

### 查看`chroot`的效果

```shell
备份
root@iZ94xwu3is8Z: cp -r /var/lib/xocker/image/busybox/ /var/lib/xocker/image/busybox.back

锁定根目录, 操作最多破坏该路径下的东西, 可以作为隔离的手段
root@iZ94xwu3is8Z: chroot /var/lib/xocker/image/busybox/ /bin/sh
/#
```

### `aufs`工作原理

```shell
root@iZ94xwu3is8Z: mkdir -p /var/lib/xocker/mnt/1
root@iZ94xwu3is8Z: mkdir -p /var/lib/xocker/mnt/1-data
root@iZ94xwu3is8Z: mkdir -p /var/lib/xocker/mnt/1-init
```

* `1-init`: 容器启动时, 初始化的配置文件.
* `1-data`: 容器启动后生成的数据.
* `1`: 是`container id`,`docker`会分配一个`uuid`.

`/var/lib/xocker/mnt/1`也是容器执行的根目录.

```shell
root@iZ94xwu3is8Z: mkdir -p /var/lib/xocker/mnt/1-init/etc/ && mkdir -p /var/lib/xocker/mnt/1-init/proc && echo "hello" > /var/lib/xocker/mnt/1-init/etc/myinit && tree /var/lib/xocker/mnt/1-init
/var/lib/xocker/mnt/1-init
|-- etc
|   `-- myinit
`-- proc
root@iZ94xwu3is8Z:/var/lib/xocker/image# mount -t aufs -o dirs=/var/lib/xocker/mnt/1-data:/var/lib/xocker/mnt/1-init:/var/lib/xocker/image/busybox none /var/lib/xocker/mnt/1

root@iZ94xwu3is8Z:~# ls -l /var/lib/xocker/mnt/1
total 4560
-rw-r--r-- 1 root   root          3 Sep 13 18:14 VERSION
drwxr-xr-x 2 root   root      12288 Aug 23 07:28 bin
drwxr-xr-x 2 sys    sys        4096 Aug 23 07:28 dev
drwxr-xr-x 2 root   root       4096 Oct  5 16:02 etc
drwxr-xr-x 2 nobody nogroup    4096 Aug 23 07:28 home
-rw-r--r-- 1 root   root       1174 Sep 13 18:14 json
-rw-r--r-- 1 root   root    4610048 Sep 13 18:14 layer.tar
drwxr-xr-x 2 root   root       4096 Aug 23 07:28 lib
lrwxrwxrwx 1 root   root          3 Aug 23 07:28 lib64 -> lib
drwxr-xr-x 2 root   root       4096 Oct  5 16:02 proc
drwxr-xr-x 2 root   root       4096 Oct  5 16:01 root
drwxrwxrwt 2 root   root       4096 Aug 23 07:28 tmp
drwxr-xr-x 3 root   root       4096 Aug 23 07:28 usr
drwxr-xr-x 4 root   root       4096 Aug 23 07:28 var
root@iZ94xwu3is8Z:~# cat /var/lib/xocker/mnt/1/etc/myinit
hello

data目录是空的
root@iZ94xwu3is8Z:~# ls /var/lib/xocker/mnt/1-data/

创建一个文件
root@iZ94xwu3is8Z:~# chroot /var/lib/xocker/mnt/1 /bin/sh
/ # touch tmp/test.data
/ # ls
VERSION    bin        dev        etc        home       json       layer.tar  lib        lib64      proc       root       tmp        usr        var
/ # ls tmp/
test.data
/ # exit

再次查看 data 目录已经写入新的文件
root@iZ94xwu3is8Z:~# ls /var/lib/xocker/mnt/1-data/
root  tmp

root@iZ94xwu3is8Z:~# chroot /var/lib/xocker/mnt/1 /bin/sh
/ # cat /etc/myinit
hello
/ # rm /etc/myinit
/ # ls /etc/myint
ls: /etc/myint: No such file or directory
/ # exit

退出后查看原始文件还在
root@iZ94xwu3is8Z:~# cat /var/lib/xocker/mnt/1-init/etc/myinit
hello
```

`data`目录是可写的, 创建的文件是在`data`目录下创建的. 删除的文件会在`rw`层(即这里的`data`层)创建一个标识.

```shell
root@iZ94xwu3is8Z:~# ls -a /var/lib/xocker/mnt/1-data/etc/
.  ..  .wh.myinit
```

## 容器网络

`docker0` `bridge`的网络.

```shell
root@iZ94xwu3is8Z:~# brctl addbr xocker0
root@iZ94xwu3is8Z:~# ip addr add 172.18.0.1/24 dev xocker0
root@iZ94xwu3is8Z:~# ip link set dev xocker0 up
root@iZ94xwu3is8Z:~# ifconfig
...
xocker0   Link encap:Ethernet  HWaddr 32:a2:f8:4b:03:7e
          inet addr:172.18.0.1  Bcast:0.0.0.0  Mask:255.255.255.0
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)

```

