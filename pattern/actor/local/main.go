package main

import (
	"fmt"
	"net/rpc"
)

type Msg struct {
	To, From, Body string
}

// Actor 单个actor的结构与处理的方法
type Actor struct {
	Name string
	Msgs []Msg
}

func (a *Actor) SendMsg(m Msg) {
	a.Msgs = append(a.Msgs, m)
}

func (a *Actor) ProcessMsgs() {
	for _, m := range a.Msgs {
		fmt.Printf("Actor %s received msg: %s\n", a.Name, m.Body)
	}
	a.Msgs = nil
}

// ActorManager 用于处理actor之间的调用, 是否需要进行rpc远程调用
type ActorManager struct {
	Actors map[string]*Actor
}

func NewActorManager() *ActorManager {
	return &ActorManager{Actors: make(map[string]*Actor)}
}

func (m *ActorManager) RegisterActor(name string) {
	m.Actors[name] = &Actor{
		Name: name,
	}
}

func (m *ActorManager) SendMsgs(msg Msg) {
	if actor, ok := m.Actors[msg.To]; ok {
		// 如果找得到对应的Actor, 证明是本地, 直接进行发送
		actor.SendMsg(msg)
	} else {
		// 如果找不到, 认为是remote, 进行远程调用
		dial, err := rpc.Dial("tcp", msg.To)
		if err != nil {
			fmt.Printf("could not connect to remote actor actor: %s: %s\n", msg.To, err)
			return
		}

		// 	golang的rpc远程调用, dial.Call(远程方法名, 传输数据, 返回数据)
		defer dial.Close()
		var reply string
		err = dial.Call("RemoteActor.ReceiveMessage", msg, &reply)
		if err != nil {
			fmt.Printf("could not receive remote actor response: %s: %s\n", msg.To, err)
			return
		}
	}
}

func main() {
	manager := NewActorManager()

	manager.RegisterActor("actor1")
	manager.RegisterActor("actor2")

	// 发送给自己的actor
	manager.SendMsgs(
		Msg{
			To:   "actor1",
			From: "actor1",
			Body: "Hello from actor1, check local msg processing",
		})

	// 发送给远程
	manager.SendMsgs(
		Msg{
			To:   "127.0.0.1:9096",
			From: "actor1",
			Body: "Hello from remote actor1, check rpc msg processing",
		})

	// 发送实现了, 手动调用process本地的actor的数据处理
	for _, actor := range manager.Actors {
		actor.ProcessMsgs()
	}

}
