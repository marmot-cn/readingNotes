# 28_02_iptables系列之基本应用及显式扩展

---

### 笔记

---

#### iptables

iptables不是服务, 但有服务脚本; 服务脚本的主要作用在于**管理保存的规则**. 装载及溢出`iptables/nerfilter`相关的内核模块.

服务启动后, 此前保存的规则也会生效.

```shell
iptables [-t TABLE] COMMAND CHAIN [num] 匹配条件 -j 处理动作
```
**匹配标准**

* 通用(普通)匹配
	* `-s, --src`: 指定源地址
	* `-d, --dst`: 指定目标地址
	* `-p {tcp|udp|icmp}`: 指定协议
	* `-i INTERFACE`: 指定数据报文流入的接口(从哪个网卡).`PREROUTING` 一般指定 `-i`.
		* `PRE_ROUTING`
		* `INPUT`
		* `FORWARD`
	* `-o INTERFACE`: 指定数据报文流出的接口(从哪个网卡).`POSTROUTING` 一般指定 `-o`.
		* `OUTPUT`
		* `POSTROUTING`
		* `FORWARD` 
* 扩展匹配(使用`netfilter`扩展模块)
	* 隐含扩展: 不用特别指明由哪个模块进行的扩展, 因为此时使用 `-p {tcp|udp|icmp}` 已经指定了协议(`-p tcp -m tcp` = `-p tcp`).
		* `-p tcp`(指定协议之后,使用特定协议的扩展):
			* `--sport PORT`: 源端口
			* `--dport PORT`: 目标端口 
			* `--tcp-flags mask comp`: `tcp`的标志位. 只检查`mask`指定的标志位, 是逗号分隔的标志位列表. `comp`: 此列表中出现在`mask`中, 且必须为`1`. `comp`中没出现, 而`mask`中出现的, 必须为`0`.
				* 标记
					* SYN
					* ACK
					* FIN
					* RST
					* URG
					* PSH
					* ALL
					* NONE
				
							---tcp-flags SYN.FIN,ACK.RST SYN,ACK
							两个列表使用逗号隔开.
							将检查tcp报文的 SYN.FIN,ACK.RST 这4个标志位,
							且只能 SYN,ACK 为1.
							剩下的都必须为0.
							
							--tcp-flags SYN,FIN,ACK,RST SYN = --syn 三次握手的第一次.
			* `--syn`: 三次握手的第一次.
		* `-p icmp`
			* `0`: `echo-replay` 响应报文
			* `8`: `echo-request` 请求报文
		* `-p udp`
			* `--sport`
			* `--dport`
	* 显示扩展(使用额外的条件匹配机制,数据包报文的速度,链接的状态,根据时间放行报文...): 必须指明由哪个模块进行的扩展(因为这些扩展和协议没有关系,所以必须手动指定).
		* `-m EXTENSION(扩展名称) --spe-opt(扩展独有选项)` 
			* `-m state`: 状态扩展, 结合`ip_conntract`追踪会话状态(根据ip追踪状态),(tcp,udp,icmp 三种协议都可以追踪).`/proc/sys/net/netfilter/`配置文件
				* `NEW`: 新请求(类似`tcp`的第一次握手,但是不限于`tcp`)
				* `ESTABLISHED`: 已建立的链接
				* `INVALID`: 非法链接请求(比如: `SYN=1,FIN=1`)
				* `RELATED`: 相关联(类似`FTP`的命令链接激活数据链接)
						
						状态为 NEW 和 ESTABLISHED 都放行
						-m state --state NEW,ESTABLISED -j ACCEPT
			* `-m multiport`: 离散多端口匹配扩展
				* `--source-ports`
				* `--destination-ports`
				* `--ports`: 即是目标端口,又是源端口 

						-m multiport --destination-ports 21,22,80 -j ACCEPT
			* `-m iprange`: 
				* `--src-range (ip-ip)`: 源地址范围
				* `--dst-range (ip-ip)`: 目标地址范围

						放行给这个地址范围之外的主机访问 
						iptables -A INPUT -p tcp -m iprange ! --src-range 172.16.100.3-172.16.100.100 --dport 22 -m state --state NEW,ESTABLISHED -j ACCEPT
			* `-m connlimit`(连接数限定,限定某一个ip地址最多可以同时发起几个链接):
				* `--connlimit-above #`: 上限,最多允许使用多少个链接. 低于`n`个开发. 通常加`!`使用

						低于两个就允许, 通常加叹号使用
						iptables -A INPUT -d 172.16.100.7 -p tcp --dport 80 -m connlimit ! --connlimit-above 2 -j ACCEPT
						
						不加叹号
						iptables -A INPUT -d 172.16.100.7 -p tcp --dport 80 -m connlimit ! --connlimit-above 2 -j DROP
						
			* `-m limit`: 限定流量上限(令牌桶控制流量上限,不控制最大上限,只控制单位时间内的流量上限和一次性蜂拥而至的上限).
				* `--limit rate` 单位时间内最多可以允许`rate`个人进来.
				* `--limit-burst number` 一次性并发涌至的请求数(第一次进来的请求)
						
						每秒钟ping请求每分钟至允许ping5次(burst默认是5个).第一次5个进来后就会按照每分钟5个速率进来.
						现在默认策略是DROP
						iptables -A INPUT -d 172.16.100.7 -p icmp --icmp-type 8 -m limit limit 5/minute -j ACCEPT
						iptables -A OUTPUT -s 172.16.100.7 -m state --state REALTED,ESTABLISHED -j ACCEPT
			
			* `-m stirng`: 内核>2.6.14
				* `-alog(必须)`: 字符串匹配算法
					* `bm`
					* `kmp`
				* `--from offset`
				* `--to offset`
				* `--string(必须)`: 匹配哪一个字符串
				* `--hex-string`: 使用十六进制
						
						只要请求包含 h7n9 就拒绝.无法检查文件内容,
						因为我们请求页面, 这个数据需要在OUTPUT响应给用户.加入我们添加到INPUT链,则只能过滤文件名字
						iptables -I OUTPUT -s 172.16.100.7 -m string --algo kmp --string "h7n9" -j REJECT
						
			* `-m recent`: 限制一段时间内的连接数
						
**条件取反**

`!`, 使用叹号. 几乎所有的条件都可以取反.

```shel
-s ! 172.16.100.6
```

**命令**

* 管理规则命令
	* `-A`: 附加一条规则(链的尾部追加一条规则), 添加在链的尾部.
	* `-I CHAIN [num]`: 插入一条规则, 指定添加在什么位置. 插入为对应`CHAIN`上的第`num`条, 省略`num`,则`num`为`1`.
	* `-D CHAIN [num]`: 删除指定链中第`num`条规则. 也可以指定匹配条件.
	* `-R CHAIN [num]`: 替换指定的规则.
* 管理链命令
	* `-F [CHAIN]`: `flush`, 清空指定规则链(`-D`删除一条, 这个命令是删除全部). 如果省略`CHAIN`, 则可以实现删除对应表中的所有链.
	* `-P CHAIN {ACCEPT|DROP}`: 设定指定链的默认策略.
	* `-N`: 自定义一条新的空链.
	* `-X`: 删除一个自定义的空链. (必须是空的, 非空的要先用 `-F` 清空)
	* `-Z`: 置零指定链中的所有规则的计数器.
	* `-E old-chain-name new-chain-name`: 重命名一条自定义的链.
* 查看类:
	* `-L`: 显示指定表中的所有规则, 分链显示.
		* `-n`: 以数字格式显示主机地址和端口号(默认会反解)
		* `-t`: 默认是`filter`链.
		* `-v`: 显示链及规则的详细信息.
		* `-x`: 显示计数器的精确值.
		* `--line-numbers`: 显示规则号码.
		
**动作**

动作(target):

* `ACCEPT`: 允许通过
* `DROP`: 丢弃
* `REJECT`: 拒绝, 返回信息
* `DNAT` 目标地址转换
* `SNAT` 源地址转换
* `REDIRECT` 端口重定向
* `MASQUERADE` 地址伪装
* `LOG` 记录日志
* `MARK` 设定标记
* `MIRROR`
* `NOTRACK`

**示例**

```shell
[root@iZ944l0t308Z ~]# iptables -L -n
Chain INPUT (policy ACCEPT)  #默认策略是 ACCEPT
target     prot opt source               destination

Chain FORWARD (policy ACCEPT)
target     prot opt source               destination

Chain OUTPUT (policy ACCEPT)
target     prot opt source               destination
```

```shell
本机 172.16.100.7, 开放的有sshd, 监听在本地 22/tcp. 放行来自于172.16.0.0对于本机的ssh访问.

放行进来的报文
iptables -t filter -A INPUT -s 172.16.0.0/16 -d 172.16.100.7 -p tcp --dport 22 -j ACCEPT

放行响应的报文
iptables -t filter -A OUTPUT -s 172.16.100.7 -d 172.16.0.0/16 -p tcp --sport 22 -j ACCEPT

设定默认策略关闭
iptables -P INPUT DROP
iptables -P OUTPUT DROP 
iptables -P FORWARD DROP 

注意: 后续我们的 -j 都是 ACCEPT, 因为我们默认都是 DROP

允许所有访问本机web服务, web服务比ssh服务访问量大, 放上面. (-s 0.0.0.0 可以不写).
ping 不通是因为 ping 请求的是 icmp 协议.
本机也 ping 不通, 虽然 127.0.0.1 就是本地回环地址, 从网路角度来讲还要从 OUTPUT 出去, 从 INPUT 进来.
iptables -I INPUT -d 172.16.100.7 -p tcp --dport 80 -j ACCEPT

这样设置, 本地就可以 ping 通了
iptables -A INPUT -s 127.0.0.1 -d 127.0.0.1 -i lo(本地回环地址) -j ACCEPT
iptables -A OUTPUT -s 127.0.0.1 -d 127.0.0.1 -o lo(本地回环地址) -j ACCEPT

允许自己ping别人, ping 出去. 但是别人不能 ping 自己
iptables -A OUTPUT -s 172.16.100.7 -p icmp --icmp-type 8 -j ACCEPT
iptables -A INPUT -d 172.16.100.7 -p icmp --icmp-type 0 -j ACCEPT


自己是DNS服务器, 负责为本地网路递归解析. 监听在 udp/tcp:53
如果不是自己域, 还要去找"根". 自己服务器作为客户端链接别人. 

1. 别人找自己解析
2. 自己响应别人的解析请求
3. 自己作为客户端去找"根"
4. 自己作为客户端接收找"根"的响应

需要 4*2(tcp,udp) 8条规则
```

**示例2**

sercie: 172.16.100.7 放行 sshd, httpd 服务. 需要追踪链接(状态监测).

服务器:

1. 一般只接受别人的请求(请求我们示例的两个模块`sshd`,`httpd`). 
2. 出去的是对某个服务的响应.

进来的链接状态:

* `NEW`
* `ESTABLISHED`, 第一次是`NEW`, 后续都是`ESTABLISHED`状态 

出去的链接状态:

* `ESTABLISHED`, 因为服务器不会去请求别人.

```shell
进来的请求只能是 NEW,ESTABLISHED.
iptables -A INPUT -d 172.16.100.7 -p tcp --dport 22 -m state --state NEW,ESTABLISHED -j ACCEPT

放行出去的ESTABLEISHED链接
iptables -A OUTPUT -s 172.16.100.7 -p tcp --sport 22 -m state --state ESTABLISED

该默认策略

iptables -P INPUT DROP
iptables -P OUTPUT DROP

设置网络链
iptables -A INPUT -d 172.16.100.7 -p tcp --dport 80 -m state --state NEW,ESTABLISHED -j ACCEPT
iptables -A OUTPUT -s 172.16.100.7 -p tcp --sport 80 -m state --state ESTABLISHED

运行这些命令会自动加载 `conntrack` 模块
[root@iZ944l0t308Z ~]# lsmod | grep ip
nf_conntrack_ipv4      14862  2
nf_defrag_ipv4         12729  1 nf_conntrack_ipv4
nf_conntrack          101024  3 xt_conntrack,nf_conntrack_netlink,nf_conntrack_ipv4
iptable_filter         12810  1
ip_tables              27239  1 iptable_filter
ipt_REJECT             12541  0
```

**示例3**

放行别人`ping`自己.

```shell
进来的
iptables -A INPUT -d 172.16.100.7 -p icmp --icmp-type 8 -m state --state NEW,ESTABLISHED -j ACCEPT
出去的
iptables -A OUTPUT -s 172.16.100.7 -p icmp --icmp-type 0 -m state --state ESTABLISHED -j ACCEPT
```

**多条链合为一条**

```shell
[root@iZ944l0t308Z ~]# iptables -L -n --line-numbers
Chain INPUT (policy ACCEPT)
num  target     prot opt source               destination
1    ACCEPT     tcp  --  0.0.0.0/0            120.25.87.35         tcp dpt:80 state NEW,ESTABLISHED
2    ACCEPT     icmp --  0.0.0.0/0            172.16.100.7         icmptype 8 state NEW,ESTABLISHED
3    ACCEPT     tcp  --  0.0.0.0/0            120.25.87.35         tcp dpt:80 state NEW,ESTABLISHED

Chain FORWARD (policy ACCEPT)
num  target     prot opt source               destination

Chain OUTPUT (policy ACCEPT)
num  target     prot opt source               destination
1               tcp  --  120.25.87.35         0.0.0.0/0            tcp spt:80 state ESTABLISHED
2    ACCEPT     icmp --  172.16.100.7         0.0.0.0/0            icmptype 0 state ESTABLISHED
3               tcp  --  120.25.87.35         0.0.0.0/0            tcp spt:80 state ESTABLISHED
```

可以看见后面三条链都是`ESTABLISHED`. 我们合为一条.

```shell
iptables -I OUTPUT -s 172.16.100.1 -m statue --state ESTABLISHED -j ACCEPT
```

#### `state` 里面的 `related`

内核首先要装载`ip_conntrack_ftp`和`ip_nat_ftp`模块.

```shel
这个命令是视频上的, 我本机cenots7 未必是这个模块
vi /etc/sysconfig/iptables-config
编辑 IPTABLES_MODULES="ip_nat_ftp ip_conntrack_ftp"
```

**放行FTP端口**

* 放行21端口 
* 本地`ftp`会去链接数据库

```shell
放行21
iptables -A INPUT -d 172.16.100.7 -p tcp --dport 21 -m state --state NEW,ESTABLISHED

放行本地端口, 可以让ftp链接数据库
iptables -A INPUT -i lo -j ACCEPT
inptables -A OUTPUT -o lo -j ACCEPT
```

但是登录成功后执行 ls 没有结果, 因为我们只开放了21端口. 如果我们使用被动模式, 需要让客户端链接我们一个随机端口. 但是我们现在没有开放随机端口. 有如下两个方法可以解决:

1. 开放1023以后的所有端口(防火墙没有意义)
2. 状态追踪(`related`)

**`related`**

只要状态是`related`和命令链接有关系, 就统统放行.

```shell
没指定端口, 是因为端口没确定
iptables -A INPUT -d 172.16.100.7 -p tcp -m state --state ESTABLISHED,RELATED -j ACCEPT


因为我们原来有规则"iptables -I OUTPUT -s 172.16.100.1 -m statue --state ESTABLISHED -j ACCEPT", 我们现在修改规则
iptables -R OUTPUT 1 -s 172.16.100.7 -m state --state ESTABLISHED,RELATED -j ACCEPT
```

#### iptables 重启会清空整个表中的每个链, 重新加载配置文件

`iptables` 启动会加载 `/etc/sysconfig/iptables`配置文件.

```shell
[root@iZ944l0t308Z ~]# iptables -t filter -A INPUT -d 120.25.87.35 -p tcp --dport 80 -j ACCEPT
[root@iZ944l0t308Z ~]# service iptables save
iptables: Saving firewall rules to /etc/sysconfig/iptables:[  OK  ]
[root@iZ944l0t308Z ~]# service iptables restart
Redirecting to /bin/systemctl restart  iptables.service
[root@iZ944l0t308Z ~]# iptables -L -n
Chain INPUT (policy ACCEPT)
target     prot opt source               destination
ACCEPT     tcp  --  0.0.0.0/0            172.16.100.7        tcp dpt:80

Chain FORWARD (policy ACCEPT)
target     prot opt source               destination

Chain OUTPUT (policy ACCEPT)
target     prot opt source               destination
```

运行`service iptables save`, 默认会保存在该配置文件. 重启后该规则仍然存在.

**`iptables-save` 和 `iptables-store`**

保存到其他文件: `iptables-save > /path/file`, 但是下次不会重新加载

读取`iptables`规则: `iptables-restore < /path/file`, 从文件加载

#### 合并链接

```shell
[root@iZ944l0t308Z ~]# iptables -L -n
Chain INPUT (policy ACCEPT)
target     prot opt source               destination
ACCEPT     tcp  --  0.0.0.0/0            172.16.100.7         tcp dpt:22 state NEW,ESTABLISHED
ACCEPT     icmp --  0.0.0.0/0            172.16.100.7         tcp dpt:80 state NEW,ESTABLISHED
ACCEPT     tcp  --  0.0.0.0/0            172.16.100.7         tcp dpt:21 state NEW,ESTABLISHED
```

我们可见前三条都有`ESTABLISHED`, 那么我们创建一条

```shel
iptables -A INPUT -d 172.16.100.7 -p tcp -m state --state RELATED,ESTABLISHED -j ACCEPT
```

这样后面的三条都可以不用匹配`ESTABLISHED`状态.

```shell
ACCEPT     tcp  --  0.0.0.0/0            172.16.100.7         tcp dpt:22 state NEW
ACCEPT     icmp --  0.0.0.0/0            172.16.100.7         tcp dpt:80 state NEW
ACCEPT     tcp  --  0.0.0.0/0            172.16.100.7         tcp dpt:21 state NEW
```

后面的三条合为一条.

```shell
iptables -I INPUT 2 -d 172.16.100.7 -p tcp -m multiport --destination-ports 21,22,80 -m state --state NEW -j ACCEPT
```

后面三条都可以删除了.

**最终**

```shel
iptables -A INPUT -d 172.16.100.7 -p tcp -m state --state RELATED,ESTABLISHED -j ACCEPT
iptables -I INPUT 2 -d 172.16.100.7 -p tcp -m multiport --destination-ports 21,22,80 -m state --state NEW -j ACCEPT
```

#### 本机测试

**另外一种禁ping方法**

```shell
iptables -t filter -A INPUT -d 120.25.87.35 -p icmp -j DROP
[root@iZ944l0t308Z ~]# iptables -L
Chain INPUT (policy ACCEPT)
target     prot opt source               destination
DROP       icmp --  anywhere             iZ944l0t308Z

Chain FORWARD (policy ACCEPT)
target     prot opt source               destination

Chain OUTPUT (policy ACCEPT)
target     prot opt source               destination

其他主机
ping 120.25.87.35
PING 120.25.87.35 (120.25.87.35): 56 data bytes
Request timeout for icmp_seq 0
Request timeout for icmp_seq 1
Request timeout for icmp_seq 2

删除
iptables -d INPUT 1

其他主机
ping 120.25.87.35
PING 120.25.87.35 (120.25.87.35): 56 data bytes
64 bytes from 120.25.87.35: icmp_seq=0 ttl=51 time=65.145 ms
64 bytes from 120.25.87.35: icmp_seq=1 ttl=51 time=44.666 ms
```

```shell
可以ping出去, 都不是能收到返回信息 icmp-8
[root@iZ944l0t308Z ~]# iptables -A INPUT -d 120.25.87.35 -p icmp --icmp-type 0 -j DROP
[root@iZ944l0t308Z ~]# ping www.baidu.com

```

#### ip_conntract 链接追踪

是一个内核模块, 实时记录客户端和服务器端彼此正在建立的链接关系, 并且能够追踪到哪一个链接和另外一个链接处于什么状态.

`nf_conntrack_ipv4`模块

```shell
[root@iZ944l0t308Z ~]# lsmod | grep ip
ipt_REJECT             12541  2
nf_conntrack_ipv4      14862  2
nf_defrag_ipv4         12729  1 nf_conntrack_ipv4
nf_conntrack          101024  3 xt_conntrack,nf_conntrack_netlink,nf_conntrack_ipv4

卸载该模块
[root@iZ944l0t308Z ~]# modprobe -r nf_conntrack_ipv4
modprobe: FATAL: Module nf_conntrack is in use.

[root@iZ944l0t308Z ~]# service iptables stop
Redirecting to /bin/systemctl stop  iptables.service
[root@iZ944l0t308Z ~]# lsmod | grep ip
ipt_REJECT             12541  0

关闭后无该文件
[root@iZ944l0t308Z ~]# cat /proc/net/nf_contrack
cat: /proc/net/nf_contrack: No such file or directory

我手动装载模块没有实现, 以下为视频流程. 视频使用的是 ip_conntrack
modprobe ip_conntrack
cat /proc/net/ip_conntrack
出现数据,且文件存在.
```


`/proc/net/ip_conntrack` 内核文件,位于内存当中,保存于当前系统上每一个客户端和当前主机所建立的链接关系.

*我的centos 7上面只有`/proc/net/nf_contrack`貌似nf_conntrack 工作在3 层，支持IPv4 和IPv6，而ip_conntrack 只支持IPv4. nf_conntrack/ip_conntrack 跟 nat 有关，用来跟踪连接条目，它会使用一个哈希表来记录 established 的记录.nf_conntrack 在 2.6.15 被引入,而 ip_conntrack 在 2.6.22 被移除*

```shell
当前主机 ssh 链接
[root@iZ944l0t308Z ~]# cat  /proc/net/nf_conntrack
ipv4     2 tcp      6 299 ESTABLISHED src=61.185.200.235 dst=120.25.87.35 sport=57360 dport=22 src=120.25.87.35 dst=61.185.200.235 sport=22 dport=57360 [ASSURED] mark=0 zone=0 use=2
```

在一次链接里面记录了`来回`两次链接.

上面示例的`来回`两次

* src=61.185.200.235 dst=120.25.87.35 sport=57360 dport=22
* src=120.25.87.35 dst=61.185.200.235 sport=22 dport=57360

**最多保存多少条目**

最多能保存多少个条目:

```shell
[root@iZ944l0t308Z ~]# cat /proc/sys/net/nf_conntrack_max
31768
```

最多可以同时追踪多少个链接, 超出这个数字, 会因为超时而丢弃.

最好不要启动这个模块, 或者调大这个`max`值.

**`iptstate`**

类似于`top`命令的`iptables`表状态的条目显示工具.

```shell
[root@iZ944l0t308Z ~]# yum install iptstate
...
[root@iZ944l0t308Z ~]# iptstate
Version: 2.2.5        Sort: SrcIP           b: change sorting   h: help
Source                                                                                   Destination                                                                              Prt State       TTL
10.170.148.109:33312                                                                     10.143.34.200:80                                                                         tcp TIME_WAIT     0:01:29
61.185.200.235:57360                                                                     120.25.87.35:22                                                                          tcp ESTABLISHED 119:59:59
120.25.87.35:37346
```

### 整理知识点

---

#### `ip_conntrack_max`

超过这个连接数, 链接会丢弃. 有两种思路

* 修改`max`值
* 减少`ESTABLISHED`的`TTL`时间

修改文件对应的值即可.

**修改`max`值**

修改链接的文件`/proc/sys/net/ipv4/ip_conntrack_max`

我服务器是是`/proc/sys/net/nf_conntrack_max`

**修改`ttl`时间**

课程示例文件地址是`/proc/sys/net/ipv4/ietfilter/ip_conntrack_tcp_timeout_established`

我服务器是`/proc/sys/net/netfilter/nf_conntrack_tcp_timeout_established`

```shell
cat /proc/sys/net/netfilter/nf_conntrack_tcp_timeout_established
432000

432000 秒等于 120 小时, 就是 iptstate 里面 ESTABLISHED 状态的 TTL 时间.
```

#### `ip_conntrack` 追踪功能丢弃链接

服务器不小心执行了`iptables -t nat -L` 激活了 `ip_conntrack` 模块.

最多可以同时追踪多少个链接, 超出这个数字, 会因为超时而丢弃. 超过`max`限制, 多余链接被丢弃.

#### http 默认长链接短连接

个人观点;

`http` 默认使用 `tcp` 链接.

```shell
61.185.200.235:60993   120.25.87.35:80   tcp TIME_WAIT     0:00:23
61.185.200.235:61126   120.25.87.35:80   tcp TIME_WAIT     0:01:38

可以看见刷新一次以后 TIME_WAIT  会变为 ESTABLISHED
61.185.200.235:61126   120.25.87.35:80   tcptcp ESTABLISHED 119:59:54
```

所以我个人观点, `tcp` 并没有断开, 会在一定时间后断开.

#### `ping`

**`icmp`**

互联网报文控制信息协议.

**`ping`别人**

`ping` 出是 `icmp` 类型 `8`.

`ping` 入是 `icmp` 类型 `0`.

**别人`ping`自己**

`ping` 入是 `8`.

`ping` 出是 `0`.