package main 

import ( 
	"fmt" 
	"os"
	"strconv"
)

func main() {

	a, _ := strconv.Atoi(os.Args[1])
	b, _ := strconv.Atoi(os.Args[2])

	g := gcd(a, b)
	fmt.Println()
	fmt.Printf("gcd is %d", g)
	fmt.Println()
}

func gcd(x, y int) int { 

	for y != 0 { 
		x, y = y, x%y 
	}

	return x 
}