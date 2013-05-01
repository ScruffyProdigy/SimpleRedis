package redis

//Set is an object that implements a basic Redis Set primitive
//see http://redis.io/commands#set for more information on redis sets
type Set struct {
	SortableKey
}

func newSet(client SafeExecutor, key string) Set {
	return Set{
		newSortableKey(client, key),
	}
}

//IsValid returns whether the underlying redis object can use the commands in this object
func (this Set) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "set")
	}()
	return c
}

//Add adds a string to the set if it isn't already there - SADD command
//returns whether or not the add succeeded
func (this Set) Add(item string) <-chan bool {
	return BoolCommand(this, this.args("sadd", item))
}

//Remove removes a string from the set if it exists - SREM command
//returns whether or not the string existed in the set
func (this Set) Remove(item string) <-chan bool {
	return BoolCommand(this, this.args("srem", item))
}

//Members returns all of the strings in the set - SMEMBERS command
func (this Set) Members() <-chan []string {
	return SliceCommand(this, this.args("smembers"))
}

//IsMember returns whether or not the string is a member of the set - SISMEMBER command
func (this Set) IsMember(item string) <-chan bool {
	return BoolCommand(this, this.args("sismember", item))
}

//Size returns the number of strings in the set - SCARD command
func (this Set) Size() <-chan int {
	return IntCommand(this, this.args("scard"))
}

//RandomMember returns a random string from the set - SRANDMEMBER command
func (this Set) RandomMember() <-chan string {
	return StringCommand(this, this.args("srandmember"))
}

//Pop removes a random string from the set and returns it - SPOP command
func (this Set) Pop() <-chan string {
	return StringCommand(this, this.args("spop"))
}

//TODO: allow variable parameters so that multiple sets can be intersected, unioned, or diffed

//Intersection returns all of the strings that are in both this set and another - SINTER command
func (this Set) Intersection(otherSet Set) <-chan []string {
	return SliceCommand(this, this.args("sinter", otherSet.key))
}

//Union returns all of the strings that are either in this set or another - SUNION command
func (this Set) Union(otherSet Set) <-chan []string {
	return SliceCommand(this, this.args("sunion", otherSet.key))
}

//Difference returns all of the strings that are in this set, but not another - SDIFF command
func (this Set) Difference(otherSet Set) <-chan []string {
	return SliceCommand(this, this.args("sdiff", otherSet.key))
}

//StoreIntersectionOf finds the intersection of two other sets and stores it in this set - SINTERSTORE command
//returns the size of the resulting set
func (this Set) StoreIntersectionOf(setA Set, setB Set) <-chan int {
	return IntCommand(this, this.args("sinterstore", setA.key, setB.key))
}

//StoreUnionOf finds the union of two other sets and stores it in this set - SUNIONSTORE command
//returns the size of the resulting set
func (this Set) StoreUnionOf(setA Set, setB Set) <-chan int {
	return IntCommand(this, this.args("sunionstore", setA.key, setB.key))
}

//StoreDifferenceOf finds the difference of two other sets and stores it in this set - SDIFFSTORE command
//returns the size of the resulting set
func (this Set) StoreDifferenceOf(setA Set, setB Set) <-chan int {
	return IntCommand(this, this.args("sdiffstore", setA.key, setB.key))
}

//MoveMemberTo removes a string from this set and adds it to another - SMOVE command
//nothing happens if the string doesn't exist in this set
func (this Set) MoveMemberTo(newSet Set, item string) <-chan bool {
	return BoolCommand(this, this.args("smove", newSet.key, item))
}

//Use allows you to use this key on a different executor
func (this Set) Use(e SafeExecutor) Set {
	this.client = e
	return this
}
