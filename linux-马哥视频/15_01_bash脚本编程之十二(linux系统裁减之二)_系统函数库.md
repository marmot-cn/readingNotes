#15_01_bash脚本编程之十二(Linux系统裁减之二) 系统函数库

###笔记

---

**裁剪系统**

1. 关机和重启
2. 主机名
3. 运行对应级别的服务脚本
4. 启动终端
5. 运行用户
6. 定义单用户级别
7. 装载网卡驱动,启用网络功能
8. 提供一个web服务器

`busybox`: 命令,可以模拟数百个常用命令的使用.

####示例

**添加硬盘**

添加2块磁盘`hda1`挂载到`/mnt/boot` `hda2`挂载到`/mnt/sysroot`.

**复制内核**

		cp /boot/vmlinuz-xxxx	/mnt/boot/vmlinuz
		
**制作initrd**
		
		#建立test目录
		mkdir test
	    cd test
	    zcat /boot/initrd-xxx.img | cpio -id
	    vi init
	    
	    #修改mkrootdev一条
	    ...
	    mkrootdev -t ext3 -o default,ro /dev/hda2
	    ...
	    
	    #创建归档文件并输出
	    find . | cpio -H newc --quiet -o | gzip -9 > /mnt/boot/initrd.gz
	    
**安装和配置grub**

	    #安装grub
	    grub-install --root-directory=/mnt /dev/hda
	    ...
	    
	    #检查/mnt/boot目录下是否有grub
	    ls /mnt/boot
	    grub initrd.gz lost+found vmlinuz
	    
	    #提供grub的配置文件
	    vim /mnt/boot/grub/grub.conf
	    #编写内容 -- 开始
	    
	    default=0
	    timeout=3
	    title test Linux(xx.xx)
	    		root(hd0,0)
	    		kernel /vmlinuz
	    		initrd /initrd.gz
	    
	    #编写内容 -- 结束
	     
**提供根文件系统所需要的目录结构**

	    #提供根文件系统
	    cd /mnt/sysroot
	    ls
	    lost+found
	    
	    #创建所需要的目录结构
	    mkdir etc/{rc.d,init.d} bin sbin proc sys dev lib root mnt media var/{log,run,lock/subsys,tmp} usr/{bin,sbin,local} tmp home opt boot -pv 
	    
**配置文件:`/etc/inittab`**

	    #创建配置文件
	    vim /etc/inittab
	    #编写内容 -- 开始
	    
	    id:3:initdefault:
	    si::sysinit:/etc/rc.d/rc.sysinit
	    
	    #编写内容 -- 结束

**`/etc/rc.d/rc.sysinit`**
	    
	    vim /etc/rc.d/rc.sysinit
	    #编写内容 -- 开始
	    
	    #!/bin/bash
	    #
	    echo -e "\tWelcome to my linux"
	    /bin/bash
	    
	    #编写内容 -- 结束
	    
	    #给执行权限
	    chmod +x etc/rc.d/rc.sysinit
	    
		cd
		#用上个章节写的脚本复制命令道/mnt/sysroot下面
		./bincp.sh 复制 init,bash,ls,touch,mkdir,rm,mv,cp,cat,mount,umount,vi,vim,chmod,chown,ping,ifconfig,insmod,modprobe,rmmod,route,halt,reboot,shutdown,hostname,sync,sleep,basename,mingetty,basename
		
		#同步磁盘
		sync
		
		#测试
		chroot /mnt/sysroot/
		bash-3.2# ...
		
**重新挂载根文件系统为读写**

		#因为现在的"根"文件系统是只读的,所以重新挂载根文件系统为可读写
		mount -n -o remount,rw /

**编写关机脚本**		

		#编写一个关机脚本
		cd /mnt/sysroot/
		vim etc/rc.d/rc.sysdone
		#编写内容 -- 开始
		
		#!/bin/bash
		#
		
		#确保所有文件都同步完成后关机
		sync
		sleep 2
		syc
		
		#因为执行halt命令是在bash下面的子进程,但是我们关机是要关闭其父进程,所以用exec替换进程(具体可查询exec命令)
		exec /sbin/halt -p
		
		#编写内容 -- 结束
		
		chmod +x etc/rc.d/rc.sysdone
		
		#再次编辑/etc/inittab
		#编写内容 -- 开始
	    
	    #新添加这一条,0级别对应的是关闭
	    l0:0:wait:/etc/rc.d/rc.sysdone
		
		#编写内容 -- 结束

**编写重启脚本**

		#编写一个重启脚本
		cd /mnt/sysroot/
		vim etc/rc.d/rc.reboot	
		
		#!/bin/bash
		#
		sync
		sleep 1
		sync
		
		exec /sbin/reboot
		
		#给执行权限
		chmod +x etc/rc.d/rc.reboot
		
		#再次编辑/etc/inittab
		
		#新添加这一条,6级别对应的是重启
	    l6:6:wait:/etc/rc.d/rc.reboot
		
**关机和重启优化为一个脚本**
		
		vim etc/rc.d/init.d/halt
		
		#!/bin/bash
		#
		
		#用$0因为是用脚本名称来区分是重启还是关机 
		case $0 in
		*reboot)
		  COMMAND='/sbin/reboot' ;;
		*halt)
		  COMMAND='/sbin/halt -p' ;;
		*)
		  echo "Only call this script by *reboot OP *halt;"
		  ;;
		esac
		
		case $1 in
		start)
			;;
		stop)
			;;
		*)
			echo "Usage: `basename $0` {start|stop}"
			
		exec $COMMAND
		
		#给予执行权限
		etc/rc.d/init.d/halt
		
		#创建连接
		cd etc/rc.d/ 
		mkdir rc0.d rc6.d
		
		cd rc0.d/
		ln -sv ../init.d/halt S99halt
		
		cd ..
		cd rc6.d/
		ln -sv ../init.d/halt S99reboot
		
		#创建rc脚本
		cd ..
		vim rc
		#编写内容 -- 开始
		
		#!/bin/bash
		#
		RUNLEVEL=$1
		
		for I in /etc/rc.d/rc$RUNLEVEL.d/K*; do
		  $I stop
		done
		
		for I in /etc/rc.d/rc$RUNLEVEL.d/S*; do
		  $I start
		done
		
		#编写内容 -- 结束
		chmod +x rc
		
		#再次编辑/etc/inittab
		l0:0:wait:/etc/rc.d/rc 0
		l6:6:wait:/etc/rc.d/rc 6
		
**在级别3下启动服务**

		#再次编辑/etc/inittab
		l3:3:wait:/etc/rc.d/rc 3

		#创建目录
		cd etc/rc.d/
		mkdir rc3.d
		vim init.d/tserver
		#编写内容 -- 开始
		
		#!/bin/bash
		#
		# chkconfig: 35 66 33
		# description: test service script
		prog=`basename $0`
		lockfile=/var/lock/subsys/$prog
		
		start(){
			echo "Starting $prog ..."
			touch $lockfile
		}
		
		stop(){
			echo "Stopping $prog..."
			rm -f $lockfile
		}
		
		status(){
			if [ -f $lockfile]; then
			  echo "Running..."
			else
			  echo "Stopped..."
		}
		
		usage(){
		  echo "Usage: $prog (start|stop|status|restart)"
		}
		
		case $1 in
		start)
		  start ;;
		stop)
		  stop ;;
		restart)
		  stop
		  start
		  ;;
		status)
		  status
		  ;;
		*)
		  usage
		  exit 1
		  ;;
		esac
		#编写内容 -- 结束
		
		#给予执行权限
		chmod +x init.d/tserver
		
		#在3级别启动
		cd rc3.d/
		ln -sv ../init.d/tserver S66tserver
		
		#在0,6级别关闭
		cd ..
		cd rc0.d/
		ln -sv ../init.d/tserver K33tserver
		cd ../rc6.d/
		ln -sv ../init.d/tserver K33tserver
		
**编辑`/etc/inittab`**

		1:2345:respawn:/sbin/mingetty --loginprog=/bin/bash tty1
		2:2345:respawn:/sbin/mingetty --loginprog=/bin/bash tty2
		
**编辑`/etc/rc.d/rc.sysinit`去掉`/bin/bash`**

因为已经在`/etc/inittab`内执行`/bin/bash`
		
		vim /etc/rc.d/rc.sysinit
	    #编写内容 -- 开始
	    
	    #!/bin/bash
	    #
	    echo -e "\tWelcome to my linux"
	    
	    #编写内容 -- 结束
	    
**创建`bash shell`链接**

		cd /mnt/sysroot/bin
		ln -sv bash sh
		
**根文件系统在`sysinit`脚本,读写方式重新挂载**

创建`/etc/fstab`文件在`/mnt/sysroot/`下
		
		cd /mnt/sysroot/
		vim etc/fstab
		#编写内容 -- 开始
				
		/dev/hda2	/	ext3	defaults	0	0
		/dev/hda1	/boot	ext3	defaults	0  0
		proc	/proc	proc	defaults	0	0
		sysfs	/sys	sysfs	defaults	0	0
						
		#编写内容 -- 结束
	    
创建主机名
		
		cd /mnt/sysroot/
		mkdir etc/sysconfig
		vim etc/sysconfig/network
		#编写内容 -- 开始
		
		HOSTNAME=chloroplast
		
	    #编写内容 -- 结束
	    
编辑`etc/rc.d/rc.sysinit`文件
		
		cd /mnt/sysroot/
		vim etc/rc.d/rc.sysinit
		#编写内容 -- 开始
		
		#!/bin/bash
	    #
	    echo -e "\tWelcome to my linux"
	    echo "Remount rootfs..."
	    mount -n -o remount,rw / 
	    
	    echo "Set the hostname..."
	    #如果文件存在用 "." 读取文件内容,如果HOSTNAME有内容,则变量$HOSTNAME有内容
	    [ -f /etc/sysconfig/network ] && . /etc/sysconfig/network
	    [ -z $HOSTNAME -o "$HOSTNAME" == '(none)' ] && HOSTNAME=localhost
	    /bin/hostname $HOSTNAME
	    
	    #编写内容 -- 结束	    
	    
编辑`/mnt/sysroot/etc/inittab`使用`agetty`(38400是速率)	    

	    cd /mnt/sysroot/
		vim etc/inittab
		#编写内容 -- 开始
		    
	    1:2345:respawn:/sbin/agetty -n -l /bin/bash 38400 tty1
		2:2345:respawn:/sbin/agetty -n -l /bin/bash 38400 tty2
		
		#编写内容 -- 结束
	     
**编写公用函数**

		cd /mnt/sysroot/
		vim etc/rc.d/init.d/functions
		#编写内容 -- 开始	    
		#!/bin/bash
		SCREEN=`stty -F /dev/console size`
		COLUMNS=${SCREEN#* }
		[ -z $COLUMNS ] && COLUMNS=80
		
		#12代表输出 [  ok  ] 包含空格的字符长度
		SPA_COL=$[$COLUMNS-14]   
		
		RED='\033[31m'
		GREEN='\033[32m'
		YELLOW='\033[33m'
		BLUE='\033\34m'
		NORMAL='\033[0m'
		
		success() {
		  string=$1
		  RT_SPA=$[$SPA_COL-${#string}]
		  echo -n "$string"
		  for I in `seq 1 $RT_SPA`;do
		    echo -n " "
		  done
		  echo -e "[   ${GREEN}OK${NORMAL}   ]"
		}
		
		failure() {
		  string=$1
		  RT_SPA=$[$SPA_COL-${#string}]
		  echo -n "$string"
		  for I in `seq 1 $RT_SPA`;do
		    echo -n " "
		  done
		  echo -e "[ ${RED}FAILER${NORMAL} ]"
		}
		
		success "starting tserver"
		failure "staring tserver"
		   
		#编写内容 -- 结束 
		
**引用公用函数文件改写`etc/init.d/tserver`**

		#!/bin/bash
		#
		# chkconfig: 35 66 33
		# description: test service script
		prog=`basename $0`
		lockfile=/var/lock/subsys/$prog
		
		#引用公共函数文件
		. /etc/rc.d/init.d/functions
		
		start(){
			touch $lockfile
			[ $? -eq 0] && success "Starting $prog" || failure "Starting $prog"
		}
		
		stop(){
			rm -f $lockfile
			[ $? -eq 0] && success "Stopping $prog" || failure "Stopping $prog"
		}
		
		status(){
			if [ -f $lockfile]; then
			  echo "Running..."
			else
			  echo "Stopped..."
		}
		
		usage(){
		  echo "Usage: $prog (start|stop|status|restart)"
		}
		
		case $1 in
		start)
		  start ;;
		stop)
		  stop ;;
		restart)
		  stop
		  start
		  ;;
		status)
		  status
		  ;;
		*)
		  usage
		  exit 1
		  ;;
		esac

**让自己制作的小系统拥有ip地址**

复制`pcnet32`模块(其依赖于`mii`模块,也要一并复制)
		
		cd /mnt/sysroot/
		mkdir lib/modules
		
		cp xxx/xxx/pcnet32.ko /mnt/sysroot/lib/modules
		cp xxx/xxx/mii.ko /mnt/sysroot/lib/modules
		
编辑`etc/rc.d/rc.sysinit`		
 
 		vim etc/rc.d/rc.sysinit
		#编写内容 -- 开始
		
		#!/bin/bash
	    #
	    echo -e "\tWelcome to my linux"
	    echo "Remount rootfs..."
	    mount -n -o remount,rw / 
	    
	    echo "Set the hostname..."
	    #如果文件存在用 "." 读取文件内容,如果HOSTNAME有内容,则变量$HOSTNAME有内容
	    [ -f /etc/sysconfig/network ] && . /etc/sysconfig/network
	    [ -z $HOSTNAME -o "$HOSTNAME" == '(none)' ] && HOSTNAME=localhost
	    /bin/hostname $HOSTNAME
	    
	    echo "Initializing network device..."
	    /sbin/insmod /lib/modules/mii.ko
	    /sbin/insmod /lib/modules/pcnet32.ko
	    
	    #编写内容 -- 结束	    
	    
编辑`etc/sysconfig/network-scripts/ifcfg-eth0`赋值ip地址

		vim etc/sysconfig/network-scripts/ifcfg-eth0
		#编写内容 -- 开始
		DEVICE=eth0
		BOOTPROTO=static
		IPADDR=172.16.100.5
		NETMASK=255.255.0.0
		GATEWAY=172.16.0.1
		ONBOOT=yes
		
		#编写内容 -- 结束

编辑`etc/rc.d/init.d/network`

		vim etc/rc.d/init.d/network
		#编写内容 -- 开始
	    
	    #!/bin/bash
	    #
	    # chkconfig: 35 09 90
	    # description: network service
	    prog=network
	    
	    . /etc/rc.d/init.d/functions
	    CONF=/etc/sysconfig/network-scripts/ifcfg-eth0
	    
	    . $CONF
	    
	    NETMASK=16
	    
	    #配置ip地址
	    start() {
	      ifconfig etho0 $IPADDR/$NETMASK up
	      [ -z $GATEWAY ] && route add default gw $GATEWAY
	    }
	    
	    stop() {
	      ifconfig eth0 down
	    }
	    
	    status() {
	      ifconfig eth0
	    }
	    
	    usage() {
	      echo "$prog: {start|stop|restart|status}"
	    }
	    
	    case $1 in 
	    start)
	      start
	      success "Config network eth0"
	      ;;
	    stop)
	      stop
	      success "Stop network card eth0"
	      ;;
	    restart)
	      stop
	      start
	      success "Restart network card eth0"
	      ;;
	    status)
	      status
	      ;;
	    *)
	      usage
	      ;;
	    esac	       
	     
	    #编写内容 -- 结束
	     
	    chmod +x etc/rc.d/init.d/network
	    
关闭开启时自动执行

		cd etc/rc.d/rc0.d/
		ln -sv ../init.d/network K90network
		cd ../rc6.d/
		ln -sv ../init.d/network K90network
		cd ../rc3.d/
		ln -sv ../init.d/network S09network
	   
####命令		
		
**`mount -a`**

挂载时不更新`/etc/mtab`文件

`cat /proc/mounts`显示当前系统所挂载的所有文件系统.			    
**`mingetty`**

启动一个终端,并启动一个登录程序

####脚本编程知识点

**变量中字符的长度**

`${#VARNAME}`

###整理知识点

---