package redis

type Set struct {
	SortableKey
}

func newSet(client SafeExecutor, key string) Set {
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
	return BoolCommand(this, this.args("sadd", item))
}

func (this Set) Remove(item string) <-chan bool {
	return BoolCommand(this, this.args("srem", item))
}

func (this Set) Members() <-chan []string {
	return SliceCommand(this, this.args("smembers"))
}

func (this Set) IsMember(item string) <-chan bool {
	return BoolCommand(this, this.args("sismember", item))
}

func (this Set) Size() <-chan int {
	return IntCommand(this, this.args("scard"))
}

func (this Set) RandomMember() <-chan string {
	return StringCommand(this, this.args("srandmember"))
}

func (this Set) Pop() <-chan string {
	return StringCommand(this, this.args("spop"))
}

func (this Set) Intersection(otherSet Set) <-chan []string {
	return SliceCommand(this, this.args("sinter", otherSet.key))
}

func (this Set) Union(otherSet Set) <-chan []string {
	return SliceCommand(this, this.args("sunion", otherSet.key))
}

func (this Set) Difference(otherSet Set) <-chan []string {
	return SliceCommand(this, this.args("sdiff", otherSet.key))
}

func (this Set) StoreIntersectionOf(setA Set, setB Set) <-chan int {
	return IntCommand(this, this.args("sinterstore", setA.key, setB.key))
}

func (this Set) StoreUnionOf(setA Set, setB Set) <-chan int {
	return IntCommand(this, this.args("sunionstore", setA.key, setB.key))
}

func (this Set) StoreDifferenceOf(setA Set, setB Set) <-chan int {
	return IntCommand(this, this.args("sdiffstore", setA.key, setB.key))
}

func (this Set) MoveMemberTo(newSet Set, item string) <-chan bool {
	return BoolCommand(this, this.args("smove", newSet.key, item))
}

func (this Set) Use(e SafeExecutor) Set {
	this.client = e
	return this
}
