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

//LLEN command -
//Length returns the number of items in this list
func (this List) Length() <-chan int {
	return IntCommand(this, this.args("llen")...)
}

//LPUSH command -
//LeftPush pushes an item onto the left side of this list
func (this List) LeftPush(items ...string) <-chan int {
	return IntCommand(this, this.args("lpush", items...)...)
}

//LPUSHX command -
//LeftPushIfExists pushes an item onto the left side of this list, but only if this list already exists
func (this List) LeftPushIfExists(item string) <-chan int {
	return IntCommand(this, this.args("lpushx", item)...)
}

//RPUSH command -
//RightPush pushes an item onto the right side of this list
func (this List) RightPush(items ...string) <-chan int {
	return IntCommand(this, this.args("rpush", items...)...)
}

//RPUSHX command -
//RightPushIfExists pushes an item onto the right side of this list, but only if this list already exists
func (this List) RightPushIfExists(item string) <-chan int {
	return IntCommand(this, this.args("rpushx", item)...)
}

//LPOP command -
//LeftPop pops an item from the left side of this list and returns it.
//If this list does not have anything in it, nothing is returned
func (this List) LeftPop() <-chan string {
	return StringCommand(this, this.args("lpop")...)
}

//BLPOP command -
//BlockUntilLeftPop pops an item from the left side of this list and returns it.
//If this list does not have anything in it, will wait until it does
func (this List) BlockUntilLeftPop() <-chan string {
	return this.BlockUntilLeftPopWithTimeout(0)
}

//BLPOP command -
//BlockUntilLeftPopWithTimeout pops an item from the left side of this list and returns it.
//If this list does not have anything in it, will wait up to "timeout" seconds for something to enter the list
func (this List) BlockUntilLeftPopWithTimeout(timeout int) <-chan string {
	return stringChannel(SliceCommand(this, this.args("blpop", itoa(timeout))...), 1)
}

//RPOP command -
//RightPop pops an item from the right side of this list and returns it.
//If this list does not have anything in it, nothing is returned
func (this List) RightPop() <-chan string {
	return StringCommand(this, this.args("rpop")...)
}

//BRPOP command -
//BlockUntilRightPop pops an item from the right side of this list and returns it.
//If this list does not have anything in it, will wait until it does
func (this List) BlockUntilRightPop() <-chan string {
	return this.BlockUntilRightPopWithTimeout(0)
}

//BRPOP command -
//BlockUntilRightPopWIthTimeout pops an item from the right side of this list and returns it.
//If this list does not have anything in it, will wait up to "timeout" seconds for something to enter the list
func (this List) BlockUntilRightPopWithTimeout(timeout int) <-chan string {
	return stringChannel(SliceCommand(this, this.args("brpop", itoa(timeout))...), 1)
}

//LINDEX command -
//Index returns the item at the specified index:
//negative numbers index from the right, with -1 being the rightmost index;
//non-negative numbers index from the left, with 0 being the leftmost index
func (this List) Index(index int) <-chan string {
	return StringCommand(this, this.args("lindex", itoa(index))...)
}

//LREM command -
//Remove removes all instances of all instances within items
func (this List) Remove(items ...string) <-chan int {
	return IntCommand(this, this.args("lrem", append([]string{"0"}, items...)...)...)
}

//LREM command -
//RemoveNFromLeft removes the first "n" instances of "item" from the list
func (this List) RemoveNFromLeft(n int, item string) <-chan int {
	return IntCommand(this, this.args("lrem", itoa(n), item)...)
}

//LREM command -
//RemoveNFromRight removes the last "n" instances of "item" from the list
func (this List) RemoveNFromRight(n int, item string) <-chan int {
	return IntCommand(this, this.args("lrem", itoa(-n), item)...)
}

//LSET command -
//Set sets the item at the specified index to "item":
//negative indexes index from the right with -1 being the rightmost;
//non-negative indexes index from the left with 0 being the leftmost
func (this List) Set(index int, item string) <-chan nothing {
	return NilCommand(this, this.args("lset", itoa(index), item)...)
}

//LINSERT BEFORE command -
//InsertBefore inserts an item before a pivot
func (this List) InsertBefore(pivot, item string) <-chan int {
	return IntCommand(this, this.args("linsert", "BEFORE", pivot, item)...)
}

//LINSERT AFTER command -
//InsertAfter inserts an item after a pivot
func (this List) InsertAfter(pivot, item string) <-chan int {
	return IntCommand(this, this.args("linsert", "AFTER", pivot, item)...)
}

//LRANGE command -
//GetFromRange returns all items from between two indices:
//negative indexes index from the right with -1 being the rightmost;
//non-negative indexes index from the left with 0 being the leftmost
func (this List) GetFromRange(left, right int) <-chan []string {
	return SliceCommand(this, this.args("lrange", itoa(left), itoa(right))...)
}

//LTRIM command -
//TrimToRange removes all items not within the two indices:
//negative indexes index from the right with -1 being the rightmost;
//non-negative indexes index from the left with 0 being the leftmost
func (this List) TrimToRange(left, right int) <-chan nothing {
	return NilCommand(this, this.args("ltrim", itoa(left), itoa(right))...)
}

//RPOPLPUSH command -
//MoveLastItemToList moves the last item on this list to the front of a new list.
//If nothing is in this list, nothing happens
func (this List) MoveLastItemToList(newList List) <-chan string {
	return StringCommand(this, this.args("rpoplpush", newList.key)...)
}

//BRPOPLPUSH command -
//BlockUntilMoveLastItemToList moves the last item on this list to the front of a new list.
//If nothing is in this list, will wait until something is
func (this List) BlockUntilMoveLastItemToList(newList List) <-chan string {
	return this.BlockUntilMoveLastItemToListWithTimeout(newList, 0)
}

//BRPOPLPUSH command -
//BlockUntilMoveLastItemToListWithTimeout moves the last item on this list to the front of a new list.
//If nothing is in this list, will wait up to "timeout" seconds for something to be there before giving up
func (this List) BlockUntilMoveLastItemToListWithTimeout(newList List, timeout int) <-chan string {
	return StringCommand(this, this.args("brpoplpush", newList.key, itoa(timeout))...)
}

//Use allows you to use this key on a different executor
func (this List) Use(e SafeExecutor) List {
	this.client = e
	return this
}
