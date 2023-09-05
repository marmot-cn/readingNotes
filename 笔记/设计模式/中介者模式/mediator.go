package main

import "fmt"

// ConcreteMediator知道每一个Colleague，并通过Send方法来协调它们之间的通信。Colleague对象与它的Mediator通信，而不是直接与其他Colleague通信。

// Mediator 接口
type Mediator interface {
	Send(message string, colleague Colleague)
}

// Colleague 接口
type Colleague interface {
	Send(message string)
	Notify(message string)
}

// ConcreteMediator 结构体
type ConcreteMediator struct {
	colleague1 Colleague
	colleague2 Colleague
}

func (cm *ConcreteMediator) Send(message string, colleague Colleague) {
	if colleague == cm.colleague1 {
		cm.colleague2.Notify(message)
	} else {
		cm.colleague1.Notify(message)
	}
}

// Colleague1 结构体
type Colleague1 struct {
	mediator Mediator
}

func (c *Colleague1) Send(message string) {
	c.mediator.Send(message, c)
}

func (c *Colleague1) Notify(message string) {
	fmt.Println("Colleague1 gets message:", message)
}

// Colleague2 结构体
type Colleague2 struct {
	mediator Mediator
}

func (c *Colleague2) Send(message string) {
	c.mediator.Send(message, c)
}

func (c *Colleague2) Notify(message string) {
	fmt.Println("Colleague2 gets message:", message)
}


// 我们首先创建了一个ConcreteMediator。
// 接着，我们创建了两个同事，Colleague1和Colleague2，并为它们都设置了相同的中介者。
// 最后，Colleague1发送了一条消息。这条消息不是直接发送给Colleague2的，而是发送给了中介者。中介者决定了如何处理这条消息（在这个例子中，它将消息发送给了Colleague2）。
func main() {
	m := &ConcreteMediator{}

	c1 := &Colleague1{mediator: m}
	c2 := &Colleague2{mediator: m}

	m.colleague1 = c1
	m.colleague2 = c2

	c1.Send("How are you, Colleague2?")
	c2.Send("Fine, thanks!")
}
