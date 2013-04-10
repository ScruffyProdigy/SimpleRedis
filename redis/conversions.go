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
