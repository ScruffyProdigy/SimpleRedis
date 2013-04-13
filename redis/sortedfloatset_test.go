package redis

import (
	"testing"
)

func TestSortedFloatSets(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	_ = r.SortedFloatSet("Test_SortedFloatTest")
}
