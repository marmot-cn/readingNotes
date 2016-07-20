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

`/etc/profile`文件的改变会涉及到系统的环境,也就是有关Linux环境变量的东西.

在文件中`:`表示并列含义,有多个的话用`:`分离.

在文件中`.`表示操作的当前目录.

**常见的环境变量**

* `PATH`: 决定了`shell`将到哪些目录中寻找命令或程序.
* `HOME`: 当前用户主目录.
* `MAIL`: 当前用户的邮件存放目录.
* `SHELL`: 当前用户用的是哪种`Shell`.
* `HISTSIZE`: 是指保存历史命令记录的条数.
* `LOGNAME`: 是指当前用户的登录名.
* `HOSTNAME`: 主机名, 许多应用程序如果要用到主机名的话, 通常是从这个环境变量中来取得的.
* `LANG/LANGUGE`: 是和语言相关的环境变量,使用多种语言的用户可以修改此环境变量.
* `PS1`: 基本提示符,对于`root`用户是`#`,对于普通用户是`$`.
* `PS2`: 附属提示符,默认是"`>`"

#####几个文件的区别

**`/etc/profile`**

用来设置`系统环境参数`,比如`$PATH`. 这里面的环境变量是对系统内`所有用户`生效的.

此文件为系统的每个用户设置环境信息,当用户第一次登录时,该文件被执行,并从`/etc/profile.d`目录的配置文件中搜集shell的设置.

**`/etc/bashrc`**

这个文件设置系统`bash shell`相关的东西,对系统内所有用户生效,只要用户运行bash命令,那么这里面的东西就在起作用.

为每一个运行bash shell的用户执行此文件,当bash shell被打开时,该文件被读取.

**`~/.bash_profile`**

用来设置一些环境变量,功能和`/etc/profile`类似,但是这个是针对用户来设定的,也就是说,你在`/home/user1/.bash_profile`中设定了环境变量只针对`user1`这个用户生效.

每个用户都可使用该文件输入专用于自己使用的shell信息,当用户登录时,该文件仅仅执行一次.默认情况下,他设置一些环境变量,执行用户的`.bashrc`文件.

交互式,`login`方式进入`bash`运行的,意思是只有用户登录时才会生效.

**`~/.bashrc`**

作用类似于`/etc/bashrc`,只是针对用户自己而言,不对其他用户生效.

文件包含专用于你的bash shell的bash信息,当登录时以及每次打开新的shell时,该文件被读取.

交互式`non-login`方式进入`bash`运行的,用户不一定登录,只要以该用户身份运行命令行就会读取该文件.

**执行顺序**

`/etc/profile -> (~/.bash_profile | ~/.bash_login | ~/.profile) -> ~/.bashrc -> /etc/bashrc -> ~/.bash_logout`

`~/.bash_profile`文件中一般会有下面的代码(执行`~/.bashrc`)执行:

		if [ -f ~/.bashrc ] ; then
		. ~/.bashrc
		fi 

`~/.bashrc`中,一般还会有以下代码(执行`/etc/bashrc`):

		if [ -f /etc/bashrc ]; then
			. /etc/bashrc
		fi
		
**`/etc/profile`和`~/.bashrc`**

`/etc/profile`中设定的变量(全局)的可以作用于任何用户,而`~/.bashrc`等中设定的变量(局部)只能继承`/etc/profile`中的变量,他们是"`父子`"关系.

**`~/.bash_profile`和`~/.bashrc`**

`~/.bash_profile`是`交互式login`方式进入 bash 运行的`~/.bashrc`,`~/.bashrc`是交互式 `non-login 方式`进入 bash 运行的通常二者设置大致相同,所以通常前者会调用后者.