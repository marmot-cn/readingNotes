# Linux 后台进程 和 Ctrl+c

---

今天测试信号 `19(SIGSTOP)` 和 `18(SIGCONT)`.

写了一个简单的脚本:

```php
<?php

while(true)
{
    sleep(1);
    echo 1, PHP_EOL;
}
```

对该脚本执行的进程发送 `kill -19 xxx`, 正常放到后台.

对该脚本执行的进程发送 `kill -18 xxx`, 正常输出,但是不能终止 `ctrl+c` 不起作用.

突然想到:

`Ctrl+c`: 终止一个正在==前台==运行的进程.

`bg`(等同于 18 信号): 重新启动一个挂起(暂停)的作业,并在==后台==运行.

参见 `APUE` 笔记 `linux_作业控制_和_nohup.md`.