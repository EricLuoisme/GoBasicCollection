package main

import (
	"fmt"
	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"net/http"
)

type GameServer struct {
	upgrader websocket.Upgrader
}

func newGameServer() actor.Receiver {
	return &GameServer{}
}

func (s *GameServer) Receive(ctx *actor.Context) {
	// 这里是根据Message()的类型进行switch, 内部的Message方法返回的是any
	switch msg := ctx.Message().(type) {
	case actor.Started:
		s.startHttp()
		_ = msg
	}
}

func (s *GameServer) startHttp() {
	fmt.Println("starting HTTP server on port 40000")
	go func() {
		http.HandleFunc("/ws", s.handleWS)
		http.ListenAndServe(":40000", nil)
	}()
}

func (s *GameServer) handleWS(w http.ResponseWriter, r *http.Request) {
	// 处理将http请求升级为websocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("ws upgrade err", err)
		return
	}
	fmt.Println("new client trying to connect")
	fmt.Println(conn)
}

func main() {
	config := actor.NewEngineConfig
	e, _ := actor.NewEngine(config())
	e.Spawn(newGameServer, "server")
	select {}
}
