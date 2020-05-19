# 26 | 基于角色的权限控制：RBAC

## 笔记

### RBAC

`Kubernetes`中所有的`API`对象, 都保存在`Etcd`里. 对这些`API`对象的操作, 一定都是通过访问`kube-apiserver`实现的. 需要`APIServer`来帮助做授权工作.

**在 Kubernetes 中, 负责完成授权(`Authorization`)工作的机制, 就是 RBAC: 基于角色的访问控制(`Role-Based Access Control`)**

基本概念:

* `Role`: 角色, 一组规则, 定义了一组对`Kubernetes API`对象的操作权限.
* `Subject`: 被作用者, 可以是"人", 也可以是机器, 也可以是在`Kubernetes`里定义的"用户"
* `RoleBinding`: "被作用者"和"角色"的绑定关系

### Role

`Role`本身就是一个`Kubernetes`的`API`对象.

```
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: mynamespace
  name: example-role
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]
```

* `Role`指定了能产生作用的`Namespace`(`mynamespace`), `Namespace`是`Kubernetes`项目里的一个逻辑管理单位. 没有指定`Namespace`, 就是使用的是默认`Namsepace: default`.
* `rules`字段, 定义的权限规则. 示例中的含义是, 允许"被作用者"对`mynamespace`下面的`Pod`对象, 进行`GET`, `WATCH`和`LIST`操作.

具体的被作用者需要通过`RoleBinding`来实现.

`RoleBinding`也是一个`Kubernetes`的`API`对象.

```
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: example-rolebinding
  namespace: mynamespace
subjects:
- kind: User
  name: example-user
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: example-role
  apiGroup: rbac.authorization.k8s.io
```

* `subjects`: 被作用者
	* `User`: `Kubernetes`里的用户, `example-user`
* `roleRef`: `RoleBinding`对象可以直接通过名字, 引用我们前面定义的`Role`对象(example-role), 从而定义了**被作用者和角色之间的绑定关系**

`Role`和`RoleBinding`对象都是`Namespaced`对象(`Namespaced Object`), 她们对权限的限制规则仅在它们自己的`Namespace`内有效, `roleRef`也只能引用当前`Namespace`里的`Role`对象.

#### User

`Kubernetes`里的`User`是一个授权系统里的"逻辑概念", 需要通过外部认证服务, 如`Keystone`来提供, 或者可以给`APIServer`指定一个用户名, 密码文件.

### 对于非`Namespaced`对象(如`Node`), 或者一个`Role`要对作用于所有的`Namespace`, 如何授权

通过使用`ClusterRole`和`ClusterRoleBinding`这两个组合, 这两个定义里没有`Namespace`字段

```
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: example-clusterrole
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]
```

```
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: example-clusterrolebinding
subjects:
- kind: User
  name: example-user
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: example-clusterrole
  apiGroup: rbac.authorization.k8s.io
```

上面例子里的`ClusterRole`和`ClusterRoleBinding`的组合, 意味着名叫`example-user`的用户, 拥有对所有`Namespace`里的`Pod`进行`GET`, `WATCH`和`LIST`操作的权限.

所有权限(Kubernetes V1.11 里所有的操作):

```
verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
```

细化操作

```
# 该规则的被作用者, 只对名叫"my-config"的ConfigMap对象, 有进行 GET 操作的权限

rules:
- apiGroups: [""]
  resources: ["configmaps"]
  resourceNames: ["my-config"]
  verbs: ["get"]
```

### Kubernetes 的内置用户

**ServiceAccount**(内置权限), 给`ServiceAccount`分配权限.

#### 1. 定义一个 ServiceAccount

```
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: mynamespace
  name: example-sa
```

#### 2. RoleBinding

```
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: example-rolebinding
  namespace: mynamespace
subjects:
- kind: ServiceAccount
  name: example-sa
  namespace: mynamespace
roleRef:
  kind: Role
  name: example-role
  apiGroup: rbac.authorization.k8s.io
```

* `subjects`里的`kind`是一个名叫`example-sa`的`ServiceAccount`

#### 3. 创建

```
$ kubectl create -f svc-account.yaml
$ kubectl create -f role-binding.yaml
$ kubectl create -f role.yaml
```

```
$ kubectl get sa -n mynamespace -o yaml
- apiVersion: v1
  kind: ServiceAccount
  metadata:
    creationTimestamp: 2018-09-08T12:59:17Z
    name: example-sa
    namespace: mynamespace
    resourceVersion: "409327"
    ...
  secrets:
  - name: example-sa-token-vmfg6
```

`Kubernetes`会为一个`ServiceAccount`自动创建并分配一个`Secret`对象.

`secret`是用于`ServiceAccount`和`APIServer`进行交互的授权文件, 一般被称为`Token`. `Token`文件的内容一般是证书或者密码, 以一个`Secret`对象的方式保存在`Etcd`当中.

用户的`Pod`就可以声明使用这个`ServiceAccount`.

```
# 定义了 Pod 要使用的 ServiceAccount 的名字是 example-sa

apiVersion: v1
kind: Pod
metadata:
  namespace: mynamespace
  name: sa-token-test
spec:
  containers:
  - name: nginx
    image: nginx:1.7.9
  serviceAccountName: example-sa
```

运行后, 该`ServiceAccount`的`token`(`secret`对象), 被`Kubernetes`挂载到了容器的`/var/run/secrets/kubernetes.io/serviceaccount`目录下.

```
$ kubectl describe pod sa-token-test -n mynamespace
Name:               sa-token-test
Namespace:          mynamespace
...
Containers:
  nginx:
    ...
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from example-sa-token-vmfg6 (ro)
```

也可以通过`kubecl exec`查看目录里的文件

```
$ kubectl exec -it sa-token-test -n mynamespace -- /bin/bash
root@sa-token-test:/# ls /var/run/secrets/kubernetes.io/serviceaccount
ca.crt namespace  token
```

容器里的应用可以使用这个`ca.crt`来访问`APIServer`. 此时它只能做`GET`,`WATCH`和 `LIST`操作. 以为`example-sa`这个`ServiceAccount`的权限, 已经被我们绑定了`Role`做了限制.

在**第15篇中**, 如果一个`Pod`没有生命`serviceAccountName`, `Kubernetes`会自动在它的`Namespace`下创建一个名叫`default`的默认`ServiceAccount`, 然后分配个这个`Pod`.

这个默认的`ServiceAccount`并没有关联任何`Role`. 也就是说, 此时它有访问`APIServer`的绝大多数权限. 这个访问所需要的`Token`, 也是默认`ServiceAccount`对应的`Secret`对象为它提供的.

```
$kubectl describe sa default
Name:                default
Namespace:           default
Labels:              <none>
Annotations:         <none>
Image pull secrets:  <none>
Mountable secrets:   default-token-s8rbq
Tokens:              default-token-s8rbq
Events:              <none>

$ kubectl get secret
NAME                  TYPE                                  DATA      AGE
kubernetes.io/service-account-token   3         82d

$ kubectl describe secret default-token-s8rbq
Name:         default-token-s8rbq
Namespace:    default
Labels:       <none>
Annotations:  kubernetes.io/service-account.name=default
              kubernetes.io/service-account.uid=ffcb12b2-917f-11e8-abde-42010aa80002

Type:  kubernetes.io/service-account-token

Data
====
ca.crt:     1025 bytes
namespace:  7 bytes
token:      <TOKEN数据>
```

**生产环境中, 为所有`Namespace`下的默认`ServiceAccount`绑定一个只读权限的`Role`**.

### 用户组

`ServiceAccount`

```
system:serviceaccount:<Namespace名字>:<ServiceAccount名字>
```

用户组名字

```
system:serviceaccounts:<Namespace名字>
```

可以在`RoleBinding`里定义如下的`subjects`

```
# Role 的权限规则, 作用于 mynamespace 里的所有 ServiceAccount. 使用到了用户组的概念

subjects:
- kind: Group
  name: system:serviceaccounts:mynamespace
  apiGroup: rbac.authorization.k8s.io
```

`Kubernetes`中已经内置了很多为系统保留的`ClusterRole`, 名字都是以`system:`开头. 可以通过`kubectl get clusterroles`查看.

`system:kube-scheduler`的`ClusterRole`, 定义的权限规则是`kube-schedule`(Kubernetes 的调度组件)运行所需要的必要权限. 查看权限列表

```
$ kubectl describe clusterrole system:kube-scheduler
Name:         system:kube-scheduler
...
PolicyRule:
  Resources                    Non-Resource URLs Resource Names    Verbs
  ---------                    -----------------  --------------    -----
...
  services                     []                 []                [get list watch]
  replicasets.apps             []                 []                [get list watch]
  statefulsets.apps            []                 []                [get list watch]
  replicasets.extensions       []                 []                [get list watch]
  poddisruptionbudgets.policy  []                 []                [get list watch]
  pods/status                  []                 []                [patch update]
```

这个`system:kube-scheduler`的`ClusterRole`, 会被绑定给`kube-system Namesapce`下名叫`kube-scheduler`的`ServiceAccount`, 正式`Kubernetes`调度器的`Pod`声明使用的`ServiceAccount`.

### 预先定义的`ClusterRole`

1. `cluster-admin`, 最高权限`verbs=*`
2. `admin`
3. `edit`
4. `view`

```

$ kubectl describe clusterrole cluster-admin -n kube-system
Name:         cluster-admin
Labels:       kubernetes.io/bootstrapping=rbac-defaults
Annotations:  rbac.authorization.kubernetes.io/autoupdate=true
PolicyRule:
  Resources  Non-Resource URLs Resource Names  Verbs
  ---------  -----------------  --------------  -----
  *.*        []                 []              [*]
             [*]                []              [*]
```

## 扩展