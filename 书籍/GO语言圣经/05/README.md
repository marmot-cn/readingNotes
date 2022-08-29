# README

## 笔记

### 5.1 函数声明

```
func name(parameter-list) (result-list) { 
	body 
}
```

没有函数体的函数声明，这表示该函数不是以Go实现的。这样的声明定义 了函数标识符。

```
package math 

func Sin(x float64) float //implemented in assembly language
```

实参包括引用类型，如

* 指针
* slice(切片)
* map
* function
* channel等类型

实参可 能会由于函数的间接引用被修改。

### 5.2 递归

### 5.3 多返回值

### 5.4 错误

### 5.5 函数值

### 5.6 匿名函数

### 5.7 可变参数

### 5.8 Deferred函数

当defer语句被执行时，跟在defer后面的函数会被延迟执行。直到包含该defer语句的函数执行完毕时， defer后的函数才会被执行，不论包含defer语句的函数是通过return正常结束，还是由于panic 导致的异常结束。你可以在一个函数中执行多条defer语句，它们的执行顺序与声明顺序相反。

defer语句经常被用于处理成对的操作，如打开、关闭、连接、断开连接、加锁、释放锁。通 过defer机制，不论函数逻辑多复杂，都能保证在任何执行路径下，资源被释放。释放资源的 defer应该直接跟在请求资源的语句后。

被延迟执行的匿名函数甚至可以修改函数返回给调用者的返回值：

```
func triple(x int) (result int) { 
defer func() { result += x }() 
return double(x) 
}

fmt.Println(triple(4)) // "12"
```

### 5.9 Panic异常

运行时错误会引起painc异常。不是所有的panic异常都来自运行时，直接调用内置的panic函数也会引发panic异常；panic函 数接受任何值作为参数。

此panic一般用于严重错误。

### 5.10 Recover捕获异常

如果在deferred函数中调用了内置函数recover，并且定义该defer语句的函数发生了panic异 常，recover会使程序从panic中恢复，并返回panic value。导致panic异常的函数不会继续运 行，但能正常返回。在未发生panic时调用recover，recover会返回nil。



