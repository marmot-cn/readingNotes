package main 

import ( 
	"fmt"
	"crypto/sha256"
)

func main() { 

	sha1 := sha256.Sum256([]byte("x")) 
	sha2 := sha256.Sum256([]byte("X"))

	count := 0
	for i := 0; i < 32; i++ {
		if sha1[i] == sha2[i] {
			count++
		}
	}

	fmt.Println(count)
}