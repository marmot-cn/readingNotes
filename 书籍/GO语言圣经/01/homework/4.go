package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("count: %d\t string:%s\n", n, line)
		}
	}
}

func countLines(f *os.File, counts map[string]int) {

	input := bufio.NewScanner(f)

	for input.Scan() {

		if counts[input.Text()] > 0 {
			fmt.Printf("duplicate: %s\t , count is %d \t in file:%s\n",
				input.Text(), counts[input.Text()], f.Name())
		}
		counts[input.Text()]++
	}
	// NOTE: ignoring potential errors from input.Err()
}

//ctrl + d 终止
// 后面接文件，则从文件读取。否则从标准输入读取
//root@daa41eb01821:/go/01# go run dup2.go testfile
//count: 2	 string:111
