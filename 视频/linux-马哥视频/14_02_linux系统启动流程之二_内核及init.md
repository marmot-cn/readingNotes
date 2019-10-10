#14_02_Linux系统启动流程之二 内核及init

###笔记

---

####查看运行级别

**`runlevel`**

查看运行级别

		[root@iZ94xwu3is8Z ~]# runlevel
		N 3
		
* `N`: 上一个级别(N,表示上次没有运行级别,直接启动了)
* `3`: 当前级别

**`who -r`**

		[root@iZ94xwu3is8Z ~]# who -r
         运行级别 3 2016-01-11 15:12
    
####查看内核release号
     
**`uname -r`**

		[root@iZ94xwu3is8Z ~]# uname -r
		3.10.0-327.4.4.el7.x86_64
		
####安装(修复)grub stage1:

**`grub`**

		# grub
		grub> root (hd1,0)   //指定内核所在的分区(hd1:硬盘,0:boot所在的分区)
		grub> setup (hd1)  //指定安装的硬盘
		
**`grub-install`**

`grub-install --root-directory=/path/to/boot's_parent_dir`(grub分区对应的根) `/PATH/TO/DEVICE`(指定硬盘)

####grub.conf grub配置文件损坏

开启自动进入`grub`命令行

		grub> find (hd0,0)/    //敲tab键,显示文件，
		grub> root (hd0,0)     //指定root
		grub> kernel /vmlinuz-xxx  //指定kernel
		grub> initrd /initrd-xxx   //指定initrd,和内核版本需要一致
		grub> boot  				//启动系统
		

####Kernel初始化的过程

1. 设备探测
2. 驱动初始化 (可能会从`initrd(initramfs)`文件中装载驱动模块)
3. 以`只读`挂载根文件系统,此时只读是为了`安全`,如果是读写挂载,任何一个bug都会导致根文件系统崩溃.所以此时只是用于读取文件.随后`init`会重新把根挂载成可读写的
4. 装载第一个进程init(PID:1)

####init

`/sbin/init`: 读取 `/etc/inittab`

红帽6.0以后使用的是`upstart`(ubuntu开发的)另外一个版本的`init`.会并行的启动很多进程,比传统的`init`快.基于`d-bus`配置,`event-driven`事件驱动

`systemd` 比 `upstart` 速度快.完全并行.

我本机阿里云使用的是centos7,使用的是`systemd`

`/etc/init`是`/etc/inittab`文件切割后的各种片(红帽6上)

`/etc/inittab`的格式(冒号隔开的4段)
		
		man inittab
		id:runlevels:action:process
		
* `id`: 标识符,唯一的即可
* `runlevels`: 在哪些运行级别下会执行此命令
* `action`: 在什么情况下执行此行
* `process`: 要执行的进程		

		id:3:initdefault:
		在运行级别3设定默认级别(initdefault),不需要运行程序
		
		si::sysinit:/etc/rc.d/rc.sysinit
		系统刚启动时候(系统初始化运行)运行/etc/rc.d/rc.sysinit
		
**action**

* `initdefault`: 设定默认运行级别
* `sysinit`: 系统初始化,只进行一次
* `wait`: 等待级别切至此级别时执行,可能会执行多次,	
* `trl-alt-del`: 同时按着三个键时触发
* `powerfail`: 停电(`UPS`)
* `powerokwait`: 停电时电力恢复
* `respawn`: 一旦程序终止会重新启动,执行多次

**`/etc/rc.d/rc.sysint`完成的任务**

1. 激活`udev`和`selinux`
2. 根据`/etc/sysctl.conf`文件,来设定内核参数
3. 设定系统时钟
4. 装载键盘映射
5. 启用交换分区
6. 设置主机名
7. 根文件系统检测,并以读写方式重新挂载
8. 激活`RAID`和`LVM`设备
9. 启用磁盘配额
10. 根据`/etc/fstab`,检查并挂载其它文件系统
11. 清理过期的锁和PID文件

**/etc/rc#.d/**

`#`代表不同的服务级别.

`K`开头文件代表关闭. 
`S`开头文件代表启动.

`K##`和`S##`.

`##`关闭或启动的优先次序,数字越小越优先被选定.线关闭以`K`开头的服务,后启动以`S`开头的服务.

所有文件都是软连接,连接到`/etc/init.d`目录下的文件

		[root@iZ94xwu3is8Z init.d]# ls -l /etc/rc3.d/
		总用量 0
		lrwxrwxrwx. 1 root root 20 11月 21 2014 K50netconsole -> ../init.d/netconsole
		lrwxrwxrwx. 1 root root 17 11月 21 2014 S10network -> ../init.d/network
		lrwxrwxrwx  1 root root 15 1月  11 14:41 S50aegis -> ../init.d/aegis
		lrwxrwxrwx  1 root root 20 1月  11 14:28 S98agentwatch -> ../init.d/agentwatch

###整理知识点

---