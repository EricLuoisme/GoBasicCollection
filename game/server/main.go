package main

import (
	"fmt"
	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"math"
	"math/rand"
	"net/http"
)

type PlayerSession struct {
	sessionID int
	clientID  int
	username  string
	inLobby   bool
	conn      *websocket.Conn
}

func newPlayerSession(sid int, conn *websocket.Conn) actor.Producer {
	return func() actor.Receiver {
		return &PlayerSession{
			sessionID: sid,
			conn:      conn,
		}
	}
}

func (p PlayerSession) Receive(context *actor.Context) {

}

type GameServer struct {
	upgrader websocket.Upgrader
	ctx      *actor.Context
}

func newGameServer() actor.Receiver {
	return &GameServer{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}}
}

func (s *GameServer) Receive(ctx *actor.Context) {
	// 这里是根据Message()的类型进行switch, 内部的Message方法返回的是any
	switch msg := ctx.Message().(type) {
	case actor.Started:
		s.startHttp()
		s.ctx = ctx
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
	// 每次新client接入, spawn复制并传输状态
	sid := rand.Intn(math.MaxInt / 10000)
	pid := s.ctx.SpawnChild(newPlayerSession(sid, conn), fmt.Sprintf("session_%d", sid))
	fmt.Printf("client with sid %d and pid %s connected\n", sid, pid)
}

func main() {
	e, _ := actor.NewEngine(actor.NewEngineConfig())
	e.Spawn(newGameServer, "server")
	select {}
}
