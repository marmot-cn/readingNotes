package main

import "fmt"

// Command 接口声明执行操作的方法
type Command interface {
	Execute()
}

// Receiver 执行与请求关联的操作
type Receiver struct{}

func (r *Receiver) Action() {
	fmt.Println("Receiver action triggered!")
}

// ConcreteCommand 将一个接收者对象与一个动作绑定，调用接收者相应的方法来实现 Execute
type ConcreteCommand struct {
	receiver *Receiver
}

func NewConcreteCommand(receiver *Receiver) *ConcreteCommand {
	return &ConcreteCommand{receiver: receiver}
}

func (c *ConcreteCommand) Execute() {
	c.receiver.Action()
}

// Invoker 要求命令执行一个请求
type Invoker struct {
	command Command
}

func (i *Invoker) SetCommand(command Command) {
	i.command = command
}

func (i *Invoker) ExecuteCommand() {
	i.command.Execute()
}

func main() {
	// 客户端代码
	receiver := &Receiver{}
	command := NewConcreteCommand(receiver)
	invoker := &Invoker{}

	invoker.SetCommand(command)
	invoker.ExecuteCommand()
}
