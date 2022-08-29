package main 

import ( 
	"fmt"
	"crypto/sha512"
	"crypto/sha256"
	"os"
)

func main() { 

	flag := os.Args[1] 

	if (flag == "384") {
		fmt.Printf("%x\n", sha512.Sum384([]byte("x")) )
	} else if (flag == "512") {
		fmt.Printf("%x\n", sha512.Sum512([]byte("x")) )
	} else {
		fmt.Printf("%x\n", sha256.Sum256([]byte("x")) )
	}
}