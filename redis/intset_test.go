package redis

import (
	"testing"
)

func TestIntSets(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	_ = r.IntSet("Test_IntSet")
}
