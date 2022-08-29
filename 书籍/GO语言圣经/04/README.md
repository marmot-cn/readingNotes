# README

## 笔记

### 4.1 数组

数组是一个由**固定**长度的特定类型元素组成的序列，一个数组可以由零个或多个元素组成。

和数组对应的类型是 Slice（切片），它是可以增长和收缩动态序列。

如果在数组的长度位置出现的是“...”省略号，则表示数组的长度是根据初始 化值的个数来计算。

```
var q [3]int = [3]int{1, 2, 3}

q := [...]int{1, 2, 3}
```

==数组的长度需要在编译时确定==

### 4.2 Slice

Slice（切片）代表变长的序列，序列中每个元素都有相同的类型。slice的语法和数组很像，只是没有固定长度而已。

slice的底层确实引用一个数组对象。一个slice由三个部分 构成：指针、长度和容量。指针指向第一个slice元素对应的底层数组元素的地址，要注意的 是slice的第一个元素并不一定就是数组的第一个元素。长度对应slice中元素的数目；长度不 能超过容量，容量一般是从slice的开始位置到底层数据的结尾位置。内置的len和cap函数分 别返回slice的长度和容量。

slice的切片操作s[i:j]，其中0 ≤ i≤ j≤ cap(s)，用于创建一个新的slice，引用s的从第i个元素开 始到第j-1个元素的子序列。新的slice将只有j-i个元素。如果i位置的索引被省略的话将使用0代 替，如果j位置的索引被省略的话将使用len(s)代替。

**make**

内置的make函数创建一个指定元素类型、长度和容量的slice。容量部分可以省略，在这种情 况下，容量将等于长度。

```
make([]T, len) 
make([]T, len, cap) // same as make([]T, cap)[:len]
```

### 4.3 Map

一个map就是一个哈希表的引用，map类型可以写为`map[K]V`，其中K和V分别 对应key和value。map中所有的key都有相同的类型，所有的value也有着相同的类型，但是 key和value之间可以是不同的数据类型。

**禁止对map元素取址**的原因是map可能随着元素数量的增长而重新分配更大的内存空间，从而 可能导致之前的地址无效。

**遍历的顺序是随机的**

### 4.4 结构体

考虑效率，结构体通常会用指针的方式传入和返回。

```
func Bonus(e *Employee, percent int) int { 
	return e.Salary * percent / 100 
}
```

如果要在函数内部修改结构体成员的话，用指针传入是必须的；在Go语言中，所有的函 数参数都是值拷贝传入的，函数参数将不再是函数调用时的原始变量。

```
//在Go语言中，对结构体进行&取地址操作时，视为对该类型进行一次 new 的实例化操作
pp := &Point{1, 2}
``` 

等价于 

```
pp := new(Point) 
*pp = Point{1, 2}
```

**new**

new(T)会为类型为T的新项分配已置零的内存空间，并返回它的地址。也就是一个类型为*T的值。在GO中就是它返回一个指针，该指针指向新分配的，类型为T的零值。

```
type SyncedBuffer struct {
    lock    sync.Mutex
    buffer  bytes.Buffer
}
 
// SyncedBuffer 类型的值在声明时就分配好内存就绪了，后续代码中p,v无需进一步处理即可正确工作
p := new(SyncedBuffer)  // type  *SyncedBuffer
var v SyncedBuffer      // type SyncedBuffer
```

### 4.5 JSON

将一个Go语言中 类似movies的结构体slice转为JSON的过程叫编组（marshaling）。编组通过调用 json.Marshal函数完成。

在编码时，默认使用Go语言结构体的成员名字作为JSON的对象（通过reflect反射技术）。只有导出的结构体成员才会被编码，这也就是我们为什么选择用大写字 母开头的成员名称。

编码的逆操作是解码，对应将JSON数据解码为Go语言的数据结构，Go语言中一般叫 unmarshaling，通过json.Unmarshal函数完成。

### 4.6 文本和HTML模板


