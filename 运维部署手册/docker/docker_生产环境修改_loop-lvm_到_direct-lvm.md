#Docker 生产环境修改 loop-lvm 到 direct-lvm

---

###检查

确认`docker`的信息.

		docker info
		...
		Data loop file: /data/docker/devicemapper/devicemapper/data
 		WARNING: Usage of loopback devices is strongly discouraged for production use. Use `--storage-opt dm.thinpooldev` to specify a custom block storage device.
 		Metadata loop file: /data/docker/devicemapper/devicemapper/metadata
 		...
 		
 我们从:
 
 * `Data loop file`
 * `Metadata loop file`
 
存放在`/var/lib/docker/devicemapper/devicemapper`文件下. 他们被回环挂在到空闲的文件上.

###配置		

####分区

		[root@rancher-agent-1 ~]# fdisk -l
		...
		磁盘 /dev/xvdc：21.5 GB, 21474836480 字节，41943040 个扇区
		Units = 扇区 of 1 * 512 = 512 bytes
		扇区大小(逻辑/物理)：512 字节 / 512 字节
		I/O 大小(最小/最佳)：512 字节 / 512 字节
		...
		
**对硬盘分区**:

		[root@rancher-agent-1 ~]# fdisk /dev/xvdc
		欢迎使用 fdisk (util-linux 2.23.2)。
		
		更改将停留在内存中，直到您决定将更改写入磁盘。
		使用写入命令前请三思。
		
		Device does not contain a recognized partition table
		使用磁盘标识符 0xe64ea89e 创建新的 DOS 磁盘标签。
		
		命令(输入 m 获取帮助)：n
		Partition type:
		   p   primary (0 primary, 0 extended, 4 free)
		   e   extended
		Select (default p): p
		分区号 (1-4，默认 1)：1
		起始 扇区 (2048-41943039，默认为 2048)：
		将使用默认值 2048
		Last 扇区, +扇区 or +size{K,M,G} (2048-41943039，默认为 41943039)：
		将使用默认值 41943039
		分区 1 已设置为 Linux 类型，大小设为 20 GiB
		
		命令(输入 m 获取帮助)：t
		已选择分区 1
		Hex 代码(输入 L 列出所有代码)：8e
		已将分区“Linux”的类型更改为“Linux LVM”
		
		命令(输入 m 获取帮助)：wq
		The partition table has been altered!
		
		Calling ioctl() to re-read partition table.
		正在同步磁盘。
		
**检查分区结果**:

		[root@rancher-agent-1 ~]# fdisk -l
		...
		磁盘 /dev/xvdc：21.5 GB, 21474836480 字节，41943040 个扇区
		Units = 扇区 of 1 * 512 = 512 bytes
		扇区大小(逻辑/物理)：512 字节 / 512 字节
		I/O 大小(最小/最佳)：512 字节 / 512 字节
		磁盘标签类型：dos
		磁盘标识符：0xe64ea89e
		
		    设备 Boot      Start         End      Blocks   Id  System
		/dev/xvdc1            2048    41943039    20970496   8e  Linux LVM
		
		磁盘 /dev/mapper/docker-202:17-786436-pool：107.4 GB, 107374182400 字节，209715200 个扇区
		Units = 扇区 of 1 * 512 = 512 bytes
		扇区大小(逻辑/物理)：512 字节 / 512 字节
		I/O 大小(最小/最佳)：65536 字节 / 65536 字节
		...
		
####创建pv

		[root@rancher-agent-1 ~]# pvcreate /dev/xvdc1
  		Physical volume "/dev/xvdc1" successfully created
  		
####创建vg

		[root@rancher-agent-1 ~]# vgcreate docker /dev/xvdc1
  		Volume group "docker" successfully created
  		
####创建thinpool
 
 * `-n` name
 * `-W|--wipesignatures {y|n}`
 
 
在名为`docker`的`卷组`上创建一个逻辑卷`thinpool`,占据整个`VG`的`%95`;
 
 		[root@rancher-agent-1 ~]# lvcreate --wipesignatures y -n thinpool docker -l 95%VG
		  Logical volume "thinpool" created.
		  
在名为`docker`的`卷组`上创建一个逻辑卷`thinpoolmeta`,占据整个`VG`的`%95`;


		[root@rancher-agent-1 ~]# lvcreate --wipesignatures y -n thinpoolmeta docker -l 1%VG
		  Logical volume "thinpoolmeta" created.
				
		
####Convert the pool to a thin pool 转换为一个精简池

`lvconvert`: 改变逻辑卷的布局. 精简池元数据存储在`docker/thinpoolmeta`中.
把数据卷,元数据合并成一个精简池,且此精简池使用原数据卷的名字


		[root@rancher-agent-1 ~]# lvconvert -y --zero n -c 512K --thinpool docker/thinpool --poolmetadata docker/thinpoolmeta
		  WARNING: Converting logical volume docker/thinpool and docker/thinpoolmeta to pool's data and metadata volumes.
		  THIS WILL DESTROY CONTENT OF LOGICAL VOLUME (filesystem etc.)
		  Converted docker/thinpool to thin pool.		
####`lvm` profile

一个LVM守护进程(dmeventd)将默认监视thin-pool的LV的数据使用量,达到指定阀值后会自动扩展.
当然自动扩展的前提是你创建thin-pool的那个卷组中依然有空闲的空间.

* `thin_pool_autoextend_threshold`: 这个参数设置了当达到磁盘使用的多少阀值后,自动扩展.如果把它设置成 100 就是关闭自动扩展.最小的值是50.
* `thin_pool_autoextend_percent `: 定义了每次达到上面定义的阀值后自动扩展多大的空间,在其当前的规模上扩展百分之几.


		vi /etc/lvm/profile/docker-thinpool.profile
		
		activation {
		    thin_pool_autoextend_threshold=80
		    thin_pool_autoextend_percent=20
		}
		
####Apply your new lvm profile

`man 5 lvm.conf`可以查看更详细的`metadataprofile`信息.

		[root@rancher-agent-1 ~]# lvchange --metadataprofile docker-thinpool docker/thinpool
  		Logical volume "thinpool" changed.
  		
#####Verify the lv is monitored

检查是否被监控

`lvs`: 报告逻辑卷的信息.`-o`: 选项.`+`代表附加选项而不是替换原来的.

		[root@rancher-agent-1 ~]# lvs -o+seg_monitor
		  LV       VG     Attr       LSize  Pool Origin Data%  Meta%  Move Log Cpy%Sync Convert Monitor
		  thinpool docker twi-a-t--- 19.00g             0.00   0.03                             monitored
		  
####如果docker之前启动过,清除graph driver directory

因为我原来的数据已经配置在`--graph=/data/docker`中,所以我对应的需要删除该目录下的文件.

		rm -rf /var/lib/docker/*

####Configure the Docker daemon with specific devicemapper options

修改我们的配置文件`/etc/systemd/system/docker.service.d/docker.conf`.

		--storage-driver=devicemapper --storage-opt=dm.thinpooldev=/dev/mapper/docker-thinpool --storage-opt dm.use_deferred_removal=true

**解释**

` dm.thinpooldev`:

If using a block device for device mapper storage, it is best to use lvm to create and manage the thin-pool volume. This volume is then handed to Docker to exclusively create snapshot volumes needed for images and containers.

As a fallback if `no thin pool is provided`, `loopback files are created`. Loopback is very slow, but can be used without any pre-configuration of storage. It is `strongly recommended that you do not use loopback in production`. Ensure your Engine daemon has a --storage-opt dm.thinpooldev argument provided.

`dm.use_deferred_removal` (延迟移除):

Deferred device removal means that if device is busy when devices are being removed/deactivated, then a deferred removal is scheduled on device. And devices automatically go away when last user of the device exits.

For example, when a container exits, its associated thin device is removed. If that device has leaked into some other mount namespace and can't be removed, the container exit still succeeds and this option causes the system to schedule the device for deferred removal. It does not wait in a loop trying to remove a busy device.

####重新加载(我们使用的是drop-in file)

		systemctl daemon-reload