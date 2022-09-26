package main

import "fmt"

func main() { 

	naturals := make(chan int) 
	squares := make(chan int) 

	// Counter 
	go func() { 
		for x := 0; ; x++ { 
			naturals <- x 

			// if x == 5 {
			// 	close(naturals)
			// 	break
			// }
		} 
	}() 

	// Squarer 
	go func() { 
		for {
			x := <-naturals 
			squares <- x * x 
		} 
	}() 

	// Printer (in main goroutine) 
	for {
		fmt.Println(<-squares) 
	} 
}