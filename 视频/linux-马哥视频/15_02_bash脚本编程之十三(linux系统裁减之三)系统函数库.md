#15_02_bash脚本编程之十三(Linux系统裁减之三)系统函数库

###笔记

---

**终端提示信息**

`/etc/issue`文件的内容

**rc.sysinit:挂载`/etc/fstab`中定义的其他文件系统**

如果一个文件系统已经挂载结束,会出现在`/proc/mounts`文件当中.

**示例**

去除/etc/fstab中所有属于swap或者proc或者sysfs的文件系统.

		grep -E -v "\<swap|proc|sysfs\>" /etc/fstab | awk '{print $1}'

判断设备有没有挂载

		grep -E -v "\<swap|proc|sysfs\>" /etc/fstab | awk '{print $1}' | while read LINE; do awk '{print $1}' /proc/mounts | grep "~$LINE$"; done
	
**设定内核参数**

`/etc/sysctl.conf`

sysctl -p 让上述文件命令生效

**用户**

PAM: Pluggable Authentication Module

/etc/pam.d/*

login: 验证

nsswitch : Network Service Switch (网络服务转换,中间层:用户登录找账户信息和密码,host,dns...,名称解析是如何工作的)

* 框架,能够完成配置去哪找用户的账户密码(/etc/passwd,/etc/shadow,/etc/group)
	* 库: libnss_file.so(去文件/etc/passwd找), libnss_nis.so, libnss_ldap.so
	* 配置文件: /etc/nsswitch.conf
	
###整理知识点

---

####nsswitch

####mingetty和agetty

####pam