# aufs 文件系统

---

`advanced multi-layered unification filesystem`. 能够将不同类型的文件系统透明地层叠在一起, 实现一个高效的分层文件系统. 简言之, 它能将不同的目录挂在到统一目录下.

## 挂载

挂载命令`mount -t aufs -o br=path/to/dir/A:path/to/dir/A none path/to/dir/uniondir`

* `-t`表示目标文件系统类型.
* `-o`表示挂载参数.
* `none`表示挂载的不是设备文件.
* `br`是aufs的参数, 表示分支, 多个分支之间用冒号分隔. 也可以使用`dirs`

## 示例

```shell
root@iZ94xwu3is8Z:~# echo 1 > a/test-1
root@iZ94xwu3is8Z:~# echo 2 > b/test-2
root@iZ94xwu3is8Z:~# mount -t aufs -o br=/root/a:/root/b: none /root/uniondir
root@iZ94xwu3is8Z:~# tree
.
|-- a
|   `-- test-1
|-- b
|   `-- test-2
`-- uniondir
    |-- test-1
    `-- test-2
```

将`a`目录的`test-1`和`b`目录的`test-2`挂载到了`uniondir`下.

在挂载的时候, 指定了`a`目录在前, 所以`a`目录是逻辑的上层. `br`参数没有加读写权限之前, 逻辑上层(`a`目录)为读写权限, 其余层(b目录)为只读权限. 所以如果我们在`uniondir`创建`test-3`,此文件会创建在a目录中.

```shell
root@iZ94xwu3is8Z:~# cd uniondir/
root@iZ94xwu3is8Z:~/uniondir# echo 3 > test-3
root@iZ94xwu3is8Z:~/uniondir# cd ..
root@iZ94xwu3is8Z:~# tree
.
|-- a
|   |-- test-1
|   `-- test-3
|-- b
|   `-- test-2
`-- uniondir
    |-- test-1
    |-- test-2
    `-- test-3

3 directories, 6 files
```

我们如果在一开始在`a`目录中创建`test-2`. 在挂载时候会优先挂载`test-2`. `mount`命令按照命令行中给出的文件夹顺序挂载, 若出现有同名文件的情况, 则以先挂载的为主, 其他的不再挂载.

对可读写分支(目录a), 修改直接作用到分支上.

```shell
root@iZ94xwu3is8Z:~# echo "test" >> uniondir/test-1
root@iZ94xwu3is8Z:~# cat uniondir/test-1
1
test
root@iZ94xwu3is8Z:~# cat a/test-1
1
test
root@iZ94xwu3is8Z:~# tree
.
|-- a
|   |-- test-1
|   `-- test-3
|-- b
|   `-- test-2
`-- uniondir
    |-- test-1
    |-- test-2
    `-- test-3

3 directories, 6 files
root@iZ94xwu3is8Z:~#
```

对只读分支(b), 修改操作会触发一次写时复制(COW, Copy On Write), 即先将待修改文件从目录b复制到目录`a`, 然后在目录`a`修改相应文件.

```shell
root@iZ94xwu3is8Z:~# echo "test" >> uniondir/test-2
root@iZ94xwu3is8Z:~# tree
.
|-- a
|   |-- test-1
|   |-- test-2
|   `-- test-3
|-- b
|   `-- test-2
`-- uniondir
    |-- test-1
    |-- test-2
    `-- test-3

3 directories, 7 files
root@iZ94xwu3is8Z:~# cat b/test-2
2
root@iZ94xwu3is8Z:~# cat uniondir/test-2
2
test
root@iZ94xwu3is8Z:~# cat a/test-2
2
test
```

### 删除文件

删除文件时，如果该文件只在`rw`目录下有，那就直接删除`rw`目录下的该文件，如果该文件在`ro`目录下有，那么aufs将会在`rw`目录里面创建一个`.wh`开头的文件，标识该文件已被删除
