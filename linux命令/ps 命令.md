### ps 命令

ps 命令是用于查看系统进程信息.

默认不带参数的,只显示运行在当前控制台下的属于当前用户的进程.

#### 常用参数

**`-A`** 

显示所有进程.

**`-e`** 

显示所有进程.

**`-a`** 

显示除守护进程和无终端进程外的所有进程.

**`-d`**

显示所有进程,除会话首进程(session leaders).

```shell
root@389444860b8e:/var/www/html# ps -d
  PID TTY          TIME CMD
    7 ?        00:00:00 php-fpm
    8 ?        00:00:00 php-fpm
   17 ?        00:00:00 php
   34 ?        00:00:00 ps
root@389444860b8e:/var/www/html# ps -A
  PID TTY          TIME CMD
    1 ?        00:00:00 php-fpm
    7 ?        00:00:00 php-fpm
    8 ?        00:00:00 php-fpm
    9 ?        00:00:00 bash
   17 ?        00:00:00 php
   18 ?        00:00:00 bash
   35 ?        00:00:00 ps
```

其中区别的是 1 号, 9 号, 18 号 没显示.

因为这些进程都是会话首进程.

```shell
 PGID   PID CMD
    1     1 php-fpm: master process (/usr/local/etc/php-fpm.conf)
    1     7 php-fpm: pool www
    1     8 php-fpm: pool www
    9     9 /bin/bash
   17    17 php priority.php
   18    18 /bin/bash
   44    44 ps -A -o pgid,pid,cmd
```

1,9,18 都是会话首进程.

**`-N`**

显示所有的程序,除了执行ps指令终端机下的程序之外.

**`-p pidlist`**

显示pid在pidlist里的进程,pidlist表示pid数组,多个用逗号分隔,如-p 123,124.

```shell
root@389444860b8e:/var/www/html# ps -p 1,7
  PID TTY          TIME CMD
    1 ?        00:00:00 php-fpm
    7 ?        00:00:00 php-fpm
```

**`-C`**

通过执行命令名称查找进程,显示执行命令在cmdlist里的进程.

```shell
root@389444860b8e:/var/www/html# ps -C php-fpm -o pgid,pid,cmd
 PGID   PID CMD
    1     1 php-fpm: master process (/usr/local/etc/php-fpm.conf)
    1     7 php-fpm: pool www
    1     8 php-fpm: pool www
```

**`-G grplist`**

通过进程组ID或进程组名称查找进程,显示进程组ID在grplist里的进程.

