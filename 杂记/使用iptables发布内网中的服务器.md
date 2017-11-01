# 使用iptables发布内网中的服务器

---

在`nat`表中添加

```
从外网网卡流入的流量到80端口, 做目标地址转换为 10.44.88.189
-A PREROUTING -i eth1 -p tcp -m tcp --dport 80 -j DNAT --to-destination 10.44.88.189

从内网网卡流出把源地址转换为 10.116.138.44, 代理主机的ip地址
-A POSTROUTING -d 10.44.88.189/32 -o eth0 -p tcp -m tcp --dport 80 -j SNAT --to-source 10.116.138.44
```

* `eth0`是内网的网卡
* `eth1`是外网的网卡
* `10.44.88.189` 目标内网服务器的地址, 里面装了nginx绑定80端口
* `10.116.138.44` 需要代理主机的内网地址
* 还需要把代理主机的转发打开`echo 1  > /proc/sys/net/ipv4/ip_forward`

这样访问代理主机的公网ip地址, 就可以访问内网机器的`nginx`上

