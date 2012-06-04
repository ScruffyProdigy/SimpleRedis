package redis

type IntList struct {
	Key
}

func newIntList(client Executor, key string) IntList {
	return IntList{
		newKey(client, key),
	}
}

func (this IntList) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "list")
	}()
	return c
}

//
func (this IntList) Length() <-chan int {
	command, output := newIntCommand(this.args("llen"))
	this.Execute(command)
	return output
}

func (this IntList) LeftPush(items ...int) <-chan int {
	command, output := newIntCommand(this.args("lpush", intsToStrings(items)...))
	this.Execute(command)
	return output
}

func (this IntList) LeftPushIfExists(item int) <-chan int {
	command, output := newIntCommand(this.args("lpushx", itoa(item)))
	this.Execute(command)
	return output
}

func (this IntList) RightPush(items ...int) <-chan int {
	command, output := newIntCommand(this.args("rpush", intsToStrings(items)...))
	this.Execute(command)
	return output
}

func (this IntList) RightPushIfExists(item int) <-chan int {
	command, output := newIntCommand(this.args("rpushx", itoa(item)))
	this.Execute(command)
	return output
}

func (this IntList) LeftPop() <-chan int {
	command, output := newIntCommand(this.args("lpop"))
	this.Execute(command)
	return output
}

//perhaps allow these commands to take extra lists
//or figure out how to just return the value, not the key
func (this IntList) BlockUntilLeftPop() <-chan int {
	return this.BlockUntilLeftPopWithTimeout(0)
}

func (this IntList) BlockUntilLeftPopWithTimeout(timeout int) <-chan int {
	command, output := newSliceCommand(this.args("blpop", itoa(timeout)))
	this.Execute(command)
	realoutput := make(chan int, 1)
	go func() {
		if slice, ok := <-output; ok {
			realoutput <- atoi(slice[1])
		}
		close(realoutput)
	}()
	return realoutput
}

func (this IntList) RightPop() <-chan int {
	command, output := newIntCommand(this.args("rpop"))
	this.Execute(command)
	return output
}

//perhaps allow these commands to take extra lists
//or figure out how to just return the value, not the key
func (this IntList) BlockUntilRightPop() <-chan int {
	return this.BlockUntilRightPopWithTimeout(0)
}

func (this IntList) BlockUntilRightPopWithTimeout(timeout int) <-chan int {
	command, output := newSliceCommand(this.args("brpop", itoa(timeout)))
	this.Execute(command)
	realoutput := make(chan int, 1)
	go func() {
		if slice, ok := <-output; ok {
			realoutput <- atoi(slice[1])
		}
		close(realoutput)
	}()
	return realoutput
}

func (this IntList) Index(index int) <-chan int {
	command, output := newIntCommand(this.args("lindex", itoa(index)))
	this.Execute(command)
	return output
}

func (this IntList) Remove(item ...int) <-chan int {
	command, output := newIntCommand(this.args("lrem", append([]string{"0"}, intsToStrings(item)...)...))
	this.Execute(command)
	return output
}

func (this IntList) RemoveNFromLeft(n int, item int) <-chan int {
	command, output := newIntCommand(this.args("lrem", itoa(n), itoa(item)))
	this.Execute(command)
	return output
}

func (this IntList) RemoveNFromRight(n int, item int) <-chan int {
	command, output := newIntCommand(this.args("lrem", itoa(-n), itoa(item)))
	this.Execute(command)
	return output
}

func (this IntList) Set(index int, item int) <-chan nothing {
	command, output := newNilCommand(this.args("lset", itoa(index), itoa(item)))
	this.Execute(command)
	return output
}

func (this IntList) InsertBefore(pivot, item int) <-chan int {
	command, output := newIntCommand(this.args("linsert", "BEFORE", itoa(pivot), itoa(item)))
	this.Execute(command)
	return output
}

func (this IntList) InsertAfter(pivot, item int) <-chan int {
	command, output := newIntCommand(this.args("linsert", "AFTER", itoa(pivot), itoa(item)))
	this.Execute(command)
	return output
}

func (this IntList) GetFromRange(left, right int) <-chan []int {
	command, output := newSliceCommand(this.args("lrange", itoa(left), itoa(right)))
	this.Execute(command)
	realoutput := make(chan []int, 1)
	go func() {
		if slice, ok := <-output; ok {
			realoutput <- stringsToInts(slice)
		}
		close(realoutput)
	}()
	return realoutput
}

func (this IntList) TrimToRange(left, right int) <-chan nothing {
	command, output := newNilCommand(this.args("ltrim", itoa(left), itoa(right)))
	this.Execute(command)
	return output
}

func (this IntList) MoveLastItemToList(newList IntList) <-chan int {
	command, output := newIntCommand(this.args("rpoplpush", newList.key))
	this.Execute(command)
	return output
}

func (this IntList) BlockUntilMoveLastItemToList(newList IntList) <-chan int {
	return this.BlockUntilMoveLastItemToListWithTimeout(newList, 0)
}

func (this IntList) BlockUntilMoveLastItemToListWithTimeout(newList IntList, timeout int) <-chan int {
	command, output := newIntCommand(this.args("brpoplpush", newList.key, itoa(timeout)))
	this.Execute(command)
	return output
}

func (this IntList) Sort() Sorter {
	return Sorter{key: this.Key}
}

func (this IntList) Use(e Executor) IntList {
	this.client = e
	return this
}
