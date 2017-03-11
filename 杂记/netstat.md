#netstat

---

用于显示与IP,TCP,UDP和ICMP协议相关的统计数据,一般用于检验本机各端口的网络连接情况
主要作用是:查看端口使用情况.

###参数含义

* `-a` (all) 显示所有内容
* `-t`  (tcp) 仅显示tcp相关内容
* `-u` (udp) 仅显示udp相关内容
* `-n` (numeric) 直接显示ip地址以及端口，不解析
* `-l`  (listen) 仅列出 Listen (监听) 的服务
* `-p` (pid) 显示出socket所属的进程PID 以及进程名字
* `-r` 显示路由信息，路由表
* `-e` 显示扩展信息，例如uid等
* `-s` 按各个协议进行统计
* `-c` 每隔一个固定时间,执行该netstat命令

**示例**

我主机上运行了rancher server,监听了`8080`端口.

		[ansible@rancher-server ~]$ netstat -atl
		Active Internet connections (servers and established)
		Proto Recv-Q Send-Q Local Address           Foreign Address         State
		tcp        0      0 0.0.0.0:17456           0.0.0.0:*               LISTEN
		tcp        0      0 iZ94xwu3is8Z:47630      100.100.35.4:squid      ESTABLISHED
		tcp        0      0 rancher-server:36498    106.11.68.13:http       ESTABLISHED
		tcp        0    316 rancher-server:17456    36.46.49.161:19842      ESTABLISHED
		tcp6       0      0 [::]:webcache           [::]:*                  LISTEN
		
		[ansible@rancher-server ~]$ netstat -antl
		Active Internet connections (servers and established)
		Proto Recv-Q Send-Q Local Address           Foreign Address         State
		tcp        0      0 0.0.0.0:17456           0.0.0.0:*               LISTEN
		tcp        0      0 10.44.88.189:47630      100.100.35.4:3128       ESTABLISHED
		tcp        0      0 120.24.3.210:36498      106.11.68.13:80         ESTABLISHED
		tcp        0    316 120.24.3.210:17456      36.46.49.161:19842      ESTABLISHED
		tcp6       0      0 :::8080                 :::*                    LISTEN
		

注意我没加 -n 参数的时候 8080 端口显示为 webcache. 端口被替换成人类可读的名字,这些是和 `/etc/services` 文件内的映射一一对应的.