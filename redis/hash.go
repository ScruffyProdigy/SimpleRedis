package redis

import (
	"strings"
)

type Hash struct {
	Key
}

func newHash(client SafeExecutor, key string) Hash {
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

func (this Hash) Float(key string) HashFloat {
	return newHashFloat(this, key)
}

func (this Hash) Size() <-chan int {
	return IntCommand(this, this.args("hlen"))
}

func (this Hash) Get() <-chan map[string]string {
	return MapCommand(this, this.args("hgetall"))
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
	return BoolCommand(this.parent, this.args("hdel"))
}

func (this HashField) Exists() <-chan bool {
	return BoolCommand(this.parent, this.args("hexists"))
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
	return StringCommand(this.parent, this.args("hget"))
}

func (this HashString) Set(val string) <-chan bool {
	return BoolCommand(this.parent, this.args("hset", val))
}

func (this HashString) SetIfEmpty(val string) <-chan bool {
	return BoolCommand(this.parent, this.args("hsetnx", val))
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
	return IntCommand(this.parent, this.args("hget"))
}

func (this HashInteger) Set(val int) <-chan bool {
	return BoolCommand(this.parent, this.args("hset", itoa(val)))
}

func (this HashInteger) SetIfEmpty(val int) <-chan bool {
	return BoolCommand(this.parent, this.args("hsetnx", itoa(val)))
}

func (this HashInteger) IncrementBy(val int) <-chan int {
	return IntCommand(this.parent, this.args("hincrby", itoa(val)))
}

func (this HashInteger) DecrementBy(val int) <-chan int {
	return IntCommand(this.parent, this.args("hincrby", itoa(-val)))
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
	return FloatCommand(this.parent, this.args("hget"))
}

func (this HashFloat) Set(val float64) <-chan bool {
	return BoolCommand(this.parent, this.args("hset", ftoa(val)))
}

func (this HashFloat) SetIfEmpty(val float64) <-chan bool {
	return BoolCommand(this.parent, this.args("hsetnx", ftoa(val)))
}

func (this HashFloat) IncrementBy(val float64) <-chan float64 {
	return FloatCommand(this.parent, this.args("hincrbyfloat", ftoa(val)))
}

func (this HashFloat) DecrementBy(val float64) <-chan float64 {
	return FloatCommand(this.parent, this.args("hincrbyfloat", ftoa(-val)))
}

func (this Hash) Use(e SafeExecutor) Hash {
	this.client = e
	return this
}
