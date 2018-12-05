# 第15讲 | HTTPS协议：点外卖的过程原来这么复杂

## 笔记

### 加密

* 对称加密
* 非对称加密

#### 对称加密

加密和解密使用的密钥是**相同**的.

#### 非对称加密

加密和解密使用的密钥是**不同**的.

* 公钥
	* 可以在互联网随意传播
* 私钥

* 客户端 -> 服务端
	* 用服务端的公钥加密
* 服务端 -> 客户端
	* 用客户端的公钥加密

### 数字证书

有权威部门证明, 权威部门颁发的称为**证书(Certificate)**.

* 公钥
* 证书所有者
* 证书发布机构
* 证书有效期

权威机构我们称为`CA`.

#### 示例

搭建了一个网站`cliu8site`.

创建私钥

```
openssl genrsa -out cliu8siteprivate.key 1024
```

创建对应的公钥

```
openssl rsa -in cliu8siteprivate.key -pubout -out cliu8sitepublic.pem
```

生成**证书请求**

```
openssl req -key cliu8siteprivate.key -new -out cliu8sitecertificate.req
```

将请求发给权威机构, 权威机构会给这个证书上一个章(**签名算法**). 使用`CA`的私钥保证是真的权威机构签名.

签名算法: 对信息做一个`Hash`计算, 得到一个`Hash`值, 这个过程是不可逆的, 也就是说无法通过`Hash`值得出原来的信息内容. 

签名

```
openssl x509 -req -in cliu8sitecertificate.req -CA cacertificate.pem -CAkey caprivate.key -out cliu8sitecertificate.pem
```

`CA`用自己的私钥给`cliu8site`网站的公钥签名.

查看证书内容

```
openssl x509 -in cliu8sitecertificate.pem -noout -text 
```

* `Issuer` 证书颁发机构
* `Subject` 证书颁发给谁
* `Validity` 证书期限
* `Public-key` 公钥内容
* `Signature Algorithm` 签名算法

这样访问的时候会从`cliu8site`等到一个证书(而不是直接一个公钥), 这个证书有个发布机构`CA`, 只要得到这个发布机构`CA`的公钥, 解密外面网站的签名, 如果解密成功了, `Hash`也对的上, 就说明这个外卖网站的公告没有啥问题.

#### CA公钥层层授信

`CA`的公钥需要更牛的`CA`给它签名, 然后形成`CA`证书. 如果想知道某个`CA`的证书是否可靠, 要看`CA`的上级证书的公钥, 能不能解开这个`CA`的签名. 类似`公安局 -> 市公安局`. 知道全球皆知的几个`CA`, 称为**`root CA`**.

**层层授信背书**

#### 自签名 

`Self-Signed Certificate`

### HTTPS 的工作模式

* 公钥私钥主要用于传输**对称加密的密钥**
* 大量数据的通信通过对称加密进行

![](./img/15_01.jpg)

* 客户端发送加密方法和客户端随机数`Client Hello`
* 服务端确认加密算法并发送服务端随机数`Server Hello`
* 服务端发送证书`Server Certificate`
* 服务端发送就这样了`Server Hello Done`
* 客户端审核证书
* 客户端计算随机数字`Pre-master`用证书中的公钥加密发给服务端`Cleint Key Exchange`
* 客户端 和 服务端 使用以下3个随机数生成相同的**对称密钥** 
	* 客户端随机数
	* 服务端随机数
	* `Pre-master`
* 客户端说以后都用这个对称密钥和加密算法通信`Change Cipher Spec`.
* 客户端将已经商定好的参数, 采用协商好的密钥加密算法进行加密通信`Encrypted Handshake Message`
* 服务端发送`Change Cipher Spec`, 没问题, 以后都采用协商的通信密钥和加密算法进行通信.
* 服务端发送`Encrypted Handshake Message`试试
* 双方握手结束, 通过对称密钥进行加密传输了

## 扩展