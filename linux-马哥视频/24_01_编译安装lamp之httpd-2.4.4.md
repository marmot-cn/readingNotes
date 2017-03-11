#24_01_编译安装LAMP之httpd-2.4.4

---

###笔记

---

####Zend

* 第一段: 词法分析,语法分析,编译为`opcode`.
* 第二段: 执行`opcode`.


`opcode`是动态编译的,下次运行直接执行已经编译好的`opcode`.

不同的php进程之间是无法共享opcode,opcode是放在内存当中.

####PHP缓存器

为了实现多个PHP进程共享opcode,引入缓存加速器:

* apc
* eAccelerator
* XCache

####PHP解释器和mysql交互

PHP程序和mysql交互,解释器自身不和mysql交互.

php解释器用于解释php脚本.

####apr

Apache Portable Rntime

httpd + apr 适配不同的环境

httpd程序相同,把底层的不同通过虚拟机(apr)来兼容起来.

####编译配置LAMP

**MPM**

* prefork
* worker
* event

###整理知识点

---