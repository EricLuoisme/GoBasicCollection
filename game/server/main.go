package main

import (
	"encoding/json"
	"fmt"
	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"github.com/roylic/gofolder/game/types"
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

func (ps *PlayerSession) Receive(ctx *actor.Context) {
	switch ctx.Message().(type) {
	case actor.Started:
		ps.readLoop()
	}
}

func (ps *PlayerSession) readLoop() {
	var msg types.WSMessage
	for {
		// 尝试将内容读到msg
		if err := ps.conn.ReadJSON(&msg); err != nil {
			fmt.Println("read error", err)
			return
		}
		// 异步处理
		go ps.handleMessage(msg)
	}
}

func (ps *PlayerSession) handleMessage(msg types.WSMessage) {
	switch msg.Type {
	case "login":
		var loginMsg types.Login
		if err := json.Unmarshal(msg.Data, &loginMsg); err != nil {
			panic(err)
		}
		fmt.Println(loginMsg)
	}
}

type GameServer struct {
	upgrader websocket.Upgrader
	ctx      *actor.Context
	sessions map[*actor.PID]struct{}
}

func newGameServer() actor.Receiver {
	return &GameServer{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		sessions: make(map[*actor.PID]struct{}),
	}
}

func (gs *GameServer) Receive(ctx *actor.Context) {
	// 这里是根据Message()的类型进行switch, 内部的Message方法返回的是any
	switch msg := ctx.Message().(type) {
	case actor.Started:
		gs.startHttp()
		gs.ctx = ctx
		_ = msg
	}
}

func (gs *GameServer) startHttp() {
	fmt.Println("starting HTTP server on port 40000")
	go func() {
		http.HandleFunc("/ws", gs.handleWS)
		http.ListenAndServe(":40000", nil)
	}()
}

func (gs *GameServer) handleWS(w http.ResponseWriter, r *http.Request) {
	// 处理将http请求升级为websocket
	conn, err := gs.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("ws upgrade err", err)
		return
	}
	fmt.Println("new client trying to connect")
	// 每次新client接入, spawn复制并传输状态
	sid := rand.Intn(math.MaxInt / 10000)
	pid := gs.ctx.SpawnChild(newPlayerSession(sid, conn), fmt.Sprintf("session_%d", sid))
	// 记录所有session, 进行内容广播
	gs.sessions[pid] = struct{}{}
	fmt.Printf("client with sid %d and pid %s connected\n", sid, pid)
}

func main() {
	e, _ := actor.NewEngine(actor.NewEngineConfig())
	e.Spawn(newGameServer, "server")
	select {}
}
