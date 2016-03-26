#Linux 修改主机名

---
		
修改前

		[chloroplast@iZ94ebqp9jtZ

####Centos 6.5
		
修改后		
		
		[chloroplast@dev-server ~]$ cat /etc/sysconfig/network
		NETWORKING=yes
		HOSTNAME=dev-server
		NETWORKING_IPV6=no
		PEERNTP=no
		GATEWAY=120.25.163.247
		
修改`HOSTNAME`


####Centos 7

**修改文件**

		[chloroplast@iZ944l0t308Z ~]$ cat /etc/hostname
		dev-server-2
		
**使用`hostnamectl`**

		[chloroplast@iZ944l0t308Z ~]$ hostnamectl status
		   Static hostname: dev-server-2
		Transient hostname: iZ944l0t308Z
		         Icon name: computer-vm
		           Chassis: vm
		        Machine ID: 45461f76679f48ee96e95da6cc798cc8
		           Boot ID: 3963391650cb4e99a68742c3a10c5f72
		    Virtualization: xen
		  Operating System: CentOS Linux 7 (Core)
		       CPE OS Name: cpe:/o:centos:centos:7
		            Kernel: Linux 3.10.0-327.4.4.el7.x86_64
		      Architecture: x86-64
		      
`修改命令`:

		hostnamectl set-hostname name
		
在通过hostname或者hostnamectl status命令查看更改是否生效