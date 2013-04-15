package redis

import (
	"io"
	"testing"
	"time"
)

func TestChannels(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}
	defer r.Close()

	messages := make(chan string)

	sender := r.Channel("Test_Channel")
	receiver1 := r.Channel("Test_Channel")
	receiver2 := r.Channel("Test_Channel")
	receiver3 := r.Channel("*Channel")
	receiver4 := r.Channel("Test*")

	var closer1 io.Closer
	start1, closer1 := receiver1.Subscribe(func(message string) {
		if message != "Test Test" {
			t.Error("receiver1 didn't get correct message (Got \"", message, "\")")
		}
		messages <- "receiver1 received a message"
		closer1.Close()
	})

	go receiver2.BlockingSubscription(func(incoming <-chan string) {
		message := <-incoming
		if message != "Test Test" {
			t.Error("receiver2 didn't get correct message (Got \"", message, "\")")
		}
		messages <- "receiver2 received a message"
	})

	var closer3 io.Closer
	start3, closer3 := receiver3.PatternSubscribe(func(message string) {
		if message != "Test Test" {
			t.Error("receiver3 didn't get correct message (Got \"", message, "\")")
		}
		messages <- "receiver3 received a message"
		closer3.Close()
	})

	go receiver4.BlockingPatternSubscription(func(incoming <-chan string) {
		message := <-incoming
		if message != "Test Test" {
			t.Error("receiver4 didn't get correct message (Got \"", message, "\")")
		}
		messages <- "receiver4 received a message"
	})

	<-start1
	<-start3

	sender.Publish("Test Test")

	timeout := time.NewTimer(2 * time.Second)
	defer timeout.Stop()

	for i := 0; i < 4; i++ {
		select {
		case m := <-messages:
			t.Log(m)
		case <-timeout.C:
			t.Error("Not All Messages Received")
			return
		}
	}

}
