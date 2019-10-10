# 02_01 Linux操作系统及常用命令

### 笔记

---

**linxu桌面**

x-window. `X` 表示图形显示协议

**IDE**

集成开发环境

**DLL(Windows)**

Windows Dynamic Link Library 动态链接库

**so(Linux)**

`.so`(shared object)共享对象. 也称为 `dso`(dynamic shared library)

**权限**

设定`访问资源能力`.

**用户ID**

用于计算机识别用户,用户ID(用户标识符)为数字.

`每个用户名都有用户id`

**认证机制**

`Authentication`. 认证的过程`鉴别一个用户的过程`.(`密码`就是一种认证机制)

**授权**

`Authorization`

**审计**

`Audtion`,通过日志来完成.

**命令提示符**

`prompt`

**#**

`#`是专属root账号的提示符,非root账号都是`$`.

**magic number**

魔数.这是放在linux的目录中的文件信息块中的一个`标识符`,一般`只有几位`,用来标识该文件是什么类型的文件,可以被什么样的应用使用.这个魔数`不是固定`的,有时候一个文件信息中的魔数可能会不断变化.

**ls**

`list` 列表,列出. 列出指定路径下所有目录和文件.

		ls -l (长格式,显示完整信息)
		ls -h (human readable,人类可读的单位换算)
		ls -a 显示所有文件(显示以.开头的隐藏文件). Linux 隐藏文件以'.'开头
		ls .  表示当前目录
		ls .. 表示父目录
		ls -A 不显示 '.' 和 '..'
		ls -d 显示目录自身属性
		ls -r 逆向显示
		ls -R 递归显示

**目录**

目录也是一种`文件`.`路径映射文件`.

`路径`:从指定`起始点到目的地所经过的位置`.

**文件系统(file system)**

层次化文件管理机制

**linux目录结构**

`倒置`的`树结构`

* 可以分叉的节点:`目录`
* 不能分叉的节点:`文件`.
* 从根开始往下找:`绝对路径`

**工作目录和当前目录**

登录系统以后,无时无刻都处于一个目录下面称为:

* 工作目录: working directory
* 当前目录: current directory

**相对路径**

相对于当前所处位置的目录.

`相对`:`根到目的地中间的一个叉`才能称为相对.

**FHS**

Filesystem Hierarchy Standard,文件系统目录标准. Linux上一些标准化目录(bin,sbin)等.

**PWD**

printing working directory. 显示当前目录.

**/**

Linux目录,最顶级用 `/` 表示. 各目录分隔也用 `/`.

**文件信息**

		-rw------- 1 root root 1371 May 17 2012 anaconda-ks.cfg
		
* `-rw-------`(10位):
	
	* 第1位: 文件类型  
		`-` : 普通文件(f)  
		`d` : 目录文件  
		`b` : 块设备文件(block)  
		`c` : 字符设备文件(character)  
		`l` : 符号链接文件(symbolic link file).软链接文件  
		`p` : 命令管道(pipe)  
		`s` : 套接字文件(socket)  
		
	* 后面9位: 文件权限,每`3`位一组,每一组:`rwx`(读,写,执行).无对应权限用`-`表示.  
			
			rw- 能读,能写,不能执行
			
* `1`: 数字表示,文件`硬连接`的次数.
* `root` : 文件属主(owner)
* `root` : 文件属组(group)
* `1371` : 文件大小(size).单位是`字节`.
* `May 17 2012` : 时间戳(timestamp).每一个文件有3个时间戳.

	* `access`: 最近一次被`访问`的时间. 
	* `modify`: 最近一次被`修改`的时间.(上例默认显示的)
	* `change`: 最近一次被`改变`的时间.
* `anaconda-ks.cfg`: 文件名

**index node**

计算机靠数字索引文件.每一个文件都有一个唯一的数字标示符.`inode`(index node).

		ls -i 显示文字索引节点号
		
**cd**

`change directory` 切换目录.

家目录,主目录,`home directory`

		cd ~USERNAME(用户的用户名),进入指定用户的家目录
		cd - 在当前目录和前一次所在目录之间来回切换
		
**命令类型**

* `内部`命令(shell 内置).内部,内建

		$ type cd
		cd is a shell builtin 
* `外部`命令.在文件系统的某个路径下有一个`与命令名称相应`的`可执行文件`.
* `type`: 显示命令属于哪种类型

**环境变量**

`变量`是命名的内存空间

* 变量赋值: 

		NAME=Jerry
		在内存中找一段空间起名叫 NAME,空间数据叫Jerry
		
* 变量声明: 申请内存使用的过程
* PATH: 一堆使用':'(冒号)隔开的路径

当我们`使用命令`时,`依次`在路径里面寻找`第一个寻找到的`.第一次执行时去找命令,然后`缓存`.第2次不用去找.

		$ hash
		hits 	command
		1		/usr/bin/printenv
		1		/bin/ls

`hash`是一种缓存. 在缓存中记录下来,`此前使用的所有命令的路径`.

**Hash**

Hash(哈希)是键值对,在`hash`中查找速度是`o(1)`.

### 整理知识点

---
