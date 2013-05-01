package redis

//Bits is an object that acts as a Redis string primitive encapsulating the functions that operate on a set of bits
//See http://redis.io/commands#string for more information on Redis Strings
type Bits struct {
	Key
}

func newBits(client SafeExecutor, key string) Bits {
	return Bits{
		newKey(client, key),
	}
}

//IsValid returns whether the underlying redis object can use the commands in this object
func (this Bits) IsValid() <-chan bool {
	c := make(chan bool, 1)
	go func() {
		defer close(c)
		c <- (<-this.Type() == "string")
	}()
	return c
}

//SetTo Sets a specific bit in the field to a specific value
func (this Bits) SetTo(index int, on bool) <-chan bool {
	if on {
		return this.On(index)
	}
	return this.Off(index)
}

//On turns on a specific bit in the field - SETBIT redis command
func (this Bits) On(index int) <-chan bool {
	return BoolCommand(this, this.args("setbit", itoa(index), "1")...)
}

//Off turns off a specific bit in the field - SETBIT redis command
func (this Bits) Off(index int) <-chan bool {
	return BoolCommand(this, this.args("setbit", itoa(index), "0")...)
}

//Get returns whether a specific bit in the field is set - GETBIT redis command
func (this Bits) Get(index int) <-chan bool {
	return BoolCommand(this, this.args("getbit", itoa(index))...)
}

//Count returns the number of bits that are set - BITCOUNT redis command
func (this Bits) Count(start, end int) <-chan int {
	return IntCommand(this, this.args("bitcount")...)
}

//StoreIntersetionOf stores the result of a BITOP AND operation of other bitfields in this bitfield
func (this Bits) StoreIntersectionOf(otherKeys ...Bits) <-chan int {
	args := []string{"BITOP", "AND", this.key}
	for _, key := range otherKeys {
		args = append(args, key.key)
	}
	return IntCommand(this, args...)
}

//StoreUnionOf stores the result of a BITOP OR operation of other bitfields in this bitfield
func (this Bits) StoreUnionOf(otherKeys ...Bits) <-chan int {
	args := []string{"BITOP", "OR", this.key}
	for _, key := range otherKeys {
		args = append(args, key.key)
	}
	return IntCommand(this, args...)
}

//StoreDifferenceOf stores the result of a BITOP XOR operation of other bitfields in this bitfield
func (this Bits) StoreDifferencesOf(otherKeys ...Bits) <-chan int {
	args := []string{"BITOP", "XOR", this.key}
	for _, key := range otherKeys {
		args = append(args, key.key)
	}
	return IntCommand(this, args...)
}

//StoreInverseOf stores the result of a BITOP NOT operation of another bitfield in this bitfield
func (this Bits) StoreInverseOf(otherKey Bits) <-chan int {
	return IntCommand(this, "BITOP", "NOT", this.key, otherKey.key)
}

//Use allows you to use this key on a different executor
func (this Bits) Use(e SafeExecutor) Bits {
	this.client = e
	return this
}
