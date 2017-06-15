# PHP PCNTL 函数

#### pcntl_alarm

为进程设置一个 `alarm` 闹钟信号.

```php
int pcntl_alarm ( int $seconds )
```

创建一个计时器,在指定的秒数后向进程发送一个`SIGALRM`信号.每次对 `pcntl_alarm()`的调用都会取消之前设置的alarm信号.

**示例**

```php
<?php

function sig_handler($signo)
{
    switch($signo) {
        case SIGALRM:
            echo "SIGALRM", PHP_EOL;
            break;
        default:
            return false;
    }
}

pcntl_signal(SIGALRM, "sig_handler");

pcntl_alarm (5);

while(true)
{
    pcntl_signal_dispatch();
    sleep(1);
}
```

输出:

```
root@543856bfc00a:/var/www/html# php test.php
SIGALRM (运行 5秒 后)
```

#### pcntl_errno

等同于 `pcntl_get_last_error` 返回错误号码在上一次 `pcntl` 函数失败后.

#### pcntl_exec

在==当前进程空间==执行指定程序.

```php
void pcntl_exec ( string $path [, array $args [, array $envs ]] )
```

* path: 必须时可执行二进制文件路径或一个在文件第一行指定了 一个可执行文件路径标头的脚本（比如文件第一行是#!/usr/local/bin/perl的perl脚本).
* args: `args`是一个要传递给程序的参数的字符串数组.
* envs: `envs`是一个要传递给程序作为环境变量的字符串数组.这个数组是 key => value格式的,key代表要传递的环境变量的名称,value代表该环境变量值.

**示例**

`backgroun.php` 背景程序,以这个程序运行 `command` 程序.

```php
<?php

pcntl_exec('./command', array('arg1','arg2'), array('env'=>'test'));

echo 222, PHP_EOL;
```

`command` 需要运行的程序. 我们赋予其 `+x` 执行权限.

```php
#!/usr/local/bin/php

<?php

echo "argv", PHP_EOL;
var_dump($argv);

echo "ENV", PHP_EOL;
var_dump($_ENV);
```

运行 `background.php`

```php
php background.php

argv
array(3) {
  [0]=>
  string(9) "./command"
  [1]=>
  string(4) "arg1"
  [2]=>
  string(4) "arg2"
}
ENV
array(1) {
  ["env"]=>
  string(4) "test"
}
```

可以看见 `background.php` 后面的输出并没有输出. `command`程序会开始执行,执行结束后退出. 同时测试了 传参 和 环境变量.

#### pcntl_fork

在当前进程当前位置产生分支(子进程).

`fork`是创建了一个子进程,父进程和子进程都从fork的位置开始向下继续执行,不同的是父进程执行过程中,得到的==fork返回值为子进程号==,而==子进程得到的是0==.

```php
int pcntl_fork ( void )
```

#### pcntl_setpriority

修改任意进程的优先级

```php
bool pcntl_setpriority ( int $priority [, int $pid = getmypid() [, int $process_identifier = PRIO_PROCESS ]] )
```

设置进程号为 `pid` 的进程的优先级.

* priority: priority 通常时-20至20这个范围内的值.默认优先级是0,值越小代表 优先级越高
* pid: 如果没有指定,默认是当前进程的进程号
* process_identifier: 
	* PRIO_PGRP: 获取进程组优先级
	* PRIO_USER: 获取用户进程优先级
	* PRIO_PROCESS: 默认值;获取进程优先级
	
#### pcntl_getpriority

获取任意进程的优先级.

```php
int pcntl_getpriority ([ int $pid = getmypid() [, int $process_identifier = PRIO_PROCESS ]] )
```

获取进程号为 `pid` 的进程的优先级.值越小代表优先级越高.(nice 值)

* pid: 如果没有指定,默认是当前进程的进程号.
* process_identifier: 
	* PRIO_PGRP: 获取进程组优先级,前方`id`是进程组id.
	* PRIO_USER: 获取用户进程优先级,前方`id`是用户id.
	* PRIO_PROCESS: 默认值;获取进程优先级

`process_identifier` 控制着`pid`该如何解释.

**查看源码**

查看源码:

```c
/usr/src/php/ext/pcntl# pwd
/usr/src/php/ext/pcntl# cat pcntl.c

PHP_FUNCTION(pcntl_getpriority)
{
	zend_long who = PRIO_PROCESS;
	zend_long pid = getpid();
	int pri;

	if (zend_parse_parameters(ZEND_NUM_ARGS(), "|ll", &pid, &who) == FAILURE) {
		RETURN_FALSE;
	}

	/* needs to be cleared, since any returned value is valid */
	errno = 0;

	pri = getpriority(who, pid);

	if (errno) {
		PCNTL_G(last_error) = errno;
		switch (errno) {
			case ESRCH:
				php_error_docref(NULL, E_WARNING, "Error %d: No process was located using the given parameters", errno);
				break;
			case EINVAL:
				php_error_docref(NULL, E_WARNING, "Error %d: Invalid identifier flag", errno);
				break;
			default:
				php_error_docref(NULL, E_WARNING, "Unknown error %d has occurred", errno);
				break;
		}
		RETURN_FALSE;
	}

	RETURN_LONG(pri);
}
```

我们可见 默认是:

* PRIO_PROCESS
* getpid()

**示例代码**

```php
<?php

//pcntl_setpriority(-1);

echo 'current process:';
$result = pcntl_getpriority();
echo $result, PHP_EOL;

echo 'user process :';
$result = pcntl_getpriority(0, PRIO_USER);
echo $result, PHP_EOL;

while(true)
{
    sleep(1);
}
```

输出:

```
current process:-1
user process :-1
```

分开解释这些值,首先我们在运行程序的时候,开启另外一个容器:

```shell
ps -efl
F S UID        PID  PPID  C PRI  NI ADDR SZ WCHAN  STIME TTY          TIME CMD
4 S root         1     0  0  80   0 - 62499 -      09:58 ?        00:00:00 php-fpm: master process (/usr/local/etc/php-fpm.conf)
5 S www-data     7     1  0  80   0 - 62499 -      09:58 ?        00:00:00 php-fpm: pool www
5 S www-data     8     1  0  80   0 - 62499 -      09:58 ?        00:00:00 php-fpm: pool www
4 S root         9     0  0  80   0 -  5069 -      09:58 ?        00:00:00 /bin/bash
4 S root        17     9  0  79  -1 - 53332 -      10:00 ?        00:00:00 php priority.php
4 S root        18     0  0  80   0 -  5069 -      10:00 ?        00:00:00 /bin/bash
0 R root        27    18  0  80   0 -  4375 -      10:01 ?        00:00:00 ps -efl
```

我们还需要开启容器的权限`privileged: true`.

```php
root@389444860b8e:/var/www/html# nice -n -1 php priority.php
current process:-1
user process :-1
```

我们通过 `nice` 命令运行程序, 并设置 `nice` 值为 `-1` .

第一个 `-1` 表示当前进程的 `nice` 为 `-1`.

第二个 `-1` 表示的是当前用户的所有进程 `nice` 值最低的是 `-1`.

The getpriority() call returns the highest priority (lowest numerical
value) enjoyed by any of the specified processes.  The setpriority()
call sets the priorities of all of the specified processes to the specified value.

`getpriority()` 源码返回最高优先级(最低 nice值), 所以可见 user 返回 -1,因为是最低,如果有更低的,则返回更低的.

#### pcntl_signal_dispatch 

#### pcntl_signal_get_handler 

#### pcntl_signal

#### pcntl_sigprocmask

#### pcntl_sigtimedwait

#### pcntl_sigwaitinfo

#### pcntl_strerror 

#### pcntl_wait

#### pcntl_waitpid

#### pcntl_wexitstatus

#### pcntl_wifexited

#### pcntl_wifsignaled

#### pcntl_wifstopped

#### pcntl_wstopsig

#### pcntl_wtermsig 



 