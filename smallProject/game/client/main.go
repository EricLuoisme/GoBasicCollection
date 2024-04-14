package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/roylic/gofolder/game/types"
	"log"
	"math"
	"math/rand"
	"time"
)

const wsServerEndpoint = "ws://localhost:40000/ws"

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

	// 规范为json格式并转换为bytes数组传输
	b, err := json.Marshal(
		types.Login{
			ClientID: c.clientID,
			UserName: c.username,
		})
	if err != nil {
		return err
	}
	// 具体types传入, 向对面传输JSON格式内容(其中data是bytes数组, 获取之后又要转json)
	return c.conn.WriteJSON(
		types.WSMessage{
			Type: "login",
			Data: b,
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
	// looping
	for {
		// 持续发送单条信息
		state := types.PlayerState{
			Health: 100,
			Pos: types.Position{
				X: rand.Intn(10),
				Y: rand.Intn(10),
			},
		}
		b, err := json.Marshal(state)
		if err != nil {
			panic(err)
		}
		// 二次序列化
		msg := types.WSMessage{
			Type: "playerstate",
			Data: b,
		}
		if err := conn.WriteJSON(msg); err != nil {
			panic(err)
		}
		// 随后等待
		time.Sleep(time.Millisecond * 200)
	}
}
