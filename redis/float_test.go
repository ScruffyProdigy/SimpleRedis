package redis

import (
	"testing"
)

func TestFloatFuncs(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	f := r.Float("Test_Float")

	<-f.Set(2.5)

	if <-f.GetSet(4.2) != 2.5 {
		t.Error("Should have been set to 2.5")
	}

	if <-f.Get() != 4.2 {
		t.Error("Should have been getset to 4.2")
	}

	if <-f.SetIfEmpty(3.7) {
		t.Error("Should not set when not empty")
	}

	<-f.Delete()

	if _, ok := <-f.Get(); ok {
		t.Error("Should not have anything to get after Delete")
	}

	if <-f.IncrementBy(2.3) != 2.3 {
		t.Error("Increment by 2.3 from nothing should give 2.3")
	}

	if <-f.IncrementBy(2.3) != 4.6 {
		t.Error("Increment by 2.3 from 2.3 should give 4.6")
	}

	if <-f.DecrementBy(1.1) != 3.5 {
		t.Error("Decrement by 1.1 from 4.6 should give 3.5")
	}
}
