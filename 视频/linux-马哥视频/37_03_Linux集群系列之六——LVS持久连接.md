# 37_03_Linux集群系列之六——LVS持久连接

---

## 笔记

---

`declare -a A`明确声明`A`为数组.

### LVS持久链接

在一定时间内, 将来自某一个客户端的请求始终定向到同一个`real server`, 持久链接不关乎调度算法.

持久链接模板: 一段内存缓冲区域. 能够将每一个第一次发起链接请求的ip, 以及给该链接选定的`real server`, 在一定时间内只要是这个客户端再次访问, 就查询会话表(尚未超时), 就会分配到同一个`real server`上.

#### 查看持久化链接:

```
[root@localhost ~]# ipvsadm -L --persistent-conn
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port            Weight    PersistConn ActiveConn InActConn
  -> RemoteAddress:Port
TCP  localhost.localdomain:http wlc
  -> localhost:http               5         0           0          4
  -> 192.168.0.179:http           2         0           0          2
  -> 192.168.0.180:http           1         0           0          1
```

#### 查看持久化链接模板:

```
[root@localhost ~]# ipvsadm -L -c
IPVS connection entries
pro expire state       source             virtual            destination
TCP 00:32  FIN_WAIT    192.168.0.201:53555 localhost.localdomain:http localhost:http
TCP 00:31  FIN_WAIT    192.168.0.201:53554 localhost.localdomain:http localhost:http
TCP 00:33  FIN_WAIT    192.168.0.201:53556 localhost.localdomain:http 192.168.0.179:http
TCP 00:28  FIN_WAIT    192.168.0.201:53551 localhost.localdomain:http 192.168.0.179:http
TCP 00:30  FIN_WAIT    192.168.0.201:53553 localhost.localdomain:http 192.168.0.180:http
TCP 00:33  FIN_WAIT    192.168.0.201:53557 localhost.localdomain:http localhost:http
TCP 00:29  FIN_WAIT    192.168.0.201:53552 localhost.localdomain:http localhost:http
```

* `source`: 源地址.
* `virtual`: 分配的虚拟服务器.
* `destination`: 目标服务器.
* `expire`: 过期时长.
* `state`: 当前链接状态.

#### 启用`ipvs`持久链接功能

定义一个集群服务时候, 添加`-p`选项.

`ipvsadm -A|E .... -p timeout`

* `timeout`: 持久链接时长, 默认为`300`秒.

##### 基于`SSL`, 需要持久链接

否则每次换链接都需要建立`ssl`会话.

##### 实现来自同一个客户端的所有请求都定向同一个`real server`上

假设后面的`real server`都提供`web`和`telnet`2类集群服务. 使用同一个`director`提供这2类集群服务的负载均衡.

同一个客户端请求到`web`请求, 在请求`telnet`时候也要定位到同一个`real server`.

* `PPC`(持久端口链接): 将来自于同一个客户端对同一个集群服务的请求, 始终定向至此前选定的`RS`.**只对一个端口进行持久, 只对集群内的一个服务持久, 集群内的不同服务不持久**
* `PCC`(`Persistent client connection`)(持久客户端链接): 将来自于同一个客户端对所有端口的请求都能够始终定向至此前选定的`RS`.
* `PNMPP`(持久防火墙标记链接), 特定机制, 定义端口之间的姻亲关系.

防火墙标记: `PREROUTING`链上对于目标端口为`80`和目标端口为`23`打相同的一个标记`10`.这时候把防火墙标记为`10`的端口定义为一个集群. **可以自有定义随意归并服务**.

```
80 : RS1
23 : 同一个RS
```

#### 示例

添加持久链接

```
[root@localhost ~]# ipvsadm -L -n
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  192.168.0.101:80 wlc
  -> 127.0.0.1:80                 Route   5      0          0
  -> 192.168.0.179:80             Route   2      0          0
  -> 192.168.0.180:80             Route   1      0          0
[root@localhost ~]# ipvsadm -E -t 192.168.0.101:80 -p 600
[root@localhost ~]# ipvsadm -L -n
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  192.168.0.101:80 wlc persistent 600
  -> 127.0.0.1:80                 Route   5      0          0
  -> 192.168.0.179:80             Route   2      0          0
  -> 192.168.0.180:80             Route   1      0          0
```

测试, 刷新访问一直定向到同一个`rs`上.

```
➜  marmot-sourcecode git:(dev) ✗ curl 192.168.0.101
server-1
➜  marmot-sourcecode git:(dev) ✗ curl 192.168.0.101
server-1
➜  marmot-sourcecode git:(dev) ✗ curl 192.168.0.101
server-1
➜  marmot-sourcecode git:(dev) ✗ curl 192.168.0.101
server-1
```

可见已经分配到一个持久链接上.
```
[root@localhost ~]# ipvsadm -L -n --persistent-conn
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port            Weight    PersistConn ActiveConn InActConn
  -> RemoteAddress:Port
TCP  192.168.0.101:80 wlc persistent 600
  -> 127.0.0.1:80                 5         0           0          0
  -> 192.168.0.179:80             2         1           0          4
  -> 192.168.0.180:80             1         0           0          0
 ```
 
#### `pcc`示例

```
清空所有规则
ipvsadm -C

ipvsadm -A -t xx.xxx.xxx.xxx:0 -s rr -p 600
ipvsadm -a -t xx.xxx.xxx.xxx:0 -r xxx.xxx.xxx.xxx -g -w 2
ipvsadm -a -t xx.xxx.xxx.xxx:0 -r xxx.xxx.xxx.xxx -g -w 1
```

把`0`端口定义为集群服务, 这意味着所有端口的服务请求都是集群服务.

只要使用`pcc`就意味着把所有服务(所有端口)都定义成集群服务, 一律像`rs`转发. 如果`rs`没有提供相应的服务, 则会报错. 
 
#### 防火墙标记

打标记示例

````
172.16.100.6 是 vip
iptables -A PREROUTING -i eth0 -t mangle -p tcp -d 172.16.100.6 --dport 80 -j MARK --set-mark 1

标记是 0-99 之间的整数
``` 

把`80`和`23`归并为同一个服务.
```
清空
ipvsad -C

iptables -A PREROUTING -i eth0 -t mangle -d 192.16.10.3 -p tcp --dport 80 -j MARK --set-mark 8
iptables -A PREROUTING -i eth0 -t mangle -d 192.16.10.3 -p tcp --dport 23 -j MARK --set-mark 8

将防火墙标记为8的定义为集群服务

ipvsadm -A -f 8 -s rr
ipvsadm -a -f 8 -r 192.168.10.7 -g -2 -w 2
ipvsadm -a -f 8 -r 192.168.10.8 -g -2 -w 5
```
 
## 整理知识点

---

### rsync+inotify

对于大文件同步效率很低.

大文件场景可以使用`sersync`.

