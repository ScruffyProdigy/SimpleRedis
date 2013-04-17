package redis

type Integer struct {
	Key
}

func newInteger(client SafeExecutor, key string) Integer {
	return Integer{
		newKey(client, key),
	}
}

func (this Integer) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "string")
	}()
	return c
}

func (this Integer) Set(val int) <-chan nothing {
	return NilCommand(this, this.args("set", itoa(val)))
}

func (this Integer) SetIfEmpty(val int) <-chan bool {
	return BoolCommand(this, this.args("setnx", itoa(val)))
}

func (this Integer) Get() <-chan int {
	return IntCommand(this, this.args("get"))
}

func (this Integer) GetSet(val int) <-chan int {
	return IntCommand(this, this.args("getset", itoa(val)))
}

func (this Integer) Increment() <-chan int {
	return IntCommand(this, this.args("incr"))
}

func (this Integer) IncrementBy(val int) <-chan int {
	return IntCommand(this, this.args("incrby", itoa(val)))
}

func (this Integer) Decrement() <-chan int {
	return IntCommand(this, this.args("decr"))
}

func (this Integer) DecrementBy(val int) <-chan int {
	return IntCommand(this, this.args("decrby", itoa(val)))
}

func (this Integer) Use(e SafeExecutor) Integer {
	this.client = e
	return this
}
