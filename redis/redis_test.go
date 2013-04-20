package redis

import (
	"bytes"
	"testing"
)

// GetRedis is meant to provide a common way for every test function to log into redis the same way
// ie redis requires a password, or you want all of the tests to run on a selected dbid so you can easily flush it later
func GetRedis(t *testing.T) *Client {
	r, err := New(DefaultConfiguration())
	//	r, err := Load(bytes.NewBuffer([]byte("{\"password\":\"password\",\"dbid\":1}")))
	if err != nil {
		t.Fatal("Can't load redis - " + err.Error())
	}
	r.SetErrorCallback(func(e error, s string) {
		t.Error(e.Error() + " - " + s)
	})
	return r
}

func TestBadCommands(t *testing.T) {
	failed := make(chan bool)
	r := GetRedis(t)
	defer r.Close()
	r.SetErrorCallback(func(e error, s string) {
		failed <- true
	})

	if _, ok := <-NilCommand(r, []string{"INVALIDCOMMAND"}); ok {
		t.Error("Should not get *ANYTHING* back")
	}
	select {
	case <-failed:
	default:
		t.Error("Using an invalid command should cause an error")
	}

	s := r.String("ErrorTest")

	var ch, ch2 <-chan nothing
	r.Pipeline(func(e SafeExecutor) {
		ch = NilCommand(e, []string{"INVALIDCOMAND"})
		ch2 = s.Use(e).Set("Test Test")
	})
	if _, ok := <-ch; ok {
		t.Error("Still should not get *ANYTHING* back")
	}
	if _, ok := <-ch2; !ok {
		t.Error("Second command should still work fine")
	}

	if res := <-s.Get(); res != "Test Test" {
		t.Error("Should have gotten 'Test Test', not ", res)
	}

}

func TestWrongPassword(t *testing.T) {
	r, err := Load(bytes.NewBuffer([]byte("{\"password\":\"wrong-password\"")))
	if err == nil {
		r.Close()
		t.Fatal("Should not work with wrong password")
	}
}
