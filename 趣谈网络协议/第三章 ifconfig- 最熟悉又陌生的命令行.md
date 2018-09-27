# 第三章 ifconfig: 最熟悉又陌生的命令行

* `ifconfig`基于`net-tools`工具
* `ip addr`基于`iproute2`工具

## 无类型域间选路

`10.100.122.2./24`这种地址表现形式就是`CIDR`. `24`的意思是, `32`位中, 前`24`位是网络号, 后`8`位是主机号.

CIDR又可以将网络划分为更小的单元.

24表示子网掩码中有多少位为1.

### 示例

`16.158.165.91/22`这个地址的`CIDR`.

`22`不是`8`的整数倍. 即前面`22`位是网络号. 除去`16.158`(16位), 剩下`16.158.<101001><01>.91`
在`10100101`中前`6`位是网络号. 后面是机器号.

所以:

* 第一个地址是: `16.158.<101001><00>.1`即`16.158.164.1`
* 子网掩码是: `255.255.<111111><00>.0`即`255.255.252.0`(22个1, 后面是0)
* 广播地址为:`16.158.<101001><11>.255`即`16.158.167.255`

### 广播地址

主机号全为"1"的网络地址用于广播之用，叫做广播地址。所谓广播，指同时向同一子网所有主机发送报文.

### 网络地址

主机号全为"0".

网络标识相同的计算机必须同属于同一链路.

#### 子网掩码

用来指明一个IP地址的哪些位标识的是主机所在的子网, 以及哪些位标识的是主机的位掩码. 子网掩码只有一个作用, 就是将某个IP地址划分为网络地址和主机地址两部分.

子网掩码的作用就是和IP地址**与(AND)**运算后得出网络地址, 子网掩码也是32bit, 并且是一串1后跟随一串0组成, **其中1表示在IP地址中的网络号对应的位数, 而0表示在IP地址中主机对应的位数**.

## 网络设备的状态标识

```
[ansible@iZ94ebqp9jtZ ~]$ ip addr
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP qlen 1000
    link/ether 00:16:3e:08:74:dc brd ff:ff:ff:ff:ff:ff
    inet 172.18.140.140/20 brd 172.18.143.255 scope global eth0
       valid_lft forever preferred_lft forever
3: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN
    link/ether 02:42:dc:02:b2:c2 brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.1/16 scope global docker0
       valid_lft forever preferred_lft forever
```

* `scope`
	* `global`: 这个网卡是对外的, 可以接受来自各个地方的包.
	* `site`: 仅支持`IPv6`, 仅允许本主机的连线.
	* `link`: 仅允许本装置自我连线.
	* `host`: 这个网卡仅仅可以供本机相互通信.
* `qlen`: 网卡队列长度`ip link set eth0 txqueuelen xxxx`
* `<BROADCAST,MULTICAST,UP,LOWER_UP`这些叫做`net_devices flags`, 网络设备状态标识.
	* `LOOPBACK`: 环回接口, 往往会被分配到`127.0.0.1`这个地址. 这个地址用于本机通信, 经过内核处理后直接返回, 不会在任何网络中出现.
	* `UP`表示网卡处于启动状态.
	* `BROADCAST`表示这个网卡有广播地址, 可以发送广播包.
	* `MULTICAST`表示王珂可以发送多播包.
	* `LOWER_UP`表示
* `LOWR UP`表示**网线插着**
* `mtu 1500`表示最大传输单元`MTU为1500`
	* `MTU`是二层`MAC`层的概念. `MAC`层有`MAC`的头. 以太网规定连`MAC`头带正文合起来, 不允许超过`1500`个字节. 正文里面有`IP`的头, `TCP`的头, `HTTP`的头. 如果放不下, 就需要分片传输.
* `qdisc`: `queueing discipline`排队规则.
	* `pfifo`, 不对进入的数据包做任何的处理, 数据包采用**先入先出**的方式通过队列.
	* `pfifo_fast`包括三个波段(`brand`), 数据包是按照服务类型`Type of Service, TOS`(`IP`头里面的一个字段 )被分配到三个波段里面.
		* `brand 0`优先级最高, 如果`brand 0`里面有数据包, 系统不会处理其他波段的.
		* `brand 1`次之
		* `brand 2`最低

## 总结

* `IP`是**地址**, 有定位功能.
* `MAC`是**身份证**, 无定位功能. 通信范围较小, 局限在一个子网里面. 一旦跨子网, `MAC`地址就不行了, 需要`IP`地址起作用了.
* `CIDR`可以用来判断**是不是本地人**

可以这么理解, 找住在xxx路xx小区xx层身份证为AAA的张三. 你直接在大马路上喊谁是`AAA`没人会应答, 所以要根据`IP`地址(xxx路xx小区xxx层)先定位, 然后到了这一层(一个子网里面)在喊(ARP)谁是`AAA`

### 其他

#### 混杂模式

混杂模式(Promiscuous Mode)是指一台机器能够接收所有经过它的数据流, 而不论其目的地址是否是他. 是相对于通常模式(又称"非混杂模式")而言的.

* ip link set eth0 promisc on   # 开启网卡的混合模式
* ip link set eth0 promisc offi # 关闭网卡的混个模式