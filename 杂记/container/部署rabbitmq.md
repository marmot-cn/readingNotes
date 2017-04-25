# 部署rabbitMq

---

		 docker run -d --hostname my-rabbit --name some-rabbit -v /data/rabbitMq:/var/lib/rabbitmq/mnesia/rabbit -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password rabbitmq:3
		 
**--hostname**

因为我是在rancher中使用,所以不用指定该参数.

**--name**

容器名字

**rabbitmq:3**

区分为 `3`和`3-management` 两个版本的镜像, 后面的版本自带web管理.

**-v**

把队列的数据库`mnesia`持久化挂载到我们的目录.

**-e RABBITMQ_DEFAULT_USER**

默认用户,替代默认的`guest`用户

**-e RABBITMQ_DEFAULT_PASS**

默认用户密码,替代默认的密码`guest`

**其他一些参数**

* /etc/rabbitmq/rabbitmq.config 配置文件

