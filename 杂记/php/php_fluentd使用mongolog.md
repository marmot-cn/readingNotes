# php fluentd 使用 monolog

---

因为暂时我的日志收集使用的是`fluentd`, 所以配套的应用层日志也想直接使用`fluentd`, 为了方便开发人员使用, 我准备单独写一个`monolog`的`handler`来在`monlog`的基础上使用日志.

## 安装 fluentd 扩展

```
{
    "require": {
        "fluent/logger": "v1.0.0"
    }
}
```

## 基于monolog的handler写一个fluentd的扩展

```php
$log = new Logger('channeltest');
$log->pushHandler(new FluentdHandler(Core::$container->get('fluentd.address'), Core::$container->get('fluentd.port'), Logger::WARNING));
$log->pushHandler(new FluentdHandler(Core::$container->get('fluentd.address'), Core::$container->get('fluentd.port'), Logger::INFO));

$log->warning('Foo',array('xxx'=>'warning'));
$log->info('Bar',array('xxx'=>'info'));
```

这样就会自动录入到`fluentd`.

## 测试关闭fluentd检查日志链接是否报错

关闭`fluentd`容器关闭后, 页面可以正常输出不报错.

## 测试效率

### 插入100000条数句比较速度

#### `MongoDBHandler`

直接插入`mongo`, 使用`mongohandler`,测试`9999条数据`.

```php
root@d934696db77a:/var/www/html# cat test.php
<?php
require './Core.php';

use Monolog\Logger;
use Monolog\Handler\MongoDBHandler;

$core = \Marmot\Core::getInstance();
$core->initCli();

$log = new Logger('name');

$log->pushHandler(new MongoDBHandler(new \MongoDB\Client('mongodb://mongo:27017'),'log','testDebug', Logger::DEBUG));

for($i=0; $i<100000; $i++) {
    $log->debug('Debug',array('xxx'=>'xxx1'));
}

root@61342609884c:/var/www/html# time php test.php

real	0m32.378s
user	0m6.580s
sys		0m8.590s
```

总共**32.378**秒.

#### `FluentdHandler`

```
root@d934696db77a:/var/www/html# cat test2.php
<?php
require './Core.php';

use Monolog\Logger;
use Monolog\Handler\MongoDBHandler;
use System\Extension\Monolog\FluentdHandler;

use Marmot\Core;

$core = \Marmot\Core::getInstance();
$core->initCli();

$log = new Logger('channeltest');
$log->pushHandler(new FluentdHandler(Core::$container->get('fluentd.address'), Core::$container->get('fluentd.port'), Logger::DEBUG));

for($i=0; $i<100000; $i++) {
    $log->debug('Debug',array('xxx'=>'xxx1'));
}

root@61342609884c:/var/www/html# time php test2.php

real	0m13.685s
user	0m1.880s
sys		0m1.200s
```

总共**13.685**秒.

#### `FluentdHandler`挂载`sock`

修改`fluentd.conf`

```
<source>
  @type unix
  path /fluentd/etc/td-agent/td-agent.sock
</source>
```

然后通过`docker`挂载目录方式挂载.

```
root@61342609884c:/var/www/html# cat test3.php
<?php
require './Core.php';

use Monolog\Logger;
use Monolog\Handler\MongoDBHandler;
use System\Extension\Monolog\FluentdHandler;

use Marmot\Core;

$core = \Marmot\Core::getInstance();
$core->initCli();

$log = new Logger('channeltest');
//port 没用, 这里只是测试, 没有修改函数
$log->pushHandler(new FluentdHandler('unix://'.S_ROOT.'sock/fluentd/td-agent.sock', Core::$container->get('fluentd.port'), Logger::DEBUG));

for($i=0; $i<100000; $i++) {
    $log->debug('Debug',array('xxx'=>'xxx1'));
}


```

? 我测试的时间和`tcp`差不多, 不过个人感觉应该比`tcp`快才对.

`https://github.com/fluent/fluentd/issues/1509`查了这个`issue`也发现类似问题. 比较诡异.

#### 总结 

我的测试还是在本机的情况下, 如果`mongo`为其他服务器, 考虑网络开销就非常大. 因为`php`访问`mongo`为短连接. 每次录入数据都要建立链接.

`fluentd`是在本地服务器, 由它负责统一录入数据到`mongo`, 经过配置文件`flush_interval 10s`, 每`10s`刷新一次到`mongo`中.
