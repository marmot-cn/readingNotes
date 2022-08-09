# README

## 笔记

### 1.1 

`main`包比较特殊。它定义了一个独立可执行的程序，而不是一个库。

`main`函数是整个程序执行时的入口。

`import`声明编译器源文件需要哪些包。

关键字

* 函数: func
* 变量: var
* 常量: const
* 类型: type

Go语言不需要在语句或者声明的末尾添加分号，除非一行上有多条语句。

gofmt 统一代码格式规范。

### 1.2

**等价while循环**

```
for condition {
}
```

**等价无限循环**

```
for {
}
```

### 1.4

常量声明和变量声明一般都会出现在包级别，所以这些 常量在整个包中都是可以共享的，或者你也可以把常量声明定义在函数体内部，那么这种常 量就只能在函数体内用。

struct是一组值或者叫字段的集合，不同的类型集合 在一个struct可以让我们以一个统一的单元进行处理

### 1.6

goroutine是一种函数的并发执行方式，而channel是用来在goroutine之间进行参数传递。

main函数本身也运行在一个goroutine中

go function则表示创建一个新的goroutine，并在 这个新的goroutine中执行这个函数。



## 代码

### 1.1 

#### 代理

```
go env -w GOPROXY="https://goproxy.cn"
```

##### 安装`goimports`

```
go install golang.org/x/tools/cmd/goimports@latest
```

#### 执行

```
[ansible@k8s-agent-1 go]$ docker exec -it go-env /bin/bash
root@daa41eb01821:/go# cd helloworld/
root@daa41eb01821:/go/helloworld# go run helloworld.go
Hello, 世界

```

#### 编译

```
root@daa41eb01821:/go/helloworld# go build helloworld.go
root@daa41eb01821:/go/helloworld# ./helloworld
Hello, 世界
```

#### gofmt

格式化代码

```
gofmt -w helloworld.go
```

#### go doc

`go doc 函数`, 可以快速读取文档

```
go doc strconv.ParseFloat
```