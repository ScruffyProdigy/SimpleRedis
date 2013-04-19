package redis

import (
	"testing"
)

func TestIntFuncs(t *testing.T) {
	r := GetRedis(t)
	defer r.Close()

	i := r.Integer("Test_Integer")

	<-i.Set(3)

	if <-i.GetSet(5) != 3 {
		t.Error("Should have set to 3")
	}

	if <-i.Get() != 5 {
		t.Error("Should have getset to 5")
	}

	if <-i.SetIfEmpty(7) {
		t.Error("SetIfEmpty should not set when not empty")
	}

	if <-i.Get() != 5 {
		t.Error("SetIfEmpty should not have affected value")
	}

	<-i.Delete()

	if _, ok := <-i.Get(); ok {
		t.Error("Should not have anything to get (Delete Fail)")
	}

	if !<-i.SetIfEmpty(7) {
		t.Error("SetIfEmpty should set when empty")
	}

	if <-i.Get() != 7 {
		t.Error("SetIfEmpty failed to set")
	}

	if <-i.Increment() != 8 {
		t.Error("Increment should give you 1 more than before")
	}

	if <-i.IncrementBy(2) != 10 {
		t.Error("IncremementBy 2 should yield 2 more than before")
	}

	if <-i.Decrement() != 9 {
		t.Error("Decrement should give you 1 less than before")
	}

	if <-i.DecrementBy(2) != 7 {
		t.Error("DecrementBy 2 should give you 2 less than before")
	}

}
