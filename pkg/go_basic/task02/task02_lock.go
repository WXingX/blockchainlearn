package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
// 启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func main05() {
	var Counter int = 0;
	var lock sync.Mutex;
	for i := 1; i <= 10; i++ {
		go func(){
			for j := 1; j <= 1000; j++ {
				lock.Lock()
				Counter++
				lock.Unlock()
			}
		}()
	}
	time.Sleep(time.Second * 10)
	fmt.Println(Counter)
}

// 使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
// 启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。

func main06() {
	// main05()
	var Counter int64 = 0
	for i := 1; i <= 10;i++{
		go func() {
			for j := 1; j <= 1000; j++ {
				atomic.AddInt64(&Counter, 1)
			}
		}()
	}
	time.Sleep(time.Second * 10)
	fmt.Println(Counter)
}