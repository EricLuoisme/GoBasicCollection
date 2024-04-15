package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
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

// readLoopFull 固定的size会导致Client传输的文件大小超出这个size时, Server会丢失部分内容
//func (f *FileServer) readLoopFull(conn net.Conn) {
//	buf := make([]byte, 2048)
//	for {
//		n, err := conn.Read(buf)
//		if err != nil {
//			log.Fatalln(err)
//		}
//		file := buf[:n]
//		fmt.Println(file)
//		fmt.Printf("received %d bytes over network\n", n)
//	}
//}

// readLoop 使用buffer进行缓存
func (f *FileServer) readLoop(conn net.Conn) {
	buf := new(bytes.Buffer)
	for {
		// 使用小端编码的特性, 将大小提取出来
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)

		// 使用copy来将connection读取到的字节copy到buffer中
		// 但需要注意, Copy会直到EOF, 如果没有send, 那么会一直阻塞
		// 这里也需要copyN
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(buf.Bytes())
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

	// 写文件流给Server, 会产生问题
	//n, err := conn.Write(file)

	// 使用小端编码, 使得size在前面, 让接收端知道要用多少bytes
	binary.Write(conn, binary.LittleEndian, int64(size))

	// copy文件为stream, 这里使用CopyN可以确认一定只传多少字节的数据, 从而给出一个EOF到最后
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		return err
	}
	fmt.Printf("writtern %d bytes over the network\n", n)
	return nil
}

func main() {

	go func() {
		time.Sleep(4 * time.Second)
		sendFile(4000) // 如果这里传入的bytes超出了server限制的大小, 会导致丢失, 所以最佳实践应该将file转换为流而不是直接传一个bytes数组
	}()

	server := &FileServer{}
	server.start()
}
