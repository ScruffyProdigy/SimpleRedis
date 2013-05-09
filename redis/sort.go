package redis

//SortableKey is a base type used to give other types the functions within here
//Redis's Sort function works with multiple types of keys, and in order to prevent too much code duplication, I've just created a base type
//See http://redis.io/commands/sort for more information on Redis's sort
type SortableKey struct {
	Key
}

func newSortableKey(client SafeExecutor, key string) SortableKey {
	return SortableKey{
		newKey(client, key),
	}
}

//SortAlphabetically will define a search in which redis sorts strings
func (this SortableKey) SortAlphabetically() *Sorter {
	return &Sorter{key: this.Key, alpha: true}
}

//SortNumerially will define a search in which redis sorts numbers
func (this SortableKey) SortNumerically() *Sorter {
	return &Sorter{key: this.Key, alpha: false}
}

type sortLimit struct {
	offset, count int
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

//Sorter keeps track of the options you want to use to sort with
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
		result = append(result, "LIMIT", itoa(this.limit.offset), itoa(this.limit.count))
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

//Limit defines the total number of results you want to receive back.
//It will skip the first "offset" results, and then only show you the next "count" results.
//
//Example: If you want to see the top 3 results, you would use field.SortNumerically().Limit(0,3)
func (this *Sorter) Limit(offset, count int) *Sorter {
	this.limit = &sortLimit{
		offset: offset,
		count:  count,
	}
	return this
}

//By allows you to use the current key as an index into a different set of keys
//
//Example: If you have a Set with {1,2,3,4,5}, and you sort By("string_*"), you will sort whatever string primitives are at string_1, string_2, string_3, string_4, and string_5
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

//Reverse will invert the order that you receive the results in
//
//Example: field.SortNumerically().Reverse() will define a descending search rather than an ascending one
func (this *Sorter) Reverse() *Sorter {
	this.reversed = !this.reversed
	return this
}

//Get will execute the search specified and return the result as a slice of strings
func (this *Sorter) Get() <-chan []string {
	return SliceCommand(this.key, this.key.args("sort", this.sortargs()...)...)
}

//GetInts will execute the search specified and return the result as a slice of integers
func (this *Sorter) GetInts() <-chan []int {
	return intsChannel(this.Get())
}

//GetFloats will execute the search specified and return the result as a slice of Floats
func (this *Sorter) GetFloats() <-chan []float64 {
	return floatsChannel(this.Get())
}

//GetFrom will execute the search, but instead of returning the results, will use the results to dig into other string primitives.
//It is the equivalent of using a GET argument in the sort
func (this *Sorter) GetFrom(pattern string) <-chan []*string {
	this.getFrom(pattern)
	return MaybeSliceCommand(this.key, this.key.args("sort", this.sortargs()...)...)
}

//GetFrom will execute the search, but instead of returning the results, will use the results to dig into other string primitives containing (hopefully) integers.
//It is the equivalent of using a GET argument in the sort
func (this *Sorter) GetIntsFrom(pattern string) <-chan []*int {
	return maybeIntsChannel(this.GetFrom(pattern))
}

//GetFrom will execute the search, but instead of returning the results, will use the results to dig into other string primitives containing (hopefully) floating point numbers.
//It is the equivalent of using a GET argument in the sort
func (this *Sorter) GetFloatsFrom(pattern string) <-chan []*float64 {
	return maybeFloatsChannel(this.GetFrom(pattern))
}

//StoreStrings will execute the sort, but instead of returning the results will store them in a list primitive.
//It is the equivalent of using a STORE argument
func (this *Sorter) StoreStrings(dest List) <-chan int {
	this.storeIn(dest.key)
	return IntCommand(this.key, this.key.args("sort", this.sortargs()...)...)
}

//StoreInts will execute the sort, but instead of returning the results will store them in a list primitive.
//It is the equivalent of using a STORE argument
func (this *Sorter) StoreInts(dest IntList) <-chan int {
	this.storeIn(dest.key)
	return IntCommand(this.key, this.key.args("sort", this.sortargs()...)...)
}

//GetFromAndStoreIn is like using both GetFrom and StoreStrings
func (this *Sorter) GetFromAndStoreIn(pattern string, dest List) <-chan int {
	return this.getFrom(pattern).StoreStrings(dest)
}

//GetIntsFromAndStoreIn is like using both GetIntsFrom and StoreInts
func (this *Sorter) GetIntsFromAndStoreIn(pattern string, dest IntList) <-chan int {
	return this.getFrom(pattern).StoreInts(dest)
}
