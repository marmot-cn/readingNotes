# 24 | 深入解析声明式API（一）：API对象的奥秘

## 笔记

### 声明式 API 的设计

`API`对象在`Etcd`里的完整资源路径, 是由:

* `Group(API 组)`
* `Version(API 版本)`
* `Resource(API 资源类型)`

![](./img/24_01.png)

**Kubernetes 里 API 对象的祖师方式, 是层层递进的**.

```
apiVersion: batch/v2alpha1
kind: CronJob
...
```

* `CronJob`: 资源类型
* `batch`: 组
* `v2alpha1`: 版本

当我们提交了这个`YAML`文件之后, `Kubernetes`就会把这个`YAML`文件里描述的内容转换成`Kubernetes`里的一个`CronJob`对象.

### Kubernetes 如何对 API 对象解析

#### 1. Kubernetes 会匹配 API 对象的组

核心`API`对象, 如`Pod, Node`, 是不需要`Group`的. (它们的`Group`是""). `Kubernetes`会直接在`/api`这个层级进行下一步的匹配过程.

非核心对象, `Kubernetes`就必须在`/apis`这个层级里查找它对应的`Group`. 如`batch`(离线业务)下就包含`Job`和`CronJob`.

#### 2. Kubernetes 匹配到 API 对象的版本号

如上面的`CronJob`示例, `Kubernetes`在`batch`这个`Group`下, 匹配到的版本号就是`v2alpha1`.

在`Kubernetes`中, 同一种`API`对象可以有多个版本, 影响到用户的变更就可以通过升级新版本来处理, 保证向后兼容.

#### 3. Kubernetes 会匹配 API 对象的资源类型

匹配完版本之后, `Kubernetes`就知道, 要创建的原来是一个`/apis/batch/v2apha1`下的`CronJob`对象.

### 创建过程

![](./img/24_02.png)

* 发起了创建`CronJob`的`POST`请求之后, 编写的`YAML`的信息就被提交给了`APIServer`.
	* 过滤请求
	* 完成前置性工作, 授权, 超时处理, 审计等.
* 请求进入`MUX`和`Routes`流程(完成`URL`和`Handler`绑定的场所). 
* `APIServer`根据这个`CronJob`类型定义, 使用用户提交的`YAML`文件里的字段, 创建一个`CronJob`对象.
	* 把用户提交的`YAML`文件转换成一个`Super Version`的对象(`API`资源类型所有版本的字段全集), 用户提交不同版本的`YAML`文件, 都可以用这个`Super Version`对象来进行处理.
* `APIServer`进行`Admission()`和`Validation()`操作
	* `Validation`负责验证这个对象里各个字段是否合法. 被验证过的`API`对象, 保存在`APIServer`里的`Registry`的数据结构中. **只要在`Reistry`结构内的`API`对象, 都是有效的`Kubernetes API`对象**.
* `APIServer`把验证的`API`对象转换成用户最初提交的版本, 进行序列化操作, 用`Etcd`把`API`保存起来.

### CRD

`Custom Resource Definition`. 允许用户在`Kubernetes`中添加一个跟`Pod`, `Node`类似的, 新的`API`资源类型, **自定义`API`资源**.

#### 示例: 添加一个`Network`的`API`资源类型

`Network`对象的`YAML`文件**example-network.yaml**

```
apiVersion: samplecrd.k8s.io/v1
kind: Network
metadata:
  name: example-network
spec:
  cidr: "192.168.0.0/16"
  gateway: "192.168.0.1"
```

* `API`资源类型: `Network`
* `API`组: `samplecrd.k8s.io`
* `API`版本: `v1`

上面这个`YAML`文件, 是一个具体的"自定义API资源", 也叫`CR(Custom Resource)`. 为了能够让`Kubernetes`认识这个`CR`, 就需要让`Kubernetes`明白这个`CR`的宏观定义, `CRD(Custom Resource Definition)`.

`Network CRD`的`YAML`文件**network.yaml**

```
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: networks.samplecrd.k8s.io
spec:
  group: samplecrd.k8s.io
  version: v1
  names:
    kind: Network
    plural: networks
  scope: Namespaced
```

* `API`信息
	* `group`: `samplecrd.k8s.io`
	* `version`: `v1`
* `CR`的资源类型: `Network`
* 复数(plural)是`networks`
* `scope`是`Namespaced`, 我们定义的这个`Network`是一个属于`Namespace`的对象, 类似于`Pod`

这样`Kubernetes`就能够认识和处理所有声明了`API`类型是"`samplecrd.k8s.io/v1/network`"的`YAML`文件了.

还需要让`Kubernetes`认识`YAML`文件里描述的"网络"部分:

* `cidr`网段
* `gateway`网关

#### 1. 在`GOPATH`下, 创建项目

```
$ tree $GOPATH/src/github.com/<your-name>/k8s-controller-custom-resource
.
├── controller.go
├── crd
│   └── network.yaml
├── example
│   └── example-network.yaml
├── main.go
└── pkg
    └── apis
        └── samplecrd
            ├── register.go
            └── v1
                ├── doc.go
                ├── register.go
                └── types.go
```

* `pkg/apis/samplecrd`是`API`组的名字
* `v1`是版本
* `types.go`定义了`Network`对象的完整描述
* `register.go`放置全局变量

```
# register.go

package samplecrd

const (
 GroupName = "samplecrd.k8s.io"
 Version   = "v1"
)
```

* `doc.go`(`Golang`的文档源文件)

```
// +k8s:deepcopy-gen=package

// +groupName=samplecrd.k8s.io
package v1
```

`+<tag_name>[=value]`的代码格式注释风格, 是`Kubernetes`进行代码生成要用的`Annotation`风格的注释.

```
# types.go


package v1
...
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Network describes a Network resource
type Network struct {
 // TypeMeta is the metadata for the resource, like kind and apiversion
 metav1.TypeMeta `json:",inline"`
 // ObjectMeta contains the metadata for the particular object, including
 // things like...
 //  - name
 //  - namespace
 //  - self link
 //  - labels
 //  - ... etc ...
 metav1.ObjectMeta `json:"metadata,omitempty"`
 
 Spec networkspec `json:"spec"`
}
// networkspec is the spec for a Network resource
type networkspec struct {
 Cidr    string `json:"cidr"`
 Gateway string `json:"gateway"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkList is a list of Network resources
type NetworkList struct {
 metav1.TypeMeta `json:",inline"`
 metav1.ListMeta `json:"metadata"`
 
 Items []Network `json:"items"`
}
```

**register.go**, `Network`资源类型在服务端注册工作, `APIServer`会完成. `register.go`文件中的`addKnownTypes()`方法是让客户端也能知道`Network`资源类型的定义.

```
# register.go

package v1
...
// addKnownTypes adds our types to the API scheme by registering
// Network and NetworkList
func addKnownTypes(scheme *runtime.Scheme) error {
 scheme.AddKnownTypes(
  SchemeGroupVersion,
  &Network{},
  &NetworkList{},
 )
 
 // register the type in the scheme
 metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
 return nil
}
```

`Network`对象的定义:

* 自定义资源类型的`API`描述，组`Group`, 版本`Version`,资源类型`Resource`. **告诉计算机兔子是哺乳动物**.
* 自定义资源类型的对象描述, `Spec, Status`等. **告诉计算机,兔子有长二度和三瓣嘴**.

使用`Kubernetes`提供的代码生成工具, 为上面定义的`Network`资源类型自动生成`clientset, informer`和`lister`.

代码生成工具`k8s.io/code-generator`

```
# 代码生成的工作目录，也就是我们的项目路径
$ ROOT_PACKAGE="github.com/resouer/k8s-controller-custom-resource"
# API Group
$ CUSTOM_RESOURCE_NAME="samplecrd"
# API Version
$ CUSTOM_RESOURCE_VERSION="v1"

# 安装k8s.io/code-generator
$ go get -u k8s.io/code-generator/...
$ cd $GOPATH/src/k8s.io/code-generator

# 执行代码自动生成，其中pkg/client是生成目标目录，pkg/apis是类型定义目录
$ ./generate-groups.sh all "$ROOT_PACKAGE/pkg/client" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"
```

代码生成完之后的目录结构

```
$ tree
.
├── controller.go
├── crd
│   └── network.yaml
├── example
│   └── example-network.yaml
├── main.go
└── pkg
    ├── apis
    │   └── samplecrd
    │       ├── constants.go
    │       └── v1
    │           ├── doc.go
    │           ├── register.go
    │           ├── types.go
    │           └── zz_generated.deepcopy.go
    └── client
        ├── clientset
        ├── informers
        └── listers
```

`zz_generated.deepcopy.go`是自动生成的`DeepCopy`文件.

`clientset`, `informers` 和 `listers`是客户端库.

### 创建 `Network`类型的`API`对象

```
$ kubectl apply -f crd/network.yaml
customresourcedefinition.apiextensions.k8s.io/networks.samplecrd.k8s.io created


$ kubectl get crd
NAME                        CREATED AT
networks.samplecrd.k8s.io   2018-09-15T10:57:12Z


$ kubectl apply -f example/example-network.yaml 
network.samplecrd.k8s.io/example-network created


$ kubectl get network
NAME              AGE
example-network   8s


$ kubectl describe network example-network
Name:         example-network
Namespace:    default
Labels:       <none>
...API Version:  samplecrd.k8s.io/v1
Kind:         Network
Metadata:
  ...
  Generation:          1
  Resource Version:    468239
  ...
Spec:
  Cidr:     192.168.0.0/16
  Gateway:  192.168.0.1
```

## 扩展