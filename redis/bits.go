package redis

type Bits struct {
	Key
}

func newBits(client SafeExecutor, key string) Bits {
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
	return BoolCommand(this, this.args("setbit", itoa(index), "1"))
}

func (this Bits) Off(index int) <-chan bool {
	return BoolCommand(this, this.args("setbit", itoa(index), "0"))
}

func (this Bits) Get(index int) <-chan bool {
	return BoolCommand(this, this.args("getbit", itoa(index)))
}

func (this Bits) Count(start, end int) <-chan int {
	return IntCommand(this, this.args("bitcount"))
}

func (this Bits) And(otherKey, resultKey Bits) <-chan int {
	return IntCommand(this, []string{"BITOP", "AND", resultKey.key, this.key, otherKey.key})
}

func (this Bits) Or(otherKey, resultKey Bits) <-chan int {
	return IntCommand(this, []string{"BITOP", "OR", resultKey.key, this.key, otherKey.key})
}

func (this Bits) Xor(otherKey, resultKey Bits) <-chan int {
	return IntCommand(this, []string{"BITOP", "XOR", resultKey.key, this.key, otherKey.key})
}

func (this Bits) Not(resultKey Bits) <-chan int {
	return IntCommand(this, []string{"BITOP", "NOT", resultKey.key, this.key})
}

func (this Bits) Use(e SafeExecutor) Bits {
	this.client = e
	return this
}
