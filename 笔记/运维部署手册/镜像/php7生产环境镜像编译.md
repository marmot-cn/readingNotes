#php7生产环境镜像编译

---

####编译过程

选择从官方下载`php:7-fpm`镜像作为原始镜像.因为镜像已经编译进去了`PDO`所以我们不在安装此扩展.

我选择了编译`memcached`和`redis`扩展.

* `memcached`: 选作用`row cache`和`page cache`
* `redis`: 暂时选作用`vector(关系型)`缓存

因为`pecl`暂未发布支持php7的最新版本,我们需要手动编译.

* `memcached`来源: https://github.com/phpredis/phpredis/ 
* `redis`来源: https://github.com/rlerdorf/php-memcached

**redis准备工作**

		git clone https://github.com/phpredis/phpredis/ 
		cd phpredis 
		git checkout php7 
		git archive --format=tar.gz -o php7redis.tar.gz HEAD
		
导出到压缩包,为我们编译镜像时使用

**memcached准备工作**

需要先下载`libmemecached`库才能正常编译.

		wget https://launchpad.net/libmemcached/1.0/1.0.18/+download/libmemcached-1.0.18.tar.gz 
		
导出`memcached`源码包

		git clone https://github.com/rlerdorf/php-memcached.git 
		cd php-memcached 
		git checkout php7 
		git archive --format=tar.gz -o php7memcached.tar.gz HEAD
		
**apcu**

为PHP框架内的容器添加apcu缓存,而不是让其依赖于memcached.	
		
		git clone https://github.com/krakjoe/apcu.git
		cd apcu
		git archive --format=tar.gz -o php7apcu.tar.gz HEAD
		
####Dockerfile编写

我编译了2个版本的镜像,一个用于生产环境,一个用于开发环境.但是来源镜像都是基于官方的php:7-fpm镜像.生产环境单独优化了开启了`hugePage`.

编译`memcached`需要`zlib1g-dev`.我在镜像里面编译了下载.同时为了优化php7,开启了`opcache`.

**开发环境Dockerfile**

		# version: v1.0.20160206

		FROM 120.25.87.35:5000/php:7-fpm-original
		ADD ./libmemcached-1.0.18.tar.gz /data/php7extension/libmemcached
		ADD ./php7redis.tar.gz /data/php7extension/redis
		ADD ./php7memcached.tar.gz /data/php7extension/memcached
		
		RUN apt-get update && apt-get install zlib1g-dev \
		&& cd /data/php7extension/libmemcached/libmemcached-1.0.18/ \
		&& ./configure \
		&& make \
		&& make install \
		&& cd /data/php7extension/memcached/ \
		&& phpize \
		&& ./configure \
			 --disable-memcached-sasl \
		&& make \
		&& make install \
		&& echo "extension=memcached.so" > /usr/local/etc/php/conf.d/memcached.ini \
		&& cd /data/php7extension/redis/ \
		&& phpize \
		&& ./configure \
		&& make \
		&& make install \
		&& echo "extension=redis.so" > /usr/local/etc/php/conf.d/redis.ini \
		&& set -ex \
		&& { \
		        echo 'zend_extension=opcache.so'; \
		        echo 'opcache.enable=1'; \
		        echo 'opcache.enable_cli=1'; \
		} | tee /usr/local/etc/php/conf.d/opcache.ini

**生产环境Dockerfile**

		# version: v1.0.20160206

		FROM 120.25.87.35:5000/php:7-fpm-original
		ADD ./libmemcached-1.0.18.tar.gz /data/php7extension/libmemcached
		ADD ./php7redis.tar.gz /data/php7extension/redis
		ADD ./php7memcached.tar.gz /data/php7extension/memcached
		
		RUN apt-get update && apt-get install zlib1g-dev \
		&& cd /data/php7extension/libmemcached/libmemcached-1.0.18/ \
		&& ./configure \
		&& make \
		&& make install \
		&& cd /data/php7extension/memcached/ \
		&& phpize \
		&& ./configure \
			 --disable-memcached-sasl \
		&& make \
		&& make install \
		&& echo "extension=memcached.so" > /usr/local/etc/php/conf.d/memcached.ini \
		&& cd /data/php7extension/redis/ \
		&& phpize \
		&& ./configure \
		&& make \
		&& make install \
		&& echo "extension=redis.so" > /usr/local/etc/php/conf.d/redis.ini \
		&& set -ex \
		&& { \
		        echo 'zend_extension=opcache.so'; \
		        echo 'opcache.enable=1'; \
		        echo 'opcache.enable_cli=1'; \
				echo 'opcache.huge_code_pages=1'; \
		} | tee /usr/local/etc/php/conf.d/opcache.ini
		
额外多了一个`opcache.huge_code_pages=1`.开启此功能需要现在服务器开启`hugepage`
		
		[chloroplast@iZ235s6fe4eZ lnmp-production]$ sudo sysctl vm.nr_hugepages=128
		[chloroplast@iZ235s6fe4eZ lnmp-production]$ cat /proc/meminfo | grep Huge
		AnonHugePages:     32768 kB
		HugePages_Total:     128
		HugePages_Free:      128
		HugePages_Rsvd:        0
		HugePages_Surp:        0
		Hugepagesize:       2048 kB
		
当开启服务后,我们验证确实使用了HugePage

		[chloroplast@iZ235s6fe4eZ php7-production]$ cat /proc/meminfo | grep Huge
		AnonHugePages:    491520 kB
		HugePages_Total:     128
		HugePages_Free:      116
		HugePages_Rsvd:       27
		HugePages_Surp:        0
		Hugepagesize:       2048 kB
		
####HugePage理解
		
				


