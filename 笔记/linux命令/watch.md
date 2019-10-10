# watch

---

## 功能

可以将命令的输出结果输出到标准输出设备, 多用于周期性执行命令/定时执行命令

## 格式

`watch[参数][命令]`

## 参数

* `-n`or`--interval` watch缺省每2秒运行一下程序,可以用`-n`或`-interval`来指定间隔的时间.
* `-d`或`--differences` 用`-d`或`--differences`选项`watch`会高亮显示变化的区域. 而`-d=cumulative`选项会把变动过的地方(不管最近的那次有没有变动)都高亮显示出来.
* `-t`或`-no-title` 会关闭watch命令在顶部的时间间隔和命令.

## 示例

### 每隔一秒高亮显示网络链接数的变化情况

`watch -n 1 -d netstat -ant`

### 每隔一秒高亮显示http链接数的变化情况

`watch -n 1 -d 'pstree|grep http'`

### 实时查看模拟攻击客户机建立起来的连接数

`watch 'netstat -an | grep:21 | \ grep<模拟攻击客户机的IP>| wc -l'`

### 监测当前目录中 scf 的文件的变化

`watch -d 'ls -l|grep scf'`

### 10秒一次输出系统的平均负载

`watch -n 10 'cat /proc/loadavg'`