#Centos7 安装docker

####确定内核版本

		[root@iZ944l0t308Z docker]# uname -r
		3.10.0-327.4.4.el7.x86_64
		
至少3.10+

####从yum安装

**更新yum安装包**

		yum -y update
		
**安装docker**

		yum -y install docker
		
**重启时自动启动docker**

		systemctl enable docker.service
		
		或
		
		chkconfig docker on
		
**启动docker**

		systemctl start docker.service
		
		或
		
		service docker start
		
需要`重启服务器`,否则不能启动docker(或者改了配置以后不能启动docker服务,但是重启后解决问题).这个问题不知道是什么情况?`Centos 6.5`没有此类问题.
		
**修改镜像存储路径**

		vi /etc/sysconfig/docker
		
		修改
		
		OPTIONS=--graph=/data/docker --selinux-enabled
		
**重启docker**

		service docker restart

**non-root用户运行docker**

如果还没有`docker group`就添加一个:

		sudo groupadd docker

用户加到docker组

		sudo usermod -aG docker your-user
		
		把用户添加到 "docker" 组
		
**验证安装**

		[root@iZ944l0t308Z ~]# docker run hello-world
		Unable to find image 'hello-world:latest' locally
		Trying to pull repository docker.io/library/hello-world ... latest: Pulling from library/hello-world
		3f12c794407e: Pull complete
		975b84d108f1: Pull complete
		Digest: sha256:8be990ef2aeb16dbcb9271ddfe2610fa6658d13f6dfb8bc72074cc1ca36966a7
		Status: Downloaded newer image for docker.io/hello-world:latest
		
		
		Hello from Docker.
		This message shows that your installation appears to be working correctly.
		
		To generate this message, Docker took the following steps:
		 1. The Docker client contacted the Docker daemon.
		 2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
		 3. The Docker daemon created a new container from that image which runs the
		    executable that produces the output you are currently reading.
		 4. The Docker daemon streamed that output to the Docker client, which sent it
		    to your terminal.
		
		To try something more ambitious, you can run an Ubuntu container with:
		 $ docker run -it ubuntu bash
		
		Share images, automate workflows, and more with a free Docker Hub account:
		 https://hub.docker.com
		
		For more examples and ideas, visit:
		 https://docs.docker.com/userguide/
		 
####从script安装

**要有root权限**

**更新yum安装包**

		yum update -y
		
**安装docker**

		curl -sSL https://get.docker.com/ | sh
		
This script adds the `docker.repo` repository and installs Docker.

这个安装暂时还没找见修改配置文件的地方



