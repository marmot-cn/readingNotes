package main 

import (
	"fmt"
)

func compare(oString string, cString string) bool {

	if (len(oString) != len(cString)) {
		return false
	}

	for _, o := range oString {
		equal := false
		for _, c := range cString {
			if ( o == c ) {
				equal = true
			}
		}

		if (!equal) {
			return false
		}
	}

	return true
}

func main() { 
	fmt.Println(compare("abc", "cba"))
	fmt.Println(compare("abc", "cbaa"))
	fmt.Println(compare("abc", "abe"))
}