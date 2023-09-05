package main

import "fmt"

// 在这个示例中，我们有一个DrawAPI接口（Implementor）和两个具体的实现：RedCircle和GreenCircle（Concrete Implementors）。我们的抽象部分是Shape，其中扩充的抽象是Circle。

// Implementor
type DrawAPI interface {
	DrawCircle(radius int, x int, y int)
}

// Concrete Implementor 1
type RedCircle struct{}

func (r *RedCircle) DrawCircle(radius int, x int, y int) {
	fmt.Printf("Drawing Circle[ color: red, radius: %d, x: %d, y: %d ]\n", radius, x, y)
}

// Concrete Implementor 2
type GreenCircle struct{}

func (g *GreenCircle) DrawCircle(radius int, x int, y int) {
	fmt.Printf("Drawing Circle[ color: green, radius: %d, x: %d, y: %d ]\n", radius, x, y)
}

// Abstraction
type Shape struct {
	drawAPI DrawAPI
}

func NewShape(d DrawAPI) *Shape {
	return &Shape{drawAPI: d}
}

// Refined Abstraction
type Circle struct {
	x, y, radius int
	*Shape
}

func NewCircle(x int, y int, radius int, drawAPI DrawAPI) *Circle {
	return &Circle{x: x, y: y, radius: radius, Shape: NewShape(drawAPI)}
}

func (c *Circle) Draw() {
	c.drawAPI.DrawCircle(c.radius, c.x, c.y)
}

func main() {
	redCircle := NewCircle(100, 100, 10, &RedCircle{})
	greenCircle := NewCircle(100, 100, 10, &GreenCircle{})

	redCircle.Draw()
	greenCircle.Draw()
}
