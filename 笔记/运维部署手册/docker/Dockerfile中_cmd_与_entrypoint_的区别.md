#Dockerfile中 CMD 与 ENTRYPOINT 的区别

---

####cmd

**命令格式**

* CMD `["executable","param1","param2"]` (exec form, `this is the preferred form`).运行一个可执行的文件并提供参数.
* CMD `["param1","param2"]` (as default parameters to `ENTRYPOINT`).为ENTRYPOINT指定参数.
* CMD `command param1 param2` (shell form).是以"`/bin/sh -c`"的方法执行的命令.

**示例**

如果指定:

		CMD [“/bin/echo”, “this is a echo test ”] 
		
build后运行(假设镜像名为ec):

		docker run ec
		
输出:

		this is a echo test
		
可以理解为`开机启动项`

**docker run命令中的参数和Dockerfile中的CMD参数**

docker run命令如果指定了参数会把CMD里的参数`覆盖`.

同样的ec镜像启动:

		docker run ec /bin/bash
		
就`不会输出`:this is a echo test 因为`CMD命令被"/bin/bash"覆盖了`.

####ENTRYPOINT

An ENTRYPOINT allows you to configure a container that will run as an executable.

它可以让你的**`容器功能表现得像一个可执行程序一样`**.

**命令格式**

* ENTRYPOINT `["executable", "param1", "param2"]` (the preferred exec form) 
* ENTRYPOINT command param1 param2 (shell form) 

也可以在docker run命令时使用`–entrypoint`指定.

**示例一**

使用下面的ENTRYPOINT构造镜像:

		ENTRYPOINT ["/bin/echo"] 
		    
那么docker build出来的镜像以后的容器功能`就像一个/bin/echo程序`.比如build出来的镜像叫imageecho, 那么可以这样使用:

		docker run -it imageecho "this is a test"
		
这里就会输出”this is a test”这串字符,而`这个imageecho镜像对应的容器表现出来的功能就像一个echo程序一样`.添加的参数"this is a test"会添加到ENTRYPOINT后面,就成了这样　/bin/echo "this is a test".

**示例二**

使用下面的ENTRYPOINT构造镜像:
		
		ENTRYPOINT ["/bin/cat"]
		
假设镜像名为st

		docker run -it st /etc/fstab 
		
这样相当: `/bin/cat /etc/fstab` 这个命令的作用	.运行之后就输出/etc/fstab里的内容.

**示例三**

ENTRYPOINT设为[“/bin/sh -c”]时候运行的情况:

		# docker run -it t2 /bin/bash 
		root@4c8549e7ce3e:/# ps 
		PID TTY TIME CMD 
		1 ? 00:00:00 　sh 
		9 ? 00:00:00 　bash 
		19 ? 00:00:00 　ps 

可以看到:
	
* `PID`为`1`的进程运行的是`sh`.
* `bash`只是`sh`的一个`子进程`.
* `/bin/bash`只是作为`/bin/sh -c`后面的`参数`.

CMD可以为ENTRYPOINT提供参数,ENTRYPOINT本身也可以包含参数.`可以把那些可能需要变动的参数写到CMD里而把那些不需要变动的参数写到ENTRYPOINT里面`.

		FROM ubuntu:14.10  
		ENTRYPOINT ["top", "-b"]   
		CMD ["-c"] 
		
把可能需要变动的参数写到CMD里面.然后你可以在docker run里指定参数,这样CMD里的参数(这里是-c)就会被覆盖掉而ENTRYPOINT里的不被覆盖.
   