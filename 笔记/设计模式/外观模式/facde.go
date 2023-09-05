package main

import "fmt"

// 子系统1
type System1 struct{}

func (s *System1) Operation1() {
	fmt.Println("System1: Operation1 executed")
}

// 子系统2
type System2 struct{}

func (s *System2) Operation2() {
	fmt.Println("System2: Operation2 executed")
}

// 外观
type Facade struct {
	system1 *System1
	system2 *System2
}

func NewFacade() *Facade {
	return &Facade{
		system1: &System1{},
		system2: &System2{},
	}
}

func (f *Facade) Execute() {
	f.system1.Operation1()
	f.system2.Operation2()
}

func main() {
	facade := NewFacade()
	facade.Execute() // 输出: System1: Operation1 executed \n System2: Operation2 executed
}
