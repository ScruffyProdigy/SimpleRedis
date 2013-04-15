package redis

type SortedIntSet struct {
	Key
}

func newSortedIntSet(client Executor, key string) SortedIntSet {
	return SortedIntSet{
		newKey(client, key),
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
	command, output := newBoolCommand(this.args("zadd", ftoa(score), itoa(item)))
	this.Execute(command)
	return output
}

func (this SortedIntSet) IncrementBy(item int, score float64) <-chan float64 {
	command, output := newFloatCommand(this.args("zincrby", ftoa(score), itoa(item)))
	this.Execute(command)
	return output
}

func (this SortedIntSet) Remove(item int) <-chan bool {
	command, output := newBoolCommand(this.args("zrem", itoa(item)))
	this.Execute(command)
	return output
}

func (this SortedIntSet) Size() <-chan int {
	command, output := newIntCommand(this.args("zcard"))
	this.Execute(command)
	return output
}

func (this SortedIntSet) IndexOf(item int) <-chan int {
	command, output := newIntCommand(this.args("zrank", itoa(item)))
	this.Execute(command)
	return output
}

func (this SortedIntSet) ReverseIndexOf(item int) <-chan int {
	command, output := newIntCommand(this.args("zrevrank", itoa(item)))
	this.Execute(command)
	return output
}

func (this SortedIntSet) ScoreOf(item int) <-chan float64 {
	command, output := newFloatCommand(this.args("zscore", itoa(item)))
	this.Execute(command)
	return output
}

func (this SortedIntSet) IndexedBetween(start, stop int) <-chan []int {
	command, output := newSliceCommand(this.args("zrange", itoa(start), itoa(stop)))
	this.Execute(command)
	realoutput := make(chan []int, 1)
	go func() {
		defer close(realoutput)
		if midway, ok := <-output; ok {
			ints, err := stringsToInts(midway)
			if err != nil {
				this.client.ErrCallback(err, "zrange")
			}
			realoutput <- ints
		}
	}()
	return realoutput
}

func (this SortedIntSet) ReverseIndexedBetween(start, stop int) <-chan []int {
	command, output := newSliceCommand(this.args("zrevrange", itoa(start), itoa(stop)))
	this.Execute(command)
	realoutput := make(chan []int, 1)
	go func() {
		defer close(realoutput)
		if midway, ok := <-output; ok {
			ints, err := stringsToInts(midway)
			if err != nil {
				this.client.ErrCallback(err, "zrevrange")
			}
			realoutput <- ints
		}
	}()
	return realoutput
}

func (this SortedIntSet) RemoveIndexedBetween(start, stop int) <-chan int {
	command, output := newIntCommand(this.args("zremrangebyrank", itoa(start), itoa(stop)))
	this.Execute(command)
	return output
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
	command, output := newIntCommand(this.key.args("zcount", this.min, this.max))
	this.key.Execute(command)
	return output
}

func (this *SortedIntSetRange) Remove() <-chan int {
	command, output := newIntCommand(this.key.args("zremrangebyscore", this.min, this.max))
	this.key.Execute(command)
	return output
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
	
	command, output := newSliceCommand(this.key.args(op, args...))
	this.key.Execute(command)
	
	realoutput := make(chan []int, 1)
	go func() {
		defer close(realoutput)
		if midway, ok := <-output; ok {
			ints, err := stringsToInts(midway)
			if err != nil {
				this.key.client.ErrCallback(err, "sorting ints")
			}
			realoutput <- ints
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
	
	command, output := newMapCommand(this.key.args(op, args...))
	this.key.Execute(command)

	realoutput := make(chan map[int]float64, 1)
	go func() {
		defer close(realoutput)
		if midway, ok := <-output; ok {
			result := make(map[int]float64, len(midway))
			for k, v := range midway {
				index, err := atoi(k)
				if err != nil {
					this.key.client.ErrCallback(err, "sorting with scores (key)")
				}

				result[index], err = atof(v)
				if err != nil {
					this.key.client.ErrCallback(err, "sorting with scores (value)")
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
	command, output := newIntCommand(this.args("MIN"))
	this.key.Execute(command)
	return output
}

func (this *SortedIntSetCombo) UseHigherScore() <-chan int {
	command, output := newIntCommand(this.args("MAX"))
	this.key.Execute(command)
	return output
}

func (this *SortedIntSetCombo) UseCombinedScores() <-chan int {
	command, output := newIntCommand(this.args("SUM"))
	this.key.Execute(command)
	return output
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

func (this SortedIntSet) Use(e Executor) SortedIntSet {
	this.client = e
	return this
}
