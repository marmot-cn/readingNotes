# docker 挂载主机目录用户映射

---

`33`是`phpfpm`内的`www-data`用户.

```shell
greoupadd www-data -g 33
useradd www-data -u 33 -g www-data
passwd www-data
...
```