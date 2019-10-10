# Linux 内存文件系统

虚拟内核文件系统(`VirtualKernel File Systems`), 是指那些是由内核产生但不存在于硬盘上(存在于内存中)的文件系统.

## proc

`proc`文件系统为操作系统本身和应用程序之间的通信提供了一个安全的接口. 通过它里面的一些文件, 可以获取系统状态信息并修改某些系统的配置信息. 当我们在内核中添加了新功能或设备驱动时, 经常需要得到一些系统状态的信息, 一般这样的功能需要经过一些像ioctl()这样的系统调用来完成.

## devfs

`/dev`目录下的每一个文件都对应的是一个设备, `devfs`也是挂载于`/dev`目录下. 在2.6内核以前使用`devfs`来提供一种类似于文件的方法来管理位于`/dev`目录下的所有设备. 但是devfs文件系统有一些缺点, 有时一个设备映射的设备文件可能不同. 例如, 我的U盘可能对应sda, 也可能对应sdb, 没有足够的主/辅设备号, 当设备过多的时候, 显然这会成为一个问题.

## sysfs

挂载于`/sys`目录下. `sysfs`文件系统把连接在系统上的**设备和总线**组织成为一个分级的文件, 用户空间的程序同样可以利用这些信息, 以实现和内核的交互. `sysfs`文件系统是当前系统上实际设备树的一个直观反映, 它是通过kobject子系统来建立这个信息的, 当一个kobject被创建的时候, 对应的文件和目录也就被创建了.

## tmpfs

`tmpfs(temporary filesystem)`是Linux特有的文件系统, 标准挂载点是`/dev/shm`, 默认大小是实际内存的一半, 如下所示. 当然, 用户也可以将`tmpfs`挂载在其他地方. tmpfs可以使用物理内存, 也可以使用swap交换空间.

```
[ansible@demo ~]$ sudo df -h
Filesystem             Size  Used Avail Use% Mounted on
/dev/vda1               99G  1.7G   92G   2% /
devtmpfs               3.9G     0  3.9G   0% /dev
tmpfs                  3.9G     0  3.9G   0% /dev/shm
tmpfs                  3.9G  388K  3.9G   1% /run
tmpfs                  3.9G     0  3.9G   0% /sys/fs/cgroup
/dev/mapper/data-data  190G   33M  190G   1% /data
tmpfs                  783M     0  783M   0% /run/user/1000
```

tmpfs有些像虚拟磁盘(`ramdisk`), 但ramdisk是一个块设备, 而且需要一个mkfs之类的命令格式化后才能使用. 而tmpfs是一个独立的文件系统, 不是块设备, 只要挂载, 就可以立即使用. 下面是`tmpfs`最主要的几个特点:

* 临时性: 由于`tmpfs`是构建在内存中的, 所以存放在`tmpfs`中的所有数据在卸载或断电后都会丢失.
* 快速读写能力: 内存的访问速度要远快于磁盘I/O操作, 即使使用了swap, 性能仍然非常卓越
* 动态收缩: tmpfs一开始使用很小的空间, 但随着文件的复制和创建, `tmpfs`文件系统会分配更多的内存, 并按照需求动态地增加文件系统的空间. 而且, 当`tmpfs`中的文件被删除时, tmpfs文件系统会动态地减小文件并释放内存资源.

### 挂载一个新的内存文件系统

挂载`100MB`权限为`0755`

```
mount -ttmpfs -o size=100M,mode=0755 tmpfs /data/cache
```

扩容到`2000MB`

```
mount -o remount,size=2000M /data/cache
```

可以放到`/etc/fstab`中.

```
tmpfs /data/cache tmpfssize=100M,mode=0755 0 0
```

我们使用`dd`写入文件前

```
[ansible@demo cache]$ free -m
              total        used        free      shared  buff/cache   available
Mem:           7822         152        6895           0         775        7396
Swap:             0           0
```

我们使用`dd`写入文件后

```
[ansible@demo cache]$ sudo dd if=/dev/zero of=test2 bs=1M count=1000
1000+0 records in
1000+0 records out
1048576000 bytes (1.0 GB) copied, 0.578326 s, 1.8 GB/s

[ansible@demo cache]$ free -h
              total        used        free      shared  buff/cache   available
Mem:           7.6G        149M        4.8G        2.0G        2.7G        5.3G
Swap:            0B          0B          0B
```