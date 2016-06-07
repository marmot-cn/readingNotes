#Centos7安装kubernetes

---

####预处理

**禁用防火墙**

		systemctl stop firewalld #停止firewall
		systemctl disable firewalld #禁止firewall开机启动
		
使用`iptables`,不过我`iptables`只是了解一个大概,需要单独编写iptables的使用情况.
		
**禁用selinux**

		vim /etc/selinux/config
		#SELINUX=enforcing
		SELINUX=disabled
		
####一些声明

* `masterIP`: master IP
* `minionIP`: minion IP
		
####master

**安装`etcd`与`kubernetes-master`**

		yum -y install etcd kubernetes-master
		
**修改`etcd`配置文件**

`/etc/etcd/etcd.conf`:

		ETCD_NAME=default
		ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
		ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:2379"
		ETCD_ADVERTISE_CLIENT_URLS="http://masterIP:2379"
		
这里如果etcd为单独的服务器配置,则使用自己的ip.		
		
**修改`/etc/kubernetes/apiserver`配置文件**

		KUBE_API_ADDRESS="--address=0.0.0.0"
		KUBE_ETCD_SERVERS="--etcd_servers=http://masterIP:2379"
		KUBE_SERVICE_ADDRESSES="--service-cluster-ip-range=10.254.0.0/16"
		KUBE_ADMISSION_CONTROL="--admission_control=NamespaceLifecycle,NamespaceExists,LimitRanger,SecurityContextDeny,ResourceQuota"
		KUBE_API_ARGS="--secure-port=0"		
		
`--secure-port=0`,是我在查看`/var/log/message`发现

		Unable to listen for secure (open /var/run/kubernetes/apiserver.crt: no such file or directory); will try again
		
		原因:
		By default, the kube-apiserver process tries to open a secure (https) server port on port 6443 using credentials from the directory /var/run/kubernetes. If you want to disable the secure port, you can pass --secure-port=0 which should make your error go away. Alternatively, you can manually create certificates for your cluster so that the process is able to successfully open the secure port. 
	
**修改`/etc/kubernetes/controller-manager`配置文件**		

		KUBE_CONTROLLER_MANAGER_ARGS="--node-monitor-grace-period=10s --pod-eviction-timeout=10s"
		
**修改`/etc/kubernetes/config`配置文件**

		KUBE_LOGTOSTDERR="--logtostderr=true"
		KUBE_LOG_LEVEL="--v=0"
		KUBE_ALLOW_PRIV="--allow_privileged=false"
		KUBE_MASTER="--master=http://masterIP:8080"

**启动服务**

		systemctl enable etcd kube-apiserver kube-scheduler kube-controller-manager
		systemctl start etcd kube-apiserver kube-scheduler kube-controller-manager
		
**定义`flannel`网络配置**

定义`flannel`网络配置到`etcd`,这个配置会推送到各个`minions`的`flannel`服务上:

		etcdctl mk /xxx/network/config '{"Network":"172.17.0.0/16"}'
		

####minon

**安装docker**

详见`centos7 安装docker`文档

**minion结点的安装**

		yum -y install kubernetes-node flannel
		
**修改`kube-node`配置**

`/etc/kubernetes/config`:

		KUBE_LOGTOSTDERR="--logtostderr=true"
		KUBE_LOG_LEVEL="--v=0"
		KUBE_ALLOW_PRIV="--allow_privileged=false"
		KUBE_MASTER="--master=http://masterIP:8080"		
`/etc/kubernetes/kubelet`:

		KUBELET_ADDRESS="--address=0.0.0.0"
		KUBELET_PORT="--port=10250"
		KUBELET_HOSTNAME="--hostname_override=minionIP"
		KUBELET_API_SERVER="--api_servers=http://masterIP:8080"
		KUBELET_ARGS="--pod-infra-container-image=kubernetes/pause"		
**修改`flannel`配置**

`/etc/sysconfig/flanneld`:

		FLANNEL_ETCD="http://masterIP:2379"
		FLANNEL_ETCD_KEY="/xxx/network"
		
这里使用的是`masterIP`,如果etcd有自己的独立服务器则使用自己的ip.  
`FLANNEL_ETCD_KEY`在上文已经定义过.

**启动服务**

		systemctl enable flanneld kubelet kube-proxy
		systemctl restart flanneld docker
		systemctl start kubelet kube-proxy	
		
**`ifconfig`**

可以在`minion`上看见2块网卡:`docker0`和`flannel0`:


		[chloroplast@dev-server-1 ~]$ ifconfig
		docker0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1472
		        inet 172.17.39.1  netmask 255.255.255.0  broadcast 0.0.0.0
		        ether 02:42:e0:85:8d:9d  txqueuelen 0  (Ethernet)
		        RX packets 235  bytes 14633 (14.2 KiB)
		        RX errors 0  dropped 0  overruns 0  frame 0
		        TX packets 213  bytes 15177 (14.8 KiB)
		        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
		 flannel0: flags=4305<UP,POINTOPOINT,RUNNING,NOARP,MULTICAST>  mtu 1472
		        inet 172.17.39.0  netmask 255.255.0.0  destination 172.17.39.0
		        unspec 00-00-00-00-00-00-00-00-00-00-00-00-00-00-00-00  txqueuelen 500  (UNSPEC)
		        RX packets 0  bytes 0 (0.0 B)
		        RX errors 0  dropped 0  overruns 0  frame 0
		        TX packets 0  bytes 0 (0.0 B)
		        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
		        
		...
		
####检查状态

在`master`上运行`kubectl get nodes`,可以发现如下信息:

		[chloroplast@dev-server-2 ~]$ kubectl get nodes
		NAME            LABELS                                 STATUS
		10.116.138.44   kubernetes.io/hostname=10.116.138.44   Ready
		
`10.116.138.44`: 这个IP是我阿里云服务器内网的ip

当我们在添加一个`minion`,服务自发现

		[chloroplast@dev-server-2 ~]$ kubectl get nodes
		NAME            LABELS                                 STATUS
		10.116.138.44   kubernetes.io/hostname=10.116.138.44   Ready
		10.44.88.189    kubernetes.io/hostname=10.44.88.189    Ready
		
