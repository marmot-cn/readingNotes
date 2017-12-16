# etcd 容器部署

---

因为业务需要etcd容器化部署.

如下是官网示例

## 单节点

		docker run -d -v /usr/share/ca-certificates/:/etc/ssl/certs -p 4001:4001 -p 2380:2380 -p 2379:2379 \
		 --name etcd quay.io/coreos/etcd \
		 -name etcd0 \
		 -advertise-client-urls http://${HostIP}:2379,http://${HostIP}:4001 \
		 -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
		 -initial-advertise-peer-urls http://${HostIP}:2380 \
		 -listen-peer-urls http://0.0.0.0:2380 \
		 -initial-cluster-token etcd-cluster-1 \
		 -initial-cluster etcd0=http://${HostIP}:2380 \
		 -initial-cluster-state new

**指定集群中的同伴信息**

		etcdctl -C http://${HostIP}:2379 member list
		etcdctl -C http://${HostIP}:4001 member list

## 多节点

`-initial-cluster`,必须包含每个etcd节点的`peer urls`.

**etcd0**

		docker run -d -v /usr/share/ca-certificates/:/etc/ssl/certs -p 4001:4001 -p 2380:2380 -p 2379:2379 \
		--name etcd quay.io/coreos/etcd \
		-name etcd0 \
		-advertise-client-urls http://192.168.12.50:2379,http://192.168.12.50:4001 \
		-listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
		-initial-advertise-peer-urls http://192.168.12.50:2380 \
		-listen-peer-urls http://0.0.0.0:2380 \
		-initial-cluster-token etcd-cluster-1 \
		-initial-cluster etcd0=http://192.168.12.50:2380,etcd1=http://192.168.12.51:2380,etcd2=http://192.168.12.52:2380 \
		-initial-cluster-state new
		
**etcd1**

		docker run -d -v /usr/share/ca-certificates/:/etc/ssl/certs -p 4001:4001 -p 2380:2380 -p 2379:2379 \
		 --name etcd quay.io/coreos/etcd \
		 -name etcd1 \
		 -advertise-client-urls http://192.168.12.51:2379,http://192.168.12.51:4001 \
		 -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
		 -initial-advertise-peer-urls http://192.168.12.51:2380 \
		 -listen-peer-urls http://0.0.0.0:2380 \
		 -initial-cluster-token etcd-cluster-1 \
		 -initial-cluster etcd0=http://192.168.12.50:2380,etcd1=http://192.168.12.51:2380,etcd2=http://192.168.12.52:2380 \
		 -initial-cluster-state new
		 
**etcd2**

		docker run -d -v /usr/share/ca-certificates/:/etc/ssl/certs -p 4001:4001 -p 2380:2380 -p 2379:2379 \
		 --name etcd quay.io/coreos/etcd \
		 -name etcd2 \
		 -advertise-client-urls http://192.168.12.52:2379,http://192.168.12.52:4001 \
		 -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
		 -initial-advertise-peer-urls http://192.168.12.52:2380 \
		 -listen-peer-urls http://0.0.0.0:2380 \
		 -initial-cluster-token etcd-cluster-1 \
		 -initial-cluster etcd0=http://192.168.12.50:2380,etcd1=http://192.168.12.51:2380,etcd2=http://192.168.12.52:2380 \
		 -initial-cluster-state new
		 
**指定集群中的同伴信息**

		etcdctl -C http://192.168.12.50:2379,http://192.168.12.51:2379,http://192.168.12.52:2379 member list
		
## 问题

### rancher 部署的问题

在`rancher`部署一开始我映射端口到私有地址, 则会提示失败. 需要一开始映射端口不指定地址. 等部署完成以后在加上服务器的私有地址即可.

### 挂载目录

在配置中发现问题如下,

配置中并没有指名具体的命令`/usr/local/bin/etcd`

需要挂载数据目录,否则重启会丢失:

* 指定数据目录: `-data-dir /data/etcd`
* 挂载目录: `-v /data/etcd:/data/etcd`

### 阿里云镜像

可以使用我在阿里云部署的镜像`registry.cn-hangzhou.aliyuncs.com/marmot/etcd`.