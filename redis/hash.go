package redis

import (
	"strings"
)

type Hash struct {
	Key
}

func newHash(client Executor, key string) Hash {
	return Hash{
		newKey(client, key),
	}
}

func (this Hash) IsValid() <-chan bool {
	c := make(chan bool, 1)
	go func() {
		defer close(c)
		c <- (<-this.Type() == "hash")
	}()
	return c
}

func (this Hash) String(key string) HashString {
	return newHashString(this, key)
}

func (this Hash) Integer(key string) HashInteger {
	return newHashInteger(this, key)
}

func (this Hash) Size() <-chan int {
	command, output := newIntCommand(this.args("hlen"))
	this.Execute(command)
	return output
}

func (this Hash) Get() <-chan map[string]string {
	command, output := newMapCommand(this.args("hgetall"))
	this.Execute(command)
	return output
}

type HashField struct {
	parent Hash
	key    string
}

func newHashField(hash Hash, key string) HashField {
	return HashField{
		parent: hash,
		key:    key,
	}
}

func (this HashField) args(command string, args ...string) []string {
	return append([]string{strings.ToUpper(command), this.parent.key, this.key}, args...)
}

func (this HashField) Delete() <-chan bool {
	command, output := newBoolCommand(this.args("hdel"))
	this.parent.Execute(command)
	return output
}

func (this HashField) Exists() <-chan bool {
	command, output := newBoolCommand(this.args("hexists"))
	this.parent.Execute(command)
	return output
}

type HashString struct {
	HashField
}

func newHashString(hash Hash, key string) HashString {
	return HashString{
		newHashField(hash, key),
	}
}

func (this HashString) Get() <-chan string {
	command, output := newStringCommand(this.args("hget"))
	this.parent.Execute(command)
	return output
}

func (this HashString) Set(val string) <-chan bool {
	command, output := newBoolCommand(this.args("hset", val))
	this.parent.Execute(command)
	return output
}

func (this HashString) SetIfEmpty(val string) <-chan bool {
	command, output := newBoolCommand(this.args("hsetnx", val))
	this.parent.Execute(command)
	return output
}

type HashInteger struct {
	HashField
}

func newHashInteger(hash Hash, key string) HashInteger {
	return HashInteger{
		newHashField(hash, key),
	}
}

func (this HashInteger) Get() <-chan int {
	command, output := newIntCommand(this.args("hget"))
	this.parent.Execute(command)
	return output
}

func (this HashInteger) Set(val int) <-chan bool {
	command, output := newBoolCommand(this.args("hset", itoa(val)))
	this.parent.Execute(command)
	return output
}

func (this HashInteger) SetIfEmpty(val int) <-chan bool {
	command, output := newBoolCommand(this.args("hsetnx", itoa(val)))
	this.parent.Execute(command)
	return output
}

func (this HashInteger) IncrementBy(val int) <-chan int {
	command, output := newIntCommand(this.args("hincrby", itoa(val)))
	this.parent.Execute(command)
	return output
}

func (this HashInteger) DecrementBy(val int) <-chan int {
	command, output := newIntCommand(this.args("hincrby", itoa(-val)))
	this.parent.Execute(command)
	return output
}

type HashFloat struct {
	HashField
}

func newHashFloat(hash Hash, key string) HashFloat {
	return HashFloat{
		newHashField(hash, key),
	}
}

func (this HashFloat) Get() <-chan float64 {
	command, output := newFloatCommand(this.args("hget"))
	this.parent.Execute(command)
	return output
}

func (this HashFloat) Set(val float64) <-chan bool {
	command, output := newBoolCommand(this.args("hset", ftoa(val)))
	this.parent.Execute(command)
	return output
}

func (this HashFloat) SetIfEmpty(val float64) <-chan bool {
	command, output := newBoolCommand(this.args("hsetnx", ftoa(val)))
	this.parent.Execute(command)
	return output
}

func (this HashFloat) IncrementBy(val float64) <-chan int {
	command, output := newIntCommand(this.args("hincrbyfloat", ftoa(val)))
	this.parent.Execute(command)
	return output
}

func (this HashFloat) DecrementBy(val float64) <-chan int {
	command, output := newIntCommand(this.args("hincrbyfloat", ftoa(-val)))
	this.parent.Execute(command)
	return output
}

func (this Hash) Use(e Executor) Hash {
	this.client = e
	return this
}
