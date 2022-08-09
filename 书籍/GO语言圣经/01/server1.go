// Server1 is a minimal "echo" server. 
package main 

import ( 
	"fmt" 
	"log" 
	"net/http" 
)

func main() { 
	http.HandleFunc("/", handler) // each request calls handler 
	log.Fatal(http.ListenAndServe("localhost:8000", nil)) 

}

// handler echoes the Path component of the request URL r. 
func handler(w http.ResponseWriter, r *http.Request) { 
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path) 
}

//root@daa41eb01821:/go/01# go run server1.go &
//[1] 6137
//root@daa41eb01821:/go/01# curl 127.0.0.1:8000/hello
//URL.Path = "/hello"