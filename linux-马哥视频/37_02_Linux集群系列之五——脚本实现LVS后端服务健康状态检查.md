# 37_02_Linux集群系列之五——脚本实现LVS后端服务健康状态检查

---

## 笔记

---

### Director是否可以作为一个realserver

实验场景

* `Director`
	* `enp0s3`, DIP: 192.168.0.178
	* `enp0s3:0`, VIP: 192.168.0.101
* `RS1`:	
	* `enp0s3 `: rip1: 192.168.0.179
	* `lo:0`: vip: 192.168.0.101
* `RS2`:
	* `eth0`: rip1: 192.168.0.180
	* `lo:0`: vip: 192.168.0.101

#### 更换arp缓存

因为我昨晚是在家测试, 今天到公司发现`192.168.0.101`(vip)已经被占用了. 只能手动清理手动绑定`mac`地址了.

```
arp -a
? (192.168.0.1) at f4:83:cd:b8:19:ba on en0 ifscope [ethernet]
? (192.168.0.101) at 74:51:ba:35:10:73 on en0 ifscope [ethernet]
...

而我们 director enp0s3:0 网卡的 mac 地址是 08:00:27:4f:4f:da
[root@localhost ~]# ifconfig
enp0s3: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 192.168.0.178  netmask 255.255.255.0  broadcast 192.168.0.255
        inet6 fe80::a00:27ff:fe4f:4fda  prefixlen 64  scopeid 0x20<link>
        ether 08:00:27:4f:4f:da  txqueuelen 1000  (Ethernet)
        RX packets 10371  bytes 1033180 (1008.9 KiB)
        RX errors 0  dropped 40  overruns 0  frame 0
        TX packets 1435  bytes 184157 (179.8 KiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

enp0s3:0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 192.168.0.101  netmask 255.255.255.255  broadcast 192.168.0.101
        ether 08:00:27:4f:4f:da  txqueuelen 1000  (Ethernet)

lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
        inet 127.0.0.1  netmask 255.0.0.0
        inet6 ::1  prefixlen 128  scopeid 0x10<host>
        loop  txqueuelen 1  (Local Loopback)
        RX packets 6  bytes 504 (504.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 6  bytes 504 (504.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

我们更改 arp 广播的 mac 地址
sudo arp -d 192.168.0.101
192.168.0.101 (192.168.0.101) deleted

sudo arp -s 192.168.0.101 08:00:27:4f:4f:da
Password:

查看后发现替换成功
rp -a
? (192.168.0.1) at f4:83:cd:b8:19:ba on en0 ifscope [ethernet]
? (192.168.0.101) at 8:0:27:4f:4f:da on en0 permanent [ethernet]

curl 192.168.0.101
server-2 is 80
server-1
server-1
```

#### director 添加 webserver

```
在director服务器上
[root@localhost ~]# service nginx start
Redirecting to /bin/systemctl start  nginx.service
[root@localhost ~]# curl 127.0.0.1
director
[root@localhost ~]# ipvsadm -a -t 192.168.0.101:80 -r 127.0.0.1 -g -w 5
[root@localhost ~]# ipvsadm -L -n
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  192.168.0.101:80 wlc
  -> 127.0.0.1:80                 Route   5      0          0
  -> 192.168.0.179:80             Route   2      0          0
  -> 192.168.0.180:80             Route   1      0          0
  
外网访问
curl 192.168.0.101
server-2 is 80
curl 192.168.0.101
server-1
curl 192.168.0.101
director
```

### 自动化服务脚本

使用`shell`脚本把`director`和`realserver`自动化处理, 因为`ipvsadm-save`这个命令只能保存`ipvsadm`的数据, 而不能保存路由信息等数据. 而且`ipvsadm`类似`iptables`服务, 没有实际用处, 只能用来添加命令用.

`director`中为了防止重复启动脚本, 可以写入一个锁文件用来判断.

`realserver`中为了防止重复启动脚本, 判断

* `ifconfig lo:0 | grep $VIP`
* `netstat -rn | grep "lo:0" | grep $VIP`

### 监控`realserver`状况

如果`realserver`出现状况移除, 如果健康在加入.

#### 检查`realserver`的健康状态

##### 使用`elinks`

我们可以根据返回状态来判断.

```
[root@localhost ~]# elinks -dump http:/192.168.0.179
   server-1
[root@localhost ~]# elinks -dump http:/192.168.0.180
   server-2 is 80
[root@localhost ~]# echo $?
0

访问不存在的地址
[root@localhost ~]# elinks -dump http:/192.168.0.181
ELinks: No route to host
[root@localhost ~]# echo $?
1
```

##### 使用隐藏页面

创建隐藏页面, 请求该隐藏页面即可.

```
elinks -dump xxxx/.xxx.html | grep "??" &> /dev/null
```

获取到特殊信息即可.

##### 使用`curl`

获取页面信息. 但可以设置超时时间.

`curl --conect-timeout 1 http://xxx`

curl:

* `-I`: 使用`head`方法发起请求, 只获取响应首部.
* `-s`: 静音模式

#### 脚本

* `FALL_BACK`: 备用服务器
* `CPORT`: 集群端口
* `RPORT`: 真实主机端口

```shell
#!/bin/bash
#

VIP=192.168.0.101

CPORT=80
FALL_BACK=127.0.0.1

RS=("192.168.0.179" "192.168.0.180")
RSTATUS=("1" "1")
RPORT=80
RW=("2" "1")

TYPE=g

add()
{
  ipvsadm -a -t $VIP:$CPORT -r $1:$RPORT -$TYPE -w $2
  [ $? -eq 0 ] && return 0 || return 1
}

del()
{
  ipvsadm -d -t $VIP:$CPORT -r $1:$RPORT
  [ $? -eq 0 ] && return 0 || return 1
}

while :; do
let COUNT=0
for I in ${RS[*]}; do
  if curl --connect-timeout 1 http://$I &> /dev/null; then
    if [ ${RSTATUS[$COUNT]} -eq 0 ]; then
      add $I ${RW[$COUNT]}
      [ $? -eq 0 ] && RSTATUS[$COUNT]=1
     fi
  else
    if [ ${RSTATUS[$COUNT]} -eq 1 ]; then
      del $I
      [ $? -eq 0 ] && RSTATUS[$COUNT]=0
     fi
  fi
  let COUNT++
done
sleep 5
done
```

## 整理知识点

---