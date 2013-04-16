package redis

type FloatSet struct {
	SortableKey
}

func newFloatSet(client Executor, key string) FloatSet {
	return FloatSet{
		newSortableKey(client, key),
	}
}

func (this FloatSet) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "set")
	}()
	return c
}

func (this FloatSet) Add(item float64) <-chan bool {
	command, output := newBoolCommand(this.args("sadd", ftoa(item)))
	this.Execute(command)
	return output
}

func (this FloatSet) Remove(item float64) <-chan bool {
	command, output := newBoolCommand(this.args("srem", ftoa(item)))
	this.Execute(command)
	return output
}

func (this FloatSet) Members() <-chan []float64 {
	command, output := newSliceCommand(this.args("smembers"))
	this.Execute(command)
	realoutput := make(chan []float64, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			floats, err := stringsToFloats(slice)
			if err != nil {
				this.client.ErrCallback(err, "smembers")
				return
			}
			realoutput <- floats
		}
	}()
	return realoutput
}

func (this FloatSet) IsMember(item float64) <-chan bool {
	command, output := newBoolCommand(this.args("sismember", ftoa(item)))
	this.Execute(command)
	return output
}

func (this FloatSet) Size() <-chan int {
	command, output := newIntCommand(this.args("scard"))
	this.Execute(command)
	return output
}

func (this FloatSet) RandomMember() <-chan float64 {
	command, output := newFloatCommand(this.args("srandmember"))
	this.Execute(command)
	return output
}

func (this FloatSet) Pop() <-chan float64 {
	command, output := newFloatCommand(this.args("spop"))
	this.Execute(command)
	return output
}

func (this FloatSet) Intersection(otherSet FloatSet) <-chan []float64 {
	command, output := newSliceCommand(this.args("sinter", otherSet.key))
	this.Execute(command)
	realoutput := make(chan []float64, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			floats, err := stringsToFloats(slice)
			if err != nil {
				this.client.ErrCallback(err, "sinter")
				return
			}
			realoutput <- floats
		}
	}()
	return realoutput
}

func (this FloatSet) Union(otherSet FloatSet) <-chan []float64 {
	command, output := newSliceCommand(this.args("sunion", otherSet.key))
	this.Execute(command)
	realoutput := make(chan []float64, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			floats, err := stringsToFloats(slice)
			if err != nil {
				this.client.ErrCallback(err, "sunion")
				return
			}
			realoutput <- floats
		}
	}()
	return realoutput
}

func (this FloatSet) Difference(otherSet FloatSet) <-chan []float64 {
	command, output := newSliceCommand(this.args("sdiff", otherSet.key))
	this.Execute(command)
	realoutput := make(chan []float64, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			floats, err := stringsToFloats(slice)
			if err != nil {
				this.client.ErrCallback(err, "sdiff")
				return
			}
			realoutput <- floats
		}
	}()
	return realoutput
}

func (this FloatSet) StoreIntersectionOf(setA FloatSet, setB FloatSet) <-chan int {
	command, output := newIntCommand(this.args("sinterstore", setA.key, setB.key))
	this.Execute(command)
	return output
}

func (this FloatSet) StoreUnionOf(setA FloatSet, setB FloatSet) <-chan int {
	command, output := newIntCommand(this.args("sunionstore", setA.key, setB.key))
	this.Execute(command)
	return output
}

func (this FloatSet) StoreDifferenceOf(setA FloatSet, setB FloatSet) <-chan int {
	command, output := newIntCommand(this.args("sdiffstore", setA.key, setB.key))
	this.Execute(command)
	return output
}

func (this FloatSet) MoveMemberTo(newSet FloatSet, item float64) <-chan bool {
	command, output := newBoolCommand(this.args("smove", newSet.key, ftoa(item)))
	this.Execute(command)
	return output
}

func (this FloatSet) Use(e Executor) FloatSet {
	this.client = e
	return this
}
