#FastCGI与mod_php详解

---

[参考链接](http://www.linuxeye.com/Linux/2824.html,"参考连接")

###背景

PHP最常用的方式:

* Apache中,以模块的方式(mod_php)运行.
* Nginx中,使用的是PHP-FPM.

###PHP处理器(PHP handlers)

首先,任何一种Web服务器(Apache、Nginx等)都是被设计成向用户发送html,图片等`静态资源`的,`Web服务器自身并不能解释任何动态脚本`(PHP、Python等).

**PHP处理器**

大多数的Web服务器都不能解析PHP代码,因此它需要一个能解析PHP代码的程序,这就是**`PHP处理器`**.

`PHP处理器`就是用来`解释Web应用中的PHP代码`,并`将它解释为HTML`或其他`静态资源`,然后`将解析的结果传给Web服务器`,最后再由Web服务器发送给用户.

**服务器与php处理器如何通信**

通过SAPI(Server Application Programming Interface 服务器端应用编程端口),简单来说,SAPI指的是PHP具体应用的编程接口.就像PC一样,无论安装哪些操作系统,只要满足了PC的接口规范都可以在PC上正常运行,PHP脚本要执行有很多种方式,通过Web服务器,或者直接在命令行下,也可以嵌入在其他程序中.

PHP最常用的SAPI提供的2种连接方法: **`mod_php和mod_fastcgi`**

###mod_php模式

**Apache是怎么能够识别php代码**

Apache的配置文件httpd.conf中加上或者修改:

		//添加
		LoadModule php5_module modules/libphp5.so
		AddType application/x-httpd-php .php
		//修改
		<IfModule dir_module>
    		DirectoryIndex index.php index.html index.htm  index.html
		</IfModule>

也即php作为Apache的一个`子模块`来运行,当通过web访问php文件时,Apache就会调用php5_module来解析php代码.

`配置加载mod_php模块后,php便是Apahce进程本身一部分,每个新的Apache子进程都会加载此模块`.

###mod_fastcgi模式

PHP-FPM - A simple and robust `FastCGI Process Manager for PHP`
PHP-FPM (FastCGI Process Manager) is an alternative PHP FastCGI implementation with some additional features useful for sites of any size, especially busier sites.

`PHP-FPM是一个PHP的FastCGI进程管理器`,解释的非常简单.这说明PHP-FPM是辅助mod_fastcgi模式进行工作的.

###什么是CGI？

`CGI(Common Gateway Interface) `是外部应用程序(CGI程序)与Web服务器之间的接口标准,是在CGI程序和Web服务器之间传递信息的规程.

`CGI规范允许Web服务器执行外部程序`,并`将它们的输出发送给Web浏览器`,CGI将Web的一组简单的静态超媒体文档变成一个完整的新的`交互式媒体`.

CGI是一种外部应用程序(CGI程序)与Web服务器的`协议`,CGI是`为了保证Server传递过来的数据是标准格式`.

###什么是FastCGI？

FastCGI像是一个`常驻(long-live)型的CGI`,它可以一直执行着,只要激活后,不会每次都要花费时间去fork一次(这是CGI最为人诟病的fork-and-execute模式).它还支持分布式的运算,即 FastCGI 程序可以在网站服务器以外的主机上执行并且接受来自其它网站服务器来的请求.

`FastCGI是语言无关的`,可伸缩架构的CGI开放扩展,其主要行为是`将CGI解释器进程保持在内存中并因此获得较高的性能`.

CGI解释器的`反复加载`是CGI`性能低下`的主要原因.如果CGI解释器保持在内存中并接受FastCGI进程管理器调度,则可以提供良好的性能、伸缩性、Fail- Over特性等等

FastCGI的整个工作流程是这样的:

1. Web Server启动时载入FastCGI进程管理器(IIS ISAPI或Apache Module)
2. FastCGI进程管理器自身初始化,启动多个`CGI解释器进程(可见多个php-cgi)`并等待WebServer的连接
3. 当客户端请求到达Web Server时,FastCGI进程管理器选择并连接到一个CGI解释器. Web server将CGI环境变量和标准输入发送到FastCGI子进程php-cgi
4. FastCGI子进程完成处理后将标准输出和错误信息从同一连接返回Web Server.当FastCGI子进程关闭连接时,请求便告处理完成,FastCGI子进程接着`等待`并处理来自FastCGI进程管理器(运行在Web Server中)的下一个连接,`在CGI模式中,php-cgi在此便已经退出`.

`FastCGI是CGI的升级版`,一种语言无关的协议,用来沟通程序(如PHP, Python, Java)和Web服务器(Apache2, Nginx),理论上任何语言编写的程序都可以通过FastCGI来提供Web服务.

FastCGI的特点是会在一个进程中依次完成多个请求,以达到提高效率的目的,大多数FastCGI实现都会维护一个进程池.

**总结**

FastCGI事先就需要启动,而且可以启动多个CGI模块,在那里一直运行等着web发请求,然后再给php解析运算,完成后生成html返回给web后,但是`完成后它不会退出`,而是`继续等着下一个web请求`.

###PHP-FPM

PHP-FPM就是针对于PHP的FastCGI的一种实现,他负责管理一个进程池,来处理来自Web服务器的请求.

PHP-FPM(FastCGI Process Manager)仅仅是个"PHP FastCGI 进程管理器".只用于PHP的.

###知识点整理

####Fail-Over

`失效转移`.

即当A无法为客户服务时,系统能够自动地切换,使B能够及时地顶上继续为客户提供服务,且客户感觉不到这个为他提供服务的对象已经更换.