# PHP 日志`%c`参数答疑

---

### 几个关键参数

* `%d`: 响应时间.
* `%C`: CPU处理时间占响应时间的百分比.
* `%M`: 峰值内存占用.

### 当时处理`%C`的几点疑惑

我在日志里面记录了`%C`, 查看英文解释如下**%CPU used by the request**, 一直以为是占用`CPU`的百分比,  查看日志就好奇为什么`CPU`动辄就占用超过`100%`.

后来查看了是**CPU处理时间占响应时间的百分比**.

#### 查看源码

[php源码](https://github.com/php/php-src), 文件路径`sapi/fpm/fpm/fpm_request.c
`

几个知识点:

##### `times`头文件

文件开头引用

```
#include <sys/times.h>
```

这个里面包含了`tms`的结构体.

```
The <sys/times.h> header shall define the structure tms, which is returned by times() and includes at least the following members:

clock_t  tms_utime  User CPU time. 
clock_t  tms_stime  System CPU time. 
clock_t  tms_cutime User CPU time of terminated child processes. 
clock_t  tms_cstime System CPU time of terminated child processes. 


clock_t times(struct tms *);
```

##### fpm_request 文件分析

在每次请求开始(`fpm_request_reading_headers`)和结束(`fpm_request_end`), 使用`times`获得当前进程`cpu`时间.

* `FPM_REQUEST_ACCEPTING`: Idle 空闲状态.
* `FPM_REQUEST_READING_HEADERS`: Reading headers 读取头信息.
* `FPM_REQUEST_INFO`: Getting request informations 获取请求信息.
* `FPM_REQUEST_EXECUTING`: Running 执行
* `FPM_REQUEST_END`: Ending 完结
* `FPM_REQUEST_FINISHED`: Finishing 结束

调用`times`函数分析:

```
//这里定义了 tms 结构体 cpu
#ifdef HAVE_TIMES
	struct tms cpu;
#endif

//这里对函数进行了调用
#ifdef HAVE_TIMES
	times(&cpu);
#endif

//这里进行了赋值, 出现在fpm_request_reading_headers中
#ifdef HAVE_TIMES
	proc->cpu_accepted = cpu;
#endif

//timersub 见下面的函数分析
//下面这段代码出现在 fpm_request_end 中

fpm_clock_get(&now); proc->accepted = now 只出现在fpm_request_reading_headers中, 即在请求开始, 可以理解为开始接受请求.
proc->tv = new 出现在每个步骤中
所以 proc->cpu_duration(持续时间) = 从到现在为止的时间(tv, 因为在最后阶段tv=now)减去, 接受请求的时间(accepted, 仅仅在接受请求时候调用 = now)

proc->last_request_cpu.tms_utime = 等于直接结束时候的时间(在fpm_request_end中调用times(&cpu);)减去请求时候记录的时间. 即为一个请求的时间. 也就是上文中我们说的 fpm_request_reading_headers 到 fpm_request_end 的时间.

#ifdef HAVE_TIMES
	timersub(&proc->tv, &proc->accepted, &proc->cpu_duration);
	proc->last_request_cpu.tms_utime = cpu.tms_utime - proc->cpu_accepted.tms_utime;
	proc->last_request_cpu.tms_stime = cpu.tms_stime - proc->cpu_accepted.tms_stime;
	proc->last_request_cpu.tms_cutime = cpu.tms_cutime - proc->cpu_accepted.tms_cutime;
	proc->last_request_cpu.tms_cstime = cpu.tms_cstime - proc->cpu_accepted.tms_cstime;
#endif
```

##### 日志部分

最后写入日志`sapi/fpm/fpm/fpm_log.c`

```
				case 'C': /* %CPU */
					if (format[0] == '\0' || !strcasecmp(format, "total")) {
						if (!test) {
							tms_total = proc.last_request_cpu.tms_utime + proc.last_request_cpu.tms_stime + proc.last_request_cpu.tms_cutime + proc.last_request_cpu.tms_cstime;
						}
					} else if (!strcasecmp(format, "user")) {
						if (!test) {
							tms_total = proc.last_request_cpu.tms_utime + proc.last_request_cpu.tms_cutime;
						}
					} else if (!strcasecmp(format, "system")) {
						if (!test) {
							tms_total = proc.last_request_cpu.tms_stime + proc.last_request_cpu.tms_cstime;
						}
					} else {
						zlog(ZLOG_WARNING, "only 'total', 'user' or 'system' are allowed as a modifier for %%%c ('%s')", *s, format);
						return -1;
					}

					format[0] = '\0';
					if (!test) {
						len2 = snprintf(b, FPM_LOG_BUFFER - len, "%.2f", tms_total / fpm_scoreboard_get_tick() / (proc.cpu_duration.tv_sec + proc.cpu_duration.tv_usec / 1000000.) * 100.);
					}
break;
```

我们在最后这个日志文件中可见`%C`里面包含3个分支:

* user
* system
* total

对应的是各个请求时间占比.

### 几个函数分析

#### `timersub`

```

subtracts the time value in b from the time value in a, and places the result in the timeval pointed to by res. The result is normalized such that res->tv_usec has a value in the range 0 to 999,999. 

void timersub(struct timeval *a, struct timeval *b,
              struct timeval *res);
```


