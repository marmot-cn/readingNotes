# 16_03_Linux日志系统syslog

---

### 笔记

---

Linux上的日志系统

* syslog
* syslog-ng: 开源

#### 日志系统: syslog()

syslog是服务,负责为各程序记录日志,可以把每个程序理解为子系统.但是类似nginx等,使用的是自己的日志系统.

syslog服务有两个进程:

* `syslogd`: 系统,专门用于非内核所产生的日志.
* `klogd`: 内核,专门用于负责内核所产生的日志.

kernel --> 物理终端(/dev/console) --> 记录到 /var/log/dmesg,(由`klogd`负责记录)

		[ansible@k8s-minon2-test ~]$ cat /var/log/dmesg
		[    0.000000] Initializing cgroup subsys cpuset
		[    0.000000] Initializing cgroup subsys cpu
		[    0.000000] Initializing cgroup subsys cpuacct
		[    0.000000] Linux version 3.10.0-123.9.3.el7.x86_64 (builder@kbuilder.dev.centos.org) (gcc version 4.8.2 20140120 (Red Hat 4.8.2-16) (GCC) ) #1 SMP Thu Nov 6 15:06:03 UTC 2014
		[    0.000000] Command line: BOOT_IMAGE=/boot/vmlinuz-3.10.0-123.9.3.el7.x86_64 root=UUID=6634633e-001d-43ba-8fab-202f1df93339 ro crashkernel=auto vconsole.font=latarcyrheb-sun16 vconsole.keymap=us biosdevname=0 rhgb quiet LANG=en_US.UTF-8

一旦系统控制权由内核转交给`/sbin/init`,日志记录转交给`klogd`,记录的日志信息会放到

* `/var/log/messages`: 非内核产生引导信息,系统初始化信息,各子系统所产生的信息.会多次滚动,以免日志文件变得过大.
* `/var/log/maillog`: 邮件系统产生的日志信息.
* `/var/log/secure`(是600权限): 任何用户登录时所产生的登录信息.

**日志需要滚动(日志切割)**:每隔一段时间,原来的日志文件重新命名为`xx.xx(数字)`,然后再把日志记录到`xx`.例如一段时间`messages`变为`messages.1`,新的日志会记录到`messages`.日志信息应该多次滚动,以免日志信息过大,难以分析.

`logrotate`: 专门用于日志切割的程序.

`/etc/cron.daily/logrotate`: 专门用于日志切割的计划任务.

`/etc/logrotate.d`: 目录下的每个文件定义每个子系统的日志格式.

**子系统**: facility,设施

**信息详细程度**:日志级别

**指定信息的存储位置**:动作

#### syslog配置文件: `/etc/syslog.conf`

配置文件定义格式为:

		facility.priority	action
		
`-action`: 加上`-`表示异步写入,先写到内存过一会在同步到磁盘上去.
		
**facility(设备)**

facility,可以理解为日志的来源或设备.目前常用的facility有以下几种:

* `auth`: 认证相关的
* `authpriv`: 认证,授权相关的
* `cron`: 任务计划相关的
* `daemon`: 守护进程相关的
* `kern`: 内核相关的
* `lpr`: 打印相关的
* `mail`: 邮件相关的
* `mark`: 标记相关的
* `news`: 新闻相关的
* `security`: 安全相关的,与`auth`类似
* `syslog`: `syslog`自己的
* `user`: 用户相关的
* `uucp`: `unix to unix cp`相关的
* `local10 到 local17`: 用户自定义使用
* `*`: 表示又有`facility`

**priority**

priority(log level)日志的级别,一般有以下几种级别(从低到高)

* `debug`: 程序或系统的调试信息
* `info`: 一般信息
* `notice`: 不影响正常功能,需要注意的消息
* `warning/warn`: 可能影响系统功能,需要提醒用户的重要事件
* `err/error`: 错误信息
* `crit`: 比较严重的
* `alert`: 必须马上处理的
* `emerg/panic(恐慌)`: 会导致系统不可用的
* `*`: 所有的日志级别
* `none`: 跟`*`相反,表示啥也没有

**action**

action(动作)日志记录的位置

* `系统上的绝对路径`: 普通文件 入: /var/log/xxx
* `|` : 管道,日志信息送给其他命令处理
* `终端`: 终端 如: /dev/console,输出到物理设备
* `@HOST`: 远程主机 如: `@10.0.0.1`,日志发送给其他主机
* `用户`: 系统用户 如: `root`,日志发送给某个用户
* `*`: 登陆到系统上的用户,一般`emerg`级别的日志是这样定义的,日志发送给所有用户

**示例**

`mail.info	/var/log/mail.info`:表示将`mail`相关的,级别为`info`即`info`以上级别的信息记录到`/var/log/mail.log`文件中.

`auth.=info @10.0.0.1`:表示将`auth`先关的,级别为`info`的信息记录到`10.0.0.1`主机上去.`.=`表示精确记录一个级别.

`user.!=error`:表示记录`user`相关的,不包括`error`级别的信息

`user.!error`:与`user.error`相反.只记录比`error`级别低的信息

`*.info`:表示记录所有的日志信息的`info`级别

`mail.*`:表示记录`mail`相关的所有级别的信息

`*.*`:所有

`cron.info;mail.info`:多个日志来源可以用`;`分隔开

`cron,mail.info`:与`cron.info;mail.info`是一个意思

`mail.*;mail.!=info`:表示记录`mail`相关的所有级别的信息,但是不包括`info`级别的


### 整理知识点

---

#### service xxx reload

发送`1`号信号`SIGUP`,让服务不用重启,就可以重读配置文件. 

`SIGHUP` as a notification about `terminal closing event` doesn't make sense for a daemon, because deamons are detached from their terminal. So the system will never send this signal to them. Then it is common practice for daemons to use it for another meaning, typically reloading the daemon's configuration. This is not a rule, just kind of a convention. That's why it's not documented in the manpage.