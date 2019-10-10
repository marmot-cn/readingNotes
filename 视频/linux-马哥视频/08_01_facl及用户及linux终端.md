#08_01 facl及用户及Linux终端

###笔记

---

**FACL**

`F`ilesystem `A`ccess `C`ontrol `L`ist: 文件系统访问控制列表.

利用`文件扩展属性保存额外的访问控制权限`.

`setfacl`:

* `-m`设定:
	* u:UID:perm
	* g:GID:perm	* d:u:UID:perm	* d:g:GID:perm. (`d`:为目录设置,则里面的文件会自动继承目录的文件系统访问控制列表).
* `-x`取消:
	* u:UID
	* g:GID
	
举例:

		# getfacl inittab
		
		# file: inittab
		# owner: root
		# group: root
		user::rw-
		group::r--
		other::r--
		
`切换到另外一个用户`:

		$ su tom
		
`没有权限写入该文件`:

		echo 123 > inittab
		bash: inittab: 权限不够

`setfacl`:

		# setfacl -m u:tom:rw inittab
		# getfacl inittab
		
		# file: inittab
		# owner: root
		# group: root
		user::rw-
		user:tom:tw-
		group::r--
		mask::rw-
		other::r--

`写入成功`:

		$ echo 123 > inittab
		
**进程访问文件安全上下文**

1. 检查进程的属主跟文件属主是否一致
2. 如果不一致 检查 facl, user 中是否有进程的属主,有按照user指定的权限
3. 检查进程的用户是否属于这个文件的属组
4. 检查 facl group , 进程的用户是否属于其中的组,有按照组指定的权限
5. 应用到other上面
		
**查看登录用户**

* `w`: 显示谁登录了,并且正在干什么. 比who更详细而已
* `who`: 显示登录到当前系统的用户有哪些(su 的用户不显示)

		# who
		用户			用户从哪个终端登录         时间     注释
		chloroplast		tty1		2015-01-27 21:49 (:0)

`终端类型`:

* `console`: 控制台,连接到显示器物理硬件设备
	* `pty`:物理终端(VGA:显卡) = 控制台
	* `tty`:虚拟终端(VGA:显卡) 
	* `ttys`:串行终端(连接进主机方式不同)
	* `pts/#`:伪终端(图形化界面登录进去,通过远程登录连接进去)
	
`who –r`: 显示运行级别  
`whoami`: 显示当前系统登录的用户  
`last`: 显示/var/log/wtmp文件,`显示用户登录历史及系统重启历史`

* `-n #`: 显示最近#次的相关信息

`lastb  /var/log/btmp文件`, 显示用户错误的登录尝试  

* `-n #`: 显示最近#次的相关信息

`lastlog`: 当前系统每一个用户最近一次成功登录信息

* `-u USERNAME`: 显示特定用户最近的登录信息

**basename**

`basename`: 取得一个路径的基名,脚本中用的比较多  
`$0`: 执行脚本时的脚本路径及名称

		#!/bin/bash
		echo `basename $0`
		
		# ./first.sh
		firsht.sh

**mail**

发邮件

**hostname**

显示当前主机名

**随机数**

`RANDOM`: 0-32768

		# echo $RANDOM
		15818
		
随机数生成器: `熵池`(某些硬件中断请求中间的时间间隔保存起来,用的时候拿走,不能复制,`复制产生重复`).

* /dev/random: 如果熵池取空停下来.更安全* /dev/urandom: 如果熵池取空`用软件模拟`生成随机数

###整理知识点

---