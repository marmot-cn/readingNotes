# uname

---

显示当前操作系统名称,

## 语法

`uname(选项)`

## 选项

* `-a`或`--all`: 显示全部的信息
* `-m`或`--machine`: 机器硬件(CPU)名
* `-n`或`--kernel-name`: 显示在网络上的主机名称
* `-r`或`--kernel-release`: 显示操作系统的发行编号
* `-s`或`--sysname`: 显示操作系统名称
* `-v`或`--kernel-version`: 显示操作系统(kernel)的版本
* `-p`或`--processor`: 输出处理器类型或"unknown"
* `-i`或`--hardware-platform`: 输出硬件平台或"unknown"
* `-o`或`--operating-system`: 输出操作系统名称

## 示例

```shell
[root@iZ94xwu3is8Z ~]# uname
Linux
[root@iZ94xwu3is8Z ~]# uname -s
Linux
[root@iZ94xwu3is8Z ~]# uname -m
x86_64
[root@iZ94xwu3is8Z ~]# uname -i
x86_64
[root@iZ94xwu3is8Z ~]# uname -n
iZ94xwu3is8Z
[root@iZ94xwu3is8Z ~]# uname -o
GNU/Linux
[root@iZ94xwu3is8Z ~]# uname -r
3.10.0-123.9.3.el7.x86_64
[root@iZ94xwu3is8Z ~]# uname -v
#1 SMP Thu Nov 6 15:06:03 UTC 2014
[root@iZ94xwu3is8Z ~]# uname -a
Linux iZ94xwu3is8Z 3.10.0-123.9.3.el7.x86_64 #1 SMP Thu Nov 6 15:06:03 UTC 2014 x86_64 x86_64 x86_64 GNU/Linux
``