package redis

type Float struct {
	Key
}

func newFloat(client Executor, key string) Float {
	return Float{
		newKey(client, key),
	}
}

func (this Float) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "string")
	}()
	return c
}

func (this Float) Set(val float64) <-chan nothing {
	command, output := newNilCommand(this.args("set", ftoa(val)))
	this.Execute(command)
	return output
}

func (this Float) SetIfEmpty(val float64) <-chan bool {
	command, output := newBoolCommand(this.args("setnx", ftoa(val)))
	this.Execute(command)
	return output
}

func (this Float) Get() <-chan float64 {
	command, output := newFloatCommand(this.args("get"))
	this.Execute(command)
	return output
}

func (this Float) GetSet(val float64) <-chan float64 {
	command, output := newFloatCommand(this.args("getset", ftoa(val)))
	this.Execute(command)
	return output
}

func (this Float) IncrementBy(val float64) <-chan float64 {
	command, output := newFloatCommand(this.args("incrbyfloat", ftoa(val)))
	this.Execute(command)
	return output
}

func (this Float) DecrementBy(val float64) <-chan float64 {
	command, output := newFloatCommand(this.args("incrbyfloat", ftoa(-val)))
	this.Execute(command)
	return output
}

func (this Float) Use(e Executor) Float {
	this.client = e
	return this
}
