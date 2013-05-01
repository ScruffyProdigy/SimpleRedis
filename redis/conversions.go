package redis

import (
	"strconv"
)

func ftoa(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func atoi(s string) (int, error) {
	i, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		return 0, e
	}
	return int(i), nil
}

func atof(s string) (float64, error) {
	f, e := strconv.ParseFloat(s, 64)
	if e != nil {
		return 0, e
	}
	return f, nil
}

func intsToStrings(ints []int) []string {
	strings := make([]string, len(ints))
	for i := range ints {
		strings[i] = itoa(ints[i])
	}
	return strings
}

func stringsToInts(strings []string) ([]int, error) {
	ints := make([]int, len(strings))
	var err error
	for i := range strings {
		ints[i], err = atoi(strings[i])
		if err != nil {
			return nil, err
		}
	}
	return ints, nil
}

func floatsToStrings(floats []float64) []string {
	strings := make([]string, len(floats))
	for i := range floats {
		strings[i] = ftoa(floats[i])
	}
	return strings
}

func stringsToFloats(strings []string) ([]float64, error) {
	floats := make([]float64, len(strings))
	var err error
	for i := range strings {
		floats[i], err = atof(strings[i])
		if err != nil {
			return nil, err
		}
	}
	return floats, nil
}

func intsChannel(in <-chan []string) <-chan []int {
	out := make(chan []int, 1)
	go func() {
		defer close(out)
		if slice, ok := <-in; ok {
			if ints, err := stringsToInts(slice); err == nil {
				out <- ints
			}
		}
	}()
	return out
}

func floatsChannel(in <-chan []string) <-chan []float64 {
	out := make(chan []float64, 1)
	go func() {
		defer close(out)
		if slice, ok := <-in; ok {
			if floats, err := stringsToFloats(slice); err == nil {
				out <- floats
			}
		}
	}()
	return out
}

func stringChannel(in <-chan []string, index int) <-chan string {
	out := make(chan string, 1)
	go func() {
		defer close(out)
		if slice, ok := <-in; ok {
			out <- slice[index]
		}
	}()
	return out
}

func intChannel(in <-chan []string, index int) <-chan int {
	out := make(chan int, 1)
	go func() {
		defer close(out)
		if slice, ok := <-in; ok {
			if i, err := atoi(slice[index]); err == nil {
				out <- i
			}
		}
	}()
	return out
}

func floatChannel(in <-chan []string, index int) <-chan float64 {
	out := make(chan float64, 1)
	go func() {
		defer close(out)
		if slice, ok := <-in; ok {
			if float, err := atof(slice[index]); err == nil {
				out <- float
			}
		}
	}()
	return out
}

func maybeIntsChannel(in <-chan []*string) <-chan []*int {
	out := make(chan []*int, 1)
	go func() {
		defer close(out)
		if strings, ok := <-in; ok {
			ints := make([]*int, len(strings))
			for i, str := range strings {
				if str != nil {
					if j, err := atoi(*str); err == nil {
						ints[i] = &j
					}
				}
			}
			out <- ints
		}
	}()
	return out
}

func maybeFloatsChannel(in <-chan []*string) <-chan []*float64 {
	out := make(chan []*float64, 1)
	go func() {
		defer close(out)
		if strings, ok := <-in; ok {
			floats := make([]*float64, len(strings))
			for i, str := range strings {
				if str != nil {
					if j, err := atof(*str); err == nil {
						floats[i] = &j
					}
				}
			}
			out <- floats
		}
	}()
	return out
}

func stringfloatMapChannel(in <-chan map[string]string) <-chan map[string]float64 {
	out := make(chan map[string]float64, 1)
	go func() {
		defer close(out)
		if strings, ok := <-in; ok {
			result := make(map[string]float64, len(strings))
			for k, v := range strings {
				if float, err := atof(v); err == nil {
					result[k] = float
				}
			}
			out <- result
		}
	}()
	return out
}

func intfloatMapChannel(in <-chan map[string]string) <-chan map[int]float64 {
	out := make(chan map[int]float64, 1)
	go func() {
		defer close(out)
		if strings, ok := <-in; ok {
			result := make(map[int]float64, len(strings))
			for k, v := range strings {
				index, err := atoi(k)
				val, err2 := atof(v)
				if err == nil && err2 == nil {
					result[index] = val
				}
			}
			out <- result
		}
	}()

	return out
}
