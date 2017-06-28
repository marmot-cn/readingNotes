# PHP SOCKET 函数

---

### socket_create

创建一个套接字.

```php
resource socket_create ( int $domain , int $type , int $protocol )
```

创建并返回一个套接字,也称作一个通讯节点.一个典型的网络连接由 2 个套接字构成,一个运行在客户端,另一个运行在服务器端.

* domain 参数指定哪个协议用在当前套接字上
	* AF_INET: IPv4
	* AF_INET6: IPv6
	* AF_UNIX: 本地通讯协议
* type 选择套接字使用的类型
	* SOCK_STREAM(流式套接字): 提供一个顺序化的,可靠的,全双工的,基于连接的字节流.支持数据传送流量控制机制.TCP 协议即基于这种流式套接字.
	* SOCK_DGRAM(报式套接字): 提供数据报文的支持.(无连接，不可靠、固定最大长度).UDP协议即基于这种数据报文套接字.

### socket_accept

### socket_bind

### socket_clear_error

### socket_cmsg_space

### socket_connect

### socket_create_listen

### socket_create_pair

### socket_get_option

### socket_getopt

### socket_getpeername

### socket_getsockname

### socket_import_stream

### socket_last_error

### socket_listen

### socket_read

### socket_recv

### socket_recvfrom

### socket_recvmsg

### socket_select

### socket_send

### socket_sendmsg

### socket_sendto

### socket_set_block

### socket_set_nonblock

### socket_set_option

### socket_setopt

### socket_shutdown

### socket_strerror

### socket_write