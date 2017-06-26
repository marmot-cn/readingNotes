### Nginx 配置

---

**user**

nginx用户及组: 用户 组.

		user nginx nginx;
		
**worker_processes**

工作进程: 数目. 根据硬件调整,通常等于CPU数量或者2倍于CPU.

		worker_processes 2;
		
		or
		
		worker_processes auto;
		设置为"auto"将尝试自动检测它.
		
**error_log**

错误日志: 存放路径.

		error_log  logs/error.log;  
		error_log  logs/error.log  notice;  
		error_log  logs/error.log  info;  
		
**pid**

pid(进程标识符): 存放路径.

		pid        /var/run/nginx.pid;
		
**worker_rlimit_nofile**

指定进程可以打开的最大描述符: 数目.

这个指令是当一个nginx进程打开的最多文件描述符数目,理论值应该是最多打开文件数(ulimt -n)与nginx进程数相除,但是nginx分配请求并不是那么均匀,所以做好与ulimit -n的值保持一致.

		worker_rlimit_nofile 65535
		
		ulimit -a
		core file size          (blocks, -c) 0
		data seg size           (kbytes, -d) unlimited
		scheduling priority             (-e) 0
		file size               (blocks, -f) unlimited
		pending signals                 (-i) 3888
		max locked memory       (kbytes, -l) 64
		max memory size         (kbytes, -m) unlimited
		open files                      (-n) 65535
		pipe size            (512 bytes, -p) 8
		POSIX message queues     (bytes, -q) 819200
		real-time priority              (-r) 0
		stack size              (kbytes, -s) 8192
		cpu time               (seconds, -t) unlimited
		max user processes              (-u) 3888
		virtual memory          (kbytes, -v) unlimited
		file locks                      (-x) unlimited

### Events模块

Events 模块中包含`Nginx`中所有处理链接的设置.


		events {
		    use epoll;
		    worker_connections  2048;
		    multi_accept on; 
		}
		
**use**

`use epoll;`

		使用epoll的I/O模型.linux建议epoll,FreeBSD建议采用kqueue,window下不指定.
		
		补充说明:
		与apache相类,nginx针对不同的操作系统,有不同的事件模型
		a) 标准事件模型
		Select,poll属于标准事件模型,如果当前系统不存在更有效的方法,nginx会选择select或poll
		b) 高效事件模型
		epoll
		
**worker_connections**		
		
`worker_connections xxx;`

		单个后台worker process进程的最大并发链接数.
		并发总数是 worker_processes 和 worker_connections 的乘积.
		即 max_clients = worker_processes * worker_connections
		在设置了反向代理的情况下,max_clients = worker_processes * worker_connections / 4 为什么上面反向代理要除以4,应该说是一个经验值.
		因为并发受IO约束,max_clients的值须小于系统可以打开的最大文件数.
		而系统可以打开的最大文件数和内存大小成正比,一般1GB内存的机器上可以打开的文件数大约是10万左右.
		
**multi_accept**

		默认是on.设置为on后,多个worker按串行方式来处理连接,也就是一个连接只有一个worker被唤醒,其他的处于休眠状态.

		设置为off后,多个worker按并行方式来处理连接,也就是一个连接会唤醒所有的worker,知道连接分配完毕,没有取得连接的继续休眠.

		当你的服务器连接数不多时,开启这个参数会让负载有一定程度的降低.但是当服务器的吞吐量很大时，为了效率,请关闭这个参数.
		
### http模块

**server_tokens**

		server_tokens off;

		错误页面的标签上是否表示 Nginx的版本.
		安全上的考虑还是off掉吧.
		
**include**

		include /etc/nginx/mime.types;

		定义MIME类型和后缀名关联的文件的位置。
		
		types {
		
			text/html html htm shtml;
			
			text/css css;
			
			text/xml xml;
			
			image/gif gif;
			
			image/jpeg jpeg jpg;
			
			application/javascript js;
			
			...
		
		}
		
		mime.types文件中大概是这个样子的.	
	
**default_type**

		default_type application/octet-stream;

		指定mime.types文件中没有记述到的后缀名的处理方法。
		默认值是text/plain.
		
**log_format**

		log_format main 'time:$time_iso8601\t'...

		log_format ltsv 'time:$time_iso8601\t'... 
		
		定义日志的格式. 可以选择main或者ltsv,后面定义要输出的内容.
		 1.$remote_addr 与 $http_x_forwarded_for 用以记录客户端的ip地址;
		 2.$remote_user:用来记录客户端用户名称;
		 3.$time_local:用来记录访问时间与时区;
		 4.$request:用来记录请求的url与http协议;
		 5.$status:用来记录请求状态;
		 6.$body_bytes_s ent:记录发送给客户端文件主体内容大小;
		 7.$http_referer:用来记录从那个页面链接访问过来的;
		 8.$http_user_agent:记录客户端浏览器的相关信息;
		 
**access_log**

		access_log /var/log/nginx/access.log main;

		连接日志的路径,上面指定的日志格式放在最后。

		access_log off;

		也可以关掉.
		
**charset**

		charset UTF-8;

		设置应答的文字格式.
		
**sendfile**

		sendfile off;
		
		指定是否使用OS的sendfile函数来传输文件.
		普通应用应该设为on,下载等IO重负荷的应用应该设为off.默认值是off.
		
**tcp_nopush**

		tcp_nopush on; 
		sendfile为on时这里也应该设为on,数据包会累积一下再一起传输,可以提高一些传输效率.
		
**tcp_nodelay**

		tcp_nodelay on;
		小的数据包不等待直接传输. 默认为on.
		看上去是和tcp_nopush相反的功能,但是两边都为on时nginx也可以平衡这两个功能的使用.
		
**keepalive_timeout**
		
		keepalive_timeout 60;
		keepalive超时时间
		
**keepalive_requests**

		keepalive_requests 100; 
		
		keepalive_timeout时效内同样的客户端超过指定数量的连接时会被强制切断.
		一般的话keepalive_timeout 5和keepalive_requests 20差不多就够了.
		默认为100.
		
**set_real_ip_from 和 real_ip_header**
	
		set_real_ip_from 10.0.0.0/8; 
		real_ip_header X-Forwarded-For; 
		
		可以防止经过代理或者负载均衡服务器时丢失源IP.
		set_real_ip_from指定代理或者负载均衡服务器的IP,可以指定复数个IP.
		real_ip_header指定从哪个header头检索出要的IP地址.
		
**client_header_timeout 和 client_body_timeout**

		client_header_timeout 10; 
		client_body_timeout 10; 
		
		读取客户端的请求head部分和客户端的请求body部分的超时时间.
		
**client_body_buffer_size 和 client_body_temp_path**

		client_body_buffer_size 32k; 
		client_body_temp_path /dev/shm/client_body_temp 1 2; 
		
		接受的请求body部分到client_body_buffer_size为止放在内存中,超出的部分输出至client_body_temp_path文件里.

**client_max_body_size**

		client_max_body_size 1m;

		客户端上传的body的最大值.超过最大值就会发生413(Request Entity Too Large)错误.
		默认为1m,最好改大一点.
		
**client_header_buffer_size和large_client_header_buffers**
		
		client_header_buffer_size 4k;
		
		large_client_header_buffers 4 8k; 

		客户端请求头部的缓冲区大小.这个可以根据你的系统分页大小来设置,一般一个请求头的大小不会超过1k,不过由于一般系统分页都要大于1k,所以这里设置为分页大小.
		
		[ansible@rancher-agent-1 ~]$ getconf PAGE_SIZE
		4096
		即 4KB
		
		但也有client_header_buffer_size超过4k的情况,但是client_header_buffer_size该值必须设置为"系统分页大小"的整倍数.
		
		一般来说默认就够了。
		发生414 (Request-URI Too Large) 错误时请增大这两个参数.
		
**limit_conn和limit_conn_zone**

		limit_conn_zone $binary_remote_addr zone=addr:10m;

		limit_conn addr 100;

		限制某条件下的同时连接数.

**open_file_cache**		
		
		open_file_cache max=65535 inactive=60s;
		
		缓存静态文件.
		这个将为打开文件指定缓存,默认是没有启用的,max指定缓存数量,建议和打开文件数一致,inactive是指经过多长时间文件没被请求后删除缓存.

**open_file_cache_valid**
		
		open_file_cache_valid 80s;

		这个是指多长时间检查一次缓存的有效信息.
		
**open_file_cache_min_uses**	
		
		open_file_cache_min_uses 1;

		open_file_cache指令中的inactive参数时间内文件的最少使用次数,如果超过这个数字,文件描述符一直是在缓存中打开的,如上例,如果有一个文件在inactive时间内一次没被使用,它将被移除.
		
**open_file_cache_errors**

		open_file_cache_errors on | off 默认值:open_file_cache_errors off
		这个指令指定是否在搜索一个文件是记录cache错误
		
**server_names_hash_bucket_size**

		server_names_hash_bucket_size 64 
		
		nginx启动时出现could not build the server_names_hash, you should increase错误时请提高这个参数的值.
		一般设成64就够了.
		保存服务器名字的hash表是由指令server_names_hash_max_size 和server_names_hash_bucket_size所控制的。参数hash bucket size总是等于hash表的大小，并且是一路处理器缓存大小的倍数。在减少了在内存中的存取次数后，使在处理器中加速查找hash表键值成为可能。如果hash bucket size等于一路处理器缓存的大小，那么在查找键的时候，最坏的情况下在内存中查找的次数为2。第一次是确定存储单元的地址，第二次是在存储单元中查找键 值。因此，如果Nginx给出需要增大hash max size 或 hash bucket size的提示，那么首要的是增大前一个参数的大小.
		
**types_hash_max_size**

		types_hash_max_size 1024; 
		
		types_hash_max_size影响散列表的冲突率.types_hash_max_size越大,就会消耗更多的内存，但散列key的冲突率会降低,检索速度就更快.types_hash_max_size越小,消耗的内存就越小,但散列key的冲突率可能上升.
		默认为1024 	
		
**types_hash_bucket_size**

		types_hash_bucket_size 64; 
		
		types_hash_bucket_size 设置了每个散列桶占用的内存大小。
		默认为64 	
		
**listen**

		listen 80 default_server; 
		
		ginx当网页服务器使用的时候写在http模块中,一般来说用作虚拟主机的情况下不会写在这里.
		
**server_name_in_redirect**

		server_name_in_redirect off; 
		重定向的时候需不需要把服务器名写入head,基本上不会设成on.
		
**port_in_redirect**

		port_in_redirect on; 
		
		设为on后,重定向的时候URL末尾会带上端口号.
		
#### upstream

		upstream bakend {
		
			server 127.0.0.1:8027;
		
			server 127.0.0.1:8028;
		
			server 127.0.0.1:8029;
		
			hash $request_uri;
	
		}
		
##### 轮询

nginx的upstream目前支持4种方式的分配.

**轮询(默认)**

每个请求按时间顺序逐一分配到不同的后端服务器,如果后端服务器down掉,能自动剔除.

**weight**

		指定轮询几率,weight和访问比率成正比,用于后端服务器性能不均的情况.
		
		例如:
		upstream bakend {
			server 192.168.0.14 weight=10;
			server 192.168.0.15 weight=10;
		}

**ip_hash**

		
      每个请求按访问ip的hash结果分配,这样每个访客固定访问一个后端服务器,可以解决session的问题.
      
      upstream bakend {
		ip_hash;
		server 192.168.0.14:88;
		server 192.168.0.15:80;
	  }
	  
**fair(第三方)**

		按后端服务器的响应时间来分配请求,响应时间短的优先分配.
		
		upstream backend {
			server server1;
			server server2;
			fair;
		}
		
**url_hash**

		按访问url的hash结果来分配请求,使每个url定向到同一个后端服务器,后端服务器为缓存时比较有效.
		
		例:在upstream中加入hash语句,server语句中不能写入weight等其他的参数,hash_method是使用的hash算法.

		upstream backend {
			server squid1:3128;
			server squid2:3128;
			hash $request_uri;
			hash_method crc32;
		}

##### 每个设备的状态设置        
        
      
**down**

		不参与负责均衡
		
**weight**

		weight越大,负载的权重就越大
		
**max_fails**

		允许请求失败的次数默认位1,当超过最大次数时,返回proxy_net_upstream模块定义的错误.
		
**fail_timemout**

		max_fails 次数失败后,暂停的时间
		
**backup**

		其他的非backup机器都down或者忙的时候才会请求到这台机器
		
##### 示例

		tips:

		upstream bakend{#定义负载均衡设备的Ip及设备状态
				ip_hash;
				server 127.0.0.1:9090 down;
				server 127.0.0.1:8080 weight=2;
				server 127.0.0.1:6060;
				server 127.0.0.1:7070 backup;
		}
		在需要使用负载均衡的server中增加
		proxy_pass http://bakend/;

### server模块

		http {

			server {

				...

       		}

     	} 
     	
		Nginx用作虚拟主机时使用.
		每一个server模块生成一个虚拟主机.
		写在http模块内部.
		
**listen**

		listen 80;
		配置监听端口
		
**server_name**

		server_name localhost; 
		server_name 指定服务器的域名
		
**root**

		root /path/public 
		定义服务器的默认网站根目录位置.
		
**rewrite**

		rewrite /(.*)/index.html $1.html permanent; 
		需要重定向的时候使用.
		
**satisfy,auth_basic**
		
		satisfy any;

		auth_basic "basic authentication";
		auth_basic_user_file /etc/nginx/.htpasswd;
		satisfy any|all 部分地址Basic认证的方式
				
						allow		Deny
		satisfy any    不认证		Basic认证
		satisfy all    Basic认证    拒绝连接
		
		auth_basic：认证的名称
		auth_basic_user_file：密码文件
			
**try_files**

		try_files $uri $uri.html $uri/index.html @unicorn; 
		
		从左边开始找指定文件是否存在.
		比如连接http://***/hoge时按hoge.html、hoge/index.html、location @unicorn {}的顺序查找.
		
### Location模块

		http {

		server {
		
		location / {
		
		...
		
		           }
		
		       }
		
		     } 		
		     
指定位置(文件或路径)时使用.
也可以用正则表达式。

		location ~ /\.(ht|svn|git) {
		
			deny all; 
		}
		
		不想让用户连接.htaccess，.svn，.git文件时用上面的设置
		
		
		location ~* \.(mp3|exe)$ {

		对以“mp3或exe”结尾的地址进行负载均衡
		
		 
			proxy_pass http://img_relay$request_uri;
		
			设置被代理服务器的端口或套接字，以及URL
		
			proxy_set_header Host $host;
		
			proxy_set_header X-Real-IP $remote_addr;
		
			proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		
			以上三行，目的是将代理服务器收到的用户的信息传到真实服务器上
		
		}
		
**stub_status**

		stub_status，allow，deny
		stub_status on;
		access_log off;
		allow 127.0.0.1;
		deny all;

		stub_status连接指定的位置时可以显示现在的连接数.一般来说不会公开.
		Allow允许指定IP的连接.
		Deny拒绝指定IP的连接.
		Nginx规则是从上到下判断的,上面的例子先判断的是allow 127.0.0.1,所以127.0.0.1允许连接.
		如果反过来设成下面这样,所有的IP都被拒绝了.(注意和Apache不一样).

		deny all;

		allow 127.0.0.1; 
		
**expires**

		expires 10d;
		使用浏览器缓存时设置.上面的例子用了10天内的浏览器缓存.
		
**add_header**

		add_header Cache-Control public;
		设置插入response header的值.
		
**break,last**

		break|last; 
		
		rewrite后接break指令,完成rewrite之后会执行完当前的location(或者是if)指令里的其他内容(停止执行当前这一轮的ngx_http_rewrite_module指令集),然后不进行新URL的重新匹配.
		rewrite后接last指令,在完成rewrite之后停止执行当前这一轮的ngx_http_rewrite_module指令集已经后续的指令,进而为新的URL寻找location匹配.
		
**internal**

		 语法：internal
		默认值：no
		使用字段： location
		internal指令指定某个location只能被“内部的”请求调用，外部的调用请求会返回”Not found” (404)
		“内部的”是指下列类型：
		
		指令error_page重定向的请求。
		ngx_http_ssi_module模块中使用include virtual指令创建的某些子请求。
		ngx_http_rewrite_module模块中使用rewrite指令修改的请求。
		一个防止错误页面被用户直接访问的例子：
		
		error_page 404 /404.html;
		location  /404.html {
		  internal;
		} 
		
### proxy相关

作为反向代理时需要用到的参数

		proxy_buffering on;

		proxy_buffer_size 8k;

		proxy_buffers 100 8k;

		proxy_cache_path /var/lib/nginx/cache levels=1:2 keys_zone=CACHE:512m inactive=1d max_size=60g; 
		
		