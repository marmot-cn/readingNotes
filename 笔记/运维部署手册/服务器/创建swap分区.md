#创建swap分区

---

SWAP就是LINUX下的虚拟内存分区,它的作用是在物理内存使用完之后,将磁盘空间(也就是SWAP分区)虚拟成内存来使用

今天在阿里云主机上跑了docker `gitlab` 内存不够, 主要在另外一个容器内编译php5.6,显示无法fork.查询后得知内存不够.所以创建`swap`分区临时使用.暂时也不想花钱加大内存.



####用文件作为Swap分区

**创建/swap路径,放置我们的文件**

		/data/swap
		
**创建用于交换分区的文件**

		sudo dd if=/dev/zero of=/data/swap/swap-file  bs=1k count=1024000
		
		1G
		
这里大小根据自己服务器的内存来选择,我的服务器是1G内存,其实应该选择2Gswap.
		
**建立`swap`文件系统**

		sudo mkswap swap-file
		
**启用交换分区文件**

		sudo swapon swap-file
		
**查看**

		[chloroplast@iZ94ebqp9jtZ ~]$ swapon -s
		Filename				Type		Size	Used	Priority
		/data/swap/swap-file                    file		1023996	528648	-1
		
		[chloroplast@iZ94ebqp9jtZ ~]$ free
		             total       used       free     shared    buffers     cached
		Mem:       1019184     954068      65116          0       4528      70044
		-/+ buffers/cache:     879496     139688
		Swap:      1023996     528648     495348
		
**使系统开机时自启用,在文件`/etc/fstab`中添加一行**

		/data/swap/swap-file swap swap defaults 0 0
		
		或者
		
		echo '/data/swap/swap-file swap swap defaults 0 0' >> /etc/fstab
		
		[chloroplast@iZ94ebqp9jtZ ~]$ sudo cat /etc/fstab

		#
		# /etc/fstab
		# Created by anaconda on Thu Aug 14 21:16:42 2014
		#
		# Accessible filesystems, by reference, are maintained under '/dev/disk'
		# See man pages fstab(5), findfs(8), mount(8) and/or blkid(8) for more info
		#
		UUID=94e4e384-0ace-437f-bc96-057dd64f42ee / ext4 defaults,barrier=0 1 1
		tmpfs                   /dev/shm                tmpfs   defaults        0 0
		devpts                  /dev/pts                devpts  gid=5,mode=620  0 0
		sysfs                   /sys                    sysfs   defaults        0 0
		proc                    /proc                   proc    defaults        0 0
		/dev/volume-group1/data /data ext4 defaults 0 0
		/dev/mapper/volume--group1-docker /docker ext4 defaults 0 0
		/dev/volume-group1/log /log ext4 defaults 0 0
		/data/swap/swap-file swap swap defaults 0 0