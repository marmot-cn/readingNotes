# README

## 第一章 初识Nginx

### NGINX 目录

* auto
	* cc: 编译
	* os: 操作系统判断
* conf, 示例配置文件
* contrib/vim, vim 配置文件
* html
	* 50x.html
	* index.html
* man, linux 帮助文件
* src, 源代码
* configure
	* --xx=PATH, 设定寻找目录
	* --with-xx, 使用哪些模块（默认不会编译进nginx）
	* --without-xx，不使用哪些模块（默认编译进nginx）
* objs, 编译时的中间文件
	* ngx_modules.c, 编译时模块
	* nginx, (make后生成的二进制文件)

### NGINX 常用信号

* TERM, INT: Quick shutdown
* QUIT: Graceful shutdown  平缓停止Nginx服务
* HUP: Configuration reload ,Start the new worker processes with a new configuration Gracefully shutdown the old worker processes. 改变配置文件,平滑的重读配置文件
* USR1: Reopen the log files 重读日志,在日志按月/日分割时有用
* USR2: Upgrade Executable on the fly 平滑的升级
* WINCH: Gracefully shutdown of worker processes 平缓停止worker进程，用于Nginx服务器平滑升级

```
nginx -g WINCH
```

### NGINX 静态资源

`$limit_rate ?`如`set $limit_rate 1k;`, 限制访问速率

`log_format`日志格式

### NGINX 缓存

* `proxy_cache_path`
	* 缓存路径
	* 命名

### 工具

`go access`查看`nginx`的`access`日志

### 证书

* DV: 域名验证
* OV: 组织验证, 申请证书验证组织机构
* EV: 扩展验证, 更严格的验证

NGINX向浏览器发送两个证书, 二级与一级证书。根证书操作系统内部自己验证。

![](img/20210928102323.jpg)

### TLS通讯过程

1. 验证身份
2. 达成安全套件共识
3. 传递密钥
4. 加密通讯

## 第二章 Nginx架构基础

### Nginx的请求处理流程

worker进程数量需要与CPU核数匹配

### Nginx的进程结构

* 单进程结构: 调试
* 多进程结构
	* 父进程: master-process，worker进程管理
	* 子进程
		* cache, 缓存在多个worker进程之间共享
			* cache manager 进程
			* cache loader 进程
		* worker, 处理真正请求，多个worker进程，希望每个worker进程占有一个CPU,

worker进程数与CPU核心数匹配，并且绑定CPU, 可以最大限度防止缓存失效。

进程间通信使用共享内存管理。

为什么多进程结构, 需要高性能，高可靠性。因为线程之间共享同一个地址空间，一个线程段错误，整个进程挂掉。

### Nginx reload 流程

1. 向`master`进程发送`HUP`信号(`reload`命令)
2. `master`进程校验配置语法是否正确
3. `master`进程打开新的监听端口(可能会打开新的端口，子进程会集成父进程端口)
4. `master`进程用新配置启动新的`worker`子进程
5. `master`进程向老`worker`子进程发送`QUIT`信号
6. 老`worker`进程关闭监听句柄，处理完当前连接后结束进程

### Nginx 热升级流程

1. 将旧`Nginx`文件换成新`Nginx`文件（注意备份）
2. 向`master`进程发送`USR2`信号
3. `master`进程修改`pid`文件名，加后缀`.oldbin`
4. `master`进程用新`Nginx`文件启动新`master`进程
5. 向老`master`进程发送`QUIT`信号，关闭老`master`进程
6. 回滚: 向老`master`发送`HUP`, 向新`master`发送`QUIT`

### `worker`进程：优雅的关闭

1. 设置定时器`worker_shutdown_timeout`
2. 关闭监听句柄
3. 关闭空闲连接
4. 在循环中等待全部连接关闭, 超过`worker_shutdown_timeout`则强制关闭
5. 退出进程

### Nginx 一线程同时处理多连接

* 用户态代码完成连接切换
* 建立减少OS进程切换

### Nginx模块

`nginx`编译后模块可通过在`objs/ngx_modules.c`中检查是否编译某个模块。

### Nginx模块的分类

* `NGX_CORE_MODULE`: 核心模块, 也会衍生定义子类型模块
* `NGX_CONF_MODULE`: 解析配置文件
* `NGX_EVENT_MODULE`: 事件模块

### Nginx 连接池

### Nginx 内存池

内存池会提前分配内存，产生内存碎片小。

链接内存，保存上下文少。`connection_pool_size`

对于请求内存，保存上下文较多。`request_pool_size`

### Nginx 进程间的通信方式

* 基础同步工具
	* 信号
	* 共享内存
* 高级通讯方式
	* 锁，多个内存共同操作一个内存，旧的操作系统使用信号量，目前操作系统使用自旋锁。使用自旋锁，要求`nginx`的`worker`必须快速使用内存。
	* Slab内存管理器

### slab 内存管理

