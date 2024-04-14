package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"time"
)

// Server 一个websocket的持续feed项目
type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{conns: make(map[*websocket.Conn]bool)}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())
	// 确认链接启动
	s.conns[ws] = true
	// 持续loop读ws的数据
	s.readLoop(ws)
}

// readLoop 持续读取来自Websocket的数据
func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read err", err)
			continue
		}
		// 读取数据
		msg := buf[:n]
		fmt.Println(string(msg))
		// 返回数据
		//ws.Write([]byte("thank for your msg!\n"))
		// 广播数据
		s.broadcast(msg)
	}
}

// broadcast 广播数据到每一条建立的websocket连接
func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		// 启用goroutine进行写操作
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error", err)
			}
		}(ws)
	}
}

// handleWSOrderbook用来处理持续的WebSocket信息流
func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client to orderbook feed:", ws.RemoteAddr())
	for {
		payload := fmt.Sprintf("data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 2)
	}
}

// main 可通过浏览器命令进行websocket连接:
// let socket = new WebSocket("ws://localhost:9099/ws")
// socket.onmessage = (event) => {console.log("received from server: ", event.data)}
// socket.send("hello from client")
func main() {
	server := NewServer()
	// 注意这里用的是handle而不是handleFunc方法
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/feed", websocket.Handler(server.handleWSOrderbook))
	http.ListenAndServe(":9099", nil)
}
