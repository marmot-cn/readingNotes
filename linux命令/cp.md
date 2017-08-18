# cp

---

`cp`命令目录基本操作`cp`命令用来将一个或多个源文件或者目录复制到指定的目地文件或目录.

### `-b`

覆盖已存在的文件目标前将目标文件备份.

```shell
[ansible@localhost ~]$ touch a
[ansible@localhost ~]$ touch b
[ansible@localhost ~]$ ls
a  b
[ansible@localhost ~]$ cp -b a b
[ansible@localhost ~]$ ls
a  b  b~
```

`b~`为备份文件.

### `-S`

在备份文件时,用指定的后缀“SUFFIX”代替文件的默认后缀.

```shell
[ansible@localhost ~]$ echo 1 > a
[ansible@localhost ~]$ echo 2 > b
[ansible@localhost ~]$ cp -S a b b.backup
[ansible@localhost ~]$ ll -s
total 12
4 -rw-rw-r--. 1 ansible ansible 2 Aug 18 01:00 a
4 -rw-rw-r--. 1 ansible ansible 2 Aug 18 01:01 b
4 -rw-rw-r--. 1 ansible ansible 2 Aug 18 01:01 b.backup
[ansible@localhost ~]$ cat b.backup
2
[ansible@localhost
```

### `-a`

参数的效果和同时指定"`-dpR`"参数相同.

### `-d`

当复制符号连接时,**把目标文件或目录也建立为符号连接**,并指向与源文件或目录连接的原始文件或目录.

### `-p`

保留源文件或目录的属性.

### `-u`

使用这项参数后只会在源文件的更改时间较目标文件更新时或是名称相互对应的目标文件并不存在时, 才复制文件.

### `-R/r`

递归处理,将指定目录下的所有文件与子目录一并处理

### `-f`

强行复制文件或目录, 不论目标文件或目录是否已存在.

### `-i`

覆盖既有文件之前先询问用户.

### `-l`

对源文件建立硬连接, 而非复制文件.

### `-s`

对源文件建立符号连接, 而非复制文件.

```shell
[ansible@localhost ~]$ cp -s a c
[ansible@localhost ~]$ ll -s
total 0
0 -rw-rw-r--. 1 ansible ansible 0 Aug 18 00:48 a
0 -rw-rw-r--. 1 ansible ansible 0 Aug 18 00:48 b
0 -rw-rw-r--. 1 ansible ansible 0 Aug 18 00:48 b~
0 lrwxrwxrwx. 1 ansible ansible 1 Aug 18 00:57 c -> a
```
