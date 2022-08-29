package main 

import ( 
	"fmt" 
	"log" 
	"os" 
	"github" 
)

func main() { 
	result, err := github.SearchIssues(os.Args[1:]) 
	if err != nil { 
		log.Fatal(err) 
	}

	fmt.Printf("%d issues:\n", result.TotalCount) 
	for _, item := range result.Items { 
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title) 
	} 
}

//复制github/def.go 和 github/search.go 到 /usr/local/go/src/github/
//go run main.go repo:golang/go is:open json decoder
//go build main.go
// ./main repo:golang/go is:open json decoder