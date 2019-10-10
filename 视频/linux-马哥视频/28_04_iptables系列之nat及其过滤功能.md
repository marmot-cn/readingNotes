# 28_04_iptables系列之nat及其过滤功能

---

### 笔记

---

#### 自定义规则链

添加了一个自定义链`clean_in`.

```shell
[root@iZ944l0t308Z ~]# iptables -N clean_in
[root@iZ944l0t308Z ~]# iptables -L -n
Chain INPUT (policy ACCEPT)
target     prot opt source               destination
ACCEPT     tcp  --  0.0.0.0/0            120.25.87.35         tcp dpt:80 state NEW,ESTABLISHED

Chain FORWARD (policy ACCEPT)
target     prot opt source               destination

Chain OUTPUT (policy ACCEPT)
target     prot opt source               destination
           tcp  --  120.25.87.35         0.0.0.0/0            tcp spt:80 state ESTABLISHED

Chain clean_in (0 references)
target     prot opt source               destination
```

`0 references`: 表示有多少条引用, 被主链所调用.

```shell
添加一条规则
[root@iZ944l0t308Z ~]# iptables -A clean_in -d 255.255.255.255 -p icmp -j DROP
[root@iZ944l0t308Z ~]# iptables -L -n --line-numbers
Chain INPUT (policy ACCEPT)
num  target     prot opt source               destination
1    ACCEPT     tcp  --  0.0.0.0/0            120.25.87.35         tcp dpt:80 state NEW,ESTABLISHED

Chain FORWARD (policy ACCEPT)
num  target     prot opt source               destination

Chain OUTPUT (policy ACCEPT)
num  target     prot opt source               destination
1               tcp  --  120.25.87.35         0.0.0.0/0            tcp spt:80 state ESTABLISHED

Chain clean_in (0 references)
num  target     prot opt source               destination
1    DROP       icmp --  0.0.0.0/0            255.255.255.255

无论到哪个地址, 只要进入input链, 先交由clean_in处理,clean_in如果没有问题在返回主链,交由第2条处理.
iptables -I INPUT -j clean_in

可见clean_in链的references从0变为1
[root@iZ944l0t308Z ~]# iptables -L -n --line-numbers
Chain INPUT (policy ACCEPT)
num  target     prot opt source               destination
1    clean_in   all  --  0.0.0.0/0            0.0.0.0/0
2    ACCEPT     tcp  --  0.0.0.0/0            120.25.87.35         tcp dpt:80 state NEW,ESTABLISHED

Chain FORWARD (policy ACCEPT)
num  target     prot opt source               destination

Chain OUTPUT (policy ACCEPT)
num  target     prot opt source               destination
1               tcp  --  120.25.87.35         0.0.0.0/0            tcp spt:80 state ESTABLISHED

Chain clean_in (1 references)
num  target     prot opt source               destination
1    DROP       icmp --  0.0.0.0/0            255.255.255.255
```

```shell
因为clean_in有引用且不是空链, 所以不能直接删除
[root@iZ944l0t308Z ~]# iptables -X clean_in
iptables: Too many links.
```

#### 利用`iptables`的`recent`模块来抵御`DOS`(拒绝服务攻击)攻击

ssh: 远程连接,

利用`connlimit`模块将单`IP`的并发设置为`3`; 会误杀使用`NAT`上网的用户, 可以根据实际情况增大该值.

```shell
iptables -I INPUT -p tcp --dport 22 -m connlimit --connlimit-above 3 -j DROP
```

`recent`模块, 将最近对服务器某个服务发起请求的链接的`IP`地址记录下来. 利用`recent`和`state`模块限制单`IP`在`300s`内只能与本机建立`3`个新连接. 被限制五分钟后即可恢复访问.

```shell
设置模板
iptables -I INPUT -p tcp --dport 22 -m state --state NEW -m recent --set --name SSH

更新模板 3/300s, 超出3次就DROP, 并且拒绝5分钟(--seconds 300 还有此意思).
iptables -I INPUT -p tcp --dport 22 -m state --state NEW -m recent --update --seconds 300 --hitcount 3 --name SSH -j DROP
```

* 第二句是记录访问`tcp 22`端口的新连接, `--set` 记录数据包的来源`IP`,记录名称为`SSH`.  如果`IP`已经存在将更新已经存在的条目.
* 第三局是指`SSH`记录中的`IP`, `300s`内发起超过`3`次连接则拒绝此`IP`的链接.
	* `--update`是指每次建立链接都更新列表
	* `--seconds`必须与`--rcheck`或者`--updates`同时使用
	* `--hitcount`必须与`--rcheck`或者`--updates`同时使用

`iptables`的记录: `/proc/net/ip_recent/ssh`

也可以使用下面的这句记录日志:

```shell
iptbles -A INPUT -p tcp --dport 22 -m state --state NEW -m recent --update --name SSH --second 300 --hitcount 3 -j LOG --log-prefix "SSH Attack"
```

#### NAT

`N`etwork `A`ddress `T`ranslation. 网络地址转换.

* `SNAT`: 源地址转换(`POSTROUTING`,`OUTPUT`). 转换`IP`报文中的源地址.
* `DNAT`: 目标地址转换. 转换`IP`报文中的目标地址.

对于`Linux`地址属于主机, 所以无论主机有多少块网卡, 都可以`ping`通. 这个==不是转发==, `ping`另外一个网卡链接的主机的时候才是转发.

转发和主机本身没有关系, 不会走`INPUT`和`OUTPUT`, 只会走`FORWARD`.

##### SNAT

**源地址转换**

`/proc/sys/net/ipv4/ip_forward`为`1`代表可以转发. 永久有效修改`/etc/sysctl.conf`中的`net.ipv4.ip_forward=1`即可,然后运行 `sysctl -p`立即生效.


主机A 和 主机C 转发通信.

```
B 扮演路由的角色转发 A -> C

主机A -> 主机B的1号网卡,主机B的2号网卡 -> 主机C

A ping C, A 的信息从1号网卡转发到2号网卡 在到 B. 在 B 来看, 就是由 A ping 过来的.
```

为什么要做源地址转换, 因为私有地址不能用在互联网, 报文可以发出去但是回不来. 

**`nat`表**

```
A -> 路由主机 -> sohu, 访问需要进行源地址转换(转换到路由主机), 否则sohu的响应会发现时私有地址链接不到.

sohu -> 路由主机 -> A, 经过路由主机的网卡需要进行目标地址转换, 转换到A, 因为路由主机并没有访问过sohu.
```

中间的链接关系记录在`nat`表(池).`nat`表内部是靠`序列号`来区分每个主机, 因为每个主机的序列号都不一样.

**`iptables`配置源地址转换**

* `-j SNAT`
	* `--to-source` 转换成哪个源地址, 地址可以不止一个, 多个地址应用场景可以用于`负载均衡`.
			
			转换源地址这4个中的任意一个, 转换给谁回应给谁
			--to-source 123.2.3.2-123.2.3.5
* `-j MASQUERADE` 外网地址是动态获取的时候使用(会自动选择一个外网地址), 效率比`SNAT`低.

```shell
任何时候如果接受一个来自 192.168.10.0/24 网段的用户请求的报文的时候, 都做源地址转换, 转换成 172.16.100.7 发出去

iptables -t nat -A POSTROUTING -s 192.168.10.0/24 -j SNAT --to-source 172.16.100.7
```

**示例1**

ADSL: 123.2.3.2

让房间内的所有主机都可以访问互联网. 房间的网段是`192.168.0.0/24`(网段的IP地址从192.168.0.1开始,到192.168.0.254结束).

```shell
iptables -t nat -A POSTROUTING -s 192.168.0.0/24 -j SNAT --to-source 123.2.3.2
```

如果重播了, 地址会发生变化.

**示例2**

* 修改`FORWARD`默认为`DROP`
* 放行内网`WEB`访问
* 放行内网`PING`
* 放行内网`FTP`

```shell
iptables -P FORWATD DROP

所有已建立的链接都放行
iptables -A FORWARD -m state --state ESTABLISHED -j ACCEPT

放行 web 访问
iptables -A FORWARD -s 192.168.10.0/24 -p tcp --dport 80 -m state --state NEW -j ACCEPT

放行 ping
iptables -A FORWARD -s 192.168.10.0/24 -p icmp --icmp-type 8 -m state --state NEW -j ACCEPT

放行 ftp
修改第一条规则, 让 RELATED 也放行, 这可以就可以放行 ftp 的数据链接
iptables -R FORWARD 1 -m state --state ESTABLISHED,RELATED -j ACCEPT
iptables -A FORWARD -S 192.168.10.0/24 -p tcp --dport 21 -m state --state NEW -j ACCEPT
```

限制`ftp`要确保`/etc/sysconfig/iptables-config`, 装载了`ip_nat_ftp` `ip_conntrack_ftp`这两个模块.

```
...
IPTABLES_MODULES="ip_nat_ftp ip_conntrack_ftp"
...
```

#### DNAT

目标地址转换.

公司内网有两台服务器:

1. WEB
2. FTP

但是只有一个公网地址在网关服务器上. 希望上面两个主机服务可以被外网访问.

网关服务器需要做目标地址转换, 开放`80`,`21`两个端口. 在做目标地址转换到内网提供真正服务的两台服务器上.

* `-j DNAT`
	* `--to-destionation IP[:port]` 

**示例**

* 网关
	* 外网地址: `172.16.100.7`
	* 内网地址: `198.108.10.6`
* 内网服务器: `192.168.10.22`, 网关指向`198.108.10.6`
* 外网服务器: `172.16.100.11`, 网关指向`172.16.100.7`

```
172.16.100.7 服务器上

为了测试关闭 172.16.100.11 网关指向. 通过地址转换来访问.

放行 WEB, WEB 是监听在 8080 端口. 访问网关是正常80端口, 但是内部转到web地址的8080
iptables -t nat -A PREROUTING -d 172.16.100.7 -p tcp --dport 80 -j DNAT --to-destionation 192.168.10.22:8080

不会用户访问带 h7n9 内容
iptables -A FORWARD -m string --algo kmp --string "h7n9" -j DROP


放行 FTP, FTP 在被动模式是大于1023的任意端口, 所以很难做 DNAT
```


#### 自动转换

目标地址转换自动进行源地址转换.

源地址转换自动进行目标地址转换.

### 整理知识点

---