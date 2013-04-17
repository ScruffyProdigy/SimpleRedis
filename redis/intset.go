package redis

type IntSet struct {
	SortableKey
}

func newIntSet(client SafeExecutor, key string) IntSet {
	return IntSet{
		newSortableKey(client, key),
	}
}

func (this IntSet) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "set")
	}()
	return c
}

func (this IntSet) Add(item int) <-chan bool {
	return BoolCommand(this, this.args("sadd", itoa(item)))
}

func (this IntSet) Remove(item int) <-chan bool {
	return BoolCommand(this, this.args("srem", itoa(item)))
}

func (this IntSet) Members() <-chan []int {
	output := SliceCommand(this, this.args("smembers"))
	realoutput := make(chan []int, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			ints, err := stringsToInts(slice)
			if err != nil {
				this.client.ErrCallback(err, "smembers")
			}
			realoutput <- ints
		}
	}()
	return realoutput
}

func (this IntSet) IsMember(item int) <-chan bool {
	return BoolCommand(this, this.args("sismember", itoa(item)))
}

func (this IntSet) Size() <-chan int {
	return IntCommand(this, this.args("scard"))
}

func (this IntSet) RandomMember() <-chan int {
	return IntCommand(this, this.args("srandmember"))
}

func (this IntSet) Pop() <-chan int {
	return IntCommand(this, this.args("spop"))
}

func (this IntSet) Intersection(otherSet IntSet) <-chan []int {
	output := SliceCommand(this, this.args("sinter", otherSet.key))
	realoutput := make(chan []int, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			ints, err := stringsToInts(slice)
			if err != nil {
				this.client.ErrCallback(err, "sinter")
			}
			realoutput <- ints
		}
	}()
	return realoutput
}

func (this IntSet) Union(otherSet IntSet) <-chan []int {
	output := SliceCommand(this, this.args("sunion", otherSet.key))
	realoutput := make(chan []int, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			ints, err := stringsToInts(slice)
			if err != nil {
				this.client.ErrCallback(err, "sunion")
			}
			realoutput <- ints
		}
	}()
	return realoutput
}

func (this IntSet) Difference(otherSet IntSet) <-chan []int {
	output := SliceCommand(this, this.args("sdiff", otherSet.key))
	realoutput := make(chan []int, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			ints, err := stringsToInts(slice)
			if err != nil {
				this.client.ErrCallback(err, "sdiff")
			}
			realoutput <- ints
		}
	}()
	return realoutput
}

func (this IntSet) StoreIntersectionOf(setA IntSet, setB IntSet) <-chan int {
	return IntCommand(this, this.args("sinterstore", setA.key, setB.key))
}

func (this IntSet) StoreUnionOf(setA IntSet, setB IntSet) <-chan int {
	return IntCommand(this, this.args("sunionstore", setA.key, setB.key))
}

func (this IntSet) StoreDifferenceOf(setA IntSet, setB IntSet) <-chan int {
	return IntCommand(this, this.args("sdiffstore", setA.key, setB.key))
}

func (this IntSet) MoveMemberTo(newSet IntSet, item int) <-chan bool {
	return BoolCommand(this, this.args("smove", newSet.key, itoa(item)))
}

func (this IntSet) Use(e SafeExecutor) IntSet {
	this.client = e
	return this
}
