#12_04_Linux软件管理之四 yum

###笔记

---

**RPM安装**

* 二进制安装(源程序->编译->二进制),有些特性是在编译时选定的.如果编译时未选定此特性,将无法使用.
* rpm包落后于源码包.

**定制**

手动编译安装,定制安装.

`编译环境`: 开发环境 (开发库,开发工具:编译器)

`gcc`: GUN C Complier, C  
`g++`: c++ 编译器

`make`: C 或 C++ 的项目管理工具, 将多个不同文件做成一个项目,编译过程通过配置文件`makefile`编译.

`makefile`: 定义了make(gcc,g++)按何种次序去编译这些源程序文件中的源程序.

`生成makfile`:

* `automake` --> 生成 makefile.in (半成品)--> makefile
* `autoconf` --> configure(脚本,配置程序如何编译,用户选择可选特性,告诉安装在什么路径下)

`make install` 安装

**手动编译安装软件包三步骤**

在源程序目录

* `./configure` 指定编译属性
	* `--prefix=/path/to/somewhere`: 指定安装路径 
	* `--sysconfigdir=/PATH/TO/CONFFILE_PATH`: 配置路径
	* `--help`: 获取脚本的使用格式
* `make`: 编译
* `make install`: 安装

**`./configure`的作用和功能**	

1. 让用户选定编译特性
2. 检查编译环境


**如果安装在非默认路径下的一些额外步骤**

1. 修改`PATH`环境变量,以能够识别此程序的二进制文件路径.
	* 修改`/etc/profile`文件
	* 在`/etc/profile.d/`的目录建立一个以`.sh`为名称后缀的文件,在里面定义`export PATH=$PATH:/path/to/somewhere`
2. 默认情况下,系统搜索库文件的路径/lib,/usr/lib; 要增添额外搜寻径路:
	* 在`/etc/ld.so.conf.d/`中创建一个以`.conf`为后缀的文件,而后把要增添的路径直接写支此文件中(下次重启有效)
	* `ldconfig` 通知系统重新搜寻库文件
		* `-v`: 显示重新搜寻库的过程
3. 头文件: (#include包含头文件)输出给系统,默认在`/usr/include`.增添头文件搜寻路径,使用`链接`进行:

		/usr/local/tengine/include/ 	/usr/include/
		ln -s /usr/local/tengine/include/* /usr/include/
		或者
		ln -s /usr/local/tengine/include/ /usr/include/tengine
	
4. man文件(帮助文件)路径: 安装在 `--prefix`指定的目录下的`man`目录: `/usr/share/man`

		man -M /PATH/TO/MAN_DIR COMMAND (指定路径访问)
		或
		在/etc/man.config中添加一条MANPATH
		

**netstat命令:网络状态**

* `-r`: 显示路由表
* `-n`: 以数字方式显示
* `-t`: 建立的tcp连接
* `-u`: 显示udp连接
* `-l`: 显示监听状态的连接
* `-p`: 显示监听指定套接字的进程的进程号与进程名


###整理知识点

---

