package main

import (
	"fmt"
	"sync"
)

func worker(workerId int, data <-chan int, wg *sync.WaitGroup) {
	for x := range data {
		fmt.Printf("worker %d received %d\n", workerId, x)
		wg.Done()
	}
}

func main() {
	data := make(chan int)
	qtdWorkers := 10
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < qtdWorkers; i++ {
		go worker(i, data, &wg)
	}

	for i := 0; i < 100; i++ {
		data <- i
	}
	wg.Wait()
}
