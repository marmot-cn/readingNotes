package main

import (
	"fmt"
	"net/http"
	"webapp/greeting"
)

func handler(writer http.ResponseWriter, request *http.Request) {
	greeting.Hello()
	fmt.Fprintf(writer, "Hello World, %s", request.URL.Path[1:])
}

func handler2(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello World 222")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/about", handler2)
	http.ListenAndServe(":8080", nil)
}