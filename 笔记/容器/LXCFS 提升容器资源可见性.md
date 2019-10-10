# LXCFS 提升容器资源可见性

## 简介

![](./img/lxcfs.png)

Linuxs利用Cgroup实现了对容器的资源限制, 但在容器内部依然缺省挂载了宿主机上的procfs的/proc目录. 其包含如: meminfo, cpuinfo, stat, uptime等资源信息。一些监控工具如free/top或遗留应用还依赖上述文件内容获取资源配置和使用情况. 当它们在容器中运行时, 就会把宿主机的资源状态读取出来, 引起错误和不便.

利用`lxcfs`来提供容器中的资源可见性.

LXCFS通过用户态文件系统, 在容器中提供下列`procfs`的文件:

```
/proc/cpuinfo
/proc/diskstats
/proc/meminfo
/proc/stat
/proc/swaps
/proc/uptime
```

## 操作

### 安装 lxcfs

```
wget https://copr-be.cloud.fedoraproject.org/results/ganto/lxd/epel-7-x86_64/00486278-lxcfs/lxcfs-2.0.5-3.el7.centos.x86_64.rpm

yum install -y lxcfs-2.0.5-3.el7.centos.x86_64.rpm
```

### 测试

按照默认限制启动容器, 查询可见内存还是主机内存(我主机内存是1G)

```
docker run -it -m 256m ubuntu:16.04 /bin/bash
root@a33a0f2467e6:/# free -h
              total        used        free      shared  buff/cache   available
Mem:           992M        106M         88M        408K        797M        709M
Swap:            0B          0B
```

挂载`lxcfs`目录启动容器, 可见容器内部使用`free`查询内存容量正常.

```
docker run -it -m 256m \
      -v /var/lib/lxcfs/proc/cpuinfo:/proc/cpuinfo:rw \
      -v /var/lib/lxcfs/proc/diskstats:/proc/diskstats:rw \
      -v /var/lib/lxcfs/proc/meminfo:/proc/meminfo:rw \
      -v /var/lib/lxcfs/proc/stat:/proc/stat:rw \
      -v /var/lib/lxcfs/proc/swaps:/proc/swaps:rw \
      -v /var/lib/lxcfs/proc/uptime:/proc/uptime:rw \
      ubuntu:16.04 /bin/bash
root@73a7409c4565:/# free -h
              total        used        free      shared  buff/cache   available
Mem:           256M        512K        254M        412K        708K        254M
Swap:          256M          0B        256M
```
      
      
      
     