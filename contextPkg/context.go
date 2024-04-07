package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// main simulate web server
func main() {

	start := time.Now()
	ctx := context.Background()
	userId := 105

	data, err := fetchUserData(ctx, userId)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("result: ", data)
	fmt.Println("took: ", time.Since(start))
}

type Response struct {
	value int
	err   error
}

// fetchUserData 模拟获取用户数据
func fetchUserData(ctx context.Context, userID int) (int, error) {

	// 设置一个timeout的ctx, 并且延迟执行cancel
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	// 设置一个channel, 方便进行异步编排
	responseChan := make(chan Response)
	go func() {
		result, err := fetch3rdPartySlow()
		responseChan <- Response{
			value: result,
			err:   err,
		}
	}()

	// 这里单个事件可以使用select就能完成, 但使用for循环的意思是如果有多个事件, 可以进行相应处理
	// 但由于该function只返回 (int,error), 并且select里面包含了return, 实际上这个for是可以删除的
	//for {
	select {
	// 证明关闭了ctx, 就是超时被cancel这些情况
	case <-ctx.Done():
		return 0, fmt.Errorf("timeout on 3rd API")
	// 正常从channel获取得到结果的情况
	case resp := <-responseChan:
		return resp.value, nil
	}
	//}
}

// fetch3rdPartySlow 模拟调用缓慢的第三方API
func fetch3rdPartySlow() (int, error) {
	time.Sleep(time.Millisecond * 500)
	return 666, nil
}
