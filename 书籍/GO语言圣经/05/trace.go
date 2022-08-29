package main 

import ( 
	"log" 
	"time"
)


func bigSlowOperation() { 
	defer trace("bigSlowOperation")() // 这里是先执行了调用函数的部分，然后defer了匿名函数

	time.Sleep(10 * time.Second) // simulate slow 
}

func trace(msg string) func() { 
	start := time.Now() 
	log.Printf("enter %s", msg) 
	return func() { 
		log.Printf("exit %s (%s)", msg,time.Since(start)) 
	} 
}

func main() { 
	bigSlowOperation()
}