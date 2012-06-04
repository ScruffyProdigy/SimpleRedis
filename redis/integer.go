package redis

type Integer struct {
	Key
}

func newInteger(client Executor, key string) Integer {
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
	command, output := newNilCommand(this.args("set", itoa(val)))
	this.Execute(command)
	return output
}

func (this Integer) SetIfEmpty(val int) <-chan bool {
	command, output := newBoolCommand(this.args("setnx", itoa(val)))
	this.Execute(command)
	return output
}

func (this Integer) Get() <-chan int {
	command, output := newIntCommand(this.args("get"))
	this.Execute(command)
	return output
}

func (this Integer) GetSet(val int) <-chan int {
	command, output := newIntCommand(this.args("getset", itoa(val)))
	this.Execute(command)
	return output
}

func (this Integer) Increment() <-chan int {
	command, output := newIntCommand(this.args("incr"))
	this.Execute(command)
	return output
}

func (this Integer) IncrementBy(val int) <-chan int {
	command, output := newIntCommand(this.args("incrby", itoa(val)))
	this.Execute(command)
	return output
}

func (this Integer) Decrement() <-chan int {
	command, output := newIntCommand(this.args("decr"))
	this.Execute(command)
	return output
}

func (this Integer) DecrementBy(val int) <-chan int {
	command, output := newIntCommand(this.args("decrby", itoa(val)))
	this.Execute(command)
	return output
}

func (this Integer) Use(e Executor) Integer {
	this.client = e
	return this
}
