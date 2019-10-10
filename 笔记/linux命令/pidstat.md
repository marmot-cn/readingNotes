# pidstat

---

## 简介

`pidstat`主要用于监控全部或这顶进程占用系统资源的情况, 如CPU, 内存, 设备IO, 任务切换, 线程等.
`pidstat`首次运行时显示自系统启动开始的各项统计信息, 之后运行`pidstat`将显示自上次运行该命令以后的统计信息. 

### 参数

* 列出`IO`统计信息: `-d`
* 列出内存使用统计: `-r`
* CPU统计信息: `-u`
* 上下文切换统计信息: `-w`

## 示例

### 默认参数

运行`pidstat`，将输出系统启动后所有活动进程的cpu统计信息. 因而即使当前某进程的cpu占用率很高，输出中的值有可能仍为0

```
ansible@demo ~]$ pidstat
Linux 3.10.0-514.26.2.el7.x86_64 (demo) 	12/15/2017 	_x86_64_	(1 CPU)

11:57:22 PM   UID       PID    %usr %system  %guest    %CPU   CPU  Command
11:57:22 PM     0         1    0.00    0.00    0.00    0.00     0  systemd
11:57:22 PM     0         2    0.00    0.00    0.00    0.00     0  kthreadd
11:57:22 PM     0         3    0.00    0.00    0.00    0.00     0  ksoftirqd/0
11:57:22 PM     0         9    0.00    0.00    0.00    0.00     0  rcu_sched
11:57:22 PM     0        10    0.00    0.00    0.00    0.00     0  watchdog/0
11:57:22 PM     0        15    0.00    0.00    0.00    0.00     0  xenbus
```

* `PID`: 进程pid.
* `%usr`: 进程在用户态运行所占cpu时间比率.
* `%system`: 进程在内核态运行所占cpu时间比率.
* `%CPU`: 进程运行所占cpu时间比率.
* `CPU`: 只是进程在哪个核运行.
* `Command`: 拉起进程对应的命令.

### 指定采样周期和采样次数

`pidstat`命令指定采样周期和采样次数, 命令形式为`pidstat [option] interval [count]`.

`pidstat`输出以2秒为采样周期, 输出10此cpu使用统计信息.

```
pidstat 2 10
```

### cpu 使用情况统计(-u)

默认是`-u`即, `pidstat -u`和`pidstat`效果一样.

### 内存使用情况统计(-r)

```
[ansible@demo ~]$ pidstat -r -p 1 1 10
Linux 3.10.0-514.26.2.el7.x86_64 (demo) 	12/16/2017 	_x86_64_	(1 CPU)

12:07:01 AM   UID       PID  minflt/s  majflt/s     VSZ    RSS   %MEM  Command
12:07:02 AM     0         1      0.00      0.00  125176   3456   0.34  systemd
12:07:03 AM     0         1      0.00      0.00  125176   3456   0.34  systemd
12:07:04 AM     0         1      0.00      0.00  125176   3456   0.34  systemd
12:07:05 AM     0         1      0.00      0.00  125176   3456   0.34  systemd
12:07:06 AM     0         1      0.00      0.00  125176   3456   0.34  systemd
12:07:07 AM     0         1      0.00      0.00  125176   3456   0.34  systemd
```

* `minflt/s`: 每秒次缺页错误次数(minor page faults), 次缺页错误次数意即虚拟内存地址映射成物理内存地址产生的page fault次数/
* `majflt/s`: 每秒主缺页错误次数(major page faults), 当虚拟内存地址映射成物理内存地址时, 相应的page在swap中, 这样的page fault为major page fault，一般在内存使用紧张时产生.
* `CSZ`: 该进程使用的虚拟内存(以kB为单位).
* `RSS`: 该进程使用的物理内存(以kB为单位).
* `%MEM`: 该进程使用内存的百分比.
* `Command`: 拉起进程对应的命令.

### IO情况统计(-d)

```
[ansible@demo ~]$ pidstat -d -p 1 1 10
Linux 3.10.0-514.26.2.el7.x86_64 (demo) 	12/16/2017 	_x86_64_	(1 CPU)

12:10:14 AM   UID       PID   kB_rd/s   kB_wr/s kB_ccwr/s  Command
12:10:15 AM     0         1      0.00      0.00      0.00  systemd
12:10:16 AM     0         1      0.00      0.00      0.00  systemd
```

* `kB_rd/s`: 每秒进程从磁盘读取的数据量(以kB为单位).
* `kB_wr/s`: 每秒进程向磁盘写的数据量(以kB为单位).
* `kB_ccwr/s`: 任务写入磁盘被取消的速率.
* `Command`: 拉起进程对应的命令.

### 针对特定进程统计(-p)

使用`-p`选项，我们可以查看特定进程的系统资源使用情况.

### 显示每个进程的上下文切换情况(-w)

```
[root@iZ94ebqp9jtZ ~]# pidstat -w -p 18060 1 10
Linux 3.10.0-514.26.2.el7.x86_64 (iZ94ebqp9jtZ) 	04/12/2018 	_x86_64_	(1 CPU)

11:44:12 AM   UID       PID   cswch/s nvcswch/s  Command
11:44:13 AM     0     18060     10.10      0.00  AliYunDun
11:44:14 AM     0     18060     10.10      0.00  AliYunDun
11:44:15 AM     0     18060     10.10      0.00  AliYunDun
11:44:16 AM     0     18060     10.10      0.00  AliYunDun
11:44:17 AM     0     18060     10.00      0.00  AliYunDun
11:44:18 AM     0     18060     10.00      0.00  AliYunDun
```

* `PID`: 进程id.
* `cswch/s`: 每秒主动任务上下文切换(自愿上下文切换)数量. 当某一任务处于阻塞等待时, 将主动让出自己的CPU资源.
* `nvcswch/s`: 每秒被动任务上下文切换(非自愿上下文切换)数量. CPU分配给某一任务的时间片已经用完, 因此将强迫该进程让出CPU的执行权.
* `Command`: 命令名.

### 显示选择任务的线程的统计信息外的额外信息(-t)

```
[root@iZ94ebqp9jtZ ~]# pidstat -t -p 18060 1 10
Linux 3.10.0-514.26.2.el7.x86_64 (iZ94ebqp9jtZ) 	04/12/2018 	_x86_64_	(1 CPU)

12:02:40 PM   UID      TGID       TID    %usr %system  %guest    %CPU   CPU  Command
12:02:41 PM     0     18060         -    0.00    0.00    0.00    0.00     0  AliYunDun
12:02:41 PM     0         -     18060    0.00    0.00    0.00    0.00     0  |__AliYunDun
12:02:41 PM     0         -     18061    0.00    0.00    0.00    0.00     0  |__AliYunDun
12:02:41 PM     0         -     18062    0.00    0.00    0.00    0.00     0  |__AliYunDun
12:02:41 PM     0         -     18073    0.00    0.00    0.00    0.00     0  |__AliYunDun
12:02:41 PM     0         -     18074    0.00    0.00    0.00    0.00     0  |__AliYunDun
```

* `TGID`: 主线程id.
* `TID`: 线程id.
* `%usr`：进程在用户空间占用cpu的百分比.
* `%system`：进程在内核空间占用cpu的百分比.
* `%guest`：进程在虚拟机占用cpu的百分比.
* `%CPU`：进程占用cpu的百分比.
* `CPU`：处理进程的cpu编号.
* `Command`：当前进程对应的命令.

### pidstat -T

```
pidstat -T TASK
pidstat -T CHILD
pidstat -T ALL
```
* `TASK`: 表示独立的`task`.
* `CHILD`: 关键字表示报告进程下所有线程统计信息.
* `AL`L 表示报告独立的task和task下面的所有线程.

```
[root@iZ94ebqp9jtZ ~]# pidstat -T ALL -p 18060 1 10
Linux 3.10.0-514.26.2.el7.x86_64 (iZ94ebqp9jtZ) 	04/12/2018 	_x86_64_	(1 CPU)

02:48:51 PM   UID       PID    %usr %system  %guest    %CPU   CPU  Command
02:48:52 PM     0     18060    0.00    0.00    0.00    0.00     0  AliYunDun

02:48:51 PM   UID       PID    usr-ms system-ms  guest-ms  Command
02:48:52 PM     0     18060         0         0         0  AliYunDun
02:48:53 PM     0     18060    1.01    0.00    0.00    1.01     0  AliYunDun
02:48:53 PM     0     18060        10         0         0  AliYunDun
02:48:54 PM     0     18060    0.00    0.00    0.00    0.00     0  AliYunDun
02:48:54 PM     0     18060         0         0         0  AliYunDun
02:48:55 PM     0     18060    0.00    1.00    0.00    1.00     0  AliYunDun
```

* `PID`: 进程id.
* `Usr-ms`: 任务和子线程在**用户级别**使用的毫秒数.
* `System-ms`: 任务和子线程在**系统级别**使用的毫秒数.
* `Guest-ms`: 任务和子线程在虚拟机(`running a virtual processor`)使用的毫秒数.
* `Command`: 命令名.

