# docker 1.10.3 添加 SYS_PTRACE

我本机是新的服务器是`docker ce`版, 所以直接添加此能力没有任何问题. 

但是在`1.10.3`版本中我加了此能力但是始终无效.

后来查阅了文档发现还有一个`Seccomp Profiles`限制了此权限. 需要添加额外参数`--security-opt seccomp:unconfined`, 让其不适用任何配置文件来限制.

```
cap_add:
  - SYS_PTRACE
security_opt:
  - seccomp:unconfined
```

加了后容器正常获取到了能力.