# 24_03_编译安装LAMP之php-5.4.13、xcache-2.0及使用ab命令实现压力测试

---

### 笔记

---

#### apache 结合 php

		AddType application/x-httpd-php .php
		AddType application/x-httped-php-source .phps

`.phps`是源码文件.

		<IfModule dir_module>
		  DirectoryIndex index.php index.html
		</ifModule>

#### xcahe

php 支持扩展功能:

* xcache
* apc
* ...

通过PHP相关命令,首先需要PHP识别并加载该扩展

		tar xf xcache-2.0.0.tar.gz
		cd xcache-2.0.0
		/usr/local/php/bin/phpize
		./configure --enable-xcache --with-php-config=/usr/local/php/bin/php-config
		make && make install
		
**phpize**

prepare a PHP-extension for compiling.

**php-config**

get information about PHP configuration and compile options.

获取 php 的配置信息 和 编译信息.

#### ab

`A`pache `B`enchmark

压力测试.

**参数**

* `-c`: 并发量(一次发起的请求数叫并发数)
* `-n`: 总共发送多少个请求

		n > c 数字
		
		c 每批多少个请求
		n 一共发出多少个请求

当`-c`过于大的时候会超过限制,linux 限制每个进程不能打开超过`1024`,可以使用`ulimit -n`查看.
也可以使用`ulimit -n #`来修改.

		ulimit -n 10000		
		
**返回结果**

* `Total transferred`: 报文和报文首部
* `HTML transferred`: HTML 数据传输大小,不包含首部

#### 编译安装apache 启用 https

		vi httpd.conf
		启用 mod_ssl.so (loadModule)
		include /etc/httpd/extra/httpd-ssl.conf
		
		编辑 httpd-ssl.conf
		修改 virtua-host
		提供证书
		
		ssl 是基于IP地址,一个IP地址只能建立一个ssl,但是SNI可以支持一个IP绑定多个ssl证书

#### 其他配置文件

**mpm.conf**

mpm配置文件

* prefork
* worker
* event

### 整理知识点

---

#### SNI

Server Name Indication

#### 压力测试工具

* http_load
* siege
* webbench
* ab

#### Mbps 和 MB/s

Mbps和MB/s表示的网络传输速度不一样,MB/S是Mbps的8倍.

`B`yte = 8*`b`it

Mbps和MB/s表示的网络传输速度不一样，MB/S是Mbps的8倍。

Mbps=Mbit/s即兆比特每秒(Million bits per second的缩写).bps（bits per second,即比特率,比特/秒,位/秒,每秒传送位数,数据传输速率的常用单位.而通常所说的文件大小的兆是指8MByte.字节(Byte)是计算机信息技术用于计量存储容量的一种计量单位,也表示一些计算机编程语言中的数据类型和语言字符.
字节和比特的换算关系如下:

		1Byte＝8bit
		1KByte=1024Byte
		1M=1024KByte
		1MB/s=8Mbps