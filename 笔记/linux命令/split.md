# split 命令

## 简介

**split**命令可以将一个大文件分割成很多个小文件, 有时需要将文件分割成更小的片段, 比如为提高可读性, 生成日志等.

## 选项

* `-b`: **值**为每一输出档案的大小, 单位为`byte`.
* `-C`: 每一输出档中, 单行的最大`byte`数.
* `-d`: 使用数字作为后缀.
* `-l`: **值**为每一输出档的列数大小.

## 实例

生成一个大小为`100KB`的测试文件:

```
[ansible@iZ94ebqp9jtZ ~]$ dd if=/dev/zero bs=100k count=1 of=data.file
1+0 records in
1+0 records out
102400 bytes (102 kB) copied, 0.000303572 s, 337 MB/s
```

使用`split`命令将上面创建的`data.file`文件分割成大小为10KB的小文件:

```
[ansible@iZ94ebqp9jtZ test]$ split -b 10k data.file
[ansible@iZ94ebqp9jtZ test]$ ls
data.file  xaa  xab  xac  xad  xae  xaf  xag  xah  xai  xaj
```

文件被分割成多个带有字母的后缀文件, 如果想用数字后缀可使用**-d参数**, 同时可以使用**-a length**来指定后缀的长度:

```
[ansible@iZ94ebqp9jtZ test]$ split -b 10k data.file -d -a 3
[ansible@iZ94ebqp9jtZ test]$ ls
data.file  x000  x001  x002  x003  x004  x005  x006  x007  x008  x009
```

为分割后的文件指定文件名的前缀:

```
[ansible@iZ94ebqp9jtZ test]$ split -b 10k data.file -d -a 3 split_file
[ansible@iZ94ebqp9jtZ test]$ ls
data.file      split_file001  split_file003  split_file005  split_file007  split_file009  
split_file000  split_file002  split_file004  split_file006  split_file008
```

使用`-l`选项根据文件的行数来分割文件, 例如把文件分割成每个包含10行的小文件:

```
split -l 10 data.file
```


