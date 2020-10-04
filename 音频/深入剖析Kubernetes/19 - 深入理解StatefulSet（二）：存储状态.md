# 19 | 深入理解StatefulSet（二）：存储状态

## 笔记

### Persistent Volume Claim 

`StatefulSet`对存储状态的管理机制. 这个机制主要使用**Persistent Volume Claim**的功能.

可以通过定义`spec.volumes`声明`Pod`的`Volume`. 如果并不知道有哪些`Volume`类型可以使用?

下列是一个声明了`Ceph RBD`类型`Volume`的`Pod`.

```
apiVersion: v1
kind: Pod
metadata:
  name: rbd
spec:
  containers:
    - image: kubernetes/pause
      name: rbd-rw
      volumeMounts:
      - name: rbdpd
        mountPath: /mnt/rbd
  volumes:
    - name: rbdpd
      rbd:
        monitors:
        - '10.16.154.78:6789'
        - '10.16.154.82:6789'
        - '10.16.154.83:6789'
        pool: kube
        image: foo
        fsType: ext4
        readOnly: true
        user: admin
        keyring: /etc/ceph/keyring
        imageformat: "2"
        imagefeatures: "layering"
```

**过度暴露**

* 如果不懂`Ceph RBD`, 看不懂`Volumes`里面的字段.
* 暴露了`Ceph RBD`对应的存储服务器的地址, 用户名和授权文件.

后来演化中, **Kubernetes项目引入了一组叫做 Persistent Volume Clain(PVC) 和 Persistent Volume(PV) 的 API 对象, 降低了用户声明和持久化 Volume 的门槛**.

#### 第一步, 定义一个 PVC, 声明想要的 Volume 属性

```
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: pv-claim
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

没有关于`Volme`细节的字段, **只有描述性的属性和定义**(如: storage: 1Gi, 大小是1GiB. ReadWriteOnce, 只能被单个节点以读写的方式映射)

#### 第二步, 在应用的 Pod 中, 声明使用这个 PVC

```
apiVersion: v1
kind: Pod
metadata:
  name: pv-pod
spec:
  containers:
    - name: pv-container
      image: nginx
      ports:
        - containerPort: 80
          name: "http-server"
      volumeMounts:
        - mountPath: "/usr/share/nginx/html"
          name: pv-storage
  volumes:
    - name: pv-storage
      persistentVolumeClaim:
        claimName: pv-claim
```

在这个`Pod`的`Volumes`定义中, 只需要:

1. 声明它的类型是`persistentVolumeClaim`
2. 指定`PVC`的名字

只要我们创建这个`PVC`对象, `Kubernetes`就会自动为它绑定一个符合条件的`Volume`, 这些符合条件的`Volume`来自`PV(Persistent Volume)`对象.

### PV

```
kind: PersistentVolume
apiVersion: v1
metadata:
  name: pv-volume
  labels:
    type: local
spec:
  capacity:
    storage: 10Gi
  rbd:
    monitors:
    - '10.16.154.78:6789'
    - '10.16.154.82:6789'
    - '10.16.154.83:6789'
    pool: kube
    image: foo
    fsType: ext4
    readOnly: true
    user: admin
    keyring: /etc/ceph/keyring
    imageformat: "2"
    imagefeatures: "layering"
```

刚才创建的`PVC`对象就会绑定这个`PV`.

`Kubernetes`中`PVC`和`PV`, **实际上类似于"接口"和"实现”的思想.

* 开发者只要知道并会使用"接口", **PVC**
* 运维人员则负责给“接口"绑定具体的实现, **PV**

这样的**解耦**, 就避免了因为向开发者暴露过多的存储系统细节而带来的隐患. 也可以**职责分离**.

### 通过 PVC, PV 的设计, 实现 StatefulSet 对存储状态的管理

```
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  serviceName: "nginx"
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.9.1
        ports:
        - containerPort: 80
          name: web
        volumeMounts:
        - name: www
          mountPath: /usr/share/nginx/html
  volumeClaimTemplates:
  - metadata:
      name: www
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi
``` 

`volumeClaimTemplates`这个字段类似于`Deployment`里`Pod`模板(`PodTemplate`). 也就是说, 凡是被这个`StatefulSet`管理的`Pod`, 都会声明一个对应的`PVC`. 这个`PVC`的定义来自于`volumeClaimTemplates`这个模板的字段. **这个PVC的名字, 会被分配一个与这个Pod完全一致的编号**.

这个自动创建的`PVC`, 与`PV`绑定成功后, 就会进入`Bound`状态, 这个意味着这个`Pod`可以挂载并使用这个`PV`了.

```
$ kubectl create -f statefulset.yaml
$ kubectl get pvc -l app=nginx
NAME        STATUS    VOLUME                                     CAPACITY   ACCESSMODES   AGE
www-web-0   Bound     pvc-15c268c7-b507-11e6-932f-42010a800002   1Gi        RWO           48s
www-web-1   Bound     pvc-15c79307-b507-11e6-932f-42010a800002   1Gi        RWO           48s
```

这些 PVC, 都以**"--< 编号 >"**的方式命名,并且处于`Bound`状态.

删除了`Pod`, `Pod`会被自动重新创建, 原来的`PV`依然会按照编号与`Pod`绑定起来.

**怎么做到的?**

删除`Pod`后:

* 对应的`PVC`和`PV`并不会被删除.
* `StatefulSet`控制器发现, 一个`Pod`消失了, 会重新创建一个新的, 同名的`Pod`.
* 新的`Pod`, 创建出来之后, `Kubernetes`会直接找到旧的同名的`PVC`, 进而绑定.

**通过这种方式, `Kubernetes`的`StatefulSet`就实现了对应用存储状态的管理**.

### 总结

* `StatefulSet` 的控制器直接管理的是 `Pod`
* `Kubernetes` 通过 `Headles Service`, 为这些有编号的 `Pod`, 在 `DNS` 服务器中生成带有同样编号的 `DNS` 记录.
* `StatefulSet` 还为每一个 `Pod` 分配并创建一个同样编号的 `PVC`.



## 扩展