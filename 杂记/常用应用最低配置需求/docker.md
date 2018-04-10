# docker

---

## 内核需求

最低需求: `3.10`内核版本.

## storage driver

`Docker`默认推荐`overlay2`. 但是我们使用的`Linux`是`CentOS`所以只能退而就其次, 使用`devicemapper`, 但是使用`direct-lvm`在生产环境.