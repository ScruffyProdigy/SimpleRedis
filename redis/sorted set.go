package redis

type SortedSet struct {
	SortableKey
}

func newSortedSet(client SafeExecutor, key string) SortedSet {
	return SortedSet{
		newSortableKey(client, key),
	}
}

func (this SortedSet) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "zset")
	}()
	return c
}

func (this SortedSet) Add(item string, score float64) <-chan bool {
	return BoolCommand(this, this.args("zadd", ftoa(score), item))
}
func (this SortedSet) IncrementBy(item string, score float64) <-chan float64 {
	return FloatCommand(this, this.args("zincrby", ftoa(score), item))
}
func (this SortedSet) Remove(item string) <-chan bool {
	return BoolCommand(this, this.args("zrem", item))
}
func (this SortedSet) Size() <-chan int {
	return IntCommand(this, this.args("zcard"))
}
func (this SortedSet) IndexOf(item string) <-chan int {
	return IntCommand(this, this.args("zrank", item))
}
func (this SortedSet) ReverseIndexOf(item string) <-chan int {
	return IntCommand(this, this.args("zrevrank", item))
}
func (this SortedSet) ScoreOf(item string) <-chan float64 {
	return FloatCommand(this, this.args("zscore", item))
}
func (this SortedSet) IndexedBetween(start, stop int) <-chan []string {
	return SliceCommand(this, this.args("zrange", itoa(start), itoa(stop)))
}
func (this SortedSet) ReverseIndexedBetween(start, stop int) <-chan []string {
	return SliceCommand(this, this.args("zrevrange", itoa(start), itoa(stop)))
}
func (this SortedSet) RemoveIndexedBetween(start, stop int) <-chan int {
	return IntCommand(this, this.args("zremrangebyrank", itoa(start), itoa(stop)))
}

type SortedSetRange struct {
	min, max      string
	fmin, fmax    float64
	limited       bool
	offset, count int
	reversed      bool

	key Key
}

func (this SortedSet) Scores() *SortedSetRange {
	return &SortedSetRange{
		min: "-inf",
		max: "+inf",
		key: this.Key,
	}
}

func (this *SortedSetRange) Above(min float64) *SortedSetRange {
	if this.min == "-inf" || this.fmin <= min {
		this.fmin = min
		this.min = "(" + ftoa(min)
	}
	return this
}

func (this *SortedSetRange) Below(max float64) *SortedSetRange {
	if this.max == "+inf" || this.fmax >= max {
		this.fmax = max
		this.max = "(" + ftoa(max)
	}
	return this
}

func (this *SortedSetRange) AboveOrEqualTo(min float64) *SortedSetRange {
	if this.min == "-inf" || this.fmin < min {
		this.fmin = min
		this.min = ftoa(min)
	}
	return this
}

func (this *SortedSetRange) BelowOrEqualTo(max float64) *SortedSetRange {
	if this.max == "+inf" || this.fmax > max {
		this.fmax = max
		this.max = ftoa(max)
	}
	return this
}

//only useful if getting or getting with scores; not useful for counting or removing
func (this *SortedSetRange) Reversed() *SortedSetRange {
	this.reversed = !this.reversed
	return this
}

//only useful if getting or getting with scores; not useful for counting or removing
func (this *SortedSetRange) Limit(offset, count int) *SortedSetRange {
	this.limited = true
	this.offset = offset
	this.count = count
	return this
}

func (this *SortedSetRange) Count() <-chan int {
	return IntCommand(this.key, this.key.args("zcount", this.min, this.max))
}

func (this *SortedSetRange) Remove() <-chan int {
	return IntCommand(this.key, this.key.args("zremrangebyscore", this.min, this.max))
}

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

	return SliceCommand(this.key, this.key.args(op, args...))
}

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

	output := MapCommand(this.key, this.key.args(op, args...))
	realoutput := make(chan map[string]float64, 1)
	go func() {
		defer close(realoutput)
		if strings, ok := <-output; ok {
			result := make(map[string]float64, len(strings))
			for k, v := range strings {
				var err error
				result[k], err = atof(v)
				if err != nil {
					this.key.client.errCallback(err, "sorting with scores")
				}
			}
			realoutput <- result
		}
	}()
	return realoutput
}

type SortedSetCombo struct {
	weighted bool
	op       string //either Union or Intersection
	sets     map[string]float64

	key Key
}

func (this SortedSet) StoreUnion() *SortedSetCombo {
	return &SortedSetCombo{
		op:  "zunionstore",
		key: this.Key,
	}
}

func (this SortedSet) StoreIntersection() *SortedSetCombo {
	return &SortedSetCombo{
		op:  "zinterstore",
		key: this.Key,
	}
}

func (this *SortedSetCombo) OfSet(otherSet SortedSet) *SortedSetCombo {
	if this.sets == nil {
		this.sets = make(map[string]float64)
	}
	this.sets[otherSet.key] = 1.0
	return this
}

func (this *SortedSetCombo) OfWeightedSet(otherSet SortedSet, weight float64) *SortedSetCombo {
	if this.sets == nil {
		this.sets = make(map[string]float64)
	}
	this.weighted = true
	this.sets[otherSet.key] = weight
	return this
}

func (this *SortedSetCombo) UseLowerScore() <-chan int {
	return IntCommand(this.key, this.args("MIN"))
}

func (this *SortedSetCombo) UseHigherScore() <-chan int {
	return IntCommand(this.key, this.args("MAX"))
}

func (this *SortedSetCombo) UseCombinedScores() <-chan int {
	return IntCommand(this.key, this.args("SUM"))
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

func (this SortedSet) Use(e SafeExecutor) SortedSet {
	this.client = e
	return this
}
