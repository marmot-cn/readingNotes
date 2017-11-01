# 36_02_Linux集群系列之二——LVS类型详解

---

## 笔记

---

### 负载均衡设备

#### 硬件负载均衡设备

`F5`: BIG IP

`Citrix`: Netscaler

`A10`

#### 软件负载均衡设备

* `4`层负载均衡设备`LVS`. 可以理解为4层交换设备, 通过`ip`地址加`端口号`来实现将用户的请求分发至后端不同的主机. 不解析更高层协议, 工作性能更好. 无法实现高级特性(类似不能根据`uri`实现负载均衡)
* `7`层实现的负载均衡设备`NGINX`,`HAPROXY`. 应用层实现. 为某些特定协议工作, 可以精确解析协议, 操作能力更强. 性能上略逊于`4`层设备.

##### `nginx`

* http
* smtp
* pop3
* imap

这些协议做负载均衡, 因为是7层代理, 只负责解析有限的7层协议.

##### `haproxy`

* http
* 基于tcp的应用协议做负载均衡(mysql, smtp)

主要应用场景在`http`协议上.

### `lvs`

`L`inux `V`irtual `S`erver

不提供任何服务, 只是将请求转发至后端服务器.

转发器`director`, 后方真正处理请求的服务器叫做`realserver`.

根据`ip`地址和`port`转发.

解析三层和四层请求, `ip`地址和`port`.

`lvs`监控在`input`链上, 一旦发现用户请求的是一个集群服务, 直接送到`postrouting`转发到其他主机上.

**工作在`input`链上**

`lvs`和`iptables`不能同时使用.

写规则是`iptables`, 真正处理请求的是`netfilter`(框架), 规则在`netfilter`生效.

`lvs`也是两段:

* `ipvsadm`: 管理集群服务的命令行工具.
* `ipvs`: 工作在内核上, 并且监控在`input`链的框架.

内核中内置`ipvs`代码.

`vip` virtual ip, 向外提供服务的ip.

客户端ip, `rsip` real server ip.

`dip` director ip.

客户端ip: `cip`

`cip` -> `vip` -> `dip` -> `rs ip`

#### lvs 类型

* `nat` 地址转换, 目标地址转换
* `dr` 直接路由
* `tun` 隧道

#### DNAT, 一般不会使用

多目标`dnat`, 工作机制和`dnat`一样.

* 源地址`cip`, 目标地址`vip`
* 通过地址转换挑选一个后端服务, 将目标地址转换. 源地址`cip`, 目标地址`rs ip`
* `rs`回应的时候, 源地址`rs ip`, 目标地址`cip`. 在将源地址转换为`vip`.

负载均衡器需要处理**进出**请求, 进出链接都要经过`director`, **压力较大**. 所能带动的`real server`有效.

`dip`通过交换机和`real server`通信. 只有`vip`才是公网地址.

##### 遵循的几个基本法则

1. 集群节点和`director`必须在同一个网络中, `dip`和`rs ip`需要通信. `rs ip`的网关必须指向`dip`. 
2. `rs`地址通常都是私有地址, 仅能用于跟`dip`通信(各集群节点之间通信).
3. `director`负责处理所有的通信(进,出).
4. 集群节点(`real server`)必须将网关指向`dip`
5. `director`支持端口映射.
6. 任何操作类型操作系统都可以用于`real server`.
7. 一个单独的`director`很可能称为集群的瓶颈.

#### DR, 一般会使用在生产环境

路由 -> 交换机 -> N个`real server`和一个`director`.

`direcotr`和`real server`只有一块网卡.

用户请求通过路由机器和交换机直接到达`director`, 根据某种挑选标准将请求直接转发至`real server`. `real server`直接响应请求给客户端.

`director`值负责调度进来的请求, `real server`直接响应客户端.

在`DR`模型中, `direcor`和`real server`都配置了`vip`. 因为如果`real server`用`rip`响应的话, 客户端会不识别的. 在`direcotr`上有两个地址`vip`和`dip`, 一个配置在网卡上, 一个配置在网卡别名上. `real server`中`vip`配置在网卡别名上, 并且是隐藏的, 不接受任何请求. 只有封装响应报文的时候用`vip`地址作为源地址响应. 真正的网卡还是本身的`rs ip`的网卡.

路由器和`director`通信的时候需要识别`vip`的`mac`地址, 使用`arp`广播, 但是`real server`不会响应广播. 因为如果`real server`和`director`对寻找`mac`地址的`arp`广播都响应, 可能谁返回最快, 谁就是`vip`对应的`mac`地址, 所以不能让`real server`不能对`mac`地址的广播请求响应.

`director`在`dr`模型中, 在转发请求给各个`real server`时, 不会通过改变目标地址来实现, 还是通过**修改mac地址**实现的. 不拆`ip`首部, 拆`mac`首部. 源`mac`改为`director`的`dmac`, 目标`mac`改为挑选到得`real server`的`mac`. 报文转发到`real server`时候, 拆掉`mac`地址时候, 目标`ip`是`vip`, 而`real server`有`vip`所有会认为到达自己的主机上的.

这时在通过`arp`找见目标`mac`对应的`ip`

请求报文小, 响应报文大. 如果不处理响应报文, 则处理能力会得到提升.

##### 遵循基本法则

1. 各集群节点必须和`director`在同一个**物理网络**中, 因为要根据`mac`地址转发. 中间没有分隔任何设备. 通过交换机链接.
2. `rs ip`可以不必是私有地址. 如果是公网地址, 一旦`director`坏了, 可以直接访问`rs ip`. 实现便捷的远程管理和通信.
3. `director`仅实现入站请求, 响应报文则由`director`直接发往客户端.
4. 集群节点**一定不能**使用`director`当做默认网关. 而是直接使用前端网关.
5. `director`不支持端口映射. 因为响应报文直接由`real server`响应. (客户端请求80端口, `real server`必须监听`80`端口, 用`80`端口响应. 而且因为修改的是`mac`地址.

#### 隧道

工作机制和`dr`一样. 转发的时候重新封装`ip`报文.

`dip`和`real server`不在同一个网络上.

转发时, 在`cip`和`vip`外面再加一个`ip`首部. 源ip`dip`目标地址`rs ip`.

通过第一层`ip`报文(在加的一层).送到目标`real server`, 拆掉第一层发现`cip`和`vip`, 自己主机用`vip`可以直接响应.

借助一个`ip`报文在发送一个`ip`报文称为隧道.

要求`director`和`real server`支持隧道机制.

##### 遵循基本法则

1. 集群节点可以跨越互联网, 没必要在同一个网络地址中.

## 整理知识点

---

### 端口映射

1. 访问端口是`80`, 后端服务器使用不同的端口响应, 如`8080`. 
2. `real server`的`ip`必须是公网地址.
3. `director`仅处理入站请求, 和`dr`一样. 响应报文直接由`real server`响应发往客户端.
4. 响应报文一定不能通过`director`.
5. 只有支持隧道功能的`os`才能用于`real server`.
6. 不支持端口映射.


### F5