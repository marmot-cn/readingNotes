#修改SSH默认端口

---

打开`/etc/ssh/sshd_config`

		sudo vi /etc/ssh/sshd_config
		
找见:

		...
		#Port 22
		#ListenAddress 0.0.0.0
		#ListenAddress ::
		...
		
		
修改如下:

		Port 22
		Port 17456
		#ListenAddress 0.0.0.0
		#ListenAddress ::
		
这里放开2个端口主要是用于测试,放置新的端口17456不行而不能ssh到服务器上.
		
找其他机器尝试用 `17456`端口`ssh`,如果可以则屏蔽掉`22`端口好.最终配置文件如下:

		#Port 22
		Port 17456
		#ListenAddress 0.0.0.0
		#ListenAddress ::
		
修改 `PermitRootLogin将yes改为no`: 用户禁止root远程登录
		
让`sshd`重新加载配置文件

		sudo /etc/init.d/sshd reload
		
或重启`sshd`服务

		sudo service sshd restart