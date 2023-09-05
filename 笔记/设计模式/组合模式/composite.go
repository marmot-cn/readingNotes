package main

import "fmt"

// Component interface
type Component interface {
	Operation() string
}

// Leaf
type Leaf struct {
	Name string
}

func (l *Leaf) Operation() string {
	return l.Name
}

// Composite
type Composite struct {
	Children []Component
}

func (c *Composite) Operation() string {
	result := "Composite: ["
	for _, child := range c.Children {
		result += child.Operation() + " "
	}
	result += "]"
	return result
}

func (c *Composite) Add(child Component) {
	c.Children = append(c.Children, child)
}

func (c *Composite) Remove(child Component) {
	for i, component := range c.Children {
		if component == child {
			c.Children = append(c.Children[:i], c.Children[i+1:]...)
			break
		}
	}
}

func main() {
	leaf1 := &Leaf{Name: "Leaf1"}
	leaf2 := &Leaf{Name: "Leaf2"}

	composite := &Composite{}
	composite.Add(leaf1)
	composite.Add(leaf2)

	fmt.Println(leaf1.Operation())
	fmt.Println(composite.Operation())
}

// 输出:
// Leaf1
// Composite: [Leaf1 Leaf2 ]
