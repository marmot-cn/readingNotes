# 文件 /etc/skel

这个目录一般是存放用户启动文件的目录, 这个目录是由`root`权限控制, **当添加用户时,这个目录下的文件自动复制到新添加的用户的家目录下**.

`/etc/skel`目录下的文件都是**隐藏文件**.

```
[ansible@iZ94ebqp9jtZ skel]$ ls -a
.  ..  .bash_logout  .bash_profile  .bashrc 
```

我们可通过修改,添加,删除`/etc/skel`目录下的文件, 来为用户提供一个统一,标准的,默认的用户环境.

`/etc/skel`目录下的文件,一般是用`useradd`和`adduser`命令添加用户(`user`)时,系统自动复制到新添加用户(`user`)的家目录下;如果我们通过修改`/etc/passwd`来添加用户时,我们可以自己创建用户的家目录,然后把`/etc/skel`下的文件复制到用户的家目录下,然后要用`chown`来改变新用户家目录的属主.


## 测试

```shell

我们在/etc/skel创建一个test文件
[ansible@iZ94ebqp9jtZ skel]$ ls -a
.  ..  .bash_logout  .bash_profile  .bashrc  test

添加一个用户
[ansible@iZ94ebqp9jtZ skel]$ sudo useradd wangcai

这里我们没有添加-a显示隐藏文件, 可以看见刚才创建用户的家目录下多了一个 test 文件
[ansible@iZ94ebqp9jtZ skel]$ sudo ls /home/wangcai/
test
```
