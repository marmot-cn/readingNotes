#17_02_Linux系统故障排除

---

###常见的系统故障排除

1. 确定问题的故障特征
2. 重现故障
3. 使用工具收集进一步信息
4. 排除不可能的原因
5. 定位故障:
	* 从简单问题入手
	* 一次尝试一种方式

####两个原则

1. 任何修改时备份源文件
2. 尽可能借助工具

####可能会出现的故障

1. 管理员密码忘记,单用户模式进入修改密码
2. 系统无法正常启动
	* grub损坏（MBR损坏, grub配置文件丢失）
	* 系统初始化故障(某文件系统无法正常挂在,驱动不兼容)
		* 进入grub的编辑模式,然后进入服务级别1(emergency)模式
	* 服务故障
	* 用户无法登陆系统
		* 用户对应shell文件损坏
3. 命令无法运行
		
		//重新手动指定PATH
		export PATH=/bin:/sbin:/usr/bin:/usr/sbi.... 
		
		退出当前登录,另启终端,重新登录
		
		使用/xxx/bin/vim 编辑配置文件 /etc/profile
		
4. 编译过程无法继续
	* 开发环境配置不正确,缺少基本组件
 
**MBR损坏**

MBR损坏:

1. 借助别的主机修复
2. 使用紧急救援模式: 进入原来的根文件系统
	* boot.iso 
	* 使用完整的系统安装光盘
			
			进入紧急救援模式
			boot: linux rescure
			/mnt/sysimage
			
chroot /mnt/sysroot 虚根,不能识别dev设备文件, 需要手动创建dev设备文件(`mknod`).

**MBR损坏示例** 
 
备份:

		dd if=/dev/sda of=/root/mbr.backup count=1 bs=512
		
损坏MBR:

		dd if=/dev/zero of=/dev/sda count=1 bs=200
		
重启系统:

		shutdown -r now
		
通过光盘(linux 系统盘)启动(光盘为第一启动),进入救援模式:		
		
	  	boot: linux rescue 
	  	
进入grub的命令行模式:

		grub
		
寻找并内核所在的分区指定根(root):
		
		//(hd0,0) 第一个硬盘第一分区,(hd2,3)表示第三个硬盘第四分区
		
		grub> find (hd0,0)/
		
		grub> root (hd0,0)
		
安装grub:

		grub> setup (hd0)
		
重启系统

		grub> reboot
		
**grub配置文件丢失**

		grub> root (hd0,0)
		grub> kernel /vmlinuz-xxx 
		grub> initrd /initrd-xxx
		grub> boot
		
`grub文件内容,grub.conf`

		default=0
		timeout-10
		title xxxx
			root (hd0,0)
			kenenl /vmlinuz-xxx ro root=xxxx
			initrd /initrd-xxx

**不小心把默认级别设定为0(停机)或6(重启)**

进入`1`级别,单用户模式,编辑inittab文件.

**不小心删除/etc/rc.d/rc3.d,3(多用户模式)**

进入单用户模式,修改目录系统

**某个服务故障导致启动停止**

例如`sendmail`,配置文件时间戳检查无法通过.

1. 重新启动,进入单用户模式,关掉`sendmail`服务,不让`sendmail`服务自动启动.
2. 修复`sendmail`时间戳.
3. 在系统启动时候按着`I`键,进行服务启动`交互模式`(每个服务是否启动与用户交互).

**系统初始化过程**

`POST`(加电对硬件进行检测)-->`BIOS`(启动设备顺序依次找其MBR中的bootloader)->`Kernel`(initrd, rootfs, /sbin/init)-->/etc/inittab

**`/etc/rc.d/rc.local`**

一般对应s99(`/etc/rc.d/rc3.d/S99local`)服务(我的阿里云centos7没有),系统启动后会执行.一般一些需要开机启动且不能作为服务的任务放在该文件内开机执行.

如果该文件内写了无法正常执行的命令都会导致系统无法正常启动.

`etc`下的`rc.local`默认指向`rc.d/rc.local`

		[chloroplast@iZ94jqmwawyZ etc]$ ll rc.local
		lrwxrwxrwx 1 root root 13 7月  13 12:28 rc.local -> rc.d/rc.local
		[chloroplast@iZ94jqmwawyZ etc]$ pwd
		/etc


**`rc.local`脚本语法错误或出现逻辑错误**

单用户模式不会启动该服务,进入单用户模式内修改该服务.

**用户对应shell文件损坏下的修复**

单用模式也需要运行shell,所以1级别也无法进入.

进入紧急救援模式,重新安装bash.


挂在光盘来安装:

		rmp --ivh --replacepkgs --root /mnt/sysimage bash-xxxx.rpm
		
基于网络来安装

**修复PATH环境变量错误引起的命令无法执行问题**

修改`PATH`

		export PATH=/data/bin
		
执行命令`ls`
		
		ls
		-bash: ls: command not found
		
重新登录即可


		
###整理知识点

---

####/etc/profile
