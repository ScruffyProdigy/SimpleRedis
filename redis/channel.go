package redis

import (
	"errors"
	"io"
)

const (
	messageBufferSize = 64
)

//A Channel is object that encapsulates the Pub/Sub redis commands
//See http://redis.io/topics/pubsub for more information on redis Pub/Sub
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

func messageLoop(conn *Connection, errCallback errCallbackFunc) <-chan string {
	output := make(chan string, messageBufferSize)
	go func() {
		defer close(output)
		defer func() {
			if rec := recover(); rec != nil {
				errCallback(getError(rec), "Closing a Channel")
			}
		}()
		working := true
		for working {
			response, err := getResponse(conn)
			if err != nil {
				errCallback.Call(err, "Message Loop Error")
				working = false
			}

			switch response.subresponses[0].val {
			case "unsubscribe":
				working = false
			case "message":
				output <- response.subresponses[2].val
			case "pmessage":
				output <- response.subresponses[3].val
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

//Subscribe calls the specified function whenever a message along this channel is published
//it returns a channel that allows you to know when the channel has started succesfully listening
//and a way to signal when you're done listening
func (this Channel) Subscribe(action func(string)) (startSignal <-chan nothing, finishSignaler io.Closer) {
	return this.subscribe(action, "subscribe", "unsubscribe")
}

//PatternSubscribe calls the specified function whenever a message along any of the channels that fit the pattern is published
//it returns a channel that allows you to know when the channel has started succesfully listening
//and a way to signal when you're done listening
func (this Channel) PatternSubscribe(action func(string)) (startSignal <-chan nothing, finishSignaler io.Closer) {
	return this.subscribe(action, "psubscribe", "punsubscribe")
}

func (this Channel) blockingSubscription(subscription func(<-chan string), sub, unsub string) {
	this.client.useNewConnection(func(conn *Connection) {
		<-NilCommand(conn, this.args(sub)...)

		defer func() {
			<-NilCommand(conn, this.args(unsub)...)
		}()

		subscription(messageLoop(conn, this.client.fErrCallback))
		return
	})
}

//BlockingSubscription sends a message through a go channel whenever a message has been published on this redis channel
//when the function terminates, the subscription is canceled
func (this Channel) BlockingSubscription(subscription func(<-chan string)) {
	this.blockingSubscription(subscription, "subscribe", "unsubscribe")
}

//BlockingPatternSubscription sends a message through a go channel whenever a message is published on any redis channel that fits the pattern
//when the function terminates, the subscription is canceled
func (this Channel) BlockingPatternSubscription(subscription func(<-chan string)) {
	this.blockingSubscription(subscription, "psubscribe", "punsubscribe")
}

//Publish publishes a message on this channel
//Use Subscribe, PatternSubscribe, BlockingSubscription, or BlockingPatternSubscription to receive the published message
func (this Channel) Publish(message string) <-chan int {
	return IntCommand(this, this.args("publish", message)...)
}

//Use allows you to use this key on a different executor
func (this Channel) Use(e SafeExecutor) Channel {
	this.Key.client = e
	return this
}
