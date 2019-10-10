# docker 不能删除 device or resource busy

## 场景描述

`docker`删除容器, 会有如下提示, 导致删除失败.

```
rm: cannot remove ‘docker/devicemapper/mnt/c1b69563d2b817b729e875f50f9f5d29206d15f65d823c864c8444aa3c6030dd’: Device or resource busy
rm: cannot remove ‘docker/containers/c1b69563d2b817b729e875f50f9f5d29206d15f65d823c864c8444aa3c6030dd/secrets’: Device or resource busy
```

我们根据docker的`PID`检索`/proc/$PID/mountinfo`

```
docker/devicemapper/mnt/xxxx

grep -l xxxx /proc/*/mountinfo
/proc/8441/mountinfo
/proc/8442/mountinfo
```

随后使用`ps`命令查看相关进程.

```
ps -f 84441
```

根据进程`id`, 关闭响应的服务器, 或者`kill`即可.