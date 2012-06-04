package redis

type FloatList struct {
	Key
}

func newFloatList(client Executor, key string) FloatList {
	return FloatList{
		newKey(client, key),
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
	command, output := newIntCommand(this.args("llen"))
	this.Execute(command)
	return output
}

func (this FloatList) LeftPush(items ...float64) <-chan float64 {
	command, output := newFloatCommand(this.args("lpush", floatsToStrings(items)...))
	this.Execute(command)
	return output
}

func (this FloatList) LeftPushIfExists(item float64) <-chan float64 {
	command, output := newFloatCommand(this.args("lpushx", ftoa(item)))
	this.Execute(command)
	return output
}

func (this FloatList) RightPush(items ...float64) <-chan float64 {
	command, output := newFloatCommand(this.args("rpush", floatsToStrings(items)...))
	this.Execute(command)
	return output
}

func (this FloatList) RightPushIfExists(item float64) <-chan float64 {
	command, output := newFloatCommand(this.args("rpushx", ftoa(item)))
	this.Execute(command)
	return output
}

func (this FloatList) LeftPop() <-chan float64 {
	command, output := newFloatCommand(this.args("lpop"))
	this.Execute(command)
	return output
}

//perhaps allow these commands to take extra lists
//or figure out how to just return the value, not the key
func (this FloatList) BlockUntilLeftPop() <-chan float64 {
	return this.BlockUntilLeftPopWithTimeout(0)
}

func (this FloatList) BlockUntilLeftPopWithTimeout(timeout int) <-chan float64 {
	command, output := newSliceCommand(this.args("blpop", itoa(timeout)))
	this.Execute(command)
	realoutput := make(chan float64, 1)
	go func() {
		if slice, ok := <-output; ok {
			realoutput <- atof(slice[1])
		}
		close(realoutput)
	}()
	return realoutput
}

func (this FloatList) RightPop() <-chan float64 {
	command, output := newFloatCommand(this.args("rpop"))
	this.Execute(command)
	return output
}

//perhaps allow these commands to take extra lists
//or figure out how to just return the value, not the key
func (this FloatList) BlockUntilRightPop() <-chan float64 {
	return this.BlockUntilRightPopWithTimeout(0)
}

func (this FloatList) BlockUntilRightPopWithTimeout(timeout int) <-chan float64 {
	command, output := newSliceCommand(this.args("brpop", itoa(timeout)))
	this.Execute(command)
	realoutput := make(chan float64, 1)
	go func() {
		if slice, ok := <-output; ok {
			realoutput <- atof(slice[1])
		}
		close(realoutput)
	}()
	return realoutput
}

func (this FloatList) Index(index int) <-chan float64 {
	command, output := newFloatCommand(this.args("lindex", itoa(index)))
	this.Execute(command)
	return output
}

func (this FloatList) Remove(item ...float64) <-chan int {
	command, output := newIntCommand(this.args("lrem", append([]string{"0"}, floatsToStrings(item)...)...))
	this.Execute(command)
	return output
}

func (this FloatList) RemoveNFromLeft(n int, item float64) <-chan int {
	command, output := newIntCommand(this.args("lrem", itoa(n), ftoa(item)))
	this.Execute(command)
	return output
}

func (this FloatList) RemoveNFromRight(n int, item float64) <-chan int {
	command, output := newIntCommand(this.args("lrem", itoa(-n), ftoa(item)))
	this.Execute(command)
	return output
}

func (this FloatList) Set(index int, item float64) <-chan nothing {
	command, output := newNilCommand(this.args("lset", itoa(index), ftoa(item)))
	this.Execute(command)
	return output
}

func (this FloatList) InsertBefore(pivot, item float64) <-chan int {
	command, output := newIntCommand(this.args("linsert", "BEFORE", ftoa(pivot), ftoa(item)))
	this.Execute(command)
	return output
}

func (this FloatList) InsertAfter(pivot, item float64) <-chan int {
	command, output := newIntCommand(this.args("linsert", "AFTER", ftoa(pivot), ftoa(item)))
	this.Execute(command)
	return output
}

func (this FloatList) GetFromRange(left, right int) <-chan []float64 {
	command, output := newSliceCommand(this.args("lrange", itoa(left), itoa(right)))
	this.Execute(command)
	realoutput := make(chan []float64, 1)
	go func() {
		if slice, ok := <-output; ok {
			realoutput <- stringsToFloats(slice)
		}
		close(realoutput)
	}()
	return realoutput
}

func (this FloatList) TrimToRange(left, right int) <-chan nothing {
	command, output := newNilCommand(this.args("ltrim", itoa(left), itoa(right)))
	this.Execute(command)
	return output
}

func (this FloatList) MoveLastItemToList(newList FloatList) <-chan float64 {
	command, output := newFloatCommand(this.args("rpoplpush", newList.key))
	this.Execute(command)
	return output
}

func (this FloatList) BlockUntilMoveLastItemToList(newList FloatList) <-chan float64 {
	return this.BlockUntilMoveLastItemToListWithTimeout(newList, 0)
}

func (this FloatList) BlockUntilMoveLastItemToListWithTimeout(newList FloatList, timeout int) <-chan float64 {
	command, output := newFloatCommand(this.args("brpoplpush", newList.key, itoa(timeout)))
	this.Execute(command)
	return output
}

func (this FloatList) Sort() Sorter {
	return Sorter{key: this.Key}
}

func (this FloatList) Use(e Executor) FloatList {
	this.client = e
	return this
}
