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

func atoi(s string) int {
	i, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		panic(e.Error() + "\"" + s + "\"")
	}
	return int(i)
}

func atof(s string) float64 {
	f, e := strconv.ParseFloat(s, 64)
	if e != nil {
		panic("Invalid Float")
	}
	return f
}

func intsToStrings(ints []int) []string {
	strings := make([]string, len(ints))
	for i := range ints {
		strings[i] = itoa(ints[i])
	}
	return strings
}

func stringsToInts(strings []string) []int {
	ints := make([]int, len(strings))
	for i := range strings {
		ints[i] = atoi(strings[i])
	}
	return ints
}

func floatsToStrings(floats []float64) []string {
	strings := make([]string, len(floats))
	for i := range floats {
		strings[i] = ftoa(floats[i])
	}
	return strings
}

func stringsToFloats(strings []string) []float64 {
	floats := make([]float64, len(strings))
	for i := range strings {
		floats[i] = atof(strings[i])
	}
	return floats
}
