package redis

type Prefix interface {
	Key(key string) Key
	String(key string) String
	Integer(key string) Integer
	Float(key string) Float
	Hash(key string) Hash
	List(key string) List
	IntList(key string) IntList
	FloatList(key string) FloatList
	Set(key string) Set
	IntSet(key string) IntSet
	FloatSet(key string) FloatSet
	SortedSet(key string) SortedSet
	SortedIntSet(key string) SortedIntSet
	SortedFloatSet(key string) SortedFloatSet
	Mutex(key string) Mutex
	Semaphore(key string, count int) Mutex
	ReadWriteMutex(key string, readers int) *ReadWriteMutex
	Channel(key string) Channel
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

func (this *prefix) Hash(key string) Hash {
	return this.parent.Hash(this.root + key)
}

func (this *prefix) List(key string) List {
	return this.parent.List(this.root + key)
}

func (this *prefix) IntList(key string) IntList {
	return this.parent.IntList(this.root + key)
}

func (this *prefix) FloatList(key string) FloatList {
	return this.parent.FloatList(this.root + key)
}

func (this *prefix) Set(key string) Set {
	return this.parent.Set(this.root + key)
}

func (this *prefix) IntSet(key string) IntSet {
	return this.parent.IntSet(this.root + key)
}

func (this *prefix) FloatSet(key string) FloatSet {
	return this.parent.FloatSet(this.root + key)
}

func (this *prefix) SortedSet(key string) SortedSet {
	return this.parent.SortedSet(this.root + key)
}

func (this *prefix) SortedIntSet(key string) SortedIntSet {
	return this.parent.SortedIntSet(this.root + key)
}

func (this *prefix) SortedFloatSet(key string) SortedFloatSet {
	return this.parent.SortedFloatSet(this.root + key)
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
