# free

---

## 功能

`linux`下在终端环境下可以使用`free`命令看到系统实际使用内存的情况. `free`命令显示系统使用和空闲的内存情况，包括物理内存、交互区内存(`swap`)和内核缓冲区内存.

## 用法

`free [选项]`

## 选项参数

* `-b`以Byte为单位显示内存使用情况.
* `-k`以KB为单位显示内存使用情况.
* `-m`以MB为单位显示内存使用情况.
* `-g`以GB为单位显示内存使用情况.
* `-h`以人类可读方式显示.
	* `--si` 以`1000`代替`1024`.
* `-l, --lohi`
* `-t, --total`显示`RAM`和`swap`的总数
* `-s N, --seconds N`每`N`秒刷新一次.
* `-cN, --count N`重复刷新`N`次, 然后退出.

## 实例

```shell
[root@pingxiang-demo ~]# free
              total        used        free      shared  buff/cache   available
Mem:        1014968      652076       62708       24464      300184      172488
Swap:             0           0           0
```

* `total`内存总数.
* `used`已经使用的内存数.
* `free`空闲的内存数.
* `buff/cache`磁盘缓存的大小

**total = used + free + buff/cache**

第二部分`Swap`指的是交换分区

## free命令buff和cache的区别

### buff 缓冲区

`Buffer`: 缓冲区, 用于存储速度不同步的设备或优先级不同的设备之间的传输数据. 通过`buffer`可以减少进程间通信需要等待的时间, 当存储速度快的设备与存储速度慢的设备进行通信时, 存储慢的数据先把数据存放到`buffer`, 达到一定程度存储块的设备再地区`buffer`的数据, 在此期间存储快的设备CPU可搜索以干其他的事情.

`Buffer`: 一般是用在写入磁盘的.

==A buffer is something that has yet to be "written" to disk.==

`buffer`是用于存放要输出到disk(块设备)的数据的.

主要用于块设备缓存,例如用户目录.inode值等(ls大目录可以看到这个值增加).

### cache 缓存区

`Cache`: 缓存区, 是高速缓存, 是位于CPU和主内存之间的容量较小但速度很快的存储器, 因为CPU的速度远远高于主内存的速度, CPU从内存中读取数据需要等待很长的时间, 而Cache保存中CPU刚用户的数据或循环使用的部分数据, 这是从Cache中读取数据会更快, 技术那好了CPU等待的时间, 提高了系统的性能.

`Cache`并不是缓存文件的, 而是缓存块的(块是I/O读写最小的单元). Cache一般会用在I/O请求上, 如果多个进程要访问某个文件, 可以把此文件读入Cache中, 这样下一个进程获取CPU控制权并访问此文件直接从Cache读取, 提高系统性能.

==A cache is something that has been "read" from the disk and stored for later use.==

`cache`是存放从disk上读出的数据.

主要用于缓存文件.

### 总结

二者是为了提高`IO`性能的.

`buffer`用于缓存将要写入磁盘的数据.

`cache`用于缓存从磁盘读入的数据, 留待后面使用.

## low and high memory(`free -l`)

`high memory`只存在于32位kernel下.

### 什么是high memory，为什么要有high memory

Linux人为的把4G虚拟地址空间(32位地址最多寻址4G)分为3G＋1G, 其中0～3G为用户程序地址空间, 3G～4G为kernel地址空间, 这就是说kernel最多寻址1G的虚拟地址空间.

当CPU启用MMU(内存管理)的paging机制后, CPU访问的是虚拟地址, 然后由MMU根据"页表"转换成物理地址. "页表"是由kernel维护的, 所以kernel可以决定1G的虚拟地址空间具体映射到什么物理地址. 但是kernel最多只有3G～4G这1G地址空间, 所以不管kernel怎么映射, 最多只能映射1G的物理内存. 所以如果一个系统有超过1G的物理内存, 在某一时刻, 必然有一部分kernel是无法直接访问到的. 另外，kernel除了访问内存外，还需要访问很多IO设备. 在现在的计算机体系结构下, 这些IO设备的资源(比如寄存器，片上内存等)一般都是通过MMIO(内存映射I/O)的方式映射到物理内存地址空间来访问的, 就是说kernel的1G地址空间除了映射内存, 还要考虑到映射这些IO资源--换句话说，kernel还需要预留出一部分虚拟地址空间用来映射这些IO设备(ioremap就是干这个的).

Linux kernel采用了最简单的映射方式来映射物理内存, 即把物理地址＋3G按照线性关系直接映射到kernel空间. 考虑到一部分kernel虚拟地址空间需要留给IO设备(以及一些其他特殊用途), Linux kernel最多直接映射896M物理内存, 而预留了最高端的128M虚拟地址空间给IO设备. 所以, 当系统有大于896M内存时, 超过896M的内存kernel就无法直接访问到了, **这部分内存就是high memory**. 那kernel就永远无法访问到超过896M的内存了吗? 不是的, kernel已经预留了128M虚拟地址, 我们可以用这个地址来动态的映射到high memory, 从而来访问high memory. 所以预留的128M除了映射IO设备外，还有一个重要的功能是提供了一种动态访问high memory的一种手段(kmap主要就是干这个的，当然还有vmalloc).

要理解high memory, 关键是把物理内存管理, 虚拟地址空间管理, 以及两者间的映射(页表管理)三个部分分开考虑, 不要把物理内存管理和虚拟地址空间管理混在一起. 比如high memory也参与kernel的物理内存分配, 你调用get_page得到的物理页有可能是low memory, 也可以是high memory, 这个物理页可以被映射到kernel, 同时也可以被映射到user space. 再比如vmalloc, 只保证返回的虚拟地址是在预留的vmalloc area里, 对应的物理内存, 可以是low memory, 也可以是high memory. 当然出于性能考虑, kernel可能会优先分配直接映射的low memory, 但我们不能假设high memory就不会被分配到.

64位系统下不会有high memory, 因为64位虚拟地址空间非常大(分给kernel的也很大), 完全能够直接映射全部物理内存.