package redis

import (
	"testing"
)

func TestSorting(t *testing.T) {
	_, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	//we have a lot to test here
}
