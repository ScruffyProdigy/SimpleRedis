package redis

import (
	"testing"
)

func TestSortedIntSets(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	_ = r.SortedIntSet("Test_SortedIntTest")
}
