# wc 

---

`Word Count`命令的功能为统计指定文件中的字节数,字数,行数,并将统计结果显示输出.

## 命令格式

`wc [选项] 文件`

## 命令功能

统计指定文件中的字节数,字数,行数,并将统计结果显示输出. 该命令统计指定文件中的字节数,字数,行数. 如果**没有给出文件名, 则从标准输入读取**.

## 命令参数

* `-c`: 统计字节数.
* `-l`: 统计行数.
* `-m`: 统计字符数. 这个标志不能与`-c`标志一起使用.
* `-w`: 统计字数. 一个字被定义为由空白,跳格或换行字符分隔的字符串.
* `-L`: 打印最长行的长度.

## 使用实例

```shell
[ansible@iZ944l0t308Z ~]$ cat test
hnlinux
peida.cnblogs.com
ubuntu
ubuntu linux
redhat
Redhat
linuxmint

[ansible@iZ944l0t308Z ~]$ wc -l test
7 test

[ansible@iZ944l0t308Z ~]$ wc -c test
70 test

[ansible@iZ944l0t308Z ~]$ wc -m test
70 test

[ansible@iZ944l0t308Z ~]$ wc -w test
8 test

[ansible@iZ944l0t308Z ~]$ wc -L test
17 test
```