package redis

type String struct {
	Key
}

func newString(client Executor, key string) String {
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
	command, output := newNilCommand(this.args("set", val))
	this.Execute(command)
	return output
}

func (this String) SetIfEmpty(val string) <-chan bool {
	command, output := newBoolCommand(this.args("setnx", val))
	this.Execute(command)
	return output
}

func (this String) Get() <-chan string {
	command, output := newStringCommand(this.args("get"))
	this.Execute(command)
	return output
}

func (this String) Clear() <-chan string {
	val := this.Get()
	<-this.Delete()
	return val
}

func (this String) Replace(val string) <-chan string {
	command, output := newStringCommand(this.args("getset", val))
	this.Execute(command)
	return output
}

func (this String) Append(val string) <-chan int {
	command, output := newIntCommand(this.args("append", val))
	this.Execute(command)
	return output
}

func (this String) Length() <-chan int {
	command, output := newIntCommand(this.args("strlen"))
	this.Execute(command)
	return output
}

func (this String) Use(e Executor) String {
	this.client = e
	return this
}
