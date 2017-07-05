# setfacl 和 getfacl

---

### ACL

ACL用于实现在原有的访问控制机制之外补充一种文件访问控制机制.

### 用户能否访问文件, 需要经过以下几重判断

1. 用户是否为文件属主?
2. 用户是否有特定的访问控制条目?
3. 用户是否属于文件属组?
4. 用户所属的组是否有特定的访问控制条目?
5. 其它

### 设置 ACL

**参数**

* `-b,--remove-all`：删除所有扩展的acl规则，基本的acl规则(所有者，群组，其他）将被保留.
* `-k,--remove-default`：删除缺省的acl规则。如果没有缺省规则,将不提示.
* `-n，--no-mask`：不要重新计算有效权限.setfacl默认会重新计算ACL mask,除非mask被明确的制定.
* `--mask`：重新计算有效权限,即使ACL mask被明确指定.
* `-d，--default`：设定默认的acl规则.
* `--restore=file`：从文件恢复备份的acl规则(这些文件可由getfacl -R产生)。通过这种机制可以恢复整个目录树的acl规则.此参数不能和除--test以外的任何参数一同执行.
* `--test`：测试模式,不会改变任何文件的acl规则,操作后的acl规格将被列出.
* `-R,--recursive`:递归的对所有文件及目录进行操作.
* `-L,--logical`：跟踪符号链接,默认情况下只跟踪符号链接文件,跳过符号链接目录.
* `-P,--physical`：跳过所有符号链接,包括符号链接文件.
* `--version`：输出setfacl的版本号并退出.
* `--help`：输出帮助信息.
* `--`：标识命令行参数结束，其后的所有参数都将被认为是文件名.
* `-`：如果文件名是-，则setfacl将从标准输入读取文件名.

**Access Entry**

`ACL`是由一系列`Access Entry`所组成的, 每一条`Access Entry`定义了特定的类别可以对文件拥有的操作权限.

`Access Entry`:

* Entry tag type
* qualifier(optional)
* permission

**Entry tag type**

* ACL_USER_OBJ: 相当于Linux里file_owner的permission
* ACL_USER: 定义了额外的用户可以对此文件拥有的permission
* ACL_GROUP_OBJ: 相当于Linux里group的permission
* ACL_GRUP: 定义了额外的组可以对此文件拥有的permission
* ACL_MASK: 定了`ACL_USER`, `ACL_GROUP_OBJ` 和 `ACL_GROUP`的最大权限(==注意只有三个==)
* ACL_OTHER: 相当于Linux里other的permission

**示例**

```shell
[root@localhost ~]# getfacl ./test.txt 
# file: test.txt 
# owner: root 
# group: admin 
user::rw- 
user:john:rw- 
group::rw- 
group:dev:r-- 
mask::rw- 
other::r--
```

```shell
user::rw- 定义了ACL_USER_OBJ,说明file owner拥有read and write permission 
user:john:rw- 定义了ACL_USER,这样用户john就拥有了对文件的读写权限,实现了我们一开始要达到的目的 group::rw- 定义了ACL_GROUP_OBJ,说明文件的group拥有read and write permission 
group:dev:r-- 定义了ACL_GROUP,使得dev组拥有了对文件的read permission 
mask::rw- 定义了ACL_MASK的权限为read and write 
other::r-- 定义了ACL_OTHER的权限为read
```

**设置ACL文件**

`Access Entry` 由三个被`:`号分隔开的字段所组成.

1. Entry tag type:
	* user 对应了 ACL_USER_OBJ 和 ACL_USER
	* group 对应了 ACL_GROUP_OBJ 和 ACL_GROUP
	* mask 对应了 ACL_MASK
	* other 对应了 ACL)OTHER
2. qualifier 就是上面例子中的`john用户`和`dev组`, 定了特定用户和用户组对于文件的权限. 只有`user`和`group`才有`qualifier`

**示例**

```shell
我们创建一个文件
[ansible@rancher-agent-1 ~]$ (umask 022; touch test)
[ansible@rancher-agent-1 ~]$ ll -s
0 -rw-r--r-- 1 ansible ansible       0 Jul  5 16:54 test

该文件只有属主才有写权限
[ansible@rancher-agent-1 ~]$ su chloroplast
Password:
[chloroplast@rancher-agent-1 ansible]$ echo 111 > test
bash: test: Permission denied
[chloroplast@rancher-agent-1 ansible]$ exit
exit
[ansible@rancher-agent-1 ~]$ getfacl test
# file: test
# owner: ansible
# group: ansible
user::rw-
group::r--
other::r--

设置acl
[ansible@rancher-agent-1 ~]$ setfacl -m user:chloroplast:rw- test
[ansible@rancher-agent-1 ~]$ getfacl test
# file: test
# owner: ansible
# group: ansible
user::rw-
user:chloroplast:rw-
group::r--
mask::rw-
other::r--

这里我们设置的目的是让 chloroplast 对该目录可以访问, 所以给予 x 和 r 权限
[ansible@rancher-agent-1 ~]$ setfacl -m u:chloroplast:r-x /home/ansible/
[ansible@rancher-agent-1 ~]$ getfacl /home/ansible/
getfacl: Removing leading '/' from absolute path names
# file: home/ansible/
# owner: ansible
# group: ansible
user::rwx
user:chloroplast:r-x
group::---
mask::r-x
other::---

[ansible@rancher-agent-1 ~]$ su chloroplast
Password:
[chloroplast@rancher-agent-1 ansible]$ ll -s
total 6904
 224 -rw-rw-r--  1 ansible ansible  228646 May 21 02:15 KeePassX-2.0.3.dmg
  32 -rw-rw-r--  1 ansible ansible   32038 May 21 02:15 KeePassX-2.0.3.dmg.1
6624 -rw-rw-r--  1 ansible ansible 6776016 Oct  9  2016 KeePassX-2.0.3.dmg.2
   8 -rw-r--r--  1 ansible ansible    5860 Jul  4 18:46 httpd.crt
  12 -rw-------  1 ansible ansible   11459 May 21 02:43 nohup.out
   4 -rw-rw-r--+ 1 ansible ansible       3 Jul  5 17:06 test
[chloroplast@rancher-agent-1 ansible]$ echo 11 > test
```

### ACL_MASK 和 Effective permission

```shell
[ansible@rancher-agent-1 ~]$ ll -s /home/
total 8
4 drwxr-x---+ 7 ansible     ansible     4096 Jul  5 16:54 ansible
[ansible@rancher-agent-1 ~]$ getfacl /home/ansible/
getfacl: Removing leading '/' from absolute path names
# file: home/ansible/
# owner: ansible
# group: ansible
user::rwx
user:chloroplast:r-x
group::---
mask::r-x
other::---
```

在权限后面多了一个`+`, 当任何一个文件拥有了`ACL_USER`或者`ACL_GROUP`的值以后我们就可以称它为`ACL`文件, 这个`+`就是用来提示我们的.

我们在观察,通过`getfacl`获取的权限 和 `ll -s`获取到的权限不一致.

`ll -s`中组权限是`r-x`, 而`getfacl`的权限是`group::---`.

==在Linux file permission里面大家都知道比如对于rw-rw-r--来说,,当中的那个rw-是指文件组的permission.但是在ACL里面这种情况只是在ACL_MASK不存在的情况下成立,如果文件有ACL_MASK值,那么当中那个rw-代表的就是mask值而不再是group permission了.==

**示例**

```shell
[root@localhost ~]# ls -l 
-rwxrw-r-- 1 root admin 0 Jul 3 23:10 test.sh
```

这里说明test.sh文件只有file owner: root拥有read,write, execute permission. 
admin组只有read and write permission,现在我们想让用户john也对test.sh具有和root一样的permission

```shell
[root@localhost ~]# setfacl -m user:john:rwx ./test.sh 
[root@localhost ~]# getfacl --omit-header ./test.sh 
user::rwx user:john:rwx 
group::rw- 
mask::rwx 
other::r--
```

这里我们看到john已经拥有了rwx的permission,mask值也被设定为rwx,那是因为它规定了`ACL_USER`,`ACL_GROUP`和`ACL_GROUP_OBJ`的最大值,现在我们再来看test.sh的Linux permission,它已经变成了:

```shell
[root@localhost ~]# ls -l 
-rwxrwxr--+ 1 root admin 0 Jul 3 23:10 test.sh
```

那么如果现在`admin`组的用户想要执行test.sh的程序会发生什么情况呢?它会被permission deny.==原因在于实际上admin组的用户只有read and write permission,这里当中显示的`rwx`是`ACL_MASK`的值而不是group的permission==.

==所以从这里我们就可以知道,如果一个文件后面有`+`标记，我们都需要用`getfacl`来确认它的`permission`,以免发生混淆.==

现在我们设置test.sh的mask为read only,那么admin组的用户还会有write permission吗?

```shell
[root@localhost ~]# setfacl -m mask::r-- ./test.sh 
[root@localhost ~]# getfacl --omit-header ./test.sh 
user::rwx 
user:john:rwx #effective:r-- 
group::rw- #effective:r-- 
mask::r-- 
other::r--
```

这时候我们可以看到`ACL_USER`和`ACL_GROUP_OBJ`旁边多了个`#effective:r--`,这是什么意思呢?让我们再来回顾一下`ACL_MASK`的定义.它规定了`ACL_USER`,`ACL_GROUP_OBJ`和`ACL_GROUP`的最大权限.那么在我们这个例子中他们的最大权限也就是`read only`.虽然我们这里给`ACL_USER`和`ACL_GROUP_OBJ`设置了其他权限,但是他们==真正有效果的只有`read`权限==. 这时我们再来查看test.sh的Linux file permission时它的group permission也会显示其mask的值(i.e. r--).

```shell
[root@localhost ~]# ls -l 
-rwxr--r--+ 1 root admin 0 Jul 3 23:10 test.sh
```

#### Default ACL

参数 `-d`.

Default ACL是指对于一个目录进行Default ACL设置,并且在此目录下建立的文件都将继承此目录的ACL.

用`root`用户建立一个`dir`目录:

```shell
[root@rancher-agent-1 data]# mkdir dir
[root@rancher-agent-1 data]# ll -s
total 36
 4 drwxr-xr-x  2 root    root      4096 Jul  5 17:35 dir
```

希望所有在此目录下建立的文件都可以被`chloroplast`用户所访问,那么我们就应该对`dir`目录设置`Default ACL`.

```shell
注意这里多了一个参数 -d
[root@rancher-agent-1 data]# setfacl -d -m user:chloroplast:rw ./dir
[root@rancher-agent-1 data]# getfacl dir/
# file: dir/
# owner: root
# group: root
user::rwx
group::r-x
other::r-x
default:user::rwx
default:user:chloroplast:rw-
default:group::r-x
default:mask::rwx
default:other::r-x
```

以看到ACL定义了default选项,`chloroplast`用户拥有了`default`的`read`,`write`,`execute`权限.
所有没有定义的default都将从file permission里copy过来,现在`root`用户在`dir`下建立一个`test`文件.

```shell
[root@rancher-agent-1 data]# cd dir/
[root@rancher-agent-1 dir]# touch test
[root@rancher-agent-1 dir]# getfacl test
# file: test
# owner: root
# group: root
user::rw-
user:chloroplast:rw-
group::r-x			#effective:r--
mask::rw-
other::r--
```

该文件`chloroplast`用户自动就有了`read and write`的权限了.

