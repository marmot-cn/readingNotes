# SOCK_RAW 与 SOCK_STREAM 、SOCK_DGRAM 的区别

---

![](./img/20190523_1.jpg)

* `SOCK_STREAM(TCP)`, `SCOK_DGRAM(UDP)` 工作在传输层
* `SOCK_RAW`工作在网络层, 可以处理`ICMP`, `IGMP`等网络报文