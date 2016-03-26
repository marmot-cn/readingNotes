#09_04 Linux压缩及归档

###笔记

---

**压缩,解压缩命令**

`压缩格式`:

* gz
* bz2
* xz
* zip
* Z

`压缩算法`: 算法不同,压缩比也会不同.

**gzip**

结尾: `.gz`

`gzip` `/PATH/TO/SOMEFILE`: 压缩完成后会删除源文件.

* `-d`: 解压缩
* `-#`: 1-9, 指定压缩比, 默认是6
* `-k`: 压缩时`保留`源文件.

**gunzip**

`gunzip` `/PATH/TO/SOMEFILE.gz`: 解压完成后会删除源文件.

**zcat**

`zcat` `/PATH/TO.SOMEFILE.gz`: `不解压`的情况,查看文件内容.

**bzip2**

结尾: `.bz2`

比`gzip`有着`更大压缩比`的工具,使用格式近似

`bzip2` `/PATH/TO/SOMEFILE`: 压缩完成后会删除源文件.

* `-d`: 解压缩
* `-#`: 1-9, 指定压缩比, 默认是6
* `-k`: 压缩时`保留`源文件.

**bunzip2**

`bunzip2` `/PATH/TO/SOMEFILE.bz`: 解压完成后会删除源文件.

**bcat**

`bcat` `/PATH/TO.SOMEFILE.bz`: `不解压`的情况,查看文件内容.

**xz**

不能压缩目录,压缩完会删除源文件.

结尾: `.xz`

`xz` /PATH/TO/SOMEFILE

* `-d`: 解压缩
* `-#`: 1-9, 指定压缩比, 默认是6
* `-k`: 压缩时`保留`源文件.

**unxz**

`unxz` `/PATH/TO/SOMEFILE.bz`

**xzcat**

`xzcat` `/PATH/TO/SOMEFILE.bz`

**zip**

即`归档`又`压缩`的工具

`zip` `FILENAME.zip` `FILE1 FILE2 …`: 压缩后删除源文件.

`unzip` `FILENAME.zip`

`archive`: 归档,归档本身并`不意味着压缩`.

**tar**

归档工具,只归档不压缩.

* `-c`: 创建新的归档文件
* `-f FILE.tar`: 使用档名,归档后叫什么名字
* `-x`: 展开归档
* `--xattrs`: 归档时,保留文件的扩展属性信息
* `-t`: 不展开归档,直接查看归档了哪些文件
* `-zcf`: 归档并调用`gzip`压缩
* `-zxf`: 调用`gzip`解压缩并展开归档

* `-jcf`: `bzip2`* `-jxf`:* `-Jcf`: xz* `-Jxf`:


归档:

		# tar -cf test.tar test*.txt
		
展开:

		# tar -xf test.tar
		
**cpio**

归档工具.


**while循环**

while 循环: 适用于循环次数未知的场景,要有退出条件.

语法:

		while CONDITION;do
			statement
			...
		done

###整理知识点

---