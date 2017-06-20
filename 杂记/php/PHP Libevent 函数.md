# PHP Libevent 函数

---

### Reactor(反应器)模式

`libevent` 是一个典型的 `reactor` 模式的实现.

**普通函数调用机制如下**:

1. 程序调用某个函数
2. 函数执行
3. 程序等待
4. 函数将结果返回给调用者(如果含有函数返回值)

顺序执行

**Reactor**:

* 引用程序需要提供相应的接口并且注册到`reactor`反应器上
* 如果相应的时间发生的话
* `reactor`将自动调用相应的注册的接口函数(类似于回调函数)通知你

### 安装

```sh
git clone https://github.com/nmathewson/Libevent.git
cd Libevent/
root@bb0aa005831b:~/Libevent# sh autogen.sh
autoreconf: Entering directory `.'
autoreconf: configure.ac: not using Gettext
autoreconf: running: aclocal -I m4 --output=aclocal.m4t
Can't exec "aclocal": No such file or directory at /usr/share/autoconf/Autom4te/FileUtils.pm line 326.
autoreconf: failed to run aclocal: No such file or directory

运行报错,因为我们需要 automake

root@bb0aa005831b:~/Libevent# apt-get update
Get:1 http://security.debian.org jessie/updates InRelease [63.1 kB]
Get:2 http://security.debian.org jessie/updates/main amd64 Packages [523 kB]
Ign http://httpredir.debian.org jessie InRelease
Get:3 http://httpredir.debian.org jessie-updates InRelease [145 kB]
Get:4 http://httpredir.debian.org jessie Release.gpg [2373 B]
Get:5 http://httpredir.debian.org jessie Release [148 kB]
Get:6 http://httpredir.debian.org jessie-updates/main amd64 Packages [17.8 kB]
Get:7 http://httpredir.debian.org jessie/main amd64 Packages [9065 kB]
Fetched 9965 kB in 24s (412 kB/s)
Reading package lists... Done
root@bb0aa005831b:~/Libevent# apt-get install automake
Reading package lists... Done
Building dependency tree
Reading state information... Done
The following extra packages will be installed:
  autotools-dev
The following NEW packages will be installed:
  automake autotools-dev
0 upgraded, 2 newly installed, 0 to remove and 57 not upgraded.
Need to get 795 kB of archives.
After this operation, 1872 kB of additional disk space will be used.
Do you want to continue? [Y/n] Y
Get:1 http://httpredir.debian.org/debian/ jessie/main automake all 1:1.14.1-4+deb8u1 [724 kB]
Get:2 http://httpredir.debian.org/debian/ jessie/main autotools-dev all 20140911.1 [70.5 kB]
Fetched 795 kB in 2s (321 kB/s)
debconf: delaying package configuration, since apt-utils is not installed
Selecting previously unselected package autotools-dev.
(Reading database ... 14099 files and directories currently installed.)
Preparing to unpack .../autotools-dev_20140911.1_all.deb ...
Unpacking autotools-dev (20140911.1) ...
Selecting previously unselected package automake.
Preparing to unpack .../automake_1%3a1.14.1-4+deb8u1_all.deb ...
Unpacking automake (1:1.14.1-4+deb8u1) ...
Setting up autotools-dev (20140911.1) ...
Setting up automake (1:1.14.1-4+deb8u1) ...
update-alternatives: using /usr/bin/automake-1.14 to provide /usr/bin/automake (automake) in auto mode
root@bb0aa005831b:~/Libevent#

继续运行,还是报错
root@bb0aa005831b:~/Libevent# sh autogen.sh
autoreconf: Entering directory `.'
autoreconf: configure.ac: not using Gettext
autoreconf: running: aclocal --force -I m4
autoreconf: configure.ac: tracing
autoreconf: configure.ac: not using Libtool
autoreconf: running: /usr/bin/autoconf --force
configure.ac:129: error: possibly undefined macro: AC_PROG_LIBTOOL
      If this token and others are legitimate, please use m4_pattern_allow.
      See the Autoconf documentation.
autoreconf: /usr/bin/autoconf failed with exit status: 1

安装libtool
root@bb0aa005831b:~/Libevent# apt-get install libtool
Reading package lists... Done
Building dependency tree
Reading state information... Done
The following extra packages will be installed:
  libltdl-dev libltdl7
Suggested packages:
  libtool-doc gfortran fortran95-compiler gcj-jdk
The following NEW packages will be installed:
  libltdl-dev libltdl7 libtool
0 upgraded, 3 newly installed, 0 to remove and 57 not upgraded.
Need to get 393 kB of archives.
After this operation, 1784 kB of additional disk space will be used.
Do you want to continue? [Y/n] Y
Get:1 http://httpredir.debian.org/debian/ jessie/main libltdl7 amd64 2.4.2-1.11+b1 [45.4 kB]
Get:2 http://httpredir.debian.org/debian/ jessie/main libtool all 2.4.2-1.11 [190 kB]
Get:3 http://httpredir.debian.org/debian/ jessie/main libltdl-dev amd64 2.4.2-1.11+b1 [157 kB]
Fetched 393 kB in 4s (91.5 kB/s)
debconf: delaying package configuration, since apt-utils is not installed
Selecting previously unselected package libltdl7:amd64.
(Reading database ... 14248 files and directories currently installed.)
Preparing to unpack .../libltdl7_2.4.2-1.11+b1_amd64.deb ...
Unpacking libltdl7:amd64 (2.4.2-1.11+b1) ...
Selecting previously unselected package libltdl-dev:amd64.
Preparing to unpack .../libltdl-dev_2.4.2-1.11+b1_amd64.deb ...
Unpacking libltdl-dev:amd64 (2.4.2-1.11+b1) ...
Selecting previously unselected package libtool.
Preparing to unpack .../libtool_2.4.2-1.11_all.deb ...
Unpacking libtool (2.4.2-1.11) ...
Setting up libltdl7:amd64 (2.4.2-1.11+b1) ...
Setting up libltdl-dev:amd64 (2.4.2-1.11+b1) ...
Setting up libtool (2.4.2-1.11) ...
Processing triggers for libc-bin (2.19-18+deb8u6) ...

root@bb0aa005831b:~/Libevent# sh autogen.sh
autoreconf: Entering directory `.'
autoreconf: configure.ac: not using Gettext
autoreconf: running: aclocal --force -I m4
autoreconf: configure.ac: tracing
autoreconf: running: libtoolize --copy --force
libtoolize: putting auxiliary files in `.'.
libtoolize: copying file `./ltmain.sh'
libtoolize: putting macros in AC_CONFIG_MACRO_DIR, `m4'.
libtoolize: copying file `m4/libtool.m4'
libtoolize: copying file `m4/ltoptions.m4'
libtoolize: copying file `m4/ltsugar.m4'
libtoolize: copying file `m4/ltversion.m4'
libtoolize: copying file `m4/lt~obsolete.m4'
autoreconf: running: /usr/bin/autoconf --force
autoreconf: running: /usr/bin/autoheader --force
autoreconf: running: automake --add-missing --copy --force-missing
configure.ac:25: installing './compile'
configure.ac:33: installing './config.guess'
configure.ac:33: installing './config.sub'
configure.ac:13: installing './install-sh'
configure.ac:13: installing './missing'
Makefile.am: installing './depcomp'
parallel-tests: installing './test-driver'
autoreconf: Leaving directory `.'

root@bb0aa005831b:~/Libevent# ./configure --prefix=/usr/local/libevent
root@bb0aa005831b:~/Libevent# make && make install

报错
```

编译到这里还是报错,后续处理.

