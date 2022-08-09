package main

import (
	"fmt"
	"os"
)

func main() {
	//等价
	// s := "", 只能用在函数内部
	//var s string 
	//var s = "" 
	//var s string = ""

	s, sep := "", "" 

	//每次循环迭代， range 产生一对值；索引以及在该索引处的元素值
	//空标识符 _ (下划线)，接收索引
	for _, arg := range os.Args[1:] { 
		s += sep + arg 
		sep = " " 
	}
	fmt.Println(s)
}

//输出
//root@daa41eb01821:/go/01# go run echo2.go a b c
//a b c