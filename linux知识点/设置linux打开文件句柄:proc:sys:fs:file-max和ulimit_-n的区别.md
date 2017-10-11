### 设置Linux打开文件句柄/proc/sys/fs/file-max和ulimit -n的区别 

---

**max-file**

表示系统级别的能够打开的文件句柄的数量.是**对整个系统的限制,并不是针对用户**的.

**ulimit -n**

控制进程级别能够打开的文件句柄的数量.提供对shell及其启动的进程的可用文件句柄的控制.这是进程级别的.



对于服务器来说,`file-max`和`ulimit`都需要设置,否则会出现文件描述符耗尽的问题.
一般如果遇到文件句柄达到上限时,会碰到"`Too many open files`"或者`Socket/File: Can’t open so many files`等错误.
为了让服务器重启之后,配置仍然有效,需要用永久生效的配置方法进行修改.

**max-file:**

查看系统级别的能够打开的文件句柄的数量，Centos7默认是794168

		# cat /proc/sys/fs/file-max
		794168
		
系统级打开最大文件句柄的数量永久生效的修改方法,修改文件,文件末尾加入配置内容:

		# vim /etc/sysctl.conf
		fs.file-max = 2000000
		
然后执行命令,使修改配置立即生效:
		
		# sysctl -p
		
		从指定的文件加载系统参数,如不指定即从/etc/sysctl.conf中加载

**ulimit:**

查看用户进程级的能够打开文件句柄的数量,Centos7默认是1024

这里设置的是当前shell的`当前用户`的打开的最大限制,如果当前用户打开多个shell,则每个shell都能打开该最大值.

		# ulimit -n
		1024
		
进程级打开文件句柄数量永久生效的修改方法,修改文件,文件末尾加入配置内容:

这里限制一个用户的所有shell能打开的最大数:

		# vim /etc/security/limits.conf
		* soft nofile 65535
		* hard nofile 65535
		
修改以后,需要重新登录才能生效.

如果需要设置当前用户session立即生效,还需要执行:

		# ulimit -n 65535 
		

		