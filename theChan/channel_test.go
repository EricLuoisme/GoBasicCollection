package main

import (
	"fmt"
	"testing"
)

func TestAddUser(t *testing.T) {
	server := TheNewSever()
	server.Start()

	// 用10条协程进行投递信息
	// 但通过调用发现, 有时候会打印出Server添加channel的信息, 有时候却直接打印loop done然后结束
	for i := 0; i < 10; i++ {
		go func(i int) {
			server.userch <- fmt.Sprintf("user_%d", i)
		}(i)
	}
	fmt.Println("loop done")

	// 最后进行printing
	server.Printing()
}
