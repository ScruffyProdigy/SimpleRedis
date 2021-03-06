package redis

//TODO: refactor to use Set code

//IntSet implements the Redis Set primitive assuming all inputs are integers (which is useful for indexes)
//see http://redis.io/commands#set for more information on redis sets
type IntSet struct {
	SortableKey
}

func newIntSet(client SafeExecutor, key string) IntSet {
	return IntSet{
		newSortableKey(client, key),
	}
}

//IsValid returns whether the underlying redis object can use the commands in this object
func (this IntSet) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "set")
	}()
	return c
}

//SADD command - 
//Add adds an integer to the set
func (this IntSet) Add(item int) <-chan bool {
	return BoolCommand(this, this.args("sadd", itoa(item))...)
}

//SREM command -
//Remove removes an integer from the set
func (this IntSet) Remove(item int) <-chan bool {
	return BoolCommand(this, this.args("srem", itoa(item))...)
}

//SMEMBERS command -
//Members lists all of the integers in the set
func (this IntSet) Members() <-chan []int {
	return intsChannel(SliceCommand(this, this.args("smembers")...))
}

//SISMEMBER command -
//IsMember returns whether or not an integer is part of the set
func (this IntSet) IsMember(item int) <-chan bool {
	return BoolCommand(this, this.args("sismember", itoa(item))...)
}

//SCARD command -
//Size returns the number of elements in the set
func (this IntSet) Size() <-chan int {
	return IntCommand(this, this.args("scard")...)
}

//SRANDMEMBER command -
//RandomMember returns a random integer from the set
func (this IntSet) RandomMember() <-chan int {
	return IntCommand(this, this.args("srandmember")...)
}

//SPOP command -
//Pop removes a random integer from the set and returns it to you
func (this IntSet) Pop() <-chan int {
	return IntCommand(this, this.args("spop")...)
}

//SINTER command -
//Intersection returns a list of all integers that this and another set have in common
func (this IntSet) Intersection(otherSets ...IntSet) <-chan []int {
	args := this.args("sinter")
	for _, set := range otherSets {
		args = append(args, set.key)
	}
	return intsChannel(SliceCommand(this, args...))
}

//SUNION command -
//Union returns a list of all integers that are either in this set or another
func (this IntSet) Union(otherSets ...IntSet) <-chan []int {
	args := this.args("sunion")
	for _, set := range otherSets {
		args = append(args, set.key)
	}
	return intsChannel(SliceCommand(this, args...))
}

//SDIFF command -
//Difference returns a list of all integers that are in this set, but not another
func (this IntSet) Difference(otherSets ...IntSet) <-chan []int {
	args := this.args("sdiff")
	for _, set := range otherSets {
		args = append(args, set.key)
	}
	return intsChannel(SliceCommand(this, args...))
}

//SINTERSTORE command -
//StoreIntersectionOf finds the intersection of multiple other sets and stores it in this one.
//It returns the number of elements in the new set
func (this IntSet) StoreIntersectionOf(sets ...IntSet) <-chan int {
	args := this.args("sinterstore")
	for _, set := range sets {
		args = append(args, set.key)
	}
	return IntCommand(this, args...)
}

//SUNIONSTORE command -
//StoreUnionOf finds the union of multiple other sets and stores it in this one.
//It returns the number of elements in the new set
func (this IntSet) StoreUnionOf(sets ...IntSet) <-chan int {
	args := this.args("sunionstore")
	for _, set := range sets {
		args = append(args, set.key)
	}
	return IntCommand(this, args...)
}

//SDIFFSTORE command -
//StoreDifferenceOf finds the difference between two other sets and stores it in this one.
//It returns the number of elements in the new set
func (this IntSet) StoreDifferenceOf(sets ...IntSet) <-chan int {
	args := this.args("sdiffstore")
	for _, set := range sets {
		args = append(args, set.key)
	}
	return IntCommand(this, args...)
}

//SMOVE command -
//MoveMemberTo removes an integer from this set if it exists, and then adds it to another set.
//Nothing happens if the integer was not a member of this set
func (this IntSet) MoveMemberTo(newSet IntSet, item int) <-chan bool {
	return BoolCommand(this, this.args("smove", newSet.key, itoa(item))...)
}

//Use allows you to use this key on a different executor
func (this IntSet) Use(e SafeExecutor) IntSet {
	this.client = e
	return this
}
