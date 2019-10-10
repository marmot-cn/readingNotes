# Linux别名登录服务器 SSH别名设置

#### 修改设置

修改`/etc/ssh/ssh_config`或`~/.ssh/config`文件,添加如下几行

		Host 别名
			HostName 服务器地址(www.xxx.com或xxx.xxx.xx)
			Port 端口号(默认22)
			User 用户名
			
**示例**

		cat ~/.ssh/config
		Host my-server
			hostname 120.25.161.1
			user	chloroplast
			
	
		➜  ~  ssh my-server
		Last login: Tue Dec 15 10:07:56 2015 from 1.86.30.54
		
		Welcome to aliyun Elastic Compute Service!
		
		[chloroplast@iZ94ebqp9jtZ ~]$