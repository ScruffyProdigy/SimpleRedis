package redis

type SortedFloatSet struct {
	SortableKey
}

func newSortedFloatSet(client Executor, key string) SortedFloatSet {
	return SortedFloatSet{
		newSortableKey(client, key),
	}
}

func (this SortedFloatSet) IsValid() <-chan bool {
	c := make(chan bool, 1)
	func() {
		defer close(c)
		c <- (<-this.Type() == "zset")
	}()
	return c
}

func (this SortedFloatSet) Add(item float64, score float64) <-chan bool {
	command, output := newBoolCommand(this.args("zadd", ftoa(score), ftoa(item)))
	this.Execute(command)
	return output
}

func (this SortedFloatSet) IncrementBy(item float64, score float64) <-chan float64 {
	command, output := newFloatCommand(this.args("zincrby", ftoa(score), ftoa(item)))
	this.Execute(command)
	return output
}

func (this SortedFloatSet) Remove(item float64) <-chan bool {
	command, output := newBoolCommand(this.args("zrem", ftoa(item)))
	this.Execute(command)
	return output
}

func (this SortedFloatSet) Size() <-chan int {
	command, output := newIntCommand(this.args("zcard"))
	this.Execute(command)
	return output
}

func (this SortedFloatSet) IndexOf(item float64) <-chan int {
	command, output := newIntCommand(this.args("zrank", ftoa(item)))
	this.Execute(command)
	return output
}

func (this SortedFloatSet) ReverseIndexOf(item float64) <-chan int {
	command, output := newIntCommand(this.args("zrevrank", ftoa(item)))
	this.Execute(command)
	return output
}

func (this SortedFloatSet) ScoreOf(item float64) <-chan float64 {
	command, output := newFloatCommand(this.args("zscore", ftoa(item)))
	this.Execute(command)
	return output
}

func (this SortedFloatSet) IndexedBetween(start, stop int) <-chan []float64 {
	command, output := newSliceCommand(this.args("zrange", itoa(start), itoa(stop)))
	this.Execute(command)
	realoutput := make(chan []float64, 1)
	go func() {
		defer close(realoutput)
		if midway, ok := <-output; ok {
			floats, err := stringsToFloats(midway)
			if err != nil {
				this.client.ErrCallback(err, "zrange")
			}
			realoutput <- floats
		}
	}()
	return realoutput
}

func (this SortedFloatSet) ReverseIndexedBetween(start, stop int) <-chan []float64 {
	command, output := newSliceCommand(this.args("zrevrange", itoa(start), itoa(stop)))
	this.Execute(command)
	realoutput := make(chan []float64, 1)
	go func() {
		defer close(realoutput)
		if midway, ok := <-output; ok {
			floats, err := stringsToFloats(midway)
			if err != nil {
				this.client.ErrCallback(err, "zrevrange")
				return
			}
			realoutput <- floats
		}
	}()
	return realoutput
}

func (this SortedFloatSet) RemoveIndexedBetween(start, stop int) <-chan int {
	command, output := newIntCommand(this.args("zremrangebyrank", itoa(start), itoa(stop)))
	this.Execute(command)
	return output
}

type SortedFloatSetRange struct {
	min, max      string
	fmin, fmax    float64
	limited       bool
	offset, count int
	reversed      bool

	key Key
}

func (this SortedFloatSet) Scores() *SortedFloatSetRange {
	return &SortedFloatSetRange{
		min: "-inf",
		max: "+inf",
		key: this.Key,
	}
}

func (this *SortedFloatSetRange) Above(min float64) *SortedFloatSetRange {
	if this.min == "-inf" || this.fmin <= min {
		this.fmin = min
		this.min = "(" + ftoa(min)
	}
	return this
}

func (this *SortedFloatSetRange) Below(max float64) *SortedFloatSetRange {
	if this.max == "+inf" || this.fmax >= max {
		this.fmax = max
		this.max = "(" + ftoa(max)
	}
	return this
}

func (this *SortedFloatSetRange) AboveOrEqualTo(min float64) *SortedFloatSetRange {
	if this.min == "-inf" || this.fmin < min {
		this.fmin = min
		this.min = ftoa(min)
	}
	return this
}

func (this *SortedFloatSetRange) BelowOrEqualTo(max float64) *SortedFloatSetRange {
	if this.max == "+inf" || this.fmax > max {
		this.fmax = max
		this.max = ftoa(max)
	}
	return this
}

//only useful if getting or getting with scores; not useful for counting or removing
func (this *SortedFloatSetRange) Reversed() *SortedFloatSetRange {
	this.reversed = !this.reversed
	return this
}

//only useful if getting or getting with scores; not useful for counting or removing
func (this *SortedFloatSetRange) Limit(offset, count int) *SortedFloatSetRange {
	this.limited = true
	this.offset = offset
	this.count = count
	return this
}

func (this *SortedFloatSetRange) Count() <-chan int {
	command, output := newIntCommand(this.key.args("zcount", this.min, this.max))
	this.key.Execute(command)
	return output
}

func (this *SortedFloatSetRange) Remove() <-chan int {
	command, output := newIntCommand(this.key.args("zremrangebyscore", this.min, this.max))
	this.key.Execute(command)
	return output
}

func (this *SortedFloatSetRange) Get() <-chan []float64 {
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

	realoutput := make(chan []float64, 1)
	go func() {
		defer close(realoutput)
		if midway, ok := <-output; ok {
			floats, err := stringsToFloats(midway)
			if err != nil {
				this.key.client.ErrCallback(err, "sort by score")
				return
			}
			realoutput <- floats
		}
	}()

	return realoutput
}

func (this *SortedFloatSetRange) GetWithScores() <-chan map[float64]float64 {
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

	realoutput := make(chan map[float64]float64, 1)
	go func() {
		defer close(realoutput)
		if midway, ok := <-output; ok {
			result := make(map[float64]float64, len(midway))
			for k, v := range midway {
				fk, err := atof(k)
				if err != nil {
					this.key.client.ErrCallback(err, "sorting with scores (key)")
					return
				}

				fv, err := atof(v)
				if err != nil {
					this.key.client.ErrCallback(err, "sorting with scores (value)")
					return
				}

				result[fk] = fv
			}
			realoutput <- result
		}
	}()

	return realoutput
}

type SortedFloatSetCombo struct {
	weighted bool
	op       string //either Union or Intersection
	mode     string //either Min, Max, or Sum
	sets     map[string]float64

	key Key
}

func (this SortedFloatSet) StoreUnion() *SortedFloatSetCombo {
	return &SortedFloatSetCombo{
		op:  "zunionstore",
		key: this.Key,
	}
}

func (this SortedFloatSet) StoreIntersection() *SortedFloatSetCombo {
	return &SortedFloatSetCombo{
		op:  "zinterstore",
		key: this.Key,
	}
}

func (this *SortedFloatSetCombo) OfSet(otherSet SortedFloatSet) *SortedFloatSetCombo {
	if this.sets == nil {
		this.sets = make(map[string]float64)
	}
	this.sets[otherSet.key] = 1.0
	return this
}

func (this *SortedFloatSetCombo) OfWeightedSet(otherSet SortedFloatSet, weight float64) *SortedFloatSetCombo {
	if this.sets == nil {
		this.sets = make(map[string]float64)
	}
	this.weighted = true
	this.sets[otherSet.key] = weight
	return this
}

func (this *SortedFloatSetCombo) UseLowerScore() <-chan int {
	command, output := newIntCommand(this.args("MIN"))
	this.key.Execute(command)
	return output
}

func (this *SortedFloatSetCombo) UseHigherScore() <-chan int {
	command, output := newIntCommand(this.args("MAX"))
	this.key.Execute(command)
	return output
}

func (this *SortedFloatSetCombo) UseCombinedScores() <-chan int {
	command, output := newIntCommand(this.args("SUM"))
	this.key.Execute(command)
	return output
}

func (this *SortedFloatSetCombo) args(mode string) []string {
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

func (this SortedFloatSet) Use(e Executor) SortedFloatSet {
	this.client = e
	return this
}
