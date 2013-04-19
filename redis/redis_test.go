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
	return r
}

func TestBadCommands(t *testing.T) {
	r := GetRedis(t)
	defer r.Close()

	//TODO: test bad commands right here to make sure error handling works

}

func TestWrongPassword(t *testing.T) {
	r, err := Load(bytes.NewBuffer([]byte("{\"password\":\"wrong-password\"")))
	if err == nil {
		r.Close()
		t.Fatal("Should not work with wrong password")
	}
}
