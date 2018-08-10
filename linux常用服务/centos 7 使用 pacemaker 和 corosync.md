# centos 7 使用 pacemaker 和 corosync

## 简介

**Pacemaker**是一款开源的集群管理软件, 可以实现最高的服务可用性. 它是由ClusterLabs分发的高级可扩展HA集群管理器.

**Corosync Cluster Engine**是一个在新的BSD许可下从OpenAIS项目派生的开源项目. 这是一个群组通信系统, 具有在应用程序中实现高可用性的附加功能.

有一些应用程序的**Pacemaker**接口. Pcsd是用于管理Pacemaker的Pacemaker命令行界面和GUI之一. 我们可以使用pcsd命令来创建, 配置或添加一个新节点到群集. 

参考[https://www.howtoing.com/how-to-set-up-nginx-high-availability-with-pacemaker-corosync-on-centos-7](https://www.howtoing.com/how-to-set-up-nginx-high-availability-with-pacemaker-corosync-on-centos-7)

## 先决条件

* web1: 192.168.0.203
* web2: 192.168.0.205
* vip: 192.168.0.206

## 映射主机文件

**所有服务器**

编辑`/etc/hosts`, 加入

```
192.168.0.203 web1
192.168.0.205 web2
```

测试主机映射配置

```
ping -c 3 web1
ping -c 3 web2
```

## 安装webserver

**所有服务器**

这里我们只是测试了`vip`, 所以简单安装了2台`nginx`.

其中`web1`的`index.html`写入`test`. `web2`的`index.html`使用默认页面.

## 安装和配置Pacemaker, Corosync和Pcsd

**所有服务器**

```
yum -y install corosync pacemaker pcs
systemctl enable pcsd
systemctl enable corosync
systemctl enable pacemaker
```

启动`pcsd`

```
systemctl start pcsd
```

设置用户`hacluster`的密码, 该用户在软件安装过程中自动创建.

```
passwd hacluster
...
```

## 创建和配置群集 

只在**web1**运行

建立集群

```
pcs cluster setup --name test_cluster web1 web2
```

启动所有集群服务并启用

```
pcs cluster start --all
pcs cluster enable --all
```

检查集群状态

```
[root@web1 ansible]# pcs status cluster
[root@web1 ansible]# pcs status cluster
Cluster Status:
 Stack: corosync
 Current DC: web2 (version 1.1.18-11.el7_5.3-2b07d5c5a9) - partition with quorum
 Last updated: Wed Aug  8 13:37:36 2018
 Last change: Wed Aug  8 12:47:41 2018 by root via cibadmin on web1
 2 nodes configured
 1 resource configured

PCSD Status:
  web1: Online
  web2: Online
```

## 禁用STONITH并忽略仲裁策略 

禁用`STONITH`

```
pcs property set stonith-enabled=false
```

忽略仲裁

```
pcs property set no-quorum-policy=ignore
```

检查属性列表并确保stonith和法定人数策略被禁用

```
pcs property list
```

## 添加浮动IP

创建`virtual_ip`

```
pcs resource create virtual_ip ocf:heartbeat:IPaddr2 ip=192.168.0.206 nic='enp0s3:0' cidr_netmask=24 op monitor interval=30s
```

检查可用的资源

```
[root@web1 ansible]# pcs status resources
 virtual_ip	(ocf::heartbeat:IPaddr2):	Started web1
```

## 将约束规则添加到群集 

这里我就没设置了, 我只是测试了`vip`.

## 测试

我们停机一台, 可以使用`vip`正常访问到`nginx`, 已经路由到`web2`了.

```
pcs cluster stop web1
```

我们直接关机

```
在web2上 reboot
```

测试也可以正常访问到对应的`nginx`页面.
