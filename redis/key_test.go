package redis

import (
	"testing"
	"time"
)

func TestKeys(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}
	defer r.Close()

	str := r.String("Test_Key")
	other_str := r.String("Other_Test_Key")
	str.Delete()
	other_str.Delete()

	if res := <-str.Type(); res != "none" {
		t.Error("Type should be none, not ", res)
	}

	if <-str.Exists() {
		t.Error("Key should not yet exist")
	}

	<-str.Set("A")
	<-str.MoveTo(other_str.Key)

	if res, ok := <-str.Get(); ok {
		t.Error("Should not have anything, and definitely not ", res)
	}
	if res, ok := <-other_str.Get(); !ok || res != "A" {
		t.Error("Should be A, not", res)
	}

	<-str.Set("B")
	if <-str.MoveToIfEmpty(other_str.Key) {
		t.Error("Should not be able to move when not empty")
	}
	if res, ok := <-str.Get(); !ok || res != "B" {
		t.Error("Should still be B, not", res)
	}
	if res, ok := <-other_str.Get(); !ok || res != "A" {
		t.Error("Should still be A, not", res)
	}

	<-other_str.Delete()
	if !<-str.MoveToIfEmpty(other_str.Key) {
		t.Error("Should be able to move now that empty")
	}
	if res, ok := <-str.Get(); ok {
		t.Error("Should not have anything anymore, instead has ", res)
	}
	if res, ok := <-other_str.Get(); !ok || res != "B" {
		t.Error("Should now be B, not", res)
	}

	if <-str.ExpireIn(time.Hour) {
		t.Error("Should not be able to set a TTL on a blank key")
	}

	<-str.Set("C")
	if !<-str.ExpireIn(time.Hour) {
		t.Error("Should be able to set a TTL now")
	}
	if res := <-str.MillisecondsToLive(); res < int(time.Hour/time.Millisecond)-50 || res > int(time.Hour/time.Millisecond) {
		t.Error("Should be about an hour's worth of milliseconds left, not", res)
	}
	if res := <-str.SecondsToLive(); res != int(time.Hour/time.Second) {
		t.Error("Should be an hour's worth of seconds left, not", res)
	}

	if !<-str.ExpireIn(time.Second) {
		t.Error("Should be able to change the TTL")
	}
	if res := <-str.MillisecondsToLive(); res < int(time.Second/time.Millisecond)-50 || res > int(time.Second/time.Millisecond) {
		t.Error("Should be about a second's worth of milliseconds left, not", res)
	}
	if res := <-str.SecondsToLive(); res != 1 {
		t.Error("Should one second left, not", res)
	}

	if res, ok := <-str.Get(); !ok || res != "C" {
		t.Error("Should still be C, not", res)
	}
	time.Sleep(time.Second)
	if res, ok := <-str.Get(); ok {
		t.Error("Should have expired, instead has ", res)
	}

	<-str.Set("D")
	<-str.ExpireAt(time.Now().Add(1 * time.Second))

	if res, ok := <-str.Get(); !ok || res != "D" {
		t.Error("Should still be D, not", res)
	}
	time.Sleep(time.Second)
	if res, ok := <-str.Get(); ok {
		t.Error("Should have expired, instead has ", res)
	}
}
