package redis

type SortedSet struct {
	SortableKey
}

func newSortedSet(client SafeExecutor, key string) SortedSet {
	return SortedSet{
		newSortableKey(client, key),
	}
}

//IsValid returns whether the underlying redis object can use the commands in this object
func (this SortedSet) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "zset")
	}()
	return c
}

//Add adds a member to a zset or updates its score if it already exists - ZADD command
//returns true when adding, false when updating
func (this SortedSet) Add(item string, score float64) <-chan bool {
	return BoolCommand(this, this.args("zadd", ftoa(score), item)...)
}

//IncrementBy adjusts the score of the member within the zset - ZINCRBY command
//returns the new score
func (this SortedSet) IncrementBy(item string, score float64) <-chan float64 {
	return FloatCommand(this, this.args("zincrby", ftoa(score), item)...)
}

//Remove removes a member from the zset if it is part of the set - ZREM command
//returns whether or not it was part of the set
func (this SortedSet) Remove(item string) <-chan bool {
	return BoolCommand(this, this.args("zrem", item)...)
}

//Size returns the number of members of the zset - ZCARD command
func (this SortedSet) Size() <-chan int {
	return IntCommand(this, this.args("zcard")...)
}

//IndexOf returns the index of a member - ZRANK command
//ie, the lowest ranked member would have an index of 0, and the next lowest an index of 1
func (this SortedSet) IndexOf(item string) <-chan int {
	return IntCommand(this, this.args("zrank", item)...)
}

//ReverseIndexOf returns the reverse index of a member - ZREVRANK command
//ie, the highest ranked member would have an reverse index of 0, and the next highest an reverse index of 1
func (this SortedSet) ReverseIndexOf(item string) <-chan int {
	return IntCommand(this, this.args("zrevrank", item)...)
}

//ScoreOf returns the score associated with a given member of the zset - ZSCORE command
func (this SortedSet) ScoreOf(item string) <-chan float64 {
	return FloatCommand(this, this.args("zscore", item)...)
}

//IndexedBetween returns a slice of all members between the indices - ZRANGE command
func (this SortedSet) IndexedBetween(start, stop int) <-chan []string {
	return SliceCommand(this, this.args("zrange", itoa(start), itoa(stop))...)
}

//ReverseIndexedBetween returns a slice of all members between the reverse indices - ZREVRANGE command
func (this SortedSet) ReverseIndexedBetween(start, stop int) <-chan []string {
	return SliceCommand(this, this.args("zrevrange", itoa(start), itoa(stop))...)
}

//IndexedBetweenWithScores returns a map of all members between the indices and their associated scores - ZRANGE command
//warning: golang maps are not ordered
func (this SortedSet) IndexedBetweenWithScores(start, stop int) <-chan map[string]float64 {
	return stringfloatMapChannel(MapCommand(this, this.args("zrange", itoa(start), itoa(stop), "WITHSCORES")...))
}

//IndexedBetweenWithScores returns a map of all members between the reverse indices and their associated scores - ZREVRANGE command
//warning: golang maps are not ordered
func (this SortedSet) ReverseIndexedBetweenWithScores(start, stop int) <-chan map[string]float64 {
	return stringfloatMapChannel(MapCommand(this, this.args("zrevrange", itoa(start), itoa(stop), "WITHSCORES")...))
}

//RemoveIndexedBetween removes all members between the indices - ZREMRANGEBYRANK command
//returns the number of members removed
func (this SortedSet) RemoveIndexedBetween(start, stop int) <-chan int {
	return IntCommand(this, this.args("zremrangebyrank", itoa(start), itoa(stop))...)
}

//SortedSetRange keeps track of all range arguments being used in a search
type SortedSetRange struct {
	min, max      string
	fmin, fmax    float64
	limited       bool
	offset, count int
	reversed      bool

	key Key
}

//Scores createa a SortedSetRange to help narrow a search to be done later
func (this SortedSet) Scores() *SortedSetRange {
	return &SortedSetRange{
		min: "-inf",
		max: "+inf",
		key: this.Key,
	}
}

//Above limits results to members who have a score above "min"
func (this *SortedSetRange) Above(min float64) *SortedSetRange {
	if this.min == "-inf" || this.fmin <= min {
		this.fmin = min
		this.min = "(" + ftoa(min)
	}
	return this
}

//Below limits results to members who have a score below "max"
func (this *SortedSetRange) Below(max float64) *SortedSetRange {
	if this.max == "+inf" || this.fmax >= max {
		this.fmax = max
		this.max = "(" + ftoa(max)
	}
	return this
}

//AboveOrEqualTo limits results to members who have a score above or equal to "min"
func (this *SortedSetRange) AboveOrEqualTo(min float64) *SortedSetRange {
	if this.min == "-inf" || this.fmin < min {
		this.fmin = min
		this.min = ftoa(min)
	}
	return this
}

//BelowOrEqualTo limits results to members who have a score below or equal to "max"
func (this *SortedSetRange) BelowOrEqualTo(max float64) *SortedSetRange {
	if this.max == "+inf" || this.fmax > max {
		this.fmax = max
		this.max = ftoa(max)
	}
	return this
}

//Reversed returns the results in reverse order
//only useful if getting or getting with scores; not useful for counting or removing
func (this *SortedSetRange) Reversed() *SortedSetRange {
	this.reversed = !this.reversed
	return this
}

//Limit limits the results you get back - it skips the first "offset" results, and then only returns the next "offset"
//only useful if getting or getting with scores; not useful for counting or removing
func (this *SortedSetRange) Limit(offset, count int) *SortedSetRange {
	this.limited = true
	this.offset = offset
	this.count = count
	return this
}

//Count returns the number of members that fit in the search criteria - ZCOUNT command
func (this *SortedSetRange) Count() <-chan int {
	return IntCommand(this.key, this.key.args("zcount", this.min, this.max)...)
}

//Remove removes all members that fit the search criteria from the zset - ZREMRANGEBYSCORE command
//returns the number of members removed
func (this *SortedSetRange) Remove() <-chan int {
	return IntCommand(this.key, this.key.args("zremrangebyscore", this.min, this.max)...)
}

//Get returns a list of all members fitting the search criteria - ZRANGEBYSCORE or ZREVRANGEBYSCORE command
func (this *SortedSetRange) Get() <-chan []string {
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

	return SliceCommand(this.key, this.key.args(op, args...)...)
}

//GetWithScores returns a map with all members fitting the search criteria and their associated scores - ZRANGEBYSCORE or ZREVRANGEBYSCORE command
func (this *SortedSetRange) GetWithScores() <-chan map[string]float64 {
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

	return stringfloatMapChannel(MapCommand(this.key, this.key.args(op, args...)...))
}

//SortedSetCombo keeps track of how you want to be combining multiple zsets
type SortedSetCombo struct {
	weighted bool
	op       string //either Union or Intersection
	sets     map[string]float64

	key Key
}

//StoreUnion sets up a combo that will be a union of other zsets - ZUNIONSTORE command
func (this SortedSet) StoreUnion() *SortedSetCombo {
	return &SortedSetCombo{
		op:  "zunionstore",
		key: this.Key,
	}
}

//StoreIntersection sets up a combo that will be an intersection of other zsets - ZINTERSTORE command
func (this SortedSet) StoreIntersection() *SortedSetCombo {
	return &SortedSetCombo{
		op:  "zinterstore",
		key: this.Key,
	}
}

//OfSet adds a zset to the combo
func (this *SortedSetCombo) OfSet(otherSet SortedSet) *SortedSetCombo {
	if this.sets == nil {
		this.sets = make(map[string]float64)
	}
	this.sets[otherSet.key] = 1.0
	return this
}

//OfWeightedSet adds a zset to the combo, and weights it to be either heavier or lighter than other zsets
func (this *SortedSetCombo) OfWeightedSet(otherSet SortedSet, weight float64) *SortedSetCombo {
	if this.sets == nil {
		this.sets = make(map[string]float64)
	}
	this.weighted = true
	this.sets[otherSet.key] = weight
	return this
}

//UseLowerScore combines the zsets, and when duplicates are found, will keep the lowest score found
func (this *SortedSetCombo) UseLowerScore() <-chan int {
	return IntCommand(this.key, this.args("MIN")...)
}

//UseHigherScore combines the zsets, and when duplicates are found, will keep the highest score found
func (this *SortedSetCombo) UseHigherScore() <-chan int {
	return IntCommand(this.key, this.args("MAX")...)
}

//UseCombinedScores combines the zsets, and when duplicates are found, will add the scores together
func (this *SortedSetCombo) UseCombinedScores() <-chan int {
	return IntCommand(this.key, this.args("SUM")...)
}

func (this *SortedSetCombo) args(mode string) []string {
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
func (this SortedSet) Use(e SafeExecutor) SortedSet {
	this.client = e
	return this
}
