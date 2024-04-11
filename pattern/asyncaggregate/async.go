package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	start := time.Now()

	userName := fetchUser()

	// 串行化
	//likes := fetchUserLikes(userName)
	//match := fetchUserMatch(userName)

	// 使用goroutine进行异步编排, 使用waitGroup进行等待
	// 由于我们用同一个channel进行消息存储, 如果这里初始化Channel的时候, 不带buffer,
	// 那么会导致我们后面goroutine向同一个channel投递数据时卡住, 只能一直等待, 从而引发runtime error
	respch := make(chan any, 2)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go fetchUserLikes(userName, respch, wg)
	go fetchUserMatch(userName, respch, wg)
	wg.Wait()
	close(respch)

	//fmt.Printf("%+v\n", likes)
	//fmt.Printf("%+v\n", match)

	// 使用for持续获取
	for resp := range respch {
		fmt.Printf("%+v\n", resp)
	}

	fmt.Println("took: ", time.Since(start))
}

func fetchUser() string {
	time.Sleep(time.Millisecond * 100)
	return "Bob"
}

func fetchUserLikes(userName string, respch chan any, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 150)
	respch <- 11
	wg.Done()
}

func fetchUserMatch(userName string, respch chan any, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 100)
	respch <- "Anna"
	wg.Done()
}
