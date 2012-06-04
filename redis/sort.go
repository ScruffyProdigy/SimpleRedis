package redis

type sortLimit struct {
	min, max int
}

type sortBy struct {
	pattern string
}

type sortGet struct {
	pattern string
}

type Sorter struct {
	limit    *sortLimit
	by       *sortBy
	get      []sortGet
	typeSet  bool
	alpha    bool
	reversed bool

	key Key
}

func (this Sorter) sortargs() []string {
	if !this.typeSet {
		panic("Sort Alphabetically or Numerically?!")
	}
	result := make([]string, 0, 10)
	if this.by != nil {
		result = append(result, "BY", this.by.pattern)
	}
	if this.limit != nil {
		result = append(result, "LIMIT", itoa(this.limit.min), itoa(this.limit.max))
	}
	for _, get := range this.get {
		result = append(result, "GET", get.pattern)
	}
	if this.reversed {
		result = append(result, "DESC")
	}
	if this.alpha {
		result = append(result, "ALPHA")
	}
	return result
}

func (this Sorter) sortstoreargs(dest string) []string {
	return append(this.sortargs(), "STORE", dest)
}

func (this Sorter) Limit(min, max int) Sorter {
	if this.limit != nil {
		panic("Sort 'Limit' already set!")
	}
	this.limit = &sortLimit{
		min: min,
		max: max,
	}
	return this
}

func (this Sorter) Alphabetically() Sorter {
	if this.typeSet {
		str := "Sort already sorting numerically"
		if this.alpha {
			str = "Sort already sorting alphabetically"
		}
		panic(str)
	}
	this.alpha = true
	this.typeSet = true
	return this
}

func (this Sorter) Numerically() Sorter {
	if this.typeSet {
		str := "Sort already sorting numerically"
		if this.alpha {
			str = "Sort already sorting alphabetically"
		}
		panic(str)
	}
	this.typeSet = true
	return this
}

//if we could figure out what kind of pattern was referred to here, we could automatically set alphabetically or numerically
func (this Sorter) By(pattern string) Sorter {
	if this.by != nil {
		panic("Sort 'By' already set")
	}
	this.by = &sortBy{
		pattern: pattern,
	}
	return this
}

//if we could figure out what kind of pattern was referred to here, we could know what kind of channels to return later
func (this Sorter) Get(pattern string) Sorter {
	this.get = append(this.get, sortGet{
		pattern: pattern,
	})
	return this
}

func (this Sorter) Reverse() Sorter {
	this.reversed = !this.reversed
	return this
}

func (this Sorter) Strings() <-chan []string {
	command, output := newSliceCommand(this.key.args("sort", this.sortargs()...))
	this.key.Execute(command)
	return output
}

func (this Sorter) Ints() <-chan []int {
	realoutput := make(chan []int, 1)
	midway := this.Strings()
	go func() {
		if output, ok := <-midway; ok {
			realoutput <- stringsToInts(output)
		}
		close(realoutput)
	}()
	return realoutput
}

func (this Sorter) Floats() <-chan []float64 {
	realoutput := make(chan []float64, 1)
	midway := this.Strings()
	go func() {
		if output, ok := <-midway; ok {
			realoutput <- stringsToFloats(output)
		}
		close(realoutput)
	}()
	return realoutput
}

//MaybeStrings, MaybeInts, and MaybeFloats should be used when you need to distinguish between 0 and nil responses
func (this Sorter) MaybeStrings() <-chan []*string {
	command, output := newMaybeSliceCommand(this.key.args("sort", this.sortargs()...))
	this.key.Execute(command)
	return output
}

func (this Sorter) MaybeInts() <-chan []*int {
	realoutput := make(chan []*int, 1)
	midway := this.MaybeStrings()
	go func() {
		if strings, ok := <-midway; ok {
			ints := make([]*int, len(strings))
			for i, str := range strings {
				if str != nil {
					j := atoi(*str)
					ints[i] = &j
				}
			}

			realoutput <- ints
		}
		close(realoutput)
	}()
	return realoutput
}

func (this Sorter) MaybeFloats() <-chan []*float64 {
	realoutput := make(chan []*float64, 1)
	midway := this.MaybeStrings()
	go func() {
		if strings, ok := <-midway; ok {
			floats := make([]*float64, len(strings))
			for i, str := range strings {
				if str != nil {
					j := atof(*str)
					floats[i] = &j
				}
			}

			realoutput <- floats
		}
		close(realoutput)
	}()
	return realoutput
}

func (this Sorter) StoreStrings(dest List) <-chan int {
	command, output := newIntCommand(this.key.args("sort", this.sortstoreargs(dest.key)...))
	this.key.Execute(command)
	return output
}

func (this Sorter) StoreInts(dest IntList) <-chan int {
	command, output := newIntCommand(this.key.args("sort", this.sortstoreargs(dest.key)...))
	this.key.Execute(command)
	return output
}

func (this Sorter) StoreFloats(dest FloatList) <-chan int {
	command, output := newIntCommand(this.key.args("sort", this.sortstoreargs(dest.key)...))
	this.key.Execute(command)
	return output
}
