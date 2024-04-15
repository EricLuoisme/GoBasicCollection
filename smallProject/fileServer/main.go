package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct {
}

func (f *FileServer) start() {

	// 监听TCP端口
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalln(err)
	}

	// 持续Accept连接
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		// 连接成功就给一个goroutine去持续读数据
		go f.readLoop(conn)
	}
}

func (f *FileServer) readLoop(conn net.Conn) {
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalln(err)
		}
		file := buf[:n]
		fmt.Println(file)
		fmt.Printf("received %d bytes over network\n", n)
	}
}

// sendFile 接收文件的方法调用
func sendFile(size int) error {

	// 读文件
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	// tcp连接
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return err
	}

	// 写文件流给Server
	n, err := conn.Write(file)
	if err != nil {
		return err
	}
	fmt.Printf("writtern %d bytes over the network\n", n)
	return nil
}

func main() {

	go func() {
		time.Sleep(4 * time.Second)
		sendFile(1000)
	}()

	server := &FileServer{}
	server.start()
}
