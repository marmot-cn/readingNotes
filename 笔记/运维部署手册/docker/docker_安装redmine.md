#docker 安装redmine

####安装redmine的数据库容器

**设置需要挂载的卷**

这里我们只挂载数据卷

		/data/mysql_redmine/
		这个是我放置redmine的mysql数据的文件夹
		
**创建数据库容器**

		sudo docker run --name mysql-redmine -v /data/mysql-redmine/:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.6
		
我们直接选用本地的mysql:5.6的镜像


####安装redmine容器并连接数据库容器

**把redmine镜像拉到本地**

		sudo docker pull sameersbn/redmine:latest
		
redmine没有官方镜像,选取了一个评分较高的镜像

**创建redmine容器**

		sudo docker run --name=redmine -p 10080:80 -d  -it --link mysql-redmine:mysql -e "DB_USER=root" -e "DB_PASS=123456" -e "DB_NAME=redmine" -v /data/redmine/:/home/redmine/data sameersbn/redmine
		
* `--link`: 连接redmine数据库容器(`mysql-redmine:mysql`),注意这里要命名成成`mysql`,这样redmine会自动和数据库连接
* `-e`: 我们事先连接到数据库(可以用容器连接,也可以实现开放一个端口在外部连接).创建`database`,`用户名`和`密码`
* `sameersbn/redmine`: 使用我们下载到本地的redmine镜像

这样可以在端口`10080`访问到`redmine`.

**docker-compose**

		redmine:
		 image: "sameersbn/redmine"
		 ports:
		  - "10080:80"
		 volumes:
		  - /data/redmine/:/home/redmine/data
		 environment:
		  - DB_USER=root
		  - DB_PASS=123456
		  - DB_NAME=redmine
		 links:
		  - mysql
		mysql:
		 image: "mysql:5.6"
		 volumes:
		  - /data/mysql-redmine/:/var/lib/mysql
		 environment:
		  - MYSQL_ROOT_PASSWORD=123456
	  
**守护进程模式启动**

		[root@iZ94ebqp9jtZ redmine]# docker-compose up -d
		Creating redmine_mysql_1
		Creating redmine_redmine_1		
		
		