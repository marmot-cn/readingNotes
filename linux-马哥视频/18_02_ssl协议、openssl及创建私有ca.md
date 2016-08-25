#18_02_ssl协议、openssl及创建私有CA

---

####用户私钥丢失,证书失效

`CA` = `C`ertificate `A`uthority

证书颁发机构(`CA`)还要实现证书吊销列表(`CRL`),保存此前曾经发出去的证书,仍未过期但是已经被撤销了.

`B`从`A`获取到证书后,也需要去`CRL`检查该证书是否存在.

`PKI`(`P`ublic `K`ey `I`nfrastructure)是一种遵循标准的利用公钥加密技术为电子商务的开展提供一套安全基础平台的技术和规范.

####证书格式

`x509`,`pkcs12`

`x509`(比较常见,标准):

* 公钥
* 有效期限
* 证书的合法拥有者
* 证书该如何被使用(用法限定,比如:`通信`过程,`加密`(用作加密解密的算法),`签名`等)
* `CA`的信息,证书颁发机构自身的信息
* `CA`签名的校验码(`CA`的签名)

**`TLS/SSL`**

使用的是`x509`证书.

`OpenGPG`: `PKI`另外的一种实现.

`PKI`有2种实现:

* `TLS`
* `OpenGPG`

**`TLS/SSL` Handshake**

`OSI`七层模型:

* 应用层
* 表示层
* 会话层
* 传输层
* 网络层
* 数据链路层
* 物理层

`TCP/IP`四层模型:

* 应用层
* `SSL`(库),半个层.
* 传输层
* 网络层
* 物理层

`https`使用的是`443`端口,套接字实现和http不相同.

在应用层和传输层之间引入半个层`SSL`(库),让应用层某种协议传输数据到`TCP`层之间,调用了`SSL`的功能,这个协议就能实现`加密`的功能.

`SSL` = `S`ecure + `S`ocket + `L`ayer

所以:

* http + ssl = https
* smtp + ssl = smtps
* ftp + ssl = ftps

现在比较流行使用的是`SSLv2`和`SSLv3`.

`http`和`https`是2种不同的协议.

`TLS`:`T`ransport `L`ayer `S`ecurity 传输层安全

* `TLSv1` 相当于 `SSLv3`

**http(tcp)**

`http` 是基于 `tcp` 协议的, 双方在建立会话之前要先3次握手.

`https`在`tcp`握手完成之后:

* ...三次握手
* 客户端向服务器端发起请求
* 服务端和客户端互相协商建立`SSL`会话(选择使用哪一种加密协议(`SSLv2`,`SSLv3`,`TLSv1`),还要协商一个对称加密的`算法`)
* 服务端将自己的证书发给客户端
* 客户端获取证书后要验证证书
* 客户端将会建立(生成)一个会话秘钥
	* https传输数据通过加密形式传输实现的(对称加密,非对称太慢) 
	* 客户端通过自己的随机数,生成一个随机的堆成秘钥,通过`server`端的公钥加密后传输给`server`端
	* 不是通过`DH`交换协议,是客户端发送给服务端
* 服务端拿秘钥加密数据传输

####实现算法的工具

Linux上不同的加密机制的工具不一致

**OpenSSL**

`OpenSSL`: `SSL`的开源实现

有三部分组成:

* `libcrypto`: 通用加密库,提供各种加密函数
* `libssl`: TLS/SSL的实现
	* 基于`会话`的,实现了身份认证(传输证书),数据机密性和会话完整性的`TLS/SSL`库文件 
* `openssl`: 多用途命令行工具	 	
	* 单向加密
	* 对称加密
	* 非堆成加密
	* 模拟实现私有证书颁发机构 	
	* `子命令`:
		* `Standard commands`: 标准命令
			* `speed`: 测试不同算法在主机的加密速度	
		* `Message Digest commands`: 信息摘要算法的命令,单向加密
		* `Cipher commands`: 加密命令
	
**私有证书颁发机构**

公司内部实现加密解密,不和外人通信.
		
		//确认安装了 openssl
		rpm -q openssl
		openssl-1.0.1e-51.el7_2.5.x86_64	
		
		//检查openssl安装了哪些文件
		rpm -ql openssl
		
		/etc/pki/CA
		/etc/pki/CA/certs
		/etc/pki/CA/crl
		...
		
		

###整理知识点

---

####whatis

whatis命令是用于查询一个命令执行什么功能,并将查询结果打印到终端上.

`whatis` = `man -f`

		[root@iZ94jqmwawyZ ~]# whatis cp		
		cp (1)               - copy files and directories
		cp (1p)              - copy files


其中的`1`和`1p`相当于man手册的section

		man 1 cp
		User Commands
		....
		
		
		man 1p cp
		POSIX Programmer's Manual
		...

####man手册

1. User Commands(用户命令和守护进程)
2. System Calls(系统调用和内核服务)
3. C Library Functions
4. Devices and Special Files(特殊文件,设备驱动程序和硬件)
5. File Formats and Conventions(配置文件)
6. Games et. Al(游戏)
7. Miscellanea(杂项命令)
8. System Administration tools and Deamons(管理命令和守护进程)

`p`常见于 opensuse 的发行版

`0p`: POSIX headers 
`1p`: POSIX utilities 
`3p`: POSIX functions

p 就是 `POSIX` 的意思