# ssh隧道

---

## 隧道是什么 

隧道是一种**把一种网络协议封装进另外一种网络协议进行传输的技术**.

这里我们研究ssh隧道, 所以所有的网络通讯都是加密. 又被称作端口转发, 因为ssh隧道通常会绑定一个本地端口, 所有发向这个端口端口的数据包, 都会被加密并透明地传输到远端系统.

## SSH隧道的类型

1. 动态端口转发
2. 本地端口转发
3. 远端端口转发

## 动态端口转发

动态端口允许通过配置一个本地端口, 把通过隧道到数据转发到远端的所有地址. 本地的应用程序需要使用Socks协议与本地端口通讯. 此时SSH充当Socks代理服务器的角色.

`ssh -D [bind_address:]port`

* `bind_address`: 指定绑定的IP地址, 默认情况会绑定在本地的回环地址(127.0.0.1), 如果空值或者为*会绑定本地所有的IP地址, 如果希望绑定的端口仅供本机使用, 可以指定为`localhost`.
* `port` 指定本地绑定的端口

### 使用场景

我们现在访问页面, 查看本机`ip`地址.

```
本机IP: 36.46.49.120陕西省西安市 电信
```

我们使用一台服务器做代理.

```
我们在本地: ssh -D 127.0.0.1:8080 ansible@120.25.161.1 -p 17456


我们在浏览器设置代理

地址:  127.0.0.1
端口: 8080

这时候访问ip: 120.25.161.1广东省深圳市 阿里云. 变为我们代理服务器ip地址.

我们访问另外一台内网服务器ip http://10.44.88.189/ 可以访问到页面

```

在`mac`下使用

我是新版的`mac`需要关闭`SIP`, 如果关闭`SIP`

1. 重启 Mac，按住 Command+R 键直到 Apple logo 出现，进入 Recovery Mode
2. 点击 Utilities > Terminal
3. 在 Terminal 中输入 csrutil disable，之后回车
4. 重启 Mac

安装`proxychains`

1. `brew install proxychains-ng`
2. 编辑配置文件`/usr/local/etc/proxychains.conf`在末尾的`[ProxyList]`下加入代理类型`socks5 127.0.0.1 8080`(地址和端口根据自己的设置)
3. 现在只要在需要使用的命令前面加上`proxychains`即可使用`socks5`代理了.

```
开启代理
ssh -D 127.0.0.1:8080 ansible@120.25.161.1 -p 17456
Last login: Sun Aug 27 20:59:37 2017 from 36.46.49.120

Welcome to Alibaba Cloud Elastic Compute Service !

-bash: warning: setlocale: LC_CTYPE: cannot change locale (UTF-8): No such file or directory
[ansible@iZ94ebqp9jtZ ~]$...


在mac下通过proxychains使用socks5代理方位内网ip地址的服务器

proxychains4 ssh root@10.44.88.189 -p 17456
[proxychains] config file found: /usr/local/etc/proxychains.conf
[proxychains] preloading /usr/local/Cellar/proxychains-ng/4.12_1/lib/libproxychains4.dylib
[proxychains] DLL init: proxychains-ng 4.12
[proxychains] Strict chain  ...  127.0.0.1:8080  ...  10.44.88.189:17456  ...  OK
The authenticity of host '[10.44.88.189]:17456 ([10.44.88.189]:17456)' can't be established.
ECDSA key fingerprint is SHA256:0eNede9Jv8kTihqKjAkx8jEG/60t1979ggxRSDfdxtQ.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '[10.44.88.189]:17456' (ECDSA) to the list of known hosts.
root@10.44.88.189's password:
Last login: Sun Aug 27 21:03:33 2017 from 10.116.138.44

Welcome to aliyun Elastic Compute Service!

-bash: warning: setlocale: LC_CTYPE: cannot change locale (UTF-8): No such file or directory
[root@iZ94xwu3is8Z ~]# ls

```


## 其他

我们有两台服务器:

* A: 安装nginx, 80端口提供web服务. 页面全是`chloroplast`
* B: 使用`iptables`做`DNAT`和`SNAT`反向代理`A`内网`80`端口

```shell
登录B服务器, 监听内网端口
[ansible@iZ94ebqp9jtZ ~]$ sudo tcpdump tcp port 80 -i eth0 -X
...
20:17:51.855622 IP 10.44.88.189.http > iZ94ebqp9jtZ.12257: Flags [P.], seq 237:304, ack 695, win 124, options [nop,nop,TS val 437370637 ecr 2825627897], length 67
	0x0000:  4500 0077 af1b 4000 3806 9bdc 0a2c 58bd  E..w..@.8....,X.
	0x0010:  0a74 8a2c 0050 2fe1 db96 1101 e14d 6caa  .t.,.P/......Ml.
	0x0020:  8018 007c d018 0000 0101 080a 1a11 bf0d  ...|............
	0x0030:  a86b a8f9 6368 6c6f 726f 706c 6173 7463  .k..chloroplastc
	0x0040:  686c 6f72 6f70 6c61 7374 6368 6c6f 726f  hloroplastchloro
	0x0050:  706c 6173 7463 686c 6f72 6f70 6c61 7374  plastchloroplast
	0x0060:  6368 6c6f 726f 706c 6173 7463 686c 6f72  chloroplastchlor
	0x0070:  6f70 6c61 7374 0a                        oplast.
...

访问页面可见是明文
```

