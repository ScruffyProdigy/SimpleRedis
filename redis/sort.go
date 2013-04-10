package redis

import (
	"errors"
)

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
		this.key.client.ErrCallback(errors.New("Argument error"), "Can't Tell If Should Sort Alphabetically or Numerically?!")
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
	this.limit = &sortLimit{
		min: min,
		max: max,
	}
	return this
}

func (this Sorter) Alphabetically() Sorter {
	this.alpha = true
	this.typeSet = true
	return this
}

func (this Sorter) Numerically() Sorter {
	this.typeSet = true
	return this
}

//if we could figure out what kind of pattern was referred to here, we could automatically set alphabetically or numerically
func (this Sorter) By(pattern string) Sorter {
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

func (this Sorter) Floats() <-chan []float64 {
	realoutput := make(chan []float64, 1)
	midway := this.Strings()
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

func (this Sorter) MaybeFloats() <-chan []*float64 {
	realoutput := make(chan []*float64, 1)
	midway := this.MaybeStrings()
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
