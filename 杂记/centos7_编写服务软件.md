# centos7 编写服务软件

---

#### 分析系统状态

**输出激活单元**

两条命令等效

		systemctl
		
		systemctl list-units
		
**输出运行失败单元**

		systemctl --failed
		
		
#### 服务文件存储位置

所有可用的单元文件存放在 `/usr/lib/systemd/system/` 和 `/etc/systemd/system/` 目录(后者优先级更高).查看所有已安装服务:

		systemctl list-unit-files
		
#### 使用单元

一个单元配置文件可以描述如下内容之一:

* 系统服务(.service),代表一个后台服务进程,比如mysqld.这是常用的一类.
* 挂载点(.mount),此类配置单元封装文件系统结构层次中的一个挂载点.Systemd 将对这个挂载点进行监控和管理.比如可以在启动时自动将其挂载;可以在某些条件下自动卸载.`Systemd 会将/etc/fstab 中的条目都转换为挂载点,并在开机时处理`.
* socket(.socket),此类配置单元封装系统和互联网中的一个套接字.
* 系统设备(.device),此类配置单元封装一个存在于 Linux 设备树中的设备.每一个使用 udev 规则标记的设备都将会在 systemd 中作为一个设备配置单元出现.
* 交换分区(.swap)
* 文件路径(.path)
* 启动目标(.target),此类配置单元为其他配置单元进行逻辑分组.它们本身实际上并不做什么,只是引用其他配置单元而已.这样便可以对配置单元做一个统一的控制.这样就可以实现大家都已经非常熟悉的运行级别概念.比如想让系统进入图形化模式,需要运行许多服务和配置命令,这些操作都由一个个的配置单元表示,将所有这些配置单元组合为一个目标(target),就表示需要将这些配置单元全部执行一遍以便进入目标所代表的系统运行状态. (例如：multi-user.target 相当于在传统使用 SysV 的系统中运行级别 3).
* 由 systemd 管理的计时器(.timer)

使用 systemctl 控制单元时,通常需要使用单元文件的全名,包括扩展名(例如 sshd.service).但是有些单元可以在systemctl中使用简写方式.

* 如果无扩展名,systemctl 默认把扩展名当作 `.service`.例如 netcfg 和 netcfg.service 是等价的
* 挂载点会自动转化为相应的 `.mount` 单元.例如 `/home` 等价于 `home.mount`
* 设备会自动转化为相应的 `.device` 单元，所以 `/dev/sda2` 等价于 `dev-sda2.device`

有一些单元的名称包含一个`@`标记,(e.g. name@string.service),这意味着它是模板单元 name@.service 的一个 实例.string 被称作实例标识符, 在 systemctl 调用模板单元时，会将其当作一个参数传给模板单元,模板单元会使用这个传入的参数代替模板中的 %I 指示符.在实例化之前，systemd 会先检查 name@string.suffix 文件是否存在(如果存在,应该就是直接使用这个文件,而不是模板实例化了).大多数情况下,包换 @ 标记都意味着这个文件是模板.如果一个模板单元没有实例化就调用，该调用会返回失败，因为模板单元中的 `%I` 指示符没有被替换.
		
**立即激活单元**

		systemctl start <单元>
		
**立即停止单元**

		systemctl stop <单元>
		
**重启单元**

		systemctl restart <单元>
		
**命令单元重新读取配置**

		systemctl reload <单元>
		
**输出单元运行状态**

		systemctl status <单元>
		
**检查单元是否配置为自动启动**

		systemctl is-enabled <单元>
		
**开机自动激活单元**

		systemctl enable <单元>
		
**取消开机自动激活单元**

		systemctl disable <单元>
		
**重新载入 systemd,扫描新的或有变动的单元**

		systemctl daemon-reload
		
#### Systemd事务

 Systemd 能保证事务完整性.Systemd 的事务概念和数据库中的有所不同,主要是为了保证多个依赖的配置单元之间没有环形引用. 存在循环依赖,那么 systemd 将无法启动任意一个服务.此时systemd 将会尝试解决这个问题,因为配置单元之间的依赖关系有两种: `required是强依赖`;`want 则是弱依赖`,systemd 将去掉 wants 关键字指定的依赖看看是否能打破循环.如果无法修复,systemd会报错.

Systemd 能够自动检测和修复这类配置错误,极大地减轻了管理员的排错负担.		
#### Target和运行级别

Systemd 用目标（target）替代了运行级别的概念,提供了更大的灵活性,如您可以继承一个已有的目标,并添加其它服务,来创建自己的目标.

		SysV 			启动级别 			Systemd目标 	 	注释
		0 				runlevel0.target, poweroff.target 	中断系统（halt）
		1, s, single 	runlevel1.target, rescue.target 	单用户模式
		2, 4 			runlevel2.target, runlevel4.target, multi-user.target 	用户自定义启动级别，通常识别为级别3。
		3 				runlevel3.target, multi-user.target 	多用户，无图形界面。用户可以通过终端或网络登录。
		5 				runlevel5.target, graphical.target 	多用户，图形界面。继承级别3的服务，并启动图形界面服务。
		6 				runlevel6.target, reboot.target 	重启
		emergency 		emergency.target 	急救模式（Emergency shell）

		
#### 处理依赖关系

使用systemd时,可通过正确编写单元配置文件来解决其依赖关系.典型的情况是,单元A要求单元B在A启动之前运行.在此情况下,向单元A配置文件中的 `[Unit]` 段添加 Requires=B 和 After=B 即可.若此依赖关系是`可选`的,可添加 `Wants=B` 和 After=B.请注意 Wants= 和 Requires= 并不意味着 After=,即`如果 After= 选项没有制定,这两个单元将被并行启动`.

依赖关系通常被用在服务(service)而不是目标(target)上.例如, `network.target 一般会被某个配置网络接口的服务引入`.所以,将自定义的单元排在该服务之后即可.因为 network.target 已经启动.

#### 服务类型

编写自定义的 service 文件时,可以选择几种不同的服务启动方式.启动方式可通过配置文件 [Service] 段中的 Type= 参数进行设置.

* Type=`simple`(默认值): systemd认为该服务将立即启动.服务进程不会fork.如果该服务要启动其他服务,不要使用此类型启动,除非该服务是socket激活型.
* Type=`forking`: systemd认为当该服务进程fork,且父进程退出后服务启动成功.对于常规的`守护进程(daemon)`,除非你确定此启动方式无法满足需求,使用此类型启动即可.使用此启动类型应同时指定 `PIDFile=`,以便systemd能够跟踪服务的主进程.
* `Type=oneshot`: 这一选项适用于只执行一项任务、随后立即退出的服务.可能需要同时设置.RemainAfterExit=yes 使得 systemd 在服务进程退出之后仍然认为服务处于激活状态.
* `Type=notify`:与 Type=simple 相同,但约定服务会在就绪后向 systemd 发送一个信号.这一通知的实现由 libsystemd-daemon.so 提供.
* `Type=dbus`: 若以此方式启动,当指定的 BusName 出现在DBus系统总线上时,systemd认为服务就绪.
* `Type=idle`: systemd会等待所有任务(Jobs)处理完成后,才开始执行idle类型的单元.除此之外,其他行为和Type=simple 类似.

#### 定时器

定时器是以 `.timer` 为后缀的配置文件,记录由system的里面由时间触发的动作,,定时器可以替代 cron 的大部分功能.

#### 日志

systemd提供了自己日志系统(logging system),称为 journal. 使用 systemd 日志,无需额外安装日志服务(syslog).读取日志的命令:

		journalctl
		
默认情况下(当 Storage= 在文件 /etc/systemd/journald.conf 中被设置为 auto),日志记录将被写入 /var/log/journal/.该目录是 systemd 软件包的一部分.若被删除,systemd 不会自动创建它,直到下次升级软件包时重建该目录.如果该目录缺失，,ystemd 会将日志记录写入 /run/systemd/journal.这意味着,系统重启后日志将丢失.

#### 文件配置 

**[Unit]**

对这个服务的说明

* Description
* After

**[Service]**

* Type
* PIDFile为存放PID的文件路径
* ExecStart为服务的具体运行命令
* ExecReload为重启命令
* ExecStop为停止命令
* PrivateTmp=True表示给服务分配独立的临时空间

[Service]部分的启动、重启、停止命令全部要求使用绝对路径,使用相对路径则会报错.

**[Install]**

是服务安装的相关设置,包含systemctlenable或者disable的命令安装信息.


#### 示例代码

`$MAINPID` 是一个特殊变量,主进程id

		[ansible@rancher-agent-2 ~]$ cat /home/ansible/helloworld.sh
		#!/bin/bash
		
		while [ : ]
		do
		    echo $(date "+%Y-%m-%d %H:%M:%S") 'hello world!' >> /var/log/helloworld.log 2>&1
		    sleep 1
		done

		[ansible@rancher-agent-2 ~]$ cat /usr/lib/systemd/system/helloworld.service
		[Unit]
		Description=helloworld service
		After=docker.service
		
		[Service]
		Type=simple
		ExecStart=/home/ansible/helloworld.sh
		ExecStop=/bin/kill -9 $MAINPID
		
		[Install]
		WantedBy=multi-user.target
