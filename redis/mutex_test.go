package redis

import (
	"testing"
)

func TestMutices(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	_ = r.Mutex("Test_Mutex")
	_ = r.Semaphore("Test_Semaphore", 2)
}
