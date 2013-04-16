package redis

type Set struct {
	SortableKey
}

func newSet(client Executor, key string) Set {
	return Set{
		newSortableKey(client, key),
	}
}

func (this Set) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "set")
	}()
	return c
}

func (this Set) Add(item string) <-chan bool {
	command, output := newBoolCommand(this.args("sadd", item))
	this.Execute(command)
	return output
}

func (this Set) Remove(item string) <-chan bool {
	command, output := newBoolCommand(this.args("srem", item))
	this.Execute(command)
	return output
}

func (this Set) Members() <-chan []string {
	command, output := newSliceCommand(this.args("smembers"))
	this.Execute(command)
	return output
}

func (this Set) IsMember(item string) <-chan bool {
	command, output := newBoolCommand(this.args("sismember", item))
	this.Execute(command)
	return output
}

func (this Set) Size() <-chan int {
	command, output := newIntCommand(this.args("scard"))
	this.Execute(command)
	return output
}

func (this Set) RandomMember() <-chan string {
	command, output := newStringCommand(this.args("srandmember"))
	this.Execute(command)
	return output
}

func (this Set) Pop() <-chan string {
	command, output := newStringCommand(this.args("spop"))
	this.Execute(command)
	return output
}

func (this Set) Intersection(otherSet Set) <-chan []string {
	command, output := newSliceCommand(this.args("sinter", otherSet.key))
	this.Execute(command)
	return output
}

func (this Set) Union(otherSet Set) <-chan []string {
	command, output := newSliceCommand(this.args("sunion", otherSet.key))
	this.Execute(command)
	return output
}

func (this Set) Difference(otherSet Set) <-chan []string {
	command, output := newSliceCommand(this.args("sdiff", otherSet.key))
	this.Execute(command)
	return output
}

func (this Set) StoreIntersectionOf(setA Set, setB Set) <-chan int {
	command, output := newIntCommand(this.args("sinterstore", setA.key, setB.key))
	this.Execute(command)
	return output
}

func (this Set) StoreUnionOf(setA Set, setB Set) <-chan int {
	command, output := newIntCommand(this.args("sunionstore", setA.key, setB.key))
	this.Execute(command)
	return output
}

func (this Set) StoreDifferenceOf(setA Set, setB Set) <-chan int {
	command, output := newIntCommand(this.args("sdiffstore", setA.key, setB.key))
	this.Execute(command)
	return output
}

func (this Set) MoveMemberTo(newSet Set, item string) <-chan bool {
	command, output := newBoolCommand(this.args("smove", newSet.key, item))
	this.Execute(command)
	return output
}

func (this Set) Use(e Executor) Set {
	this.client = e
	return this
}
