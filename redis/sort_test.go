package redis

import (
	"testing"
)

func TestSorting(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}
	defer r.Close()

	//we have a lot to test here

}
