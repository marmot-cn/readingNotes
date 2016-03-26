#04_01 Linux用户管理命令详解

###笔记

---

**useradd**

`useradd` `[options]` `USERNAME`
		
`-u`: UID 其他用户尚未使用的id号,且>500  
`-g`: GID(可以写组id,组名) 基本组,这个组要事先存在

`用户默认id`是从`passwd`文件内最大`uid+1`

`-G`: GID,...(附加组)  
`-c "COMMENT"`: 注释信息  
`-d`: /path/to/directory 指定家目录
`-s SHELL路径`: 指定用户shell  

* 创建一个`不能登录`的用户

		# useradd -s /sbin/nologin user5
		# su - user5
		This account is currently not available.
		
* `SHELL(环境变量)`: 用于保存当前用的默认shell

`-m -k`: 创建家目录,复制 `/etc/skel 配置文件`到`家目录`  
`-M`: 不创建家目录
		
		# useradd -M user7
		# su - user7
		su: warning: cannot change directory to /home/user7: No such file
`-r`: 添加系统用户(uid 1-499)
		
**/etc/login.defs**

设置用户帐号限制的文件,在这里我们可配置密码的`最大过期天数`,`密码的最大长度约束`等内容.该文件里的配置`对root用户无效`.如果/etc/shadow文件里有相同的选项,则以/etc/shadow里的设置为准,也就是说`/etc/shadow的配置优先级高`于/etc/login.defs

**userdel**

`userdel` `[option]` `USERNAME`

删除用户不指定任何选项,删除用户家目录`不会`被删除.

`-r`: 同时删除用的家目录

**修改用户账号属性**

`usermod`:

* `-u`: UID(uid)
* `-g`: GID(基本组)
* `-G`: GID(修改了附加组,以前的附加组没了)
* `-a -G`: GID(为用户在原有的附加组上添加新的附加组)
* `-c`: 修改注释
* `-d`: 为用户指定新的家目录
* `-d -m`: 指定新的家目录,移动此前文件到新的家目录中去
* `-s`: 改用户的shell
* `-l`: 改用户的登录名
* `-L`: 锁定账号
* `-U`: 解锁账号

**查看用户信息**

`id`: 查看用户信息
`finger`: 检索用户信息

		finger USERNAME

**修改用户的默认**

`chsh`: 修改用户的`默认shell`

**修改用户注释信息**

`chfn`

**密码管理**

`passwd` `[USERNAME]` 管理员可以修改其他用户的密码,非管理员只能修改自己的

`--stdin`: 从标准输入读取密码

		# echo "redhat" | passwd --stdin user3
		Changing password for user user3
		
`-l`: 锁定用户账户  
`-u`: 解锁用户账户  
`-d`: 删除用户密码  

**pwck**

检查用户账号完整性

**组管理**

`groupadd`: 创建组

* `-g`:GID
* `-r`:添加系统组(id 1-499)

`groupmod`: 修改组属性

* `-g`:GID
* `-n`:GROUPNAME

`groupdel`: 删除组

`gpasswd`: 为组设定密码 

`newgrp` `GROUPNAME`: 登陆到一个新的基本组,这样创建进程的属组就变了.如果想切换回去,使用 `exit`.

**修改用户密码过期信息**

`change`:

* `-d`: 最近一次的修改时间
* `-E`: 过期时间* `-I`: 非活动时间* `-m`: 最短使用期限* `-M`: 最长使用期限* `-W`: 警告时间

###整理知识点

---