# nginx-php-memcached 数据访问

---

## 应用场景

因为最近要准备些文件服务, 排除`Linux`的服务器配置, 需要在应用层考虑怎么通过缓存来提升访问量. 如果不加缓存, 则实现流程是:

```
用户通过id -> nginx -> php -> 缓存 -> 如果缓存不存在, 从mysql获取数据
```

但是如果我的`php`已经缓存了数据, 其实我可以让`nginx`直接到`memcached`获取数据看看.

## 解决方案

### 首先尝试非集群的缓存获取数据

#### 准备测试环境

```
nginx:
 image: "registry.cn-hangzhou.aliyuncs.com/nginx-phpfpm/nginx-front"
 ports:
  - "80:80"
 links:
  - "phpfpm"
  - "memcached-1"
 volumes:
  - ./:/var/www/html/

phpfpm:
  image: "registry.cn-hangzhou.aliyuncs.com/phpfpm/phpfpm-end"
  volumes:
   - ./:/var/www/html/
  links:
   - "memcached-1"

memcached-1:
  image: "registry.aliyuncs.com/marmot/memcached:1.0"
```

我们期望结果:

1. `php`在`memcached`写入`key=test`和`value=test-value`.
2. 通过`nginx`可以直接访问到.

#### 准备php代码

```php
<?php

$memcached = new \Memcached();
$memcached->addServers([['memcached-1','11211']]);

$memcached->set($_GET['key'], 'hello'.$_GET['key']);

var_dump($memcached->get($_GET['key']));

var_dump('hello web');
```

这样我们如果从浏览器访问`php`(`http://127.0.0.1/?key=a`), 则会输出

```
string(6) "helloa" string(9) "hello web" 
```

#### 准备nginx

现在我们期望访问到`http://127.0.0.1/?key=a`输出:

```
string(6) "helloa"
```

即不会输出`php`的输出页面.

但是如果我们访问到`http://127.0.0.1/?key=b`输出:

```
string(6) "hellob" string(9) "hello web" 
```

即还是会跳到`php`.

`nginx`的配置文件:

```
location / {
    set            $memcached_key "$uri?$args";
    memcached_pass memcached-1:11211;
    error_page     404 502 504 = @fallback;
}

location @fallback {
    index  index.html index.php;
    try_files $uri $uri index.php?$args;
}
```

#### 测试结果

第一次获取数据时候回出现`hello web`. 第二次就从`nginx`直接连接到`memcached`了.

#### 测试`memcached`关闭后的超时

我们测试当把`memcached`关闭后会出现什么情况.

```
            Name                          Command               State          Ports
-------------------------------------------------------------------------------------------
nginxmemcached_memcached-1_1   docker-entrypoint.sh memcached   Exit 0
nginxmemcached_nginx_1         nginx -g daemon off;             Up       0.0.0.0:80->80/tcp
nginxmemcached_phpfpm_1        docker-php-entrypoint php-fpm    Up       9000/tcp
```

我们访问页面发现会卡主一段时候再次到`php`, 我们需要减少超时时间.

```
location / {
    set            $memcached_key "$uri?$args";
    memcached_connect_timeout 1s;
    memcached_read_timeout 1s;
    memcached_pass memcached-1:11211;
    error_page     404 502 504 = @fallback;
}

location @fallback {
    index  index.html index.php;
    try_files $uri $uri index.php?$args;
}
```

### 尝试集群的缓存获取数据

集群环境需要使用额外的第三方扩展, 但是考虑到后期准备使用:

`openresty` + `redis`来实现缓存, 且也能使用`lua`语言进行动态的扩展.