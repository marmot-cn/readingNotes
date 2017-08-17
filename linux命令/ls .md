# ls 

---

`ls`命令用于列出目录中的文件

## `ls -F`

* 对于**目录**, 以`/`结尾.
* 对于**普通文件**, 结尾没有任何特殊字符.
* 对于**软链接**文件, 后缀以`@`结尾.
* 对于**可执行**文件, 以`*`结尾.

```shell
[ansible@iZ944l0t308Z ~]$ ls -F /usr/bin/
Mail@                curl*                  geoiplookup6* ...

ls -F /etc/
DIR_COLORS               bash_completion.d/       csh.login ...
```

