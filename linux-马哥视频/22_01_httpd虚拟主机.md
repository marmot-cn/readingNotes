#22_01_httpd虚拟主机

---

###笔记

---

####apache虚拟主机

apache: 运行服务器,通常称为Host,主机.

虚拟主机:

* IP
* 80 端口

资源有限.

* 基于IP的(多个网卡)
	* IP1:80
	* IP2:80
* 基于端口
	* IP:80
	* IP:8080
* 基于域名
	* 主机名不同
	* 解析到同一个IP上 

**示例**
			
		NameVirtualHost	172.16.100.2:80
		
		<VirtualHost 172.16.100.1:80>
			ServerName hello.magedu.com
			DocumentRoot "/www/magedu.com"
			CustomLog	/var/log/httpd/magedu.com/access_log combined
		</VirtualHost>
		
		<VirtualHost 172.16.100.2:80>
			ServerName www.a.org
			DocumentRoot "/www/a.org"
			CustomLog	/var/log/httpd/a.org/access_log combined
		</VirtualHost>
		
		<VirtualHost 172.16.100.1:8080>
			ServerName www.b.net
			DocumentRoot "/www/b.net"
			CustomLog	/var/log/httpd/b.net/access_log combined
		</VirtualHost>

		我们创建好相应的DocumentRoot文件目录,创建好相应的log目录.
		
		1.给网卡添加IP
		ip addr add 172.16.100.2/16 dev eth0
		2.监听 8080 端口
		编辑 httpd.conf, 添加 Listen 8080
		
		访问 172.16.100.1 会访问到 /www/magedu.com/index.html
		访问 172.16.100.2 会访问到 /www/a.org/index.html
		访问 172.16.100.1:8080 会访问到 /www/b.net/index.html
		
**基于域名**

		<VirtualHost 172.16.100.2:80>
			ServerName www.a.org
			DocumentRoot "/www/a.org"
		</VirtualHost>
		
		<VirtualHost 172.16.100.2:80>
			ServerName www.d.gov
			DocumentRoot "/www/d.gov"
			<Directory	"/www/d.gov">
				Options none
				AllowOverride none
				Order deny,allow
				Deny from 172.16.100.177
			</Directory>
		</VirtualHost>
	
		172.16.100.177 来源的ip 禁止访问 www.d.gov
		
		访问 172.16.100.2 会返回符合的第一个 即 /www/a.org/index.html
		访问 www.a.org 会访问到 /www/a.org/index.html
		访问 www.d.gov 会访问到 /www/d.gov/index.html

###整理知识点

---