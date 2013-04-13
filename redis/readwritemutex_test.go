package redis

import (
	"testing"
)

func TestReadWriteMutices(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	_ = r.ReadWriteMutex("Test_ReadWriteMutex", 2)
}
