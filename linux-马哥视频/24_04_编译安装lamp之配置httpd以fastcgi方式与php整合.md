# 24_04_编译安装LAMP之配置httpd以FastCGI方式与php整合

---

### 笔记

---

#### httpd 和 php 结合的方式

* cgi
* module
* fastcgi(fpm)

httpd 必须要提供 fastcgi 模块. 2.4以后就自带了,fcgi.

`--enable-modules=most`,`--enable-mpms-shared`

#### phpfpm 模式安装php

**配置时**

		--with-apxs2=/usr/local/apache/bin/apxs 改为 --enable-fpm
		
**为php提供配置文件**

		cp php.ini-production /etc/php.ini
		
**配置php-fpm**

为php-fpm提供Sysv init脚本,并将其添加至服务列表
		
		cp sapi/fom/init.d.php-fpm	/etc/rc.d/init.d/php-fpm
		chmod +x /etc/rc.d/init.d/php-fpm
		
		加入到服务列表
		chkconfig --add php-fpm
		chkconfig php-fpm on

**为php-fpm提供配置文件**

		cp xxx/php-fpm.conf.default /usr/local/php/etc/php-fpm.conf

**编辑php-fpm的配置文件**

		vim /usr/local/php/etc/php-fpm.conf
		配置fpm的相关选项为你所需要的值,并启用pid文件(最后一行)
		
		#最多有多少个子进程
		pm.max_children = 50
		#刚开始启动多少个空闲进程
		pm.start_servers = 5
		#最少有个几个空闲进程
		pm.min_spare_severs = 2
		#最多有几个空闲进程
		pm.max_spare_severs = 8
		#因为/etc/rc.d/init.d/php-fpm这里默认是在下面路径寻找pid文件
		pid = /usr/local/php/var/run/php-fpm.pid 

**启动php-fpm**

		service php-fpm start
		
**验证php-fpm是否启动**

		ps aux | grep php-fpm
		
**验证监听端口**

		netstat -tnpl | grep php-fpm
		tcp 	0	0	127.0.0.1:9000	0.0.0.0:*	LISTEN	689/php-fpm

#### 配置httpd-2.4.4

**启用httpd的相关模块**

在Apache httpd 2.4以后已经专门又一个模块针对FastCGI的实现,此模块为mod_proxy_fcgi.so,它其实是作为mod_proxy.so模块的补充.因此,这两个模块都要加载

* LoadModule proxy_module modules/mod_proxy.so
* LoadModule proxy_fcgi_module modules/mod_proxy_fcgi.so

**AddType**

		AddType application/x-httpd-php .php
		AddType application/x-httpd-php-source .phps
		
		#主页面支持 index.php
		<IfModule dir_module>
			DirectoryIndex index.php index.html
		</IfModule>

**配置虚拟主机支持使用fcgi**

在相应的虚拟主机中添加类似如下两行

		关闭正向代理
		ProxyRequests Off 
		
		当用请求一个uri,转换到另外一个主机,以/开头中间跟任意字符以.php结尾
		反向代理到 fcgi 协议
		$1 代表 \1,引用第一个括号中表示的内容,有第2个括号用$2表示,这里表示请求哪个php文件都反向代理到对应地址的php文件
		把以.php结尾的文件请求发送到php-fpm进程,php-fpm进程至少需要知道运行的目录和URI,所以这里直接在fcig://127.0.0.1:9000后面指明了这两个参数,其它的参数ude传递已经被mod_proxy_fcgi.so进行了封装,不需要手动指定
		
		ProxyPassMatch ^/(.*\.php)$ fcgi://127.0.0.1:9000/PATH/TO/DOCUMENT_ROOT/$1
		
例如:
		
		<VirtualHost *:80>
			DocumentRoot "/www/magedu.com"
			ServerName magedu.com
			ServerAlias www.magedu.com
			
			ProxyRequests Off 
			ProxyPassMatch ^/(.*\.php)$ fcgi://127.0.0.1:9000/www/magedu.com/$1
		
				<Directory "/www/magedu.com">
						Options none
						AllowOverride none
						Require all granted
				</Directory>
			
		</VirtualHost>

### 整理知识点

---

#### 反向代理

当客户端请求一个内容时,服务器自身没有,服务器到另外一台服务器获取内容,先缓存到本地在返回给客户端.

**ProxyPass**

		ProxyPass /images/a.jpg  http://172.16.100.2/images/a.jpg
		
		当前主机并没有 /images 这个目录和 a.jpg 这个文件,到对应的主机获取后返回给客户端
		
**ProxyPassMatch**

支持了正则
	
		