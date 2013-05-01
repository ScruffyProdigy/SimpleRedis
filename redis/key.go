package redis

import (
	"strings"
	"time"
)

//Key is used as a base to give all types of keys the same basic commands
//see http://redis.io/commands#generic for more information on these types of commands
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

//Exists returns whether or not the key already exists - EXISTS command
func (this Key) Exists() <-chan bool {
	return BoolCommand(this, this.args("exists"))
}

//Delete removes a key from Redis - DEL command
func (this Key) Delete() <-chan bool {
	return BoolCommand(this, this.args("del"))
}

//Type returns the type of the underlying key - TYPE command
//specifically, it will be one of: none, string, list, set, zset and hash.
func (this Key) Type() <-chan string {
	return StringCommand(this, this.args("type"))
}

//MoveTo transfers this key to a different one - RENAME command
func (this Key) MoveTo(other Key) <-chan nothing {
	return NilCommand(this, this.args("rename", other.key))
}

//MoveToIfEmpty transfers this key to a different one, but only if the new one is empty - RENAMENX command
func (this Key) MoveToIfEmpty(other Key) <-chan bool {
	return BoolCommand(this, this.args("renamenx", other.key))
}

//ExpireIn sets the key to expire after a specified duration
//Currently, if the duration is less than an hour, it will set the duration to the nearest millisecond - PEXPIRE command
//if the duration is greater than or equal to an hour, it will set the duration to the nearest second instead - EXPIRE command
func (this Key) ExpireIn(duration time.Duration) <-chan bool {
	//if the time to expire is in a time range larger than an hour, the number of milliseconds probably is not particularly important, so we can use a regular expire
	if duration >= time.Hour {
		return BoolCommand(this, this.args("expire", itoa(int(duration/time.Second))))
	}
	//otherwise use pexpire to get down to the nearest millisecond
	return BoolCommand(this, this.args("pexpire", itoa(int(duration/time.Millisecond))))
}

//ExpireAt sets the key to expire at a specific time - EXPIREAT command
func (this Key) ExpireAt(timestamp time.Time) <-chan bool {
	return BoolCommand(this, this.args("expireat", itoa(int(timestamp.Unix()))))
}

//SecondsToLive returns to number of seconds until this key is set to expire - TTL command
func (this Key) SecondsToLive() <-chan int {
	return IntCommand(this, this.args("ttl"))
}

//MillisecondsToLive returns the number of milliseconds left until this key is set to expire - PTTL command
func (this Key) MillisecondsToLive() <-chan int {
	return IntCommand(this, this.args("pttl"))
}

//Execute allows the Key to be an Executor, which makes things quicker to code
func (this Key) Execute(command command) {
	this.client.Execute(command)
}

//Use allows you to use this key on a different executor
func (this Key) Use(e SafeExecutor) Key {
	this.client = e
	return this
}
