# tee

---

`tee`指令会从标准输入设备读取数据,将其内容输出到标准输出设备,同时保存成文件(注意是即输出到STDOUT, 有保存成文件).

**参数**

* `-a`或`--append` 附加到既有文件的后面, 而非覆盖它.
* `-i-i`或`--ignore-interrupts` 忽略中断信号.

从标注输入

```shell
root@iZ944l0t308Z ~]# tee t222
111
111
333
333
^C
[root@iZ944l0t308Z ~]# cat t222
111
333
[root@iZ944l0t308Z ~]#
```

从管道

```shel
[root@iZ944l0t308Z ~]# who | tee who.out
root     pts/0        Jul 19 14:30 (219.144.242.44)
[root@iZ944l0t308Z ~]# cat who.out
root     pts/0        Jul 19 14:30 (219.144.242.44)
```

```
[root@iZ944l0t308Z ~]# echo 222 > ttt
[root@iZ944l0t308Z ~]# cat ttt | tee test
222
[root@iZ944l0t308Z ~]# cat test
222
```

#### vim

`vim`里面使用`tee`来执行`sudo`保存到当前没有编辑权限的文件.

`:w !sudo tee % > /dev/null`

**`%`**

`%` 代表当前文件.

**`:w`**

如果打开`file1.txt`编辑后执行`:w file2.txt`它会类似于保存到`file2.txt`.

`file1.txt`不会被修改, 但是该文件当前的缓冲内容会被送到`file2.txt`.

代替`file2`, 我们可以使用一个`shell`命令来接收这个缓冲内容.`:w !cat`会显示当前的内容.

所以如果我们没有权限编辑文件,`:w`不能编辑一个只读文件. 但是可以把缓冲内容输出到`shell`，所以我们用`tee`接收这个缓冲输出.

**`tee`**

在上述命令, 使用`tee`接收到输出, 但是我们忽略标准输出, 直接保存到文件.



