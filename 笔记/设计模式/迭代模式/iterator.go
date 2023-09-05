package main

import "fmt"

// Iterator interface
type Iterator interface {
	HasNext() bool
	Next() interface{}
}

// Aggregate interface
type Aggregate interface {
	Iterator() Iterator
}

// ConcreteAggregate
type ConcreteAggregate struct {
	items []string
}

func (ca *ConcreteAggregate) Iterator() Iterator {
	return &ConcreteIterator{aggregate: ca}
}

// ConcreteIterator
type ConcreteIterator struct {
	index     int
	aggregate *ConcreteAggregate
}

func (ci *ConcreteIterator) HasNext() bool {
	if ci.index < len(ci.aggregate.items) {
		return true
	}
	return false
}

func (ci *ConcreteIterator) Next() interface{} {
	if ci.HasNext() {
		item := ci.aggregate.items[ci.index]
		ci.index++
		return item
	}
	return nil
}

func main() {
	items := []string{"apple", "banana", "cherry"}
	aggregate := &ConcreteAggregate{items: items}
	iterator := aggregate.Iterator()

	for iterator.HasNext() {
		fmt.Println(iterator.Next())
	}
	// 输出:
	// apple
	// banana
	// cherry
}
