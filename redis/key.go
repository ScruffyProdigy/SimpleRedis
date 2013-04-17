package redis

import (
	"strings"
	"time"
)

type Key struct {
	key    string
	client SafeExecutor
}

func newKey(client SafeExecutor, key string) Key {
	return Key{
		key:    key,
		client: client,
	}
}

func (this Key) args(command string, arguments ...string) []string {
	return append([]string{strings.ToUpper(command), this.key}, arguments...)
}

func (this Key) Exists() <-chan bool {
	return BoolCommand(this, this.args("exists"))
}

func (this Key) Delete() <-chan bool {
	return BoolCommand(this, this.args("del"))
}

func (this Key) Type() <-chan string {
	return StringCommand(this, this.args("type"))
}

func (this Key) MoveTo(other Key) <-chan nothing {
	return NilCommand(this, this.args("rename", other.key))
}

func (this Key) MoveToIfEmpty(other Key) <-chan bool {
	return BoolCommand(this, this.args("renamenx", other.key))
}

func (this Key) ExpireIn(duration time.Duration) <-chan bool {
	//if the time to expire is in a time range larger than an hour, the number of milliseconds probably is not particularly important, so we can use a regular expire
	if duration >= time.Hour {
		return BoolCommand(this, this.args("expire", itoa(int(duration/time.Second))))
	}
	//otherwise use pexpire to get down to the nearest millisecond
	return BoolCommand(this, this.args("pexpire", itoa(int(duration/time.Millisecond))))
}

func (this Key) ExpireAt(timestamp time.Time) <-chan bool {
	return BoolCommand(this, this.args("expireat", itoa(int(timestamp.Unix()))))
}

func (this Key) SecondsToLive() <-chan int {
	return IntCommand(this, this.args("ttl"))
}

func (this Key) MillisecondsToLive() <-chan int {
	return IntCommand(this, this.args("pttl"))
}

func (this Key) Execute(command command) {
	this.client.Execute(command)
}

func (this Key) Use(e SafeExecutor) Key {
	this.client = e
	return this
}
