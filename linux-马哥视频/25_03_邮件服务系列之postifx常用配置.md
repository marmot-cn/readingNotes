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

MTA: 邮件传送代理

* sendmail
* qmail
* postfix,模块化设计
* exim

MRA: 邮件访问代理 

* cyrus-imap
* dovecot

dovecot依赖mysql客户端,支持4种协议(pop3, imap4, pops, imaps),有SASL认证能力(虚拟账户).支持两种邮箱格式.

pop3: 110端口/tcp
imap4: 143端口/tcp

全部是以明文方式工作.

#### 邮箱格式

* mbox: 一个文件存储所有邮件
* maildir: 一个文件存储一封邮件,所有邮件存储在一个目录中.

红帽邮件服务器默认是`mbox`.

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

如果目标收件人地址@后面的为mydestination所标识的,就表示邮件服务器就是目的地.否则就是中继

		mydestination= $参数名称,来引用响应参数的值
		
		mydestination = $myhostname, localhost.$mydomiain, localhost, $mydomain
		
**mynetworks**

允许中继的地址
		
		mynetwork = 192.168.1.0/24 127.0.0.0/8		    这几个网段可以中继

配合`smtpd_recipient_restrictions = permit_mynetworks`一起实现.

#### postfix 结合 SASL 用户认证

1. 启动 SASL 服务
		
		/etc/rc.d/init.d/saslauthd 服务脚本
			/etc/sysconfig/saslauthd 配置文件
			
			saslauthd -v: 显示当前主机saslauthd服务所支持的认证机制, 默认为pam(shadow 机制是去 /etc/shadow 文件认证)
			
			testsaslauthd -u 用户名 -p 密码
			测试认证功能是否正常
		

##### 实现postfix基于客户端的访问控制

**基于客户端的访问控制概览**

postfix内置了多种反垃圾邮件的机制,其中就包括"客户端"发送邮件机制,客户端判别机制可以设定一系列客户信息的判别条件:

* smtpd_client_restrictions
* smtpd_data_restrictions
* smtpd_helo_restrictions
* smtpd_recipient_restrictions
* smtpd_sender_restrictions

上面的每一项参数分别用于检查`SMTP会话过程中的特定阶段`,即客户端提供相应信息的阶段,如当客户端发起链接请求时,postfix就可以根据配置文件中定义的`smtpd_client_restrictions`参数来判别此客户端ip的访问权限.相应地,`smtpd_helo_restrictions`则用于根据用户的`helo`信息判别客户端的访问能力等等.

**默认配置如下**

		smtpd_client_restrictions = 
		smtpd_data_restrictions = 
		smtpd_end_of_data_restrictions = 
		smtpd_etrn_restrictions = 
		smtpd_helo_restrictions = 
		smtpd_recipient_restrictions = permit_mynetworks,reject_unauth_destination
		smtpd_sender_restrictions = 

		
		permit_mynetworks: 限制了只有mynetworks参数中定义的本地网络中的客户端才能通过postfix转发邮件,其他客户端则不被允许,从而关闭了开放中继（open relay）的功能.
		reject_unauth_destination: 拒绝未经认证无法到达的目标.

##### 检查表格式说明

`hash`类的检查表都是用类似如下格式:

		pattern		action
		
检查表文件中,空白行,仅包含空白字符串的行和以`#`开头的行都会被忽略.以空白字符开头后跟其他非空白字符的行会被认为是前一行的延续,是一行的组成部分.

**pattern**

`pattern`通常有两类地址: 邮件地址和主机名称/地址.

邮件地址的pattern格式如下:

* user@domain: 用于匹配指定邮件地址;
* domian.tld: 用于匹配以此域名作为邮件地址中的域名部分的所有邮件地址;
* user@: 用于匹配以此作为邮件地址中的用户名部分的所有邮件地址;

主机名称/地址的pattern格式如下:

* domain.tld 用于匹配指定域及其子域内的所有主机;
* .domian.tld 用于匹配指定域的子域内的所有主机;
* net 用于匹配特定的ip地址或网络内的所有主机;
* network/mask CIDR格式，匹配指定网络内的所有主机;

**action**

接受类的动作:

OK 接受其pattern匹配的邮件地址或主机名称/地址;
全部由数字组成的action 隐式表示ok

拒绝类的动作:

* 4NN	text
* 5NN	text
	
其中4NN类表示过一会儿重试;5NN类表示严重错误,将停止重试邮件发送;`421`和`521`对于`postfix`来说有特殊意义,尽量不要自定义两个代码

* REJECT optional text... 拒绝; text为可选信息;
* DEFER optional text... 拒绝; text为可选信息;

**查找表**

访问控制文件,定义谁可以访问谁不能访问.文件数据过多需要转成hash格式 --> /etc/postfix.access.db(二进制格式)

		/etc/postfix/access
		obama@aol.com	reject
		microsoft.com	ok

也可以放到mysql库内.

		smtpd_client_restrictions = check_client_access hash:/etc/postfix/access
		smtpd_helo_restrictions = check_helo_access mysql:/etc/postfix/mysql_user
		smtpd_recipient_restrictions =
		smtpd_sender_restrictions = 

##### smtp 工作流程

connection:(建立链接的过程): smtpd_client_restrictions
helo(指令): smtpd_helo_restrictions,限定只有符合条件的用户发送helo指令
mail from(指令): smtpd_sender_restrictions,限定只有符合条件的用户发送mail from 指令
rcpt to(指令): smtpd_recipient_restrictions,...
data(指令): smtpd_data_restrictions,...

所有错误信息都在`rcpt to`段显示是否允许发送.

限制`ip`地址,限制`发件人`

#### 示例

禁止172.16.100.66这台主机通过工作在172.16.100.1上的postfix服务发送邮件为例演示说明其实现过程.访问表使用hash的格式.

1. 编辑`/etc/postfix/access`文件,以之作为客户端检查的控制文件,定义如下一行:

		172.16.100.200		REJECT
		
2. 将此文件转换为`hash`格式

		postmap /etc/postfix/access

3. 配置postfix使用此文件对客户端进行检查

		编辑/etc/postfix/main.cf文件,添加如下参数:
		smtpd_client_restrictions = check_client_access hash:/etc/postfix/access
		
4. 让postfix重新载入配置文件即可进行发信控制效果测试了



### 整理知识点

---