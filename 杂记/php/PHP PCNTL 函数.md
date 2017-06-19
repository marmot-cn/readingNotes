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

我们在源码内看见使用的是`pri = getpriority(who, pid);`,其实我个人感觉这里写的很容易让人误解,因为真正底层头文件是`int getpriority(int which,int who);`,第一个参数是`which`,第二个参数是`who`.而 PHP 源码里面第一个参数是 `who`,第二个是 `pid` 很容易让人误解.

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

#### pcntl_signal

安装一个信号处理器.

```php
bool pcntl_signal ( int $signo , callback $handler [, bool $restart_syscalls = true ] )
```

函数 `pcntl_signal()` 为 `signo` 指定的信号安装一个新的信号处理器.

* signo: 信号编号
* handler: 
	* 信号处理器可以是用户创建的函数或方法名字
	* 系统常量 `SIG_IGN`(忽略信号处理程序)
	* `SIG_DFL`(默认信号处理程序)
* restart_syscalls: 系统调用被信号打断时,系统调用是否从开始处重新开始. 

#### pcntl_signal_dispatch

调用等待信号处理器.

```php
bool pcntl_signal_dispatch ( void )
```

`pcntl_signal_dispatch()` 调用每个等待信号通过 `pcntl_signal()` 安装的处理器.

使用 `pcntl_signal()` 设置信号处理器后,调用该函数触发信号处理器.

#### pcntl_signal_get_handler

获取指定信号的处理器. 这个函数只有在`7.1`可用.

`int|string pcntl_signal_get_handler ( int $signo )`

* signo: 信号编号
* 返回:
	* 整数: SIG_DFL(0), SIG_IGN(1)
	* 字符: 函数名字

```php
<?php
// Outputs: int(0)
var_dump(pcntl_signal_get_handler(SIGUSR1)); 

function pcntl_test($signo) {}
pcntl_signal(SIGUSR1, 'pcntl_test');
// Outputs: string(10) "pcntl_test"
var_dump(pcntl_signal_get_handler(SIGUSR1)); 

pcntl_signal(SIGUSR1, SIG_DFL);
// Outputs: int(0)
var_dump(pcntl_signal_get_handler(SIGUSR1)); 

pcntl_signal(SIGUSR1, SIG_IGN);
// Outputs: int(1)
var_dump(pcntl_signal_get_handler(SIGUSR1)); 
```

#### pcntl_sigprocmask

设置或检索阻塞信号.

```php
bool pcntl_sigprocmask ( int $how , array $set [, array &$oldset ] )
```

`pcntl_sigprocmask()`用来增加,删除或设置阻塞信号,具体行为依赖于参数 how.

* how 设置pcntl_sigprocmask()函数的行为
	* SIG_BLOCK: 把信号加入到当前阻塞信号中
	* SIG_UNBLOCK: 从当前阻塞信号中移出信号
	* SIG_SETMASK: 用给定的信号列表替换当前阻塞信号列表
* set: 信号列表
* oldset: 是一个输出参数,用来返回之前的阻塞信号列表数组.

#### pcntl_sigwaitinfo

等待信号.

```php
int pcntl_sigwaitinfo ( array $set [, array &$siginfo ] )
```

函数暂停调用脚本的执行直到接收到 `set` 参数中列出的某个信号.只要其中的一个信号已经在等待状态(通过 `pcntl_sigprocmask()`函数阻塞),函数pcntl_sigwaitinfo()就回立刻返回.

* set : 要等待的信号数组.
* siginfo: 是一个输出参数,用来返回信号的信息.
	* 以下元素会为所有信号设置:
		*  signo: 信号编号
		*  errno: 错误编号
		*  code: 信号代码
	* 下面元素可能会为 `SIGCHLD`(在一个进程终止或者停止时,将SIGCHLD信号发送给其父进程)信号设置
		* status: 退出的值或信号 
		* utime: 用户小号的时间
		* stime: 系统(内核)小号的时间
		* pid: 发送进程id
		* uid: 发送进程的实际用户id
	* 信号`SIGILL`, `SIGFPE`, `SIGSEGV` 和 `SIGBUS` 可能会被设置的元素
		* addr: 发生故障的内存位置 
	* 信号`SIGPOLL` 可能会被设置的元素.
		* band: Band event
		* fd: 文件描述符 

**一些信号解释**

* SIGFPE: 在发生致命的算术运算错误时发出. 不仅包括浮点运算错误, 还包括溢出及除数为0等其它所有的算术的错误.
* SIGILL: 执行了非法指令. 通常是因为可执行文件本身出现错误, 或者试图执行数据段. 堆栈溢出时也有可能产生这个信号.
* SIGBUS: 非法地址,包括内存地址对齐(alignment)出错. 比如访问一个四个字长的整数, 但其地址不是4的倍数.它与`SIGSEGV`的区别在于后者是由于对合法存储地址的非法访问触发的(如访问不属于自己存储空间或只读存储空间).
* SIGSEGV: 试图访问未分配给自己的内存, 或试图往没有写权限的内存地址写数据.

**示例**

```php
<?php

echo "Blocking SIGHUP signal\n";
pcntl_sigprocmask(SIG_BLOCK, array(SIGHUP));

echo "Waiting for signals\n";
$info = array();
pcntl_sigwaitinfo(array(SIGHUP), $info);

var_dump($info);
```

输出:

```php
root@389444860b8e:/var/www/html# php test.php
Blocking SIGHUP signal
Waiting for signals

我们开启另外一个终端

root@389444860b8e:/var/www/html# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 Jun18 ?        00:00:00 php-fpm: master process (/usr/local/etc/php-fpm.conf)
www-data     7     1  0 Jun18 ?        00:00:00 php-fpm: pool www
www-data     8     1  0 Jun18 ?        00:00:00 php-fpm: pool www
root         9     0  0 01:35 ?        00:00:00 /bin/bash
root        16     0  0 01:35 ?        00:00:00 /bin/bash
root        22     9  0 01:35 ?        00:00:00 php test.php
root        23    16  0 01:35 ?        00:00:00 ps -ef
root@389444860b8e:/var/www/html# kill -1 22

返回第一个终端输出:
root@389444860b8e:/var/www/html# php test.php
Blocking SIGHUP signal
Waiting for signals
array(3) {
  ["signo"]=>
  int(1)
  ["errno"]=>
  int(0)
  ["code"]=>
  int(0)
}
```

**示例子进程**

```php
<?php

echo "Blocking SIGHUP signal\n";
pcntl_sigprocmask(SIG_BLOCK, array(SIGHUP,SIGCHLD));

$pid = pcntl_fork();

if ($pid) {
    echo "Waiting for signals\n";
    $info = array();
    pcntl_sigwaitinfo(array(SIGHUP,SIGCHLD), $info);
} else {
    for ($i=0; $i<5; $i++) {
        sleep(1);
    }
    exit();
}
var_dump($info);
```

输出:

```php
root@389444860b8e:/var/www/html# php test.php
Blocking SIGHUP signal
Waiting for signals
等待 5 秒
array(8) {
  ["signo"]=>
  int(17)
  ["errno"]=>
  int(0)
  ["code"]=>
  int(1)
  ["status"]=>
  int(0)
  ["utime"]=>
  float(0)
  ["stime"]=>
  float(1)
  ["pid"]=>
  int(45)
  ["uid"]=>
  int(0)
}

我们可见返回了 pid 进程id(子进程id)
```

#### pcntl_sigtimedwait

带超时机制的信号等待.

```
int pcntl_sigtimedwait ( array $set [, array &$siginfo [, int $seconds = 0 [, int $nanoseconds = 0 ]]] )
```

和 `pcntl_sigwaitinfo()` 的行为一致, 不同在于多了两个参数 `seconds` 和 `nanoseconds` 这使得脚本等待的事件有了一个时间的上限.

#### pcntl_strerror 

获取给定 `errno` 的错误信息.

```php
string pcntl_strerror ( int $errno )
```

`pcntl_get_last_error()`, 该函数会返回上次 `pcntl` 函数失败后的错误号.

**示例**

```php
<?php

var_dump(pcntl_strerror(0));
var_dump(pcntl_strerror(1));
var_dump(pcntl_strerror(2));
var_dump(pcntl_strerror(3));
```

输出:

```php
root@389444860b8e:/var/www/html# php test2.php
string(7) "Success"
string(23) "Operation not permitted"
string(25) "No such file or directory"
string(15) "No such process"
```

#### pcntl_wait

等待或返回 fork 的子进程状态.

```php
int pcntl_wait ( int &$status [, int $options = 0 ] )
```

`wait` 函数挂起当前进程的执行直到一个子进程退出或接收到一个信号要求中断当前进程或调用一个信号处理函数.如果一个子进程在调用次函数时已经退出(俗称僵尸进程),此函数立刻返回.进程使用的所有系统资源将被释放.

等同于 `-1` 作为参数 `pid` 的值并且没有 `options` 参数来调用的 `pcntl_waitpid()` 函数.

* status: `pcntl_wait()` 将会存储状态信息到 `status` 参数上,这个通过`status`参数返回的状态信息可以用以下函数`pcntl_wifexited()`, `pcntl_wifstopped()`, `pcntl_wifsignaled()`, `pcntl_wexitstatus()`, `pcntl_wtermsig()`以及 `pcntl_wstopsig()`获取其具体的值.
* options: 
	* 0
	* WNOHANG 如果没有子进程退出立刻返回.
	* WUNTRACED 子进程已经退出并且其状态未报告时返回
	* WNOHANG | WUNTRACED (或运算,两个常量代表意义都有效).

#### pcntl_waitpid

等待或返回fork的子进程状态

```php
int pcntl_waitpid ( int $pid , int &$status [, int $options = 0 ] )
```

挂起当前进程的执行直到参数pid指定的进程号的进程退出， 或接收到一个信号要求中断当前进程或调用一个信号处理函数.

和 `pcntl_wait` 区别是制定了进程号.

* pid:
	* < -1: 等待任意==进程组ID==等于参数 pid 给定值得绝对值的进程 
	* -1: 等待任意子进程,与`pcntl_wait`函数行为一致
	* 0: 等待==任意与调用进程组ID相同==的子进程.
	* >0: 等待进程号等于参数 `pid` 值得子进程

#### pcntl_wifexited

检查状态代码是否代表一个正常的退出.

```php
bool pcntl_wifexited ( int $status )
```

* status: 调用 pcntl_waitpid() 时的状态参数.

当子进程状态代码代表正常退出时返回 TRUE , 其他情况返回 FALSE .

#### pcntl_wexitstatus

返回一个中断的子进程的返回代码.

```php
int pcntl_wexitstatus ( int $status )
```

这个函数仅在函数 `pcntl_wifexited(`)返回 `TRUE` 时有效,

**正常等待退出示例**

```php
$pid = pcntl_fork();

if ($pid == -1) {
    die("could not fork");
} elseif ($pid) {
    echo "I'm the Parent", PHP_EOL;
    pcntl_wait($status);
    echo "Waite done", PHP_EOL;
    $result = pcntl_wifexited($status);
    echo "pcntl_wifexited: ". (boolval($result) ? "true" : "false"), PHP_EOL;
    echo "pcntl_wexitstatus: ".pcntl_wexitstatus($status), PHP_EOL;
} else {
    echo "I'm the children", PHP_EOL;
    sleep(1);
    exit(2);
}
```
输出:

```php
root@5b419151a444:/var/www/html# php test.php
I'm the Parent
I'm the children
1 秒后
Waite done
pcntl_wifexited: true
pcntl_wexitstatus: 2
```

**非正常等待退出示例**

```php
<?php

$pid = pcntl_fork();

if ($pid == -1) {
    die("could not fork");
} elseif ($pid) {
    echo "I'm the Parent", PHP_EOL;
    pcntl_wait($status);
    echo "Waite done", PHP_EOL;
    $result = pcntl_wifexited($status);
    echo "pcntl_wifexited: ". (boolval($result) ? "true" : "false"), PHP_EOL;
} else {
    echo "I'm the children", PHP_EOL;
    sleep(100);
    exit(2);
}
```
`"pcntl_wexitstatus` 只有 `pcntl_wifexited` 返回 `TRUE` 才可以调用.

```php
root@5b419151a444:/var/www/html# php test.php
I'm the Parent
I'm the children

在另外一个终端:
10s内,杀死子进程
root@5b419151a444:/var/www/html# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 05:13 ?        00:00:00 php-fpm: master process (/usr/local/etc/php-fpm.conf)
www-data     7     1  0 05:13 ?        00:00:00 php-fpm: pool www
www-data     8     1  0 05:13 ?        00:00:00 php-fpm: pool www
root         9     0  0 05:13 ?        00:00:00 /bin/bash
root        22     0  0 05:19 ?        00:00:00 /bin/bash
root        69     9  1 05:39 ?        00:00:00 php test.php
root        70    69  0 05:39 ?        00:00:00 php test.php
root        71    22  0 05:40 ?        00:00:00 ps -ef
root@5b419151a444:/var/www/html# kill 70

第一个终端输出:
root@5b419151a444:/var/www/html# php test.php
I'm the Parent
I'm the children
Waite done
pcntl_wifexited: false
```

#### pcntl_wifsignaled

检查子进程是否是由于某个==未捕获==的信号退出的.

```php
bool pcntl_wifsignaled ( int $status )
```

如果子进程是由于某个==未捕获==的信号退出的返回 TRUE, 其他情况返回 FALSE.

#### pcntl_wtermsig 

返回导致子进程中断的信号.

返回导致子进程中断的信号编号. 这个函数仅在 `pcntl_wifsignaled()` 返回 `TRUE` 时有效.

返回整型的信号编号.

**示例**

```php
<?php

$pid = pcntl_fork();

if ($pid == -1) {
    die("could not fork");
} elseif ($pid) {
    echo "I'm the Parent", PHP_EOL;
    pcntl_waitpid($pid, $status);
    echo "Waite done", PHP_EOL;
    $result = pcntl_wifsignaled($status);
    echo "pcntl_wifsignaled: ". (boolval($result) ? "true" : "false"), PHP_EOL;
    echo "pcntl_wtermsig: ". pcntl_wtermsig($status), PHP_EOL;
} else {
    echo "I'm the children", PHP_EOL;
    sleep(100);
    exit(2);
}
```

输出:

```php

另外一个终端:
root@5b419151a444:/var/www/html# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 05:13 ?        00:00:00 php-fpm: master process (/usr/local/etc/php-fpm.conf)
www-data     7     1  0 05:13 ?        00:00:00 php-fpm: pool www
www-data     8     1  0 05:13 ?        00:00:00 php-fpm: pool www
root         9     0  0 05:13 ?        00:00:00 /bin/bash
root        22     0  0 05:19 ?        00:00:00 /bin/bash
root        95     0  0 05:52 ?        00:00:00 /bin/bash
root       125     9  1 06:17 ?        00:00:00 php test.php
root       126   125  0 06:17 ?        00:00:00 php test.php
root       127    22  0 06:17 ?        00:00:00 ps -ef
root@5b419151a444:/var/www/html# kill -9 126

上一个终端:
root@5b419151a444:/var/www/html# php test.php
I'm the children
I'm the Parent
Waite done
pcntl_wifsignaled: true
pcntl_wtermsig: 9
```

#### pcntl_wifstopped

检查子进程当前是否已经==停止==.

```php
bool pcntl_wifstopped ( int $status )
```

仅查子进程当前是否停止; 此函数只有作用于使用了 `WUNTRACED` 作为 `option` 的`pcntl_waitpid()` 函数调用产生的 `status` 时才有效.

如果子进程当前是停止的返回 `TRUE`, 其他情况返回 `FALSE`.

#### pcntl_wstopsig

返回导致子进程停止的信号.

返回导致子进程停止的信号编号,

这个函数仅在 `pcntl_wifstopped()` 返回 TRUE 时有效,

返回信号编号.

**示例**

```php
root@5b419151a444:/var/www/html# cat test.php
<?php

$pid = pcntl_fork();

if ($pid == -1) {
    die("could not fork");
} elseif ($pid) {
    echo "I'm the Parent", PHP_EOL;
    pcntl_waitpid($pid, $status, WUNTRACED);
    echo "Waite done", PHP_EOL;
    $result = pcntl_wifstopped($status);
    echo "pcntl_wifstopped: ". (boolval($result) ? "true" : "false"), PHP_EOL;
    echo "pcntl_wstopsig: ". pcntl_wstopsig($status), PHP_EOL;
} else {
    echo "I'm the children", PHP_EOL;
    sleep(100);
    exit(2);
}
```

输出:

```php
另外一个终端

oot@5b419151a444:/var/www/html# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 05:13 ?        00:00:00 php-fpm: master process (/usr/local/etc/php-fpm.conf)
www-data     7     1  0 05:13 ?        00:00:00 php-fpm: pool www
www-data     8     1  0 05:13 ?        00:00:00 php-fpm: pool www
root         9     0  0 05:13 ?        00:00:00 /bin/bash
root        22     0  0 05:19 ?        00:00:00 /bin/bash
root        95     0  0 05:52 ?        00:00:00 /bin/bash
root       147     9  0 06:37 ?        00:00:00 php test.php
root       148   147  0 06:37 ?        00:00:00 php test.php
root       149    22  0 06:37 ?        00:00:00 ps -ef

19 信号是  SIGSTOP
root@5b419151a444:/var/www/html# kill -19 148

第一个终端返回:

root@5b419151a444:/var/www/html# php test.php
I'm the children
I'm the Parent
Waite done
pcntl_wifstopped: true
pcntl_wstopsig: 19
```

这个示例我们使用了`pcntl_waitpid($pid, $status, WUNTRACED);`. 因为 `pcntl_wifstopped` 只作用于使用了 `WUNTRACED` 作为 `option` 的 `pcntl_waitpid()` 函数调用产生的 `status` 时才有效.

* `WUNTRACED`: 子进程已经退出并且其状态未报告时返回.

查看`C`函数 `WIFSTOPPED`: `WIFSTOPPED(status)`, 当子进程接收到==停止信号时== true.

其实 `pcntl_wifstopped` 该函数也描述的是检查子进程是否已经 ==停止==.

所以该示例我们发了 19(停止) 信号.

 