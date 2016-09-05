#18_02_ssl协议、openssl及创建私有CA

---

###笔记

---

####用户私钥丢失,证书失效

`CA` = `C`ertificate `A`uthority

`证书颁发机构`(`CA`)还要实现`证书吊销列表`(`CRL`),保存此前曾经发出去的证书,仍未过期但是已经被撤销了.

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
	
`openssl 加密解密示例`:

使用`des3算法` `加密`文件`xxxxx`

		openssl enc -des3 -salt -a -in xxxx -out xxxx.des3		enter des-ede3-cbc encryption password:
		Verifying - enter des-ede3-cbc encryption password:
		
使用`des3算法` `解密`文件`xxxx.des3(我们刚加密过的)`	
		
		openssl enc -des3 -d -salt -a -in xxx.des3 -out xxx
		enter des-ede3-cbc decryption password:
	
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

####openssl实现私有CA

1. 生成一对秘钥
2. 生成自签署证书

私钥权限`600`

* `openssl genrsa -out /PATH/TO/KEYFILENAME NUMBITS`: 生成私钥
* `openssl rsa -in /PATH/TO/KEYFILENAME -pubout`: 输出公钥
* `openssl req -new -x509 -key 秘钥文件 -out 证书保存位置  -days 保存天数`: 生成证书
* `openssl x509 -text -in 证书`: 读取证书

**生成私钥**


		[chloroplast@iZ94jqmwawyZ ~]$ openssl genrsa
		Generating RSA private key, 1024 bit long modulus
		..............++++++
		........................................++++++
		e is 65537 (0x10001)
		-----BEGIN RSA PRIVATE KEY-----
		MIICXgIBAAKBgQDNs7x0gIHKOBz8lIFY1/2esTXWHi7fzAmi4+93iPMBZ1CzG4V1
		fYIpU5n1C7zyM4UCq6h41TN6TjeUlscTGxQjv2ohkDVtce+kuEDE2TIb18BmGgSS
		MwBi/b3TtWnHcEhaJyXNgYCpfyymIZwd2kYtR3x6u9f5YwFGMPIKtHzTywIDAQAB
		AoGAUaj1uoY9gCrQjxDhXIS6YWJWTf9DeoLEnI7CRQDv/3GlXsUhMSg3IPLYXqhf
		RQNg3VOKGRYCTp54gBtvQk1wq525JHPkqA/QMNkGnw4uHanBCsMjPo+fDM9m79dc
		5X8UIE3KQtJLsdzyVVmqBYF9NwKcblbnXK+bgbzZ8AfDvZECQQD1XAjDNHzFTxkI
		MrGySSqycWw3G4nXrHWbGZ42aLW54Un+RQRD+TdEexlIGGzdRGJiuFjhoy82vzo/
		4DjsdjUJAkEA1p9thTtPTuGyoZYCfsvvnzIuFZWz+P5WT8f72iHlWG/drEMl0559
		QNdpF5btJLF+2p325hDmAn8PMfIyHXnrMwJBALFL8PUFr4dwUblP0IHxRw4s0bK8
		jo2vjEgoaeANKAwKlMpNGvj3VA2DGlCzfa8iJCoL5gYeQhbAdhoEL34HKOkCQQC3
		r3cyohI8dtpFhXfZQX1yCKZ8fsWrgzIn0gbxKDV7vTJBrq5/MZQNnM8rC1cnImpp
		fOzE9w2EcW511s2hgKkVAkEAmEgn7rHnRzN3bmsDTm8UeLhrp6UV6NYJL1Uydb/r
		PxRzMghGTo388uIXZi9M/XAikN/EsBDYCO/7tOFAA6/8XA==
		-----END RSA PRIVATE KEY-----

使用输出重定向保存结果.
		
		[chloroplast@iZ94jqmwawyZ ~]$ openssl genrsa > server.key
		Generating RSA private key, 1024 bit long modulus
		......++++++
		.................++++++
		e is 65537 (0x10001)
		
`()`代表在子shell中执行,创建`600`权限的文件.		
		
		(umask 077; openssl genrsa -out server1024.key 1024)		
**输出公钥**

`-pubout`,从私钥当中把公钥提取出来.

		openssl rsa -in server1024.key -pubout		 	
		[chloroplast@iZ94jqmwawyZ ~]$ openssl rsa -in server1024.key -pubout
		writing RSA key
		-----BEGIN PUBLIC KEY-----
		MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCxCbO3PDNADoMI23XzXWlg2cEx
		08jVtrKbToKqO1f/C4gpfIhiEjTPUs/U4+0SCMBJ/CtqmZcDe4S65eiZ4R9D4gP1
		y60ZqWgzaXTKomXsxvLGN8j5flLt4yRrceMD2Ah/xafGKircGCDJtzm/ucJOHnOD
		dVKru0IlyJErLv2TswIDAQAB
		-----END PUBLIC KEY-----

**生成自签署证书**

`req`: 生成证书的工具.

		openssl req -new -x509 -key 秘钥文件 -out 证书保存位置  -days 保存天数
		
		[chloroplast@iZ94jqmwawyZ ~]$ openssl req -new -x509 -key server1024.key -out server.crt -days 365
		You are about to be asked to enter information that will be incorporated
		into your certificate request.
		What you are about to enter is what is called a Distinguished Name or a DN.
		There are quite a few fields but you can leave some blank
		For some fields there will be a default value,
		If you enter '.', the field will be left blank.
		-----
		Country Name (2 letter code) [XX]: CN
		State or Province Name (full name) []:SHANNXI
		Locality Name (eg, city) [Default City]:XIAN
		Organization Name (eg, company) [Default Company Ltd]:marmot
		Organizational Unit Name (eg, section) []:Tech
		Common Name (eg, your name or your server's hostname) []:ca.marmot.com
		Email Address []:41893204@qq.com

`Common Name (eg, your name or your server's hostname) []`比较重要.

		[chloroplast@iZ94jqmwawyZ ~]$ cat server.crt
		-----BEGIN CERTIFICATE-----
		MIIC3DCCAkWgAwIBAgIJAIdQsJG0CSWjMA0GCSqGSIb3DQEBCwUAMIGGMQswCQYD
		VQQGEwJDTjEQMA4GA1UECAwHU0hBTk5YSTENMAsGA1UEBwwEWElBTjEPMA0GA1UE
		CgwGbWFybW90MQ0wCwYDVQQLDARUZWNoMRYwFAYDVQQDDA1jYS5tYXJtb3QuY29t
		MR4wHAYJKoZIhvcNAQkBFg80MTg5MzIwNEBxcS5jb20wHhcNMTYwOTA1MDYwMzQy
		WhcNMTcwOTA1MDYwMzQyWjCBhjELMAkGA1UEBhMCQ04xEDAOBgNVBAgMB1NIQU5O
		WEkxDTALBgNVBAcMBFhJQU4xDzANBgNVBAoMBm1hcm1vdDENMAsGA1UECwwEVGVj
		aDEWMBQGA1UEAwwNY2EubWFybW90LmNvbTEeMBwGCSqGSIb3DQEJARYPNDE4OTMy
		MDRAcXEuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCxCbO3PDNADoMI
		23XzXWlg2cEx08jVtrKbToKqO1f/C4gpfIhiEjTPUs/U4+0SCMBJ/CtqmZcDe4S6
		5eiZ4R9D4gP1y60ZqWgzaXTKomXsxvLGN8j5flLt4yRrceMD2Ah/xafGKircGCDJ
		tzm/ucJOHnODdVKru0IlyJErLv2TswIDAQABo1AwTjAdBgNVHQ4EFgQUbhGVJDzJ
		FQtySUe4pb7QC+fBg4kwHwYDVR0jBBgwFoAUbhGVJDzJFQtySUe4pb7QC+fBg4kw
		DAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOBgQAgWSncOJn4CXxP62j1GUiv
		a/NEk7Or1nz4uv1hb9o/E2ZIlz/7Ix0HQzFoMwv1xyEdtXgx43kIaojrU563mtng
		XB3LK+x5DJWJxHd74YZ+MAmAI+aFY9c8M/y3Khd4UodR1gdryVNmC+iTqOnUAXtB
		IGg5AX5w+btDx+LR+ZSHxA==
		-----END CERTIFICATE-----

**根据配置文件生成自签署证书**

配置文件: `/etc/pki/tls/openssl.cnf`

`生成证书`:

证书名字叫`cakey.pem`是因为配置文件的要求.

		[root@iZ94jqmwawyZ private]# pwd
		/etc/pki/CA
		[root@iZ94jqmwawyZ private]# (umask 077; openssl genrsa -out private/cakey.pem 2048)
		Generating RSA private key, 2048 bit long modulus
		............................+++
		...........................+++
		e is 65537 (0x10001)
		[root@iZ94jqmwawyZ private]# ls
		cakey.pem
				
`生成子签证书`:

		[root@iZ94jqmwawyZ private]# openssl req -new -x509 -key private/cakey.pem -out cacert.pem
		You are about to be asked to enter information that will be incorporated
		into your certificate request.
		What you are about to enter is what is called a Distinguished Name or a DN.
		There are quite a few fields but you can leave some blank
		For some fields there will be a default value,
		If you enter '.', the field will be left blank.
		-----
		Country Name (2 letter code) [XX]:CN
		State or Province Name (full name) []:SHANXI
		Locality Name (eg, city) [Default City]:XIAN
		Organization Name (eg, company) [Default Company Ltd]:marmot
		Organizational Unit Name (eg, section) []:Tech
		Common Name (eg, your name or your server's hostname) []:ca.marmot.com
		Email Address []:41893204@qq.com		
		[root@iZ94jqmwawyZ CA]# ls
		cacert.pem  certs  crl  newcerts  private

`生成目录`:
		
		确认有目录 certs crl newcerts
		
`创建文件`:

		根据配置文件
		database	= $dir/index.txt	# database index file.
		
		touch index.txt
		
		根据配置文件
		serial		= $dir/serial 		# The current serial number
		
		touch serial
		起始号
		echo 01 > serial
		
		[root@iZ94jqmwawyZ CA]# ls
		cacert.pem  certs  crl  index.txt  newcerts  private  serial

**申请证书**

我们先生成自己的私钥,让CA签署

		[root@iZ94jqmwawyZ httpd]# mkdir ssl
		[root@iZ94jqmwawyZ httpd]# cd ssl
		[root@iZ94jqmwawyZ ssl]# pwd
		/etc/httpd/ssl
		
		[root@iZ94jqmwawyZ ssl]# (umask 077; openssl genrsa -out httpd.key 1024)
		Generating RSA private key, 1024 bit long modulus
		................++++++
		.....++++++
		e is 65537 (0x10001)
		
		申请证书(csr: Certificate Signature Request 证书签字申请)
		私钥证书要和当初CA的证书要保持到组织机构名等一致,否则不能签署
		
		[root@iZ94jqmwawyZ ssl]# openssl req -new -key httpd.key -out httpd.csr
		You are about to be asked to enter information that will be incorporated
		into your certificate request.
		What you are about to enter is what is called a Distinguished Name or a DN.
		There are quite a few fields but you can leave some blank
		For some fields there will be a default value,
		If you enter '.', the field will be left blank.
		-----
		Country Name (2 letter code) [XX]:CN
		State or Province Name (full name) []:SHANXI
		Locality Name (eg, city) [Default City]:XIAN
		Organization Name (eg, company) [Default Company Ltd]:marmot
		Organizational Unit Name (eg, section) []:Tech
		Common Name (eg, your name or your server's hostname) []:www.test.com
		Email Address []:admin@qq.com
		
		Please enter the following 'extra' attributes
		to be sent with your certificate request
		A challenge password []:
		An optional company name []:
		
		让CA签署,如果不是同一台主机(我们这里测试是同一台主机)需要发给CA
		[root@iZ94jqmwawyZ ssl]# openssl ca -in httpd.csr -out httpd.crt -days 365
		Using configuration from /etc/pki/tls/openssl.cnf
		Check that the request matches the signature
		Signature ok
		Certificate Details:
		        Serial Number: 1 (0x1)
		        Validity
		            Not Before: Sep  5 06:45:24 2016 GMT
		            Not After : Sep  5 06:45:24 2017 GMT
		        Subject:
		            countryName               = CN
		            stateOrProvinceName       = SHANXI
		            organizationName          = marmot
		            organizationalUnitName    = Tech
		            commonName                = www.test.com
		            emailAddress              = admin@qq.com
		        X509v3 extensions:
		            X509v3 Basic Constraints:
		                CA:FALSE
		            Netscape Comment:
		                OpenSSL Generated Certificate
		            X509v3 Subject Key Identifier:
		                B4:4C:80:96:19:58:16:A4:66:2C:A6:15:EF:CC:F0:37:0E:4E:21:3C
		            X509v3 Authority Key Identifier:
		                keyid:CF:D2:8B:89:CF:74:BE:AF:3D:98:75:54:50:A6:65:CD:DA:B7:D5:D8
		
		Certificate is to be certified until Sep  5 06:45:24 2017 GMT (365 days)
		Sign the certificate? [y/n]:y


		1 out of 1 certificate requests certified, commit? [y/n]y
		Write out database with 1 new entries
		Data Base Updated


		[root@iZ94jqmwawyZ ssl]# cd /etc/pki/CA
		
		我们数据库已经更新,签署了证书, 01 这个数字是从serial当中读取出来的
		[root@iZ94jqmwawyZ CA]# cat index.txt
		V	170905064524Z		01	unknown	/C=CN/ST=SHANXI/O=marmot/OU=Tech/CN=www.test.com/emailAddress=admin@qq.com
		
		序列数字已经从01涨为02
		[root@iZ94jqmwawyZ CA]# cat serial
		02
		
**测试证书**

有`Makefile`,通过`make`命令生成测试证书.

		[root@iZ94jqmwawyZ CA]# cd /etc/pki/tls/certs/
		[root@iZ94jqmwawyZ certs]# ls
		ca-bundle.crt  ca-bundle.trust.crt  make-dummy-cert  Makefile  renew-dummy-cert		
		[root@iZ94jqmwawyZ certs]# make httpd.pem
		umask 77 ; \
		PEM1=`/bin/mktemp /tmp/openssl.XXXXXX` ; \
		PEM2=`/bin/mktemp /tmp/openssl.XXXXXX` ; \
		/usr/bin/openssl req -utf8 -newkey rsa:2048 -keyout $PEM1 -nodes -x509 -days 365 -out $PEM2 -set_serial 0 ; \
		cat $PEM1 >  httpd.pem ; \
		echo ""    >> httpd.pem ; \
		cat $PEM2 >> httpd.pem ; \
		rm -f $PEM1 $PEM2
		Generating a 2048 bit RSA private key
		..............................................+++
		.....................................................................................+++
		writing new private key to '/tmp/openssl.VBYdvJ'
		-----
		You are about to be asked to enter information that will be incorporated
		into your certificate request.
		What you are about to enter is what is called a Distinguished Name or a DN.
		There are quite a few fields but you can leave some blank
		For some fields there will be a default value,
		If you enter '.', the field will be left blank.
		-----
		Country Name (2 letter code) [XX]:CN
		State or Province Name (full name) []:SHANXI
		Locality Name (eg, city) [Default City]:XIAN
		Organization Name (eg, company) [Default Company Ltd]:marmo
		Organizational Unit Name (eg, section) []:Tech
		Common Name (eg, your name or your server's hostname) []:ca.test.com
		Email Address []:41893204@qq.com
		[root@iZ94jqmwawyZ certs]# ls
		ca-bundle.crt  ca-bundle.trust.crt  httpd.pem  make-dummy-cert  Makefile  renew-dummy-cert
		
		
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

`p` 就是 `POSIX` 的意思

####公钥由私钥生成

####`()`子shell