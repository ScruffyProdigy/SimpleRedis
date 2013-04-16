package redis

type IntSet struct {
	SortableKey
}

func newIntSet(client Executor, key string) IntSet {
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
	command, output := newBoolCommand(this.args("sadd", itoa(item)))
	this.Execute(command)
	return output
}

func (this IntSet) Remove(item int) <-chan bool {
	command, output := newBoolCommand(this.args("srem", itoa(item)))
	this.Execute(command)
	return output
}

func (this IntSet) Members() <-chan []int {
	command, output := newSliceCommand(this.args("smembers"))
	this.Execute(command)
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
	command, output := newBoolCommand(this.args("sismember", itoa(item)))
	this.Execute(command)
	return output
}

func (this IntSet) Size() <-chan int {
	command, output := newIntCommand(this.args("scard"))
	this.Execute(command)
	return output
}

func (this IntSet) RandomMember() <-chan int {
	command, output := newIntCommand(this.args("srandmember"))
	this.Execute(command)
	return output
}

func (this IntSet) Pop() <-chan int {
	command, output := newIntCommand(this.args("spop"))
	this.Execute(command)
	return output
}

func (this IntSet) Intersection(otherSet IntSet) <-chan []int {
	command, output := newSliceCommand(this.args("sinter", otherSet.key))
	this.Execute(command)
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
	command, output := newSliceCommand(this.args("sunion", otherSet.key))
	this.Execute(command)
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
	command, output := newSliceCommand(this.args("sdiff", otherSet.key))
	this.Execute(command)
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
	command, output := newIntCommand(this.args("sinterstore", setA.key, setB.key))
	this.Execute(command)
	return output
}

func (this IntSet) StoreUnionOf(setA IntSet, setB IntSet) <-chan int {
	command, output := newIntCommand(this.args("sunionstore", setA.key, setB.key))
	this.Execute(command)
	return output
}

func (this IntSet) StoreDifferenceOf(setA IntSet, setB IntSet) <-chan int {
	command, output := newIntCommand(this.args("sdiffstore", setA.key, setB.key))
	this.Execute(command)
	return output
}

func (this IntSet) MoveMemberTo(newSet IntSet, item int) <-chan bool {
	command, output := newBoolCommand(this.args("smove", newSet.key, itoa(item)))
	this.Execute(command)
	return output
}

func (this IntSet) Use(e Executor) IntSet {
	this.client = e
	return this
}
