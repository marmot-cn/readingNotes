# Nginx的keepalive_timeout和TCP的keepalive 

---

`keep alive`分为`HTTP`层的`Keep-Alive`和`TCP`层的`KeepAlive`, 两者是不同的概念. (注意写法不一样)

### `HTTP Keep-Alive`

`http`早期, 每个`http`请求都要求打开一个`tcp socket`连接, 并且使用一次之后就断开`tcp`连接.

使用`keep-alive`可以改善这种状态, 即在一次`TCP`连接中可以持续发送多分数据而不会断开连接. 通过使用`keep-alive`机制, 可以减少`tcp`连接建立次数, 也就意味着可以减少`TIME_WAIT`状态连接, 以此提高性能和提高`httpd`服务器的吞吐率(更少的`tcp`连接意味着更少的系统内核调用, `socket`的`accept()`和`close()`调用).

它主要是用于客户端告诉服务端，这个连接我还会继续使用，在使用完之后不要关闭.

配置不当的`keep-alive`, 有时比重复利用连接带来的损失更大.

`Httpd`守护进程, 一般都提供了`keep-alive timeout`时间设置参数. 这个`keepalive timeout`时间意味着: ==一个`http`产生的`tcp`连接在传送完最后一个响应后, 还需要`hold`主`keepalive timeout`秒后, 才开始关闭这个链接==.

当`httpd`守护进程发送完一个响应后, 理应马上主动关闭响应的`tcp`链接, 设置`keepalive timeout`后, `httpd`守护进程回想说:"再等等吧, 看看浏览器还有没有请求过来", 这一等, 便是`keepalive timeout`时间.如果这个守护进程在这个等待的时间里, 一直没有收到浏览发过来`http`请求, 则关闭这个`http`链接.

### `TCP KEEPALIVE`

链接建立之后, 如果应程序或者上次协议一直不发送数据, 或者隔很长时间才发送一次数据, 当链接很久没有数据报文传输时如何去确定对方还在线, 到底是掉线了还是确实没有数据传输, 链接还需不需要保持.

`TCP`协议解决问题的方式:

1. 超过一段时间后, `TCP`自动发送一个数据为空的报文给对方.
2. 如果对方回应了这个报文, 说明对方还在线, 链接可以继续保持.
3. 如果对方没有报文返回, 并且重试了多次后则认为链接丢失, 没有必要保持链接.

### `http keep-alive` 和 `tcp keepalive`

两个不是同一回事, 意图不一样.

* `http keep-alive`是为了让`tcp`活的更久一点, 以便在同一个连接上传送多个`http`, 提高`socket`效率.
* `tcp keepalive`是`tcp`的一种**检测tcp链接状况的保险机制**

**`tcp keepalive`保险定时器**

* `/proc/sys/net/ipv4/tcp_keepalive_time`// 距离上次传送数据多少时间未收到判断为开始检测
* `/proc/sys/net/ipv4/tcp_keepalive_intvl`// 检测开始每多少时间发送心跳包
* `/proc/sys/net/ipv4/tcp_keepalive_probes`// 发送几次心跳包对方未响应则close连接

`keepalive`是`TCP`保鲜定时器, 当网络两端建立`tcp`链接之后, 闲置`idle`(双方没有任何数据流发送往来)了`tcp_keepalive_time`后, 服务器内核就会尝试向客户端发送侦测包, 来判断`tcp`链接状况(有可能客户端崩溃, 强制关闭了应用, 主机不可达等等). 如果没有收到对方的回答(`ack`包), 则会在`tcp_keepalive_intvl`后再次尝试发送侦测包, 知道收到对方的`ack`, 如果一直没有收到对方的`ack`, 一共会尝试`tcp_keepalive_probes`次, 每次的间隔时间在这里分别是`15s,30s,45s,60s,75s`. 如果尝试`tcp_keepalive_probes`依然没有收到对方的`ack`包, 则会丢弃该`tcp`链接. `tcp`链接默认限制时间是`2小时`, 一般设置为`30分钟`足够了.

```shell
[root@iZ944l0t308Z ~]# cat /proc/sys/net/ipv4/tcp_keepalive_time
7200 (2小时)
```

