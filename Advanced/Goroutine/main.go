package main

import (
	"fmt"
	"sync"
	"time"
)

type Task func()

func main() {
	// var wg sync.WaitGroup
	// wg.Add(2)

	// go func() {
	// 	defer wg.Done()
	// 	printOdds()
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	printEvens()
	// }()

	// wg.Wait()

	tasks := []Task{
		func() { time.Sleep(2 * time.Second) },
		func() { time.Sleep(4 * time.Second) },
		func() { time.Sleep(6 * time.Second) },
	}
	durations := scheduler(tasks)
	for i, v := range durations {
		fmt.Printf("task %d took %d\n", i, v)
	}
}

func scheduler(task []Task) []time.Duration {
	var wg sync.WaitGroup
	wg.Add(len(task))
	duraions := make([]time.Duration, len(task))

	for i, v := range task {
		i, v := i, v //closures capture
		go func() {
			defer wg.Done()
			now := time.Now()
			v()
			duraions[i] = time.Since(now)
		}()
	}
	wg.Wait()
	return duraions
}

func printOdds() {
	for i := 1; i <= 10; i += 2 {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}

func printEvens() {
	for i := 2; i <= 10; i += 2 {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}
