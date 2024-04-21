package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func dowork(d time.Duration) {
	fmt.Println("doing work...")
	time.Sleep(d)
	fmt.Println("work is done!")
}

func doworkWg(d time.Duration, wg *sync.WaitGroup) {
	fmt.Println("doing work...")
	time.Sleep(d)
	fmt.Println("work is done!")
	wg.Done()
}

func doworkCh(d time.Duration, res chan string) {
	fmt.Println("doing work...")
	time.Sleep(d)
	res <- "work is done!"
}

func Test_WaitGroup(t *testing.T) {

	start := time.Now()

	// 顺序版本
	//dowork(time.Second * 2)
	//dowork(time.Second * 4)

	// waitGroup版本, 先进行Add, 需要用新goroutine进行执行的部分传入Wg并调用Done, 最后主线程需要调用Wait()等一下
	//wg := sync.WaitGroup{}
	//wg.Add(2)
	//go doworkWg(time.Second*2, &wg)
	//go doworkWg(time.Second*4, &wg)
	//wg.Wait()

	// channel版本
	//resultch := make(chan string)
	//go doworkCh(time.Second*2, resultch)
	//go doworkCh(time.Second*4, resultch)
	//
	//res1 := <-resultch
	//res2 := <-resultch
	//fmt.Println(res1)
	//fmt.Println(res2)
	//fmt.Printf("work took %v seconds\n", time.Since(start))

	// 复合版本
	resultch := make(chan string)
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go doworkCh(time.Second*2, resultch)
	go doworkCh(time.Second*4, resultch)
	go doworkCh(time.Second*6, resultch)

	go func() {
		for res := range resultch {
			fmt.Println(res)
			wg.Done()
		}
		// 如果最下面一行不加上 time.Sleep(), 下面这一行会来不及打印, 就因为上面done已经完成, 主线程也就结束退出
		fmt.Printf("work took %v seconds\n", time.Since(start))
	}()

	wg.Wait()
	close(resultch)
	time.Sleep(time.Second)
}
