package redis

import (
	"testing"
)

func TestFloatLists(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	_ = r.FloatList("Test_FloatList")
}
