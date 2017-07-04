# 26_03_邮件服务系列之pop3s、maildrop

---

### 笔记

---

postfix + sasl + courier-authlib + MySQL(实现了虚拟用户,虚拟域) + dovecot + 
Webmail (exmail,extmain)

#### 加密用户邮件传输

SMTP, POP3, IMAP4 都是明文传输的 `plaintext`, 纯文本.

http --> https (ssl/tls)

smtp(25) -> smtps (465端口)

**S/MIME**

Security MIME, 安全的多用途用户邮件扩展.

USER <-> USER

每个用户都有一个用户证书.

mail, hash(单向加密,计算特征码 finger print), 使用私钥加密指纹存放在邮件后面, 生成一次性会话密钥,加密这段数据. 再用对方的公钥加密这段密码.

工具: `OpenSSL`, `GPG`(PGP的一种实现).

PKI(`P`ublic `K`ey `I`nfrastructure)公钥基础设施: CA(`C`ertificate `A`uthority)认证中心.

发送前,投递之后即邮件本身是加密的. 到了用户本地在解密.

**pop3s, imaps**

邮件下载到本地是安全的.

和 smtps 一样, 不能保证邮件存储是加密的,仅仅是传输过程是加密的. 

#### `s_client`

测试建立`ssl`会话不能使用`telnet`.

```shell
openssl s_client -connect mail.magedu.com:995 -CA filte(验证服务器端证书)

进入类似 telnet 交互模式
```

#### tcpdump

协议报文分析器.

可以分析多种协议的多种报文格式. 可以抓包并解码.

`wireshark`(GUI) 图形化工具.

```shell
tcpdump [options] 过滤条件

获取报文的条件:

tcp src host 172.16.100.1
tcp src or dst port 21
```

使用非ssl,抓取到的报文是明文的.

**参数**

* -i: 指定在哪个网卡设备抓包
* -n: 不反解主机名
* -nn: 不反解主机名和端口号
* -X: 显示报文的内容以`16`(hex)进制和`ASCII`
* -XX: 同 `-X` 还会显示 ethernet hader (以太网首部)
* -V, -VV: 显示更详细信息

```shell
eth0 网口,端口 110
tcpdump -i eth0 -X -nn -vv tcp port 110


添加条件 源地址是 172.16.100.1
tcpdump -i eth0 -X -nn -vv tcp port 110 and ip src 172.16.100.1

添加条件 源地址,目标地址都是 172.16.100.1
tcpdump -i eth0 -X -nn -vv tcp port 110 and ip host 172.16.100.1
```

**tcpdump的语法**

```shell
tcpdump [options] [Protocol] [Direction] [Host(s)] [Value] [Logical Operations] [Other expression]
```

**Protocol协议**

Values(取值):

* ether
* fddi
* ip (三层)
* arp
* rarp
* decnet
* lat
* sca
* moprc
* mopdl
* tcp and udp (四层)

If no protocol is specified, all the protocols are used.

**Direction(流向)**

Values(取值):

* src
* dst
* src and dst
* src or dst

默认是"src or dst"

**Host(s)主机**

Values(取值):

* net
* port
* host
* portrange 端口范围

默认即`host`

```shell
src 10.1.1.1 等同于 src host 10.1.1.1

dst net 172.16.0.0./16 目标地址网络
```

**Logical Operations**

1. AND : and or &&
2. OR : or or ||
3. EXCEPT : not or !

**示例**

```shell
只分析 IP 报文, IP 源地址必须是 172.16.100.1
ip src 172.16.100.1 

只分析 IP 报文, IP 源地址必须是 172.16.100.1, 目标地址必须是 172.16.100.1
ip src and dst 172.16.100.1 

只分析 IP 报文, IP 源地址或目标地址 是 172.16.100.1
ip src or dst 172.16.100.1 
```

#### 垃圾邮件处理

**内容过滤器**

APACHE: Spamassassin Perl语言开发,垃圾邮件分拣器.根据其垃圾邮件特征码库进行分析.

**RBL**

`R`ealtime `B`lack `L`ist.

实时黑名单列表.

**关闭OpenRelay**

关闭开放中继.

#### 病毒邮件网关

**clamav**

开元杀毒软件. 做病毒邮件服务器网关.

#### 呼叫器

邮件服务器通过插件**呼叫器**调用垃圾邮件过滤器`Spamassassin`,病毒邮件网关`clamav`

caller

`MIMEDefang`, `Mailscanner`, `Amavisd-new`(新).

`Amavisd-new` 调用 `Spamassassin`,`clamav`.

#### maildrop

属于 MDA, Courier的组件. 负责邮件投递.

### 整理知识点

---