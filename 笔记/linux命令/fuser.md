# fuser

## 概述

`fuser`命令用于报告进程使用的文件和网络套接字.

`fuser`命令列出了本地进程的进程号, 那些本地进程使用file, 参数指定的本地或远程文件. 对于阻塞特别设备, 此命令列出了使用该设备上任何文件的进程.

每个进程号后面都跟随一个字母, 该字母指示进程如何使用文件.

* `c`：指示进程的工作目录.
* `e`：指示该文件为进程的可执行文件(即进程由该文件拉起).
* `f`：指示该文件被进程打开, 默认情况下f字符不显示.
* `F`：指示该文件被进程打开进行写入, 默认情况下F字符不显示.
* `r`：指示该目录为进程的根目录.
* `m`：指示进程使用该文件进行内存映射, 抑或该文件为共享库文件, 被进程映射进内存.

## 语法

```
fuser (选项) (参数)
```

## 选项

* `-a`：显示命令行中指定的所有文件.
* `-k`：杀死访问指定文件的所有进程.
* `-i`：杀死进程前需要用户进行确认.
* `-l`：列出所有已知信号名.
* `-m`：指定一个被加载的文件系统或一个被加载的块设备.
* `-n`：选择不同的名称空间.
* `-u`：在每个进程后显示所属的用户名.

## 示例

我们以根目录挂载的分区作为示例.

```
[ansible@demo ~]$ df -h
Filesystem             Size  Used Avail Use% Mounted on
/dev/vda1               99G  2.0G   92G   3% /

[ansible@demo ~]$ sudo fuser -um /dev/vda1
/dev/vda1:               1rce(root)     2rc(root)     3rc(root)     5rc(root)     7rc(root)     8rc(root)     9rc(root)    10rc(root)    11rc(root)    12rc(root)    13rc(root)    15rc(root)    18rc(root)    19rc(root)    20rc(root)    21rc(root)    22rc(root)    23rc(root)    24rc(root)    31rc(root)    32rc(root)    33rc(root)    34rc(root)    42rc(root)    44rc(root)    45rc(root)    47rc(root)    66rc(root)    98rc(root)   235rc(root)   238rc(root)   239rc(root)   240rc(root)   241rc(root)   248rc(root)   260rc(root)   261rc(root)   266rc(root)   267rc(root)   335rce(root)   358rce(root)   445rc(root)   466rce(polkitd)   468rce(root)   470rce(root)   471rce(dbus)   483rce(root)   484rce(root)   499rce(root)   722rce(root)   810rce(root)   831rce(ntp)  7997rc(root)  9143rc(root)  9293rc(root)  9384rc(root)  9438rc(root)  9476rce(root)  9478rce(ansible)  9577rc(root)  9584rce(root)  9586rce(ansible)  9587rce(ansible)  9622rc(root)  9651rce(root) 10501rce(root) 16256rce(root) 18855rce(root) 18862rce(root) 19329rce(root) 20770rc(root) 20785rc(root) 20786rc(root) 20787rc(root) 20789rc(root) 20795rc(root) 20796rc(root) 20797rc(root) 20798rc(root) 20799rc(root) 20800rc(root) 20802rce(root) 21602rc(root) 21603rc(root) 21792rc(root) 21793rc(root) 21794rc(root) 21797rc(root) 21801rc(root) 21802rc(root) 21803rc(root) 21804rc(root) 21805rc(root) 21806rc(root) 21807rc(root) 23876rce(root) 23937rce(root) 23941rce(root) 23961rce(root) 24843rce(root) 24848rce(root) 25551rce(root)
```