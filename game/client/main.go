package main

import (
	"github.com/gorilla/websocket"
	"log"
	"math"
	"math/rand"
)

const wsServerEndpoint = "ws://localhost:40000"

type Login struct {
	ClientID int    `json:"clientID"`
	UserName string `json:"username"`
}

type GameClient struct {
	conn     *websocket.Conn
	clientID int
	username string
}

func newGameClient(conn *websocket.Conn, username string) *GameClient {
	return &GameClient{
		conn:     conn,
		clientID: rand.Intn(math.MaxInt),
		username: username,
	}
}

func (c *GameClient) login() error {
	return c.conn.WriteJSON(
		Login{
			ClientID: c.clientID,
			UserName: c.username,
		})
}

func main() {

	// Client使用WebSocket与Server进行通讯
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, _, err := dialer.Dial(wsServerEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 建立连接
	c := newGameClient(conn, "Tom")
	if err := c.login(); err != nil {
		log.Fatal(err)
	}
}
