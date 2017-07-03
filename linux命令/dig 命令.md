# dig 命令

---

`d`omain `i`nformation `g`roper. 它是一个用来灵活探测`DNS`的工具.

除非被告知请求特定域名服务器，`dig`将尝试`/etc/resolv.conf`中列举的所有服务器.

当未指定任何命令行参数或选项时,`dig`将对“`.`"(根)执行 `NS` 查询.

**常用查询选项**

每个查询选项被带前缀`+`的关键字标识.一些关键字设置或复位一个选项.通常前缀是求反关键字含义的字符串`no`.

* `+[no]tcp` 查询域名服务器时使用[不使用] TCP.缺省行为是使用 UDP,除非是 AXFR(全区域传输) 或 IXFR (增量区域传输)请求,才使用 TCP 连接.
* `+[no]vc` 查询域名服务器时使用[不使用]TCP.+[no]tcp 的备用语法提供了向下兼容.
* `+[no]ignore` 忽略 UDP 响应的中断,而不是用 TCP 重试.缺省情况运行 TCP 重试.
* `+domain=somename` 设定包含单个域 somename 的搜索列表.
* `+[no]search` 使用[不使用]搜索列表.缺省情况不使用搜索列表.
* `+[no]trace` 切换为待查询名称`从根名称服务器开始的代理路径跟踪`.缺省情况不使用跟踪.一旦启用跟踪,dig 使用迭代查询解析待查询名称.它将按照从根服务器的参照,显示来自每台使用解析查询的服务器的应答.
* `+[no]short` 提供简要答复.缺省值是以冗长格式显示答复信息.
* `+[no]identify` 当启用 `+short` 选项时,显示 [或不显示] 提供应答的 IP 地址和端口号.如果请求简短格式应答,缺省情况不显示提供应答的服务器的源地址和端口号.
* `+[no]stats` 该查询选项设定显示统计信息:查询进行时,应答的大小等等.缺省显示查询统计信息.
* `+[no]qr` 显示 [不显示] 发送的查询请求.缺省不显示.
* `+[no]question` 当返回应答时,显示 [不显示] 查询请求的问题部分.缺省作为注释显示问题部分.
* `+[no]answer` 显示 [不显示] 应答的回答部分.
* `+time=T` 为查询设置超时时间为 `T` 秒.缺省是5秒.如果将 T 设置为小于1的数,则以1秒作为查询超时时间.
* `+tries=A` 设置向服务器发送 UDP 查询请求的重试次数为 A,代替缺省的 3 次.如果把 A 小于或等于 0,则采用 1 为重试次数.
* `+[no]all` 设置或清除所有显示标志.

**示例**

一个典型的 dig 调用类似:

```shell
dig @server name type
```

* `server`: 待查询名称服务器的名称或 IP 地址.
* `name`: 将要查询的资源记录的名称.
* `type`: 显示所需的查询类型 － ANY、A、MX、SIG，以及任何有效查询类型等.如果不提供任何类型参数,dig 将对记录 A 执行查询.

```shell
dig www.dameiwang.cc A +trace

; <<>> DiG 9.8.3-P1 <<>> www.dameiwang.cc A +trace
;; global options: +cmd
.			479853	IN	NS	b.root-servers.net.
.			479853	IN	NS	c.root-servers.net.
.			479853	IN	NS	d.root-servers.net.
.			479853	IN	NS	e.root-servers.net.
.			479853	IN	NS	f.root-servers.net.
.			479853	IN	NS	g.root-servers.net.
.			479853	IN	NS	h.root-servers.net.
.			479853	IN	NS	i.root-servers.net.
.			479853	IN	NS	j.root-servers.net.
.			479853	IN	NS	k.root-servers.net.
.			479853	IN	NS	l.root-servers.net.
.			479853	IN	NS	m.root-servers.net.
.			479853	IN	NS	a.root-servers.net.
;; Received 508 bytes from 218.30.19.40#53(218.30.19.40) in 26 ms

cc.			172800	IN	NS	ac4.nstld.com.
cc.			172800	IN	NS	ac2.nstld.com.
cc.			172800	IN	NS	ac3.nstld.com.
cc.			172800	IN	NS	ac1.nstld.com.
;; Received 291 bytes from 192.228.79.201#53(192.228.79.201) in 215 ms

dameiwang.cc.		172800	IN	NS	ns16.xincache.com.
dameiwang.cc.		172800	IN	NS	ns15.xincache.com.
;; Received 84 bytes from 192.42.174.30#53(192.42.174.30) in 384 ms

www.dameiwang.cc.	3600	IN	A	121.199.19.173
;; Received 50 bytes from 113.17.175.215#53(113.17.175.215) in 127 ms
```

