#11_04 Linux网络配置之四 ifconfig及ip命令详解

###笔记

---

####主机接入网络

* IP
* NETMASK(掩码)
* GATEWAY(网关)
* HOSTNAME(主机名),标识主机
* DNS:帮助主机解析其他主机名,额外的DNS地址备用
* 路由
* DHCP: (Dynamic Host Condiguration Protocol)动态主机配置协议

**DNS备用地址**

DNS全球

DNS1,DNS2,DNS3

DNS1连不上,找不见服务器则使用备用的.

`不是`第一个解析不到,使用备用的(因为是全球的,第一个解析不到,第二个肯定也解析不到)

**DHCP**

让一台DHCP服务器为其他主机提供地址,或者其他网络配置信息的一种服务.

####Linux网络

网络是属于`内核`的功能. 

给Linux主机网卡配置地址,`地址属于内核`,`不属于`网卡. 看上去是配置在网卡上的.

无论地址配置哪个网卡,当前主机都会相应.

**网络接口**

每一个网络接口都有`名称`

* `lo`: 本地回环 (本地又做服务器端,又做客户端),数据报文在内存中,不到网络
* `以太网网卡`: eth(0,1,2..) 
* `ppp(0,1,2...)`: 点对点链接

`名称的定义`:

RHEL5: `/etc/modprobe.conf`,通过`alias`定义

RHEL6: `/etc/udev/rules.d/70-persistent-net.rules` 

使用网卡直接通过`名称`来使用.

####ifconfig

网络配置命令

显示当前主机`活动`状态的网卡信息.

`示例`:

		[root@dev-server /]# ifconfig
		...

		eth0      Link encap:Ethernet  HWaddr 00:16:3E:00:41:F2
		          inet addr:10.116.138.44  Bcast:10.116.143.255  Mask:255.255.248.0
		          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
		          RX packets:170443 errors:0 dropped:0 overruns:0 frame:0
		          TX packets:1336 errors:0 dropped:0 overruns:0 carrier:0
		          collisions:0 txqueuelen:1000
		          RX bytes:7798305 (7.4 MiB)  TX bytes:114577 (111.8 KiB)
		          Interrupt:74
		 ...
		 
* `Link encap`: 二次网络使用的协议,`Ethernet`以太网
* `HWaddr`: mac地址,物理地址
* `inet addr`: ip地址
* `Bcast`: 广播地址
* `MASK`: 子网掩码
* `UP`: 启用状态
* `BROADCAST`: 允许广播
* `RUNNING`: 运行
* `MULTICAST`: 多播
* `MTU`: 最大传输单元
* `Metric`: 度量
* `RX`: 接受到的报文个数
* `errors`: 错误的个数
* `dropped`: 丢弃的个数
* `overruns`: 溢出的个数
* `frame`: 帧数
* `TX`: 传出去数据包的个数
* `collisions`: 有多少次冲突发生
* `txqueuelen`: 传输队列的长度
* `RX bytes`: 收到的`字节`数 
* `TX bytes`: 传出去的`字节`数
* `Interrupt`: 中断号

**参数**

ifconfig [网络接口的名称...] (跟上参数显示指定网络接口的信息):

* `-a`: 显示所有接口的配置信息

**配置地址**

`ifconfig ethX IP/MASK [up|down]`

* `up`: 启用
* `down`: 禁用(不需要指定ip地址)

		ifconfig eth1 down
		
配置`立即生效`,但重启网络服务或主机,都会失效.

**网络服务**
 
RHEL5: 
 
`/etc/init.d/network`

* `start`: 读取网卡配置文件让网卡生效	
* `stop`
* `status`
* `restart`

RHEL6: 

`/etc/init.d/NetworkManager`

####添加路由

**route(路由)**

`route` 管理路由

* `add`: 添加路由
	* `-host`: 主机路由
	* `-net:` 网络路由
		* `-net 0.0.0.0`: 默认路由 
* `del`: 删除路由
* `-n`: 以数字方式显示各主机或端口等相关信息 


`route add -net|-host DEST gw NEXTHOP`

添加`默认路由`: `route add default gw NEXTHOP`  
删除`默认路由`: `route del default`

做出的改动重启网络服务或主机后失效.

默认route命令查看本地路由表

		[root@dev-server /]# route
		Kernel IP routing table
		Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
		default         120.25.163.247  0.0.0.0         UG    0      0        0 eth1
		10.0.0.0        10.116.143.247  255.0.0.0       UG    0      0        0 eth0
		10.116.136.0    *               255.255.248.0   U     0      0        0 eth0
		100.64.0.0      10.116.143.247  255.192.0.0     UG    0      0        0 eth0
		120.25.160.0    *               255.255.252.0   U     0      0        0 eth1
		link-local      *               255.255.0.0     U     1002   0        0 eth0
		link-local      *               255.255.0.0     U     1003   0        0 eth1
		172.16.0.0      10.116.143.247  255.240.0.0     UG    0      0        0 eth0
		192.168.42.0    *               255.255.255.0   U     0      0        0 docker0


`Flags`:

* `U`: 启用状态
* `G`: 网关路由 (有下一跳)		
 		
`示例`
		
		添加一条路由 通过  192.168.10.1 到 10.0.0.0
		
		route add -net 10.0.0.0/8 gw 192.168.10.1 
 		
 		gw: 表示网关(gate way)	
 		
 		route del -net 10.0.0.0/8
 		
**配置文件**

网络配置文件: `/etc/sysconfig/network`

网络接口配置文件: `/etc/sysconfig/network-scripts/ifcfg-接口名称(eth0,...)`

`示例`:

		[root@dev-server /]# cat /etc/sysconfig/network-scripts/ifcfg-eth0
		DEVICE=eth0
		ONBOOT=yes
		BOOTPROTO=static
		IPADDR=10.116.138.44
		NETMASK=255.255.248.0
		
* `DEVICE`: 关联的设备,要于文件名的后半部一致(`接口名称`)
* `ONBOOT`: 是否开机时自动激活此网络接口
* `BOOTPROTO {static|none|dhcp|bootp}`: 引导协议,要使用静态地址,使用`static`或`none`; `dhcp`表示使用`DHCP`服务器获取地址.
* `IPADDR`: ip地址
* `NETMASK`: 子网掩码
* `GATEWAY`: 设定默认网关
* `HWADDR`: 硬件地址,要与硬件中的地址保持一致,可省
* `USERCTL {yes|no}`: 是否允许普通用户控制此接口
* `PEERNDS {yes|no}`: 是否在`BOOTPROTP`为`dhcp`时接受由DHCP服务器指定的DNS地址

不会立即生效,但重启网络服务或主机都会生效.
 				 				
**路由**

`/etc/sysconfig/network-scripts/route-ethX` 				 				
`添加格式一`:

`DEST`	`via`	`NEXTHOP` 	

	192.168.10.0/24	via	10.10.10.254			 	
`添加格式二`:

`ADDRESS0=`  
`NETMASK0=`  
`GATEWAY0=`  
`...`
		
		ADDRESS0=192.168.10.0
		NETMASK0=255.255.255.0
		GATEWAY0=10.10.10.254
		

两种格式不能混合使用
 		
**DNS服务器指定方法**

`/etc/resolv.conf`

nameserver DNS_IP_1
nameserver DNS_IP_2
nameserver DNS_IP_3

最多只能有3个 

**指定本地解析**

`/etc/hosts` 本机解析不适用dns

`主机IP`		`主机名`		`主机别名`

**配置主机名**

`hostname` HOSTNAME

立即生效,但不是永久有效

配置文件 `/etc/sysconfig/network` 

* `NETWORKING`: 是否启用本机网络功能
* `NETWORKING_IPV6`: 是否启用ipv6
* `HOSTNAME`: 主机名

**setup**

`setup`

`system-config-network-tui`

文字图形化界面配置
 				 	
**ip**

ip:

* `link`: 配置网络接口属性
	* `show`
		* `ip -s link show`: 显示统计信息
	* `set`: 设定网络接口属性 (`man ip`查看详细)
		* `ip link set DEV {up|down}`
* `addr`: 协议地址
	* `add`
		* `ip addr add ADDRESS dev DEV`
		* `label` 别名
	* `del`
		* `ip addr del ADDRESS dev DEV`
	* `show`
	* `flush` 清空
		* `ip addr flush dev DEV to PREFIX`
* `route`: 路由
 				 	
`示例:禁用eth1`

		ip link set eh1 down

`示例:启用eth1`

		ip link set eh1 up 	
		
`示例:ip addr add`
		
		ip addr add 10.2.2.2/8 dev eth1 label eth1:0
		
		ifconfig 看不见第2个地址(辅助地址),用ip addr show 可以看见
		
`示例: ip addr show`						
		
		[chloroplast@dev-server ~]$ ip addr show
		1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN
		    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
		    inet 127.0.0.1/8 scope host lo
		       valid_lft forever preferred_lft forever
		2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP qlen 1000
		    link/ether 00:16:3e:00:41:f2 brd ff:ff:ff:ff:ff:ff
		    inet 10.116.138.44/21 brd 10.116.143.255 scope global eth0
		       valid_lft forever preferred_lft forever
		 ...
		 
`127.0.0.1/8` 后面的 `/8` 代表子网掩码.8带表掩码是255.0.0.0 的掩码.24代表掩码是255.255.255.0.掩码位数.因为掩码以2进制表示就是:11111111.00000000.00000000.00000000 用十进制就表示是255.0.0.0 .那个8,24,就代表的1有几个. 
		
**网卡多地址**		
		
一个网卡可以使用多个地址:

网络设备可以别名, ehtX:X (eh0:0,eth0:1...)

配置方法: `ifconfig ethX:X IP/NETMASK`

		ifconfig eth0:0 172.16.200.33/16
		
配置文件: `/etc/sysconfig/network-scripts/ifcfg-ethX:X`
		
非主要地址(别名)不能使用`DHCP`动态获取.

				 				
###整理知识点

---

**网关**

网关(Gateway)又称网间连接器,协议转换器.网关在网络层以上实现网络互连,网关是一种充当转换重任的计算机系统或设备.使用在不同的通信协议,数据格式或语言,甚至体系结构完全不同的两种系统之间,网关是一个翻译器.

**路由**

路由就是指导IP数据包发送的路径信息.路由协议是在路由指导IP数据包发送过程中事先约定好的规定和标准.

**MTU**

最大传输单元(Maximum Transmission Unit,MTU)是指一种`通信协议`的某一层上面所能`通过的最大数据包大小`(以`字节`为单位).

**网络设备别名**

		ip addr add 172.25.215.40/24 dev etho label eth0:0

		后面的eth0:0 表示，我们给eth0这块网卡增加一个IP别名，后面那个0表示别名号，第二个别名就可以写成 eth0:1
		
eth0:0 网卡,其实这个就是eth0网卡的一个IP别名,eth0上就有两个IP地址了,此时我们通过外面的主机ping这两个IP地址都是可以ping通的.
