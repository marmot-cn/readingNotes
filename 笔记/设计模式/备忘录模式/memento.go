package main

import "fmt"

// Originator类可以创建一个备忘录，Memento类存储Originator的状态，而Caretaker类存储了所有的备忘录，允许用户将Originator的状态回滚到之前的状态。

// Memento 存储Originator的状态
type Memento struct {
	state string
}

func (m *Memento) getState() string {
	return m.state
}

// Originator 创建并存储Memento
type Originator struct {
	state string
}

func (o *Originator) setState(state string) {
	o.state = state
}

func (o *Originator) getState() string {
	return o.state
}

func (o *Originator) saveStateToMemento() *Memento {
	return &Memento{o.state}
}

func (o *Originator) getStateFromMemento(memento *Memento) {
	o.state = memento.getState()
}

// Caretaker 负责从Memento恢复Originator的状态
type Caretaker struct {
	mementoList []*Memento
}

func (c *Caretaker) add(memento *Memento) {
	c.mementoList = append(c.mementoList, memento)
}

func (c *Caretaker) get(index int) *Memento {
	return c.mementoList[index]
}

func main() {
	originator := &Originator{}
	caretaker := &Caretaker{}

	originator.setState("State #1")
	caretaker.add(originator.saveStateToMemento())

	originator.setState("State #2")
	caretaker.add(originator.saveStateToMemento())

	originator.setState("State #3")
	fmt.Println("Current State:", originator.getState())

	originator.getStateFromMemento(caretaker.get(0))
	fmt.Println("First saved State:", originator.getState())

	originator.getStateFromMemento(caretaker.get(1))
	fmt.Println("Second saved State:", originator.getState())
}
