package main

import (
	"fmt"
	"sync"
)

// 小练习，正确使用waitGroup并发调用goroutine
func main() {

	ch := make(chan int)

	var waitGroup sync.WaitGroup

	feeds := []int{1, 2, 3, 4, 5}
	waitGroup.Add(len(feeds))

	for _, feed := range feeds {
		go func(feed int) {
			fmt.Println("send", feed)
			ch <- feed
			waitGroup.Done()
		}(feed)
	}

	go func() {
		fmt.Println("wait")
		waitGroup.Wait()
		close(ch)
	}()

	for data := range ch {
		fmt.Println("receive", data)
	}
}
