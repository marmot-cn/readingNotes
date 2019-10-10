# 30_03_nss&pam

---

### 笔记

---

#### 认证服务

* Authertication: 认证. 验证用户身份的手段和机制.
* Authorization: 授权. 认证完用户, 决定一个用户能否访问一个服务的过程.
* Audition: 审计

认证: 某个来获取资源的申请者, 是否他自己声称的用户.

#### 名称解析

`libnss`

操作系统识别用户 是靠 用户id, 识别用户组是靠 组id. 不是依靠用户名和组名.

需要将`username-->uid`, `groupname-->gid`. 需要将名称转换成数字的还有:`FQDN->IP`, `HTTP->80`.

如上整个过程叫做==名称解析==. 将用户识别的字符串, 转换成系统识别的`id`.

名称解析是==双向==的.

名称解析通常就是, 在一个数据库当中查找条目中对应的关系. 正向解析以名称为索引, 反向解析以数字为索引.

`FQND`->`IP`多种名称解析方式(两种存放解析库的不同手段).

* `DNS`
* `/etc/hosts`
* `mysql`
* `nis` 网络信息服务
* `ldap` 轻量级目录访问协议, 查找效率比`mysql`高.

##### login 程序示例

用户名 -> uid. 

早期`login`程序内置去`/etc/passwd`里面去找. 如果用户太多, 从文件查找效率太低. 如果需要换`ldap`, 如果`login`程序不支持, 则无法替换.

#### nsswitch

将一个应用程序的名称解析功能独立出去. 有一个专门的程序负责名称解析.

`App -> nsswitch -> resolve_lib`

`nsswitch`网络服务转换. 本身也是一堆库文件(把通用服务做成库, 谁用谁调用).

##### 配置文件`/etc/nsswitch.conf`

定义某种名称解析服务, 通过哪种手段来实现.

```shell
[root@localhost backend-crew]# cat /etc/nsswitch.conf
#
# /etc/nsswitch.conf
#
# An example Name Service Switch config file. This file should be
# sorted with the most-used services at the beginning.
#
# The entry '[NOTFOUND=return]' means that the search for an
# entry should stop if the search in the previous entry turned
# up nothing. Note that if the search failed due to some other reason
# (like no NIS server responding) then the search continues with the
# next entry.
#
# Valid entries include:
#
#	nisplus			Use NIS+ (NIS version 3)
#	nis			Use NIS (NIS version 2), also called YP
#	dns			Use DNS (Domain Name Service)
#	files			Use the local files
#	db			Use the local database (.db) files
#	compat			Use NIS on compat mode
#	hesiod			Use Hesiod for user lookups
#	[NOTFOUND=return]	Stop searching if not found so far
#

# To use db, put the "db" in front of "files" for entries you want to be
# looked up first in the databases
#
# Example:
#passwd:    db files nisplus nis
#shadow:    db files nisplus nis
#group:     db files nisplus nis

passwd:     files sss
shadow:     files sss
group:      files sss
#initgroups: files

#hosts:     db files nisplus nis dns
hosts:      files dns myhostname

# Example - obey only what nisplus tells us...
#services:   nisplus [NOTFOUND=return] files
#networks:   nisplus [NOTFOUND=return] files
#protocols:  nisplus [NOTFOUND=return] files
#rpc:        nisplus [NOTFOUND=return] files
#ethers:     nisplus [NOTFOUND=return] files
#netmasks:   nisplus [NOTFOUND=return] files

bootparams: nisplus [NOTFOUND=return] files

ethers:     files
netmasks:   files
networks:   files
protocols:  files
rpc:        files
services:   files sss

netgroup:   files sss

publickey:  nisplus

automount:  files sss
aliases:    files nisplus
```

`sss`: `sssd`是一款用以取代ldap和AD的软件.

参数示例:

```shell
解析passwd去文件里面找
passwd:file
解析主机名先去文件里面找, 找不见在去dns
hostname:file dns
```

```
protocols:  files
对应文件在
[root@localhost lib64]# cat /etc/protocols
# /etc/protocols:
# $Id: protocols,v 1.11 2011/05/03 14:45:40 ovasik Exp $
#
# Internet (IP) protocols
#
#	from: @(#)protocols	5.1 (Berkeley) 4/17/89
#
# Updated for NetBSD based on RFC 1340, Assigned Numbers (July 1992).
# Last IANA update included dated 2011-05-03
#
# See also http://www.iana.org/assignments/protocol-numbers

ip	0	IP		# internet protocol, pseudo protocol number
hopopt	0	HOPOPT		# hop-by-hop options for ipv6
icmp	1	ICMP		# internet control message protocol
igmp	2	IGMP		# internet group management protocol
...
emcon	14	EMCON		# EMCON
xnet	15	XNET		# Cross Net Debugger
chaos	16	CHAOS		# Chaos
udp	17	UDP		# user datagram protocol
mux	18	MUX		# Multiplexing protocol
dcn	19	DCN-MEAS		# DCN Measurement Subsystems
hmp	20	HMP		# host monitoring protocol
...

[root@localhost lib64]# cat /etc/services
....
```

`passd`定了去`file`找, 但是具体去文件查找的功能是在库文件里面实现的.

`.so`动态链接库.

```shell
[root@localhost lib64]# pwd
/usr/lib64
[root@localhost lib64]# ls | grep nss
libevent_openssl-2.0.so.5
libevent_openssl-2.0.so.5.1.9
libjansson.so.4
libjansson.so.4.4.0
libnss3.so
libnss_compat-2.17.so
libnss_compat.so
libnss_compat.so.2
libnss_db-2.17.so
libnss_db.so
libnss_db.so.2
libnss_dns-2.17.so
libnss_dns.so
libnss_dns.so.2
libnss_files-2.17.so
libnss_files.so
libnss_files.so.2
libnss_hesiod-2.17.so
libnss_hesiod.so
libnss_hesiod.so.2
libnss_myhostname.so.2
libnss_mymachines.so.2
libnss_nis-2.17.so
libnss_nis.so
libnss_nis.so.2
libnss_nisplus-2.17.so
libnss_nisplus.so
libnss_nisplus.so.2
libnss_sss.so.2
libnssckbi.so
libnssdbm3.chk
libnssdbm3.so
libnsspem.so
libnsssysinit.so
libnssutil3.so
libsss_nss_idmap.so.0
libsss_nss_idmap.so.0.1.0
nss
openssl
openssl098e
```

##### `[NOTFOUND=return]`

`[NOTFOUND=return]` 如果前面没有找见, 直接返回不在继续查找.

`hosts: files dns`

查找后的处理机制:

* `SUCCESS` 服务正常(比如上面的示例`files`文件存在), 并且找到了名称, 正常转换. 
* `NOTFOUND` 服务正常, 没有找到对应名称.
* `UNAVAIL` 服务不可用.
* `TRYAGAIN` 服务有临时性故障.

默认情况下, 出现`SUCCESS`就`return`(执行动作)(上面的例子查找主机名, 如果在文件中找见就不继续去dns查找), 否则就继续.

```
passwd: nis[NOTFOUND=return] files
```

passwd在`nis`中找不见就`return`, 如果`UNAVAIL`(`nis`不可用)就去文件查找.

##### `getent`

从某个管理库中获取相应条目.

```shell
passwd /etc/nsswitch.conf 里面定义的是文件, getent passwd 就在 /etc/passwd 文件里面查找.

[root@localhost lib64]# getent passwd
root:x:0:0:root:/root:/bin/bash
bin:x:1:1:bin:/bin:/sbin/nologin
daemon:x:2:2:daemon:/sbin:/sbin/nologin
adm:x:3:4:adm:/var/adm:/sbin/nologin
lp:x:4:7:lp:/var/spool/lpd:/sbin/nologin
sync:x:5:0:sync:/sbin:/bin/sync
shutdown:x:6:0:shutdown:/sbin:/sbin/shutdown
halt:x:7:0:halt:/sbin:/sbin/halt
...
查找特定用户
[root@localhost lib64]# getent passwd root
root:x:0:0:root:/root:/bin/bash

查找 hosts
[root@localhost lib64]# getent hosts
127.0.0.1       localhost localhost.localdomain localhost4 localhost4.localdomain4
127.0.0.1       localhost localhost.localdomain localhost6 localhost6.localdomain6
127.0.0.1       www.skynetcredit.com
[root@localhost lib64]# cat /etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
127.0.0.1   www.skynetcredit.com
查找具体域名
[root@localhost lib64]# getent hosts www.skynetcredit.com
127.0.0.1       www.skynetcredit.com
文件里面没有定义, 去DNS查找
[root@localhost lib64]# getent hosts www.baidu.com
220.181.111.188 www.a.shifen.com www.baidu.com
220.181.112.244 www.a.shifen.com www.baidu.com
```

#### pam 认证机制

`Pluggable Authentication Modules` 可插入的认证模块.

认证之前一般需要名称解析(虚拟用户就不需要解析).

认证 和 名称解析 是两套各自独立的机制.

名称解析: `shadow: files(/etc/shadow)`找密码是在`/etc/shadow`里面去找, 但是不管验证.

##### 操作系统用户登录过程

* 用户输入用户名`root` -> `nsswitch.conf` -> `passwd:files`(查找用户的id) 
* 用户输入密码`xxx` ->`nsswitch.conf` -> `shadow: files`
* 根据认证机制把`xxx + salt`计算出密码 -> `compare` (与`shadows`文件中的密码串比较)

认证本身也可以不用借助名称解析服务去查找用于原始存放的密码. `login`借用了认证机制找`shadow`, 很多服务的认证机制不依赖于名称解析.

##### Authentication

* md5: /etc/shadow
* mysql
* ldap
* nis
* kerberos

##### 中间层`pam`

`APP -> PAM -> Authenication`

`PAM` 本身不做认证, 是一个认证框架(Framework), 借助于对应服务的库文件来认证.

```shell
[root@localhost lib64]# ls /lib64/security/
pam_access.so    pam_echo.so       pam_fprintd.so ...
...
```

`pam_unix.so` 是实现传统意义的用户认证. `/etc/shadow`比对认证.

`pam_winbind.so` 到 `windows ad(活动目录)(ldap的服务器)`去认证

##### pam 的配置文件

配置文件在`/etc/pam.d/`目录下.

```
该文件适用于验证 login 的
[root@iZ94ebqp9jtZ ~]# cat /etc/pam.d/login
#%PAM-1.0
auth [user_unknown=ignore success=ok ignore=ignore default=bad] pam_securetty.so
auth       substack     system-auth
auth       include      postlogin
account    required     pam_nologin.so
account    include      system-auth
password   include      system-auth
# pam_selinux.so close should be the first session rule
session    required     pam_selinux.so close
session    required     pam_loginuid.so
session    optional     pam_console.so
# pam_selinux.so open should only be followed by sessions to be executed in the user context
session    required     pam_selinux.so open
session    required     pam_namespace.so
session    optional     pam_keyinit.so force revoke
session    include      system-auth
session    include      postlogin
-session   optional     pam_ck_connector.so
```

`pam`的配置文件有两段组成:

* `/etc/pam.conf` 主配置文件,集中式配置(redhat没有)
	* `Service type control module-path [module-arguments]` 
* `/etc/pam.d/service` 只用该文件(可以把`pam.d`下的文件合并成一个`pam.conf`)
	* `type control module-path [module-arguments]`
	
参数:
	
* `service`: 服务对应的文件名, 所有字母==必须小写==.
	* `other`: 定义默认规则(`/etc/pam.d/other`文件).
* `type`: 类型(auth, passwd, session, account)	
	* `auth`: 用来检查用户输入的账户和密码是否匹配
	* `account`: 审核用户账户是否依然有效(`passwd -l xxx`可以锁定用户, 密码对了但是不能登录进入)
	* `password`: 在用户修改密码时是否符合密码修改动作允许(复杂要求, 修改太频繁)
	* `session`: 是登录认证都成功后, 会话的相关属性(用户一次登录只可以20分钟)
* `control`: 当类型有多个条目的时候, 彼此之间如何互相作用.
	* `required`: 必须要过(有一票否决权:不过就是不通过了, 没有一票肯定权), 过了继续过同一组的其他检查, 就算`required`条目检查不同, 同组其他检查还要继续检查.
	* `requisite`: 必须的, 过了继续过同一组的其他检查, 如果不过后面不需要检查(和`reuqired`的==区别==).
	* `sufficient`: 充分条件, 如果过了就一定过了, 后面不需要检查(一票肯定权). 如果没通过(没有否决权), 其他过了也过了.
	* `optional`: 陪衬,过与不过不受影响.
	* `include`: 弃权票, 让其他文件确定.
	* `substack`:将指定配置文件中type作为参数包含在此控制语句中，和Include不同的是，在substack中完成任务或者die，只会影响substack内控制命令，不会影响完整的stack桟
* `module-path`: 完成该功能的模块是什么.
* `module-arguments`: 模块的参数.

###### control

```shell
[value1=action1 value2=action2] 如果值为1,是action1结果. 如果值为2, 是action2结果.
```

`action`:

* `ok`: 模块过了, 继续检查. 没有一票通过权.
* `done`: 一票通过权, 模块过了, 不继续检查, 并返回最终结果.
* `bad`: 结果失败了, 继续检查.
* `die`: 结果失败了, 不继续检查, 并返回最终结果. 一票否决权.
* `ignore`: 忽略
* `reset`: 忽略此前所有结果.(除了`done`和`die`返回的).


`required`: `[success=ok new_authtok_reqd=ok ignore=ignore defualt=bad]`

`requisite`: `[success=ok new_authtok_reqd=ok ignore=ignore defualt=die]` (只要不是成功, 就`die`一票否决)

`sufficient `: `[success=done new_authtok_reqd=done defualt= ignore]`

`optional `: `[success=ok new_authtok_reqd=ok defualt= ignore]`


##### 示例

至少要提供前两种服务类型.

```shell
vsftpd 只提供了三种.

[root@iZ94ebqp9jtZ ~]# cat /etc/pam.d/vsftpd
#%PAM-1.0
session    optional     pam_keyinit.so    force revoke
auth       required	pam_listfile.so item=user sense=deny file=/etc/vsftpd/ftpusers onerr=succeed
auth       required	pam_shells.so
auth       include	password-auth
account    include	password-auth
session    required     pam_loginuid.so
session    include	password-auth
```

多行, 意味着一个功能可以有多种手段.



### 整理知识点

---

#### `lib` 和 `lib64`

`lib`是放`32`位的库文件.

`lib64`是放`64`位的库文件.

#### `passwd -l xxx`

锁定一个账户.