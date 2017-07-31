# 31_02_bash脚本编程之在bash脚本中使用选项

---

## 笔记

---

### `mkscript`

使用`mkscript`脚本自动生成脚本, 可以接受其他脚本名作为参数. 可以自动生成首部信息, 并且使用`vim`编辑器打开这个文件, 光标处于行首.

如果参数的脚本有内容, 直接打开. 不修改. 只有脚本没有内容才追加信息.

#### 示例

```shell
[root@iZ944l0t308Z ~]# cat mkscript.sh
#!/bin/bash
# Name: mkscript
# Description: Create Script
# Author: kevin
# Version: 0.0.1
# Datetime: 20170730
# Usage: mkscript FILENAME

while getopts ":d:" OPT; do
  case $OPT in
    d)
      DESC=$OPTARG ;;
    *)
      echo "Usage: mkscript [-d DESCRIPTION] FILENAME" ;;
  esac
done

#需要把前面的 -d 这些选项排除走, 否则$1就指到选项上而不是文件名上面
shift $((OPTIND-1))

#check $1 file is empty
if ! grep "[^[:space:]]" $1 &> /dev/null; then
cat > $1 << EOF
# Name: `basename $1`
# Description: $DESC
# Author: kevin
# Version: 0.0.1
# Datetime: `date '+%F %T'`
# Usage: `basename $1`

EOF
fi

vim + $1

# check syntax
until bash -n $1 &> /dev/null; do
  read -p "Syntax error, q|Q for quit, else for edit:" OPT
  case $OPT in
    q|Q)
      echo "Quit."
      exit 8 ;;
    *)
      vim + $1
      ;;
  esac
done

chmod +x $1
```

### `getopts`

获取`bash`脚本选项, 并且能获取到后续内容.

```shell
[root@iZ944l0t308Z ~]# help getopts
getopts: getopts optstring name [arg]
    Parse option arguments.

    Getopts is used by shell procedures to parse positional parameters
    as options.

    OPTSTRING contains the option letters to be recognized; if a letter
    is followed by a colon, the option is expected to have an argument,
    which should be separated from it by white space.

    Each time it is invoked, getopts will place the next option in the
    shell variable $name, initializing name if it does not exist, and
    the index of the next argument to be processed into the shell
    variable OPTIND.  OPTIND is initialized to 1 each time the shell or
    a shell script is invoked.  When an option requires an argument,
    getopts places that argument into the shell variable OPTARG.

    getopts reports errors in one of two ways.  If the first character
    of OPTSTRING is a colon, getopts uses silent error reporting.  In
    this mode, no error messages are printed.  If an invalid option is
    seen, getopts places the option character found into OPTARG.  If a
    required argument is not found, getopts places a ':' into NAME and
    sets OPTARG to the option character found.  If getopts is not in
    silent mode, and an invalid option is seen, getopts places '?' into
    NAME and unsets OPTARG.  If a required argument is not found, a '?'
    is placed in NAME, OPTARG is unset, and a diagnostic message is
    printed.

    If the shell variable OPTERR has the value 0, getopts disables the
    printing of error messages, even if the first character of
    OPTSTRING is not a colon.  OPTERR has the value 1 by default.

    Getopts normally parses the positional parameters ($0 - $9), but if
    more arguments are given, they are parsed instead.

    Exit Status:
    Returns success if an option is found; fails if the end of options is
    encountered or an error occurs.
```

`optstring`: 表示接收的选项是什么, 只能接收短选项.

如果选项多余一个, 用"`,`"分隔. 如果选项后面有内容, 则跟"`:`"即可.

#### 示例

##### 测试基本使用

```shell
[root@iZ944l0t308Z ~]# cat opttest.sh
#!/bin/bash

getopts "bd" OPT
echo $OPT

[root@iZ944l0t308Z ~]# ./opttest.sh -d
d

只能获取一个选项
[root@iZ944l0t308Z ~]# ./opttest.sh -b -d
b

非法选项
[root@iZ944l0t308Z ~]# ./opttest.sh -a
./opttest.sh: illegal option -- a
?
```

##### 测试带参数

`$OPTARG`内置变量

选项后面添加"`:`".

```shell
d 可以带参数

[root@iZ944l0t308Z ~]# cat opttest.sh
#!/bin/bash

getopts "bd:" OPT
echo $OPT
echo $OPTARG
[root@iZ944l0t308Z ~]# ./opttest.sh -d haha
d
haha
[root@iZ944l0t308Z ~]# ./opttest.sh -b haha
b


```

##### `getopts` 不输出错误信息

参数前面加"`:`"

```shell
[root@iZ944l0t308Z ~]# cat opttest.sh
#!/bin/bash

getopts ":bd:" OPT
echo $OPT
echo $OPTARG

不像以前输出./opttest.sh: illegal option -- a
[root@iZ944l0t308Z ~]# ./opttest.sh -a
?
a
```

##### 获取所有选项

`OPTIND` 选项索引

```shell
[root@iZ944l0t308Z ~]# cat opttest2.sh
#!/bin/bash

while getopts ":b:d:" SWITCH; do
  case $SWITCH in
    b) echo "The option is b."
       echo $OPTARG 
       echo $OPTIND
       ;;
    d) echo "The option is d."
       echo $OPTARG 
       echo $OPTIND
       ;;
    *) echo "Wrong option." ;;
  esac
done
[root@iZ944l0t308Z ~]# chmod +x opttest2.sh
[root@iZ944l0t308Z ~]# ./opttest2.sh -b haha
The option is b.
haha
3
[root@iZ944l0t308Z ~]# ./opttest2.sh -d hoho
The option is d.
hoho
3
[root@iZ944l0t308Z ~]# ./opttest2.sh -c
Wrong option.

同时给两个选项
[root@iZ944l0t308Z ~]# ./opttest2.sh -b haha -d hoho
The option is b.
haha
3
The option is d.
hoho
5
```

`OPTIND` 指到当前选项

`OPTIND`初始值为1, 其含义是下一个待处理的参数的索引

#### `shift`

```
./scripts a b c d

$1 = a
$2 = b
$3 = c
$4 = d

也可以说

$1 = a
shift 1个
$1 = b
```

### 示例

写一个脚本`getinterface.sh`, 脚本可以接受选项(i,I,a), 完成以下任务:

1. 使用以下形式: getinterface.sh [-i interface|-I IP|-a]
2. 当用使用`-i`选项时, 显示其指定网卡的`ip`地址;
3. 当用户使用`-I`选项时, 显示其后面的`IP`地址所属的网络接口.
4. 当用户单独使用`-a`选项时, 显示所有网络接口及其`IP`地址(`lo`除外).

#### 筛选出网卡

```shell
[root@iZ944l0t308Z ~]# ifconfig | grep "^[^[:space:]]\{1,\}"
eth0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
eth1: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
[root@iZ944l0t308Z ~]# ifconfig | grep "^[^[:space:]]\{1,\}" | awk -F ":" '{print $1}'
eth0
eth1
lo
```

#### 筛选一个网卡的ip地址

```shell
[root@iZ944l0t308Z ~]# ifconfig eth0 | grep -o "inet [0-9\.]\{1,\}" |cut -d ' ' -f 2
10.170.148.109
```

#### 根据ip筛选网卡

`grep -B 1`: 除了显示符合范本样式的那一列之外，并显示该列之前(1行)的内容.

```shell
[root@iZ944l0t308Z ~]# ifconfig
eth0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 10.170.148.109  netmask 255.255.248.0  broadcast 10.170.151.255
        ether 00:16:3e:00:30:5d  txqueuelen 1000  (Ethernet)
        RX packets 1017871  bytes 889264174 (848.0 MiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 959641  bytes 144965603 (138.2 MiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

eth1: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 120.25.87.35  netmask 255.255.252.0  broadcast 120.25.87.255
        ether 00:16:3e:00:51:a9  txqueuelen 1000  (Ethernet)
        RX packets 99845  bytes 29621205 (28.2 MiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 70033  bytes 13184361 (12.5 MiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
        inet 127.0.0.1  netmask 255.0.0.0
        loop  txqueuelen 0  (Local Loopback)
        RX packets 366  bytes 88582 (86.5 KiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 366  bytes 88582 (86.5 KiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
[root@iZ944l0t308Z ~]# ifconfig | grep -B 1 "10.170.148.109" | grep -o "^[^[:space:]]\{1,\}" | cut -d ':' -f 1
eth0
```

## 整理知识点

---

### `vim + num`

打开文件, 光标并定位到第`num`行.

如果是`vim +`则光标直接定位到结尾行.

### `bash -n FILENAME`

"`-n`"选项进行shell脚本的语法检查.

只读取shell脚本, 但不实际执行.

### `bash -x FILENAME`

"`-x`"选项实现shell脚本逐条语句的跟踪.

进入跟踪方式，显示所执行的每一条命令.

### `bash -c string`

"`-c`"选项使shell解释器从一个字符串中而不是从一个文件中读取并执行shell命令。当需要临时测试一小段脚本的执行结果时，可以使用这个选项，如下所示:

```shell
bash -c 'a=1;b=2;let c=$a+$b;echo "c=$c"' 

执行结果 c=3
```