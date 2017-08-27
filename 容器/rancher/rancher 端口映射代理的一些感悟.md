# rancher 端口映射代理的一些感悟

---

今天查看服务器,发现使用`netstat -tnlp`没有找见有监听端口. 因为原来版本使用`docker`需要外键访问如果绑定端口都会映射到宿主机上. 但是在服务器上我映射很多端口作为给开发人员测试使用, 可是使用`netstat`后确只发现有我的`ssh`端口.

后来查看`iptables`倒出规则表, 在`nat`下:

```shel
...
-A CATTLE_PREROUTING ! -i docker0 -p tcp -m tcp --dport 8003 -j MARK --set-xmark 0x1068/0xffffffff
-A CATTLE_PREROUTING ! -i docker0 -p tcp -m tcp --dport 8003 -j DNAT --to-destination 10.42.162.18:80
-A CATTLE_PREROUTING -p tcp -m tcp --dport 8003 -m addrtype --dst-type LOCAL -j DNAT --to-destination 10.42.162.18:80
```
发现这3条规则.

```shell
[ansible@sandbox-service-2 ~]$ sudo iptables -t nat -L -n -v
...
   10   640 MARK       tcp  --  !docker0 *       0.0.0.0/0            0.0.0.0/0            tcp dpt:8003 MARK set 0x1068
   10   640 DNAT       tcp  --  !docker0 *       0.0.0.0/0            0.0.0.0/0            tcp dpt:8003 to:10.42.155.146:80
    1    60 DNAT       tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            tcp dpt:8003 ADDRTYPE match dst-type LOCAL to:10.42.155.146:80
...
```

其中`-A CATTLE_PREROUTING ! -i docker0 -p tcp -m tcp --dport 8003 -j MARK --set-xmark 0x1068/0xffffffff`和`-A CATTLE_PREROUTING ! -i docker0 -p tcp -m tcp --dport 8003 -j DNAT --to-destination 10.42.162.18:80`表示**不是**从`docker0`流入的流量且目标端口是`8003`需要做目标地址转换`DNAT`.所以测试从我本机访问`8003`端口流量会被这两条规则匹配上.

`-A CATTLE_PREROUTING -p tcp -m tcp --dport 8003 -m addrtype --dst-type LOCAL -j DNAT --to-destination 10.42.162.18:80`

这条规则使用了`addrtype`模块,且需要匹配目标地址是本地`--dst-type LOCAL`. 经过测试在**同一个网段**的容器间访问会映射到该规则上(自身容器内访问应该会回环到本地,所以不匹配该规则).