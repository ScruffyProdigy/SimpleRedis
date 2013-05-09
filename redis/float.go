package redis

//Float is an object that acts as a Redis string primitive encapsulating the functions that operate on a floating point number
//See http://redis.io/commands#string for more information on Redis Strings
type Float struct {
	Key
}

func newFloat(client SafeExecutor, key string) Float {
	return Float{
		newKey(client, key),
	}
}

//IsValid returns whether the underlying redis object can use the commands in this object
func (this Float) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "string")
	}()
	return c
}

//SET command - 
//Set sets the object to a specific floating point value
func (this Float) Set(val float64) <-chan nothing {
	return NilCommand(this, this.args("set", ftoa(val))...)
}

//SETNX command - 
//SetIfEmpty sets the object to a specfic floating point value, but only if it was empty beforehand
func (this Float) SetIfEmpty(val float64) <-chan bool {
	return BoolCommand(this, this.args("setnx", ftoa(val))...)
}

//GET command - 
//Get gets the floating point value stored in the object
func (this Float) Get() <-chan float64 {
	return FloatCommand(this, this.args("get")...)
}

//GETSET command - 
//GetSet gets the current floating point value stored in an object, and sets the value to a new one
func (this Float) GetSet(val float64) <-chan float64 {
	return FloatCommand(this, this.args("getset", ftoa(val))...)
}

//INCRBYFLOAT command - 
//IncrementBy increments the floating point value stored in an object by a set amount and returns the new amount
func (this Float) IncrementBy(val float64) <-chan float64 {
	return FloatCommand(this, this.args("incrbyfloat", ftoa(val))...)
}

//INCRBYFLOAT command - 
//DecrementBy decreases the floating point value stored in an object by a set amount and returns the new amount
func (this Float) DecrementBy(val float64) <-chan float64 {
	return FloatCommand(this, this.args("incrbyfloat", ftoa(-val))...)
}

//Use allows you to use this key on a different executor
func (this Float) Use(e SafeExecutor) Float {
	this.client = e
	return this
}
