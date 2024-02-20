package main

import "fmt"

func main() {
	canal := make(chan string)

	go func() {
		canal <- "Olá Mundo!"
	}()

	msg := <-canal
	fmt.Println(msg)
}
