package main

import (
	"fmt"
	"sync"
	"time"
)

func printOddNumbers(oddChan chan struct{}) {
	for i := 1; i <= 10; i++ {
		if i%2 != 0 {
			fmt.Printf("printOddNumbers : %d\n", i)
			// time.Sleep(time.Second)
		}
	}
	close(oddChan)
}

func printEvenNumbers(evenChan chan struct{}) {
	for i := 2; i <= 10; i++ {
		if i%2 == 0 {
			fmt.Printf("printEvenNumbers : %d\n", i)
			// time.Sleep(time.Second)
		}
	}
	close(evenChan)
}

// 设计一个任务调度器，接收一组任务（可以用函数表示），
// 并使用协程并发执行这些任务，同时统计每个任务的执行时间。
// 任务 计算1 到 n的和
func taskAdd1ton(n int) int {
	// return n * (n + 1) / 2
	var sum int = 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum
}

type TaskFunc func(i int) int
type Task struct {
	taskId   int
	taskFunc TaskFunc
	n        int
}

type runResult struct {
	taskId    int
	result    int
	startTime time.Time
	endTime   time.Time
	runTime   time.Duration
}

func TaskScheduler(tasks []Task, resChan chan runResult) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			// fmt.Println("task is ", task)
			startTime := time.Now()
			result := task.taskFunc(task.n)
			endTime := time.Now()
			resChan <- runResult{
				taskId:    task.taskId,
				result:    result,
				startTime: startTime,
				endTime:   endTime,
				runTime:   endTime.Sub(startTime),
			}
		}(task)
	}
	wg.Wait()
	close(resChan)
}

func main() {
	var results []runResult
	var resChan chan runResult = make(chan runResult, 3)
	TaskScheduler([]Task{
		{taskId: 1, taskFunc: taskAdd1ton, n: 10000},
		{taskId: 2, taskFunc: taskAdd1ton, n: 255000},
		{taskId: 3, taskFunc: taskAdd1ton, n: 90000},
	}, resChan)
	for {
		r, ok := <-resChan
		if !ok {
			break
		}
		results = append(results, r)
	}

	fmt.Println(results)

	// var oddChan chan struct{} = make(chan struct{})
	// var evenChan chan struct{} = make(chan struct{})
	// go printOddNumbers(oddChan)
	// go printEvenNumbers(evenChan)
	// <-oddChan
	// <-evenChan
	// fmt.Print("hello goroutine\n")
	// time.Sleep(10 * time.Second)

}
