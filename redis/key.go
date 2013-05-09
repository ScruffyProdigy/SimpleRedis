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

//EXISTS command - 
//Exists returns whether or not the key already exists
func (this Key) Exists() <-chan bool {
	return BoolCommand(this, this.args("exists")...)
}

//DEL command - 
//Delete removes a key from Redis
func (this Key) Delete() <-chan bool {
	return BoolCommand(this, this.args("del")...)
}

//TYPE command - 
//Type returns the type of the underlying key,
//specifically, it will be one of: none, string, list, set, zset and hash.
func (this Key) Type() <-chan string {
	return StringCommand(this, this.args("type")...)
}

//RENAME command - 
//MoveTo transfers this key to a different one
func (this Key) MoveTo(other Key) <-chan nothing {
	return NilCommand(this, this.args("rename", other.key)...)
}

//RENAMENX command - 
//MoveToIfEmpty transfers this key to a different one, but only if the new one is empty
func (this Key) MoveToIfEmpty(other Key) <-chan bool {
	return BoolCommand(this, this.args("renamenx", other.key)...)
}

//PEXPIRE or EXPIRE command - 
//ExpireIn sets the key to expire after a specified duration. 
//Currently, if the duration is less than an hour, it will set the duration to the nearest millisecond;
//if the duration is greater than or equal to an hour, it will set the duration to the nearest second instead
func (this Key) ExpireIn(duration time.Duration) <-chan bool {
	//if the time to expire is in a time range larger than an hour, the number of milliseconds probably is not particularly important, so we can use a regular expire
	if duration >= time.Hour {
		return BoolCommand(this, this.args("expire", itoa(int(duration/time.Second)))...)
	}
	//otherwise use pexpire to get down to the nearest millisecond
	return BoolCommand(this, this.args("pexpire", itoa(int(duration/time.Millisecond)))...)
}

//EXPIREAT command - 
//ExpireAt sets the key to expire at a specific time
func (this Key) ExpireAt(timestamp time.Time) <-chan bool {
	return BoolCommand(this, this.args("expireat", itoa(int(timestamp.Unix())))...)
}

//TTL command - 
//SecondsToLive returns to number of seconds until this key is set to expire
func (this Key) SecondsToLive() <-chan int {
	return IntCommand(this, this.args("ttl")...)
}

//PTTL command - 
//MillisecondsToLive returns the number of milliseconds left until this key is set to expire
func (this Key) MillisecondsToLive() <-chan int {
	return IntCommand(this, this.args("pttl")...)
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
