package main

import (
	"fmt"
	"os"
)

func main() {
	for index, arg := range os.Args[1:] { 
		fmt.Println(index, arg)
	}
}

//输出
//root@daa41eb01821:/go/01/homework# go run 2.go a b c
//0 a
//1 b
//2 c