package main

import (
	"fmt"
)

type Printer interface {
	PrintWithPointer()
	PrintWithValue()
}

type MyStruct struct {
	data string
}

// 使用指针接收者
func (m *MyStruct) PrintWithPointer() {
	fmt.Println("Print with Pointer:", m.data)
}

// 使用值接收者
func (m MyStruct) PrintWithValue() {
	fmt.Println("Print with Value:", m.data)
}

func main() {
	s := MyStruct{data: "Hello"}

	// 定义一个接口变量
	var p Printer

	// 使用指针赋值
	p = &s
	p.PrintWithPointer()  // 正常工作
	p.PrintWithValue()    // 正常工作

	// 使用值赋值
	p = s
	p.PrintWithValue()    // 正常工作
	// p.PrintWithPointer()  // 这行会引发编译错误
}
