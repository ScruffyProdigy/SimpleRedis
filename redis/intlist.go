package redis

//TODO: refactor to reuse List code

//IntList implements the Redis List primitive assuming all inputs are ints (which is useful for indexes)
//See http://redis.io/commands#list for more information on Redis Lists
type IntList struct {
	SortableKey
}

func newIntList(client SafeExecutor, key string) IntList {
	return IntList{
		newSortableKey(client, key),
	}
}

//IsValid returns whether the underlying redis object can use the commands in this object
func (this IntList) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "list")
	}()
	return c
}

//LLEN command - 
//Length returns the number of items in this list
func (this IntList) Length() <-chan int {
	return IntCommand(this, this.args("llen")...)
}

//LPUSH command - 
//LeftPush pushes an integer onto the left side of this list
func (this IntList) LeftPush(items ...int) <-chan int {
	return IntCommand(this, this.args("lpush", intsToStrings(items)...)...)
}

//LPUSHX command - 
//LeftPushIfExists pushes an integer onto the left side of the list if the list exists
func (this IntList) LeftPushIfExists(item int) <-chan int {
	return IntCommand(this, this.args("lpushx", itoa(item))...)
}

//RPUSH command - 
//RightPush pushes an integer onto the right side of this list - RPUSH command
func (this IntList) RightPush(items ...int) <-chan int {
	return IntCommand(this, this.args("rpush", intsToStrings(items)...)...)
}

//RPUSHX command -
//RightPushIfExists pushes an integer onto the right side of this list if the list exists - RPUSHX command
func (this IntList) RightPushIfExists(item int) <-chan int {
	return IntCommand(this, this.args("rpushx", itoa(item))...)
}

//LPOP command -
//LeftPop pops the leftmost integer off of the list and returns it.
//If there is nothing in the list, it returns nothing
func (this IntList) LeftPop() <-chan int {
	return IntCommand(this, this.args("lpop")...)
}

//BLPOP command -
//BlockUntilLeftPop pops the leftmost integer off of the list and returns it.
//If there is nothing in the list, it will wait until something gets placed in the list
func (this IntList) BlockUntilLeftPop() <-chan int {
	return this.BlockUntilLeftPopWithTimeout(0)
}

//BLPOP command -
//BlockUntilLeftPopWithTimeout pops the leftmost integer off of the list and returns it.
//If there is nothing in the list, it will wait up to "timeout" seconds for something to be placed in the list
func (this IntList) BlockUntilLeftPopWithTimeout(timeout int) <-chan int {
	return intChannel(SliceCommand(this, this.args("blpop", itoa(timeout))...), 1)
}

//RPOP command -
//RightPop pops the rightmost integer off of the list and returns it.
//If there is nothing in the list, it returns nothing
func (this IntList) RightPop() <-chan int {
	return IntCommand(this, this.args("rpop")...)
}

//BRPOP command -
//BlockUntilRightPop pops the rightmost integer off of the list and returns it.
//If there is nothing in the list, it will wait for something to be placed in the list
func (this IntList) BlockUntilRightPop() <-chan int {
	return this.BlockUntilRightPopWithTimeout(0)
}

//BRPOP command -
//BlockUntilRightPopWithTimeout pops the rightmost integer off of the list and returns it.
//If there is nothing in the list, it will wait up to "timeout" seconds for something to be placed in it
func (this IntList) BlockUntilRightPopWithTimeout(timeout int) <-chan int {
	return intChannel(SliceCommand(this, this.args("brpop", itoa(timeout))...), 1)
}

//LINDEX command -
//Index returns the integer waiting at the specified index:
//negative indexes index from the right with -1 being the rightmost;
//non-negative indexes index from the left with 0 being the leftmost
func (this IntList) Index(index int) <-chan int {
	return IntCommand(this, this.args("lindex", itoa(index))...)
}

//LREM command -
//Remove removes all instances of all instances within items
func (this IntList) Remove(items ...int) <-chan int {
	return IntCommand(this, this.args("lrem", append([]string{"0"}, intsToStrings(items)...)...)...)
}

//LREM command -
//Removes the first "n" instances of "item" from the list
func (this IntList) RemoveNFromLeft(n int, item int) <-chan int {
	return IntCommand(this, this.args("lrem", itoa(n), itoa(item))...)
}

//LREM command -
//Removes the last "n" instances of "item" from the list
func (this IntList) RemoveNFromRight(n int, item int) <-chan int {
	return IntCommand(this, this.args("lrem", itoa(-n), itoa(item))...)
}

//LSET command -
//Set sets the integer at the specified index to "item":
//negative indexes index from the right with -1 being the rightmost;
//non-negative indexes index from the left with 0 being the leftmost
func (this IntList) Set(index int, item int) <-chan nothing {
	return NilCommand(this, this.args("lset", itoa(index), itoa(item))...)
}

//LINSERT command -
//InsertBefore inserts an item before a pivot
func (this IntList) InsertBefore(pivot, item int) <-chan int {
	return IntCommand(this, this.args("linsert", "BEFORE", itoa(pivot), itoa(item))...)
}

//LINSERT command -
//InsertAbove inserts an item after a pivot
func (this IntList) InsertAfter(pivot, item int) <-chan int {
	return IntCommand(this, this.args("linsert", "AFTER", itoa(pivot), itoa(item))...)
}

//LRANGE command -
//GetFromRange returns all items from between two indices:
//negative indexes index from the right with -1 being the rightmost;
//non-negative indexes index from the left with 0 being the leftmost
func (this IntList) GetFromRange(left, right int) <-chan []int {
	return intsChannel(SliceCommand(this, this.args("lrange", itoa(left), itoa(right))...))
}

//LTRIM command -
//TrimToRange removes all items not within the two indices:
//negative indexes index from the right with -1 being the rightmost;
//non-negative indexes index from the left with 0 being the leftmost
func (this IntList) TrimToRange(left, right int) <-chan nothing {
	return NilCommand(this, this.args("ltrim", itoa(left), itoa(right))...)
}

//RPOPLPUSH command -
//MoveLastItemToList moves the last item on this list to the front of a new list.
//If nothing is in this list, nothing happens
func (this IntList) MoveLastItemToList(newList IntList) <-chan int {
	return IntCommand(this, this.args("rpoplpush", newList.key)...)
}

//BRPOPLPUSH command -
//BlockUntilMoveLastItemToList moves the last item on this list to the front of a new list.
//If nothing is in this list, will wait until something is
func (this IntList) BlockUntilMoveLastItemToList(newList IntList) <-chan int {
	return this.BlockUntilMoveLastItemToListWithTimeout(newList, 0)
}

//BRPOPLPUSH command -
//BlockUntilMoveLastItemToListWithTimeout moves the last item on this list to the front of a new list.
//If nothing is in this list, will wait up to "timeout" seconds for something to be there before giving up
func (this IntList) BlockUntilMoveLastItemToListWithTimeout(newList IntList, timeout int) <-chan int {
	return IntCommand(this, this.args("brpoplpush", newList.key, itoa(timeout))...)
}

//Use allows you to use this key on a different executor
func (this IntList) Use(e SafeExecutor) IntList {
	this.client = e
	return this
}
