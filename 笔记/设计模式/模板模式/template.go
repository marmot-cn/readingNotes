package main

import "fmt"

// AbstractClass 定义算法的结构
type AbstractClass interface {
	PrimitiveOperation1()
	PrimitiveOperation2()
	TemplateMethod()
}

// ConcreteClass 实现抽象类定义的操作
type ConcreteClass struct{}

func (c *ConcreteClass) PrimitiveOperation1() {
	fmt.Println("ConcreteClass: Primitive Operation 1 executed")
}

func (c *ConcreteClass) PrimitiveOperation2() {
	fmt.Println("ConcreteClass: Primitive Operation 2 executed")
}

func (c *ConcreteClass) TemplateMethod() {
	fmt.Println("Algorithm Structure:")
	c.PrimitiveOperation1()
	c.PrimitiveOperation2()
}

func main() {
	template := &ConcreteClass{}
	template.TemplateMethod() 
	// 输出: 
	// Algorithm Structure:
	// ConcreteClass: Primitive Operation 1 executed
	// ConcreteClass: Primitive Operation 2 executed
}
