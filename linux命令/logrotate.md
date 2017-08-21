# logrotate

---

可以自动对日志进行截断(或轮循), 压缩以及删除旧的日志文件.

## 安装

一般都会安装在内.

```shell
yum install logrotate crontabs 
```

## 配置文件

`logrotate`的配置文件是`/etc/logrotate.conf`. 日志文件的轮循设置在独立的配置文件中, 它(们)放在`/etc/logrotate.d/`目录下.

## 参数

* `minisize 1M`: 必须大于1MB才会轮转.
* `size 50M`: 超过50MB后轮转日志.
* `copytruncate`: 清空原有文件, 而不是创建一个新文件.
* `monthly`: 日志文件将按月轮循. 其它可用值为‘`daily`’，‘`weekly`’或者‘`yearly`’.
* `rotate 5`:  一次将存储5个归档日志. 对于第六个归档, **时间最久的归档将被删除**.
* `compress`: 在轮循任务完成后, 已轮循的归档将使用`gzip`进行压缩.
* `delaycompress`: 总是与`compress`选项一起用,`delaycompress`选项指示`logrotate`不要将最近的归档压缩,压缩将在下一次轮循周期进行.这在你或任何软件仍然需要读取最新归档时很有用.
* `missingok`: 在日志轮循期间,任何错误将被忽略,例如“文件无法找到”之类的错误.
* `notifempty`: 如果日志文件为空,轮循不会进行. 
* `create mode owner group `:  以指定的权限创建全新的日志文件,同时`logrotate`也会重命名原始日志文件.
* `prerotate,endscript`: 在logrotate之前执行的命令. 如`/usr/bin/charrt -a /var/log/logfile`.
* `postrotate,endscript`: 在所有其它指令完成后, `postrotate`和`endscript`里面指定的命令将被执行.在这种情况下, `rsyslogd`进程将立即再次读取其配置并继续运行.如`/usr/bin/charrt +a /var/log/logfile`.
* `sharedscripts`: 共享脚本，表示切换时只执行一次脚本.
* `dateext`: 增加日期作为后缀，不然会是一串无意义的数字.
* `dateformat .%s`: 切换后文件名，必须配合dateext使用.

## 手动运行`logrotate`

要调用为`/etc/lograte.d/`下配置的所有日志调用`logrotate`：

```shell
logrotate /etc/logrotate.conf 
```

要为某个特定的配置调用`logrotate`:

```shell
logrotate /etc/logrotate.d/log-file 
```

## 计划任务

```shell
[ansible@iZ944l0t308Z ~]$ sudo cat /etc/cron.daily/logrotate
#!/bin/sh

/usr/sbin/logrotate -s /var/lib/logrotate/logrotate.status /etc/logrotate.conf
EXITVALUE=$?
if [ $EXITVALUE != 0 ]; then
    /usr/bin/logger -t logrotate "ALERT exited abnormally with [$EXITVALUE]"
fi
exit 0
```

每天执行一次.

## 示例

```shell
[root@localhost ansible]# cat /etc/logrotate.d/syslog
/var/log/cron
/var/log/maillog
/var/log/messages
/var/log/secure
/var/log/spooler
{
    missingok
    sharedscripts #–表示切换时脚本只执行一次
    postrotate #–表示rotate后执行的脚本
	/bin/kill -HUP `cat /var/run/syslogd.pid 2> /dev/null` 2> /dev/null || true # 其中“-HUP”使syslogd关闭所有日志文件，重读/etc/syslog.conf配置文件后重新开始记录日志
    endscript  #–表示脚本结束
}
```

## 测试

```shel
[root@localhost ansible]# touch /var/log/log-file
[root@localhost ansible]# echo 111 > /var/log/log-file 
[root@localhost ansible]# cat /etc/logrotate.d/log-file
/var/log/log-file {
	monthly
	rotate 5
	dateext
	create 644 root root
	postrotate
	/usr/bin/killall -HUP rsyslogd
	endscript
}

排障过程中的最佳选择是使用‘-d’选项以预演方式运行logrotate。要进行验证，不用实际轮循任何日志文件，可以模拟演练日志轮循并显示其输出
[root@localhost ansible]# logrotate -d /etc/logrotate.d/log-file
reading config file /etc/logrotate.d/log-file
Allocating hash table for state file, size 15360 B

Handling 1 logs

rotating pattern: /var/log/log-file  monthly (5 rotations)
empty log files are rotated, old logs are removed
considering log /var/log/log-file
  log does not need rotating (log has been already rotated)
  
正如我们从上面的输出结果可以看到的，logrotate判断该轮循是不必要的。如果文件的时间小于一天，这就会发生了

我们也可以通过使用‘-f’选项来强制logrotate轮循日志文件,‘-v’参数提供了详细的输出
[root@localhost log]# logrotate -vf /etc/logrotate.d/log-file
reading config file /etc/logrotate.d/log-file
Allocating hash table for state file, size 15360 B

Handling 1 logs

rotating pattern: /var/log/log-file  forced from command line (5 rotations)
empty log files are rotated, old logs are removed
considering log /var/log/log-file
  log needs rotating
rotating log /var/log/log-file, log->rotateCount is 5
dateext suffix '-20170820'
glob pattern '-[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]'
glob finding old rotated logs failed
renaming /var/log/log-file to /var/log/log-file-20170820
creating new /var/log/log-file mode = 0644 uid = 0 gid = 0
running postrotate script

[root@localhost log]# ls /var/log/log-file*
/var/log/log-file  /var/log/log-file-20170820
[root@localhost log]# cat /var/log/log-file-20170820
111
```
