# openssl 命令 (私有CA办法证书)

---

#### 流程图

![openssl-ca](./img/openssl-ca.jpg)

#### 示例

#### 修改配置文件

```shell
cat  /etc/pki/tls/openssl.cnf
 

##### 创建私有证书签发机构CA步骤

```shell
cat  /etc/pki/tls/openssl.cnf
...
dir		= ./demoCA		# TSA root directory

读取配置文件,知道默认路径是 ./demoCA

[root@rancher-server ansible]# (umask 077; openssl genrsa -out /etc/pki/CA/private/cakey.pem 4096)
Generating RSA private key, 4096 bit long modulus
.................................................................................................................++
................................................................................................................................................................................................................................................++
e is 65537 (0x10001)

[root@rancher-server ansible]# ll -s /etc/pki/CA/private/
total 4
4 -rw------- 1 root root 3243 Jul  4 17:05 cakey.pem
[root@rancher-server ansible]# ll -s /etc/pki/CA
total 16
...
4 drwx------. 2 root root 4096 Jul  4 17:05 private
```

`umask 077`, 建立的文件默认权限是 `-rw-------`, 文件夹权限 `drwx------`.

`umask` 是从权限中"拿走"相应的位即可.


* r 读取: 4
* w 写入: 2
* x 执行: 1

umask  | 文件  | 目录
-------| -----| -------------
0  	   | 6 	  | 	7
1  	   | 6 	  | 	6
2  	   | 4 	  | 	5
3  	   | 4 	  | 	4
4  	   | 2 	  | 	3
5  	   | 2 	  | 	3
6  	   | 0 	  | 	1
7  	   | 0 	  | 	0

综上所述我们创建的文件时 属主具有 读写权限. 文件夹 属主具有 读写执行权限.

**从私钥创建公钥**

```shell
openssl rsa -in /etc/pki/CA/private/cakey.pem -pubout
```

##### 生成CA自签证书

```shell
[root@rancher-server ansible]# openssl req -new  -x509 -key /etc/pki/CA/private/cakey.pem -out /etc/pki/CA/cacert.pem -days 3655
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:CN 			# 两个字符表示的国家代码, CN为中国
State or Province Name (full name) []:Shannxi	# 省或洲的完整名称
Locality Name (eg, city) [Default City]:xi'an	# 所在位置的名称(默认为城市)
Organization Name (eg, company) [Default Company Ltd]:marmot  # 组织机构名称(默认为公司)
Organizational Unit Name (eg, section) []:Tech # 组织机构单元名称(eg.部门)
Common Name (eg, your name or your server's hostname) []:www.marmot.com # 持有者名或者所在服务器主机名(即域名)
Email Address []:41893204@qq.com # 管理员邮件地址,可以省略
```

* `-new`: 生成新证书签署请求.
* `-x509`: 生成自签格式证书,专用于创建私有CA时
* `-key`: 生成请求时用到的私有文件路径.
* `-out`: 生成的请求文件路径,如果自签操作将直接生成签署过的证书
* `-days`: 证书的有效时长,单位是day

`-key /etc/pki/CA/private/cakey.pem` 指明的是私钥的位置,会自动抽取出私钥中的公钥.

**pem文件**

Openssl使用`PEM`(Privacy Enhanced Mail)格式来存放各种信息,它是 openssl 默认采用的信息存放方式.

它的形式为“`——-BEGIN XXXX ——`",与结尾的“`——END XXXX——`"对应.

使用PEM格式存储的证书:

```shell
—–BEGIN CERTIFICATE—–
MIIF5TCCA82gAwIBAgIJAJUN43kzLnaAMA0GCSqGSIb3DQEBCwUAMIGIMQswCQYD
...
-----END CERTIFICATE-----
```

使用PEM格式存储的私钥:

```shell
-----BEGIN RSA PRIVATE KEY-----
MIIJKQIBAAKCAgEAvgVMq9oCogd6FbLjJFsa6XqVi20GkcW24k1JlIwowWlujgAE
...
-----END RSA PRIVATE KEY-----
```

##### 为CA提供所需的目录及文件

```shell
创建签发证书, 吊销证书, 新证书目录
mkdir  -pv  /etc/pki/CA/{certs,crl,newcerts}
```

```shell
创建证书序列号文件, 证书索引文件
touch  /etc/pki/CA/{serial,index.txt}

第一次创建的时候需要给予证书序列号
echo  01 > /etc/pki/CA/serial          
```

##### 服务申请证书签署实现SSL安全通信

要用到证书进行安全通信的服务器, 需要向CA请求签署证书.

需要签署的服务==无需==和CA证书签署机构主机在同一台服务器上.

**在用到证书的服务器生成私钥(非CA服务器)**

```shell
ssh ansible@120.25.87.35 -p 17456
...
切换 root

[root@rancher-agent-1 ansible]# mkdir -p /etc/httpd/ssl
[root@rancher-agent-1 ansible]# cd  /etc/httpd/ssl
生成私钥
[root@rancher-agent-1 ssl]# (umask  077; openssl  genrsa -out  /etc/httpd/ssl/httpd.key  2048)
```

**生成证书签署请求**

`csr` 是证书请求文件,用于申请证书. 在制作`csr`文件的时,必须使用自己的私钥来签署申,还可以设定一个密钥.

```shell
[root@rancher-agent-1 ssl]#openssl  req  -new  -key /etc/httpd/ssl/httpd.key  -out /etc/httpd/ssl/httpd.csr  -days  365
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:CN
State or Province Name (full name) []:Shannxi
Locality Name (eg, city) [Default City]:xi'an
Organization Name (eg, company) [Default Company Ltd]:marmot
Organizational Unit Name (eg, section) []:Tech
Common Name (eg, your name or your server's hostname) []:www.marmot.com
Email Address []:41893204@qq.com

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
```

* `.csr` 表示证书签署请求文件
* 要保证和签署机构CA签署机构信息一致

**将请求通过可靠方式发送给CA主机**

把`httpd.csr`发到`CA`所在服务器

```shell
scp -P 17456 /etc/httpd/ssl/httpd.csr ansible@10.44.88.189:
ssh: connect to host 10.44.88.189 port 22: Connection refused
lost connection
[root@rancher-agent-1 ssl]# scp -P 17456 /etc/httpd/ssl/httpd.csr ansible@10.44.88.189:
The authenticity of host '[10.44.88.189]:17456 ([10.44.88.189]:17456)' can't be established.
ECDSA key fingerprint is 1c:d1:06:ef:36:0d:6b:c7:72:51:73:b0:2d:38:02:b9.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '[10.44.88.189]:17456' (ECDSA) to the list of known hosts.
ansible@10.44.88.189's password:
httpd.csr
```

**在CA主机上签署证书**

`crt`是CA认证后的证书文件, 签署人用自己的key给你签署的凭证.

```shell
[root@rancher-server ansible]#openssl ca -in /home/ansible/httpd.csr -out /etc/pki/CA/certs/httpd.crt -days 365
Using configuration from /etc/pki/tls/openssl.cnf
Check that the request matches the signature
Signature ok
Certificate Details:
        Serial Number: 1 (0x1)
        Validity
            Not Before: Jul  4 10:43:50 2017 GMT
            Not After : Jul  4 10:43:50 2018 GMT
        Subject:
            countryName               = CN
            stateOrProvinceName       = Shannxi
            organizationName          = marmot
            organizationalUnitName    = Tech
            commonName                = www.marmot.com
            emailAddress              = 41893204@qq.com
        X509v3 extensions:
            X509v3 Basic Constraints:
                CA:FALSE
            Netscape Comment:
                OpenSSL Generated Certificate
            X509v3 Subject Key Identifier:
                E1:5C:D4:EA:9D:6B:2F:B9:16:6A:BC:5A:87:40:AB:40:8C:B1:FF:C0
            X509v3 Authority Key Identifier:
                keyid:DC:64:7F:4D:9E:DF:29:8B:CE:44:D7:07:4E:96:F7:AD:A9:F2:82:38

Certificate is to be certified until Jul  4 10:43:50 2018 GMT (365 days)
Sign the certificate? [y/n]:y


1 out of 1 certificate requests certified, commit? [y/n]y
Write out database with 1 new entries
Data Base Updated
```

**查看所签署的证书信息**

```shell
[root@rancher-server ansible]# cat  /etc/pki/CA/index.txt
V	180704104350Z		01	unknown	/C=CN/ST=Shannxi/O=marmot/OU=Tech/CN=www.marmot.com/emailAddress=41893204@qq.com
```

* `V` 表示已经签署的
* `01` 表示证书序列号

**将CA签署机构的`.crt`证书发送给服务器**

```shell
[root@rancher-server ansible]#scp -P 17456 /etc/pki/CA/certs/httpd.crt ansible@10.170.148.109:
```

**删除服务器和CA主机上签署前的`*.csr`文件,确保安全**

```shell
申请证书主机
rm  -rf  /etc/httpd/ssl/httpd.csr
```

```shell
CA主机
rm  -rf /home/ansible/httpd.csr
```

#### openssl 加密解密过程

![ssl](./img/ssl.jpg)

#### ca工作流程

![ca](./img/ca.jpg)
