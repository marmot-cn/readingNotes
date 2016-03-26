#Mac安装docker

---

####下载Docker Toolbox

[https://www.docker.com/products/docker-toolbox](https://www.docker.com/products/docker-toolbox "https://www.docker.com/products/docker-toolbox")


####启动

* `Docker Quickstart Terminal` 启动docker
* `Kitematic (Beta)` docker hub 管理工具
* `docker-compose` docker 编排工具,内置

####配置私有仓库

因为墙的原因,外网镜像比较难于访问,所以一般公司会内部构建一个私有仓库.我们首先还需要排至自己的客户端,配置方法如下:

		docker-machine ssh default
				                ##         .
		                  ## ## ##        ==
		               ## ## ## ## ##    ===
		           /"""""""""""""""""\___/ ===
		      ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~ /  ===- ~~~
		           \______ o           __/
		             \    \         __/
		              \____\_______/
		 _                 _   ____     _            _
		| |__   ___   ___ | |_|___ \ __| | ___   ___| | _____ _ __
		| '_ \ / _ \ / _ \| __| __) / _` |/ _ \ / __| |/ / _ \ '__|
		| |_) | (_) | (_) | |_ / __/ (_| | (_) | (__|   <  __/ |
		|_.__/ \___/ \___/ \__|_____\__,_|\___/ \___|_|\_\___|_|
		Boot2Docker version 1.9.1, build master : cef800b - Fri Nov 20 19:33:59 UTC 2015
		Docker version 1.9.1, build a34a1d5
		docker@default:~$ sudo vi /var/lib/boot2docker/profile
		#添加如下信息
		EXTRA_ARGS="--insecure-registry 120.25.87.35:5000"
		#重启docker服务
		sudo /etc/init.d/docker restart
		exit
		
好了现在就可以访问私有仓库了.

`120.25.87.35:5000`是我自己的私有仓库,后期可以根据公司的服务器配置自己的内部仓库地址.

####使用docker-compose编排镜像

**示例安装lnmp+memcached环境**

`docker-compose.yaml`:

		➜  lnmp  cat docker-compose.yml
		nginx:
		 image: "120.25.87.35:5000/nginx:1.9"
		 ports:
		  - "80:80"
		 links:
		  - "phpfpm"
		 volumes:
		  - ~/data/html:/usr/share/nginx/html
		  - ~/log/nginx:/var/log/nginx
		 container_name: web-nginx
		
		phpfpm:
		  image: "120.25.87.35:5000/php:5.6-fpm"
		  volumes:
		   - ~/data/html/:/var/www/html/
		  links:
		   - "memcached-1:memcached_1"
		   - "memcached-2:memcached_2"
		   - "mysql-master:dbw"
		  container_name: web-phpfpm
		
		mysql-master:
		  image: "120.25.87.35:5000/mysql:mac-5.6"
		  volumes:
		  - ~/data/mysql/:/var/lib/mysql
		  environment:
		   - MYSQL_ROOT_PASSWORD=123456
		  container_name: web-mysql-master
		
		memcached-1:
		  image: "120.25.87.35:5000/memcached:latest"
		  command: memcached -m 128
		  container_name: web-memcached1
		
		memcached-2:
		  image: "120.25.87.35:5000/memcached:latest"
		  command: memcached -m 64
		  container_name: web-memcached2
		
		phpmyadmin:
		  image: "120.25.87.35:5000/phpmyadmin:nazarpc-latest"
		  links:
		   - "mysql-master:mysql"
		  ports:
		   - "10081:80"
		  environment:
		   - UPLOAD_SIZE=1G
		  container_name: web-phpmyadmin
		  
以上安装了如下东西:

1. `php`: 5.6 编译了memcached 和 数据库 pdo 进去
2. `mysql`: 5.6 
3. `memcached` X 2
4. `phpmyadmin`: 检查内部数据库

在`phpfpm`使用了`links`,则可以直接使用别名`memcached_1`,`memcached_2`来访问`memcached`.使用`dbw`访问数据库.

整体对外暴漏`80`端口对外网站访问,`10081`端口为phpmyadmin.

**在mac环境容器挂载卷权限问题**

我在所有挂在卷前都指向自己的`home`目录下的文件夹.`~`符号.否则会有创建文件夹的权限问题.

**在mac环境的端口映射问题**

container使用的端口通过 docker -P或者 -p映射到了 VM里的 LinuxHost上,但是mac本机里是没有的.从本机可以用VM的ip访问到container.

比如我们通过`docker-compose up -d`启动了容器,对外映射了`80`和`10081`端口.虚拟机的IP地址为:

			docker is configured to use the default machine with IP 192.168.99.100
			For help getting started, check out the docs at https://docs.docker.com
			
我们通过`192.168.99.100`(默认`80`)端口可以访问到我们的网站目录,`192.168.99.100:10081`访问到`phpmyadmin`.

如果我们想在本机访问到容器的网站目录和`phpmyadmin`,我们可以在mac主机与vm虚拟机上的nat端口映射建立永久性的映射.

		VBoxManage controlvm "default"  natpf1 "tcp-port80,tcp,,8000,,80"
		
		VBoxManage controlvm "default"  natpf1 "tcp-port80,tcp,,10081,,10081"
		
我们把主机的端口`8000`映射到`80`,`10081`映射到`10081`.

现在我们通过主机地址`127.0.0.1:8000`,和`127.0.0.1:10081`访问网站目录和`phpmyadmin`.

**在mac环境使用mysql镜像的权限问题**

docker容器运行时,容器的当前用户mysql,由于权限不足无法mac本机上创建文件.
挂在卷容器内用户映射容器外用户权限问题.
我们重新编辑`mysql`的镜像:

		FROM 120.25.87.35:5000/mysql:5.6

		RUN usermod -u 1000 mysql \
		&& mkdir -p /var/run/mysqld \
		&& chmod -R 777 /var/run/mysqld
		
然后清空`~/data/mysql/`内的内容,这样运行mysql容器则一切正常.

