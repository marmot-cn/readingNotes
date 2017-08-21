# killall

---

用于杀死指定名字的进程.

我们可以使用`kill`命令杀死指定进程`PID`的进程,如果要找到我们需要杀死的进程,我们还需要在之前使用`ps`等命令再配合`grep`来查找进程,而`killall`把这两个过程合二为一.

`root`用户将影响用户的进程, 非`root`用户只能影响自己的进程.

```shell
把所有该用户登录的bash杀死.
[ansible@localhost ~]$ killall -9 bash
```

## cenots7 精简版安装killall命令

```shell
yum install psmisc
```

`psmisc`:

* `killall`
* `fuser`: 显示使用指定文件或者文件系统的进程的PID.
* `pstree`: 树型显示当前运行的进程.