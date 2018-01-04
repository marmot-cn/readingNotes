# 生产环境docker部署mongo副本集

## 版本

部署了`mongo:3.6`, 通过镜像部署集成.

## 生产环境调整参数

参数调整部署在`ansible`剧本内.

## 部署

### 下载ca证书

### 创建文件夹

### mongo配置文件

`mongo`的配置文件.

```
net:
  ssl:
    mode: requireSSL
    PEMKeyFile: /etc/ssl/client.pem
    PEMKeyPassword: '!&,wa~/M8@=\*Mw>'
    CAFile: /etc/ssl/cacert.pem
    clusterFile: /etc/ssl/mongo-1.pem
    clusterPassword: 'N;N.C4qD"jw`3M,g'
    allowInvalidHostnames: true
security:
    clusterAuthMode: x509
replication:
    replSetName: "marmot"
```



简单部署完成后, 查看日志发现:

```
WARNING: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine

WARNING: Access control is not enabled for the database.
Read and write access to data and configuration is unrestricted.

WARNING: /sys/kernel/mm/transparent_hugepage/enabled is 'always'.
We suggest setting it to 'never'

WARNING: /sys/kernel/mm/transparent_hugepage/defrag is 'always'.
We suggest setting it to 'never'
```

查看手册, 还有如下约定:

```
Disable the tuned tool if you are running RHEL 7 / CentOS 7 in a virtual environment.

When RHEL 7 / CentOS 7 run in a virtual environment, the tuned tool automatically invokes a performance profile derived from performance throughput, which automatically sets the readahead settings to 4MB. This can negatively impact performance.
```

docker run --name mongo-1 -p "27017:27017" -v /data/mongo-1:/data/db -v /data/mongo-1/mongod.conf:/etc/mongod.conf -d mongo:3.6 --config /etc/mongod.conf

docker run --name mongo-2 -p "27018:27017" -v /data/mongo-2:/data/db -v /data/mongo-2/mongod.conf:/etc/mongod.conf -d mongo:3.6 --config /etc/mongod.conf

docker run --name mongo-3 -p "27019:27017" -v /data/mongo-3:/data/db -v /data/mongo-3/mongod.conf:/etc/mongod.conf -d mongo:3.6 --config /etc/mongod.conf


rs.initiate( {
   _id : "marmot",
   members: [
      { _id: 0, host: "10.170.148.109:27017" },
      { _id: 1, host: "10.170.148.109:27018" },
      { _id: 2, host: "10.170.148.109:27019" }
   ]
})


mongo --ssl --sslAllowInvalidHostnames --sslCAFile /data/db/cacert.pem --sslPEMKeyFile /data/db/mongo.pem


### 部署副本集带auth

先启动不带`auth`的节点

```
docker run --name mongo-1 -p "27017:27017" -v /data/mongo-1:/data/db -d mongo:3.6 --replSet "marmot" --keyFile /data/db/keyfile
docker run --name mongo-2 -p "27018:27017" -v /data/mongo-2:/data/db -d mongo:3.6 --replSet "marmot" --keyFile /data/db/keyfile
docker run --name mongo-3 -p "27019:27017" -v /data/mongo-3:/data/db -d mongo:3.6 --replSet "marmot" --keyFile /data/db/keyfile
```

在其中一台节点先生成用户

```
use admin
db.createUser(
  {
    user: "myUserAdmin",
    pwd: "abc123",
    roles: [ { role: "userAdminAnyDatabase", db: "admin" } ]
  }
)
db.createUser( {
    user: "siteRootAdmin",
    pwd: "abc123",
    roles: [ { role: "root", db: "admin" } ]
 });

use admin
db.auth("myUserAdmin", "abc123")
db.auth("siteRootAdmin", "abc123")

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

关闭所有节点, 添加`--auth`, 然后在重新启动

```
docker run --name mongo-1 -p "27017:27017" -v /data/mongo-1:/data/db -d mongo:3.6 --auth --replSet "marmot" --keyFile /data/db/keyfile

docker run --name mongo-2 -p "27018:27017" -v /data/mongo-2:/data/db -d mongo:3.6 --auth --replSet "marmot" --keyFile /data/db/keyfile

docker run --name mongo-3 -p "27019:27017" -v /data/mongo-3:/data/db -d mongo:3.6 --auth --replSet "marmot" --keyFile /data/db/keyfile
```

在到主节点测试添加数据, 现在好了

```
docker exec -it mongo-1 mongo admin
db.auth("siteRootAdmin", "abc123")

use demo
db.createUser(
  {
    user: "myTester",
    pwd: "xyz123",
    roles: [ { role: "readWrite", db: "demo" }]
  }
)
```

db.getMongo().setSlaveOk()

先不加

```
security:
    authorization: enabled
```

docker run --name mongo-1 -p "27017:27017" -v /data/mongo-1:/data/db -d mongo:3.6 --config /data/db/mongod.conf

docker run --name mongo-2 -p "27018:27017" -v /data/mongo-2:/data/db -d mongo:3.6 --config /data/db/mongod.conf

docker run --name mongo-3 -p "27019:27017" -v /data/mongo-3:/data/db -d mongo:3.6 --config /data/db/mongod.conf

```
mongo --ssl  --sslAllowInvalidHostnames --sslCAFile /data/db/cacert.pem --sslPEMKeyFile /data/db/mongo.pem
```

集群不加`--auth`