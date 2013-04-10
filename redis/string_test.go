package redis

import (
	"testing"
)

func TestStringFuncs(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't Load Redis")
	}

	s := r.String("Test_String")

	<-s.Set("Blah")

	if <-s.Get() != "Blah" {
		t.Error("Didn't get what we set")
	}

	if <-s.SetIfEmpty("Blah Blah") {
		t.Error("Shouldn't 'Set if empty' when not empty")
	}

	if <-s.Clear() != "Blah" {
		t.Error("Should still be blah")
	}

	if _, ok := <-s.Get(); ok {
		t.Error("Getting something after we clear")
	}

	if !<-s.SetIfEmpty("Blah Blah") {
		t.Error("Should 'Set if empty' when empty")
	}

	if <-s.Replace("Blah") != "Blah Blah" {
		t.Error("Should have been set")
	}

	if <-s.Get() != "Blah" {
		t.Error("Should have been replaced")
	}

	if <-s.Append(" Blah") != 9 {
		t.Error("Append should return strlen")
	}

	if <-s.Get() != "Blah Blah" {
		t.Error("Append Failed")
	}

	if <-s.Length() != 9 {
		t.Error("Length incorrect")
	}
}
