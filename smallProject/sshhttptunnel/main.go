package main

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

type Tunnel struct {
	w      io.Writer
	donech chan struct{}
}

var tunnels = map[int]chan Tunnel{}

// main: 将SSH连接的与Http请求进行捆绑, 传输内容
func main() {

	go func() {
		http.HandleFunc("/", handleRequest)
		log.Fatalln(http.ListenAndServe(":3000", nil))
	}()

	// 处理SSH请求
	ssh.Handle(func(s ssh.Session) {
		id := rand.Intn(math.MaxInt)
		tunnels[id] = make(chan Tunnel)

		fmt.Println("tunnel ID ->", id)

		tunnel := <-tunnels[id] // 开始阻塞channel
		fmt.Println("tunnel is ready")

		_, err := io.Copy(tunnel.w, s)
		if err != nil {
			log.Fatalln(err)
		}

		// 关闭阻塞channel
		close(tunnel.donech)
		s.Write([]byte("channel closed"))
	})

	// 开启serve端口后, 这里就会被阻塞
	log.Fatalln(ssh.ListenAndServe(":2222", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	tunnel, ok := tunnels[id]
	if !ok {
		w.Write([]byte("tunnel not found"))
		return
	}

	donech := make(chan struct{})
	tunnel <- Tunnel{
		w:      w,
		donech: donech,
	}

	// 阻塞
	<-donech
}
