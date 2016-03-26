#docker-compose的中间容器

---

####重建容器

重建容器操作事实上就是删除并重新创建一个同名容器,删除并新建的同名容器不但不具有旧容器的volume内容,还会间接导致引用它的其他容器volume内容丢失`,因为docker-compose会按照拓扑关系重新引用新容器的volume`.

####中间容器

dokcer-compose中使用了中间容器(intermediate container)来暂时"记下"旧容器的volume.

**步骤**

1. 中间容器先引用酒容器volume
2. 删除旧容器(只是删除容器,旧容器volume里的内容是保留的: `-v`和`--rm`会删除非从宿主机挂载的卷)
3. 重新创建容器
4. 新容器可以通过引用中间容器获得旧容器volume内容

####示例

**下载debian镜像用于测试**

		[chloroplast@dev-server ~]$ docker pull debian:7
		7: Pulling from debian
		2c788329cf71: Already exists
		c1661b87f436: Already exists
		debian:7: The image you are pulling has been verified. Important: image verification is a tech preview feature and should not be relied on to provide security.
		Digest: sha256:4f5d72ea21bdbf700f7d6a2f252b1af7b9d044cd3c66c687cc3bf9bdb1b5327e
		Status: Downloaded newer image for debian:7
		
		[chloroplast@dev-server ~]$ docker images
		...
		debian                                7                   c1661b87f436        3 weeks ago         84.89 MB
		...
		
**创建用于被挂载的容器`container-test`**

创建一个名字为`container-test`的容器:

		[chloroplast@dev-server ~]$ docker run -it -d --name container-test -h CONTAINER -v /data debian:7

在`/data`(我们挂载的)创建文件`test-file`,并写入内容1:
		
		//进入容器
		[chloroplast@dev-server ~]$ docker exec -it container-test /bin/bash
		root@CONTAINER:/# ls
		bin  boot  data  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  selinux  srv  sys	tmp  usr  var
		root@CONTAINER:/# cd /data/
		root@CONTAINER:/data# touch test-file
		root@CONTAINER:/data# echo 1 > test-file
		root@CONTAINER:/data# cat test-file
		1
		
		//查看宿主机的目录
		[chloroplast@dev-server ~]$ docker inspect -f {{.Volumes}} container-test
		map[/data:/docker/volumes/3af5434182aa49063a72cefc18875f85e2a6de6baed792bfdc6848875c53c415/_data]
		[chloroplast@dev-server ~]$ su root
		[root@dev-server chloroplast]# cat /docker/volumes/3af5434182aa49063a72cefc18875f85e2a6de6baed792bfdc6848875c53c415/_data/test-file
		1 //检查到文件和容器内文件内容一致
			
**停止`container-test`容器**

		[root@dev-server chloroplast]# docker stop container-test
		container-test		
		
我们验证了无论容器是否启动都可以使用`--volumes-from`			
			
**创建容器`container`挂载`container-test`**

		[root@dev-server chloroplast]# docker run -it -d -h NEWCONTAINER --name container --volumes-from container-test debian:7
		f8cc2619f0ae22ec87cf00c04507e59a3703f63c6e417884f0ffdb1a38a0507c

检查从`container-test`挂载的卷:

		[root@dev-server chloroplast]# docker exec -it container /bin/bash
		root@NEWCONTAINER:/# cd /data/
		root@NEWCONTAINER:/data# cat test-file
		1
		
**删除`container-test`容器,并重建**

		[root@dev-server chloroplast]# docker rm container-test
		container-test
		[root@dev-server chloroplast]# docker run -it -d --name container-test -h CONTAINER -v /data debian:7
		33e1131ce70dc5cdd24af534755e2e228a9bcc8c72094643a1c459c00f93c464
		
其实我们此刻进入`container-test`也会发现`/data`为空.		
		
**重建`container`容器,并重建**

		[root@dev-server chloroplast]# docker rm container
container
		[root@dev-server chloroplast]# docker run -it -d -h NEWCONTAINER --name container --volumes-from container-test debian:7
f30e19bd552c01f41da04b4e07070c665903d6a06c7da608f21131d328dbe53f

**两个容器的`/data`都为空**

`container`挂载`container-test`的容器. `container-test`的`/data`目录为空,则肯定为空

		[root@dev-server chloroplast]# docker exec -it container /bin/bash
		root@NEWCONTAINER:/#
		root@NEWCONTAINER:/# cd /data/
		root@NEWCONTAINER:/data# ls
		
####使用docker-compose示例

**编写docker-compose.yml**

		[root@dev-server test-intermediate-container]# cat docker-compose.yml
		container:
		 image: "debian:7"
		 volumes_from:
		  - "container-test"
		 stdin_open: true
		
		container-test:
		  image: "debian:7"	
		  stdin_open: true	

添加`stdin_open: true`,保证容器在运行否则容器在主进程结束会退出 :`A docker container exits when its main process finishes`

`stdin_open: true	` = `dokcer run -i`	 Keep STDIN open even if not attached

**启动容器**

		[root@dev-server test-intermediate-container]# docker-compose up -d
		
		[root@dev-server test-intermediate-container]# docker ps
		CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
		c6375e4360eb        debian:7            "/bin/bash"         3 minutes ago       Up 3 minutes                            testintermediatecontainer_container_1
		7b086321d9d8        debian:7            "/bin/bash"         3 minutes ago       Up 3 minutes                            testintermediatecontainer_container-test_1


我们检查了容器启动时间.

`--no-recreate`(If containers already exist, don't recreate them.)参数不会重建容器,就算配置文件有变化.

		[root@dev-server test-intermediate-container]# docker-compose stop
		Stopping testintermediatecontainer_container_1 ... done
		Stopping testintermediatecontainer_container-test_1 ... done
		[root@dev-server test-intermediate-container]# docker-compose up -d --no-recreate
		Starting testintermediatecontainer_container-test_1
		Starting testintermediatecontainer_container_1
		[root@dev-server test-intermediate-container]# docker ps
		CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
		c6375e4360eb        debian:7            "/bin/bash"         4 minutes ago       Up 1 seconds                            testintermediatecontainer_container_1
		7b086321d9d8        debian:7            "/bin/bash"         4 minutes ago       Up 1 seconds                            testintermediatecontainer_container-test_1		


检查CREATED时间还是为旧的时间

`--force-recreate`会重建容器,就算配置文件无变化.

**测试写入文件**

进入`container-test`容器:

		[root@dev-server test-intermediate-container]# docker exec -it testintermediatecontainer_container-test_1 /bin/bash
		root@7117431f8a3f:/# echo 1 > /data/touch-file
		root@7117431f8a3f:/# cat /data/touch-file
		1

进入`container`容器:

		[root@dev-server test-intermediate-container]# docker exec -it testintermediatecontainer_container_1 /bin/bash
		root@0491e55774cf:/# cat /data/touch-file
		1

**重建容器**

在没有使用docker-compose时,我们重建后内容会丢失.

		[root@dev-server test-intermediate-container]# docker-compose stop
		Stopping testintermediatecontainer_container_1 ...
		Stopping testintermediatecontainer_container-test_1 ...
		
		[root@dev-server test-intermediate-container]# docker-compose up -d --force-recreate
		Recreating testintermediatecontainer_container-test_1
		WARNING: Service "container-test" is using volume "/data" from the previous container. Host mapping "None" has no effect. Remove the existing containers (with `docker-compose rm container-test`) to use the host volume mapping.
		Recreating testintermediatecontainer_container_1
		[root@dev-server test-intermediate-container]# docker ps
		CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
		d2d6538f6586        debian:7            "/bin/bash"         3 seconds ago       Up 2 seconds                            testintermediatecontainer_container_1
		3b534882e627        debian:7            "/bin/bash"         3 seconds ago       Up 2 seconds                            testintermediatecontainer_container-test_1

看`ps`结果中的`CREATED`时间,证明我们是新建的容器.
		
**检查文件**

		[root@dev-server test-intermediate-container]# docker exec -it testintermediatecontainer_container-test_1 /bin/bash
		root@3b534882e627:/# cat /data/touch-file
		1
		
		[root@dev-server test-intermediate-container]# docker exec -it testintermediatecontainer_container_1 /bin/bash
		root@d2d6538f6586:/# cat /data/touch-file
		1
		
docker-compose 使用中间容器(可以理解为`--rm`的容器)来保存了数据.让新建的容器引用中间容器获取旧容器的内容.

但是会产生`warning`?,这点查看compose源码这里是会发出warning.但是官方说1.52已经解决了,这里查了还是没有解决.

这种挂载方式我们一般也不会常用,我们一般会挂载到宿主机上.这里只是测试docker-compose的中间容器.