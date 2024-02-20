package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Message struct {
	id  int64
	Msg string
}

func main() {
	c1 := make(chan Message)
	c2 := make(chan Message)
	var i int64 = 0
	go func() {
		for {
			atomic.AddInt64(&i, 1)
			msg := Message{i, "hello from RabbitMQ"}
			time.Sleep(time.Second)
			c1 <- msg
		}
	}()

	go func() {
		for {
			i++
			atomic.AddInt64(&i, 1)
			time.Sleep(time.Second * 2)
			msg := Message{i, "Hello from kafka"}
			c2 <- msg
		}
	}()

	//for i := 0; i < 3; i++ {
	for {
		select {
		case msg1 := <-c1:
			fmt.Printf("received from rabbitMQ: id: %d %s\n", msg1.id, msg1.Msg)
		case msg2 := <-c2:
			fmt.Printf("received from kafka: Id:%d %s\n", msg2.id, msg2.Msg)
		case <-time.After(time.Second * 3):
			println("timeout")
			//default:
			//	println("default")
		}
	}
	//}
}
