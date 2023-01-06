package main

import "fmt"

type matcher interface {
	HelloWorld()
}

type match struct {
	value bool
}

func (m *match) HelloWorld() {
	fmt.Println("hello world")
}

func main() {
	var dm match
	var ma matcher = &dm

	ma.HelloWorld()
	dm.HelloWorld()
}
