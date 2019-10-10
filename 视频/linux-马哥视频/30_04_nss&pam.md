# 30_04_nss&pam

---

### 笔记

---

#### `system-auth`

```shell
[root@iZ94xwu3is8Z pam.d]# pwd
/etc/pam.d
[root@iZ94xwu3is8Z pam.d]# ll -s
...
0 lrwxrwxrwx. 1 root root  14 Nov 21  2014 system-auth -> system-auth-ac
...
[root@iZ94xwu3is8Z pam.d]# cat system-auth
#%PAM-1.0
# This file is auto-generated.
# User changes will be destroyed the next time authconfig is run.
auth        required      pam_env.so
auth        sufficient    pam_unix.so nullok try_first_pass
auth        requisite     pam_succeed_if.so uid >= 1000 quiet_success
auth        required      pam_deny.so

account     required      pam_unix.so
account     sufficient    pam_localuser.so
account     sufficient    pam_succeed_if.so uid < 1000 quiet
account     required      pam_permit.so

password    requisite     pam_pwquality.so try_first_pass local_users_only retry=3 authtok_type=
password    sufficient    pam_unix.so sha512 shadow nullok try_first_pass use_authtok
password    required      pam_deny.so

session     optional      pam_keyinit.so revoke
session     required      pam_limits.so
-session     optional      pam_systemd.so
session     [success=1 default=ignore] pam_succeed_if.so service in crond quiet use_uid
session     required      pam_unix.so
```

##### 模块

`module-path module-arguments` 模块路径和模块参数.

* `pam_unix.so`: password 到`shadow`中验证密码
	* `nullok` 密码允许为空
	* `try_first_pass` 在提示用户输入密码, 先尝试用户之前在其他模块输入过密码.
	* `use_fitst_pass` 上面的是尝试, 这个是直接使用在其他模块输入过密码.
	* `shadow` 读写都以`shadow`形式
	* `md5` 用户密码加密使用`md5`
	* `use_authok`: 当修改密码的时候要求上个模块输入的密码.
* `pam_permit.so`: 允许访问
* `pam_deny.so`: 拒绝访问, 通常用在`other`当中.
* `pam_cracklib.so`: 依据字典中的词检查密码, 如果密码在字典中存在则不过(同时检查规则和字典). 用于在改密码时检查弱口令.
	* `minlen`: 最短长度
	* `difok`: 验证密码和此前是否相同
	* `dcredit` 最少包含几个数字
	* `ucredit`: 要包含几位大写字母
	* `lcredit`: 要包含几个小写字母
	* `ocredit`: 要包含几个其他字符
	* `retry`: 最多尝试几次 
* `pam_shells.so`: 要求用户登录时候`shell`必须是`/etc/shells`文件中列出来的. 要求用户使用安全`shell`.
* `pam_securetty.so`: 安全的`tty`. 限定管理员只能登录到默写特殊设备.(`/etc/securetty`,`root`用户可以使用的`tty`)
* `pam_listfile.so`: 到一个文件中验证用户账户是否合法.
		
		可见vsftpd利用 pam_listfile 使/etc/vsftpd/ftpusers的用户都被拒绝.
		[root@iZ944l0t308Z ~]# cat /etc/pam.d/vsftpd
		#%PAM-1.0
		session    optional     pam_keyinit.so    force revoke
		auth       required	pam_listfile.so item=user sense=deny file=/etc/vsftpd/ftpusers onerr=succeed(没有匹配到,或者故障的默认处理机制)
		auth       required	pam_shells.so
		auth       include	password-auth
		account    include	password-auth
		session    required     pam_loginuid.so
		session    include	password-auth
* `pam_rootok.so`: 只要是`root`(`uid`=0)就通过

		[root@iZ94xwu3is8Z ~]# cat /etc/pam.d/su
		#%PAM-1.0
		
		只要id=0, 就通过(为什么管理员su到其他用户就直接通过, 如果不是则使用 system-auth 认证)
		auth		sufficient	pam_rootok.so
		# Uncomment the following line to implicitly trust users in the "wheel" group.
		#auth		sufficient	pam_wheel.so trust use_uid
		# Uncomment the following line to require a user to be in the "wheel" group.
		#auth		required	pam_wheel.so use_uid
		auth		substack	system-auth
		auth		include		postlogin
		account		sufficient	pam_succeed_if.so uid = 0 use_uid quiet
		account		include		system-auth
		password	include		system-auth
		session		include		system-auth
		session		include		postlogin
		session		optional	pam_xauth.so
* `pam_limits.so`: 在一次用户会话里面, 设定能够使用系统资源的限定. 就算是管理员也受此限制. 默认情况下使用的配置文件`/etc/security/limits.conf`或者是`/etc/security/limits.d/`目录下所有`*.conf`
* `pam_env.so`: 设置或撤销环境变量. 在用户登录之前根据`/etc/security/pam_env.conf`设置用户的环境变量.
* `pam_wheel.so`: 限定`wheel`组的用户可以`su`到`root`用户. 默认是注释掉的.
* `pam_lastlog.so`: 显示上次用户登录信息.
* `pam_issue.so`: 设定用户登录的时候, 显示`issue`文件信息
* `pam_motd.so`: 用户登录时候的欢迎信息`/etc/motd`, 用户定义是否显示该文件.
* `pam_succed_if`: 检查用户是否符合某个条件. (比如用户uid>500才可以登录系统)
		
		大于等于 1000 才可以登录, quiet_success 是不记录日志到系统日志中
		auth        requisite     pam_succeed_if.so uid >= 1000 quiet_success

* `pam_time`: 根据时间来限定登录, 让用户在特定时间登录操作系统.配置文件`/etc/security/time.conf`

		services;ttys;users;times
		
		所有用户除了root用户都拒绝通过控制台登录系统
		login ; tty* & !ttyp* ; !root ; !AL0000-2400

##### `limits.conf`

限定一个用户使用本地资源, 编辑该文件.

* `domain`
	* 用户名
	* 组名(`@group`)
	* `*`所有的, 默认.
	* `%`(?)
* `type`
	* `soft` 软限制
	* `hard` 硬限制
* `item`
	* `core`: 能打开的核心文件(内核中运作的)大小.
	* `nofile`: 能打开最大的文件数.
	* `rss`: 最大实际内存(最多使用多少实际内存), 单位是`KB`.
	* `as`: 地址空间限制
	* `cpu`: 最多使用`CPU`时间,单位是`MIN(分钟)`(设定一个用户`cpu`为50分钟, 如果该用户登录系统进行编译, 超过50分钟就失败)
	* `nproc`: 最大进程数, 最多可以打开多少个进程.
	* `nice`: 最大`nice`值可以提升.


```shell
[root@iZ94xwu3is8Z ~]# cat /etc/security/limits.conf
# /etc/security/limits.conf
#
#This file sets the resource limits for the users logged in via PAM.
#It does not affect resource limits of the system services.
#
#Also note that configuration files in /etc/security/limits.d directory,
#which are read in alphabetical order, override the settings in this
#file in case the domain is the same or more specific.
#That means for example that setting a limit for wildcard domain here
#can be overriden with a wildcard setting in a config file in the
#subdirectory, but a user specific setting here can be overriden only
#with a user specific setting in the subdirectory.
#
#Each line describes a limit for a user in the form:
#
#<domain>        <type>  <item>  <value>
#
#Where:
#<domain> can be:
#        - a user name
#        - a group name, with @group syntax
#        - the wildcard *, for default entry
#        - the wildcard %, can be also used with %group syntax,
#                 for maxlogin limit
#
#<type> can have the two values:
#        - "soft" for enforcing the soft limits
#        - "hard" for enforcing hard limits
#
#<item> can be one of the following:
#        - core - limits the core file size (KB)
#        - data - max data size (KB)
#        - fsize - maximum filesize (KB)
#        - memlock - max locked-in-memory address space (KB)
#        - nofile - max number of open files
#        - rss - max resident set size (KB)
#        - stack - max stack size (KB)
#        - cpu - max CPU time (MIN)
#        - nproc - max number of processes
#        - as - address space limit (KB)
#        - maxlogins - max number of logins for this user
#        - maxsyslogins - max number of logins on the system
#        - priority - the priority to run user process with
#        - locks - max number of file locks the user can hold
#        - sigpending - max number of pending signals
#        - msgqueue - max memory used by POSIX message queues (bytes)
#        - nice - max nice priority allowed to raise to values: [-20, 19]
#        - rtprio - max realtime priority
#
#<domain>      <type>  <item>         <value>
#

#*               soft    core            0
#*               hard    rss             10000

学员组 最多进程个数 20
#@student        hard    nproc           20
#@faculty        soft    nproc           20
#@faculty        hard    nproc           50
#ftp             hard    nproc           0

学员组 软硬都限制 最多允许登录4次, 登录4次以后就不能登录了
#@student        -       maxlogins       4

# End of file
* soft nofile 65535
* hard nofile 65535
```

#### 示例

```shell
system-auth 文件中 auth 开头的行都包含进来
auth include system-auth
```

```shell
这里的 use_authok 是当用户的密码验证通过以后(第一步), 在设置密码时必须是上次输入的密码(第一步验证通过的密码).
password required pam_cracklib.so retry=3
password required pam_unix.so use_authok
```

创建一个文件, 只能让属于一个组的用户可以登录.

```shell
修改 system-auth, 添加一行

检测用户环境变量
auth        required      pam_env.so

新添加的一行, 必须要过
auth		required		pam_listfile.so item=group sense=allow file=/etc/pam_allowgroups

auth        sufficient    pam_unix.so nullok try_first_pass
auth        requisite     pam_succeed_if.so uid >= 1000 quiet_success
auth        required      pam_deny.so
```

### 整理知识点

---

#### 软限制 硬限制

普通用户只能使用软限制`ulimit`.

软限制可以超出. 可以使用`ulimit`调整限制.

`ulimit`
 
* `-u`: 调整最大进程数
* `-n`: 调整最大文件打开数

#### `export`

使用`export`命令可以看见当前用户的环境变量.

#### `pam`用户登录流程

```shell
auth        required      pam_env.so
auth        sufficient    pam_unix.so nullok try_first_pass
auth        requisite     pam_succeed_if.so uid >= 500 quiet
auth        required      pam_deny.so
```

第一部分表示,当用户登录的时候,首先会通过auth类接口对用户身份进行识别和密码认证.所以在该过程中验证会经过几个带auth的配置项.

`auth        required      pam_env.so`
其中的第一步是通过pam_env.so模块来定义用户登录之后的环境变量, pam_env.so允许设置和更改用户登录时候的环境变量,默认情况下,若没有特别指定配置文件,将依据/etc/security/pam_env.conf进行用户登录之后环境变量的设置.

`auth        sufficient    pam_unix.so nullok try_first_pass`
然后通过pam_unix.so模块来提示用户输入密码,并将用户密码与/etc/shadow中记录的密码信息进行对比,如果密码比对结果正确则允许用户登录,而且该配置项的使用的是“sufficient”控制位,即表示只要该配置项的验证通过,用户即可完全通过认证而不用再去走下面的认证项.不过在特殊情况下,用户允许使用空密码登录系统,例如当将某个用户在/etc/shadow中的密码字段删除之后,该用户可以只输入用户名直接登录系统.

`auth        requisite     pam_succeed_if.so uid >= 500 quiet`
下面的配置项中,通过pam_succeed_if.so对用户的登录条件做一些限制,表示允许uid大于500的用户在通过密码验证的情况下登录,在Linux系统中,一般系统用户的uid都在500之内,所以该项即表示允许使用useradd命令以及默认选项建立的普通用户直接由本地控制台登录系统.

`auth        required      pam_deny.so`
最后通过pam_deny.so模块对所有不满足上述任意条件的登录请求直接拒绝,pam_deny.so是一个特殊的模块,该模块返回值永远为否,类似于大多数安全机制的配置准则,在所有认证规则走完之后,对不匹配任何规则的请求直接拒绝.

```shell
account     required      pam_unix.so
account     sufficient    pam_succeed_if.so uid < 500 quiet
account     required      pam_permit.so
```

第二部分的三个配置项主要表示通过account账户类接口来识别账户的合法性以及登录权限.

`account     required      pam_unix.so`
第一行仍然使用pam_unix.so模块来声明用户需要通过密码认证.

`account     sufficient    pam_succeed_if.so uid < 500 quiet`
第二行承认了系统中uid小于500的系统用户的合法性.

`account     required      pam_permit.so`
之后对所有类型的用户登录请求都开放控制台.    


```shell
password    requisite     pam_cracklib.so try_first_pass retry=3
password    sufficient    pam_unix.so md5 shadow nullok try_first_pass use_authtok
password    required      pam_deny.so
```

第三部分会通过password口另类接口来确认用户使用的密码或者口令的合法性.

`password    requisite     pam_cracklib.so try_first_pass retry=3`第一行配置项表示需要的情况下将调用`pam_cracklib`来验证用户密码复杂度.如果用户输入密码不满足复杂度要求或者密码错,最多将在三次这种错误之后直接返回密码错误的提示,否则期间任何一次正确的密码验证都允许登录.需要指出的是,pam_cracklib.so是一个常用的控制密码复杂度的pam模块,关于其用法举例我们会在之后详细介绍.

`password    sufficient    pam_unix.so md5 shadow nullok try_first_pass use_authtok`
`password    required      pam_deny.so`
之后带`pam_unix.so`和`pam_deny.so`的两行配置项的意思与之前类似.都表示需要通过密码认证并对不符合上述任何配置项要求的登录请求直接予以拒绝.不过用户如果执行的操作是单纯的登录,则这部分配置是不起作用的.


```shell
session     optional      pam_keyinit.so revoke
session     required      pam_limits.so
session     [success=1 default=ignore] pam_succeed_if.so service in crond quiet use_uid
session     required      pam_unix.so
```
第四部分主要将通过session会话类接口为用户初始化会话连接.其中几个比较重要的地方包括,使用pam_keyinit.so表示当用户登录的时候为其建立相应的密钥环,并在用户登出的时候予以撤销.不过该行配置的控制位使用的是optional,表示这并非必要条件.之后通过pam_limits.so限制用户登录时的会话连接资源,相关pam_limit.so配置文件是/etc/security/limits.conf,默认情况下对每个登录用户都没有限制.关于该模块的配置方法在后面也会详细介绍.