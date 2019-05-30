# 04 | 快速上手几个Linux命令：每家公司都有自己的黑话

## 学习

### 用户与密码

修改密码`passwd`.

创建用户`useradd`.

创建用户, 默认创建一个同名组.

**Linux 是"命令行 + 文件"模式**, 通过命令创建的用户其实是放在`/etc/passwd`文件里面的.

组的信息我们放在`/etc/group`文件中.

`/etc/passwd`中的`/bin/bash`的位置是用于配置登陆后的默认交互命令行. `Windows`登录进去是界面, 其实就是`explorer.exe`, 而`Linux`登录后的交互命令行是一个解析脚本的程序, 这里配置的是`/bin/bash`.

### 安装软件

`Linux`下面的安装包有常用的两大体系, 一个是`CentOS`体系, 一个是`Ubuntu`体系, 前者使用`rpm`, 后者使用`deb`.

`rpm -qa`和`dpkg -l`查看安装的软件列表.

* `-q`是`query`
* `-a`是`all`
* `-l`是`lis`

```
rpm -qa | grep xxx
```

`|`是管道, 用于连接两个程序. 前面`rpm -qa`的输出就放进管道里面, 然后作为`grep`的输入, `grep`将在里面进行搜索带关键词`xxx`的行, 并且输出出来.

* `more`是分页后只能往后翻页, 翻到最后一页自动结束返回命令行
* `less`是往前往后都能翻页, 需要输入`q`返回命令行

**删除**, 可以用`rpm -e`和`dpkg -r`.

* `-e`是`erase`
* `-r`是`remove`

也可以通过`yum`或`apt-get`安装和卸载

* `yum install xxx` 和 `yum erase xxx`
* `apt-get instal xxx` 和 `apt-get purge xxx`

`Linux`允许我们配置从哪里下载软件, 地点就在配置文件里面.

对于`Centos`, 配置文件在`/etc/yum.repos.d/CentOS-Base.repo`里.

```
[base]
name=CentOS-$releasever - Base - 163.com
baseurl=http://mirrors.163.com/centos/$releasever/os/$basearch/
gpgcheck=1
gpgkey=http://mirrors.163.com/centos/RPM-GPG-KEY-CentOS-7

```

对于与`Ubuntu`来讲, 配置文件在`/etc/apt/sources.list`里

```
deb http://mirrors.163.com/ubuntu/ xenial main restricted universe multiverse
deb http://mirrors.163.com/ubuntu/ xenial-security main restricted universe multiverse
deb http://mirrors.163.com/ubuntu/ xenial-updates main restricted universe multiverse
deb http://mirrors.163.com/ubuntu/ xenial-proposed main restricted universe multiverse
deb http://mirrors.163.com/ubuntu/ xenial-backports main restricted universe multiverse
```

**其实无论是先下载再安装, 还是通过软件关键进行安装, 都是下载一些文件, 然后将这些文件放在某个路径下, 然后在相应的配置文件中配置一下**

* `Windows`会变成`C:\Program Files`下面的一个文件夹以及注册表里面的一些配置.
* `Linux`里面会放的更散一点.
	* 主执行文件会放在`/usr/bin`或者`/usr/sbin`下
	* 其他的库文件会放在`/var`下
	* 配置文件会放在`/etc`下

#### 通过`zip`下载安装



## 扩展