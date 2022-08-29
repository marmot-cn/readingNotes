# README

## 笔记

### 6.1 方法声明

在函数声明时，在其名字之前放上一个变量，即是一个方法。这个附加的参数会将该函数附 加到这种类型上，即相当于为这种类型定义了一个独占的方法。

```
package geometry 

import "math" 

type Point struct{ X, Y float64 } 

// traditional function 

func Distance(p, q Point) float64 { 
	return math.Hypot(q.X-p.X, q.Y-p.Y) 
}

// same thing, but as a method of the Point type 
func (p Point) Distance(q Point) float64 { 
	return math.Hypot(q.X-p.X, q.Y-p.Y) 
}
```

上面的代码里那个附加的参数p，叫做方法的接收器(receiver)，早期的面向对象语言留下的遗 产将调用一个方法称为“向一个对象发送消息”。

这种`p.Distance`的表达式叫做选择器，因为他会选择合适的对应p这个对象的Distance方法来 执行。

选择器也会被用来选择一个struct类型的字段。

### 6.2 基于指针对象的方法

* 不管你的method的receiver是指针类型还是非指针类型，都是可以通过指针/非指针类型 进行调用的，编译器会帮你做类型转换。
* 在声明一个method的receiver该是指针还是非指针类型时，你需要考虑两方面的内部，第 一方面是这个对象本身是不是特别大，如果声明为非指针变量时，调用会产生一次拷 贝；第二方面是如果你用指针类型作为receiver，那么你一定要注意，这种指针类型指向 的始终是一块内存地址，就算你对其进行了拷贝。熟悉C或者C艹的人这里应该很快能明 白。

### 6.3 通过嵌入结构体来扩展类型

### 6.4 方法值和方法表达式

当T是一个类型时，方法表达式可能会写作 T.f 或者 (*T).f ，会返回一个函数"值"，这种函 数会将其第一个参数用作接收器。

```
p := Point{1, 2} 
q := Point{4, 6} 

distance := Point.Distance // method expression 
fmt.Println(distance(p, q)) // "5"
 fmt.Printf("%T\n", distance) // "func(Point, Point) float64" 
 
 scale := (*Point).ScaleBy 
 scale(&p, 2) fmt.Println(p) // "{2 4}" 
 fmt.Printf("%T\n", scale) // "func(*Point, float64)" 
 
 // 译注：这个Distance实际上是指定了Point对象为接收器的一个方法func (p Point) Distance()， 
 // 但通过Point.Distance得到的函数需要比实际的Distance方法多一个参数， 
 // 即其需要用第一个额外参数指定接收器，后面排列Distance方法的参数。 
 // 看起来本书中函数和方法的区别是指有没有接收器，而不像其他语言那样是指有没有返回值。
```

### 6.5 示例:Bit数组

### 6.6 封装

