package redis

//String is an object which implements a basic Redis String primitive
type String struct {
	Key
}

func newString(client SafeExecutor, key string) String {
	return String{
		newKey(client, key),
	}
}

//IsValid returns whether the underlying redis object can use the commands in this object
func (this String) IsValid() <-chan bool {
	c := make(chan bool, 1)
	go func() {
		defer close(c)
		c <- (<-this.Type() == "string")
	}()
	return c
}

//Set sets the value of the key or updates it if it already exists
//returns true if setting, false if updating
func (this String) Set(val string) <-chan nothing {
	return NilCommand(this, this.args("set", val)...)
}

//SetIfEmpty sets the value of the key, but does nothing if it already exists
//returns true if setting, false if skipping
func (this String) SetIfEmpty(val string) <-chan bool {
	return BoolCommand(this, this.args("setnx", val)...)
}

//Get returns the value of the key
func (this String) Get() <-chan string {
	return StringCommand(this, this.args("get")...)
}

//Replace sets the value of the key and returns its old value
func (this String) Replace(val string) <-chan string {
	return StringCommand(this, this.args("getset", val)...)
}

//Append appends the value to the end of the key
func (this String) Append(val string) <-chan int {
	return IntCommand(this, this.args("append", val)...)
}

//Length returns the number of characters in the value of the key
func (this String) Length() <-chan int {
	return IntCommand(this, this.args("strlen")...)
}

//Use allows you to use this key on a different executor
func (this String) Use(e SafeExecutor) String {
	this.client = e
	return this
}
