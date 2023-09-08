package main

import (
	"fmt"
	"time"
)

// 模拟了NewComplexStruct的初始化过程是耗时的。而通过克隆得到的对象会复制原始对象的状态，但不会经历耗时的初始化过程。


// 一个模拟数据的复杂结构
type ComplexStruct struct {
	Data   map[string]map[int]float64
	Loaded time.Time
}

func NewComplexStruct() *ComplexStruct {
	// 模拟一个耗时的加载过程
	time.Sleep(2 * time.Second)
	data := make(map[string]map[int]float64)
	data["example"] = map[int]float64{1: 0.9, 2: 0.8}
	return &ComplexStruct{
		Data:   data,
		Loaded: time.Now(),
	}
}

func (cs *ComplexStruct) Clone() *ComplexStruct {
	// 这里进行了浅拷贝，如果你需要深拷贝，可以使用其他方法
	data := make(map[string]map[int]float64)
	for key, value := range cs.Data {
		data[key] = value
	}
	return &ComplexStruct{
		Data:   data,
		Loaded: cs.Loaded,
	}
}

func main() {
	start := time.Now()

	original := NewComplexStruct() // 模拟耗时的加载过程
	fmt.Println("Time taken to initialize:", time.Since(start))

	start = time.Now()
	clone := original.Clone() // 这个过程应该更快
	fmt.Println("Time taken to clone:", time.Since(start))

	fmt.Println("Original Loaded Time:", original.Loaded)
	fmt.Println("Clone Loaded Time:", clone.Loaded)
}