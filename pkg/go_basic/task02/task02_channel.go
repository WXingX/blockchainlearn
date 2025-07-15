package main

import (
	"fmt"
	"time"
)

// 写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，
// 并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。

func generateNum(intChan chan int) {
	for i := 1; i <= 10; i++ {
		intChan <- i
	}
	close(intChan)
}

func printNum(intChan chan int) {
	for v := range intChan {
		fmt.Printf("printNum: %d\n", v)
	}
}

// 实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，
// 消费者协程从通道中接收这些整数并打印。

func main04() {
	var intChan chan int = make(chan int)
	// 有缓冲channel
	// var intChan2 chan int = make(chan int, 100)
	go generateNum(intChan)
	go printNum(intChan)
	time.Sleep(5 * time.Second)
}
