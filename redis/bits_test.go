package redis

import (
	"testing"
)

func TestBitsFuncs(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}
	defer r.Close()

	b := r.Bits("Test_Bits")

	<-b.On(2)
	<-b.Off(3)
	<-b.SetTo(4, true)
	<-b.SetTo(5, false)
	<-b.On(15)

	if !<-b.Off(2) {
		t.Error("Off should return old value")
	}
	if <-b.On(2) {
		t.Error("On should return old value")
	}
	if <-b.Off(3) {
		t.Error("Off should return old value")
	}
	if !<-b.On(2) {
		t.Error("On should return old value")
	}
	if !<-b.Get(2) {
		t.Error("Get should return true for a set bit")
	}
	if <-b.Get(3) {
		t.Error("Get should return false for an unset bit")
	}
	if !<-b.Get(4) {
		t.Error("Get should return true for a set bit")
	}
	if <-b.Get(5) {
		t.Error("Get should return false for an unset bit")
	}
	if <-b.Count(0, 32) != 3 {
		t.Error("There should be 3 bits set")
	}
	<-b.Off(4)

	a := r.Bits("Test_Bits_2")
	c := r.Bits("Test_Bits_3")

	<-a.On(2)
	<-a.Off(3)
	<-a.Off(15)
	<-a.On(5)

	x := <-a.And(b, c)
	if x != 2 {
		t.Error("we're using 2 characters, not", x)
	}
	if !<-c.Get(2) {
		t.Error("2nd bit should be set")
	}
	if <-c.Get(3) {
		t.Error("3rd bit should not be set")
	}
	if <-c.Get(15) {
		t.Error("15th bit should not be set")
	}
	if <-c.Get(5) {
		t.Error("5th bit should not be set")
	}

	x = <-a.Or(b, c)
	if x != 2 {
		t.Error("we're using 2 characters, not", x)
	}
	if !<-c.Get(2) {
		t.Error("2nd bit should be set")
	}
	if <-c.Get(3) {
		t.Error("3rd bit should not be set")
	}
	if !<-c.Get(15) {
		t.Error("15th bit should be set")
	}
	if !<-c.Get(5) {
		t.Error("5th bit should be set")
	}

	x = <-a.Xor(b, c)
	if x != 2 {
		t.Error("we're using 2 characters, not", x)
	}
	if <-c.Get(2) {
		t.Error("2nd bit should not be set")
	}
	if <-c.Get(3) {
		t.Error("3rd bit should not be set")
	}
	if !<-c.Get(15) {
		t.Error("15th bit should be set")
	}
	if !<-c.Get(5) {
		t.Error("5th bit should be set")
	}

	x = <-a.Not(c)
	if x != 2 {
		t.Error("we're using 2 characters, not", x)
	}
	if <-c.Get(2) {
		t.Error("2nd bit should not be set")
	}
	if !<-c.Get(3) {
		t.Error("3rd bit should be set")
	}

}
