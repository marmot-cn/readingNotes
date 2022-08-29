package main 

import (
	"fmt"
	"bytes"
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
	fmt.Println(comma("12345678911133344"))
}