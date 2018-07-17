# nf_conntrack: table full,dropping packet

---

`ip_conntrack`是Linux NAT一个跟踪连接条目的模块记录着允许的跟踪连接条目`ip_conntrack`模块会记录`tcp`通讯协议的`established connection`记录.

This is normal.
Docker uses iptables connection tracking.
Connection tracking entries will subside for a while, even after the connection is closed.
There are two ways to work around this:

1. increase nf_conntrack_max (it just uses more memory, but AFAIR it's about 100 bytes per connection, so even if you set it to 2000000, you end up using 200 MB of RAM)
2. reduce the nf_conntrack_tcp_timeout values so that connections are removed faster after they are closed.

**参数调优**

理论最大值 `CONNTRACK_MAX = RAMSIZE (in bytes) / 16384 / (ARCH / 32)`,以64G的64位操作系统为例: `CONNTRACK_MAX = 64*1024*1024*1024/16384/2 = 2097152`

`CONNTRACK_MAX = RAMSIZE (以bytes记) / 16384 =RAMSIZE (以MegaBytes记) * 64`, 因此, 一个`32`位的带`512M`内存的PC在默认情况下能够处理`512*1024^2/16384 = 512*64 = 32768`个并发的netfilter连接.

我们服务器`8G`(部分`4G`) = `8*1024*1024*1024/16384/2 = 262144`, 排除使用的内存我们统一按照`4G`来计算,所以最后值是`131072`.

* `1024*1024*1024` = `1G`
* `ARCH / 32` = `64 /32` = `2`

**调优示例**

```shell
修改文件  /etc/sysctl.conf

net.netfilter.nf_conntrack_max  =  131072   
net.netfilter.nf_conntrack_tcp_timeout_established   =   300   
net.netfilter.nf_conntrack_tcp_timeout_close_wait  =   60   
net.netfilter.nf_conntrack_tcp_timeout_fin_wait  =   120   
net.netfilter.nf_conntrack_tcp_timeout_time_wait  =   120

sysctl -p 刷新生效

即时生效:
sysctl –w net.netfilter.nf_conntrack_max = xxx
```

```shell
[root@iZ944l0t308Z net]# cat /proc/sys/net/nf_conntrack_max
31768
[root@iZ944l0t308Z net]# cat /proc/sys/net/netfilter/nf_conntrack_tcp_timeout_established
432000 (120小时,5天)
[root@iZ944l0t308Z net]# cat /proc/sys/net/netfilter/nf_conntrack_tcp_timeout_close_wait
60
[root@iZ944l0t308Z net]# cat /proc/sys/net/netfilter/nf_conntrack_tcp_timeout_fin_wait
120
[root@iZ944l0t308Z net]# cat /proc/sys/net/netfilter/nf_conntrack_tcp_timeout_time_wait
120
```

#### 示例

docker 待测试服务器
```shell
登录服务(docker 部署),这台服务器部署了容器 nginx

编辑该文件 把net.netfilter.nf_conntrack_max设为200
[ansible@rancher-agent-2 ~]$ vi /etc/sysctl.conf
...
net.netfilter.nf_conntrack_max = 200
...
[ansible@rancher-agent-2 ~]$ cat /proc/sys/net/nf_conntrack_max
200
```

登录另外一台服务器(ab 测试服务器)
```shell
登录另外一台服务器
[root@iZ944l0t308Z ~]# ab -c 300 -n 1000 http://10.116.138.44/
```

docker 待测试服务器

```shell
可以看见我们在做ab测试的时候, 会出现 nf_conntrack 丢包
[ansible@rancher-agent-2 ~]$ sudo tail /var/log/messages
Jul 17 14:55:13 rancher-agent-2 kernel: nf_conntrack: table full, dropping packet
Jul 17 14:55:13 rancher-agent-2 kernel: nf_conntrack: table full, dropping packet
Jul 17 14:55:13 rancher-agent-2 kernel: nf_conntrack: table full, dropping packet
Jul 17 14:55:13 rancher-agent-2 kernel: nf_conntrack: table full, dropping packet
Jul 17 14:55:14 rancher-agent-2 kernel: nf_conntrack: table full, dropping packet
Jul 17 14:55:14 rancher-agent-2 kernel: nf_conntrack: table full, dropping packet

我们在修改net.netfilter.nf_conntrack_max把值扩大
[ansible@rancher-agent-2 ~]$ vi /etc/sysctl.conf
...
net.netfilter.nf_conntrack_max = 200
...
```