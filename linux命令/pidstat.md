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