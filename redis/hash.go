package redis

import (
	"strings"
)

//Hash is a object that implements the redis Hash primitive
//See http://redis.io/commands#hash for more information on Redis Hashes
type Hash struct {
	Key
}

func newHash(client SafeExecutor, key string) Hash {
	return Hash{
		newKey(client, key),
	}
}

//IsValid returns whether the underlying redis object can use the commands in this object
func (this Hash) IsValid() <-chan bool {
	c := make(chan bool, 1)
	go func() {
		defer close(c)
		c <- (<-this.Type() == "hash")
	}()
	return c
}

//String defines a field within the Hash that will be treated as a basic string
func (this Hash) String(key string) HashString {
	return newHashString(this, key)
}

//Integer defines a field within the Hash that will be treated as a basic integer
func (this Hash) Integer(key string) HashInteger {
	return newHashInteger(this, key)
}

//Float defines a field within the Hash that will be treated as a basic float
func (this Hash) Float(key string) HashFloat {
	return newHashFloat(this, key)
}

//HLEN command - 
//Size returns the number of fields that currently exist in the Hash
func (this Hash) Size() <-chan int {
	return IntCommand(this, this.args("hlen")...)
}

//HGETALL command - 
//Get returns a map that contains all of the values in the hash
func (this Hash) Get() <-chan map[string]string {
	return MapCommand(this, this.args("hgetall")...)
}

//HashField implements basic functions that apply to Hash Fields
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

//HDEL command - 
//Delete removes this field from the Hash if it exists; 
//returns whether or not the delete suceeded
func (this HashField) Delete() <-chan bool {
	return BoolCommand(this.parent, this.args("hdel")...)
}

//HEXISTS command - 
//Exists returns whether or not this field exists within the hash
func (this HashField) Exists() <-chan bool {
	return BoolCommand(this.parent, this.args("hexists")...)
}

//HashString implements the basic functions on hash fields that are basic strings
type HashString struct {
	HashField
}

func newHashString(hash Hash, key string) HashString {
	return HashString{
		newHashField(hash, key),
	}
}

//HGET command - 
//Get returns the string that is in this field
func (this HashString) Get() <-chan string {
	return StringCommand(this.parent, this.args("hget")...)
}

//HSET command - 
//Set sets this field to a specific string
func (this HashString) Set(val string) <-chan bool {
	return BoolCommand(this.parent, this.args("hset", val)...)
}

//HSETNX command - 
//SetIfEmpty sets this field to a specific string if there isn't anything in it yet; 
//returns whether or not the command succeeded
func (this HashString) SetIfEmpty(val string) <-chan bool {
	return BoolCommand(this.parent, this.args("hsetnx", val)...)
}

//HashInteger implements the basic functions on hash fields that are basic integers
type HashInteger struct {
	HashField
}

func newHashInteger(hash Hash, key string) HashInteger {
	return HashInteger{
		newHashField(hash, key),
	}
}

//HGET command - 
//Get returns the integer that is in this field
func (this HashInteger) Get() <-chan int {
	return IntCommand(this.parent, this.args("hget")...)
}

//HSET command - 
//Set sets this field to an integer
func (this HashInteger) Set(val int) <-chan bool {
	return BoolCommand(this.parent, this.args("hset", itoa(val))...)
}

//HSETNX command - 
//SetIfEmpty sets this field to an integer but only if it was empty before
func (this HashInteger) SetIfEmpty(val int) <-chan bool {
	return BoolCommand(this.parent, this.args("hsetnx", itoa(val))...)
}

//HINCRBY command - 
//IncremementBy increments the integer in this field by "val"
func (this HashInteger) IncrementBy(val int) <-chan int {
	return IntCommand(this.parent, this.args("hincrby", itoa(val))...)
}

//HINCRBY command - 
//DecrementBy decreases the integer in this field by "val"
func (this HashInteger) DecrementBy(val int) <-chan int {
	return IntCommand(this.parent, this.args("hincrby", itoa(-val))...)
}

//HashFloat is an object that implements the Hash functions that apply to float fields
type HashFloat struct {
	HashField
}

func newHashFloat(hash Hash, key string) HashFloat {
	return HashFloat{
		newHashField(hash, key),
	}
}

//HGET command - 
//Get gets the float in this field
func (this HashFloat) Get() <-chan float64 {
	return FloatCommand(this.parent, this.args("hget")...)
}

//HSET command - 
//Set sets this field to a float
func (this HashFloat) Set(val float64) <-chan bool {
	return BoolCommand(this.parent, this.args("hset", ftoa(val))...)
}

//HSETNX command - 
//SetIfEmpty sets this field to a float if nothing is already in it;
//returns whether or not it succeeded
func (this HashFloat) SetIfEmpty(val float64) <-chan bool {
	return BoolCommand(this.parent, this.args("hsetnx", ftoa(val))...)
}

//HINCRYBYFLOAT command - 
//IncrementBy increases the float in this field by "val"
func (this HashFloat) IncrementBy(val float64) <-chan float64 {
	return FloatCommand(this.parent, this.args("hincrbyfloat", ftoa(val))...)
}

//HINCRBYFLOAT command - 
//DecrementBy decreases the float in this field by "val"
func (this HashFloat) DecrementBy(val float64) <-chan float64 {
	return FloatCommand(this.parent, this.args("hincrbyfloat", ftoa(-val))...)
}

//Use allows you to use this key on a different executor
func (this Hash) Use(e SafeExecutor) Hash {
	this.client = e
	return this
}
