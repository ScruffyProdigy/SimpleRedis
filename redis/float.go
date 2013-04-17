package redis

type Float struct {
	Key
}

func newFloat(client SafeExecutor, key string) Float {
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
	return NilCommand(this, this.args("set", ftoa(val)))
}

func (this Float) SetIfEmpty(val float64) <-chan bool {
	return BoolCommand(this, this.args("setnx", ftoa(val)))
}

func (this Float) Get() <-chan float64 {
	return FloatCommand(this, this.args("get"))
}

func (this Float) GetSet(val float64) <-chan float64 {
	return FloatCommand(this, this.args("getset", ftoa(val)))
}

func (this Float) IncrementBy(val float64) <-chan float64 {
	return FloatCommand(this, this.args("incrbyfloat", ftoa(val)))
}

func (this Float) DecrementBy(val float64) <-chan float64 {
	return FloatCommand(this, this.args("incrbyfloat", ftoa(-val)))
}

func (this Float) Use(e SafeExecutor) Float {
	this.client = e
	return this
}
