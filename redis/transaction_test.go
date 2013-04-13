package redis

import (
	"testing"
)

func TestPipeline(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	r.Pipeline(func(e Executor) {

	})
}

func TestTransaction(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	r.Transaction(func(e Executor) {

	})

	r.Transaction(func(e Executor) {
		panic(nil)
	})
}
