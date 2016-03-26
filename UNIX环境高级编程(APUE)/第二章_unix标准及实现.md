#第二章 UNIX标准及实现

###2.2 UNIX标准化

####ISO C

`提供C程序的可移植性`, 使其能适合于大量不同的操作系统, 而不只是适合UNIX系统. 此标准不仅定义了C程序设计语言的语法和语义,还定义了其标准库.

####IEEE POSIX

`POSIX`指的是可移植操作系统接口(`Portable Operating System Interface`).

该标准的目的是`提升应用程序在各种UNIX系统环境之间的可移植性`.

该标准说明了一个`接口(interface)`而不是一种`实现(implementation)`.

####Single UNIX Specification

`Single UNIX Specification(SUS, 单一UNIX规范)`是`POSIX.1`标准的一个`超集`,定义了一些附加接口扩展了POSIX.1规范提供的功能.

`POSIX.1`相当于`Single UNIX Specification`中的基本规范部分.

**XSI**

`XSI(XIS conforming)`,只有`遵循` `XSI`的实现才能称为`UNIX`系统.

遵循XSI的实现必须支持POSIX.1的部分:

* 文件同步
* 线程栈地址
* 长度属性
* 线程进程共享同步
* `_XOPEN_UNIX_`符号常量

###2.3 UNIX系统实现

####SVR4

####4.4BSD

####FreeBSD

####Linux

####Mac OS X

####Solaris

###2.4 标准和实现的关系

####ISO C 限制

`编译`时的限制

####POSIX 限制

操作系统实现限制的常量

* 数值限制
* 最小值
* 最大值
* 运行时可以增加的值
* 运行时不变值
* 其他不变量
* 路径名可变值

####XSI 限制

####函数`sysconf` `pathconf` `fpathconf`

**sysconf**

`long int sysconf(int name);`

用于获得与文件或目录无关的限制值以及系统特征选项.

以`_SC_`开始的常量用作标识`运行时限制`的`sysconf`参数.

**pathconf**

`long int pathconf(const char *pathname,int name);`

获得与文件或目录有关的限制值,作用于文件名pathname.

以`_PC_`开始的常量用作标识`运行时限制`的`pathconf`参数.

**fpathconf**

`long int fpathconf(int filedes,int name);`

获得与文件或目录有关的限制值,作用于已打开的文件描述符 filedes.

以`_PC_`开始的常量用作标识`运行时限制`的`fpathconf`参数.

###2.5 限制

###2.6 选项

**处理选项的方法**

* 编译时选项定义在`<unistd.h>`中.
* 与文件或目录无关的运行时选项用`sysconf`函数来判断.
* 与文件或目录有关的运行时选项通过调用`pathconf`或`fpathconf`函数来判断.

###2.7 功能测试宏

用来在编译时控制一些头文件的版本.

`_POSIX_C_SOURCE`

`_XOPEN_SOURCE`

###2.8 基本系统数据类型

**基本系统数据类型**

定义在头文件`<sys/types.h>`中.

* `clock_t`: 时钟滴答计数器
* `comp_t`: 压缩的时钟滴答
* `dev_t`: 设备号(主和次)
* `fd_set`: 文件描述符集
* `fpos_t`: 文件位置
* `gid_t`: 数值组ID
* `ino_t`: i节点编号
* `mode_t`: 文件类型,文件创建模式
* `nlink_t`: 目录项的链接计数
* `off_t`: 文件长度和偏移量
* `pid_t`: 进程ID和进程组ID
* `pthread_t`: 线程ID
* `ptrdiff_t`: 两个指针相减的结果(带符号的)
* `rlim_t`: 资源限制
* `sig_atomic_t`: 能原子性地访问的数据类型
* `sigset_t`: 信号集
* `size_t`: 对象长度(不带符号)
* `ssize_t`: 返回字节计数的函数(带符号的)
* `time_t`: 日历时间的秒计数器
* `uid_t`: 数值用户ID
* `wchar_t`: 能表示所有不同的字符码

###2.9 标准之间的冲突

如果出现冲突,`POSIX.1` 服从 `ISO C`标准.

###2.10 小结