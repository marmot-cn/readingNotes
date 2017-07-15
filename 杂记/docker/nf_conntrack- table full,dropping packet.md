# nf_conntrack: table full,dropping packet

---

This is normal.
Docker uses iptables connection tracking.
Connection tracking entries will subside for a while, even after the connection is closed.
There are two ways to work around this:

1. increase nf_conntrack_max (it just uses more memory, but AFAIR it's about 100 bytes per connection, so even if you set it to 2000000, you end up using 200 MB of RAM)
2. reduce the nf_conntrack_tcp_timeout values so that connections are removed faster after they are closed.

**调优示例**

```shell
修改文件  /etc/sysctl.conf

net.netfilter.nf_conntrack_max  =   1048576   
net.netfilter.nf_conntrack_tcp_timeout_established   =   300   
net.netfilter.nf_conntrack_tcp_timeout_close_wait  =   60   
net.netfilter.nf_conntrack_tcp_timeout_fin_wait  =   120   
net.netfilter.nf_conntrack_tcp_timeout_time_wait  =   120

sysctl -p 刷新生效
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