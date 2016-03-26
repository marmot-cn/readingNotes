#Kubernetes 集群配置文件解析

---

该文件用于理解配置Kubernetes相关集群`配置`时和`使用`时的一些配置文件参数的理解.里面每个单独的组件`etcd`,`flannel`在单独罗列.

###master

####`etcd`

**ETCD_NAME**

etcd节点名称,如果etcd集群只有一个node,这一项可以注释不用配置,默认名称为`default`.

**ETCD_DATA_DIR**

etcd存储数据的目录

**ETCD_LISTEN_CLIENT_URLS**

client通信端口

**ETCD_ADVERTISE_CLIENT_URLS**

client广播端口

####`kubernetes/apiserver`

**KUBE_API_ADDRESS**

监听的接口,如果配置为127.0.0.1则只监听localhost,配置为0.0.0.0会监听所有接口.

**KUBE_API_PORT**

`KUBE_API_PORT="--port=8080"`

apiserver的监听端口,默认8080,不用修改.

**KUBELET_PORT**

`KUBELET_PORT="--kubelet_port=10250"`

minion上kubelet监听的端口,默认10250,不用修改.

**KUBE_ETCD_SERVERS**

etcd服务地址,前面已经启动了etcd服务,这里配置上面`etcd`的`ETCD_ADVERTISE_CLIENT_URLS`.

**KUBE_SERVICE_ADDRESSES**

kubernetes可以分配的ip的范围,kubernetes启动的每一个`pod`以及`serveice`都会分配一个`ip`地址,将从这个范围分配.

**KUBE_ADMISSION_CONTROL**

??

**KUBE_API_ARGS**

需要额外添加的配置项,简单地启用一个集群无需配置.

####`kubernetes/controller-manager`

**KUBE_CONTROLLER_MANAGER_ARGS**

需要额外添加的参数

####`kubernetes/config`

**KUBE_LOGTOSTDERR**

`KUBE_LOGTOSTDERR="--logtostderr=true"`

表示错误日志记录到`文件`还是输出到`stderr`.

**KUBE_LOG_LEVEL**

`KUBE_LOG_LEVEL="--v=0"`

日志等级

**KUBE_ALLOW_PRIV**

`KUBE_ALLOW_PRIV="--allow_privileged=false"`

是否允许运行特权容器。

**KUBE_MASTER**

`apiserver` 的地址


###kube-node

####`kubernetes/config`

和上文一致

####`kubernetes/kubelet`

**KUBELET_ADDRESS**

`KUBELET_ADDRESS="--address=0.0.0.0"`

minion监听的地址,每个minion根据实际的ip配置.

?? 0.0.0.0 的作用

**KUBELET_PORT**

`KUBELET_PORT="--port=10250"`

监听端口,不要修改,如果修改,同时需要修改master上配置文件中涉及的配置项.

**KUBELET_HOSTNAME**

`KUBELET_HOSTNAME="--hostname_override=10.116.138.44"`

kubernetes看到的minion的名称,设置和ip地址一样便于识别.

**KUBELET_API_SERVER**

`KUBELET_API_SERVER="--api_servers=http://10.170.148.109:8080"`

`apiserver` 的地址

**KUBELET_ARGS**

`KUBELET_ARGS="--pod-infra-container-image=kubernetes/pause"`

额外增加参数

指定`kubernetes/pause`镜像,如果不指定该镜像默认是放在GCE的镜像仓库里了，`gcr.io/google_containers/puase:0.8.0`,因为google被墙了,访问不了.所以启动pod的时候会因为不能下载该镜像而报错.解决方案有2种:

1. 下载镜像`docker.io/kubernetes/pause`,然后重新`tag`为`gcr.io/google_containers/pause:0.8.0`

		[root@iZ94xwu3is8Z ~]# docker images
		REPOSITORY                       TAG                 IMAGE ID            CREATED             VIRTUAL SIZE
		120.25.87.35:5000/php-test       latest              970c5a38b506        43 hours ago        444.2 MB
		<none>                           <none>              d7f0e75cf11f        8 months ago        360.3 MB
		docker.io/kubernetes/pause       latest              6c4579af347b        18 months ago       239.8 kB
		gcr.io/google_containers/pause   0.8.0               6c4579af347b        18 months ago       239.8 kB


2. 修改启动参数,从docker下载`docker.io/kubernetes/pause`

####`/etc/sysconfig/flanneld`

单独阐述flanneld在kubernetes的应用

**FLANNEL_ETCD**

etcd的地址

**FLANNEL_ETCD_KEY**

etcd的key






