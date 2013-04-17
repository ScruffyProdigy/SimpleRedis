package redis

type FloatSet struct {
	SortableKey
}

func newFloatSet(client SafeExecutor, key string) FloatSet {
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
	return BoolCommand(this, this.args("sadd", ftoa(item)))
}

func (this FloatSet) Remove(item float64) <-chan bool {
	return BoolCommand(this, this.args("srem", ftoa(item)))
}

func (this FloatSet) Members() <-chan []float64 {
	output := SliceCommand(this, this.args("smembers"))
	realoutput := make(chan []float64, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			if floats, err := stringsToFloats(slice); err != nil {
				this.client.ErrCallback(err, "smembers")
				return
			} else {
				realoutput <- floats
			}
		}
	}()
	return realoutput
}

func (this FloatSet) IsMember(item float64) <-chan bool {
	return BoolCommand(this, this.args("sismember", ftoa(item)))
}

func (this FloatSet) Size() <-chan int {
	return IntCommand(this, this.args("scard"))
}

func (this FloatSet) RandomMember() <-chan float64 {
	return FloatCommand(this, this.args("srandmember"))
}

func (this FloatSet) Pop() <-chan float64 {
	return FloatCommand(this, this.args("spop"))
}

func (this FloatSet) Intersection(otherSet FloatSet) <-chan []float64 {
	output := SliceCommand(this, this.args("sinter", otherSet.key))
	realoutput := make(chan []float64, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			if floats, err := stringsToFloats(slice); err != nil {
				this.client.ErrCallback(err, "sinter")
				return
			} else {
				realoutput <- floats
			}
		}
	}()
	return realoutput
}

func (this FloatSet) Union(otherSet FloatSet) <-chan []float64 {
	output := SliceCommand(this, this.args("sunion", otherSet.key))
	realoutput := make(chan []float64, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			if floats, err := stringsToFloats(slice); err != nil {
				this.client.ErrCallback(err, "sunion")
				return
			} else {
				realoutput <- floats
			}
		}
	}()
	return realoutput
}

func (this FloatSet) Difference(otherSet FloatSet) <-chan []float64 {
	output := SliceCommand(this, this.args("sdiff", otherSet.key))
	realoutput := make(chan []float64, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			if floats, err := stringsToFloats(slice); err != nil {
				this.client.ErrCallback(err, "sdiff")
				return
			} else {
				realoutput <- floats
			}
		}
	}()
	return realoutput
}

func (this FloatSet) StoreIntersectionOf(setA FloatSet, setB FloatSet) <-chan int {
	return IntCommand(this, this.args("sinterstore", setA.key, setB.key))
}

func (this FloatSet) StoreUnionOf(setA FloatSet, setB FloatSet) <-chan int {
	return IntCommand(this, this.args("sunionstore", setA.key, setB.key))
}

func (this FloatSet) StoreDifferenceOf(setA FloatSet, setB FloatSet) <-chan int {
	return IntCommand(this, this.args("sdiffstore", setA.key, setB.key))
}

func (this FloatSet) MoveMemberTo(newSet FloatSet, item float64) <-chan bool {
	return BoolCommand(this, this.args("smove", newSet.key, ftoa(item)))
}

func (this FloatSet) Use(e SafeExecutor) FloatSet {
	this.client = e
	return this
}
