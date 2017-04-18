# 25_03_邮件服务系列之Postifx常用配置

---

### 笔记

---

#### smtp

smtp 协议,明文传输数据.

功能增强`ESMTP`,也是明文传输.

升级`ssl`为加密.

smtp 只能用于邮件传输, pop3(邮局协议,3是版本号)协议用于接收邮件.

IMAP4: Internet Mail Access Protocol

`SASL`(让smtp实现用户认证):简单认证安全层, Simple Authentication Secure Layer

* v1
* v2

MDA: 邮件投递代理.

* procmail
* maildrop

MUA: 邮件用户代理

* mutt, linux下文本界面
* mail命令可以可以

		tom@a.org --> c.com(MX) --> jerry@b.net
		本来应该发到 b.net(MX记录),但是提交到 c.com(MX)
		
		c.com 涉及到邮件中继,是否中继取决于c.com是否开放中继.
		b.net 不涉及到邮件中继

MTA:

* sendmail
* qmail
* postfix,模块化设计
* exim

#### postfix 配置文件关键参数

**myhostname**

`myhostname=xxx`

当前邮件服务器自己的主机名称,也就是`hostname(命令)`的结果,最好保持一致.

**myorigin**

起源,邮件地址伪装(重写).

		myorigin=magedu.com
		
		root@mail.magedu.com
		=> 改写为 root@magedu.com 所有称为邮件地址伪装

**mydomain**

自己当前主机所在的域,如果没有定义`mydomian`,则把`myhostname`的第一段去了,剩下的当做本机的域名(假如 myhostname=mail.magedu.com,去掉第一段后为 magedu.com)

		mydomain=magedu.com
		
**mydestination**

所有目标收件人地址@后面的为mydestination所标识的,就表示邮件服务器就是目的地.否则就是中继

		mydestination= $参数名称,来引用响应参数的值
		
		mydestination = $myhostname, localhost.$mydomiain, localhost, $mydomain
		
**mynetworks**

允许中继的地址
		
		mynetwork = 192.168.1.0/24 127.0.0.0/8		    这几个网段可以中继

### 整理知识点

---