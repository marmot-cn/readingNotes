package main

import (
	"fmt"
)

// 定义了一个Prototype接口和两个实现该接口的ConcretePrototype类。每个ConcretePrototype类都有其自己的克隆方法来创建其实例的一个克隆。

// Prototype 接口声明了克隆方法
type Prototype interface {
	Clone() Prototype
}

// ConcretePrototype1 一个具体的原型类
type ConcretePrototype1 struct {
	name string
}

// Clone 实现Prototype接口的克隆方法
func (cp1 *ConcretePrototype1) Clone() Prototype {
	return &ConcretePrototype1{name: cp1.name}
}

// ConcretePrototype2 另一个具体的原型类
type ConcretePrototype2 struct {
	id int
}

// Clone 实现Prototype接口的克隆方法
func (cp2 *ConcretePrototype2) Clone() Prototype {
	return &ConcretePrototype2{id: cp2.id}
}

func main() {
	proto1 := &ConcretePrototype1{name: "Prototype1"}
	clone1 := proto1.Clone()

	proto2 := &ConcretePrototype2{id: 2021}
	clone2 := proto2.Clone()

	fmt.Printf("Original 1: %v, Clone 1: %v\n", proto1, clone1)
	fmt.Printf("Original 2: %v, Clone 2: %v\n", proto2, clone2)
}
