# lsof

---

##  简介

`lsof`(list open files)命令用于查看你进程开打的文件, 打开文件的进程, 进程打开的端口(TCP、UDP). 找回/恢复删除的文件. 是十分方便的系统监视工具, 因为lsof命令需要访问核心内存和各种文件, 所以需要root用户执行.

在linux环境下, 任何事物都以文件的形式存在, 通过文件不仅仅可以访问常规数据, 还可以访问网络连接和硬件. 所以如传输控制协议(TCP)和用户数据报协议(UDP)套接字等, 系统在后台都为该应用程序分配了一个文件描述符, 无论这个文件的本质如何, 该文件描述符为应用程序与基础操作系统之间的交互提供了通用接口. 因为应用程序打开文件的描述符列表提供了大量关于这个应用程序本身的信息, 因此通过lsof工具能够查看这个列表对系统监测以及排错将是很有帮助的.

## 选项

## 示例

```
lsof

COMMAND     PID   TID    USER   FD      TYPE             DEVICE  SIZE/OFF       NODE NAME
systemd       1          root  cwd       DIR              253,1      4096          2 /
systemd       1          root  rtd       DIR              253,1      4096          2 /
systemd       1          root  txt       REG              253,1   1523568    1053845 /usr/lib/systemd/systemd
systemd       1          root  mem       REG              253,1     20040    1050452 /usr/lib64/libuuid.so.1.3.0
systemd       1          root  mem       REG              253,1    261336    1051899 /usr/lib64/libblkid.so.1.1.0
systemd       1          root  mem       REG              253,1     90664    1050435 /usr/lib64/libz.so.1.2.7
systemd       1          root  mem       REG              253,1    157424    1050447 /usr/lib64/liblzma.so.5.2.2
systemd       1          root  mem       REG              253,1     23968    1050682 /usr/lib64/libcap-ng.so.0.0.0
systemd       1          root  mem       REG              253,1     19888    1050666 /usr/lib64/libattr.so.1.1.0
systemd       1          root  mem       REG              253,1     19776    1049995 /usr/lib64/libdl-2.17.so
systemd       1          root  mem       REG              253,1    402384    1050423 /usr/lib64/libpcre.so.1.2.0
systemd       1          root  mem       REG              253,1   2127336    1049989 /usr/lib64/libc-2.17.so
systemd       1          root  mem       REG              253,1    144792    1050015 /usr/lib64/libpthread-2.17.so
systemd       1          root  mem       REG              253,1     88720    1048596 /usr/lib64/libgcc_s-4.8.5-20150702.so.1
systemd       1          root  mem       REG              253,1     44448    1050019 /usr/lib64/librt-2.17.so
systemd       1          root  mem       REG              253,1    269416    1052057 /usr/lib64/libmount.so.1.1.0
systemd       1          root  mem       REG              253,1     91800    1051447 /usr/lib64/libkmod.so.2.2.10
systemd       1          root  DEL       REG              253,1              1050684 /usr/lib64/libaudit.so.1.0.0;5b63ccad
systemd       1          root  mem       REG              253,1     61672    1052186 /usr/lib64/libpam.so.0.83.1
systemd       1          root  mem       REG              253,1     20032    1050671 /usr/lib64/libcap.so.2.22
systemd       1          root  DEL       REG              253,1              1050432 /usr/lib64/libselinux.so.1;5b63ccad
systemd       1          root  mem       REG              253,1    164264    1049982 /usr/lib64/ld-2.17.so
systemd       1          root    0u      CHR                1,3       0t0       1031 /dev/null
systemd       1          root    1u      CHR                1,3       0t0       1031 /dev/null
systemd       1          root    2u      CHR                1,3       0t0       1031 /dev/null
systemd       1          root    3u     unix 0xffff880231f93800       0t0     371466 socket
systemd       1          root    4u  a_inode                0,9         0       5880 [eventpoll]
systemd       1          root    5u  a_inode                0,9         0       5880 [signalfd]
systemd       1          root    6r      DIR               0,21         0       7413 /sys/fs/cgroup/systemd
```

* `COMMAND`：进程的名称
* `PID`：进程标识符
* `USER`：进程所有者
* `FD`：文件描述符，应用程序通过文件描述符识别该文件.
	* `cwd`：表示`current work dirctory`，即: 应用程序的当前工作目录, 这是该应用程序启动的目录, 除非它本身对这个目录进行更改.