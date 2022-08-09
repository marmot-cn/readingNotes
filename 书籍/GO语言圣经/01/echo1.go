package main

import (
	"fmt"
	"os"
)

func main() {
	//定义了两个string类型的变量s和sep, 变量没有显示初始化，隐式地赋予其类型的零值
	//数值是0, 字符串是空
	var s, sep string

	for i := 1; i < len(os.Args); i++ {
		//+ 链接字符串
		s += sep + os.Args[i]
		sep = " "
	}

	fmt.Println(s)
}

//输出
//root@daa41eb01821:/go/01# go run echo1.go a b c
//a b c