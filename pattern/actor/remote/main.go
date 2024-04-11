package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type Msg struct {
	To, From, Body string
}

// RemoteActor 处理远程调用actor时对面的内容
type RemoteActor struct {
}

func (a *RemoteActor) ReceiveMessage(m Msg, reply *string) error {
	fmt.Printf("remote actor %s receive msg: %s\n", m.To, m.Body)
	*reply = "Message received"
	return nil
}

func main() {
	actor := new(RemoteActor)
	rpc.Register(actor)
	listen, err := net.Listen("tcp", ":9096")
	if err != nil {
		fmt.Printf("erro listening", err)
		return
	}
	defer listen.Close()

	fmt.Println("remote actor start listening...")
	rpc.Accept(listen)
}
