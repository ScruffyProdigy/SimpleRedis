package redis

type SortedIntSet struct {
	SortableKey
}

func newSortedIntSet(client SafeExecutor, key string) SortedIntSet {
	return SortedIntSet{
		newSortableKey(client, key),
	}
}

func (this SortedIntSet) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "zset")
	}()
	return c
}

func (this SortedIntSet) Add(item int, score float64) <-chan bool {
	return BoolCommand(this, this.args("zadd", ftoa(score), itoa(item)))
}
func (this SortedIntSet) IncrementBy(item int, score float64) <-chan float64 {
	return FloatCommand(this, this.args("zincrby", ftoa(score), itoa(item)))
}
func (this SortedIntSet) Remove(item int) <-chan bool {
	return BoolCommand(this, this.args("zrem", itoa(item)))
}
func (this SortedIntSet) Size() <-chan int {
	return IntCommand(this, this.args("zcard"))
}
func (this SortedIntSet) IndexOf(item int) <-chan int {
	return IntCommand(this, this.args("zrank", itoa(item)))
}
func (this SortedIntSet) ReverseIndexOf(item int) <-chan int {
	return IntCommand(this, this.args("zrevrank", itoa(item)))
}
func (this SortedIntSet) ScoreOf(item int) <-chan float64 {
	return FloatCommand(this, this.args("zscore", itoa(item)))
}
func (this SortedIntSet) IndexedBetween(start, stop int) <-chan []int {
	output := SliceCommand(this, this.args("zrange", itoa(start), itoa(stop)))
	realoutput := make(chan []int, 1)
	go func() {
		defer close(realoutput)
		if strings, ok := <-output; ok {
			if ints, err := stringsToInts(strings); err != nil {
				this.client.errCallback(err, "zrange")
			} else {
				realoutput <- ints
			}
		}
	}()
	return realoutput
}

func (this SortedIntSet) ReverseIndexedBetween(start, stop int) <-chan []int {
	output := SliceCommand(this, this.args("zrevrange", itoa(start), itoa(stop)))
	realoutput := make(chan []int, 1)
	go func() {
		defer close(realoutput)
		if strings, ok := <-output; ok {
			if ints, err := stringsToInts(strings); err != nil {
				this.client.errCallback(err, "zrevrange")
			} else {
				realoutput <- ints
			}
		}
	}()
	return realoutput
}

func (this SortedIntSet) RemoveIndexedBetween(start, stop int) <-chan int {
	return IntCommand(this, this.args("zremrangebyrank", itoa(start), itoa(stop)))
}

type SortedIntSetRange struct {
	min, max      string
	fmin, fmax    float64
	limited       bool
	offset, count int
	reversed      bool

	key Key
}

func (this SortedIntSet) Scores() *SortedIntSetRange {
	return &SortedIntSetRange{
		min: "-inf",
		max: "+inf",
		key: this.Key,
	}
}

func (this *SortedIntSetRange) Above(min float64) *SortedIntSetRange {
	if this.min == "-inf" || this.fmin <= min {
		this.fmin = min
		this.min = "(" + ftoa(min)
	}
	return this
}

func (this *SortedIntSetRange) Below(max float64) *SortedIntSetRange {
	if this.max == "+inf" || this.fmax >= max {
		this.fmax = max
		this.max = "(" + ftoa(max)
	}
	return this
}

func (this *SortedIntSetRange) AboveOrEqualTo(min float64) *SortedIntSetRange {
	if this.min == "-inf" || this.fmin < min {
		this.fmin = min
		this.min = ftoa(min)
	}
	return this
}

func (this *SortedIntSetRange) BelowOrEqualTo(max float64) *SortedIntSetRange {
	if this.max == "+inf" || this.fmax > max {
		this.fmax = max
		this.max = ftoa(max)
	}
	return this
}

//only useful if getting or getting with scores; not useful for counting or removing
func (this *SortedIntSetRange) Reversed() *SortedIntSetRange {
	this.reversed = !this.reversed
	return this
}

//only useful if getting or getting with scores; not useful for counting or removing
func (this *SortedIntSetRange) Limit(offset, count int) *SortedIntSetRange {
	this.limited = true
	this.offset = offset
	this.count = count
	return this
}

func (this *SortedIntSetRange) Count() <-chan int {
	return IntCommand(this.key, this.key.args("zcount", this.min, this.max))
}

func (this *SortedIntSetRange) Remove() <-chan int {
	return IntCommand(this.key, this.key.args("zremrangebyscore", this.min, this.max))
}

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

	output := SliceCommand(this.key, this.key.args(op, args...))
	realoutput := make(chan []int, 1)
	go func() {
		defer close(realoutput)
		if strings, ok := <-output; ok {
			if ints, err := stringsToInts(strings); err != nil {
				this.key.client.errCallback(err, "sorting ints")
			} else {
				realoutput <- ints
			}
		}
	}()

	return realoutput
}

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

	output := MapCommand(this.key, this.key.args(op, args...))
	realoutput := make(chan map[int]float64, 1)
	go func() {
		defer close(realoutput)
		if midway, ok := <-output; ok {
			result := make(map[int]float64, len(midway))
			for k, v := range midway {
				index, err := atoi(k)
				if err != nil {
					this.key.client.errCallback(err, "sorting with scores (key)")
				}

				result[index], err = atof(v)
				if err != nil {
					this.key.client.errCallback(err, "sorting with scores (value)")
				}
			}
			realoutput <- result
		}
	}()

	return realoutput
}

type SortedIntSetCombo struct {
	weighted bool
	op       string //either Union or Intersection
	sets     map[string]float64

	key Key
}

func (this SortedIntSet) StoreUnion() *SortedIntSetCombo {
	return &SortedIntSetCombo{
		op:  "zunionstore",
		key: this.Key,
	}
}

func (this SortedIntSet) StoreIntersection() *SortedIntSetCombo {
	return &SortedIntSetCombo{
		op:  "zinterstore",
		key: this.Key,
	}
}

func (this *SortedIntSetCombo) OfSet(otherSet SortedIntSet) *SortedIntSetCombo {
	if this.sets == nil {
		this.sets = make(map[string]float64)
	}
	this.sets[otherSet.key] = 1.0
	return this
}

func (this *SortedIntSetCombo) OfWeightedSet(otherSet SortedIntSet, weight float64) *SortedIntSetCombo {
	if this.sets == nil {
		this.sets = make(map[string]float64)
	}
	this.weighted = true
	this.sets[otherSet.key] = weight
	return this
}

func (this *SortedIntSetCombo) UseLowerScore() <-chan int {
	return IntCommand(this.key, this.args("MIN"))
}

func (this *SortedIntSetCombo) UseHigherScore() <-chan int {
	return IntCommand(this.key, this.args("MAX"))
}

func (this *SortedIntSetCombo) UseCombinedScores() <-chan int {
	return IntCommand(this.key, this.args("SUM"))
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

func (this SortedIntSet) Use(e SafeExecutor) SortedIntSet {
	this.client = e
	return this
}
