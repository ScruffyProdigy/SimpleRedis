package redis

import (
	"testing"
)

func TestCommands(t *testing.T) {
	_, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	//TODO: test bad commands right here to make sure error handling works
}
