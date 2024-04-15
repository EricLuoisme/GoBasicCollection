package main

import (
	"fmt"
	"time"
)

type Server struct {
	users  map[string]string
	userch chan string
	quitch chan struct{} // 还是使用struct{}来作为卡住不允许shutdown的情况
}

func TheNewSever() *Server {
	return &Server{
		users:  make(map[string]string),
		userch: make(chan string),
		quitch: make(chan struct{}),
	}
}

func (s *Server) Start() {
	go s.loop()
}

// loop 方法类似单线程for(true), 持续的消费channel的内容, 没有内容的情况下进行阻塞
// 通过结合select, 使得同时处理2条channel的内容, 1) 需要consume的数据channel, 2) 需要关闭资源的quitChannel
func (s *Server) loop() {
	for {
		// 使用select来接收quitChannel的信号
		select {
		case msg := <-s.userch:
			fmt.Printf("adding new user %s\n", msg)
		case <-s.quitch:
			return
		}
	}
}

// Printing 用来遍历当前map存在的元素
func (s *Server) Printing() {
	// 这里的m是map的interface, 类似entrySet
	for _, m := range s.users {
		// 要进入到这一层, 才是key和value
		for k, v := range m {
			fmt.Println(k, "value is", v)
		}
	}
}

func (s *Server) addUser(user string) {
	s.users[user] = user
}

// main notes:
// - buffered & unbuffered channels
// - read & write protection
// - when to use them
func main() {

	//// 如果使用unBuffered-channel, 这里会引发deadlock问题, 因为 userCh <- "Bob" 这一行就引发了block,
	//// 因为没有缓存的channel, 这里需要一直等人consume传入的数据
	//userCh := make(chan string)
	//// 但如果这里使用另一个goroutine, 则不会发生block问题
	//go func() {
	//	userCh <- "Bob"
	//}()
	////userCh <- "Bob"
	//user := <-userCh
	//fmt.Println(user)
	//
	//// 如果使用buffered-channel, 下面情况也可以直接进行执行
	//user2Ch := make(chan string, 2)
	//user2Ch <- "Alice"
	//user2 := <-user2Ch
	//fmt.Println(user2)
	//
	//// 传输chan进行添加与获取
	//sendMessage(user2Ch)
	//consumeMessag(user2Ch)

	server := TheNewSever()
	server.Start()
	go func() {
		time.Sleep(time.Second * 3)
		close(server.quitch)
	}()

}

func sendMessage(userCh chan string) {
	userCh <- "Charlie"
}

func consumeMessag(userCh chan string) {

	// 使用Channel获取数据时, 可以使用2个返回参数的, 进行判断, 是否真的获取到数据(e.g. channel关闭这里会是false)
	user, ok := <-userCh
	if ok {
		fmt.Println("from chan:", user)
	}

}
