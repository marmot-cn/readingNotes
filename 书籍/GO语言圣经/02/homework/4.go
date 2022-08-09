package main 

import ( 
	"fmt"
	"os"
	"strconv" 
)

// PopCount returns the population count (number of set bits) of x. 
func PopCount(x uint64) int { 

	sum := 0
	for i := x; i > 0; i = i/2 {
		sum += int(byte(i&1))
	}
	return sum 
}

func main() { 
	x,_ := strconv.Atoi(os.Args[1])
	fmt.Println(PopCount(uint64(x)))
}