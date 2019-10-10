# filebeat 和 ES 对接 docker 日志

---

## 在本机测试

启动`nginx`镜像.

```shell
[ansible@localhost ~]$ systemctl status docker
● docker.service - Docker Application Container Engine
   Loaded: loaded (/usr/lib/systemd/system/docker.service; enabled; vendor preset: disabled)
   Active: active (running) since Thu 2017-10-05 00:02:15 CST; 8h ago
     Docs: https://docs.docker.com
 Main PID: 1009 (dockerd)
   Memory: 91.8M
   CGroup: /system.slice/docker.service
           ├─1009 /usr/bin/dockerd --graph=/data/docker
...
```

获取`filebeat`

```shell
docker pull docker.elastic.co/beats/filebeat:5.6.
```

启动`es`

```shell
docker run -d -p 9200:9200 --name elasticsearch elasticsearch
```

启动`kibana`

```
docker run --name kibana --link elasticsearch:elasticsearch  -e ELASTICSEARCH_URL="http://elasticsearch:9200" -p 5601:5601 -d kibana
```

启动`filebeat`

```
docker run -d --name filebeat -v /root/filebeat/filebeat.yml:/usr/share/filebeat/filebeat.yml -v /data/docker/containers:/data/docker/containers --link elasticsearch:elasticsearch docker.elastic.co/beats/filebeat:5.6.3
```

curl -XDELETE 'http://127.0.0.1:9200/filebeat-*'

curl http://127.0.0.1:9200/_cat/indices?v

curl -XGET http://127.0.0.1:9200/filebeat-2017.10.05/_search?pretty