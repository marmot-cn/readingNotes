#Rancher

----

###概述

Rancher 只需要每个host的`CPU`,`memory`,`local disk storage`(硬盘)和`network connectivity`(网络连通性).

`Cattle`是`Rancher`自己独有的容器编排工具.

###Quick Start Guide

机器配制需求:

* kernel of 3.10+
* 1GB 内存

		[chloroplast@iZ94jqmwawyZ ~]$ uname -r
		3.10.0-123.9.3.el7.x86_64

启动`rancher server`:

		 sudo docker run -d --restart=always -p 8080:8080 rancher/server
		 sudo docker logs -f containerid
		 
`UI`暴漏在`8080`端口.`http://server_ip:8080`

添加`host`,host的ip需要通过添加`CATTLE_AGENT_IP`来添加.

**Network Agent容器**

`Network Agent`容器是是`Rancher`创建的系统容器,提供:

* 跨主机访问
* 健康检查
* ...

**overlay网络**

无论主机ip是什么,容器的ip和network agent的ip都在`10.42.*.*`范围内.这个是由`Rancher`创建的`overlay`网络,这样容器可以跨主机访问容器.

**自己创建容器并且在rancher的`overlay`网络内**

		 docker run -d -it --label io.rancher.container.network=true ubuntu:14.04.2
		 
`Rancher`将会启动这个容器去方位`overlay`网络.

###安装rancher server

####单节点

配制外置数据库:

		> CREATE USER 'cattle'@'%' IDENTIFIED BY 'xxxx'; 
		> CREATE DATABASE IF NOT EXISTS cattle COLLATE = 'utf8_general_ci' CHARACTER SET = 'utf8';
		> GRANT ALL ON cattle.* TO 'cattle'@'%' IDENTIFIED BY 'cattle';
		> GRANT ALL ON cattle.* TO 'cattle'@'localhost' IDENTIFIED BY 'cattle';
		
启动`rancher server`:

		sudo docker run -d --restart=always -p 8080:8080 \
		    -e CATTLE_DB_CATTLE_MYSQL_HOST=<hostname or IP of MySQL instance> \
		    -e CATTLE_DB_CATTLE_MYSQL_PORT=<port> \
		    -e CATTLE_DB_CATTLE_MYSQL_NAME=<Name of Database> \
		    -e CATTLE_DB_CATTLE_USERNAME=<Username> \
		    -e CATTLE_DB_CATTLE_PASSWORD=<Password> \
		    rancher/server

####多节点

**端口需求**

每个`Node`的端口需求:

* `Global Access`: Tcp Ports `22`,`80`,`443`,`18080`
* `Node之间`: 
	* UDP Ports: `500`, `4500`
	* TCP Ports: `2181`, `2376`, `2888`, `3888`, `6379`
	
**大规模部署需求**	
	
大部署需求:

* 每个rancher server节点必须至少有 `4GB` 或者 `8GB` 的 `heap size`,需要至少 `8GB` 或者 `16GB`的RAM
* Mysql 备份

**高可用**

`HA`支持3种节点数量配置:

* `1 Node`: 不是高可用
* `3 Nodes`: 可以允许任何`1`个节点宕机
* `5 Nodes`: 可以允许任何`2`个节点宕机

###更新

####RANCHER SERVER

* `rancher/server:latest`,`latest`标记表示最新版.
* `rancher/server:stable`,`stable`稳定版.

#####更新: 有external数据库,建议使用这种模式部署server端

**停止容器**

**更新镜像**

**使用新版本的镜像启动容器(external数据库)**

**移除旧版本的容器**

因为旧容器标记着`--restart=always`,所以机器在重启的时候,容器也会被重启.

#####更新: 有外挂data数据目录

**停止容器**

**更新镜像**

**使用新版本的镜像启动容器(挂在相同目录)**

**移除旧版本的容器**

#####更新: 没有外挂data数据目录

**停止容器**

		docker stop <container_name_of_original_server>
		
**创建rancher-data容器**

		 docker create --volumes-from <container_name_of_original_server> --name rancher-data rancher/server:<tag_of_previous_rancher_server>
		 
**更新镜像**

		docker pull rancher/server:latest
		
**启动容器**

不要在启动中中断,可能会持续一段时间.


		docker run -d --volumes-from rancher-data --restart=always -p 8080:8080 rancher/server:latest

#####更新: server HA 模式
	
	
		sudo docker rm -f $(sudo docker ps -a | grep rancher | awk {'print $1'}) 
		
		//使用server的最新版本
		sudo sh rancher-ha.sh rancher/server:v1.1.0
		
###RANCHER 基础服务

####网络

Rancher通过实现`覆盖网络(overlay network)`的网络隧道实现跨主机通信.如果一个容器需要加入rancher的网络,需要添加`label`:

		--label io.rancher.container.network=true
		
在Rancher的网络内,容器会被同时赋予一个 Docker bridge IP 段(172.17.0.0/16)和一个 Rancher 管理的IP (10.42.0.0/16)在默认的 docker0 网关.

容器在同一个环境内是可以被路由和访问的通过rancher的网络.

####负载均衡 Load Balancer

Rancher 通过 HAProxy 实现负载均衡.

Our load balancer has HAProxy software installed on the load balancer agent containers. The load balancer uses a round robin algorithm from HAProxy to select the target services.

只能用在rancher的网络内.

端口`42`不能被用在`source port`, 因为`rancher`使用这个端口用于`health checks`.

`示例`:

		Source Port		Target Port
		80				8080
		81				8081
		
		Target Services
		Service1
		Service2
		Service3
		
* 任何通往主机`80`端口的数据,会被`round-robin-ed`到`Service1, Service2, Service3`的端口`8080`和`8081`.
* 任何通往主机`81`端口的数据,会被`round-robin-ed`到`Service1, Service2, Service3`的端口`8080`和`8081`.

**内部负载均衡**	

选择`Intenal`. 所有Internal ports 只能被同一个环境的服务访问到.

###Cattle

####Stack

**docker-compose.yml**

描述一个服务,传统的docker-compose编排文件.

**rancher-compose.yml**

管理rancher启动服务的额外信息.

####Service

**service高可用**

类似`k8s`,我们描述一个服务的最终数量状态,rancher会监控其健康状态自动调节服务为这个数量.如:我们期望服务有3个分布在4台服务器上,如果我们手动关掉一个服务的服务器,rancher也会自动调度另外的一台服务器启动服务,保证有3个服务可用.

