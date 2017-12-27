# 在linux下使用noatime提升文件系统性能

---

## 简介

今天在编写`mongo`在生产环境的配置文件时, 看见生产环境的配置下推荐`monut`时使用`noatime`. 查了下资料发现还有此优化.

默认的方式下linux会把文件访问的时间`atime`做记录, 这在绝大部分的场合都是没有必要的, 如果遇到机器IO负载高或是`CPU WAIT`高的情况, 可以尝试在挂载文件系统时使用`noatime`和`nodiratime`.

* `noatime`: 当访问一个文件时`access time`不会更新.
* `nodiratime`: 当访问一个目录时`access time`不会更新.

## 测试

### 指定 noatime,nodiratime

可见我的`/data`目录挂载时候指定了`noatime,nodiratime`

```
[ansible@demo data]$ cat /etc/fstab

#
# /etc/fstab
# Created by anaconda on Mon Jul 10 12:22:03 2017
#
# Accessible filesystems, by reference, are maintained under '/dev/disk'
# See man pages fstab(5), findfs(8), mount(8) and/or blkid(8) for more info
#
UUID=b7792c31-ad03-4f04-a650-a72e861c892d /                       ext4    defaults        1 1
/dev/data/data /data xfs noatime,nodiratime 0 0
swap /data/swap/swap-file swap defaults 0 0
```

```
[root@demo data]# pwd
/data

[root@demo data]# touch test ; stat test ;
  File: 'test'
  Size: 0         	Blocks: 0          IO Block: 4096   regular empty file
Device: fd03h/64771d	Inode: 98          Links: 1
Access: (0644/-rw-r--r--)  Uid: (    0/    root)   Gid: (    0/    root)
Access: 2017-12-26 16:10:30.944756283 +0800
Modify: 2017-12-26 16:10:30.944756283 +0800
Change: 2017-12-26 16:10:30.944756283 +0800
 Birth: -

[root@demo data]#  echo hello >> test ; stat test;
  File: 'test'
  Size: 6         	Blocks: 8          IO Block: 4096   regular file
Device: fd03h/64771d	Inode: 98          Links: 1
Access: (0644/-rw-r--r--)  Uid: (    0/    root)   Gid: (    0/    root)
Access: 2017-12-26 16:10:30.944756283 +0800
Modify: 2017-12-26 16:10:41.166766644 +0800
Change: 2017-12-26 16:10:41.166766644 +0800
 Birth: -
 
[root@demo data]# cat test ;stat test
hello
  File: 'test'
  Size: 6         	Blocks: 8          IO Block: 4096   regular file
Device: fd03h/64771d	Inode: 98          Links: 1
Access: (0644/-rw-r--r--)  Uid: (    0/    root)   Gid: (    0/    root)
Access: 2017-12-26 16:10:30.944756283 +0800
Modify: 2017-12-26 16:10:41.166766644 +0800
Change: 2017-12-26 16:10:41.166766644 +0800
 Birth: -
```

#### 总结

1. `read`文件的时候不会导致`atime、mtime、ctime`改变.
2. `write`文件只会导致`mtime`和`ctime`更新, 不会导致`atime`更新.

### 未指定 noatime,nodiratime

```
[root@demo data]# cd ~
[root@demo ~]# pwd
/root

[root@demo ~]# touch test ; stat test ;
  File: 'test'
  Size: 0         	Blocks: 0          IO Block: 4096   regular empty file
Device: ca01h/51713d	Inode: 262153      Links: 1
Access: (0644/-rw-r--r--)  Uid: (    0/    root)   Gid: (    0/    root)
Access: 2017-12-26 16:12:26.909874502 +0800
Modify: 2017-12-26 16:12:26.909874502 +0800
Change: 2017-12-26 16:12:26.909874502 +0800
 Birth: -
 
[root@demo ~]# echo hello >> test ; stat test;
  File: 'test'
  Size: 6         	Blocks: 8          IO Block: 4096   regular file
Device: ca01h/51713d	Inode: 262153      Links: 1
Access: (0644/-rw-r--r--)  Uid: (    0/    root)   Gid: (    0/    root)
Access: 2017-12-26 16:12:26.909874502 +0800
Modify: 2017-12-26 16:12:37.682885558 +0800
Change: 2017-12-26 16:12:37.682885558 +0800
 Birth: -
 
[root@demo ~]# cat test ;stat test
hello
  File: 'test'
  Size: 6         	Blocks: 8          IO Block: 4096   regular file
Device: ca01h/51713d	Inode: 262153      Links: 1
Access: (0644/-rw-r--r--)  Uid: (    0/    root)   Gid: (    0/    root)
Access: 2017-12-26 16:12:49.468897669 +0800
Modify: 2017-12-26 16:12:37.682885558 +0800
Change: 2017-12-26 16:12:37.682885558 +0800
 Birth: -
```

#### 总结

1. `read`文件的时候会导致`atime`更新, 不会导致`mtime`和`ctime`更新.
2. `write`文件只会导致`mtime`和`ctime`更新, 不会导致`atime`更新.