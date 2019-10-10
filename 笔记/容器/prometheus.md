# prometheus

## 需要开发事项

1. 同步时间为中国时间, 否则时间序列监控会出现问题.
2. 监控主机`node_exporter`.
	* 网络流量监控?
	* 服务器监控
	* 运行结果统计分析(prometheus)
3. 监控容器`cadvisor`.
	* 应用服务监控
	* 运行结果统计分析(prometheus)
4. 运维服务审计`auditd`,https://linux.cn/article-4907-1.html
	* 运维服务记录
5. 数据库监控`https://hub.docker.com/r/prom/mysqld-exporter/`,`https://github.com/percona/grafana-dashboards`(在dashboards文件夹下可以找见对应的`json`文件)
6. `memcached`监控`https://hub.docker.com/r/prom/memcached-exporter/`
7. `mongodb-grafana`: https://grafana.com/dashboards/2583

## 还未解决的问题

* `prometheus`的数据存储.

```
-v /data/prometheus/data:/prometheus
```

## 准备流程

### 安装 Prometheus

镜像`prom/prometheus`.

配置文件在容器内的位置`/etc/prometheus/prometheus.yml`.

需要挂载目录`/prometheus`在容器内, 收集采集信息.

### 样例配置文件

```
# my global config
global:
  scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
  - static_configs:
    - targets:
      # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
      - targets: ['localhost:9090']

  # Scrape the Node Exporter every 5 seconds.
  - job_name: 'node'
    static_configs:
      - targets: ['your_server_ip:9100']

  - job_name: 'docker'
  	static_configs:
      - targets: ['xxxx:8080','xxxx:8080']

  - job_name: 'mysql global status'
    static_configs:
      - targets: ['192.168.0.143:9104']
    params:
      collect[]:
        - global_status

  - job_name: 'mysql performance'
    scrape_interval: 1m
    static_configs:
      - targets:
        - '192.168.0.143:9104'
    params:
      collect[]:
        - perf_schema.tableiowaits
        - perf_schema.indexiowaits
        - perf_schema.tablelocks
```

### 启动

```
docker run -p 9090:9090 --name=prometheus -v /etc/localtime:/etc/localtime -v /data/prometheus/config/prometheus.yml:/etc/prometheus/prometheus.yml -d registry.cn-hangzhou.aliyuncs.com/marmot/prometheus --config.file=/etc/prometheus/prometheus.yml
```

* `/etc/localtime:/etc/localtime`: 挂载时区.
* `-storage.local.path`: 配置存储路径.
* `-storage.local.memory-chunks`: 调整`Prometheus`在宿主机的内存使用. 大致可以存储1000个时间序列(每个序列占用10个chunks). **这个参数要根据生产环境调整**

### 测试访问

* `http://xxx:9090/`可以访问到`prometheus`是否正常.

### 设置 Node Exporter

`Node Exporter`是一个用于采集服务器数据的服务. 可以采集的信息包括:

* 服务器文件系统
* 网络设备
* 处理器使用情况
* 内存使用情况

#### 启动 Node Exporter

```
docker run -d -p 9100:9100 -v /etc/localtime:/etc/localtime -v "/proc:/host/proc" -v "/sys:/host/sys" -v "/:/rootfs" --net="host" prom/node-exporter --collector.procfs /host/proc -collector.sysfs /host/proc --collector.filesystem.ignored-mount-points "^/(sys|proc|dev|host|etc)($|/)"
```

### 设置对数据库的监控

#### 设置数据库账户

```
CREATE USER 'exporter'@'%' IDENTIFIED BY '123456';
GRANT PROCESS, REPLICATION CLIENT ON *.* TO 'exporter'@'%';
GRANT SELECT ON performance_schema.* TO 'exporter'@'%';
```

#### 安装Mysql Server Exporter

```
docker pull prom/mysqld-exporter

docker run -d -p 9104:9104 --link=mysql-master:mysql  \
        -e DATA_SOURCE_NAME="exporter:123456@(mysql:3306)/database" prom/mysqld-exporter
```

我第一次收集不到数据是发现设置账户的权限有问题.

### 设置对docker容器的监控

#### 安装`cadvisor`

```
docker run \
  --volume=/:/rootfs:ro \
  --volume=/var/run:/var/run:rw \
  --volume=/sys:/sys:ro \
  --volume=/var/lib/docker/:/var/lib/docker:ro \
  --volume=/dev/disk/:/dev/disk:ro \
  --publish=8080:8080 \
  --detach=true \
  --name=cadvisor \
  registry.cn-hangzhou.aliyuncs.com/marmot/cadvisor:latest
```

## `Grafana`

### 设置`Grafana`

我们使用`Grafana`做统计`Prometheus`的数据仪表盘. 可以使用外接的数据库如`mysql`.默认使用本地的`SQLite3`.

```
docker run -d -p 3000:3000 -e "GF_SECURITY_ADMIN_PASSWORD=admin_password" -v /etc/localtime:/etc/localtime -v /data/grafana_db:/var/lib/grafana registry.cn-hangzhou.aliyuncs.com/marmot/grafana
```

默认`Grafana`会自动创建并初始化`SQLite3`数据库在`/var/lib/grafana/grafana.db`.

使用环境变量`GF_SECURITY_ADMIN_PASSWORD`设置密码, 覆盖默认的密码`admin`.

访问`http://xxxx:3000`可以访问到`Grafana`.

添加新的`Data sources`.

* `Type`: `Prometheus`.
* `Http`: xxx:9090

选在`+`, `import`新的模板, 这里

* `node`测试我选用的是`1860(https://grafana.com/dashboards/1860)`.
* `docker`测试我选用的是`395`
* `mysql`我选用的是(https://github.com/percona/grafana-dashboards/tree/master/dashboards)
	* `MySQL_InnoDB_Metrics.json`
	* `MySQL_Overview.json`

## 测试安装

准备环境:

* server1: 安装容器, nginx,php,...
* server2: 安装容器, nginx,php,...
* server3: 安装`Prometheus`.