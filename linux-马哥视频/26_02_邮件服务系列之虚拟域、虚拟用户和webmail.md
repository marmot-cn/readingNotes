# 26_02_邮件服务系列之虚拟域、虚拟用户和webmail

---

### 笔记

---

#### ltdl

动态模块加载器(GNU Libtool Dynamic Mdoule Loader)

### dovecot 

宏:

* `%d` : 域名
* `%/n`: 用户名

**配置**

```shell
vim /etc/dovecot.conf
mail_location = maildir:/var/mailbox/%d/%n/Maildir
...
auth default {
	mechanisms = plain
	passdb sql {
		args = /etc/dovecot-mysql.conf
	}
	userdb sql {
		args = /etc/dovecot-mysql.conf
	}
	...
```

```shell
数据库 数据库用户名 密码 都是我们提前建立好的

vim /etc/dovecot-mysql.conf
driver = mysql
connect = host=locolost(若mysql.sockr不是默认的,可以使用host="sock文件的路径"来指定新位置) dbname=extmail user=extmail passowrd=extmail
defualt_pass_scheme = CRYPT
password_query = SELECT username AS user, password as password FROM mailbox WHERE username = '%u'
user_query = SELECT maildir, uidnumber as uid, gidnumber as gid FROM mailbox WHERE username = '%u'
```


### 整理知识点

---