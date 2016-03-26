#13_04_Linux进程管理之二

###笔记

---

**进程优先级关系**

有`140`种优先级.数字越小,优先级越高.

`100`-`139`: 用户可控制  
`0`-`99`: 内核调整

**O(大o)标准**

我们从队列中挑出来一个的时间随着队列的变化曲线,函数表现.

* `O(1)`: 无论队列有多长,从中挑出来一个时间是一样的.
* `O(n)`: 线性增长,队列越长,挑出来的时间越长.
* `O(log n)`:
* `O(n^2)`:
* `O(2^n)`:

**Linux队列优先级算法为O(1)**

Linux 2.6内核从队列挑选为`O(1)`,140个队列,从0开始扫描到140.找出优先级最高的.

**进程优先级高的特性**

1. 获得更多的CPU运行时间
2. 更优先获得CPU运行的机会

**进程调整优先级**

每个进程都有一个`nice`值.

`nice`: `-20`-`19` 对应 `100`-`139`. 值越小优先级越高.

默认情况下,每一个进程的`nice`值为`0`.

普通用户仅能够调大自己的进程的`nice`值.

**进程ID:pid**

pid(process id).

每个进程都有父进程,除了`init`进程(`init`进程是所有进程的父进程).

每个进程的额属性信息在`/proc`目录下(`内核映射`).

进程是`不连续`的.每一个进程的号码是唯一的,如果进程退出了,其他进程也不能占用.`init`进程号为`1`.

`最大`进程号(pid):

		[chloroplast@dev-server ~]$ cat /proc/sys/kernel/pid_max
		32768

**进程的分类**

* 跟终端相关的进程(从终端启动)
* 跟终端无关的进程

**进程的状态**

* `D`: 不可终端的睡眠
* `R`: 运行或就绪
* `S`: 可中断的睡眠
* `T`: 停止
* `Z`: 僵死态

* `<`: 高优先级
* `N`: 低优先级
* `+`: 前台进程组中的进程
* `l`: 多线程进程
* `s`: 会话进程的首进程

`示例`:
		
		[chloroplast@dev-server 1]$ ps aux | head
		USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
		root         1  0.0  0.0  19284   688 ?        Ss    2015   0:00 /sbin/init
		root         2  0.0  0.0      0     0 ?        S     2015   0:00 [kthreadd]
		root         3  0.0  0.0      0     0 ?        S     2015   0:07 [ksoftirqd/0]
		root         5  0.0  0.0      0     0 ?        S<    2015   0:00 [kworker/0:0H]
		root         7  0.0  0.0      0     0 ?        S     2015   0:00 [migration/0]
		root         8  0.0  0.0      0     0 ?        S     2015   0:00 [rcu_bh]
		root         9  0.0  0.0      0     0 ?        S     2015   0:34 [rcu_sched]
		root        10  0.0  0.0      0     0 ?        S     2015   0:03 [watchdog/0]
		root        11  0.0  0.0      0     0 ?        S<    2015   0:00 [khelper]

* `USER`: 发起进程的用户
* `PID`: 进程号
* `%CPU`: 占用CPU时间的半分比
* `%MEM`: 占据物理内存的半分比
* `VSZ`: 虚拟内存集
* `RSS`: 常驻内存集
* `TTY`: 和哪个终端相关联
* `STAT`: 状态
* `START`: 启动时间
* `TIME`: 占据CPU的运行时间
* `COMMAND`: 哪个命令启动的,如果command外面有`[]`表示是内核线程.


####相关命令

**`ps`**

`Process State`,查看进程状态.

`显示风格`:

* `SysV`风格: 选项加`-`
* `BSD`风格 

`参数(BSD)aux`:

* `a`: 显示所有跟终端有关的进程.
* `u`: 显示进程由哪个用户启动.
* `x`: ps to list all processes owned by you

`参数(SysV)elF`:

* `-elF`
* `-ef`
* `-eF`

`ps -o 属性1,属性2`: 指定显示属性

		ps -o pid,comm,ni


**`pstree`**

显示进程树

**`pgrep`**

以`grep`风格显示哪些进程, 过滤寻找进程

**`pidof`**

根据程序名查找进程号(`pid`)

		[chloroplast@dev-server 1]$ pidof sshd
		1022 831 829

**`top`**

		top - 14:24:19 up 8 days,  2:48,  1 user,  load average: 0.00, 0.01, 0.05
		Tasks:  77 total,   2 running,  75 sleeping,   0 stopped,   0 zombie
		Cpu(s):  0.0%us,  0.0%sy,  0.0%ni,100.0%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
		Mem:   1019184k total,   949296k used,    69888k free,   285388k buffers
		Swap:  1023996k total,        8k used,  1023988k free,   118268k cached
		
		  PID USER      PR  NI  VIRT  RES  SHR S %CPU %MEM    TIME+  COMMAND
		  999 root      20   0  222m 5448 3348 S  0.3  0.5  18:50.73 AliYunDun
		    1 root      20   0 19284  688  396 S  0.0  0.1   0:00.67 init
		    2 root      20   0     0    0    0 S  0.0  0.0   0:00.00 kthreadd

`信息显示`:

* `top - 14:24:19`: 当前系统时间
* `up 8 days,  2:48`: 运行时长
* `1 user`: 登录的用户
* `load average: 0.00, 0.01, 0.05`: 平均负载(过去`1分钟`的平均队列长度,过去`5分钟`的平均队列长度,过去`15分钟`的平均队列长度),值越小cpu负载越低
* `Tasks:  77 total,   2 running,  75 sleeping`: 所有进程的相关信息(一共`total`个进程,`running`个正在运行,`sleeping`正在睡眠,`stopped`停止,`zombie`僵死)
* `Cpu(s)`: cpu负载情况,按`1`cpu按编号显示(如果有多个cpu)
	* `%us`: 用户空间(`user space`)的用户进程占据
	* `%sy`: 内核空间的内核进程占据
	* `%ni`: 调整`nice`值影响的cpu比例
	* `%id`: 空闲百分比
	* `%wa`: 等待I/O完成所占据的时间
	* `%hi`: 硬件中断占据的时间
	* `%si`: 软件中断占据的时间
	* `%st`: 被偷走的时间(虚拟化场景可能被偷走)
* `Mem`: 内存
	* `total`: 总数
	* `used`: 已用数
	* `free`: 空闲数
	* `buffers`:
* `Swap`: 交换空间
	* `total`: 总数
	* `used`: 已用数
	* `free`: 空闲数
	* `cached`:

`列表显示`:

* `PID`: 进程号
* `USER`: 进程所属用户     
* `PR`: 优先级
* `NI`:  `nice`值
* `VIRT`: 虚拟内存集
* `RES`: 常驻内存集
* `SHR`: 共享内存大小(共享库)
* `S`: 状态
* `%CPU`: 占cpu半分比
* `%MEM`: 内存半分比
* `TIME+`: 真正占据cpu时长
* `COMMAND`: 进程名

`交互式子命令`:

* `M`: 根据主流内存大小排序
* `P`: 根据占据cpu百分比排序
* `T`: 根据累计时间进行排序
* `l`: 是否显示平均负载和启动时间
* `t`: 是否显示进行和cpu状态相关信息
* `m`: 是否显示内存相关信息
* `c`: 是否显示完整的命令行信息
* `q`: 退出top
* `k`: 终止某个进程

`top`:

* `-d #`: 指定刷新时长,每隔几秒显示

		top -d 1

* `-b`: 以批处理显示
* `-n #`: 在批模式下,共显示多少批 

**进程间通信(IPC: Inter Process Communication)**

* 共享内存
* 信号 Signal
* Semaphore: 信号量数组

**kill**

发送信号

`kill -l` 显示信号列表

重要的信号:

* `1`: SIGHUP,让一个进程不用重启,就可以重读其配置文件,并让新的配置信息生效.
* `2`: SIGINT,`Ctrl+c` 中断一个进程
* `9`: SIGKILL, 杀死一个进程
* `15`: SIGTERM, 终止一个进程, 给足够的时间关闭文件,释放资源,`kill默认信号`

`指定一个信号`:

* 信号号码: kill -1
* 信号名称: kill -SIGKILL
* 信号名称简写: kill -KILL

`killall COMMAND`: 杀死同一进程名的进程.指定信号的方法和`kill`一样.

 		killall COMMAND
 		
**调整nice值**

调整已经启动的进程的`nice`值:

`renice NI PID`

		renice 3 3704
		
在启动时指定`nice`值:

		nice -n NI COMMAND

**`vmstat`**

Report virtual memory statistics.

`vmstat #`: 每隔`#`秒刷新一次

`vmstat #1 #2`: 每隔`#1`秒刷新一次,但是只显示`#2`次

系统状态查看命令

		[chloroplast@dev-server ~]$ vmstat
		procs -----------memory---------- ---swap-- -----io---- --system-- -----cpu-----
		 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
		 0  0      8  69352 286032 118364    0    0     4     3   53   22  0  0 100  0  0
		 
* `procs`: 
	* `r`: 运行队列长度
	* `b`: 阻塞队列长度	 
* `memory`: 内存统计数据
	* `swpd`: 交换内存大小
	* `free`: 空闲内存大小
	* `buff`: 缓冲
	* `cache`: 缓存
* `swap`: 动态
	* `si`: swap in, 从物理内从放到交换内存
	* `so`: swap out, 从交换内存放到物理内存
* `io`:
	* `bi`: `blocks in`, 有多少个磁盘上的块放到物理内存
	* `bo`: `blocks out`, 有多少个磁盘上的块放到同步到磁盘上
* `system`:
	* `in`: 中断的个数
	* `cs`: 上下文(进程)切换的个数
* `cpu`:
	* `us`: 用户空间(`user space`)的用户进程占据
	* `sy`: 内核空间的内核进程占据
	* `id`: 空闲百分比
	* `wa`: 等待I/O完成所占据的时间
	* `st`: 被偷走的

**`free`**

**`kill`**

**`pkill`**

杀死特殊进程

**前台和后台**

`前台->后台`:

* 命令后面加`&`符号,让命令在后台`执行`
* `Ctrl+z`把正在前台的作业放到后台,在后台处于`Stopped`状态.默认是一个`Stopped`信号,停止运行.

`前台作业`: 占据命令行终端的  
`后台作业`: 启动之后,立即释放命令提示符,后续的操作在后台完成

**`jobs`**

查看后台的所有作业,每个作业有`作业号`,作业号`不同于`进程号.

* `+`: 命令将默认操作的作业
* `-`: 命令将第二个默认操作的作业

**`bg`**

让后台的停止(`Stopped`)作业继续运行.

`bg [%JOBID]`,省略`JOBID`表示默认执行`+`作业.`%`可以省略.

**`fg`**

把作业调回前台

`fg [%JOBID]`,省略`JOBID`表示默认执行`+`作业.`%`可以省略.

**杀死后台作业**

终止某作业, `kill %JOBID`.

但是`bg`和`fg`不能操作进程.

**`uptime`**

和`top`命令的信息显示一样.

		[chloroplast@dev-server ~]$ uptime
 		20:29:19 up 8 days,  8:53,  1 user,  load average: 0.00, 0.01, 0.05
 		
 **`/proc/cpuinfo`**
 
 cpu信息
 
 		[chloroplast@dev-server ~]$ cat /proc/cpuinfo
		processor	: 0
		vendor_id	: GenuineIntel
		cpu family	: 6
		model		: 62
		model name	: Intel(R) Xeon(R) CPU E5-2650 v2 @ 2.60GHz
		stepping	: 4
		microcode	: 0x428
		cpu MHz		: 2600.058
		cache size	: 20480 KB
		physical id	: 0
		siblings	: 1
		core id		: 0
		cpu cores	: 1
		apicid		: 0
		initial apicid	: 0
		fpu		: yes
		fpu_exception	: yes
		cpuid level	: 13
		wp		: yes
		flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat clflush mmx fxsr sse sse2 ht syscall nx rdtscp lm constant_tsc rep_good nopl pni ssse3 cx16 sse4_1 sse4_2 popcnt aes hypervisor lahf_lm
		bogomips	: 5200.11
		clflush size	: 64
		cache_alignment	: 64
		address sizes	: 46 bits physical, 48 bits virtual
		power management:
 
 **`/proc/meminfo`**
 
 内存信息
 
 		[chloroplast@dev-server ~]$ cat /proc/meminfo
		MemTotal:        1019184 kB
		MemFree:           69064 kB
		Buffers:          286092 kB
		Cached:           118416 kB
		SwapCached:            8 kB
		Active:           297788 kB
		Inactive:         135024 kB
		Active(anon):       2340 kB
		Inactive(anon):    26168 kB
		Active(file):     295448 kB
		Inactive(file):   108856 kB
		Unevictable:           0 kB
		Mlocked:               0 kB
		SwapTotal:       1023996 kB
		SwapFree:        1023988 kB
		Dirty:                 0 kB
		Writeback:             0 kB
		AnonPages:         28320 kB
		Mapped:            11344 kB
		...

###整理知识点

---

**buffer**

A buffer is something that has yet to be "`written`" to disk.

缓冲(buffers)是根据磁盘的读写设计的,把分散的写操作集中进行,减少磁盘碎片和硬盘的反复寻道,从而提高系统性能.

**cache** 

A cache is something that has been "`read`" from the disk and stored for later use

缓存(cached)是把读取过的数据保存起来,重新读取时若命中(找到需要的数据)就不要去读硬盘了,若没有命中就读硬盘.其中的数据会根据读取频率进行组织,把最频繁读取的内容放在最容易找到的位置,把不再读的内容不断往后排,直至从中删除.

