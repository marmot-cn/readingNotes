# README

## 8.1 Goroutines

当一个程序启动时，其主函数即在一个单独的goroutine中运行，我们叫它main goroutine。主函数返回时，所有的goroutine都会被直接打断，程序退出。

新的goroutine会用go语句来创建。在语法上，go语句是一个普通的函数或方法调用前加上关键字`go`。

`spinner`和`fib`分别在独立的函数中，但两个函数会同时执行。

```
package main

import (
	"fmt"
	"time"
)

func main() { 
	go spinner(100 * time.Millisecond) 
	const n = 45 
	fibN := fib(n) // slow 
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN) 
}

func spinner(delay time.Duration) { 
	for {
		for _, r := range `-\|/` { 
			fmt.Printf("\r%c", r) 
			time.Sleep(delay) 
		} 
	} 
}

func fib(x int) int { 
	if x < 2 { 
		return x 
		}
	return fib(x-1) + fib(x-2) 
}
```

## 8.4 Channel

range循环可直接在 channels上面迭代。使用range循环是上面处理模式的简洁语法，它依次从channel接收数 据，当channel被关闭并且没有值可接收时跳出循环。

### 8.4.3 单方向的Channel

类型 chan<- int 表示一个只发送int的channel，只能发送不能接 收。相反，类型 <-chan int 表示一个只接收int的channel，只能接收不能发送。

### 8.4.4. 带缓存的Channels

**cap**

获取`channel`内部缓存的容量

```
fmt.Println(cap(ch))
```

**len**

将返回channel内部缓存队列中有效 元素的个数

```
fmt.Println(len(ch))
```

## 8.5 并发的循环

## 8.7 基于select的多路复用

```
select { 
	case <-ch1: 
		// ... 
	case x := <-ch2: 
		// ...use x... 
	case ch3 <- y: 
		// ... 
	default: 
	// ... 
}
```

select会等待case中有能够执行的case时去执行。当条件满足时，select才会去通信并执行 case之后的语句；**这时候其它通信是不会执行的**。

如果多个case同时就绪时，**select会随机地选择一个执行**，这样来保证每一个channel都有平 等的被select的机会。

在select语句中操作nil的channel永远都不会被select到。

### 8.9 并发的退出