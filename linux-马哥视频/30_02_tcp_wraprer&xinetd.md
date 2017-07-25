# 30_02_tcp_wraprer&xinetd

---

### 笔记

---

#### xinetd

配置文件`/etc/xinetd.conf` = 主配置文件 + 各种配置文件片段.

各种配置文件片段: `/etc/xinetd.d/*`.

```shell
[root@iZ944l0t308Z ~]# cat /etc/xinetd.conf
#
# This is the master xinetd configuration file. Settings in the
# default section will be inherited by all service configurations
# unless explicitly overridden in the service configuration. See
# xinetd.conf in the man pages for a more detailed explanation of
# these attributes.

defaults
{
# The next two items are intended to be a quick access place to
# temporarily enable or disable services.
#
#	enabled		=
#	disabled	=

# Define general logging characteristics.
	log_type	= SYSLOG daemon info
	log_on_failure	= HOST
	log_on_success	= PID HOST DURATION EXIT
...

[root@iZ944l0t308Z ~]# ls /etc/xinetd.d/
chargen-dgram  chargen-stream  daytime-dgram  daytime-stream  discard-dgram  discard-stream  echo-dgram  echo-stream  tcpmux-server  time-dgram  time-stream
```

##### 配置文件

两部分组成:

* 全局配置(服务的默认配置): 对`/etc/xinetd.d/*`都生效的.`defaults`段. 具体服务自己没有定义, 从此继承.
* 服务配置段

		service <service_name>
		{
			<attribute> <assign_op> <value> <value> ...
		}
		
`assign_op`:

* `=`: 覆盖默认选项.
* `+=`: 在默认选项新增.
* `-=`: 在默认选项移除.

##### 配合文件参数

```shell
man xinetd.conf 查看配置文件说明
```

* `log_type`: 日志类型
	* `SYSLOG`, 通过`SYSLOG`系统服务来记录日志(syslogd:系统日志, klogd:内核日志). 日志通常记录在`/var/log/messages`中.
	* `FILE`
* `log_on_failure`: 当访问失败记录日志格式 
* `log_on_sucess`: 当访问成功记录日志格式

`chkconfig`改变瞬时守护进程就是修改配置文件的`disable`参数的值.

**套接字类型**

* tcp: stream
* udp: dgram
* rpc(`portmap`提供`rpc`服务)

**访问控制**

* `only_from = `: 进允许哪些客户机, 访问控制法则.
	* `IP`: 172.16.100.200
	* `NETWORK`: 172.26.0.0/16, 172.16.0.0/255.255.0.0
	* `HOSTNAME`: FQDN www.xxx.com
	* `DOMAIN`: *.xxx.com
* `on_access = `: 拒绝谁访问. 参数和`only_from`一样.

**时间控制**

* `access_time = hh:mm-hh:mm`: 哪个时间可以访问.
	* `hh`: 0-23
	* `mm`: 0-59

**监听的地址(提供服务的地址)**

* `bind = `: 绑定在哪个`ip`上.
* `interface = `: 和`bind`一样.

		本机有多个ip地址, 只允许一个ip地址可以访问.
		bind = 172.16.100.1 只监听该地址.
		
		现在 netstat 查看会从 0.0.0.0 监听地址变为 172.16.100.1
		
**资源访问控制**

服务所能接收的连接数(最多只允许多少个用户同时连进来)

* `cps = `(connection per second): 
	* 第一个参数: ==每秒钟==允许链接进来最多的请求个数, 如果超了一段时间服务会禁用一段时间.
	* 第二个参数: 临时禁用多长时间.
			
			只允许一个链接进来, 后续在进来需要等待10s
			cps = 1 10 
* `per_source = `某一个特定`ip`最多发起几个链接请求
		
		每一个单独的ip只允许发起一个链接请求
		per_source = 1

* `instances = `: 某个特定服务最大同时链接数.
			
**向启动的额server传递参数**

* `server_args = `

#### 示例

设定本地的`rsync`服务(非独立守护进程), 满足如下需求:

1. 仅监听本地172.16.x.1的地址上提供服务
2. 仅允许172.16.0.0/16网络内的主机访问; 但不允许172.16.0.1访问.
3. 仅允许同时运行最多3个实例, 而且每个ip最多只允许发起两个连接请求.

```shell
/etc/xinetd.d/rsync

serivce rsync
{
	disable = no
	socket_type = stream
	wait = no
	user = root
	server = /user/bin/rsync
	server_args = --daemon
	log_on_failure += USERID
	only_from = 172.16.0.0/16
	no_access = 172.16.0.1
	bind	= 172.16.100.1
	instances = 3
	psr_source = 2
}
```
		
### 整理知识点

--- 