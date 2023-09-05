package main

import "fmt"

// Subject Interface
type Subject interface {
	Request()
}

// RealSubject
type RealSubject struct{}

func (real *RealSubject) Request() {
	fmt.Println("RealSubject: Handling Request.")
}

// Proxy
type Proxy struct {
	realSubject *RealSubject
}

func (proxy *Proxy) Request() {
	if proxy.realSubject == nil {
		proxy.realSubject = &RealSubject{}
	}
	fmt.Println("Proxy: Intercepting and forwarding request to RealSubject.")
	proxy.realSubject.Request()
}

func main() {
	var subject Subject
	subject = &Proxy{}

	subject.Request()
}
