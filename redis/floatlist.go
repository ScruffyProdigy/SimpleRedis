package redis

type FloatList struct {
	SortableKey
}

func newFloatList(client SafeExecutor, key string) FloatList {
	return FloatList{
		newSortableKey(client, key),
	}
}

func (this FloatList) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "list")
	}()
	return c
}

//
func (this FloatList) Length() <-chan int {
	return IntCommand(this, this.args("llen"))
}

func (this FloatList) LeftPush(items ...float64) <-chan float64 {
	return FloatCommand(this, this.args("lpush", floatsToStrings(items)...))
}

func (this FloatList) LeftPushIfExists(item float64) <-chan float64 {
	return FloatCommand(this, this.args("lpushx", ftoa(item)))
}

func (this FloatList) RightPush(items ...float64) <-chan float64 {
	return FloatCommand(this, this.args("rpush", floatsToStrings(items)...))
}

func (this FloatList) RightPushIfExists(item float64) <-chan float64 {
	return FloatCommand(this, this.args("rpushx", ftoa(item)))
}

func (this FloatList) LeftPop() <-chan float64 {
	return FloatCommand(this, this.args("lpop"))
}

//perhaps allow these commands to take extra lists
//or figure out how to just return the value, not the key
func (this FloatList) BlockUntilLeftPop() <-chan float64 {
	return this.BlockUntilLeftPopWithTimeout(0)
}

func (this FloatList) BlockUntilLeftPopWithTimeout(timeout int) <-chan float64 {
	output := SliceCommand(this, this.args("blpop", itoa(timeout)))
	realoutput := make(chan float64, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			f, err := atof(slice[1])
			if err != nil {
				this.client.errCallback(err, "blpop")
				return
			}
			realoutput <- f
		}
	}()
	return realoutput
}

func (this FloatList) RightPop() <-chan float64 {
	return FloatCommand(this, this.args("rpop"))
}

//perhaps allow these commands to take extra lists
//or figure out how to just return the value, not the key
func (this FloatList) BlockUntilRightPop() <-chan float64 {
	return this.BlockUntilRightPopWithTimeout(0)
}

func (this FloatList) BlockUntilRightPopWithTimeout(timeout int) <-chan float64 {
	output := SliceCommand(this, this.args("brpop", itoa(timeout)))
	realoutput := make(chan float64, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			f, err := atof(slice[1])
			if err != nil {
				this.client.errCallback(err, "brpop")
				return
			}
			realoutput <- f
		}
	}()
	return realoutput
}

func (this FloatList) Index(index int) <-chan float64 {
	return FloatCommand(this, this.args("lindex", itoa(index)))
}

func (this FloatList) Remove(item ...float64) <-chan int {
	return IntCommand(this, this.args("lrem", append([]string{"0"}, floatsToStrings(item)...)...))
}

func (this FloatList) RemoveNFromLeft(n int, item float64) <-chan int {
	return IntCommand(this, this.args("lrem", itoa(n), ftoa(item)))
}

func (this FloatList) RemoveNFromRight(n int, item float64) <-chan int {
	return IntCommand(this, this.args("lrem", itoa(-n), ftoa(item)))
}

func (this FloatList) Set(index int, item float64) <-chan nothing {
	return NilCommand(this, this.args("lset", itoa(index), ftoa(item)))
}

func (this FloatList) InsertBefore(pivot, item float64) <-chan int {
	return IntCommand(this, this.args("linsert", "BEFORE", ftoa(pivot), ftoa(item)))
}

func (this FloatList) InsertAfter(pivot, item float64) <-chan int {
	return IntCommand(this, this.args("linsert", "AFTER", ftoa(pivot), ftoa(item)))
}

func (this FloatList) GetFromRange(left, right int) <-chan []float64 {
	output := SliceCommand(this, this.args("lrange", itoa(left), itoa(right)))
	realoutput := make(chan []float64, 1)
	go func() {
		defer close(realoutput)
		if slice, ok := <-output; ok {
			floats, err := stringsToFloats(slice)
			if err != nil {
				this.client.errCallback(err, "lrange")
				return
			}
			realoutput <- floats
		}
	}()
	return realoutput
}

func (this FloatList) TrimToRange(left, right int) <-chan nothing {
	return NilCommand(this, this.args("ltrim", itoa(left), itoa(right)))
}

func (this FloatList) MoveLastItemToList(newList FloatList) <-chan float64 {
	return FloatCommand(this, this.args("rpoplpush", newList.key))
}

func (this FloatList) BlockUntilMoveLastItemToList(newList FloatList) <-chan float64 {
	return this.BlockUntilMoveLastItemToListWithTimeout(newList, 0)
}

func (this FloatList) BlockUntilMoveLastItemToListWithTimeout(newList FloatList, timeout int) <-chan float64 {
	return FloatCommand(this, this.args("brpoplpush", newList.key, itoa(timeout)))
}

func (this FloatList) Use(e SafeExecutor) FloatList {
	this.client = e
	return this
}
