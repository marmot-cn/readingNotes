package main 

import ( 
	"fmt" 
	"io/ioutil" 
	"os" 
	"strings" 
)

func main() { 
	counts := make(map[string]int) 
	for _, filename := range os.Args[1:] { 
		data, err := ioutil.ReadFile(filename) 
		if err != nil { 
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err) 
			continue 
		}
		for _, line := range strings.Split(string(data), "\n") { 
			counts[line]++ 
		} 
	}
	for line, n := range counts { 
		if n > 1 { 
			fmt.Printf("count: %d\t string:%s\n", n, line)
		} 
	} 
}

//root@daa41eb01821:/go/01# go run dup3.go testfile
//count: 2	 string:111