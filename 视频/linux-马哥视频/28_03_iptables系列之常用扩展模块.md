# 28_03_iptables系列之常用扩展模块

---

### 笔记

---

详情见`28_02`

使用`iptables -L -n -v` 来检查是否有包匹配到.

```shell
[root@iZ944l0t308Z ~]# iptables -L -n -v
Chain INPUT (policy ACCEPT 84488 packets, 8041K bytes)
 pkts bytes target     prot opt in     out     source               destination
 1488  135K ACCEPT     tcp  --  *      *       0.0.0.0/0            120.25.87.35         tcp dpt:80 state NEW,ESTABLISHED
    0     0 ACCEPT     icmp --  *      *       0.0.0.0/0            172.16.100.7         icmptype 8 state NEW,ESTABLISHED
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            120.25.87.35         tcp dpt:80 state NEW,ESTABLISHED

Chain FORWARD (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination

Chain OUTPUT (policy ACCEPT 83041 packets, 12M bytes)
 pkts bytes target     prot opt in     out     source               destination
 1118  702K            tcp  --  *      *       120.25.87.35         0.0.0.0/0            tcp spt:80 state ESTABLISHED
    0     0 ACCEPT     icmp --  *      *       172.16.100.7         0.0.0.0/0            icmptype 0 state ESTABLISHED
 1009  694K            tcp  --  *      *       120.25.87.35         0.0.0.0/0            tcp spt:80 state ESTABLISHED
```

**示例**
 
```shell
我们的响应不能包含 h7n9
iptables -I OUTPUT -s 172.16.100.7 -m string --algo kmp --string "h7n9" -j REJECT
 
touch test.html
echo "h7n9" > test.html
请求服务器卡主
 
 [root@iZ944l0t308Z html]# iptables -L -n -v
Chain INPUT (policy ACCEPT 465 packets, 33338 bytes)
 pkts bytes target     prot opt in     out     source               destination
 1526  139K ACCEPT     tcp  --  *      *       0.0.0.0/0            120.25.87.35         tcp dpt:80 state NEW,ESTABLISHED
    0     0 ACCEPT     icmp --  *      *       0.0.0.0/0            172.16.100.7         icmptype 8 state NEW,ESTABLISHED
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            120.25.87.35         tcp dpt:80 state NEW,ESTABLISHED

Chain FORWARD (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination

Chain OUTPUT (policy ACCEPT 304 packets, 84237 bytes)
 pkts bytes target     prot opt in     out     source               destination
   46 15346 REJECT     all  --  *      *       120.25.87.35         0.0.0.0/0            STRING match  "h7n9" ALGO name kmp TO 65535 reject-with icmp-port-unreachable
 1138  718K            tcp  --  *      *       120.25.87.35         0.0.0.0/0            tcp spt:80 state ESTABLISHED
    0     0 ACCEPT     icmp --  *      *       172.16.100.7         0.0.0.0/0            icmptype 0 state ESTABLISHED
 1029  711K            tcp  --  *      *       120.25.87.35         0.0.0.0/0            tcp spt:80 state ESTABLISHED
 
 可见我们添加规则有匹配到报文.
```
 
#### 目标
 
`-j`
 
**LOG**

`LOG `和`DROP`,`ACCEPT`,`REJECT`一起用的时候一定要放到前面, 放到后面前面的规则就先匹配放行了.

这个规则并不做最后用户报文的裁定(拒绝和同意).

* `--log-level` 日志级别
* `--log-prefix` 日志前缀
* `--log-tcp-sequence` 记录tcp序列号
* `--log-tcp-options` `tcp`报文的选项
* `--log-ip-options` `ip`报文的选项
* `--log-uid` 哪个用户请求/响应的

日志记录会太多,导致`IO`太多. 需要加上速率限制.

```shell
记录ping成功记录日志,
iptables -I INPUT -d 120.25.87.35 -p icmp --icmp-type 8 -j LOG --log-prefix "firewall log for icmp"

ping 主机
ping 120.25.87.35
PING 120.25.87.35 (120.25.87.35): 56 data bytes
64 bytes from 120.25.87.35: icmp_seq=0 ttl=50 time=42.590 ms
64 bytes from 120.25.87.35: icmp_seq=1 ttl=50 time=41.440 ms
64 bytes from 120.25.87.35: icmp_seq=2 ttl=50 time=44.137 ms
64 bytes from 120.25.87.35: icmp_seq=3 ttl=50 time=39.516 ms
^C
--- 120.25.87.35 ping statistics ---
4 packets transmitted, 4 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 39.516/41.921/44.137/1.686 ms

[root@iZ944l0t308Z ~]# tail /var/log/messages
Jul 15 21:40:01 iZ944l0t308Z systemd: Starting Session 1425 of user root.
Jul 15 21:48:38 iZ944l0t308Z systemd: Started Session 1426 of user root.
Jul 15 21:48:38 iZ944l0t308Z systemd-logind: New session 1426 of user root.
Jul 15 21:48:38 iZ944l0t308Z systemd: Starting Session 1426 of user root.
Jul 15 21:50:01 iZ944l0t308Z systemd: Started Session 1427 of user root.
Jul 15 21:50:01 iZ944l0t308Z systemd: Starting Session 1427 of user root.
Jul 15 21:51:36 iZ944l0t308Z kernel: firewall log for icmpIN=eth1 OUT= MAC=00:16:3e:00:51:a9:ee:ff:ff:ff:ff:ff:08:00 SRC=36.46.50.188 DST=120.25.87.35 LEN=84 TOS=0x14 PREC=0x00 TTL=51 ID=40571 PROTO=ICMP TYPE=8 CODE=0 ID=11705 SEQ=0
Jul 15 21:51:37 iZ944l0t308Z kernel: firewall log for icmpIN=eth1 OUT= MAC=00:16:3e:00:51:a9:ee:ff:ff:ff:ff:ff:08:00 SRC=36.46.50.188 DST=120.25.87.35 LEN=84 TOS=0x14 PREC=0x00 TTL=51 ID=58758 PROTO=ICMP TYPE=8 CODE=0 ID=11705 SEQ=1
Jul 15 21:51:38 iZ944l0t308Z kernel: firewall log for icmpIN=eth1 OUT= MAC=00:16:3e:00:51:a9:ee:ff:ff:ff:ff:ff:08:00 SRC=36.46.50.188 DST=120.25.87.35 LEN=84 TOS=0x14 PREC=0x00 TTL=51 ID=13303 PROTO=ICMP TYPE=8 CODE=0 ID=11705 SEQ=2
Jul 15 21:51:39 iZ944l0t308Z kernel: firewall log for icmpIN=eth1 OUT= MAC=00:16:3e:00:51:a9:ee:ff:ff:ff:ff:ff:08:00 SRC=36.46.50.188 DST=120.25.87.35 LEN=84 TOS=0x14 PREC=0x00 TTL=51 ID=13084 PROTO=ICMP TYPE=8 CODE=0 ID=11705 SEQ=3
```

### 整理知识点

---