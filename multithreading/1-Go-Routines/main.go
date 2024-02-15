package main

import (
	"fmt"
	"sync"
	"time"
)

func task(name string) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%d: Task %s is running \n", i, name)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(25)

	go task("A")
	go task("B")
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("%d: Task %s is running \n", i, "anonymous")
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(15 * time.Second)
}
