#Docker run -it 作用

---

**`-it`**

* `-i`参数的作用是保证容器中`STDIN`是开启的.
* `-t`告诉Docker为要创建的容器分配一个伪tty终端.

`-it`可以提供一个交互式的shell

####示例

**不带`-it`**

我们先不带`-it`运行容器.

		[root@dev-server-2 chloroplast]# docker run --name=my_debian docker.io/debian
		
什么反应都没有就退出了.因为我们没有启动任何进程.因为想要保持Docker容器的活跃状态,需要其中运行的进程不能中断.

**`-i`**

		[root@dev-server-2 chloroplast]# docker run -i --rm --name=my_debian docker.io/debian
		
		#运行ls
		ls
		#得到如下
		bin
		boot
		dev
		etc
		home
		lib
		lib64
		media
		mnt
		opt
		proc
		root
		run
		sbin
		srv
		sys
		tmp
		usr
		var
		
输入`exit`退出进程,则容器退出

**`-t`**

		[root@dev-server-2 chloroplast]# docker run -t --rm --name=my_debian docker.io/debian
		root@25e240762bfc:/# ls
		
		#无任何反应,因为STDIN没有开启
	


