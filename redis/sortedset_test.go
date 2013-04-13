package redis

import (
	"testing"
)

func TestSortedSets(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	_ = r.SortedSet("Test_SortedSetTest")
}
