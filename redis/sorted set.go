package redis

type SortedSet struct {
	Key
}

func newSortedSet(client Executor, key string) SortedSet {
	return SortedSet{
		newKey(client, key),
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
	command, output := newBoolCommand(this.args("zadd", ftoa(score), item))
	this.Execute(command)
	return output
}

func (this SortedSet) IncrementBy(item string, score float64) <-chan float64 {
	command, output := newFloatCommand(this.args("zincrby", ftoa(score), item))
	this.Execute(command)
	return output
}

func (this SortedSet) Remove(item string) <-chan bool {
	command, output := newBoolCommand(this.args("zrem", item))
	this.Execute(command)
	return output
}

func (this SortedSet) Size() <-chan int {
	command, output := newIntCommand(this.args("zcard"))
	this.Execute(command)
	return output
}

func (this SortedSet) IndexOf(item string) <-chan int {
	command, output := newIntCommand(this.args("zrank", item))
	this.Execute(command)
	return output
}

func (this SortedSet) ReverseIndexOf(item string) <-chan int {
	command, output := newIntCommand(this.args("zrevrank", item))
	this.Execute(command)
	return output
}

func (this SortedSet) ScoreOf(item string) <-chan float64 {
	command, output := newFloatCommand(this.args("zscore", item))
	this.Execute(command)
	return output
}

func (this SortedSet) IndexedBetween(start, stop int) <-chan []string {
	command, output := newSliceCommand(this.args("zrange", itoa(start), itoa(stop)))
	this.Execute(command)
	return output
}

func (this SortedSet) ReverseIndexedBetween(start, stop int) <-chan []string {
	command, output := newSliceCommand(this.args("zrevrange", itoa(start), itoa(stop)))
	this.Execute(command)
	return output
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
	if this.min == "-inf" || this.fmin >= min {
		this.fmin = min
		this.min = ftoa(min)
	}
	return this
}

func (this *SortedSetRange) Below(max float64) *SortedSetRange {
	if this.max == "+inf" || this.fmax <= max {
		this.fmax = max
		this.max = ftoa(max)
	}
	return this
}

func (this *SortedSetRange) AboveOrEqualTo(min float64) *SortedSetRange {
	if this.min == "-inf" || this.fmin > min {
		this.fmin = min
		this.min = "(" + ftoa(min)
	}
	return this
}

func (this *SortedSetRange) BelowOrEqualTo(max float64) *SortedSetRange {
	if this.max == "+inf" || this.fmax < max {
		this.fmax = max
		this.max = "(" + ftoa(max)
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
	command, output := newIntCommand(this.key.args("zcount", this.min, this.max))
	this.key.Execute(command)
	return output
}

func (this *SortedSetRange) Remove() <-chan int {
	command, output := newIntCommand(this.key.args("zremrangebyrank", this.min, this.max))
	this.key.Execute(command)
	return output
}

func (this *SortedSetRange) Get() <-chan []string {
	op := "zrangebyscore"
	if this.reversed {
		op = "zrevrangebyscore"
	}
	args := make([]string, 2, 5)
	args[0] = this.min
	args[1] = this.max
	if this.limited {
		args = append(args, "LIMIT", itoa(this.offset), itoa(this.count))
	}
	command, output := newSliceCommand(this.key.args(op, args...))
	this.key.Execute(command)
	return output
}

func (this *SortedSetRange) GetWithScores() <-chan map[string]float64 {
	op := "zrangebyscore"
	if this.reversed {
		op = "zrevrangebyscore"
	}
	args := make([]string, 2, 6)
	args[0] = this.min
	args[1] = this.max
	if this.limited {
		args = append(args, "WITHSCORES", "LIMIT", itoa(this.offset), itoa(this.count))
	}
	command, output := newMapCommand(this.key.args(op, args...))
	this.key.Execute(command)
	realoutput := make(chan map[string]float64, 1)
	go func() {
		defer close(realoutput)
		if midway, ok := <-output; ok {
			result := make(map[string]float64, len(midway))
			for k, v := range midway {
				var err error
				result[k], err = atof(v)
				if err != nil {
					this.key.client.ErrCallback(err, "sorting with scores")
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
	mode     string //either Min, Max, or Sum
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
	command, output := newIntCommand(this.args("MIN"))
	this.key.Execute(command)
	return output
}

func (this *SortedSetCombo) UseHigherScore() <-chan int {
	command, output := newIntCommand(this.args("MAX"))
	this.key.Execute(command)
	return output
}

func (this *SortedSetCombo) UseCombinedScores() <-chan int {
	command, output := newIntCommand(this.args("SUM"))
	this.key.Execute(command)
	return output
}

func (this *SortedSetCombo) args(mode string) []string {
	result := make([]string, 0, 10)
	weights := make([]string, 0, 2)
	for set, weight := range this.sets {
		result = append(result, set)
		weights = append(weights, ftoa(weight))
	}
	if this.weighted {
		result = append(result, "WEIGHTS")
		result = append(result, weights...)
	}
	if mode != "SUM" {
		result = append(result, "AGGREGATE", mode)
	}
	return this.key.args(this.op, result...)
}

func (this SortedSet) Use(e Executor) SortedSet {
	this.client = e
	return this
}
