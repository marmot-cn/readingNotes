#  08 | 白话容器基础（四）：重新认识Docker容器

## 笔记

### Linux 容器

* Docker on Mac
* Windows Docker 

都是基于虚拟化技术实现的.

### Dockerfile 的设计思想

使用一些**标准**的原语(大写高亮词语), 描述我们索要构建的`Docker`镜像. 这些原语, 都是**按顺序处理**的.

* `FROM`: 使用xx作为基础镜像
* `EXPOSE xx`: 允许外界访问容器`xx`端口
* `ENV name value`: 设置环境变量`name`值为`value`
* `RUN` 就是在容器里执行`shell`命令的意思.
* `WORKDIR`, 在这一句之后, `Dockerfile`后面的操作都以这一句指定的目录作为当前目录.
* `CMD` 意思是`Dockerfile`指定`xxx`为这个容器的进程(即 启动命令).
* `ENTRYPOINT`
	* `Docker`提供隐含的`ENTRYPOINT`即`/bin/sh -c xxx(这里的xxx就是上面的CMD)`

**Dockerfile 中的每个原语执行后, 都会生成一个对应的镜像层**. 即使原语本身并没有明显地修改文件的操作(如, ENV 原语), 它对应的层也会存在. 只不过在外界来看, **这个层是空的**.

### 不映射端口访问容器方式

需要使用`docker inspect`命令查看容器的`IP`地址, 然后访问.

```
[root@borad ~]# docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED              STATUS              PORTS               NAMES
57b507814f55        nginx               "nginx -g 'daemon ..."   About a minute ago   Up About a minute   80/tcp              kickass_saha

[root@borad ~]# docker inspect kickass_sahadock
...
            "Networks": {
                "bridge": {
                    "IPAMConfig": null,
                    "Links": null,
                    "Aliases": null,
                    "NetworkID": "58e192173b5b9b36be9f859f87ebdd4884c4960dd4cb3db464a98800bd860298",
                    "EndpointID": "065ef02a67c15f0b20744f7740f7988e869ac2063aa0743a371f7b9883f7134d",
                    "Gateway": "172.17.0.1",
                    "IPAddress": "172.17.0.2",
                    "IPPrefixLen": 16,
                    "IPv6Gateway": "",
                    "GlobalIPv6Address": "",
                    "GlobalIPv6PrefixLen": 0,
                    "MacAddress": "02:42:ac:11:00:02"
                }
            }
        }
    }
]
...

[root@borad ~]# curl 172.17.0.2
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

#### docker commit

把一个**正在运行**的容器, 直接提交为一个镜像.

实际上就是把容器运行起来后,把最上层的"**可读写层"(COW), 加上原先容器镜像的只读层**, 打包组成了一个新的镜像. 下面这些只读层在宿主机上是共享的, 不会占用额外的空间.

`Init`层的存在，就是为了避免你执行 docker cmmit 时, 把 Docker 自己对 `/etc/hosts` 等文件做的修改, 也一起提交掉.

#### 查看当前正在运行的`Docker`容器的进程号

这个的容器 id 使用的是示例的容器的id.

```
$ docker inspect --format '{{ .State.Pid }}'  4ddf4638572d
25686
```

(18907)是我做例子的时候的真实`PID`, 查看宿主机的`proc`文件, 查该进程所有的`Namespace`对应的文件.

```
[root@borad ~]# ll -s /proc/18907/ns
total 0
0 lrwxrwxrwx 1 root root 0 Oct 24 15:51 ipc -> ipc:[4026532236]
0 lrwxrwxrwx 1 root root 0 Oct 24 15:47 mnt -> mnt:[4026532234]
0 lrwxrwxrwx 1 root root 0 Oct 24 15:47 net -> net:[4026532239]
0 lrwxrwxrwx 1 root root 0 Oct 24 15:51 pid -> pid:[4026532237]
0 lrwxrwxrwx 1 root root 0 Oct 24 15:51 user -> user:[4026531837]
0 lrwxrwxrwx 1 root root 0 Oct 24 15:51 uts -> uts:[4026532235]
```

**一个进程的每种 Linux Namespace, 都在它对应的`/proc/[进程号]/ns` 下有一个对应的虚拟文件, 并且链接到一个真实的 Namespace 文件上**.

#### docker exec 的实现原理

**一个进程, 可以选择加入到某个进程已有的`Namespace`当中, 从而达到"进入"这个进程所在容器的目的, 这正是 `docker exec` 的实现原理**.

依赖一个`setns()`的 Linux 系统调用.

```
#define _GNU_SOURCE
#include <fcntl.h>
#include <sched.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>

#define errExit(msg) do { perror(msg); exit(EXIT_FAILURE);} while (0)

int main(int argc, char *argv[]) {
    int fd;
    
    fd = open(argv[1], O_RDONLY);
    if (setns(fd, 0) == -1) {
        errExit("setns");
    }
    execvp(argv[2], &argv[2]); 
    errExit("execvp");
}
```

通过`open()`系统调用打开了指定的`Namespace`文件, 并且把这个文件的描述符`fd`交给`setns()`使用. 在`setns()`执行后, 当前进程就加入了这个文件对应的 Linux Namespace 当中了.

示例: 金融到如期进程(PID=25686)的 Network Namespace 中:

```
$ gcc -o set_ns set_ns.c 
$ ./set_ns /proc/25686/ns/net /bin/bash 
$ ifconfig
eth0      Link encap:Ethernet  HWaddr 02:42:ac:11:00:02  
          inet addr:172.17.0.2  Bcast:0.0.0.0  Mask:255.255.0.0
          inet6 addr: fe80::42:acff:fe11:2/64 Scope:Link
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:12 errors:0 dropped:0 overruns:0 frame:0
          TX packets:10 errors:0 dropped:0 overruns:0 carrier:0
	   collisions:0 txqueuelen:0 
          RX bytes:976 (976.0 B)  TX bytes:796 (796.0 B)

lo        Link encap:Local Loopback  
          inet addr:127.0.0.1  Mask:255.0.0.0
          inet6 addr: ::1/128 Scope:Host
          UP LOOPBACK RUNNING  MTU:65536  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
	  collisions:0 txqueuelen:1000 
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)
```

在宿主机, 用`ps`指令查看`set_ns`程序执行的`/bin/bash`进程, 真实的`PID`是`28499`

```
# 在宿主机上
ps aux | grep /bin/bash
root     28499  0.0  0.0 19944  3612 pts/0    S    14:15   0:00 /bin/bash
```

查看 PID=28499 的进程的 Namespace, 就会发现

```
$ ls -l /proc/28499/ns/net
lrwxrwxrwx 1 root root 0 Aug 13 14:18 /proc/28499/ns/net -> net:[4026532281]

$ ls -l  /proc/25686/ns/net
lrwxrwxrwx 1 root root 0 Aug 13 14:05 /proc/25686/ns/net -> net:[4026532281]
```

在`/proc/[PID]/ns/net`目录下, 这个 PID=28499 进程, 与我们前面的 DOcker 容器进程 (PID=25686)指向的 Network Namespace 文件完全一样. **这说明这两个进程, 共享了这个名叫`net:[4026532281]`的Network Namespace**.

#### docker run --net

启动一个容器**并"加入"**到另一个容器的 Network Namespace 里.

指定`--net=host`, 就意味着这个容器**不会为进程启用 Network Namesapce**. 这就意味着,这个容器拆除了 Network Namespace 的"隔离墙". 它回合宿主机上的其他普通进程一样, 直接共享宿主机的网络栈. **容器可以直接操作和使用宿主机网络**.

### Volume 数据卷

#### 问题

* 容器里进程新建的文件, 怎么才能让宿主机获取到?
* 宿主机上的文件和目录, 怎么才能让容器里的进程访问到?

#### Volume 机制

允许你将宿主机指定的目录或者文件, 挂在到容器里面进行读取和修改.

```
$ docker run -v /test ...
$ docker run -v /home:/test ...
```

第一种, 名优显示声明宿主机目录, 俺么 Docker 就会默认在宿主机上创建一个临时目录 `/var/lib/docker/volumes/[VOLUME_ID]/_data`,然后把它挂在到容器的`/test`目录上.

第二种, Docker 直接把宿主机的`/home`目录挂载到容器的`/test`目录上.

##### 挂载原理

只需要**在`rootfs`准备好之后, 在执行`chroot`之前**, 把`Volume`指定的宿主机目录(比如`/home`目录). 挂载到指定的容器目录(比如`/test`目录)在宿主机上对应的目录(即`/var/lib/docker/aufs/mnt/[可读写层ID]/test`)上, 这个 Volume 的挂载工作就完成了.

由于执行这个挂载操作时, "容器进程"已经创建了, 也就意味着此时`Mount Namespace`已经开启了. 所以**这个挂载事件只在这个容器里可见**. 在宿主机上, 是看不见容器内部的这个挂载点的. **保证了容器的隔离性不会被 Volume 打破**.

#### 容器初始化进程 docker init

`dockerinit`会负责完成:

* 根目录的准备
* 挂载设备
* 挂载目录
* 配置`hostname`
* ... 一些需要在容器内进行初始化操作.

最后通过`execv()`系统调用, 让**应用进程取代自己**, 称为容器里的 PID=1 的进程.

#### Linux 的绑定挂载(bind mount) 机制

主要作用就是, **允许你将一个目录或者文件,而不是整个设备, 挂载到一个指定的目录上**.

这时你在该挂载点上进行的任何操作, 只是发生在被挂载的目录或者文件上, 而原挂载点的内容则会被隐藏起来且不受影响.

###### 原理

**绑定挂载**实际上是一个`inode`替换的过程. 在 Linux 操作系统中:

* `inode`可以理解为存放**文件内容的"对象"**(有关该文件的组织和管理的信息主要存放inode里面，它记录着文件在存储介质上的位置与分布, inode代表的是物理意义上的文件，通过inode可以得到一个数组，这个数组**记录了文件内容的位置**, 如该文件位于硬盘的第3，8，10块，那么这个数组的内容就是3,8,10).
* `dentry`, 也叫目录项, 就是访问这个`inode`所使用的**"指针"**(dentry记录着文件名，上级目录等信息，正是它形成了我们所看到的树状结构).

![](./img/08_01.png)

`mount --bind /home /test`, 会将`/home`挂在到`/test`上. 其实相当于将`/test`的`dentry`, 重定向到了`/home`的`inode`. 这样当我们修改`/test`目录时, 实际修改的是`/home`目录的`inode`.

一旦执行`umount`命令, `/test`目录原先的内容就会恢复: **因为修改真正发生在的, 是`/home`目录里**.

##### 挂载内容会不会被 `docker commit`提交

由于`Mount Namespace`的隔离作用, **宿主机并不知道这个绑定挂载的存在**. 所以, 在宿主机看来,容器中可读写层的`/test`目录(/var/lib/docker/aufs/mnt/[可读写层ID]/test), **始终是空的**.

由于 Docker 一开始还是要创建`/test`这个目录作为挂载点, 所以执行了 `docker commit`之后, 会发现新产生的镜像里, 会多出一个空的`/test`目录. 因为 **新建目录**, **不是挂载操作**, **Mount Namespace**对它可起不到"障眼法"的作用.

##### 示例

启动一个 helloworld 容器, 声明一个 Volume, 挂载在容器里的`/test`目录上.

```
$ docker run -d -v /test helloworld
cf53b766fa6f
```

查看这个`Volume`的ID.

```
$ docker volume ls
DRIVER              VOLUME NAME
local               cb1c2f7221fa9b0971cc35f68aa1034824755ac44a034c0c0a1dd318838d3a6d
```

找到在 Docker 工作目录下的 volumes 路径.

```
$ ls /var/lib/docker/volumes/cb1c2f7221fa/_data/
```

在容器的 Volume 里, 添加一个文件 

```
$ docker exec -it cf53b766fa6f /bin/sh
cd test/
touch text.txt
```

宿主机, 发现文件已经出现在宿主机上对应的临时目录里

```
$ ls /var/lib/docker/volumes/cb1c2f7221fa/_data/
text.txt
```

在宿主机上查看该容器的可读写层, 虽然可以看到这个`/test`目录, 但其内容是空的.

```
$ ls /var/lib/docker/aufs/mnt/6780d0778b8a/test
```

可确认:

* 容器 Volume 里的信息, 不会被 docker commit 提交掉.
* 这个挂载点目录 /test 本身, 会出现在新的镜像当中

### 总结

![](./img/08_02.png)

示例中的容器进程`python app.py`运行在由Linux Namesapce 和 Cgroups 构成的隔离环境里. 而它运行所需要的各种文件, 比如 python, app.py, 以及整个操作系统文件, 则由多个联合挂载在一起的`rootfs`层提供.

容器声明的`Volume`挂载点, 也在**可读写层**.

## 扩展

### `/bin/sh -c xxx`

```
-c string  If the -c option is present, then commands are read from  string.   If  there  are  arguments  after  the string, they are assigned to the positional parameters, starting with $0.
```