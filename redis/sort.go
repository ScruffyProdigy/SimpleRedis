package redis

type SortableKey struct {
	Key
}

func newSortableKey(client Executor, key string) SortableKey {
	return SortableKey{
		newKey(client, key),
	}
}

func (this SortableKey) SortAlphabetically() *Sorter {
	return &Sorter{key: this.Key, alpha: true}
}

func (this SortableKey) SortNumerically() *Sorter {
	return &Sorter{key: this.Key, alpha: false}
}

type sortLimit struct {
	min, max int
}

type sortBy struct {
	pattern string
}

type sortGet struct {
	pattern string
}

type sortStore struct {
	dest string
}

type Sorter struct {
	limit    *sortLimit
	by       *sortBy
	get      *sortGet
	store    *sortStore
	alpha    bool
	reversed bool

	key Key
}

func (this Sorter) sortargs() []string {
	result := make([]string, 0, 10)
	if this.by != nil {
		result = append(result, "BY", this.by.pattern)
	}
	if this.limit != nil {
		result = append(result, "LIMIT", itoa(this.limit.min), itoa(this.limit.max))
	}
	if this.get != nil {
		result = append(result, "GET", this.get.pattern)
	}
	if this.reversed {
		result = append(result, "DESC")
	}
	if this.alpha {
		result = append(result, "ALPHA")
	}
	if this.store != nil {
		result = append(result, "STORE", this.store.dest)
	}
	return result
}

func (this *Sorter) Limit(min, max int) *Sorter {
	this.limit = &sortLimit{
		min: min,
		max: max,
	}
	return this
}

func (this *Sorter) By(pattern string) *Sorter {
	this.by = &sortBy{
		pattern: pattern,
	}
	return this
}

func (this *Sorter) getFrom(pattern string) *Sorter {
	this.get = &sortGet{
		pattern: pattern,
	}
	return this
}

func (this *Sorter) storeIn(dest string) *Sorter {
	this.store = &sortStore{
		dest: dest,
	}
	return this
}

func (this *Sorter) Reverse() *Sorter {
	this.reversed = !this.reversed
	return this
}

func (this *Sorter) Get() <-chan []string {
	command, output := newSliceCommand(this.key.args("sort", this.sortargs()...))
	this.key.Execute(command)
	return output
}

func (this *Sorter) GetInts() <-chan []int {
	realoutput := make(chan []int, 1)
	midway := this.Get()
	go func() {
		defer close(realoutput)
		if output, ok := <-midway; ok {
			ints, err := stringsToInts(output)
			if err != nil {
				this.key.client.ErrCallback(err, "sorting ints")
			}
			realoutput <- ints
		}
	}()
	return realoutput
}

func (this *Sorter) GetFloats() <-chan []float64 {
	realoutput := make(chan []float64, 1)
	midway := this.Get()
	go func() {
		defer close(realoutput)
		if output, ok := <-midway; ok {
			floats, err := stringsToFloats(output)
			if err != nil {
				this.key.client.ErrCallback(err, "sorting floats")
			}
			realoutput <- floats
		}
	}()
	return realoutput
}

func (this *Sorter) GetFrom(pattern string) <-chan []*string {
	this.getFrom(pattern)
	command, output := newMaybeSliceCommand(this.key.args("sort", this.sortargs()...))
	this.key.Execute(command)
	return output
}

func (this *Sorter) GetIntsFrom(pattern string) <-chan []*int {
	realoutput := make(chan []*int, 1)
	midway := this.GetFrom(pattern)
	go func() {
		defer close(realoutput)
		if strings, ok := <-midway; ok {
			ints := make([]*int, len(strings))
			for i, str := range strings {
				if str != nil {
					j, err := atoi(*str)
					if err != nil {
						this.key.client.ErrCallback(err, "sorting ints")
					}
					ints[i] = &j
				}
			}

			realoutput <- ints
		}
	}()
	return realoutput
}

func (this *Sorter) GetFloatsFrom(pattern string) <-chan []*float64 {
	realoutput := make(chan []*float64, 1)
	midway := this.GetFrom(pattern)
	go func() {
		defer close(realoutput)
		if strings, ok := <-midway; ok {
			floats := make([]*float64, len(strings))
			for i, str := range strings {
				if str != nil {
					j, err := atof(*str)
					if err != nil {
						this.key.client.ErrCallback(err, "sorting floats")
					}
					floats[i] = &j
				}
			}

			realoutput <- floats
		}
	}()
	return realoutput
}

func (this *Sorter) StoreStrings(dest List) <-chan int {
	this.storeIn(dest.key)
	command, output := newIntCommand(this.key.args("sort", this.sortargs()...))
	this.key.Execute(command)
	return output
}

func (this *Sorter) StoreInts(dest IntList) <-chan int {
	this.storeIn(dest.key)
	command, output := newIntCommand(this.key.args("sort", this.sortargs()...))
	this.key.Execute(command)
	return output
}

func (this *Sorter) StoreFloats(dest FloatList) <-chan int {
	this.storeIn(dest.key)
	command, output := newIntCommand(this.key.args("sort", this.sortargs()...))
	this.key.Execute(command)
	return output
}

func (this *Sorter) GetFromAndStoreIn(pattern string, dest List) <-chan int {
	return this.getFrom(pattern).StoreStrings(dest)
}

func (this *Sorter) GetIntsFromAndStoreIn(pattern string, dest IntList) <-chan int {
	return this.getFrom(pattern).StoreInts(dest)
}

func (this *Sorter) GetFloatsFromAndStoreIn(pattern string, dest FloatList) <-chan int {
	return this.getFrom(pattern).StoreFloats(dest)
}
