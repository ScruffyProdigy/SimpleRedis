package redis

import (
	"testing"
)

func TestPipeline(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}
	defer r.Close()

	a := r.String("Pipeline_Test_A")
	b := r.String("Pipeline_Test_B")
	c := r.String("Pipeline_Test_C")

	<-a.Delete()
	<-b.Delete()
	<-c.Delete()

	r.Pipeline(func(e Executor) {
		a.Use(e).Set("A")
		b.Use(e).Set("B")
		c.Use(e).Set("C")

		if _, ok := <-a.Get(); ok {
			t.Error("a should not be set yet")
		}
		if _, ok := <-b.Get(); ok {
			t.Error("b should not be set yet")
		}
		if _, ok := <-c.Get(); ok {
			t.Error("c should not be set yet")
		}
	})

	if <-a.Get() != "A" {
		t.Error("a should be A")
	}

	if <-b.Get() != "B" {
		t.Error("b should be B")
	}

	if <-c.Get() != "C" {
		t.Error("c should be C")
	}

}

func TestTransaction(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}
	defer r.Close()

	a := r.String("Transaction_Test_A")
	b := r.String("Transaction_Test_B")
	c := r.String("Transaction_Test_C")

	<-a.Delete()
	<-b.Delete()
	<-c.Delete()

	r.Transaction(func(e Executor) {
		a.Use(e).Set("A")
		b.Use(e).Set("B")
		c.Use(e).Set("C")
		if _, ok := <-a.Get(); ok {
			t.Error("a should not be set yet")
		}
		if _, ok := <-b.Get(); ok {
			t.Error("b should not be set yet")
		}
		if _, ok := <-c.Get(); ok {
			t.Error("c should not be set yet")
		}

	})

	if <-a.Get() != "A" {
		t.Error("a should be A")
	}

	if <-b.Get() != "B" {
		t.Error("b should be B")
	}

	if <-c.Get() != "C" {
		t.Error("c should be C")
	}

	r.Transaction(func(e Executor) {
		a.Use(e).Set("D")
		b.Use(e).Set("E")
		c.Use(e).Set("F")

		if <-a.Get() != "A" {
			t.Error("a should be A")
		}
		if <-b.Get() != "B" {
			t.Error("b should be B")
		}
		if <-c.Get() != "C" {
			t.Error("c should be C")
		}

		panic("let's just discard these actions")
	})

	if <-a.Get() != "A" {
		t.Error("a should be A")
	}

	if <-b.Get() != "B" {
		t.Error("b should be B")
	}

	if <-c.Get() != "C" {
		t.Error("c should be C")
	}
}
