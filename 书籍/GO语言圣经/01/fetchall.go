package main 

import ( 
	"fmt" 
	"io" 
	"io/ioutil" 
	"net/http" 
	"os" 
	"time" 
)

func main() { 
	start := time.Now() 
	ch := make(chan string) //make函数创建了一个传递string类型参数的channel
	for _, url := range os.Args[1:] { 
		go fetch(url, ch) // start a goroutine 
	}

	//通道 ch 是可以进行遍历的，遍历的结果就是接收到的数据。数据类型就是通道的数据类型。
	for range os.Args[1:] { 
		fmt.Println(<-ch) // receive from channel ch 
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds()) 
}

func fetch(url string, ch chan<- string) { 
	start := time.Now() 
	resp, err := http.Get(url) 
	if err != nil { 
		ch <- fmt.Sprint(err) // send to channel ch 
		return 
	}

	//ioutil.Discard 可以把 这个变量看作一个垃圾桶，可以向里面写一些不需要的数据
	nbytes, err := io.Copy(ioutil.Discard, resp.Body) 
	resp.Body.Close() // don't leak resources 
	if err != nil { 
		ch <- fmt.Sprintf("while reading %s: %v", url, err) 
		return 
	}

	secs := time.Since(start).Seconds() 
	//每当请求返回内容时，fetch函数都会往ch这个channel里写 入一个字符串
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url) 
}