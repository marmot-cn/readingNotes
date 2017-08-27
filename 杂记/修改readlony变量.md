# 修改readlony变量

---

今天在`/etc/profile.d/xx.sh`中设置变量的时候:

```shell
export TMOUT=900    # 设置900秒内用户无操作就字段断开终端

readonly TMOUT     # 将值设置为readonly 防止用户更改
```

设置了`readonly`之后在当前shell下是无法取消的,需要先将`/etc/profile`中设置`readonly`行注释起来或直接删除,`logout`后重新`login`.

或者通过如下方式在当前`shell`中修改.

```
$ readonly PI=3.14
$ unset PI
-bash: unset: PI: cannot unset: readonly variable
$ cat << EOF| sudo gdb
attach $$
call unbind_variable("PI")
detach
EOF
$ echo $PI

```