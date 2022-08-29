package main 

import ( 
	"fmt" 
)

func add(x int, y int) int {return x + y} 
func sub(x, y int) (z int) { z = x - y; return} 
func first(x int, _ int) int { return x } 
func zero(int, int) int { return 0 } 


func main() { 
	fmt.Printf("%T\n", add)
	fmt.Printf("%T\n", sub)
	fmt.Printf("%T\n", first)
	fmt.Printf("%T\n", zero) 
}
