package redis

type List struct {
	SortableKey
}

func newList(client Executor, key string) List {
	return List{
		newSortableKey(client, key),
	}
}

func (this List) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "list")
	}()
	return c
}

func (this List) Length() <-chan int {
	command, output := newIntCommand(this.args("llen"))
	this.Execute(command)
	return output
}

func (this List) LeftPush(items ...string) <-chan int {
	command, output := newIntCommand(this.args("lpush", items...))
	this.Execute(command)
	return output
}

func (this List) LeftPushIfExists(item string) <-chan int {
	command, output := newIntCommand(this.args("lpushx", item))
	this.Execute(command)
	return output
}

func (this List) RightPush(items ...string) <-chan int {
	command, output := newIntCommand(this.args("rpush", items...))
	this.Execute(command)
	return output
}

func (this List) RightPushIfExists(item string) <-chan int {
	command, output := newIntCommand(this.args("rpushx", item))
	this.Execute(command)
	return output
}

func (this List) LeftPop() <-chan string {
	command, output := newStringCommand(this.args("lpop"))
	this.Execute(command)
	return output
}

func (this List) BlockUntilLeftPop() <-chan string {
	return this.BlockUntilLeftPopWithTimeout(0)
}

func (this List) BlockUntilLeftPopWithTimeout(timeout int) <-chan string {
	command, output := newSliceCommand(this.args("blpop", itoa(timeout)))
	this.Execute(command)
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
	command, output := newStringCommand(this.args("rpop"))
	this.Execute(command)
	return output
}

func (this List) BlockUntilRightPop() <-chan string {
	return this.BlockUntilRightPopWithTimeout(0)
}

func (this List) BlockUntilRightPopWithTimeout(timeout int) <-chan string {
	command, output := newSliceCommand(this.args("brpop", itoa(timeout)))
	this.Execute(command)
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
	command, output := newStringCommand(this.args("lindex", itoa(index)))
	this.Execute(command)
	return output
}

func (this List) Remove(item string) <-chan int {
	command, output := newIntCommand(this.args("lrem", "0", item))
	this.Execute(command)
	return output
}

func (this List) RemoveNFromLeft(n int, item string) <-chan int {
	command, output := newIntCommand(this.args("lrem", itoa(n), item))
	this.Execute(command)
	return output
}

func (this List) RemoveNFromRight(n int, item string) <-chan int {
	command, output := newIntCommand(this.args("lrem", itoa(-n), item))
	this.Execute(command)
	return output
}

func (this List) Set(index int, item string) <-chan nothing {
	command, output := newNilCommand(this.args("lset", itoa(index), item))
	this.Execute(command)
	return output
}

func (this List) InsertBefore(pivot, item string) <-chan int {
	command, output := newIntCommand(this.args("linsert", "BEFORE", pivot, item))
	this.Execute(command)
	return output
}

func (this List) InsertAfter(pivot, item string) <-chan int {
	command, output := newIntCommand(this.args("linsert", "AFTER", pivot, item))
	this.Execute(command)
	return output
}

func (this List) GetFromRange(left, right int) <-chan []string {
	command, output := newSliceCommand(this.args("lrange", itoa(left), itoa(right)))
	this.Execute(command)
	return output
}

func (this List) TrimToRange(left, right int) <-chan nothing {
	command, output := newNilCommand(this.args("ltrim", itoa(left), itoa(right)))
	this.Execute(command)
	return output
}

func (this List) MoveLastItemToList(newList List) <-chan string {
	command, output := newStringCommand(this.args("rpoplpush", newList.key))
	this.Execute(command)
	return output
}

func (this List) BlockUntilMoveLastItemToList(newList List) <-chan string {
	return this.BlockUntilMoveLastItemToListWithTimeout(newList, 0)
}

func (this List) BlockUntilMoveLastItemToListWithTimeout(newList List, timeout int) <-chan string {
	command, output := newStringCommand(this.args("brpoplpush", newList.key, itoa(timeout)))
	this.Execute(command)
	return output
}

func (this List) Use(e Executor) List {
	this.client = e
	return this
}
