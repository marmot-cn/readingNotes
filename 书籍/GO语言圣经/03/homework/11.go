package main 

import (
	"fmt"
	"bytes"
	"strings"
)

func comma(s string) string { 

	var buf bytes.Buffer

	for i,v := range s {

		fmt.Fprintf(&buf, "%v", string(v))

		if (i == len(s)-1) {
			break
		}

		if i%3 == 2 && i>0 {
			buf.WriteByte(',')
		}
	}

	return buf.String()
}

func main() { 

	num := "123456789111333444.321312412"

	index := strings.LastIndex(num, ".")

	fmt.Println(comma(num[:index]) + "." + num[index+1:])
}