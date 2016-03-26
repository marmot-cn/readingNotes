#Open vSwitch组建安装

---

以下安装在本机阿里云服务器安装

下载源文件:

		wget http://openvswitch.org/releases/openvswitch-2.4.0.tar.gz
		
解压:

		tar -xzvf openvswitch-2.4.0.tar.gz

####Centos 6.5
		
编译:

		cd openvswitch-2.4.0
		./configure --with-linux=/lib/modules/`uname -r`/build
		make
		sudo make install
		
####Centos 7

如果直接编译会报错

**安装依赖包**

		yum -y install openssl-devel wget kernel-devel
		
**添加用户,跳至用户文件夹**

		adduser ovswitch
		su - ovswitch
		cd ~
		
**下载源码**

		wget http://openvswitch.org/releases/openvswitch-2.4.0.tar.gz
		
**解压**
		
		tar -xzvf openvswitch-2.4.0.tar.gz
		
**创建编译目录**

		mkdir -p ~/rpmbuild/SOURCES
		
**从`spec`文件中删除`openvswitch-kmod`的依赖包,并创建一个新的`spec`文件**

		sed 's/openvswitch-kmod, //g' openvswitch-2.4.0/rhel/openvswitch.spec > openvswitch-2.4.0/rhel/openvswitch_no_kmod.spec
		
**开始编译**

		cp openvswitch-2.4.0.tar.gz rpmbuild/SOURCES

		rpmbuild -bb --without check ~/openvswitch-2.4.0/rhel/openvswitch_no_kmod.spec
		
**退出当前用户使用`root`**

**安装编译生成的`rpm`文件**

		yum localinstall /home/ovswitch/rpmbuild/RPMS/x86_64/openvswitch-2.4.0-1.x86_64.rpm
		
**启动服务**

		systemctl start openvswitch.service
		
**开机启动**

		systemctl enable openvswitch.service		
		
**查看服务状态**

		[root@dev-server-2 openvswitch-2.4.0]# systemctl -l status openvswitch.service
		● openvswitch.service - LSB: Open vSwitch switch
		   Loaded: loaded (/etc/rc.d/init.d/openvswitch)
		   Active: active (running) since 六 2016-01-09 17:01:55 CST; 5min ago
		   ...