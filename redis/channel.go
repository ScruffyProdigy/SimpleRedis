package redis

import (
	"errors"
	//	"fmt"
	"io"
)

const (
	messageBufferSize = 64
)

type Channel struct {
	Key
	client *Client
}

func newChannel(client *Client, key string) Channel {
	return Channel{
		Key:    newKey(client, key),
		client: client,
	}
}

type subscription chan<- bool

func (this *subscription) Close() error {
	if *this == nil {
		return errors.New("Already closed this subscription")
	}
	*this <- true
	close(*this)
	*this = nil
	return nil
}

func getError(rec interface{}) error {
	if err, ok := rec.(error); ok {
		return err
	}
	if str, ok := rec.(string); ok {
		return errors.New(str)
	}
	return errors.New("Unknown Error:" /*+fmt.Sprintf(rec)*/)
}

func messageLoop(conn Connection) <-chan string {
	output := make(chan string, messageBufferSize)
	go func() {
		defer close(output)
		working := true
		for working {
			response := getResponse(conn)

			switch response.subresponses[0].val {
			case "unsubscribe":
				working = false
			case "message":
				output <- response.subresponses[2].val
			}
		}
	}()
	return output
}

func (this Channel) subscribe(action func(string), sub, unsub string) (startSignal <-chan nothing, finishSignaler io.Closer) {
	closer := make(chan bool, 1)
	happened := make(chan nothing, 1)
	go this.blockingSubscription(func(messages <-chan string) {
		happened <- nothing{}
		for {
			select {
			case m := <-messages:
				action(m)
			case <-closer:
				return
			}
		}
	}, sub, unsub)
	subsc := (subscription)(closer)
	return happened, &subsc
}

func (this Channel) Subscribe(action func(string)) (startSignal <-chan nothing, finishSignaler io.Closer) {
	return this.subscribe(action, "subscribe", "unsubscribe")
}

func (this Channel) PatternSubscribe(action func(string)) (startSignal <-chan nothing, finishSignaler io.Closer) {
	return this.subscribe(action, "psubscribe", "punsubscribe")
}

func (this Channel) blockingSubscription(subscription func(<-chan string), sub, unsub string) {
	this.client.useNewConnection(func(conn Connection) {
		subscriber, result := newNilCommand(this.args(sub))
		conn.Execute(subscriber)
		<-result

		defer func() {
			unsubscriber, _ := newNilCommand(this.args(unsub))
			// we can't get a response, because another gorouting is already listening. 
			// We'll just use the input, and have the other side get the output
			conn.input(unsubscriber)
		}()

		output := messageLoop(conn)
		subscription(output)

		return
	})
}

func (this Channel) BlockingSubscription(subscription func(<-chan string)) {
	this.blockingSubscription(subscription, "subscribe", "unsubscribe")
}

func (this Channel) BlockingPatternSubscription(subscription func(<-chan string)) {
	this.blockingSubscription(subscription, "psubscribe", "punsubscribe")
}

func (this Channel) Publish(message string) <-chan int {
	command, result := newIntCommand(this.args("publish", message))
	this.Key.client.Execute(command)
	return result
}

func (this Channel) Use(e Executor) Channel {
	this.Key.client = e
	return this
}
