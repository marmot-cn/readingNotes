#21_03_httpd属性配置

---

###笔记

---

####URL

URL 是相对于 DocumentRoot 的路径而言的.

####Options

* `None`: 不支持任何选项.
* `Indexes`: 允许索引目录.
* `FollowSynLinks`: 允许访问符号链接指向文件.
* `Includes`: 允许服务器端包含(SSI).
* `SymLinksifOwnerMatch`: 允许访问符号链接,但是必须执行 `httpd`进程的的用户和文件所属用户一致.
* `ExecCGI`: 允许运行`cgi`脚本.
* `MultiViews`: 内容协商机制,多功能视图. 可以根据客户端来源的语言和文字来判断显示哪一个网页.(有多个默认页,不同地域的用户访问不同的页面).
* `All`: 启用所有选项.

####AllowOveride

允许覆盖`服务器访问控制列表`.

在 AllowOverride 设置为 None 时,`.htaccess` 文件将被完全忽略。当此指令设置为 All 时,所有具有 "`.htaccess`" 作用域的指令都允许出现在 `.htaccess` 文件中.

还可以对它指定如下一些能被重写的指令类型:

* `AuthConfig`: 允许使用所有的权限指令,他们包括 AuthDBMGroupFile  AuthDBMUserFile  AuthGroupFile  AuthName  AuthTypeAuthUserFile 和 Require.
* `FileInfo`
* `Indexes`
* `Limit`
* `Options`

####Order

`Order`: 用于定义基于主机的访问功能的, IP, 网络地址或主机定义访问控制机制.`Order`本身不控制,只是说明默认机制.

		Order allow,deny
		allow from
		deny from
		
		示例: 只允许 192.168.0.0/24 访问
		Order allow, deny
		Allow from 192.168.0.0/24
		处理明确allow的,默认都deny
		
		不允许 192.168.0.0/24 访问
		Order deny, allow
		Deny from 192.168.0.0/24
		这里需要修改deny和allow的顺序,处理明确deny的默认都allow
		
		拒绝 192.168.0.1 和 172.16.100.177
		Order deny,allow
		Deny from 172.16.100.177 192.168.0.1
		用空格区分多个
		
		
地址的表示方法:

* IP
* network/nermask
* HOSTNAME: www.a.com
* DOMAINNAME: magedu.com
* Partial IP: 172.16, 172.16.0.0/16
				
`AuthConfig` 基于验证一个文件中的账户和密码才能访问.	

		AuthTypes Basic
		AuthName "Restricted Files"
		AuthUserFile /usr/local/apache/passwd/passwords
		AuthGroupFile /usr/local/apache/passwd/groups
		Require user marion (只有某个用户才能登陆)
		Require group GroupName (只有某个组才能登陆)
		Require valid-user (只要出现在文件中的用户都可以登陆)


**Order**

基于主机的访问控制.

**AuthConfig**

基于用户或组的访问控制.

####elinks

纯文本浏览器

		elinks http://172.16.100.1

* `-dump`: 显示后立即退出
* `-source`: 显示网页源码

####Alias

路径别名

###整理知识点

---