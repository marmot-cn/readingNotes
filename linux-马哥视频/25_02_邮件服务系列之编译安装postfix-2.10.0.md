# 25_02_邮件服务系列之编译安装Postfix-2.10.0

---

### 笔记

---

#### MTA

邮件传输代理,统称. SMTP 服务器.

SMTP 是一种协议.

smtp + ssl = smtps 默认明文传输

**sendmail**

提供 MTA 协议的软件 `sendmail`.`sandmial`最早的协议是`UUCP`.

sendmail,单体结构.会用到 `suid`,配置文件语法复杂(m4编写).

**qmail**

qmail, 性能好,现在应用不多.

**postfix**

新贵, 采用模块化设计, 避免使用`suid`, 一开始就注入邮件服务器安全的观念, 安全. 跟`sendmail`兼容性好, 投递效率是`sendmail`的4倍.

**exim**

英国剑桥

**Exchange**

Windows用的邮件服务器, 异步消息协作平台. 必须根windows的AD整合. 重量级.

#### MDA

邮件投递代理.

postfix 自带投递代理 (本地用户邮件投递, 虚拟用户邮件投递), 也可以使用`maildrop`.

邮件投递代理程序:

* procmail
* maildrop

#### MRA

邮件检索代理, 用户链接检索邮件.

实现协议

* pop3, pop3+ssl = pop3 默认明文传输
* imap4, imaps

实现MRA的软件:

* cyrus-imap
* dovecot

#### MUA

邮件用户代理

* OutLook Express(简装版), Outlook(专业版)
* Foxmail
* Thunderbird
* Evolution
* mutt(文本界面)

#### Webmail

邮箱服务程序

* Openwebmail
* squirrelmail(red hat 自带)
* Extmail (Extman),自己在Centos基础上做了一个 EMOS,自带 Extmail

#### SASL

用户认证,目前用的 V2 版协议.

cyrus-sasl 提供 saslauthd.

SASL 是一个认证框架, 要真正实现认证还需要认证机制.

		[root@rancher-agent-1 ansible]# yum list all | grep sasl
		Repodata is over 2 weeks old. Install yum-cron? Or run: yum makecache fast
		cyrus-sasl-lib.x86_64                   2.1.26-20.el7_2                @base
		cyrus-sasl-plain.x86_64                 2.1.26-20.el7_2                @base
		cyrus-sasl.i686                         2.1.26-20.el7_2                base
		cyrus-sasl.x86_64                       2.1.26-20.el7_2                base
		cyrus-sasl-devel.i686                   2.1.26-20.el7_2                base
		cyrus-sasl-devel.x86_64                 2.1.26-20.el7_2                base
		cyrus-sasl-gs2.i686                     2.1.26-20.el7_2                base
		cyrus-sasl-gs2.x86_64                   2.1.26-20.el7_2                base
		cyrus-sasl-gssapi.i686                  2.1.26-20.el7_2                base
		cyrus-sasl-gssapi.x86_64                2.1.26-20.el7_2                base
		cyrus-sasl-ldap.i686                    2.1.26-20.el7_2                base
		cyrus-sasl-ldap.x86_64                  2.1.26-20.el7_2                base
		cyrus-sasl-lib.i686                     2.1.26-20.el7_2                base
		cyrus-sasl-md5.i686                     2.1.26-20.el7_2                base
		cyrus-sasl-md5.x86_64                   2.1.26-20.el7_2                base
		cyrus-sasl-ntlm.i686                    2.1.26-20.el7_2                base
		cyrus-sasl-ntlm.x86_64                  2.1.26-20.el7_2                base
		cyrus-sasl-plain.i686                   2.1.26-20.el7_2                base
		cyrus-sasl-scram.i686                   2.1.26-20.el7_2                base
		cyrus-sasl-scram.x86_64                 2.1.26-20.el7_2                base
		cyrus-sasl-sql.i686                     2.1.26-20.el7_2                base
		cyrus-sasl-sql.x86_64                   2.1.26-20.el7_2                base

courier-authlib 实现到数据库(mysql...)的认证, 实现虚拟用户.

#### 部署

**发邮件**

Postfix + SASL + Mysql(courier-authlib, 虚拟用户)

postfix: rpm 不实现虚拟用户.

**收邮件**

Dovecot + MySql

**Webmail**

Extmail+Extman+httpd

##### 安装postfix
		
		创建独立的用户名和组,并且不让其登陆(-s 新账户的登录 shell),不创建用户的主目录(-M)
		 
		groupadd -g 2525 postfix
		useradd -g postfix -u 2525 -s /sbin/nologin -M postfix
		groupadd -g 2526 postdrop
		useradd -g postdrop -u 2526 -s /sbin/nologn -M postdrop
		
		从www.postfix.org下载源码.
		tar xf postfix-2.10.0.tar.gz
		cd postfix-2.10.0
		编译方法比较独特,CCARGS 指定编译选项,到哪去找系统的头文件. AUXLIBS 辅助的库文件路径. -DUS 使用响应功能.
		make makefiles 'CCARGS=-DHAS_MYSQL -I/usr/local/mysql/include(根据自己的路径) -DUSE_SASL_AUTH(启用SASL认证) -DUSE_CYRUS_SASL -I/usr/include/sasl -DUSE_TLS(支持TLS,使用smtps协议)' 'AUXLIBS=-L/usr/local/mysql/lib -lmysqlclient -lz(压缩库文件) -lm(模块文件) -L/usr/lib/sasl2(检查 sasl dev包是否安装,yum list all | grep sasl) -lsasl2 -lssl -lcrypto'
		make
		make install
		
		按照以下的提示输入相关的路径([]中的是缺省值,"]"后的是输入值,省略的表示采用默认值)
		...
		
		生成别名二进制文件
		
		发邮件
		telnet localhost 25 (发送)
		220 localhost.localdomian ESMTP Postfix (响应)
		helo(命令就是少了一个l) localhost (发送)
		250 localhost.localdomian (响应)
		mail from:obama@whitehouse.com(随便填写发件人) (发送)
		250 2.1.0 Ok
		rcpt to:ansible (要先通过newaliases生成/etc/aliases.db) (发送)
		data
		354 End data with <CR><LF>. <CR><LF>
		Subject:How are you these datys?
		Are you gua le ma?
		. (一个空白行加一个点表示邮件结束)
		250 2.0.0 Ok: queued as B8A1D1E745F (响应,邮件已经进入发送队列)
		
		检查日志是否发送成功
		tail /var/log/maillog
		Apr 13 07:38:12 rancher-agent-1 postfix[27317]: fatal: unknown inet_protocols value "IPv4" in "IPv4"
		Apr 13 07:38:51 rancher-agent-1 postfix/postfix-script[27429]: starting the Postfix mail system
		Apr 13 07:38:51 rancher-agent-1 postfix/master[27431]: daemon started -- version 2.10.1, configuration /etc/postfix
		Apr 13 07:38:57 rancher-agent-1 postfix/smtpd[27447]: connect from localhost[127.0.0.1]
		Apr 13 07:39:26 rancher-agent-1 postfix/smtpd[27447]: 9D29380F3D: client=localhost[127.0.0.1]
		Apr 13 07:39:39 rancher-agent-1 postfix/cleanup[27452]: 9D29380F3D: message-id=<20170412233926.9D29380F3D@rancher-agent-1.localdomain>
		Apr 13 07:39:39 rancher-agent-1 postfix/qmgr[27433]: 9D29380F3D: from=<obama@whitehouse.com>, size=370, nrcpt=1 (queue active)
		Apr 13 07:39:40 rancher-agent-1 postfix/local[27453]: 9D29380F3D: to=<ansible@rancher-agent-1.localdomain>, orig_to=<ansible>, relay=local, delay=23, delays=23/0.04/0/0.01, dsn=2.0.0, status=sent (delivered to mailbox)
		Apr 13 07:39:40 rancher-agent-1 postfix/qmgr[27433]: 9D29380F3D: removed
		Apr 13 07:39:51 rancher-agent-1 postfix/smtpd[27447]: disconnect from localhost[127.0.0.1]
		
		登陆ansible用户查看邮件:
		[ansible@rancher-agent-1 ~]$ mail
		Heirloom Mail version 12.5 7/5/10.  Type ? for help.
		"/var/spool/mail/ansible": 1 message 1 new
		>N  1 obama@whitehouse.com  Thu Apr 13 07:39  14/523   "How are you these datys?"
		& 1
		Message  1:
		From obama@whitehouse.com  Thu Apr 13 07:39:39 2017
		Return-Path: <obama@whitehouse.com>
		X-Original-To: ansible
		Delivered-To: ansible@rancher-agent-1.localdomain
		Subject:How are you these datys?
		Date: Thu, 13 Apr 2017 07:39:16 +0800 (CST)
		From: obama@whitehouse.com
		Status: R
		
		Are you gua le ma?
		
		&		
		
**postfix的配置文件**

postfix 模块化:
	* master(主进程).`/etc/postfix/master.cf` 核心配置文件.
	* 辅助进程,用于完成不同的功能.`/etc/postfix/main.cf` 其他模块的配置文件.
			
			参数(必须顶格写,写在行的绝对行首,以空白开头的行被认为是上一行的延续) = 值			 
`postconf`,配置`postfix`的命令.

* `-d`: 显示默认配置选项	
* `-n`: 显示修改了的配置
* `-m`: 显示所有查找表类型
* `-A`: 显示支持的SASL客户端插件类型

		postconf -A
		cyrus
* `-e PARAMETER=VALUE`: 更改某参数配置信息,并保存至`main.cf`文件. 
	
#### smtp状态码

* 1xx: 纯信息
* 2xx: 正确
* 3xx: 上一步操作尚未完成, 仍需要继续补充
* 4xx: 暂时性错误
* 5xx: 永久性错误	
	
#### smtp协议源语	

* helo: (smtp协议发送hello信息)
* ehlo: (esmtp协议发送hello信息)	
* mail from: 指定发件人
* rcpt to: 指定收件人
	
#### alias: 邮件别名	
	
abc@magedu.com 都转发给 postmaster@magedu.com	

看上去是 abc@magedu.com 实际是给 postmaster@magedu.com 这就是邮件别名

`/etc/aliases` --> hash编码,转换为 --> /etc/aliases.db

通过 `newaliases` 命令转换上述hash.
	
#### postfix 默认中继

postfix 会默认把当前主机所在的 ip 地址所在的网段都认为是本地客户端, 所有本地客户端都会默认被中继.	
	
### 整理知识点

---

#### nslookup

#### txt记录

TXT记录一般是为某条记录设置说明,比如你新建了一条a.ezloo.com的TXT记录,TXT记录内容"this is a test TXT record.",然后你用 nslookup -qt=txt a.ezloo.com,你就能看到"this is a test TXT record"的字样.

**验证域名的所有**

除外,TXT还可以用来验证域名的所有,比如你的域名使用了Google的某项服务,Google会要求你建一个TXT记录,然后Google验证你对此域名是否具备管理权限.

在命令行下可以使用nslookup -q=txt a.ezloo.com来查看TXT记录。

#### AAAA记录

AAAA记录是一个指向IPv6地址的记录.

可以使用nslookup -q=aaaa a.ezloo.com来查看AAAA记录

**测试**
		
		我额外给 qixinyun.prototype.qixinyun.com	text 记录添加了一条 TXT记录 标注为 3213124123123

		[ansible@rancher-agent-1 ~]$ nslookup -q=TXT(txt也可以) qixinyun.prototype.qixinyun.com
		Server:		10.143.22.118
		Address:	10.143.22.118#53
		
		Non-authoritative answer:
		qixinyun.prototype.qixinyun.com	text = "3213124123123"
		
		Authoritative answers can be found from:
		qixinyun.com	nameserver = dns9.hichina.com.
		qixinyun.com	nameserver = dns10.hichina.com.
		dns9.hichina.com	internet address = 140.205.81.15
		dns9.hichina.com	internet address = 140.205.81.25
		dns9.hichina.com	internet address = 140.205.228.15
		dns9.hichina.com	internet address = 140.205.228.25
		dns9.hichina.com	internet address = 42.120.221.15
		dns9.hichina.com	internet address = 42.120.221.25
		dns10.hichina.com	internet address = 140.205.228.26
		dns10.hichina.com	internet address = 42.120.221.16
		dns10.hichina.com	internet address = 42.120.221.26
		dns10.hichina.com	internet address = 140.205.81.16
		dns10.hichina.com	internet address = 140.205.81.26
		dns10.hichina.com	internet address = 140.205.228.16
