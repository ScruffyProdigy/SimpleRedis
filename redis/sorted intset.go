package redis

//TODO: refactor to use SortedSet code

//SortedIntSet is an object which implements the Redis ZSet Primitive assume all inputs are ints (which is useful for indexes)
//See http://redis.io/commands#sorted_set for more info on ZSets
type SortedIntSet struct {
	SortableKey
}

func newSortedIntSet(client SafeExecutor, key string) SortedIntSet {
	return SortedIntSet{
		newSortableKey(client, key),
	}
}

//IsValid returns whether the underlying redis object can use the commands in this object
func (this SortedIntSet) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "zset")
	}()
	return c
}

//Add adds an integer to a zset or updates its score if it already exists
//returns true when adding, false when updating
func (this SortedIntSet) Add(item int, score float64) <-chan bool {
	return BoolCommand(this, this.args("zadd", ftoa(score), itoa(item))...)
}

//IncrementBy adjusts the score of the member within the zset
//returns the new score
func (this SortedIntSet) IncrementBy(item int, score float64) <-chan float64 {
	return FloatCommand(this, this.args("zincrby", ftoa(score), itoa(item))...)
}

//Remove removes a member from the zset if it is part of the set
//returns whether or not it was part of the set
func (this SortedIntSet) Remove(item int) <-chan bool {
	return BoolCommand(this, this.args("zrem", itoa(item))...)
}

//Size returns the number of members of the zset
func (this SortedIntSet) Size() <-chan int {
	return IntCommand(this, this.args("zcard")...)
}

//IndexOf returns the index of a member - 
//ie, the lowest ranked member would have an index of 0, and the next lowest an index of 1
func (this SortedIntSet) IndexOf(item int) <-chan int {
	return IntCommand(this, this.args("zrank", itoa(item))...)
}

//ReverseIndexOf returns the reverse index of a member - 
//ie, the highest ranked member would have an reverse index of 0, and the next highest an reverse index of 1
func (this SortedIntSet) ReverseIndexOf(item int) <-chan int {
	return IntCommand(this, this.args("zrevrank", itoa(item))...)
}

//ScoreOf returns the score associated with a given member of the zset
func (this SortedIntSet) ScoreOf(item int) <-chan float64 {
	return FloatCommand(this, this.args("zscore", itoa(item))...)
}

//IndexedBetween returns a slice of all members between the indices
func (this SortedIntSet) IndexedBetween(start, stop int) <-chan []int {
	return intsChannel(SliceCommand(this, this.args("zrange", itoa(start), itoa(stop))...))
}

//ReverseIndexedBetween returns a slice of all members between the reverse indices
func (this SortedIntSet) ReverseIndexedBetween(start, stop int) <-chan []int {
	return intsChannel(SliceCommand(this, this.args("zrevrange", itoa(start), itoa(stop))...))
}

//IndexedBetweenWithScores returns a map of all members between the indices and their associated scores
//warning: golang maps are not ordered
func (this SortedIntSet) IndexedBetweenWithScores(start, stop int) <-chan map[int]float64 {
	return intfloatMapChannel(MapCommand(this, this.args("zrange", itoa(start), itoa(stop), "WITHSCORES")...))
}

//IndexedBetweenWithScores returns a map of all members between the reverse indices and their associated scores
//warning: golang maps are not ordered
func (this SortedIntSet) ReverseIndexedBetweenWithScores(start, stop int) <-chan map[int]float64 {
	return intfloatMapChannel(MapCommand(this, this.args("zrevrange", itoa(start), itoa(stop), "WITHSCORES")...))
}

//RemoveIndexedBetween removes all members between the indices
//returns the number of members removed
func (this SortedIntSet) RemoveIndexedBetween(start, stop int) <-chan int {
	return IntCommand(this, this.args("zremrangebyrank", itoa(start), itoa(stop))...)
}

//SortedIntSetRange keeps track of all range arguments being used in a search
type SortedIntSetRange struct {
	min, max      string
	fmin, fmax    float64
	limited       bool
	offset, count int
	reversed      bool

	key Key
}

//Scores createa a SortedIntSetRange to help narrow a search to be done later
func (this SortedIntSet) Scores() *SortedIntSetRange {
	return &SortedIntSetRange{
		min: "-inf",
		max: "+inf",
		key: this.Key,
	}
}

//Above limits results to members who have a score above "min"
func (this *SortedIntSetRange) Above(min float64) *SortedIntSetRange {
	if this.min == "-inf" || this.fmin <= min {
		this.fmin = min
		this.min = "(" + ftoa(min)
	}
	return this
}

//Below limits results to members who have a score below "max"
func (this *SortedIntSetRange) Below(max float64) *SortedIntSetRange {
	if this.max == "+inf" || this.fmax >= max {
		this.fmax = max
		this.max = "(" + ftoa(max)
	}
	return this
}

//AboveOrEqualTo limits results to members who have a score above or equal to "min"
func (this *SortedIntSetRange) AboveOrEqualTo(min float64) *SortedIntSetRange {
	if this.min == "-inf" || this.fmin < min {
		this.fmin = min
		this.min = ftoa(min)
	}
	return this
}

//BelowOrEqualTo limits results to members who have a score below or equal to "max"
func (this *SortedIntSetRange) BelowOrEqualTo(max float64) *SortedIntSetRange {
	if this.max == "+inf" || this.fmax > max {
		this.fmax = max
		this.max = ftoa(max)
	}
	return this
}

//Reversed returns the results in reverse order
//only useful if getting or getting with scores; not useful for counting or removing
func (this *SortedIntSetRange) Reversed() *SortedIntSetRange {
	this.reversed = !this.reversed
	return this
}

//Limit limits the results you get back - it skips the first "offset" results, and then only returns the next "offset"
//only useful if getting or getting with scores; not useful for counting or removing
func (this *SortedIntSetRange) Limit(offset, count int) *SortedIntSetRange {
	this.limited = true
	this.offset = offset
	this.count = count
	return this
}

//Count returns the number of members that fit in the search criteria
func (this *SortedIntSetRange) Count() <-chan int {
	return IntCommand(this.key, this.key.args("zcount", this.min, this.max)...)
}

//Remove removes all members that fit the search criteria from the zset
//returns the number of members removed
func (this *SortedIntSetRange) Remove() <-chan int {
	return IntCommand(this.key, this.key.args("zremrangebyscore", this.min, this.max)...)
}

//Get returns a list of all members fitting the search criteria
func (this *SortedIntSetRange) Get() <-chan []int {
	op := "zrangebyscore"
	args := make([]string, 2, 5)

	if this.reversed {
		op = "zrevrangebyscore"
		args[0] = this.max
		args[1] = this.min
	} else {
		args[0] = this.min
		args[1] = this.max
	}

	if this.limited {
		args = append(args, "LIMIT", itoa(this.offset), itoa(this.count))
	}

	return intsChannel(SliceCommand(this.key, this.key.args(op, args...)...))
}

//GetWithScores returns a map with all members fitting the search criteria and their associated scores
func (this *SortedIntSetRange) GetWithScores() <-chan map[int]float64 {
	op := "zrangebyscore"
	args := make([]string, 3, 6)

	if this.reversed {
		op = "zrevrangebyscore"
		args[0] = this.max
		args[1] = this.min
	} else {
		args[0] = this.min
		args[1] = this.max
	}

	args[2] = "WITHSCORES"

	if this.limited {
		args = append(args, "LIMIT", itoa(this.offset), itoa(this.count))
	}

	return intfloatMapChannel(MapCommand(this.key, this.key.args(op, args...)...))
}

//SortedIntSetCombo keeps track of how you want to be combining multiple zsets
type SortedIntSetCombo struct {
	weighted bool
	op       string //either Union or Intersection
	sets     map[string]float64

	key Key
}

//StoreUnion sets up a combo that will be a union of other zsets
func (this SortedIntSet) StoreUnion() *SortedIntSetCombo {
	return &SortedIntSetCombo{
		op:  "zunionstore",
		key: this.Key,
	}
}

//StoreIntersection sets up a combo that will be an intersection of other zsets
func (this SortedIntSet) StoreIntersection() *SortedIntSetCombo {
	return &SortedIntSetCombo{
		op:  "zinterstore",
		key: this.Key,
	}
}

//OfSet adds a zset to the combo
func (this *SortedIntSetCombo) OfSet(otherSet SortedIntSet) *SortedIntSetCombo {
	if this.sets == nil {
		this.sets = make(map[string]float64)
	}
	this.sets[otherSet.key] = 1.0
	return this
}

//OfWeightedSet adds a zset to the combo, and weights it to be either heavier or lighter than other zsets
func (this *SortedIntSetCombo) OfWeightedSet(otherSet SortedIntSet, weight float64) *SortedIntSetCombo {
	if this.sets == nil {
		this.sets = make(map[string]float64)
	}
	this.weighted = true
	this.sets[otherSet.key] = weight
	return this
}

//UseLowerScore combines the zsets, and when duplicates are found, will keep the lowest score found
func (this *SortedIntSetCombo) UseLowerScore() <-chan int {
	return IntCommand(this.key, this.args("MIN")...)
}

//UseHigherScore combines the zsets, and when duplicates are found, will keep the highest score found
func (this *SortedIntSetCombo) UseHigherScore() <-chan int {
	return IntCommand(this.key, this.args("MAX")...)
}

//UseCombinedScores combines the zsets, and when duplicates are found, will add the scores together
func (this *SortedIntSetCombo) UseCombinedScores() <-chan int {
	return IntCommand(this.key, this.args("SUM")...)
}

func (this *SortedIntSetCombo) args(mode string) []string {
	result := make([]string, 1, 11)
	result[0] = itoa(len(this.sets))

	weights := make([]string, 1, 3)
	weights[0] = "WEIGHTS"

	for set, weight := range this.sets {
		result = append(result, set)
		weights = append(weights, ftoa(weight))
	}

	if this.weighted {
		result = append(result, weights...)
	}

	if mode != "SUM" {
		result = append(result, "AGGREGATE", mode)
	}

	return this.key.args(this.op, result...)
}

//Use allows you to use this key on a different executor
func (this SortedIntSet) Use(e SafeExecutor) SortedIntSet {
	this.client = e
	return this
}
