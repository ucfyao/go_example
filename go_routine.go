package main

import (
	"fmt"
	"math/rand"
	"time"
)

func producer(header string, pipe chan string) {
	for {
		pipe <- fmt.Sprintf("%s: %v", header, rand.Int31())
		time.Sleep(time.Second)
	}
}

func consumer(pipe chan string) {
	for {
		message := <-pipe
		fmt.Println(message)
	}
}

func main() {
	channel := make(chan string)
	go producer("dog", channel)
	go producer("cat", channel)
	consumer(channel)
}
