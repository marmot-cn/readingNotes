package main

import (
	"fmt"
)

// 当我们尝试获取具有相同键的蝇量对象时，FlyweightFactory确保我们得到的是已经存在的对象，而不是创建一个新的对象。这样，我们可以节省内存和其他资源。

// Flyweight 接口
type Flyweight interface {
	Operation(extrinsicState int)
}

// ConcreteFlyweight 结构体
type ConcreteFlyweight struct{}

func (cf *ConcreteFlyweight) Operation(extrinsicState int) {
	fmt.Println("ConcreteFlyweight state:", extrinsicState)
}

// FlyweightFactory 结构体
type FlyweightFactory struct {
	flyweights map[string]Flyweight
}

func NewFlyweightFactory() *FlyweightFactory {
	return &FlyweightFactory{
		flyweights: make(map[string]Flyweight),
	}
}

func (ff *FlyweightFactory) GetFlyweight(key string) Flyweight {
	if _, ok := ff.flyweights[key]; !ok {
		ff.flyweights[key] = &ConcreteFlyweight{}
	}
	return ff.flyweights[key]
}

func main() {
	factory := NewFlyweightFactory()

	flyweight1 := factory.GetFlyweight("A")
	flyweight1.Operation(1)

	flyweight2 := factory.GetFlyweight("A")
	flyweight2.Operation(2)

	flyweight3 := factory.GetFlyweight("B")
	flyweight3.Operation(3)

	// 输出 flyweights 的数量
	fmt.Println("Number of flyweights:", len(factory.flyweights))
}
