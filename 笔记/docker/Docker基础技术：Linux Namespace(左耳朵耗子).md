# Docker基础技术：Linux Namespace(左耳朵耗子)

## Linux Namespace

`LinuxNamespace`是`Linux`提供的一种内核界别环境隔离的方法.

`chroot`提供了一种简单的隔离模式: `chroot`内部的文件系统无法访问外部的内容.

`Linux Namespace`在此基础上, 提供对`UTS, IPC, mount, PID, network, User`等的隔离机制.

`Linux`下的父进程的`PID`是`1`. **同chroot一样, 如果我们可以把用户的进程空间xxx到某个进程分支下, 并像chroot那样让其下面的进程看到的那个超级父进程的PID为1**, 于是就可以达到资源隔离的效果了(不同的PID namespace中的进程无法看到彼此).

### Linux Namespace 种类

分类   | 系统调用参数
------------- | -------------
Mount namespaces  | CLONE_NEWNS
UTS namespaces  | CLONE_NEWUTS
IPC namespaces  | CLONE_NEWIPC
PID namespaces  | CLONE_NEWPID
Network namespaces  | CLONE_NEWNET
User namespaces  | CLONE_NEWUSER

主要三个系统调用

* `clone()` 实现线程的系统调用, 用来创建一个进程, 并可以通过设计上述参数达到隔离.
* `unshare()` 使某进程脱离某个`namespace`
* `setns()` 把某进程加入到某个`namespace`

### clone() 系统调用

```C
int container_pid = clone(container_main, container_stack+STACK_SIZE, SIGCHLD, NULL);
```

### UTS Namespace

**CLONE_NEWUTS**

the UTS namespaces feature allows each container to have its own hostname and NIS domain name.

```C
int container_main(void * arg)
{
...
 sethostname("container",10);
...
}

int main()
{
...
int container_pid = clone(container_main, container_stack+STACK_SIZE, 
            CLONE_NEWUTS | SIGCHLD, NULL);
...
}
```

子进程的hostname变成了`container`.

```shell
hchen@ubuntu:~$ sudo ./uts
Parent - start a container!
Container - inside the container!
root@container:~# hostname
container
root@container:~# uname -n
container
```

### IPC Namespace

**CLONE_NEWIPC**

`Inter-Process Communication`, 是`Unix/Linux`下进程间通信的一种方式. `IPC`有共享内存, 信号量, 消息队列等方法.

为了隔离, 我们也需要把`IPC`个隔离开来, 这样, 只有在同一个`Namespce`下的进程才能相互通信.

`IPC`需要一个全局的`ID`, 我们的`Namespace`需要对这个`ID`隔离, 不能让别的`Namespace`的进程看到.

```C
int container_pid = clone(container_main, container_stack+STACK_SIZE, 
            CLONE_NEWUTS | CLONE_NEWIPC | SIGCHLD, NULL);
```

先创建一个IPC的Queue

```shell
hchen@ubuntu:~$ ipcmk -Q 
Message queue id: 0
 
hchen@ubuntu:~$ ipcs -q
------ Message Queues --------
key        msqid      owner      perms      used-bytes   messages    
0xd0d56eb2 0          hchen      644        0            0
```

如果我们运行没有`CLONE_NEWIPC`的程序, 我们会看到, 在子进程中还是能看到这个全启的IPC Queue.

```shell
hchen@ubuntu:~$ sudo ./uts
Parent - start a container!
Container - inside the container!
 
root@container:~# ipcs -q
 
------ Message Queues --------
key        msqid      owner      perms      used-bytes   messages    
0xd0d56eb2 0          hchen      644        0            0
```

但是, 如果我们运行加上了`CLONE_NEWIPC`的程序, 我们就会下面的结果:

```shell
root@ubuntu:~$ sudo./ipc
Parent - start a container!
Container - inside the container!
 
root@container:~/linux_namespace# ipcs -q
 
------ Message Queues --------
key        msqid      owner      perms      used-bytes   messages
```

可以看到`IPC`已经被隔离了.

### PID Namespace

**CLONE_NEWPID**

```C
int container_pid = clone(container_main, container_stack+STACK_SIZE, 
            CLONE_NEWUTS | CLONE_NEWPID | SIGCHLD, NULL); 
```

此时子进程的`PID`为`1`. 在传统的`Unix`系统中, `PID`为`1`的进程是`init`. 是所有进程的父进程, 有很多特权(比如: 屏蔽信号等), 另外, 其还会检查所有进程的状态.

如果某个子进程脱离了父进程(父进程没有`wait`它), 那么`init`就会负责回收资源并结束这个进程.

在子进程的`shell`里输入`ps,top`等命令, 可以看得到所有进程. 因为像`ps, top`这些命令回去读`/proc`文件系统, 所以, 因为`/proc`文件系统在父进程和子进程都是一样的, 所以这些命令显示的东西都是一样的.

### Mount Namespace

**Mount namespaces**

```C
int container_main(void* arg)
{
...
system("mount -t proc proc /proc");
...
}

int main()
{
...
int container_pid = clone(container_main, container_stack+STACK_SIZE, 
            CLONE_NEWUTS | CLONE_NEWPID | CLONE_NEWNS | SIGCHLD, NULL);
...
}
```

`pid=1`的进程是我们的`/bin/bash`.

在子进程中的`top`命令只看得到两个进程了.

在通过`CLONE_NEWNS`创建`mount namespace`后, 父进程会把自己的文件结构复制给子进程中. 而子进程中新的namespace中的所有`mount`操作都只影响自身的文件系统, 而不对外界产生任何影响.

### Docker 的 Mount Namespace

制作一个所谓的镜像.

制作`rootfs`

```C
#define _GNU_SOURCE
#include <sys/types.h>
#include <sys/wait.h>
#include <sys/mount.h>
#include <stdio.h>
#include <sched.h>
#include <signal.h>
#include <unistd.h>
 
#define STACK_SIZE (1024 * 1024)
 
static char container_stack[STACK_SIZE];
char* const container_args[] = {
    "/bin/bash",
    "-l",
    NULL
};
 
int container_main(void* arg)
{
    printf("Container [%5d] - inside the container!\n", getpid());
 
    //set hostname
    sethostname("container",10);
 
    //remount "/proc" to make sure the "top" and "ps" show container's information
    if (mount("proc", "rootfs/proc", "proc", 0, NULL) !=0 ) {
        perror("proc");
    }
    if (mount("sysfs", "rootfs/sys", "sysfs", 0, NULL)!=0) {
        perror("sys");
    }
    if (mount("none", "rootfs/tmp", "tmpfs", 0, NULL)!=0) {
        perror("tmp");
    }
    if (mount("udev", "rootfs/dev", "devtmpfs", 0, NULL)!=0) {
        perror("dev");
    }
    if (mount("devpts", "rootfs/dev/pts", "devpts", 0, NULL)!=0) {
        perror("dev/pts");
    }
    if (mount("shm", "rootfs/dev/shm", "tmpfs", 0, NULL)!=0) {
        perror("dev/shm");
    }
    if (mount("tmpfs", "rootfs/run", "tmpfs", 0, NULL)!=0) {
        perror("run");
    }
    /* 
     * 模仿Docker的从外向容器里mount相关的配置文件 
     * 你可以查看：/var/lib/docker/containers/<container_id>/目录，
     * 你会看到docker的这些文件的。
     */
    if (mount("conf/hosts", "rootfs/etc/hosts", "none", MS_BIND, NULL)!=0 ||
          mount("conf/hostname", "rootfs/etc/hostname", "none", MS_BIND, NULL)!=0 ||
          mount("conf/resolv.conf", "rootfs/etc/resolv.conf", "none", MS_BIND, NULL)!=0 ) {
        perror("conf");
    }
    /* 模仿docker run命令中的 -v, --volume=[] 参数干的事 */
    if (mount("/tmp/t1", "rootfs/mnt", "none", MS_BIND, NULL)!=0) {
        perror("mnt");
    }
 
    /* chroot 隔离目录 */
    if ( chdir("./rootfs") != 0 || chroot("./") != 0 ){
        perror("chdir/chroot");
    }
 
    execv(container_args[0], container_args);
    perror("exec");
    printf("Something's wrong!\n");
    return 1;
}
 
int main()
{
    printf("Parent [%5d] - start a container!\n", getpid());
    int container_pid = clone(container_main, container_stack+STACK_SIZE, 
            CLONE_NEWUTS | CLONE_NEWIPC | CLONE_NEWPID | CLONE_NEWNS | SIGCHLD, NULL);
    waitpid(container_pid, NULL, 0);
    printf("Parent - container stopped!\n");
    return 0;
}
```

### User Namespace

**CLONE_NEWUSER**

使用了这个参数后, 内部看到的`UID`和`GID`已经与外部不同了, 默认显示为`65534`. 那是因为容器找不到其真正的`UID`所以, 设置上了最大的`UID`(`/proc/sys/kenrel/overflowuid`).


要把容器中的`uid`和真实系统的`uid`给映射在一起, 需要修改`/proc/<pid>/uid_map`和`/proc/<pid>/gid_map`这两个文件. 格式为

`ID-inside-ns ID-outside-ns length`

* ID-inside-ns表示在容器显示的UID或GID
* ID-outside-ns表示容器外映射的真实的UID或GID
* 第三个字段表示映射的范围, 一般填1, 表示一一对应

```shell
#把真实的uid=1000映射成容器内的uid=0

$ cat /proc/2465/uid_map
         0       1000          1
```

```shell
#表示把namespace内部从0开始的uid映射到外部从0开始的uid, 其最大范围是无符号32位整形
$ cat /proc/$$/uid_map
         0          0          4294967295
```

* 写这两个的进程需要这个`namspace`中的`CAP_SETUID(CAP_SETGID)`权限
* 写入的进程必须是此`user namespace`的父或子的`user namespace`进程
* 满足如下两个条件之一:
	* 父进程将`effective uid/gid`映射到子进程的`user namespace`中
	* 父进程如果有`CAP_SETUID/CAP_SETGID`权限, 那么它将可以因三个火到父进程中的任一`uid/gid`

这样容器里是`root`, 但其实这个容器的`/bin/bash`进程是以一个普通用户来运行的. 这样, 我们容器的安全性会得到提高.

**User Namespace**是以普通用户运行, 但是别的**Namespace**需要`root`权限, 如果通知使用多个**Namespace**, 一般:

* 先用一般用户创建**User Namespace**
* 然后把这个而一般用户映射成`root`
* 在容器用`root`来创建其他的**Namespace**

### Network Namespace

一般用`ip`命令创建**Network Namespace**

![](./img/01.jpg)

上图中, Docker使用了一个私有网段, `172.40.1.0`.

使用ip link show或ip addr show来查看当前宿主机的网络情况. 会看到一个`docker0`还有一个`vethxxxx`的虚拟网卡.

```
## 首先，我们先增加一个网桥lxcbr0，模仿docker0
brctl addbr lxcbr0
brctl stp lxcbr0 off
ifconfig lxcbr0 192.168.10.1/24 up #为网桥设置IP地址
 
## 接下来，我们要创建一个network namespace - ns1
 
# 增加一个namesapce 命令为 ns1 （使用ip netns add命令）
ip netns add ns1 
 
# 激活namespace中的loopback，即127.0.0.1（使用ip netns exec ns1来操作ns1中的命令）
ip netns exec ns1   ip link set dev lo up 
 
## 然后，我们需要增加一对虚拟网卡
 
# 增加一个pair虚拟网卡，注意其中的veth类型，其中一个网卡要按进容器中
ip link add veth-ns1 type veth peer name lxcbr0.1
 
# 把 veth-ns1 按到namespace ns1中，这样容器中就会有一个新的网卡了
ip link set veth-ns1 netns ns1
 
# 把容器里的 veth-ns1改名为 eth0 （容器外会冲突，容器内就不会了）
ip netns exec ns1  ip link set dev veth-ns1 name eth0 
 
# 为容器中的网卡分配一个IP地址，并激活它
ip netns exec ns1 ifconfig eth0 192.168.10.11/24 up
 
 
# 上面我们把veth-ns1这个网卡按到了容器中，然后我们要把lxcbr0.1添加上网桥上
brctl addif lxcbr0 lxcbr0.1
 
# 为容器增加一个路由规则，让容器可以访问外面的网络
ip netns exec ns1     ip route add default via 192.168.10.1
 
# 在/etc/netns下创建network namespce名称为ns1的目录，
# 然后为这个namespace设置resolv.conf，这样，容器内就可以访问域名了
mkdir -p /etc/netns/ns1
echo "nameserver 8.8.8.8" > /etc/netns/ns1/resolv.conf
```
