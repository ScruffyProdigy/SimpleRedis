package redis

import (
	"testing"
)

func TestHashes(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	_ = r.Hash("Test_Hash")
}
