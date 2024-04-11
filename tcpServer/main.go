package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{} // 使用 chan struct 可以使得内存使用为0, 目的是该chan只是作为signal的作用, 不涉及传输数据
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
	}
}

func (s *Server) Start() error {
	// 面向TCP的监听设置
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	fmt.Println("server start")
	// 额外goroutine开始接收调用
	go s.AcceptLoop()
	// 延后close
	defer ln.Close()
	s.ln = ln
	// 核心trick, 这个地方使用的是chan的阻塞, 如果不进行chan的退出(也就是向chan里面发送一个消息), 那么该goroutine就会卡住在这里, 相当于挂起
	<-s.quitch
	return nil
}

func (s *Server) AcceptLoop() {
	for {
		// 接受新建立的connection
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error", err)
			continue
		}
		// 处理connection -> 读
		fmt.Println("new connection to the server", conn.RemoteAddr())
		go s.ReadLoop(conn)
	}
}

func (s *Server) ReadLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		// 从connection中读
		read, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error", err)
			continue
		}
		// 处理读到的bytes, 以slice的方式获取
		msg := buf[:read]
		fmt.Println("server read", string(msg))
	}
}

func main() {
	// 使用命令 telnet localhost 9099 就可以建立连接进行通信
	server := NewServer(":9099")
	log.Fatalln(server.Start())
}
