# 测试通过路由设置外网ping内网地址

---

## 场景

* `A`主机
	* `enp0s3`: 192.168.0.178
	* `enp0s8`: 192.168.56.102
* `B`主机
	* `enp0s3`: 192.168.0.180
	* `enp0s8`: 192.168.56.103
* `C`主机
	* `enp0s3`: 192.168.0.179
	* `enp0s8`: 192.168.56.101

## 测试

把`C`主机的`enp0s8`网卡关掉

```shell
ifconfig enp0s8 down
```

这是`ping 192,168.56.102`地址会`ping`不通.

再给`C`主机设定路由规则

```
ip route add 192.168.56.102 via 192.168.0.178

[root@localhost ~]# route -n
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
0.0.0.0         192.168.0.1     0.0.0.0         UG    100    0        0 enp0s3
192.168.0.0     0.0.0.0         255.255.255.0   U     100    0        0 enp0s3
192.168.56.102  192.168.0.178   255.255.255.255 UGH   0      0        0 enp0s3
```

此时就可以`ping`通.