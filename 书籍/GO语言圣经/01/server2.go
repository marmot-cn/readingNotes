// Server2 is a minimal "echo" and counter server. 

package main 

import ( 
	"fmt" 
	"log" 
	"net/http" 
	"sync" 
)

var mu sync.Mutex 
var count int 

func main() { 
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter) 
	log.Fatal(http.ListenAndServe("localhost:8000", nil)) 
}

// handler echoes the Path component of the requested URL. 
func handler(w http.ResponseWriter, r *http.Request) { 
	mu.Lock() 
	count++ 
	mu.Unlock() 
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path) 
}

// counter echoes the number of calls so far. 
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock() 
	fmt.Fprintf(w, "Count %d\n", count) 
	mu.Unlock() 
}

// root@daa41eb01821:/go/01# go run server2.go &
// [1] 6221
// root@daa41eb01821:/go/01# curl 127.0.0.1:8000/hello
// URL.Path = "/hello"
// root@daa41eb01821:/go/01# curl 127.0.0.1:8000/hello
// URL.Path = "/hello"
// root@daa41eb01821:/go/01# curl 127.0.0.1:8000/hello
// URL.Path = "/hello"
// root@daa41eb01821:/go/01# curl 127.0.0.1:8000/count
// Count 3