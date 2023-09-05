package main

import "fmt"

// Target 是客户端期望的接口
type Target interface {
	Request() string
}

// Adaptee 是需要被适配的类
type Adaptee struct{}

func (a *Adaptee) SpecificRequest() string {
	return "Specific request from Adaptee"
}

// Adapter 使得 Adaptee 能够与 Target 接口一起工作
type Adapter struct {
	adaptee *Adaptee
}

func NewAdapter(adaptee *Adaptee) *Adapter {
	return &Adapter{adaptee: adaptee}
}

func (a *Adapter) Request() string {
	return a.adaptee.SpecificRequest()
}

func main() {
	adaptee := &Adaptee{}
	adapter := NewAdapter(adaptee)
	fmt.Println(adapter.Request()) // 输出: Specific request from Adaptee
}
