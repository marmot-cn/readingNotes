#21_02_httpd安装配置

---

###笔记

---

**apache**

A Patchy Server, 充满补丁的服务.

软件基金会 `ASF`: `A`pache `S`oftware `F`oundation.

`FSF`: GNU, GPL

一般用`apache`称呼`httpd`.

####httpd

Web Server, 开源.

是受SELinux控制的.

		getenforce 查看 selinux状态
		
		setenforce 关闭 selinux

#####httpd特性

**事先创建进程**

进程控制上事先创建进程,临时创建进程响应用户请求比较慢.事先把进程创建好,做为空闲进程.如果空闲进程不够,创建新的空闲进程.

**按需维护适当的进程**

**模块化设计,核心比较小,各种功能都通过模块添加,模块可以在运行时启用**

类似`Linux`内核,平时只启用核心功能,需要额外功能装载额外模块.

`PHP`也是以模块化模式整合进去.

需要用到额外的功能,单独编译模块即可.`运行时配置`

**支持多种方式虚拟主机配置**

**支持https协议**

使用`mod_ssl`模块来实现.

**支持用户认证**

**支持访问控制机制**

* 基于IP或主机名的`ACL`.
* 支持每目录的访问控制.

**支持URL重写**

对客户端是透明的.伪静态.

####nginx

多进程响应,N个用户请求.一般nginx用作反向代理.

####虚拟主机

* 物理服务器只有一个
* Web程序只有一个,却可以服务多个不同的站点,`www.a.com`,`www.b.com`会访问不同的站点.

		客户端IP,端口都是一个.独立的套接字(IP+端口=套接字).
		
**虚拟主机方式**

* 基于IP的虚拟主机,不同的IP对应不同的主机.但是IP是紧缺资源.
* 基于端口的虚拟主机,不同的端口对应不同的主机.不适用标准`80`端口,访问较麻烦.
* 基于域名的虚拟主机: IP 和 端口 一样, DNS 主机名不同,打开不同的虚拟主机.

**基于域名的虚拟主机实现方式**

客户端请求报文:

		Method URL version
		header
		
		body

标准的请求方式: `protocol://HOST:PORT/path/to/source`

		GET /download/linux.tar.bz2  HTTP/1.0
		Host: www.magedu.com
		
		组合起来就是标准的请求防范.

`Host`是一个域名,这样服务器可以基于`Host`来区分.这也是为什么请求报文中必须包含`Host`.如果请求时候是基于`ip`地址,只能`返回默认的虚拟主机`.

####Httpd安装

* rpm包
* 源码编译


		httpd -l, 列出当前服务器编译支持的模块
		Compiled in modules:
		  core.c 核心模块
		  mod_so.c 动态模块加载
		  http_core.c http核心
		  prefork.c
  
**httpd多道处理模式**

`MPM(Multi Path Modules)`:`prefork`, 默认. 事先启动多个进程.

* `mpm_winnt`, windows 专用的
* `prefork`,一个请求用一个进程响应.
* `worker`, 基于线程的,一个请求用一个线程响应,(web服务器启动多个进程,每个进程生成多个线程).
	* 多个线程共享一个进程的资源,如果一个线程打开过一个资源,第2个线程不用再次打开,但是在同时写的时候,为了避免资源竞争,需要加锁.因为linux默认不支持线程,效率不如`prefork`,所以默认使用`prefork`.
* `event`, 基于事件的驱动,一个进程直接处理多个请求.`Nginx`默认是`event`模型.

执行程序是`httpd`,启动后由多个`httpd`进程,在众多进程当中有一个`httpd`进程的属主和属组是`root`.其他进程的属主和属组是`apache`,`apache`.

1. `httpd`: `root`,`root` 主导进程 `master process`
2. `httpd`: `apache`,`apache` 工作进程 `work process`,真正处理用户请求

在`linux`使用小于1024的端口必须有管理员程序,所以启动这个程序的必须是`root`.第一个进程不负责响应客户请求,值负责创建和回收空闲进程.

**占据端口**

`80/tcp`, `ssl: 443/tcp`

**目录**

* `/etc/httpd`: 工作根目录,相当于程序安装目录.
* `/etc/httpd/conf`: 配置文件目录.
	* `httpd.conf`: 主配置文件.
	* `/etc/httpd/conf.d/*.conf`: 主配置文件使用`include`包含其他的配置文件,组合起来使用.
* `/etc/httpd/modules`: 模块目录.
* `/etc/httpd/logs` --> `/var/logs/httpd`: 日志目录.
	* 访问日志,`access_log`,发起请求和响应的结果.
	* 错误日志,`err_log`,错误信息.
* `/var/www`
	* `/var/www/html`,存放静态文件.
	* `/var/www/cgi-bin`,存放动态文件.

**如何处理动态内容**

web服务器不处理任何动态内容.

通过某种`协议`,调用额外的其他程序,来运行这个程序并将结果取回后响应给客户端.

####cgi

让web服务器可以和额外的应用程序通信的一种机制.必要的时候启动一个额外的程序来处理动态内容.

`C`ommon `G`ateway `I`nterface : 通用网关接口,是一种协议.

client --> httpd (index.cgi) --> 发起一个进程(和cgi程序语言相关的进程)`Spawn Process` --> 处理响应给httpd --> 返回给client

**fastcgi**

`cgi`每次都是`web`创建和销毁,`fascgi`专门又一个`master`进程,预先创建好`work`进程.并且`work`进程由`master`进程负责创建和销毁.

专门的进程工作在专门的`socket`上来和`web`服务器进程通信.

可以将`web`服务器和`动态处理`分开.所处处理动态内容的服务又叫做`应用程序服务器`.

####程序

程序是由`指令`和`数据`组成的.

####CPU-bound

cpu密集型,对CPU占用高.


###整理知识点

---