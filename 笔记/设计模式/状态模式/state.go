package main

import "fmt"

// State interface
type State interface {
	Handle(context *Context)
}

// ConcreteStateA
type ConcreteStateA struct{}

func (s *ConcreteStateA) Handle(context *Context) {
	fmt.Println("Handling in ConcreteStateA")
	context.SetState(&ConcreteStateB{})
}

// ConcreteStateB
type ConcreteStateB struct{}

func (s *ConcreteStateB) Handle(context *Context) {
	fmt.Println("Handling in ConcreteStateB")
	context.SetState(&ConcreteStateA{})
}

// Context
type Context struct {
	state State
}

func NewContext(state State) *Context {
	return &Context{state: state}
}

func (c *Context) SetState(state State) {
	c.state = state
}

func (c *Context) Request() {
	c.state.Handle(c)
}

func main() {
	context := NewContext(&ConcreteStateA{})
	context.Request()  // Handling in ConcreteStateA
	context.Request()  // Handling in ConcreteStateB
	context.Request()  // Handling in ConcreteStateA
}

