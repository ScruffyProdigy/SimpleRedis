package redis

import (
	"testing"
)

func TestCommands(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}
	defer r.Close()

	//TODO: test bad commands right here to make sure error handling works

}
