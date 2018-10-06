# 43_03_配置Ngnix作为Web Server详解

---

## 笔记

### nginx

* File AIO, 文件异步IO
* Asynchronous
* Event-driven edge trigger
* 模块化设计
* 重新加载配置以及在线升级时, 不需要中断正在处理的请求(热部署).
	* 一个线程(轻量级进程)响应多个请求
		* 旧的链接使用旧的配置
		* 新的链接使用新的配置
* 日志
	* 自定义访问日志格式
	* 带缓存的日志写操作
		* 将日志通过缓冲区, 一段时间在flush
	* 快速日志轮转
* 重写(rwite)模块, 使用正则表达式改变URI.
* 支持验证`HTTP_referer`
* 基于客户端IP地址和HTTP基本认证机制的访问控制.
* 来自同一地址的同时连接数和请求数限制.
	* 速度限制

用处:

* web服务器
* 轻量级反向代理
	* web
	* mail
* 图像缩放(节约网络带宽)
* 缓存
	* disk
	* 缓存文件描述符
	* memcached(模块化支持)

#### 架构和扩展性

一个主进程和多个工作进程, 工作进程以非特权用户运行.

`nginx`服务器启动时内部会生成一个`master`主进程. 

* `master`会启动多个`worker`线程(轻量级进程), 以普通用户启动, 处理用户请求. `master`监控`worker`进程.
* `cacheloader`进程, 装载, 清理缓存.


```
root@3fcd1647825d:/var/www/html# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 Sep30 ?        00:00:00 nginx: master process nginx -g daemon off;
nginx        9     1  0 Sep30 ?        00:00:00 nginx: worker process
nginx       10     1  0 Sep30 ?        00:00:00 nginx: worker process
nginx       11     1  0 Sep30 ?        00:00:00 nginx: worker process
nginx       12     1  0 Sep30 ?        00:00:00 nginx: worker process
root        13     0  0 22:21 pts/0    00:00:00 /bin/bash
root       270    13  0 22:27 pts/0    00:00:00 ps -ef

root@3fcd1647825d:/var/www/html# netstat -ntlp
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
tcp        0      0 0.0.0.0:80              0.0.0.0:*               LISTEN      1/nginx: master pro
```

`matser`是管理员权限启动, 监听`80`端口. 只有管理员才有权限监听小于`1023`的端口. 负责状态主配置文件. 装载成功后, 新连接装载配置文件.

`worker`进程用到哪个模块, 装载哪个模块. 除了基本功能都转交给模块处理, 模块以流水线方式工作.

#### master进程

* 读取并验证配置信息
* 创建, 绑定及关闭套接字
* 启动, 终止及维护`worker`进程的个数
* 无需终止服务而重新配置工作的特性.
* 控制非中断式程序升级(平滑升级), 启用新的二进制程序并在需要时回滚至老版本.
* 重新打开日志文件, 实现日志滚动.
* 编译嵌入式`perl`脚本

#### worker进程

* 接收, 传入并处理来自客户端的链接
* 提供反向代理及过滤功能
* `nginx`任何能完成的其它任务.

#### cache loader进程

* 检查缓存存储中的缓存对象
* 使用缓存元数据建立内存数据库
	* 文件句柄
	* ....

#### cache manager进程

* 缓存的失效即过期检验

#### 配置上下文

配置分段

* main: 核心配置, 对任何功能都生效配置
* http: 对web服务器有效
	* server: 必须属于`http`或`mail`中
	* upstream: 定义反向代理
	* location: 
* mail: 对邮件服务器有效

#### sendfile

请求流程

* worker进程工作在用户空间
* 用户请求静态页面, 页面在磁盘某个分区某个文件系统上面
* 用户请求进来, 到网卡, 内核处理交给监听在80空间的应用程序(worker进程) **内核空间**
* worker发现用户请求静态页面, 向内核发出系统调用(IO)
* 内核准备一个缓冲
* 内核从磁盘加载文件到缓冲当中
* 将缓冲内容复制到进程自己的地址空间
* worker将程序封装成响应报文
	* http请求首部(真正的封装使在内核中完成)
* 将数据发给内核, 有内核封装tcp首部, IP首部, mac首部
* 响应给客户端


`sendfile`避免两次复制, 在内核中读过来就响应

* 第一次从内核空间复制到用户空间
* 最后输出的时候从用户空间复制到内核空间

#### Accept-filters

链接过滤器

#### 10000个非活跃的HTTP keep-alive链接仅占用约2.5M内存

事件驱动只扫描活跃链接, 非活跃链接不做管理.

`nginx`自身只需要很小内存为每一个链接维持一个 文件描述符.

#### 内核和进程交互

**复制**, 除非**共享内存**

### 支持事件驱动的`I/O`框架

* epoll(Linux 2.6+)
* kqueue(FreeBSD 4.1+)
* /dev/poll(Solaris 7 11/99+)

### 配置文件

#### worker个数设定

* 如果负载以`CPU`密集型应用为主, 如`SSL`或压缩应用, 则`worker`数应与`cpu`数相同. 
	* 本地应用占据更多的`CPU`时间处理.
* 如果负载以`IO`密集型为主, 如响应大量内容给客户端, 则`worker`数应该为`cpu`个数的`1.5`倍或`2`倍.

#### events

```
events {
	worker_connections 1024;
}
```

事件驱动中每个`worker`支持的链接数.

最大的链接数`worker_connections`*`worker_processes`.

#### server

每个`server`段定义一个虚拟主机.

#### location

基于`URI`路径定义访问属性.

```
http://www.abc.com/从这里开始是URI路径
```

```
location /URI/ {
   root "本地文件目录";
   index 默认主页面;
}
```

* `root`定义对于`URI`映射在本地文件系统的路径.

```
# 对当前路径及子路径下的所有对象都生效
location URI {}; 

# 精确匹配指定路径, 只对当前资源生效, 不包括子路径
location = URI {}; 

# 模式匹配 区分字符大小写, URI可以使用正则表达式做通配
location ~ URI {};

# 模式匹配 不区分字符大小写, URI可以使用正则表达式做通配
location ~* URI {};

# 不使用正则表达式
location ^~ URI {};
```

* 优先级由高到低
	* `=`精确匹配
	* `^~`不适用正则表达式
	* `正则表达式`
	* 不适用任何符号



#### error_page

```
error_page 500 502 503 504 /50x.html;

location = /50x.html {
      root html;#相对路径, 相对于nginx的默认安装目录, prefix指定
}
```

如果返回错误代码是`500`,`502`,`503`,`504` 则读取根下`50x.html`

### 访问控制

#### 基于IP访问控制

**deny(拒绝)** 和 **allow(允许)**

```
location / {
  deny all;
  allow xxx.xx.xx.xx/xx;
}
```

#### 基于用户访问控制

Auth Basic

```
lication / {
  auth_basic "xxx"; #服务器用于提示
  auth_basic_user_file htpasswd;
}

```

使用`apache`的`htpasswd`命令, 创建`htpasswd`文件.

`htpasswd`: `-c`创建文件.

### LNMP

#### FastCGI

`Nginx`不支持模块化整合`php`.

`php-fpm`: 默认监听`9000`端口. 

#### 参数

**fastcgi_pass** fastcgi模式的反向代理.

**fastcgi_index** fastcgi的主页面.

**fastcgi_param** fastcgi额外参数. 当做`fastcgi`的参数, 传递给`fpm`服务器.

```
fastcgi_param SCRIPT_FILENAME xxxx$fastcgi_script_name;

#SCRIPT_FILENAME 脚本文件名称
```

## 整理知识点

---

### specs

