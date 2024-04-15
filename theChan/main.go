package main

import "fmt"

// main notes:
// - buffered & unbuffered channels
// - read & write protection
// - when to use them
func main() {

	// 如果使用unBuffered-channel, 这里会引发deadlock问题, 因为 userCh <- "Bob" 这一行就引发了block,
	// 因为没有缓存的channel, 这里需要一直等人consume传入的数据
	userCh := make(chan string)
	// 但如果这里使用另一个goroutine, 则不会发生block问题
	go func() {
		userCh <- "Bob"
	}()
	//userCh <- "Bob"
	user := <-userCh
	fmt.Println(user)

	// 如果使用buffered-channel, 下面情况也可以直接进行执行
	user2Ch := make(chan string, 2)
	user2Ch <- "Alice"
	user2 := <-user2Ch
	fmt.Println(user2)

	// 传输chan进行添加与获取
	sendMessage(user2Ch)
	consumeMessag(user2Ch)
}

func sendMessage(userCh chan string) {
	userCh <- "Charlie"
}

func consumeMessag(userCh chan string) {
	fmt.Println("from chan:", <-userCh)
}
