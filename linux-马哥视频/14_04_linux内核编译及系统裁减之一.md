#14_04_Linux内核编译及系统裁减之一

###笔记

---

内核由`核心`(`/boot/vmlinuz-version`)和`内核模块`(ko Kernel Object)(`/lib/modules/version`)组成.

**内核设计**

* 单内核(Linux是单内核设计,但是引用了微内核的设计思想)
* 微内核

**装载模块**

* `insmod`
* `modprobe`

**用户控件访问,监控内核的方式**

* `/proc`
* `/sys`

都是伪文件系统(实际上不存在),文件是内核参数.

* `/proc/sys`: 此目录中的文件很多是可读写的.
* `/sys/`: 某些文件可写

**设定内核参数值的方法**

		echo VALUE > /proc/sys/TO/SOMEFILE
		
`sysctl -w xxxx.xxx`(路径,用`.`分隔)
		
		echo xxxx > /proc/sys/kernel/hostname
		sysctl -w kernel.hostname=xxxx 修改主机名 
		
立即生效,但不能永久有效.

永久有效,但不能立即生效: `/etc/sysctl.conf`  
修改文件完成之后,执行如下命令可立即生效:

		sysctl -p
		
`sysctl -a`: 显示内核参数

**内核模块管理**

* `lsmod`: 查看
* `modprobe`:
	* `modprobe MOD_NAME`: 装载某模块
	* `modprobe -r MOD_NAME`: 卸载某模块
* `modinfo MOD_NAME`: 查看模块的具体信息
* `insmod /PATH/TO/MODULE_FILE`: 装载模块
* `rmmod MOD_NAME`: 移除模块
* `depmod /PATH/TO/MODULES_DIR`: 生成模块依赖到指定文件夹

内核中的功能除了核心功能之外,在编译时,达多功能都有三种选择:

1. 不使用此功能
2. 编辑成内核模块
3. 编译进内核

如何手动编译内核

`/boot/config-xxx`是红帽编译内核的官方配置文件.

内核文件夹内的`/.config`隐藏文件,是自己的配置文件.

在内核目录下运行`make menuconfig`命令,选择内核特性.

运行`make`编译

`make modules_install`安装内核模块

`make install`安装内核

**screen命令**

* `screen -ls`: 显示已经建立的屏幕
* `screen`: 直接打开一个新的屏幕
* `Ctrl+a, d`:拆除屏幕
* `screen -r ID`:还原回某屏幕
* `exit`: 退出

`screen`可以防止远程连接中断

**二次编译时清理,清理前,如果有需要,请备份配置文件.config**

`make clean`

`make mrproper`

**mkinitrd**

mkinitrd initrd文件路径  内核版本号

		mkinitrd /boot/initrd-`uname -r`.img `uname -r`
		
**复制依赖库(ldd)文件的脚本**

字符串截取的一些示例:

		[root@iZ94xwu3is8Z ~]# echo $FILE
		/usr/local/src
		[root@iZ94xwu3is8Z ~]# echo ${FILE#*/}
		usr/local/src
		[root@iZ94xwu3is8Z ~]# echo ${FILE%*/}
		/usr/local/src

脚本文件代码:

		#!/bin/bash
		#
		DEST=/xxx/xxx
		libcp(){
		  LIBPATH=${1%/*}
		  [! -d $DEST$LIBPATH] && mkdir -p $DEST$LIBPATH
		  [! -e $DEST${1}] && cp $1 $DEST$LIBPATH && echo "copy lib  $1 finished."
		}
		
		bincp(){
		  CMDPATH=${1%/*}
		  [! -d $DEST$CMDPATH] && mkdir -p $DEST$CMDPATH
		  [! -e $DEST${1}] && cp $1 $DEST$CMDPATH
		  
		  #把依赖库文件提取出来 
		  for LIB in `ldd $1 | grep -o "/.*lib\(64\)\{0,1\}/[^[:space:]]\{1,\}"`
		    libcp $LIB
		  done
		}
		
		read -p "Your command:" CMD
		until [$CMD == 'q']; do
		  ! which $CMD && echo "wrong command" && read -p "Input again" CMD && continue
		  COMMAND=`which $CMD | grep -v "^alias" | grep -o "[^[:space:]]\{1,\}"`
		  bincp $COMMAND
		  echo "copy $COMMAND finished."
		  read -p "Continue: " CMD
	    done

###整理知识点

---