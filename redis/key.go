package redis

import (
	"strings"
	"time"
)

type Key struct {
	key    string
	client Executor
}

func newKey(client Executor, key string) Key {
	return Key{
		key:    key,
		client: client,
	}
}

func (this Key) args(command string, arguments ...string) []string {
	return append([]string{strings.ToUpper(command), this.key}, arguments...)
}

func (this Key) Exists() <-chan bool {
	command, output := newBoolCommand(this.args("exists"))
	this.Execute(command)
	return output
}

func (this Key) Delete() <-chan bool {
	command, output := newBoolCommand(this.args("del"))
	this.Execute(command)
	return output
}

func (this Key) Type() <-chan string {
	command, output := newStringCommand(this.args("type"))
	this.Execute(command)
	return output
}

func (this Key) MoveTo(other Key) <-chan nothing {
	command, output := newNilCommand(this.args("rename", other.key))
	this.Execute(command)
	return output
}

func (this Key) MoveToIfEmpty(other Key) <-chan bool {
	command, output := newBoolCommand(this.args("renamenx", other.key))
	this.Execute(command)
	return output
}

func (this Key) ExpireIn(seconds time.Duration) <-chan bool {
	command, output := newBoolCommand(this.args("expire", itoa(int(seconds))))
	this.Execute(command)
	return output
}

func (this Key) ExpireAt(timestamp time.Time) <-chan bool {
	command, output := newBoolCommand(this.args("expireat", itoa(int(timestamp.Unix()))))
	this.Execute(command)
	return output
}

func (this Key) TimeToLive() <-chan int {
	command, output := newIntCommand(this.args("ttl"))
	this.Execute(command)
	return output
}

func (this Key) Execute(command command) {
	this.client.Execute(command)
}

func (this Key) Use(e Executor) Key {
	this.client = e
	return this
}
