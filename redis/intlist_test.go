package redis

import (
	"testing"
)

func TestIntLists(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	_ = r.IntList("Test_IntList")
}
