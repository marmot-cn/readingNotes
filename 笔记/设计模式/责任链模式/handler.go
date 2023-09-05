package main

import "fmt"

// 我们创建了两个具体的处理者：ConcreteHandlerA 和 ConcreteHandlerB。当我们调用HandleRequest时，如果处理者A不能处理，它会将请求传递给处理者B。如果都不能处理，请求就不会被处理。

// Handler 定义处理者接口
type Handler interface {
	SetNext(Handler)
	HandleRequest(request string) bool
}

// ConcreteHandlerA 具体处理者A
type ConcreteHandlerA struct {
	next Handler
}

func (h *ConcreteHandlerA) SetNext(next Handler) {
	h.next = next
}

func (h *ConcreteHandlerA) HandleRequest(request string) bool {
	if request == "RequestA" {
		fmt.Println("HandlerA handles RequestA")
		return true
	}

	if h.next != nil {
		return h.next.HandleRequest(request)
	}

	return false
}

// ConcreteHandlerB 具体处理者B
type ConcreteHandlerB struct {
	next Handler
}

func (h *ConcreteHandlerB) SetNext(next Handler) {
	h.next = next
}

func (h *ConcreteHandlerB) HandleRequest(request string) bool {
	if request == "RequestB" {
		fmt.Println("HandlerB handles RequestB")
		return true
	}

	if h.next != nil {
		return h.next.HandleRequest(request)
	}

	return false
}

func main() {
	handlerA := &ConcreteHandlerA{}
	handlerB := &ConcreteHandlerB{}

	handlerA.SetNext(handlerB)

	handlerA.HandleRequest("RequestA")
	handlerA.HandleRequest("RequestB")
	handlerA.HandleRequest("RequestC")
}
