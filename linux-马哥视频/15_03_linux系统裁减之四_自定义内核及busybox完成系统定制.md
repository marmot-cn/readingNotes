#15_03_Linux系统裁减之四 自定义内核及busybox完成系统定制

###笔记

---

###busybox

www.busybox.net

靠一个二进制程序模拟实现许多命令.

debian 和 ubuntu使用 busybox 完成系统启动过程.

####查看本机硬件设备信息

**查看cpu信息**

`cat /proc/cpuinfo`

		[chloroplast@dev-server-2 ~]$ cat /proc/cpuinfo
		processor	: 0
		vendor_id	: GenuineIntel
		cpu family	: 6
		model		: 62
		model name	: Intel(R) Xeon(R) CPU E5-2650 v2 @ 2.60GHz
		stepping	: 4
		microcode	: 0x428
		cpu MHz		: 2593.810
		cache size	: 20480 KB
		physical id	: 0
		siblings	: 1
		core id		: 0
		cpu cores	: 1
		apicid		: 0
		initial apicid	: 0
		fpu		: yes
		fpu_exception	: yes
		cpuid level	: 13
		wp		: yes
		flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat clflush mmx fxsr sse sse2 ht syscall nx rdtscp lm constant_tsc rep_good nopl pni ssse3 cx16 sse4_1 sse4_2 popcnt aes hypervisor lahf_lm
		bogomips	: 5187.62
		clflush size	: 64
		cache_alignment	: 64
		address sizes	: 46 bits physical, 48 bits virtual
		power management:

**列出本机的usb接口**

`lsusb`

		[chloroplast@dev-server-2 ~]$ lsusb
		Bus 001 Device 002: ID 0627:0001 Adomax Technology Co., Ltd
		Bus 001 Device 001: ID 1d6b:0001 Linux Foundation 1.1 root hub

**列出pci总线设备**

`lspci`

		[chloroplast@dev-server-2 ~]$ lspci
		00:00.0 Host bridge: Intel Corporation 440FX - 82441FX PMC [Natoma] (rev 02)
		00:01.0 ISA bridge: Intel Corporation 82371SB PIIX3 ISA [Natoma/Triton II]
		00:01.1 IDE interface: Intel Corporation 82371SB PIIX3 IDE [Natoma/Triton II]
		00:01.2 USB controller: Intel Corporation 82371SB PIIX3 USB [Natoma/Triton II] (rev 01)
		00:01.3 Bridge: Intel Corporation 82371AB/EB/MB PIIX4 ACPI (rev 01)
		00:02.0 VGA compatible controller: Cirrus Logic GD 5446
		00:03.0 Unassigned class [ff80]: XenSource, Inc. Xen Platform Device (rev 01)

**查看每一个硬件的详细信息**

`hal-device`: 硬件抽象层.Hadeare Abstract Layer

####编译内核:

**配置**

* make menuconfig
* make gconfig
* make kconfig
* make oldconfig
* make config

配置完成以后保存成`.config`

**make**

**make modules_install:安装模块**

**make install**

模块的安装位置: `/lib/modules/KERNEL_VERSION/`

####如何实现部分编译:

**只编译某子目录下的相关代码**

* `make dir/` 指定某个目录,只编译该目录下的代码
	
		make arch/ 只编译内核核心(不包括模块的内核核心)
		make drivers/ 只编译驱动
		make drivers/net/ 只编译和网络相关的驱动

**只编译部分模块**

		make M=drivers/net/

**只编译某一个模块**

		make drivers/net/pcnet32.ko

**将编译完成的结果放置于别的目录中**

		make O=/tmp/kernel ,将编译完成的结果放在/tmp/kernel目录内

####如何编译busybox

`mdev -s` 探测额外的硬件.

`/sysroot`下的busybox移至另一个目录,以实现与真正的根文件系统分开制作.我们这里选择使用`/mnt/temp`目录:

		# mkdir -pv /tmp/busybox
		# cp -r /mnt/sysroot/*	/tmp/busybox
		
#####制作initrd

		# cd /tmp/busybox
		
**建立rootfs**

		#mkdir -pv proc	sys	etc/init.d	tmp	dev	mnt/sysroot
		
**创建两个必要的设备文件**

		#mknod	dev/console  c	5	1 (c:字符设备 5:主设备号 1:副设备号 )
		#mknod	dev/null	c	1	3

**为initrd制作init程序,此程序的主要任务是实现rootfs的切换,因此,可以以脚本的方式来实现它:**

		#rm linuxrc
		#vim init
		添加如下内容
		#!/bin/sh
		mount -t proc proc /proc
		mount -t sysfs	sysfs	/sys
		mdev -s
		mount -t ext3	/dev/hda2	/mnt/sysroot
		exec	switch_root	/mnt/sysroot	/sbin/init
		
给此脚本执行权限

		chmod +x init
		
**制作initrd**

		#find . | cpio --quit -H newc -o | gzip -9 -n > /mnt/boot/initrd.gz

###整理知识点

---

####什么是pci总线
