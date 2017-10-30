# fluentd 收集 docker 日志

---

## 安装fluentd

### 下载镜像

```
[ansible@localhost ~]$ docker pull fluent/fluentd:v0.12-debian
v0.12-debian: Pulling from fluent/fluentd
bc95e04b23c0: Already exists
06d84b97b8fc: Pull complete
73aff3a6c1ef: Pull complete
dbc928c93438: Pull complete
d138a36e5b8d: Pull complete
3d48a36daec4: Pull complete
89f3f689a3f9: Pull complete
Digest: sha256:fd2a0ebe20a9e04e16bff0e03b568dd6dde66fea559f762dda10745c175e8ce4
Status: Downloaded newer image for fluent/fluentd:v0.12-debian
```

### 创建配置文件

```
[ansible@localhost ~]$ cat fluentd/fluentd.conf
<source>
  @type http
  port 9880
  bind 0.0.0.0
</source>
<match **>
  @type stdout
</match>
```

从`http`接收日志, 并且输出到标准输出.

### 运行并测试

```shell
[ansible@localhost ~]$ docker run -d -p 9880:9880 -v /home/ansible/fluentd:/fluentd/etc -e FLUENTD_CONF=fluentd.conf fluent/fluentd:v0.12-debian
b6f24f4bf280ca09bda19287d65b34692729f33c9896d60e71e6db58291d5f54
[ansible@localhost ~]$ docker ps
CONTAINER ID        IMAGE                                    COMMAND                  CREATED             STATUS              PORTS                                         NAMES
b6f24f4bf280        fluent/fluentd:v0.12-debian              "/bin/entrypoint.s..."   12 seconds ago      Up 10 seconds       5140/tcp, 24224/tcp, 0.0.0.0:9880->9880/tcp   upbeat_euclid
```

通过`http`发送一条日志信息.

```shell
[ansible@localhost ~]$ curl -X POST -d 'json={"json":"message"}' http://localhost:9880/sample.test
[ansible@localhost ~]$ docker logs upbeat_euclid
2017-10-05 15:23:43 +0000 [info]: reading config file path="/fluentd/etc/fluentd.conf"
2017-10-05 15:23:43 +0000 [info]: starting fluentd-0.12.40
2017-10-05 15:23:43 +0000 [info]: gem 'fluentd' version '0.12.40'
2017-10-05 15:23:43 +0000 [info]: adding match pattern="**" type="stdout"
2017-10-05 15:23:43 +0000 [info]: adding source type="http"
2017-10-05 15:23:43 +0000 [info]: using configuration file: <ROOT>
  <source>
    @type http
    port 9880
    bind 0.0.0.0
  </source>
  <match **>
    @type stdout
  </match>
</ROOT>
2017-10-05 15:24:32 +0000 sample.test: {"json":"message"}
```

可见最后一条我们已经接收到日志.

## 结合容器测试

### 修改配置文件

```
nginx:
  image: "nginx"
  ports:
    - "80:80"
  volumes:
    - ./:/var/www/html/
  log_driver: "fluentd"
  log_opt:
     fluentd-address: "127.0.0.1:24224"
     tag: "demo.nginx-1"
  external_links:
    - "fluentd:fluentd"
  container_name: nginx-1
```

    
### nginx 修改日志格式为json

```
log_format json '{ "@timestamp": "$time_local", '
    	    '"remote_ip": "$remote_addr", '
	    '"remote_user": "$remote_user", '
	    '"request": "$request", '
	    '"response": "$status", '
	    '"bytes": "$body_bytes_sent", '
	    '"referrer": "$http_referer", '
	    '"agent": "$http_user_agent" }';

access_log  /var/log/nginx/access.log  json;
```

### 安装mongo插件

```shell
apt-get update

# 因为里面的ruby版本是2.3
apt-get install ruby2.3-dev
apt-get install make
apt-get install gcc
```

### 最终编译文件

```dockerfile
FROM fluent/fluentd:v0.12-debian

RUN apt-get update && apt-get install -y gcc make ruby2.3-dev \
&& fluent-gem install fluent-plugin-mongo \
&& apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false gcc make \
&& set -ex \
&& ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone
```

* 安装`gcc`,`make`,`ruby2.3-dev`
* 安装`mongo`插件
* 清理安装包
* 修改时区为国内时区

## 配置文件

实现目标:

1. 绑定主机端口, 使容器可以使用`fluentd`日志存储引擎倒出数据.
2. 日志使用`json`格式.
3. 日志存储到`mongo`, 集合使用`capped`属性

```
<source>
  @type forward
  port 24224
  bind 0.0.0.0
</source>

<filter *.**>
  @type parser
  format json
  key_name log
  reserve_data true
</filter>

<match *.**>
    @type copy
    <store>
    @type mongo
    host dds-uf68cf9433726fa42.mongodb.rds.aliyuncs.com
    port 3717
    database sandboxlog
    collection services

    user sandboxlog
    password hWcExTEbz4cA2yh2

    capped
    capped_size 1024m

    time_key time

    flush_interval 10s
    </store>
    <store>
    @type stdout
    </store>
 </match>
```

### 参数解释

#### source

```
<source>
  @type forward
  port 24224
  bind 0.0.0.0
</source>
```

`in_forward`插件监听在`TCP`端口收取一个事件流. 也同样监听在一个`UDP`端口上接收心跳信息.

* `@type`必须是`forward`.
* `port`监听端口号, 默认值是`24224`.
* `bind`绑定监听地址.

#### filter

```
<filter *.**>
  @type parser
  format json
  key_name log
  reserve_data true
</filter>
```

`filter_parser`过滤插件, 过滤数据并修改事件日志格式.

* `format`: `json`格式数据.
* `key_name`: 需要过滤的字段.
* `reserve_data`: 保留原始数据, 默认为`false`.

#### match

```
<match *.**>
    @type copy
    <store>
    @type mongo
    host dds-uf68cf9433726fa42.mongodb.rds.aliyuncs.com
    port 3717
    database sandboxlog
    collection services

    user sandboxlog
    password hWcExTEbz4cA2yh2

    capped
    capped_size 1024m

    time_key time

    flush_interval 10s
    </store>
    <store>
    @type stdout
    </store>
 </match>
```

```
<store>
@type stdout
</store>
```

* `sotre`至少一个该标签. 描述存储目标.
* `stdout`把日志输出到标准输出.

```
<store>
@type mongo
host dds-uf68cf9433726fa42.mongodb.rds.aliyuncs.com
port 3717
database sandboxlog
collection services

user sandboxlog
password hWcExTEbz4cA2yh2

capped
capped_size 1024m

time_key time

flush_interval 10s
</store>
```

* `@type mongo`, 使用`mongo`存储.
* `host`, `mongo`地址.
* `port`, `mongo`端口.
* `database`, 数据库.
* `collection`集合名称.
* `user`, `mongo`用户名.
* `password`, `mongo`密码.
* `capped`, 使用`mongo`中`capped`属性, 但是我设置了并没与自动设定集合这个属性.
* `capped_size`, 集合大小, 因为我上个参数设置了没用, 这个参数属性也没生效.
* `time_key` 时间戳的名字
* `flush_interval` 刷新时间间隔, 默认是`60s`.