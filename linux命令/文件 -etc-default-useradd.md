# 文件 /etc/default/useradd

---

通过`useradd`添加用户时的规则文件.

```shell
[ansible@iZ94ebqp9jtZ ~]$ cat /etc/default/useradd
# useradd defaults file
GROUP=100
HOME=/home
INACTIVE=-1
EXPIRE=
SHELL=/bin/bash
SKEL=/etc/skel
CREATE_MAIL_SPOOL=yes
```

* `GROUP=100`, 默认用户组100. 我们创建用户的时候默认会创建一个和用户名相同的用户组.但是当我们使用参数`useradd -N(--no-user-group 不创建和用户名相同的用户组)`时, 创建用户的用户组为该配置文件的用户组.

	```shell
	[ansible@iZ94ebqp9jtZ ~]$ sudo useradd jerry -N
	[ansible@iZ94ebqp9jtZ ~]$ id jerry
	uid=1002(jerry) gid=100(users) groups=100(users)
	```

* `HOME=/home`: 把用户的家目录建在/home中.
* `INACTIVE=-1`：是否启用帐号过期停权, `-1`表示不启用.
* `EXPIRE=`：帐号终止日期, 不设置表示不启用.
* `SHELL=/bin/bash`：所用SHELL的类型.
* `SKEL=/etc/skel`: 默认添加用户的目录默认文件存放位置;也就是说,当我们用`adduser`添加用户时,用户家目录下的文件,都是从这个目录中复制过去的.


## adduser

`adduser`是`useradd`的软连.

```shell
可见其前面属性l
[ansible@iZ94ebqp9jtZ ~]$ sudo ls -l /sbin/adduser
lrwxrwxrwx. 1 root root 7 Jul 10 20:22 /sbin/adduser -> useradd
[ansible@iZ94ebqp9jtZ ~]$ sudo file /sbin/adduser
/sbin/adduser: symbolic link to `useradd'
```