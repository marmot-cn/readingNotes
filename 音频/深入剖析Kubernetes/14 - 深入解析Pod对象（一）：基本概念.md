# 14 | 深入解析Pod对象（一）：基本概念

## 笔记

`Pod`是`Kubernetes`项目中的最小编排单位. 容器`Container`就成了`Pod`属性里的一个普通的字段.

**Pod扮演的是传统部署环境里"虚拟机"的角色**.

**调度, 网络, 存储, 以及安全相关的属性, 基本上是 Pod 级别的**.

它们描述的是"机器"这个整体, 不是里面运行的"程序".

* "Pod"的网络定义: 配置这个"机器"的网卡
* "Pod"的存储定义: 机器的磁盘
* "Pod"的安全定义: 机器的防火墙
* "Pod"的调度: 运行在哪个服务器上

### Pod 重要字段

* `NodeSelector`: 供用户将`Pod`与`Node`进行绑定的字段(通过标签调度)
* `NodeName`: 一旦`Pod`的这个字段被赋值, `Kubernetes`会被认为这个`Pod`已经进过了调度, 调度的结果就是赋值的节点名字. **这个字段一般由调度器负责设置, 用户可以设置它来"骗过"调度器**.
* `HostAliases`: 定了`Pod`的`hosts`文件(如`/etc/hosts`)里面的内容.

**凡是跟容器的Linux Nmamespace相关的属性, 也一定是Pod级别的**.

**Pod的设计**, 就是要让它里面的容器尽可能多地共享Linux Nmaespace, 仅保留必要的隔离和限制能力. 这样, `Pod`模拟出的效果, 就是跟虚拟机里程序间的关系非常类似了.

`shareProcessNamespace=true`, 意味着这个`Pod`里的容器要共享`Pid Namespace`.

**凡是Pod中的容器要共享宿主机的Namespace, 也一定是Pod级别的定义**

`init Containers`的生命周期, 会限于所有的`Containers`, 并且严格按照定义的顺序执行.

`ImagePullPolicy`

* `Always`: 总是拉取.
* `Never`: 永不拉取.
* `IfNotPresent`: 不存在这个镜像时才拉取

`Lifecycle`, 在容器状态发生变化时触发一系列"钩子"

* `postStart`: 在容器启动后, 立刻执行一个指定的操作. 不严格保证顺序, `postStart`启动时, `ENTRYPOINT`有可能没有结束.
* `preStop`: 在容器被杀死之前, 它会阻塞当前的容器杀死流程, 直到这个`Hook`定义操作完成之后, 才允许容器被杀死.

`Pod`的`Status`字段, 是`Pod`的当前状态:

* `Pending`: `Pod`的`YAML`文件已经提交给`Kubernetes`, `API`对象已经被创建并保存在`Etcd`当中.
* `Running`: `Pod`已经调度成功, 跟一个具体的节点绑定. 它包含的容器都已经创建成功, 并且至少有一个正在运行中.
* `Succeeded`: `Pod`里的所有容器都正常运行完毕, 并且已经退出了, 在运行一次性任务时最为常见.
* `Failed`: 至少有一个容器以不正常的状态(非0的返回码)退出. 需要想办法`Debug`这个容器的应用, 查看`Pod`的`Events`和日志.
* `Unknows`: 这是一个异常状态. 意味着`Pod`的状态不能持续地被`kubelet`汇报给`kube-apiserver`. 有可能是主从节点(`Matser`和`Kubelet`)间的通信出现了问题.

`Pod`的`Conditions`, 主要用于描述造成当前`Status`的具体原因是什么.

* `Unscheduable`: 调度出现了问题.
* `Ready`: 意味着`Pod`不仅已经正常启动(`Running`状态), 而且已经可以对外提供服务了. 这两者之间(`Running`和`Ready`)是有区别的.

## 扩展

### `Pod`状态是`Ready`, 实际上不能提供服务的情况

1. 程序本身`bug`
2. 程序内存问题, 已经僵死, 但是进程还在, 但无响应
3. `Dockerfile`写的不规范, 应用程序不是主进程
4. 程序出现死循环