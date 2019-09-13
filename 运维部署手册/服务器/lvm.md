# LVM

---

### 简介

* 创建物理卷`pvcreate`
* 创建卷组并给卷组增加分区`vgcreate`
* 创建新的逻辑卷使用lvcreate


![lvm](./img/lvm.jpg "lvm")


### 示例

这里我用的是阿里云的云盘.挂载一块新的云盘后

#### 准备存盘分区
		
		[root@iZ94ebqp9jtZ /]# fdisk -l
		Disk /dev/xvda: 21.5 GB, 21474836480 bytes
		255 heads, 63 sectors/track, 2610 cylinders
		Units = cylinders of 16065 * 512 = 8225280 bytes
		Sector size (logical/physical): 512 bytes / 512 bytes
		I/O size (minimum/optimal): 512 bytes / 512 bytes
		Disk identifier: 0x00078f9c
		
		    Device Boot      Start         End      Blocks   Id  System
		/dev/xvda1   *           1        2611    20970496   83  Linux
		
		Disk /dev/xvdb: 53.7 GB, 53687091200 bytes
		255 heads, 56 sectors/track, 7342 cylinders
		Units = cylinders of 14280 * 512 = 7311360 bytes
		Sector size (logical/physical): 512 bytes / 512 bytes
		I/O size (minimum/optimal): 512 bytes / 512 bytes
		Disk identifier: 0x452bd090
		
		    Device Boot      Start         End      Blocks   Id  System
		/dev/xvdb1               1         734     5240732   83  Linux
		[root@iZ94ebqp9jtZ /]# fdisk /dev/xvdb
		
`fdisk`,因为这个磁盘分区以前使用过(文章`阿里云挂载云盘`),现在重新创建.

		[root@iZ94ebqp9jtZ /]# fdisk /dev/xvdb

		WARNING: DOS-compatible mode is deprecated. It's strongly recommended to
		         switch off the mode (command 'c') and change display units to
		         sectors (command 'u').
		
		Command (m for help): n
		Command action
		   e   extended
		   p   primary partition (1-4)
		p
		Partition number (1-4): 1
		Partition 1 is already defined.  Delete it before re-adding it.
		
		Command (m for help): n
		Command action
		   e   extended
		   p   primary partition (1-4)
		p
		Partition number (1-4): 1
		Partition 1 is already defined.  Delete it before re-adding it.
		
		Command (m for help): d
		Selected partition 1
		
		Command (m for help): n
		Command action
		   e   extended
		   p   primary partition (1-4)
		p
		Partition number (1-4): 1
		First cylinder (1-7342, default 1):
		Using default value 1
		Last cylinder, +cylinders or +size{K,M,G} (1-7342, default 7342):
		Using default value 7342
		
		Command (m for help): t ##改变类型
		Selected partition 1
		Hex code (type L to list codes): 8e ##LVM的分区代码
		Changed system type of partition 1 to 8e (Linux LVM)
		
		Command (m for help): wq
		The partition table has been altered!
		
		Calling ioctl() to re-read partition table.
		Syncing disks.	
		
#### 准备物理卷(PV)

**创建物理卷(`pvcreate`)**

		[root@iZ94ebqp9jtZ /]# pvcreate /dev/xvdb1
  		Physical volume "/dev/xvdb1" successfully created
  		
**检查物理卷的创建情况(`pvdisplay`)**

		[root@iZ94ebqp9jtZ /]# pvdisplay
		 "/dev/xvdb1" is a new physical volume of "49.99 GiB"
		  --- NEW Physical volume ---
		  PV Name               /dev/xvdb1
		  VG Name
		  PV Size               49.99 GiB
		  Allocatable           NO
		  PE Size               0
		  Total PE              0
		  Free PE               0
		  Allocated PE          0
		  PV UUID               NjzZIO-AkLn-nbOD-73OV-oEaR-9mTo-YZAlna

**删除物理卷(`pvremove`)**
 		
`pvremove /dev/xdb1`
 			
#### 准备卷组(VG)

**创建卷组(`vgcreate`)**

		[root@iZ94ebqp9jtZ /]# vgcreate volume-group1 /dev/xvdb1
  		Volume group "volume-group1" successfully created
  		
创建名为`volume-group1`的卷组,使用`/dev/xvdb1`

**查看卷组信息**
	
		[chloroplast@iZ94ebqp9jtZ ~]$ sudo vgs
		[root@iZ94ebqp9jtZ chloroplast]# vgs
  		VG            #PV #LV #SN Attr   VSize  VFree
  		volume-group1   1   2   0 wz--n- 49.99g 7.99g
  		
**删除卷组(`vgremove`)**

`vgremove volume-group1`

#### 创建逻辑卷(LV)

**创建逻辑卷(`lvcreate`)**

		[root@iZ94ebqp9jtZ /]# lvcreate -L 30G -n data volume-group1
  		Logical volume "data" created		
		
创建一个名为`data`,大小为`30G`的逻辑卷.

可以使用`lvdisplay`查看逻辑卷使用情况.

		[root@iZ94ebqp9jtZ /]# lvdisplay
		  --- Logical volume ---
		  LV Path                /dev/volume-group1/data
		  LV Name                data
		  VG Name                volume-group1
		  LV UUID                HH6apj-R8LZ-NpV5-CBJU-tmef-jqBV-QU0TXr
		  LV Write Access        read/write
		  LV Creation host, time iZ94ebqp9jtZ, 2015-12-15 12:36:31 +0800
		  LV Status              available
		  # open                 0
		  LV Size                30.00 GiB
		  Current LE             7680
		  Segments               1
		  Allocation             inherit
		  Read ahead sectors     auto
		  - currently set to     256
		  Block device           253:0
		  
`一些参数`:

`-l`:
	* `%VG`: a percentage of the total space in the Volume Group with the suffix %VG
	* `%FREE`: a percentage of the remaining free space in the Volume Group with the suffix %FREE
	* `%PVS	`: a percentage of the remaining free space for the specified PhysicalVolume(s) with the suffix %PVS		
	* `%ORIGIN`: a percentage of the total space in the Origin  Logical  Volume  with  the  suffix  %ORIGIN  

**格式化和挂载逻辑卷**

		[root@iZ94ebqp9jtZ /]# mkfs.ext4 /dev/volume-group1/data
		mke2fs 1.41.12 (17-May-2010)
		文件系统标签=
		操作系统:Linux
		块大小=4096 (log=2)
		分块大小=4096 (log=2)
		Stride=0 blocks, Stripe width=0 blocks
		1966080 inodes, 7864320 blocks
		
		393216 blocks (5.00%) reserved for the super user
		第一个数据块=0
		Maximum filesystem blocks=4294967296
		240 block groups
		32768 blocks per group, 32768 fragments per group
		8192 inodes per group
		Superblock backups stored on blocks:
			32768, 98304, 163840, 229376, 294912, 819200, 884736, 1605632, 2654208,
			4096000
		
		正在写入inode表: 完成
		Creating journal (32768 blocks): 完成
		Writing superblocks and filesystem accounting information: 完成
		
		This filesystem will be automatically checked every 24 mounts or
		180 days, whichever comes first.  Use tune2fs -c or -i to override.
		
		[root@iZ94ebqp9jtZ /]# mount /dev/volume-group1/data /data
		[root@iZ94ebqp9jtZ /]# df -h
		Filesystem            Size  Used Avail Use% Mounted on
		/dev/xvda1             20G  4.0G   15G  22% /
		tmpfs                 498M     0  498M   0% /dev/shm
		/dev/mapper/volume--group1-data
		                       30G   44M   28G   1% /data
		[root@iZ94ebqp9jtZ home]# echo '/dev/volume-group1/data /data ext4 defaults 0 0' >> /etc/fstab
		
#### 扩容

**从30G扩容为32G**

		lvextend -L32G /dev/volume-group1/data

也可以执行
		
		lvextend -L+2G /dev/lvm_test/test

		
**umount文件系统**

		[root@iZ94ebqp9jtZ data]# umount /dev/volume-group1/data
umount: /data: device is busy.
        (In some cases useful info about processes that use
         the device is found by lsof(8) or fuser(1))
         
 查看到有进程在使用
 
 		[root@iZ94ebqp9jtZ data]# fuser -m /dev/volume-group1/data
		/dev/volume-group1/data:  1172c
		
kill掉正在访问的进程

		[root@iZ94ebqp9jtZ data]# fuser -m -v -i -k /dev/volume-group1/data
		
umount

		 umount /dev/volume-group1/data
		 
调整被加载的文件系统大小(ext2/ext3/ext4)使用`resize2fs`(增大和减小都支持)

		 e2fsck -f /dev/volume-group1/data
		 resize2fs /dev/volume-group1/data
		 
调整被加载的文件系统大小(xfs)使用`xfs_growfs`(只支持增大), **可以不用卸载文件系统, 直接扩展文件系统, xfs_growfs /mountpoint**

		xfs_growfs /dev/volume-group1/data

mount		 

		 mount -a
		 [root@iZ94ebqp9jtZ home]# df -h
		 Filesystem            Size  Used Avail Use% Mounted on
		 /dev/xvda1             20G  4.0G   15G  22% /
		 tmpfs                 498M     0  498M   0% /dev/shm
		 /dev/mapper/volume--group1-data
		                       32G   48M   30G   1% /data
		                       
#### 添加新的物理卷到卷组中

当系统安装了新的磁盘或新建分区并创建了新的物理卷,而要将其添加到已有卷组时,就需要使用vgextend命令:

		[root@iZ94ebqp9jtZ chloroplast]# vgextend volume-group1 /dev/xvdc1
		  Volume group "volume-group1" successfully extended
		[root@iZ94ebqp9jtZ chloroplast]# vgs
		  VG            #PV #LV #SN Attr   VSize  VFree
		  volume-group1   2   2   0 wz--n- 55.99g 13.99g
		  
`/dev/xvdc1`是新的物理卷,制作方式和`/dev/xvdb1`一致.
				                       
