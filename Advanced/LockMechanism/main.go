package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 题目：
// 题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ： sync.Mutex 的使用、并发数据安全。
// 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ：原子操作、并发数据安全。

const (
	goroutines = 10
	increments = 1000
)

func main() {
	fmt.Println(couterWithMutex())
	fmt.Println(counterWithAtomic())
}

func couterWithMutex() int {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
		count int
	)

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				mutex.Lock()
				count++
				mutex.Unlock()
			}
		}()
	}
	wg.Wait()
	return count
}

func counterWithAtomic() int64 {
	var (
		wg    sync.WaitGroup
		count int64
	)

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				atomic.AddInt64(&count, 1)
			}
		}()
	}
	wg.Wait()
	return count
}
