package redis

import (
	"testing"
)

func TestKeys(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}
	defer r.Close()

	_ = r.String("Test_Key")
}
