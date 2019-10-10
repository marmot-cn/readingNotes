# Linux中set,env和export这三个命令的区别

## 简介

* `set`: 显示当前shell的变量，包括当前用户的变量.
* `env`: 显示当前用户的变量 .
* `export`: 显示当前导出成用户变量的shell变量 .

## 示例

```
[oracle@zhou3 ~]$ aaa=bbb --shell变量设定     
[oracle@zhou3 ~]$ echo $aaa      
bbb     
[oracle@zhou3 ~]$ env| grep aaa --设置完当前用户变量并没有     
[oracle@zhou3 ~]$ set| grep aaa  --shell变量有     
aaa=bbb     
[oracle@zhou3 ~]$ export| grep aaa --这个指的export也没导出，导出变量也没有     
[oracle@zhou3 ~]$ export aaa   --那么用export 导出一下     
[oracle@zhou3 ~]$ env| grep aaa  --发现用户变量内 存在了     
aaa=bbb  
```
 
## 总结

每个shell都有自己特有的变量, 这和用户变量是不同的. 当前用户变量和你用什么shell无关, 不管你用什么shell都是存在的. 比如HOME,SHELL等这些变量, 但shell自己的变量, 不同的shell是不同的, 比如BASH_ARGC, BASH等, 这些变量只有set才会显示, 是bash特有的. export不加参数的时候, 显示哪些变量被导出成了用户变量, 因为一个shell自己的变量可以通过export"导出"变成一个用户变量. 

## `unset`清除变量

使用`unset`命令来清除环境变量.

`set`, `env`, `export`设置的变量, 都可以用`unset`来清除的.

```
[ansible@demo ~]$ export TEST="Test..."
[ansible@demo ~]$ env|grep TEST
TEST=Test...
[ansible@demo ~]$ unset TEST
[ansible@demo ~]$ env|grep TEST
[ansible@demo ~]$
```

## `readonly`命令设置只读变量

```
$ export TEST="Test..." #增加一个环境变量TEST  
$ readonly TEST #将环境变量TEST设为只读  
$ unset TEST #会发现此变量不能被删除  
-bash: unset: TEST: cannot unset: readonly variable  
$ TEST="New" #会发现此也变量不能被修改  
-bash: TEST: readonly variable  
```

## 变量的配置文件

* `~/.bash_profile` 用户登录时被读取, 其中包含的命令被执行.
* `~/.bashrc` 启动新的shell时被读取, 并执行.
* `~/.bash_logout` shell登录退出时被读取.

### `bash`的初始化过程

`bash`会依次检测下列文件, 如果存在就读取, 否则, 跳过.

* 主目录下的文件`.bash_profile`.
* 主目录下的`.bash_login`.
* 主目录下的文件`.profile`.