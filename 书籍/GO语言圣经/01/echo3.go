package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}

//输出
//root@daa41eb01821:/go/01# go run echo3.go a b c
//a b c