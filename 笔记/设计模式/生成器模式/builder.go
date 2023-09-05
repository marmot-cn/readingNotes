package main

import "fmt"

// Product
type Car struct {
	Wheels    int
	Color     string
	TopSpeed  int
}

// Builder 
type CarBuilder interface {
	SetWheels() CarBuilder
	SetColor() CarBuilder
	Build() *Car
}

// ConcreteBuilder
type SportsCarBuilder struct {
	car *Car
}

func NewSportsCarBuilder() *SportsCarBuilder {
	return &SportsCarBuilder{&Car{}}
}

func (s *SportsCarBuilder) SetWheels() CarBuilder {
	s.car.Wheels = 4
	return s
}

func (s *SportsCarBuilder) SetColor() CarBuilder {
	s.car.Color = "Red"
	return s
}

func (s *SportsCarBuilder) Build() *Car {
	s.car.TopSpeed = 200
	return s.car
}

// Director
func NewCar(builder CarBuilder) *Car {
	builder.SetWheels().SetColor()
	return builder.Build()
}

func main() {
	builder := NewSportsCarBuilder()
	car := NewCar(builder)
	fmt.Printf("Car: %+v\n", car)
}
