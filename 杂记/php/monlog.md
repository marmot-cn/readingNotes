# monlog

## 安装 monlog 扩展

```
composer require monolog/monolog
```

## 安装 mongodb

创建一个新目录, 配置`docker-compose.yml`文件, 下载镜像并启动.

```
mongo:
  image: "registry.cn-hangzhou.aliyuncs.com/marmot/mongo-3.2:1.0"
  volumes:
   - ./:/data/db
  container_name: mongo
```

## 开发环境链接 mongodb

外链`mongo`.

```
  external_links:
   - "mongo:mongo"
```

## 测试代码

### 日志级别

* DEBUG：详细的debug信息
* INFO：感兴趣的事件。像用户登录，SQL日志
* NOTICE：正常但有重大意义的事件。
* WARNING：发生异常，使用了已经过时的API。
* ERROR：运行时发生了错误，错误需要记录下来并监视，但错误不需要立即处理。
* CRITICAL：关键错误，像应用中的组件不可用。
* ALETR：需要立即采取措施的错误，像整个网站挂掉了，数据库不可用。这个时候触发器会通过SMS通知你，

### 样例代码

#### 存储日志

```php
use Monolog\Logger;
use Monolog\Handler\MongoDBHandler;

$log = new Logger('name');
$log->pushHandler(new MongoDBHandler(new \MongoDB\Client('mongodb://mongo:27017'),'log','testWarning', Logger::WARNING));
$log->pushHandler(new MongoDBHandler(new \MongoDB\Client('mongodb://mongo:27017'),'log','testError', Logger::ERROR));
$log->pushHandler(new MongoDBHandler(new \MongoDB\Client('mongodb://mongo:27017'),'log','testInfo', Logger::INFO));
$log->pushHandler(new MongoDBHandler(new \MongoDB\Client('mongodb://mongo:27017'),'log','testDebug', Logger::DEBUG));

$log->warning('Foo',array('xxx'=>'xxx1'));
$log->info('Bar',array('xxx'=>'xxx1'));
$log->debug('Debug',array('xxx'=>'xxx1'));
    
`warning`的日志会存储在自己这个级别以及更低级别的日志内.

testWarning 集合存储1条 Foo
testError 不存储因为级别高
testInfo 集合存储2条 Bar, Foo
testDebu 集合存储3条 Bar, Foo, Debug
```

`MongoDBHandler`:

```php
public function __construct($mongo, $database, $collection, $level = Logger::DEBUG, $bubble = true)

$bubble 属性，这个属性定义了handler是否拦截记录不让它流到下一个handler, 可以理解为冒泡, 默认可以流入到多个.
```

#### 检测从mongo查询日志

进入`mongo`库查询, 可以找见数据.

```shell
docker exec -it mongo /bin/bash
root@afe77aae2050:/# mongo
MongoDB shell version: 3.2.10
connecting to: test
Welcome to the MongoDB shell.
For interactive help, type "help".
For more comprehensive documentation, see
	http://docs.mongodb.org/
Questions? Try the support group
	http://groups.google.com/group/mongodb-user
>
> show dbs;
distributionProduct  0.000GB
local                0.000GB
log                  0.000GB
product              0.000GB
userGroup            0.000GB
workflow             0.000GB
> use log;
switched to db log
> show collections;
testDebug
testInfo
testWarning
>> db.testDebug.find().limit(10);
{ "_id" : ObjectId("59f6f897d64d750008191c92"), "message" : "Foo", "context" : { "xxx" : "xxx1" }, "level" : 300, "level_name" : "WARNING", "channel" : "name", "datetime" : "2017-10-30 18:01:59", "extra" : [ ] }
{ "_id" : ObjectId("59f6f897d64d750008191c95"), "message" : "Bar", "context" : { "xxx" : "xxx1" }, "level" : 200, "level_name" : "INFO", "channel" : "name", "datetime" : "2017-10-30 18:01:59", "extra" : [ ] }
{ "_id" : ObjectId("59f6f897d64d750008191c97"), "message" : "Debug", "context" : { "xxx" : "xxx1" }, "level" : 100, "level_name" : "DEBUG", "channel" : "name", "datetime" : "2017-10-30 18:01:59", "extra" : [ ] }
```

#### 从php获取日志

```php

$mongo = new \MongoDB\Client('mongodb://mongo:27017');

//查询某个库
$db = $mongo->log;

//查询某个集合
$collection = $db->testDebug;

//获取10条, 按照时间倒叙查询
$cursor = $collection->find([], ['limit'=>10, 'sort'=>['datetime'=>-1]]);

foreach ($cursor as $document) {
	echo '<pre>';
	print_r($document);
}
```