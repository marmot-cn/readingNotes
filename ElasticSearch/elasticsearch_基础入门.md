### ElasticSearch

---

我们使用容器安装ES,首先使用最简单的方式进行安装测试:
		
下载并启动镜像:
		
		docker run -d -p 9200:9200 --name elasticsearch elasticsearch
				
测试:

		curl http://120.25.87.35:9200/
		{
		  "name" : "wq5sX1D",
		  "cluster_name" : "elasticsearch",
		  "cluster_uuid" : "5FiKl_6NTx-vfP3ROY0ufw",
		  "version" : {
		    "number" : "5.2.2",
		    "build_hash" : "f9d9b74",
		    "build_date" : "2017-02-24T17:26:45.835Z",
		    "build_snapshot" : false,
		    "lucene_version" : "6.4.1"
		  },
		  "tagline" : "You Know, for Search"
		}		
				
检查健康: 

		[ansible@rancher-agent-1 ~]$ curl http://127.0.0.1:9200/_cat/health
		1489739804 08:36:44 elasticsearch green 1 1 0 0 0 0 0 0 - 100.0%
		
#### 节点

单个`节点`可以作为一个运行中的 Elasticsearch 的实例.而一个`集群`是一组拥有相同 cluster.name 的节点， 他们能一起工作并共享数据，还提供容错与可伸缩性


#### kibana

**安装**

		docker run --name kibana --link elasticsearch:elasticsearch -p 5601:5601 -d kibana
		
**安装 sense 插件**

教程里面提示需要安装,但是检查后发现已经默认安装在内了(5.2).

### 和ElasticSearch 交互

#### Java 使用方式

Elasticsearch 内置的两个客户端:

* 节点客户端(Node client)
* 传输客户端(Transport client)

**节点客户端**

节点客户端作为一个非数据节点加入到本地集群中.它本身不保存任何数据,但是它知道数据在集群中的哪个节点中,并且可以把请求转发到正确的节点.

**传输客户端**

轻量级的传输客户端可以可以将请求发送到远程集群.它本身不加入集群,但是它可以将请求转发到集群中的一个节点上.


两个 Java 客户端都是通过 `9300` 端口并使用本地 Elasticsearch 传输 协议和集群交互.集群中的节点通过端口 `9300` 彼此通信。如果这个端口没有打开,节点将无法形成一个集群.

#### Restful 方式

curl -X<VERB> `<PROTOCOL>://<HOST>:<PORT>/<PATH>?<QUERY_STRING>` -d `<BODY>`

* VERB : 适当的 HTTP 方法 或 谓词 : GET`、 `POST`、 `PUT`、 `HEAD 或者 `DELETE`.
* PROTOCOL : http 或者 https (如果你在 Elasticsearch 前面有一个 https 代理).
