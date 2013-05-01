package redis

type Prefix interface {
	//Key creates a basic key; you probably won't use this directly very often
	Key(key string) Key

	//String creates the definition for a basic Redis String primitive
	String(key string) String

	//Integer creates the definition for a Redis String primitive that contains an integer
	Integer(key string) Integer

	//Float creates the definition for a Redis String primitive that contains a floating point number
	Float(key string) Float

	//Bits creates the definition for a Redis String primitive that contains a bitfield
	Bits(key string) Bits

	//Hash creates the definition for a basic Redis Hash primitive
	Hash(key string) Hash

	//List creates the definition for a basic Redis List primitive
	List(key string) List

	//IntList creates the definition for a Redis List primitive that contains only integers
	IntList(key string) IntList

	//Set creates the definition for a basic Redis Set primitive
	Set(key string) Set

	//IntSet creates the definition for a Redis Set primitive that contains only integers
	IntSet(key string) IntSet

	//SortedSet creates the definition for a basic Redis ZSet primitive
	SortedSet(key string) SortedSet

	//SortedIntSet creates the definition for a Redis ZSet primitive that contains only integers
	SortedIntSet(key string) SortedIntSet

	//Mutex creates a Mutex within redis
	Mutex(key string) Mutex

	//Semaphore creates a Semaphore within redis
	Semaphore(key string, count int) Mutex

	//ReadWriteMutex creates a Read/Write Mutex within redis
	ReadWriteMutex(key string, readers int) *ReadWriteMutex

	//Channel defines a pub/sub channel within redis
	Channel(key string) Channel

	//Prefix allows you to create a namespace for other redis primitives to help make sure there are no duplication conflicts
	Prefix(key string) Prefix
}

type prefix struct {
	parent Prefix
	root   string
}

func (this *prefix) Key(key string) Key {
	return this.parent.Key(this.root + key)
}

func (this *prefix) String(key string) String {
	return this.parent.String(this.root + key)
}

func (this *prefix) Integer(key string) Integer {
	return this.parent.Integer(this.root + key)
}

func (this *prefix) Float(key string) Float {
	return this.parent.Float(this.root + key)
}

func (this *prefix) Bits(key string) Bits {
	return this.parent.Bits(this.root + key)
}

func (this *prefix) Hash(key string) Hash {
	return this.parent.Hash(this.root + key)
}

func (this *prefix) List(key string) List {
	return this.parent.List(this.root + key)
}

func (this *prefix) IntList(key string) IntList {
	return this.parent.IntList(this.root + key)
}

func (this *prefix) Set(key string) Set {
	return this.parent.Set(this.root + key)
}

func (this *prefix) IntSet(key string) IntSet {
	return this.parent.IntSet(this.root + key)
}

func (this *prefix) SortedSet(key string) SortedSet {
	return this.parent.SortedSet(this.root + key)
}

func (this *prefix) SortedIntSet(key string) SortedIntSet {
	return this.parent.SortedIntSet(this.root + key)
}

func (this *prefix) Mutex(key string) Mutex {
	return this.parent.Mutex(this.root + key)
}

func (this *prefix) Semaphore(key string, count int) Mutex {
	return this.parent.Semaphore(this.root+key, count)
}

func (this *prefix) ReadWriteMutex(key string, readers int) *ReadWriteMutex {
	return this.parent.ReadWriteMutex(this.root+key, readers)
}

func (this *prefix) Channel(key string) Channel {
	return this.parent.Channel(this.root + key)
}

func (this *prefix) Prefix(key string) Prefix {
	return newPrefix(this, key)
}

func newPrefix(parent Prefix, key string) Prefix {
	p := new(prefix)
	p.parent = parent
	p.root = key
	return p
}
