package redis

//Float is an object that acts as a Redis string primitive encapsulating the functions that operate on a floating point number
//See http://redis.io/commands#string for more information on Redis Strings
type Integer struct {
	Key
}

func newInteger(client SafeExecutor, key string) Integer {
	return Integer{
		newKey(client, key),
	}
}

//IsValid returns whether the underlying redis object can use the commands in this object
func (this Integer) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "string")
	}()
	return c
}

//Set sets the Integer to "val" - SET command
func (this Integer) Set(val int) <-chan nothing {
	return NilCommand(this, this.args("set", itoa(val)))
}

//SetIfEmpty sets the integer to "val", but only if it was empty before - SETNX command
func (this Integer) SetIfEmpty(val int) <-chan bool {
	return BoolCommand(this, this.args("setnx", itoa(val)))
}

//Get returns the the value of this integer - GET command
func (this Integer) Get() <-chan int {
	return IntCommand(this, this.args("get"))
}

//Gets the value of this integer before setting it to something else - GETSET command
func (this Integer) GetSet(val int) <-chan int {
	return IntCommand(this, this.args("getset", itoa(val)))
}

//Increment increases the value of this integer and returns the new value - INCR command
func (this Integer) Increment() <-chan int {
	return IntCommand(this, this.args("incr"))
}

//IncrementBy increases the value of this integer by "val", and returns the new value - INCRBY command
func (this Integer) IncrementBy(val int) <-chan int {
	return IntCommand(this, this.args("incrby", itoa(val)))
}

//Decrement decrements this integer and returns the new value - DECR command
func (this Integer) Decrement() <-chan int {
	return IntCommand(this, this.args("decr"))
}

//DecrementBy decreases this integer by "val", and returns the new value - DECRBY command
func (this Integer) DecrementBy(val int) <-chan int {
	return IntCommand(this, this.args("decrby", itoa(val)))
}

//Use allows you to use this key on a different executor
func (this Integer) Use(e SafeExecutor) Integer {
	this.client = e
	return this
}
