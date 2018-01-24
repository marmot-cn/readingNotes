# 生产环境docker部署mongo副本集

## 版本

部署了`mongo:3.6`, 通过镜像部署集成.

## 生产环境调整参数

参数调整部署在`ansible`剧本内.

## 部署

### 下载ca证书

### 创建文件夹

* `ca`, 
* `config`
* `data`

```
[ansible@demo mongo-1]$ tree
.
|-- ca
|	|-- client.pem 客户端访问证书
| 	|--	cacert.pem CA证书
|   `-- mongo-1.pem 数据同步使用的整数
|-- config
|   `-- mongod.conf 配置文件
`-- data
```

### mongo配置文件

`mongo`的配置文件.

每个节点都需要一个配置文件, `clusterFile`会根据不同的节点使用不同的证书.
其他节点的整数都已经集成到镜像内. 

```
net:
  ssl:
    mode: requireSSL
    PEMKeyFile: /etc/ssl/private/client.pem
    PEMKeyPassword: 'xxxxx'
    CAFile: /etc/ssl/private/cacert.pem
    clusterFile: /etc/ssl/private/mongo-?.pem
    clusterPassword: 'xxxx'
    allowInvalidHostnames: true
security:
    clusterAuthMode: x509
replication:
    replSetName: "marmot"
```

### 启动mongo节点

```
docker run --name mongo -p "ip:端口:27017" -v /data/mongo-?/data:/data/db -v /data/mongo-?/config/mongod.conf:/etc/mongod.conf -v /data/mongo-?/ca:/etc/ssl/private -d registry.cn-hangzhou.aliyuncs.com/marmot/mongo-replica-set --config /etc/mongod.conf
```

### 配置副本集

进入其中一台节点.

```
rs.initiate( {
   _id : "marmot",
   members: [
      { _id: 0, host: "xxx.xxx.xxx.xxx:端口" },
      { _id: 1, host: "xxx.xxx.xxx.xxx:端口" },
      { _id: 2, host: "xxx.xxx.xxx.xxx:端口" }
   ]
})
```

创建用户管理员.

```
use admin
db.createUser(
  {
    user: "admin",
    pwd: "xxxxxx",
    roles: [ { role: "userAdminAnyDatabase", db: "admin" } ]
  }
) 
```

创建root.

```
db.createUser( {
    user: "root",
    pwd: "xxxx",
    roles: [ { role: "root", db: "admin" } ]
 });
```

验证用户名和密码.

```
use admin
db.auth("admin", "xxxx")
db.auth("root", "xxx")
```

创建`test`库用于测试, 开发环境直接给予`root`账户也可以.

```
use test
db.createUser(
  {
    user: "myTester",
    pwd: "xyz123",
    roles: [ { role: "readWrite", db: "test" },
             { role: "readWrite", db: "reporting" } ]
  }
)
db.auth("myTester", "xyz123")
db.demo.insert( { item: "card", qty: 15 } )
db.demo.find()
```

查看集群状态

```
use admin
db.auth("root", "xxxx")
rs.status()
```

### 链接客户端

```
mongo --ssl --sslAllowInvalidHostnames --sslCAFile /etc/ssl/private/cacert.pem --sslPEMKeyFile /etc/ssl/private/client.pem
```

### 访问从节点数据

需要设置 

```
db.getMongo().setSlaveOk()
```