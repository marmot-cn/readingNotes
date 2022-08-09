package main 

import ( 
	"fmt"
	"os"
	"strconv" 
)

// pc[i] is the population count of i. 
var pc [256]byte 


func init() { 
	for i := range pc { 
		//byte(i&1) 是最后一位，pc[i/2]是排除最后一位的1的个数
		//如4是 100
		//100 & 1 最后一位是0, 的0
		// i/2 即 4/2 是 10 
		pc[i] = pc[i/2] + byte(i&1) 
	} 
}

// PopCount returns the population count (number of set bits) of x. 
func PopCount(x uint64) int { 
	return int(pc[byte(x>>(0*8))] + 
	pc[byte(x>>(1*8))] + 
	pc[byte(x>>(2*8))] + 
	pc[byte(x>>(3*8))] + 
	pc[byte(x>>(4*8))] + 
	pc[byte(x>>(5*8))] + 
	pc[byte(x>>(6*8))] + 
	pc[byte(x>>(7*8))]) 
}

func main() { 
	x,_ := strconv.Atoi(os.Args[1])
	fmt.Println(PopCount(uint64(x)))
}