# 15 | 深入解析Pod对象（二）：使用进阶

## 笔记

### Projected Volume

投射数据卷

在`Kubernetes`中, 几种特殊的`Volume`, 存在的意义不是为了存放容器里的数据, 也不是用来进行容器和宿主机之间的数据交换. 这些特殊的`Volume`的作用, 是为容器提供预先定义好的数据. 从容器角度来看, 这些`Volume`里的信息就是仿佛**被Kubernetes"投射"(Project)进入容器当中的**.

`Projected Volume`有四种:

* `Secret`
* `ConfigMap`
* `Download API`
* `ServiceAccountToken`

### Secret

把`Pod`想要访问的加密数据, 存放到`Etcd`中. 可以通过在`Pod`的容器挂载`Volume`的方式, 访问到这些`Secret`里保存的信息.

可以通过`kubectl create secret`指令创建, 也可以通过`YAML`创建. 这些数据必须是经过`Base64`转码的, 以免出现明文密码的安全隐患.

创建完`Pod`之后, 保存在`Etcd`里的用户名和密码信息, 已经以文件形式出现在了容器的`Volume`目录里. 文件的名字就是, 指定的`Key`.

一旦其对应的`Etcd`里的数据被更新, 这些`Volume`里的文件内容, 同样也会被更新. **更新会有一定的延时, 需要编写好重试和超时的逻辑**.

### ConfigMap

与`Secret`的区别在于, `ConfigMap`保存的是不需要加密的, 应用所需的配置信息.

可以通过`kubectl create configmap`命令创建, 也可以通过`YAML`创建.

### Download API

让`Pod`里的容器能够直接获取到这个`Pod API`对象本身的信息.

## 扩展