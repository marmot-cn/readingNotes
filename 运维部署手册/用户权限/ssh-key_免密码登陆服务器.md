#SSH-KEY 免密码登陆服务器

####在本地服务器生成密钥和公钥

`openssh`的`ssh-keygen`命令用来产生这样的私钥和公钥

**一些参数说明**

* `-t dsa`: 采用dsa加密方式的公钥/私钥对,除了dsa还有rsa方式,rsa方式最短不能小于768字节长度
* `-C`: 对这个公钥/私钥对的一个注释和说明,一般用所有人的邮件代替.可以省略不写
* `-b xx`: 设定xx字节的公钥/私钥对,本示例没用,最长4096字节,一般1024或2048就可以了,太长的话加密解密需要的时间也长.

**示例**
	
		ssh-keygen -t rsa -C "41893204@qq.com"

		Generating public/private rsa key pair.
		Enter file in which to save the key (/Users/chloroplast1983/.ssh/id_rsa):
		/Users/chloroplast1983/.ssh/id_rsa already exists.
		Overwrite (y/n)? y
		Enter passphrase (empty for no passphrase):
		Enter same passphrase again:
		Your identification has been saved in /Users/chloroplast1983/.ssh/id_rsa.
		Your public key has been saved in /Users/chloroplast1983/.ssh/id_rsa.pub.
		The key fingerprint is:
		SHA256:iBleI+2QDZCC/uVxPJlxpJDysndhDB+gJTCz1fm72go 41893204@qq.com
		The key's randomart image is:
		+---[RSA 2048]----+
		|.+o=o++ ..       |
		|o *.+X.o..       |
		|.o .B @.*        |
		| . o % #         |
		|  . O * S        |
		|   o o o         |
		|   E. . .        |
		|    . ..         |
		|     oo.         |
		+----[SHA256]-----+
		
`cat /Users/chloroplast1983/.ssh/id_rsa.pub`

		ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCu7IB+uJUiSIW2MLhRaJ0s3s852K/G/7y5Rfc0xTAYuGu2+5CAzEOkzlNsSbOLx/Lkg3V+Dy5dkUT8rH3UtOr/oN+Hg0H5XXLn4JPqnMvWaQuks5fe++dsXm4QJJ+DHywsJZjkpYqElKpDjb5bj6vCMiSZYhal4iQZcyJ4KBEuYrOAiP4dx9f7yjIW3AZmYEmT2doJo/SYd7jKufBgg33e+TFKzuGVQZlGV5TuOMPcUudZj5nJz7eNre3db8bIaFfi2c/qiRKZlzGNwpxQm7Io+Tl5yP1Y6GmqlxwC+cWkc4pn0m3tLKwWgAV7FH7GfMUG5ChZYCDrUH7HQgsfXTaJ 41893204@qq.com
		
`id_dsa.pub`id_dsa.pub是公钥

####设置目标服务器

如果没有新建文件和文件夹

		[chloroplast@iZ94ebqp9jtZ ~]$ pwd
		/home/chloroplast
		[chloroplast@iZ94ebqp9jtZ ~]$ mkdir .ssh
		[chloroplast@iZ94ebqp9jtZ ~]$ cd .ssh/
		[chloroplast@iZ94ebqp9jtZ .ssh]$ vi authorized_keys
		
把本地服务器(上小节生成的)的`.ssh/id_rsa.pub`内容复制到`authorized_keys`内

		[chloroplast@iZ94ebqp9jtZ .ssh]$ cat authorized_keys
		sh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCu7IB+uJUiSIW2MLhRaJ0s3s852K/G/7y5Rfc0xTAYuGu2+5CAzEOkzlNsSbOLx/Lkg3V+Dy5dkUT8rH3UtOr/oN+Hg0H5XXLn4JPqnMvWaQuks5fe++dsXm4QJJ+DHywsJZjkpYqElKpDjb5bj6vCMiSZYhal4iQZcyJ4KBEuYrOAiP4dx9f7yjIW3AZmYEmT2doJo/SYd7jKufBgg33e+TFKzuGVQZlGV5TuOMPcUudZj5nJz7eNre3db8bIaFfi2c/qiRKZlzGNwpxQm7Io+Tl5yP1Y6GmqlxwC+cWkc4pn0m3tLKwWgAV7FH7GfMUG5ChZYCDrUH7HQgsfXTaJ 41893204@qq.com
		
**修改权限**

为了安全起见，要保证.ssh和authorized_keys都只有用户自己有写权限，否则验证无效

		chmod 700 ~/.ssh
		chmod 600 ~/.ssh/authorized_keys		
		
**验证登录**

无需密码即可登录

		ssh chloroplast@120.25.161.1
		Last login: Sun Dec 13 12:59:55 2015 from 1.86.230.8
		
		Welcome to aliyun Elastic Compute Service!
		
		[chloroplast@iZ94ebqp9jtZ ~]$
