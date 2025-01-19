package main

import (
	"fmt"
	"unsafe"
)

type ExampleNoAlign struct {
	A bool // 1 byte
	B int64 // 8 bytes
	C int32 // 4 bytes
}

type ExampleWithAlign struct {
	B int64 // 8 bytes
	C int32 // 4 bytes
	A bool // 1 byte
}

func main() {
	noAlign := ExampleNoAlign{}
	withAlign := ExampleWithAlign{}

	fmt.Println("NoAlign Size:", unsafe.Sizeof(noAlign))
	fmt.Println("WithAlign Size:", unsafe.Sizeof(withAlign))
}


// ExampleNoAlign 可能会有额外的内存填充（padding），尤其是在 A 和 B 之间，以确保 B 对齐到 8 字节边界。这会导致整个结构体占用的内存比字段总和更大。
// ExampleWithAlign 通过将最大的字段（在这个例子中是 int64 类型的 B）放在前面，减少了因对齐产生的额外填充，从而可能有更小的总体内存占用。

// 输出
// NoAlign Size: 24
// WithAlign Size: 16