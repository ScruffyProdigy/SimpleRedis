package redis

import (
	"testing"
)

func TestFloatSets(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	_ = r.FloatSet("Test_FloatSet")
}
