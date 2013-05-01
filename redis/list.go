package redis

type List struct {
	SortableKey
}

func newList(client SafeExecutor, key string) List {
	return List{
		newSortableKey(client, key),
	}
}

//IsValid returns whether the underlying redis object can use the commands in this object
func (this List) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "list")
	}()
	return c
}

func (this List) Length() <-chan int {
	return IntCommand(this, this.args("llen"))
}

func (this List) LeftPush(items ...string) <-chan int {
	return IntCommand(this, this.args("lpush", items...))
}

func (this List) LeftPushIfExists(item string) <-chan int {
	return IntCommand(this, this.args("lpushx", item))
}

func (this List) RightPush(items ...string) <-chan int {
	return IntCommand(this, this.args("rpush", items...))
}

func (this List) RightPushIfExists(item string) <-chan int {
	return IntCommand(this, this.args("rpushx", item))
}

func (this List) LeftPop() <-chan string {
	return StringCommand(this, this.args("lpop"))
}

func (this List) BlockUntilLeftPop() <-chan string {
	return this.BlockUntilLeftPopWithTimeout(0)
}

func (this List) BlockUntilLeftPopWithTimeout(timeout int) <-chan string {
	output := SliceCommand(this, this.args("blpop", itoa(timeout)))
	realoutput := make(chan string, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			realoutput <- slice[1]
		}
	}()
	return realoutput
}

func (this List) RightPop() <-chan string {
	return StringCommand(this, this.args("rpop"))
}

func (this List) BlockUntilRightPop() <-chan string {
	return this.BlockUntilRightPopWithTimeout(0)
}

func (this List) BlockUntilRightPopWithTimeout(timeout int) <-chan string {
	output := SliceCommand(this, this.args("brpop", itoa(timeout)))
	realoutput := make(chan string, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			realoutput <- slice[1]
		}
	}()
	return realoutput
}

func (this List) Index(index int) <-chan string {
	return StringCommand(this, this.args("lindex", itoa(index)))
}
func (this List) Remove(item string) <-chan int {
	return IntCommand(this, this.args("lrem", "0", item))
}
func (this List) RemoveNFromLeft(n int, item string) <-chan int {
	return IntCommand(this, this.args("lrem", itoa(n), item))
}
func (this List) RemoveNFromRight(n int, item string) <-chan int {
	return IntCommand(this, this.args("lrem", itoa(-n), item))
}
func (this List) Set(index int, item string) <-chan nothing {
	return NilCommand(this, this.args("lset", itoa(index), item))
}
func (this List) InsertBefore(pivot, item string) <-chan int {
	return IntCommand(this, this.args("linsert", "BEFORE", pivot, item))
}
func (this List) InsertAfter(pivot, item string) <-chan int {
	return IntCommand(this, this.args("linsert", "AFTER", pivot, item))
}
func (this List) GetFromRange(left, right int) <-chan []string {
	return SliceCommand(this, this.args("lrange", itoa(left), itoa(right)))
}
func (this List) TrimToRange(left, right int) <-chan nothing {
	return NilCommand(this, this.args("ltrim", itoa(left), itoa(right)))
}
func (this List) MoveLastItemToList(newList List) <-chan string {
	return StringCommand(this, this.args("rpoplpush", newList.key))
}
func (this List) BlockUntilMoveLastItemToList(newList List) <-chan string {
	return this.BlockUntilMoveLastItemToListWithTimeout(newList, 0)
}

func (this List) BlockUntilMoveLastItemToListWithTimeout(newList List, timeout int) <-chan string {
	return StringCommand(this, this.args("brpoplpush", newList.key, itoa(timeout)))
}

//Use allows you to use this key on a different executor
func (this List) Use(e SafeExecutor) List {
	this.client = e
	return this
}
