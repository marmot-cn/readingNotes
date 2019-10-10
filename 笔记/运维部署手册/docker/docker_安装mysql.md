#docker 安装mysql

####安装mysql从Dockerfile

**编写Dockerfile**

		# version: v1.0.20151209

		FROM mysql:5.6
		MAINTAINER chloroplast "41893204@qq.com"
		
拉取mysql官方镜像,5.6版本. mysql:5.6

**创建镜像**

		sudo docker build -t chloroplast/mysql:v1.0.20151209 ./
		
`./`指的是Dockerfile的路径,会自动识别`./`下面的`Dockerfile`

`:v1.0.20151209`: tag标签

使用以下命令,就可以看到镜像:

		sudo docker images
		
**设置需要挂载的卷**

首先我们挂载一个外置的目录存放mysql数据,保证容器和数据分离:

		/data/mysql
		这个是我放置mysql数据的文件夹
		
创建一个mysql配置文件的文件夹,我们需要开启二进制日志

		/myconfig/mysql/config-file.cnf

我们挂载该目录后,mysql会自动加载其中的 `*.cnf`配置文件,并且和其内部的`my.cnf`合并,该文件内容如下:

		[mysqld]
		log-bin = /var/lib/mysql/binlog
		
**创建容器**

		sudo docker run --name daemon_mysql_server -d -p 3306:3306 -v /data/mysql:/var/lib/mysql -v /myconfig/mysql:/etc/mysql/conf.d -e MYSQL_ROOT_PASSWORD=123456 -it chloroplast/mysql:v1.0.20151209
		
* `--name`: 命名容器
* `-d`: 守护进程模式打开
* `-p`: 3306:3306 把服务器的3306端口和容器的3306端口对接在一起.其实也可以不用暴漏该端口,如果对外发布,则外部也可以访问数据库
* `-v /data/mysql:/var/lib/mysql`: 挂载数据卷
* `-v /myconfig/mysql:/etc/mysql/conf.d`: 配置文件
* `-e`: 初始化mysql密码
* `-it`:
	* `-i`:
	* `-t`:
* `chloroplast/mysql:v1.0.20151209`: 使用的镜像和版本 	
**外部访问该数据库**

		mysql -h120.25.161.1 -uroot -p123456
		
		Warning: Using a password on the command line interface can be insecure.
		Welcome to the MySQL monitor.  Commands end with ; or \g.
		Your MySQL connection id is 97
		Server version: 5.6.28-log MySQL Community Server (GPL)
		
		Copyright (c) 2000, 2015, Oracle and/or its affiliates. All rights reserved.
		
		Oracle is a registered trademark of Oracle Corporation and/or its
		affiliates. Other names may be trademarks of their respective
		owners.
		
		Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
		
		mysql>	
		
上面默认绑定了主机的3306端口号

**使用容器访问数据库**

		sudo docker run -it --link container_name:container_name_alias --rm mysql:5.6 sh -c 'exec mysql -hcontainer_name_alias  -uroot -p123456'
		
* `--link`: 连接我们的`daemon_mysql_server`容器
* `--rm`: 使用后自动删除容器

我们可以连接到我们的数据库容器,而不用对外发布端口.更加安全.

**fig**

继续补充