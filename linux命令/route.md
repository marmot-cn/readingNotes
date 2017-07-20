# route

---

网关实质上是一个网络通向其他网络的`IP`地址. 即使是两个网络连接在同一台交换机(或集线器)上, `TCP/IP`协议也会根据子网掩码(255.255.255.0)判定两个网络中的主机处在不同的网络里. 而要实现这两个网络之间的通信, 则必须通过网关.

如果网络A中的主机发现数据包的目的地不在本地网络中, 就把数据包转发给自己的网关, 再由网关转发给网络B的网关, 网络B的网关再转发给网络B的某个主机.

网关的IP地址是具有路由功能的设备的IP地址.

#### route 命令

route命令用于显示和操作IP路由表。要实现两个不同的子网之间的通信，需要一台连接两个网络的路由器，或者同时位于两个网络的网关来实现。在Linux系统中，设置路由通常是 为了解决以下问题：该Linux系统在一个局域网中，局域网中有一个网关，能够让机器访问Internet，那么就需要将这台机器的IP地址设置为 Linux机器的默认路由。要注意的是，直接在命令行下执行route命令来添加路由，不会永久保存，当网卡重启或者机器重启之后，该路由就失效了；要想永久保存，有如下方法：

* 在`/etc/rc.local`里添加
* 在`/etc/sysconfig/network`里添加到末尾
* 修改`/etc/sysconfig/static-router`

##### 命令参数

`route [-nee]`

`route add [-net|-host] [网域或主机] netmask [mask] [gw|dev]`

`route del [-net|-host] [网域或主机] netmask [mask] [gw|dev]`

参数:

* `-n`: 不要使用通讯协定或主机名称,直接使用`IP`或`port number`.
* `-ee`：使用更详细的资讯来显示.

增加(`add`)与删除(`del`)路由的相关参数:

* `-net` ：表示后面接的路由为一个网域
* `-host`：表示后面接的为连接到单部主机的路由
* `netmask`：与网域有关，可以设定 `netmask` 决定网域的大小
* `gw` ：gateway 的简写，后续接的是 `IP` 的数值喔，与 `dev` 不同
* `dev`：如果只是要指定由那一块网路卡连线出去，则使用这个设定，后面接 `eth0` 等

`route add default gw {IP-ADDRESS} {INTERFACE-NAME}`用于设置默认路由

* `{IP-ADDRESS}`: 用于指定路由器(网关)的IP地址.
* `{INTERFACE-NAME}`: 用于指定接口名称，如eth0. 使用`/sbin/ifconfig -a`可以显示所有接口信息.

`route add -net {NETWORK-ADDRESS} netmask {NETMASK} dev {INTERFACE-NAME}`添加到指定网络的路由规则.

* `{NETWORK-ADDRESS}`: 用于指定网络地址
* `{NETMASK}`: 用于指定子网掩码
* `{INTERFACE-NAME}`: 用于指定接口名称，如eth0

`route add -net {NETWORK-ADDRESS} netmask {NETMASK} reject`设置到指定网络为不可达，避免在连接到这个网络的地址时程序过长时间的等待，直接就知道该网络不可达.

#### 示例(我自己没有环境测试过)

A电脑：外网(eth0) ip为192.168.1.1，内网(eth1) ip为10.1.1.1

B电脑：内网(eth0) ip为10.1.1.2

##### 配置方法如下:

**启用Linux的ip转发功能**

首先查看ip转发功能是否已经打开

`sysctl -a |grep 'net.ipv4.ip_forward'`

如果值为1，则说明已经打开，否则执行以下操作

`echo "1" >/proc/sys/net/ipv4/ip_forward`

这是临时修改，系统重启后失效

或者直接修改配置文件

`vi /etc/sysctl.conf`

找到

`net.ipv4.ip_forward = 0`

改成

`net.ipv4.ip_forward = 1`

然后运行 `sysctl -p`

 

**设置iptables规则**

在A电脑上执行

`iptables -t nat -A POSTROUTING -s 10.1.1.0/24 -j SNAT --to-source 192.168.1.1`

保存配置

`service iptables save`

这个语句就是告诉系统把即将要流出本机的数据的source ip address修改成为192.168.1.1。这样，数据包在达到目的机器以后，目的机器会将包返回到192.168.1.1。如果不做这个操作，那么数据包在传递的过程中，reply的包肯定会丢失。


假如当前系统用的是ADSL/3G/4G动态拨号方式，那么每次拨号，出口ip都会改变，snat是把源地址转换为固定的ip地址，这样就会有局限性。这时可以使用：


`iptables -t nat -A POSTROUTING -s 10.1.1.0/24 -o eth0 -j MASQUERADE`


与snat不同的是，masquerade可以自动读取外网卡获得动态ip地址，然后进行地址转换。

POSTROUTING：在通过Linux路由器之后做的策略，也就是路由器的外网接口。

-s 10.1.1.0/24：源数据所来自这个网段，也可以是单个ip，不写表示所有内网ip。

-o eth0 -j MASQUERADE：表示在eth0这个外网接口上使用IP伪装。

 

**在B电脑上设置网关**

vi/etc/sysconfig/network-scripts/ifcfg-eth0

添加或修改

GATEWAY=10.1.1.1

保存后重启网络服务

service network restart

 

**测试**

在B电脑上执行

ping www.163.com

[root@www ~]# ping www.163.com
PING 163.xdwscache.ourglb0.com (58.216.21.93) 56(84) bytes of data.
64 bytes from 58.216.21.93: icmp_seq=1 ttl=128 time=8.23 ms
64 bytes from 58.216.21.93: icmp_seq=2 ttl=128 time=9.72 ms
64 bytes from 58.216.21.93: icmp_seq=3 ttl=128 time=8.78 ms
64 bytes from 58.216.21.93: icmp_seq=4 ttl=128 time=8.97 ms


curl www.163.com -I

[root@www ~]# curl www.163.com -I
HTTP/1.0 200 OK
Expires: Tue, 09 May 2017 15:21:02 GMT
Date: Tue, 09 May 2017 15:19:42 GMT
Server: nginx
Content-Type: text/html; charset=GBK
Vary: Accept-Encoding,User-Agent,Accept
Cache-Control: max-age=80
Via: 1.1 cache.163.com:80 (squid)
X-Via: 1.0 czdx87:3 (Cdn Cache Server V2.0), 1.0 adxxz41:6 (Cdn Cache Server V2.0)
Connection: keep-alive


结果显示已经可以正常上网