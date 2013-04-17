package redis

type String struct {
	Key
}

func newString(client SafeExecutor, key string) String {
	return String{
		newKey(client, key),
	}
}

func (this String) IsValid() <-chan bool {
	c := make(chan bool, 1)
	go func() {
		defer close(c)
		c <- (<-this.Type() == "string")
	}()
	return c
}

func (this String) Set(val string) <-chan nothing {
	return NilCommand(this, this.args("set", val))
}
func (this String) SetIfEmpty(val string) <-chan bool {
	return BoolCommand(this, this.args("setnx", val))
}
func (this String) Get() <-chan string {
	return StringCommand(this, this.args("get"))
}

func (this String) Replace(val string) <-chan string {
	return StringCommand(this, this.args("getset", val))
}
func (this String) Append(val string) <-chan int {
	return IntCommand(this, this.args("append", val))
}
func (this String) Length() <-chan int {
	return IntCommand(this, this.args("strlen"))
}
func (this String) Use(e SafeExecutor) String {
	this.client = e
	return this
}
