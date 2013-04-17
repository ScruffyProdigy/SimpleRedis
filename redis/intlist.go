package redis

type IntList struct {
	SortableKey
}

func newIntList(client SafeExecutor, key string) IntList {
	return IntList{
		newSortableKey(client, key),
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
	return IntCommand(this, this.args("llen"))
}

func (this IntList) LeftPush(items ...int) <-chan int {
	return IntCommand(this, this.args("lpush", intsToStrings(items)...))
}

func (this IntList) LeftPushIfExists(item int) <-chan int {
	return IntCommand(this, this.args("lpushx", itoa(item)))
}

func (this IntList) RightPush(items ...int) <-chan int {
	return IntCommand(this, this.args("rpush", intsToStrings(items)...))
}

func (this IntList) RightPushIfExists(item int) <-chan int {
	return IntCommand(this, this.args("rpushx", itoa(item)))
}

func (this IntList) LeftPop() <-chan int {
	return IntCommand(this, this.args("lpop"))
}

//perhaps allow these commands to take extra lists
//or figure out how to just return the value, not the key
func (this IntList) BlockUntilLeftPop() <-chan int {
	return this.BlockUntilLeftPopWithTimeout(0)
}

func (this IntList) BlockUntilLeftPopWithTimeout(timeout int) <-chan int {
	output := SliceCommand(this, this.args("blpop", itoa(timeout)))
	realoutput := make(chan int, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			if res, err := atoi(slice[1]); err != nil {
				this.client.ErrCallback(err, "blpop")
				return
			} else {
				realoutput <- res
			}
		}
	}()
	return realoutput
}

func (this IntList) RightPop() <-chan int {
	return IntCommand(this, this.args("rpop"))
}

//perhaps allow these commands to take extra lists
//or figure out how to just return the value, not the key
func (this IntList) BlockUntilRightPop() <-chan int {
	return this.BlockUntilRightPopWithTimeout(0)
}

func (this IntList) BlockUntilRightPopWithTimeout(timeout int) <-chan int {
	output := SliceCommand(this, this.args("brpop", itoa(timeout)))
	realoutput := make(chan int, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			if res, err := atoi(slice[1]); err != nil {
				this.client.ErrCallback(err, "brpop")
			} else {
				realoutput <- res
			}
		}
	}()
	return realoutput
}

func (this IntList) Index(index int) <-chan int {
	return IntCommand(this, this.args("lindex", itoa(index)))
}

func (this IntList) Remove(item ...int) <-chan int {
	return IntCommand(this, this.args("lrem", append([]string{"0"}, intsToStrings(item)...)...))
}

func (this IntList) RemoveNFromLeft(n int, item int) <-chan int {
	return IntCommand(this, this.args("lrem", itoa(n), itoa(item)))
}

func (this IntList) RemoveNFromRight(n int, item int) <-chan int {
	return IntCommand(this, this.args("lrem", itoa(-n), itoa(item)))
}

func (this IntList) Set(index int, item int) <-chan nothing {
	return NilCommand(this, this.args("lset", itoa(index), itoa(item)))
}

func (this IntList) InsertBefore(pivot, item int) <-chan int {
	return IntCommand(this, this.args("linsert", "BEFORE", itoa(pivot), itoa(item)))
}

func (this IntList) InsertAfter(pivot, item int) <-chan int {
	return IntCommand(this, this.args("linsert", "AFTER", itoa(pivot), itoa(item)))
}

func (this IntList) GetFromRange(left, right int) <-chan []int {
	output := SliceCommand(this, this.args("lrange", itoa(left), itoa(right)))
	realoutput := make(chan []int, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			if ints, err := stringsToInts(slice); err != nil {
				this.client.ErrCallback(err, "lrange")
			} else {
				realoutput <- ints
			}
		}
	}()
	return realoutput
}

func (this IntList) TrimToRange(left, right int) <-chan nothing {
	return NilCommand(this, this.args("ltrim", itoa(left), itoa(right)))
}

func (this IntList) MoveLastItemToList(newList IntList) <-chan int {
	return IntCommand(this, this.args("rpoplpush", newList.key))
}

func (this IntList) BlockUntilMoveLastItemToList(newList IntList) <-chan int {
	return this.BlockUntilMoveLastItemToListWithTimeout(newList, 0)
}

func (this IntList) BlockUntilMoveLastItemToListWithTimeout(newList IntList, timeout int) <-chan int {
	return IntCommand(this, this.args("brpoplpush", newList.key, itoa(timeout)))
}

func (this IntList) Use(e SafeExecutor) IntList {
	this.client = e
	return this
}
