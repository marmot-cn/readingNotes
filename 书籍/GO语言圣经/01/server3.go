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
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto) 
	for k, v := range r.Header { 
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v) 
	}

	fmt.Fprintf(w, "Host = %q\n", r.Host) 
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr) i
	if err := r.ParseForm(); err != nil { 
		log.Print(err) 
	}

	for k, v := range r.Form { 
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v) 
	}
}
