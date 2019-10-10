# 38_04_Linux集群系列之十——高可用集群之heartbeat安装配置

---

## 笔记

---

### heartbeat v2

* node1, node2
	* 需要靠节点名称来识别对方, 每个节点都需要识别其他节点, 对于名称的解析不能依赖于IP, 应该使用本地`hosts`文件, 所有节点的解析文件应该一样.
	* 节点名称必须和`uname -n`命令结果一致.
	* 各个节点需要`ssh`互信通信.
	* 各节点时间需要同步.

### 示例 

#### 配置两台主机ip

#### 配置两台主机hostname

* 临时生效
	* `hostname`
* 永久生效
	* `/etc/sysconfig/network`

#### ssh 互信通信

* ssh-keygen
* ssh-copy-id

#### 时间同步

`ntp`

#### 同步`hosts`文件

#### heartbeat 配置文件

1. 密钥文件, 600 权限, authkeys
2. heartbeat服务的配置文件`ha.cf`

## 整理知识点

---