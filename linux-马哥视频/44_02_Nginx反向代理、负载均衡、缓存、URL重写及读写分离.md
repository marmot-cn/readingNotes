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
proxy_set_header X-Real-IP
```

## 整理知识点