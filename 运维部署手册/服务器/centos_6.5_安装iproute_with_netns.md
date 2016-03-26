#CentOS 6.5 安装iproute with netns

---

**检查本机centos版本**

		[chloroplast@dev-server ~]$ lsb_release -a
		LSB Version:	:base-4.0-amd64:base-4.0-noarch:core-4.0-amd64:core-4.0-noarch
		Distributor ID:	CentOS
		Description:	CentOS release 6.5 (Final)
		Release:	6.5
		Codename:	Final
		
**检查本机`iproute`版本**

		[chloroplast@dev-server ~]$ rpm -qa | grep iproute
		iproute-2.6.32-45.el6.x86_64
		
**下载rpm包**

		wget https://repos.fedorapeople.org/openstack/EOL/openstack-grizzly/epel-6/iproute-2.6.32-130.el6ost.netns.2.x86_64.rpm
		
		在服务器chloroplast家目录下
		
**安装**

		rpm -ivh --replacefiles ./iproute-2.6.32-130.el6ost.netns.2.x86_64.rpm 
		