package redis

type Bits struct {
	Key
}

func newBits(client Executor, key string) Bits {
	return Bits{
		newKey(client, key),
	}
}

func (this Bits) IsValid() <-chan bool {
	c := make(chan bool, 1)
	go func() {
		defer close(c)
		c <- (<-this.Type() == "string")
	}()
	return c
}

func (this Bits) SetTo(index int, on bool) <-chan bool {
	if on {
		return this.On(index)
	}
	return this.Off(index)
}

func (this Bits) On(index int) <-chan bool {
	command, output := newBoolCommand(this.args("setbit", itoa(index), "1"))
	this.Execute(command)
	return output
}

func (this Bits) Off(index int) <-chan bool {
	command, output := newBoolCommand(this.args("setbit", itoa(index), "0"))
	this.Execute(command)
	return output
}

func (this Bits) Get(index int) <-chan bool {
	command, output := newBoolCommand(this.args("getbit", itoa(index)))
	this.Execute(command)
	return output
}

func (this Bits) Count(start, end int) <-chan int {
	command, output := newIntCommand(this.args("bitcount"))
	this.Execute(command)
	return output
}

func (this Bits) And(otherKey, resultKey Bits) <-chan int {
	command, output := newIntCommand([]string{"BITOP", "AND", resultKey.key, this.key, otherKey.key})
	this.Execute(command)
	return output
}

func (this Bits) Or(otherKey, resultKey Bits) <-chan int {
	command, output := newIntCommand([]string{"BITOP", "OR", resultKey.key, this.key, otherKey.key})
	this.Execute(command)
	return output
}

func (this Bits) Xor(otherKey, resultKey Bits) <-chan int {
	command, output := newIntCommand([]string{"BITOP", "XOR", resultKey.key, this.key, otherKey.key})
	this.Execute(command)
	return output
}

func (this Bits) Not(resultKey Bits) <-chan int {
	command, output := newIntCommand([]string{"BITOP", "NOT", resultKey.key, this.key})
	this.Execute(command)
	return output
}

func (this Bits) Use(e Executor) Bits {
	this.client = e
	return this
}
