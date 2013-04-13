package redis

import (
	"testing"
)

func TestSets(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	_ = r.Set("Test_Set")
}
