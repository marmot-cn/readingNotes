#Docker Centos7 使用direct-lvm代替loopback

示例配置文件位置 `/usr/lib/docker-storage-setup/docker-storage-setup`

新添加硬盘`/dev/xvdc`,最早硬盘大小为`5G`,报错.后来扩容为`10G`一切正常.

**停止docker服务**

		systemctl stop docker # 停止当前运行的 docker
		
**创建配置文件**

`/dev/xvdc` 是我的新挂载的硬盘

		vi /etc/sysconfig/docker-storage-setup 
		
		DEVS="/dev/xvdc"
		VG=docker-vg
		
**删除旧的镜像位置**

		rm -rf /data/docker
		
我本机镜像位置在 /data/docker
		
**运行脚本**

		docker-storage-setup
		
**运行`lvs`**

![lvs](./img/01.png "lvs")

**启动docker**

		systemctl start docker
		
**查看`docker info`**

		[root@dev-server-2 chloroplast]# docker info
		Containers: 0
		Images: 0
		Storage Driver: devicemapper
 		Pool Name: docker--vg-docker--pool  # 此处已经变为相关的设备文件
 		...



