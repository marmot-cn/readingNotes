# NOTE

---

### `cli_set_process_title`

设置进程名称,可以在 top, ps 中看见.

`cli_get_process_title` 获取当前进程 `cli_set_process_title()` 设置后的名称.

**示例**

```php
<?php
$title = "My Amazing PHP Script";
$pid = getmypid(); // you can use this to see your process title in ps

if (!cli_set_process_title($title)) {
    echo "Unable to set process title for PID $pid...\n";
    exit(1);
} else {
    echo "The process title '$title' for PID $pid has been set for your process!\n";
    sleep(20);
}

root@bb0aa005831b:/var/www/html# php test2.php
The process title 'My Amazing PHP Script' for PID 94 has been set for your process!

另外一个终端:
root@bb0aa005831b:/var/www/html# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 Jun26 ?        00:00:07 php-fpm: master process (/usr/local/etc/php-fpm.conf)
www-data     7     1  0 Jun26 ?        00:00:00 php-fpm: pool www
www-data     8     1  0 Jun26 ?        00:00:00 php-fpm: pool www
root        68     0  0 06:40 ?        00:00:00 /bin/bash
root        84     0  0 06:58 ?        00:00:00 /bin/bash
root        94    68  0 07:02 ?        00:00:00 My Amazing PHP Script
root        95    84  0 07:02 ?        00:00:00 ps -ef

可以看见 进程 CMD 一栏中的名字已经修改了.
```

### `spl_object_hash`

```php
string spl_object_hash ( object $obj );
```

本函数为指定对象返回一个唯一标识符.这个标识符可用于作为保存对象或区分不同对象的`hash key`.

* object: Any object

**示例**

比较两个不同对象.

```php
?php

class A {

    public $a = 0;
}

$a = new A();
$b = new A();

var_dump(spl_object_hash($a));
var_dump(spl_object_hash($b));

输出 (结果不同):
string(32) "00000000120daf1e0000000042af9d96"
string(32) "00000000120daf1d0000000042af9d96"
```

比较两个相同对象.

```php
<?php

class A {

    public $a = 0;
}

$a = new A();

var_dump(spl_object_hash($a));
var_dump(spl_object_hash($a));

输出 (结果相同):
root@bb0aa005831b:/var/www/html# php test.php
string(32) "00000000103a223300000000703eec5d"
string(32) "00000000103a223300000000703eec5d"
```

### posix_kill

发送一个信号给一个进程

```php
bool posix_kill ( int $pid , int $sig )
```

* pid: 进程id
* sig: 信号
