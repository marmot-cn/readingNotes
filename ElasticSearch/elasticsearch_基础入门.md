### ElasticSearch

---

我们使用容器安装ES,首先使用最简单的方式进行安装测试:
		
下载并启动镜像:
		
		docker run -d -p 9200:9200 --name elasticsearch elasticsearch
				
测试:

		curl http://120.25.161.1:9200/
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

		[ansible@rancher-agent-1 ~]$ curl http://120.25.161.1:9200/_cat/health
		1490857379 07:02:59 elasticsearch yellow 1 1 1 1 0 0 1 0 - 50.0%
		
因为副本中只有一个节点,所以显示`yellow`
		
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
* HOST Elasticsearch集群中的任何一个节点的主机名,如果是在本地的节点,那么就叫localhost.
* PORT Elasticsearch HTTP服务所在的端口,默认为9200
* PATH API路径(例如_count将返回集群中文档的数量)，PATH可以包含多个组件，例如_cluster/stats或者_nodes/stats/jvm
* QUERY_STRING 一些可选的查询请求参数,例如?pretty参数将使请求返回更加美观易读的JSON数据
* BODY 一个JSON格式的请求主体

**计算集群中的文档数量**

		curl -XGET 'http://120.25.161.1:9200/_count?pretty' -d '
		{
		    "query": {
		        "match_all": {}
		    }
		}
		'
返回:

		{
		  "count" : 1,
		  "_shards" : {
		    "total" : 1,
		    "successful" : 1,
		    "failed" : 0
		  }
		}

**显示http头**

		curl -i -XGET http://120.25.161.1:9200
		HTTP/1.1 200 OK
		content-type: application/json; charset=UTF-8
		content-length: 327
		
		{
		  "name" : "Ex-i2hy",
		  "cluster_name" : "elasticsearch",
		  "cluster_uuid" : "YsrfRW-OS8eUsi96cBBwug",
		  "version" : {
		    "number" : "5.2.2",
		    "build_hash" : "f9d9b74",
		    "build_date" : "2017-02-24T17:26:45.835Z",
		    "build_snapshot" : false,
		    "lucene_version" : "6.4.1"
		  },
		  "tagline" : "You Know, for Search"
		}

**在sense中**

![sense](./img/01.png "sense")

#### 文档

Elasticsearch是面向文档(`document oriented`)的,这意味着它可以存储整个对象或文档(document).然而它不仅仅是存储,还会索引(index)每个文档的内容使之可以被搜索.在Elasticsearch中,你可以对文档(而非成行成列的数据)进行索引、搜索、排序、过滤.

##### json

Lasticsearch使用Javascript对象符号(JavaScript Object Notation),也就是JSON,作为文档序列化格式.

#### 索引

索引(indexing)、搜索(search)以及聚合(aggregations).

在Elasticsearch中,文档归属于一种类型(type),而这些类型存在于索引(index)中.

		Relational DB -> Databases -> Tables -> Rows -> Columns
		Elasticsearch -> Indices   -> Types  -> Documents -> Fields
		
Elasticsearch集群可以包含多个索引(indices)(数据库),每一个索引可以包含多个类型(types)(表),每一个类型包含多个文档(documents)(行),然后每个文档包含多个字段(Fields)(列).

**索引的含义**

* 一个索引(index)就像是传统关系数据库中的数据库,它是相关文档存储的地方,index的复数是indices 或indexes.
* 索引(动词)「索引一个文档」表示把一个文档存储到索引(名词)里,以便它可以被检索或者查询.这很像SQL中的INSERT关键字,差别是,如果文档已经存在,新的文档将覆盖旧的文档.
* 倒排索引 传统数据库为特定列增加一个索引,例如B-Tree索引来加速检索.Elasticsearch和Lucene使用一种叫做倒排索引(inverted index)的数据结构来达到相同目的.

默认情况下,文档中的所有字段都会被索引(拥有一个倒排索引),只有这样他们才是可被搜索的.

![createDocument](./img/02.png "createDocument")

* megacorp 索引名
* employee 类型名
* 1 这个员工

#### 搜索

只要执行HTTP GET请求并指出文档的"地址"——索引,类型和ID既可.根据这三部分信息,我们就可以返回原始JSON文档:

		GET /megacorp/employee/1
		
响应内容:

John Smith的原始JSON文档包含在`_source`字段中.

		{
		  "_index": "megacorp",
		  "_type": "employee",
		  "_id": "1",
		  "_version": 1,
		  "found": true,
		  "_source": {
		    "first_name": "John",
		    "last_name": "Smith",
		    "age": 25,
		    "about": "I love to go rock climbing",
		    "interests": [
		      "sports",
		      "music"
		    ]
		  }
		}

* 通过HTTP方法GET来检索文档
* 使用DELETE方法删除文档
* 使用HEAD方法检查某文档是否存在
* 如果想更新已存在的文档,我们只需再PUT一次

##### 简单搜索

搜索全部员工的请求:

		GET /megacorp/employee/_search
		
你可以看到我们依然使用megacorp索引和employee类型,但是我们在结尾使用关键字`_search`来取代原来的文档ID.


		{
		  "took": 59,
		  "timed_out": false,
		  "_shards": {
		    "total": 5,
		    "successful": 5,
		    "failed": 0
		  },
		  "hits": {
		    "total": 3,
		    "max_score": 1,
		    "hits": [
		      {
		        "_index": "megacorp",
		        "_type": "employee",
		        "_id": "2",
		        "_score": 1,
		        "_source": {
		          "first_name": "Jane",
		          "last_name": "Smith",
		          "age": 32,
		          "about": "I like to collect rock albums",
		          "interests": [
		            "music"
		          ]
		        }
		      },
		      {
		        "_index": "megacorp",
		        "_type": "employee",
		        "_id": "1",
		        "_score": 1,
		        "_source": {
		          "first_name": "John",
		          "last_name": "Smith",
		          "age": 25,
		          "about": "I love to go rock climbing",
		          "interests": [
		            "sports",
		            "music"
		          ]
		        }
		      },
		      {
		        "_index": "megacorp",
		        "_type": "employee",
		        "_id": "3",
		        "_score": 1,
		        "_source": {
		          "first_name": "Douglas",
		          "last_name": "Fir",
		          "age": 35,
		          "about": "I like to build cabinets",
		          "interests": [
		            "forestry"
		          ]
		        }
		      }
		    ]
		  }
		}


搜索姓氏中包含"Smith"的员工.要做到这一点,我们将在命令行中使用轻量级的搜索方法.这种方法常被称作查询字符串(query string)搜索,因为我们像传递URL参数一样去传递查询语句:


请求中依旧使用`_search`关键字,然后将查询语句传递给参数`q=`.这样就可以得到所有姓氏为Smith的结果.


		GET /megacorp/employee/_search?q=last_name:Smith

		{
		  "took": 128,
		  "timed_out": false,
		  "_shards": {
		    "total": 5,
		    "successful": 5,
		    "failed": 0
		  },
		  "hits": {
		    "total": 2,
		    "max_score": 0.2876821,
		    "hits": [
		      {
		        "_index": "megacorp",
		        "_type": "employee",
		        "_id": "2",
		        "_score": 0.2876821,
		        "_source": {
		          "first_name": "Jane",
		          "last_name": "Smith",
		          "age": 32,
		          "about": "I like to collect rock albums",
		          "interests": [
		            "music"
		          ]
		        }
		      },
		      {
		        "_index": "megacorp",
		        "_type": "employee",
		        "_id": "1",
		        "_score": 0.2876821,
		        "_source": {
		          "first_name": "John",
		          "last_name": "Smith",
		          "age": 25,
		          "about": "I love to go rock climbing",
		          "interests": [
		            "sports",
		            "music"
		          ]
		        }
		      }
		    ]
		  }
		}
		
		
##### 使用DSL语句查询

DSL(Domain Specific Language特定领域语言)以JSON请求体的形式出现.我们可以这样表示之前关于"Smith"的查询:

		GET /megacorp/employee/_search
		{
		    "query" : {
		        "match" : {
		            "last_name" : "Smith"
		        }
		    }
		}
		
不再使用查询字符串(`query string`)做为参数.而是使用请求体代替.其中使用了`match`语句.

##### 更复杂的搜索

我们依旧想要找到姓氏为"Smith"的员工,但是我们只想得到年龄大于30岁的员工.我们的语句将添加过滤器(filter),它使得我们高效率的执行一个结构化搜索:


		GET /megacorp/employee/_search
		{
		    "query" : {
		        "filtered" : {
		            "filter" : {
		                "range" : {
		                    "age" : { "gt" : 30 } 
		                }
		            },
		            "query" : {
		                "match" : {
		                    "last_name" : "smith" 
		                }
		            }
		        }
		    }
		}
		
这里是旧版本,使用`filter`,新的使用`BoolQuery`.

[BoolQuery](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-bool-query.html "bool")


		POST /megacorp/employee/_search
		{
		    "query" : {
		        "bool" : {
		            "filter" : {
		                "range" : {
		                    "age" : { "gt" : 30 } 
		                }
		            },
		            "must" : {
		                "term" : {
		                    "last_name" : "smith" 
		                }
		            }
		        }
		    }
		}

#### 全文搜索

搜索所有喜欢“rock climbing”的员工:

		GET /megacorp/employee/_search
		{
		    "query" : {
		        "match" : {
		            "about" : "rock climbing"
		        }
		    }
		}
		
返回结果:		
		
		{
		  "took": 24,
		  "timed_out": false,
		  "_shards": {
		    "total": 5,
		    "successful": 5,
		    "failed": 0
		  },
		  "hits": {
		    "total": 2,
		    "max_score": 0.53484553,
		    "hits": [
		      {
		        "_index": "megacorp",
		        "_type": "employee",
		        "_id": "1",
		        "_score": 0.53484553, ### 这里是相关性评分
		        "_source": {
		          "first_name": "John",
		          "last_name": "Smith",
		          "age": 25,
		          "about": "I love to go rock climbing",
		          "interests": [
		            "sports",
		            "music"
		          ]
		        }
		      },
		      {
		        "_index": "megacorp",
		        "_type": "employee",
		        "_id": "2",
		        "_score": 0.26742277, ### 这里是相关性评分
		        "_source": {
		          "first_name": "Jane",
		          "last_name": "Smith",
		          "age": 32,
		          "about": "I like to collect rock albums",
		          "interests": [
		            "music"
		          ]
		        }
		      }
		    ]
		  }
		}
		
默认情况下,Elasticsearch根据结果相关性评分来对结果集进行排序,所谓的「结果相关性评分」就是文档与查询条件的匹配程度.很显然,排名第一的John Smith的about字段明确的写到"rock climbing".

但是为什么Jane Smith也会出现在结果里呢?原因是“rock”在她的abuot字段中被提及了.因为只有"rock"被提及而"climbing"没有,所以她的_score要低于John.

Elasticsearch如何在各种文本字段中进行全文搜索,并且返回相关性最大的结果集.相关性(relevance)的概念在Elasticsearch中非常重要,而这个概念在传统关系型数据库中是不可想象的,因为传统数据库对记录的查询只有匹配或者不匹配.

##### 短语搜索

想要确切的匹配若干个单词或者短语(phrases).例如我们想要查询同时包含"rock"和"climbing"(并且是相邻的)的员工记录.

只要将`match`查询变更为`match_phrase`查询即可.

		GET /megacorp/employee/_search
		{
		    "query" : {
		        "match_phrase" : {
		            "about" : "rock climbing"
		        }
		    }
		}
		
该查询返回John Smith的文档:

		{
		  "took": 39,
		  "timed_out": false,
		  "_shards": {
		    "total": 5,
		    "successful": 5,
		    "failed": 0
		  },
		  "hits": {
		    "total": 1,
		    "max_score": 0.53484553,
		    "hits": [
		      {
		        "_index": "megacorp",
		        "_type": "employee",
		        "_id": "1",
		        "_score": 0.53484553,
		        "_source": {
		          "first_name": "John",
		          "last_name": "Smith",
		          "age": 25,
		          "about": "I love to go rock climbing",
		          "interests": [
		            "sports",
		            "music"
		          ]
		        }
		      }
		    ]
		  }
		}
		
##### 高亮我们的搜索

从每个搜索结果中高亮(highlight)匹配到的关键字,这样用户可以知道为什么这些文档和查询相匹配

在之前的语句上增加`highlight`参数:

		GET /megacorp/employee/_search
		{
		    "query" : {
		        "match_phrase" : {
		            "about" : "rock climbing"
		        }
		    },
		    "highlight": {
		        "fields" : {
		            "about" : {}
		        }
		    }
		}
		
当我们运行这个语句时,会命中与之前相同的结果,但是在返回结果中会有一个新的部分叫做highlight,这里包含了来自about字段中的文本,并且用<em></em>来标识匹配到的单词.

		{
		  "took": 137,
		  "timed_out": false,
		  "_shards": {
		    "total": 5,
		    "successful": 5,
		    "failed": 0
		  },
		  "hits": {
		    "total": 1,
		    "max_score": 0.53484553,
		    "hits": [
		      {
		        "_index": "megacorp",
		        "_type": "employee",
		        "_id": "1",
		        "_score": 0.53484553,
		        "_source": {
		          "first_name": "John",
		          "last_name": "Smith",
		          "age": 25,
		          "about": "I love to go rock climbing",
		          "interests": [
		            "sports",
		            "music"
		          ]
		        },
		        "highlight": {
		          "about": [
		            "I love to go <em>rock</em> <em>climbing</em>" #在这里高亮
		          ]
		        }
		      }
		    ]
		  }
		}
		
#### 聚合

##### 分析

Elasticsearch有一个功能叫做聚合(aggregations),它允许你在数据上生成复杂的分析统计.它很像SQL中的GROUP BY但是功能更强大.
 
让我们找到所有职员中最大的共同点(兴趣爱好)是什么:

