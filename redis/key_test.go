package redis

import (
	"testing"
)

func TestKeys(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	_ = r.String("Test_Key")
}
