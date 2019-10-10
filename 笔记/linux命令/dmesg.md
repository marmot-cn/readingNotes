# dmesg

---

20180203在询问周老师他们查看`rancher`网络不通时通过`dmesg`查知内核报错而发现`arp表`写满的问题. 今天就记录一下`dmesg`命令的用法.

## 简介

`dmesg`命令用来显示开机信息, `kernel`会将开机信息存储在`ring buffer`中. 开机时来不及查看信息, 可利用`dmesg`来查看. 开机信息亦保存在`/var/log/dmesg`.

## 语法

```
dmesg [-cn][-s <缓冲区大小>]
```

* `-c`: 显示信息后, 清除`ring buffer`中的内容.
* `-s<缓冲区大小>`: 默认设置为8196, 与内核的默认`syslog`缓冲区大小一致.
* `-n`设置记录信息的层级; 比如: `-n 1`为最低级, 只向控制台显示内核恐慌信息.

## 作用

1. `dmesg`命令显示Linux内核的环形缓冲区信息, 我们可以从中获得诸如系统架构、CPU、挂载的硬件, RAM等多个运行级别的大量的系统信息. 当计算机启动时, 系统内核(操作系统的核心部分)将会被加载到内存中. 在加载的过程中会显示很多的信息, 在这些信息中我们可以看到内核检测硬件设备.
2. `dmesg`命令设备故障的诊断是非常重要的. 在`dmesg`命令的帮助下进行硬件的连接或断开连接操作时, 我们可以看到硬件的检测或者断开连接的信息.

![dmesg](./img/dmesg.png)

## 备注

`dmesg`用来显示内核环缓冲区(`kernel-ring buffer`)内容, 内核将各种消息存放在这里. 在系统引导时, 内核将与硬件和模块初始化相关的信息填到这个缓冲区内. 内核环缓冲区中的消息对于诊断系统问题通常非常有用.

## 示例

括号里面的数字为`timestamp`, 时间戳, 该时间为系统从开机到现在的运行时间, 单位为秒.

```
[ansible@localhost ~]$ dmesg
[    0.000000] Initializing cgroup subsys cpuset
[    0.000000] Initializing cgroup subsys cpu
[    0.000000] Initializing cgroup subsys cpuacct
[    0.000000] Linux version 3.10.0-514.el7.x86_64 (builder@kbuilder.dev.centos.org) (gcc version 4.8.5 20150623 (Red Hat 4.8.5-11) (GCC) ) #1 SMP Tue Nov 22 16:42:41 UTC 2016
[    0.000000] Command line: BOOT_IMAGE=/vmlinuz-3.10.0-514.el7.x86_64 root=/dev/mapper/cl-root ro crashkernel=auto rd.lvm.lv=cl/root rd.lvm.lv=cl/swap rhgb quiet LANG=en_US.UTF-8
[    0.000000] e820: BIOS-provided physical RAM map:
[    0.000000] BIOS-e820: [mem 0x0000000000000000-0x000000000009fbff] usable
[    0.000000] BIOS-e820: [mem 0x000000000009fc00-0x000000000009ffff] reserved
[    0.000000] BIOS-e820: [mem 0x00000000000f0000-0x00000000000fffff] reserved
```