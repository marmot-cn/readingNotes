package main 

import ( 
	"fmt"
	"os"
	"strconv" 
)

// PopCount returns the population count (number of set bits) of x. 
func PopCount(x uint64) int { 
	sum := 0
	//直到x等于0
	for x != 0 {
		//逐步清0
		x = x&(x-1)
		sum++
	}
	return sum 
}

func main() { 
	x,_ := strconv.Atoi(os.Args[1])
	fmt.Println(PopCount(uint64(x)))
}