package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Args[0])
}

//执行
//root@daa41eb01821:/go/01/homework# go run 1.go
///tmp/go-build1696579466/b001/exe/1
//root@daa41eb01821:/go/01/homework# go build 1.go
//root@daa41eb01821:/go/01/homework# ./1
//./1