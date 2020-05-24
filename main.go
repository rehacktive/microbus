package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/rehacktive/microbus/microbus"
)

// pubsub like reactive programming?

type message struct {
	content string
}

func main() {
	mb := microbus.New()

	subs1 := mb.Subscribe("test")

	go func() {
		for msg := range subs1 {
			fmt.Println("subs1: ", msg)
		}
	}()

	subs2 := mb.Subscribe("test")

	go func() {
		for msg := range subs2 {
			fmt.Println("subs2: ", msg)
		}
	}()

	go func() {
		for i := 1; i <= 5; i++ {
			mb.Publish("test", message{content: "hello worldz"})
			time.Sleep(2 * time.Second)
		}
		mb.CloseAll()
		mb.Publish("test", message{content: "hello worldz"})
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
}
