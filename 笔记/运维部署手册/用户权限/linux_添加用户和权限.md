# Linux 添加用户和权限

---

#### 添加自用用户

添加用户名为`chloroplast`的用户:

		adduser chloroplast
		
修改用户密码:		

		passwd chloroplast
		
两次输入密码

#### 让chloroplast赋予root权限

修改`/etc/sudoers`权限

		chmod 755 /etc/sudoers

修改`/etc/sudoers`:

		## Allow root to run any commands anywhere
		root	ALL=(ALL) 	ALL
		chloroplast	ALL=(ALL)	ALL
		
在`root`下添加`chloroplast`
	
还原`/etc/sudoers`权限

		chmod 440 /etc/sudoers
		
		
		