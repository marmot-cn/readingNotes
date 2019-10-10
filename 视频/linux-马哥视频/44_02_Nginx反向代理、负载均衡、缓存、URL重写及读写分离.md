# 44_02_Nginx反向代理、负载均衡、缓存、URL重写及读写分离

---

## 笔记

### 反向代理

`proxy_pass`

### Nginx

* main
	* worker_process
	* error_log
	* user
	* group
* events: 事件驱动相关
* http: 关于http相关的配置
* server: 虚拟主机
* location: UIR的属性

```
location [op] UIR {

}

op: ~, ~*, ^~, =


location @name
```

### 反向代理路由匹配

#### 情况1

```
location /uri {
  proxy_pass http://xxx.com:8080/newuri;
}
```

`location`的`/uri`将被替换为上游服务器上的`/newuri`.

#### 情况2(使用了模式匹配)

如果`location`的`URI`是通过模式匹配定义的, 其`URI `将直接被传递至上游服务器, 而不能为其指定转换的另一个`URI`.

```
location ~ ^/forum {
  proxy_pass http://xxx.com:8080;
}
```

`/form`会被转换为`http://xxx.com:8080/forum`.

#### 情况3

如果使用`URL`重定向, 那么`nginx`将使用重定向后的`URI`处理请求.

```
location / {
  rewrite /(.*)$ /index.php?page=$1 break;
  proxy_pass http://xxx.com:8080/index;
}
```

传给上游服务器的是`index.php?page=<match>`, 而不是`/index`.

#### 传递真实IP

```
proxy_set_header X-Real-IP $remote_addr
```

#### upstream

需要在`sever`之外定义, 每个`upstream`需要有一个独立的名称.

* `server`: 定义每一个后端服务器. 后面只能跟**地址**, 不能有协议.
* `weight`: 权重, 没有权重就是默认为`0`.
* `backup`: 备用服务器, 所有`server`都宕掉后访问该服务.
* 健康状况检查
	* `max_fails`
	* `fail_timeout`

### nginx 缓存

* cache: 
	* 共享内存: 存储键和缓存对象元数据
	* 磁盘空间: 存储数据
* proxy_cache_path(定义缓存): 不能定义在`server{}`上下文中.
	* levels: 缓存几级目录(levels=1:2, 代表有2级目录,一级目录名字符1个, 二级目录名字符2个)
	* keys_zone = name(键名字):size(多大) (keys_zone=first:20m)
	* max_size: 最大空间
* cache_manager: LRU, 单独进程

#### 使用缓存

```
proxy_cache_path /nginx/cache/first levels=1:2 keys_zone first:20m max_size=1g;

location / {
	proxy_pass http://xxx/;
	proxy_cache first;
}
```

#### 常用缓存

* open_log_cache: 日志缓存, 日志先保存到内存中在同步到文件中
* open_file_cache: 打开文件缓存, 
* fastcgi_cache: 缓存后端处理结果.

### nginx limit

`nginx`的`limit`限制也给予共享内存实现.

### nginx rewrite URL重写

```
if (condition) {

}
```

`if`只能用在`server`和`location`中.

测试:

* 双目测试
* 单目测试

**支持正则表达式**

* `last`: 本次重写完成之后, 重启下一轮检查.
* `break`: 本次重写完成之后, 直接执行后续操作.
* `redirect`: 302
* `permanent`: 301

### WebDAV

Web-based Distributed Authoring and Versioning.

基于`HTTP 1.1`协议的通信协议, 它扩展了`HTTP 1.1`, 在`GET, POST, HEAD`等几个`HTTP`标准方法以外添加了一些新的方法, 使应用程序可以直接对`Web Server`直接读写, 并支持文件锁定(`Locking`)及解锁(`Unlock`), 还可以支持文件的版本控制.

引入`PUT`和`DELETE`.



## 整理知识点