package redis

import (
	"testing"
)

func TestIntSets(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}
	defer r.Close()

	set := r.IntSet("Test_IntSet")
	<-set.Delete()

	if res := <-set.Size(); res != 0 {
		t.Error("Set should have 0 members to start out, not", res)
	}
	if res := <-set.Members(); len(res) != 0 {
		t.Error("Should get an empty slice from an empty set")
	}
	if <-set.IsMember(1) {
		t.Error("1 should not be a member of the empty set")
	}

	if <-set.Remove(1) {
		t.Error("Should not be able to remove non-existant 1")
	}
	if !<-set.Add(1) {
		t.Error("Should be able to add 1")
	}
	if <-set.Add(1) {
		t.Error("1 should already exist, should not be able to add")
	}

	if res := <-set.Size(); res != 1 {
		t.Error("Set should have 1 member")
	}
	if res := <-set.Members(); len(res) != 1 || res[0] != 1 {
		t.Error("Should get a slice with 1 item, which is a 1")
	}
	if !<-set.IsMember(1) {
		t.Error("1 should be a member of the set")
	}

	if !<-set.Remove(1) {
		t.Error("1 should be removable now")
	}

	a := set.Add(1)
	b := set.Add(2)
	c := set.Add(3)
	<-a
	<-b
	<-c

	for i := 0; i < 100; i++ {
		if res := <-set.RandomMember(); res != 1 && res != 2 && res != 3 {
			t.Error("found incorrect member of set: ", res)
		}
	}

	for i := 0; i < 3; i++ {
		res := <-set.Pop()
		if res != 1 && res != 2 && res != 3 {
			t.Error("found incorrect member of set: ", res)
		}
		if <-set.IsMember(res) {
			t.Error(res, "should've been popped out of set")
		}
		if <-set.Size()+i != 2 {
			t.Error("There should only be", 2-i, "items left")
		}
	}

	otherset := r.IntSet("Other_Test_IntSet")
	<-otherset.Delete()

	a = set.Add(1)
	b = set.Add(2)
	c = otherset.Add(1)
	d := otherset.Add(3)
	<-a
	<-b
	<-c
	<-d

	done := make(chan bool)

	go func() {
		res := <-set.Intersection(otherset)
		if len(res) != 1 {
			t.Error("There should only be 1 item in the intersection, not", res)
		}
		if len(res) == 1 && res[0] != 1 {
			t.Error("1 should be in the intersection of the sets")
		}

		done <- true
	}()

	go func() {
		res := <-set.Union(otherset)
		if len(res) != 3 {
			t.Error("There should be 3 items in the union, not", len(res), ":", res)
		}
		if len(res) == 3 && res[0] != 1 && res[1] != 1 && res[2] != 1 {
			t.Error("1 should be in the union of the sets:", res)
		}
		if len(res) == 3 && res[0] != 2 && res[1] != 2 && res[2] != 2 {
			t.Error("2 should be in the union of the sets:", res)
		}
		if len(res) == 3 && res[0] != 3 && res[1] != 3 && res[2] != 3 {
			t.Error("3 should be in the union of the sets:", res)
		}

		done <- true
	}()

	go func() {
		res := <-set.Difference(otherset)
		if len(res) != 1 {
			t.Error("There should only be 1 item in the difference, not", len(res), ":", res)
		}
		if len(res) == 1 && res[0] != 2 {
			t.Error("2 should be in the difference of the sets:", res)
		}

		done <- true
	}()

	go func() {
		inter := r.IntSet("Intersection_IntSet")

		if res := <-inter.StoreIntersectionOf(set, otherset); res != 1 {
			t.Error("There should only be 1 item in the intersection, not", res)
		}

		res := <-inter.Members()
		if len(res) != 1 {
			t.Error("Result slice should have 1 item, not", len(res))
		}
		if len(res) == 1 && res[0] != 1 {
			t.Error("1 should be in the intersection of the sets:", res)
		}

		done <- true
	}()

	go func() {
		union := r.IntSet("Union_IntSet")

		if res := <-union.StoreUnionOf(set, otherset); res != 3 {
			t.Error("There should be 3 items in the union, not", res)
		}

		res := <-union.Members()
		if len(res) != 3 {
			t.Error("Result slice should have 3 members, not", len(res))
		}
		if len(res) == 3 && res[0] != 1 && res[1] != 1 && res[2] != 1 {
			t.Error("1 should be in the union of the sets:", res)
		}
		if len(res) == 3 && res[0] != 2 && res[1] != 2 && res[2] != 2 {
			t.Error("2 should be in the union of the sets:", res)
		}
		if len(res) == 3 && res[0] != 3 && res[1] != 3 && res[2] != 3 {
			t.Error("3 should be in the union of the sets:", res)
		}

		done <- true
	}()

	go func() {
		diff := r.IntSet("Diff_IntSet")

		if res := <-diff.StoreDifferenceOf(set, otherset); res != 1 {
			t.Error("There should only be 1 item in the difference, not", res)
		}

		res := <-diff.Members()
		if len(res) != 1 {
			t.Error("Result slice should have 1 member, not", len(res))
		}
		if len(res) == 1 && res[0] != 2 {
			t.Error("2 should be in the difference of the sets:", res)
		}

		done <- true
	}()

	for i := 0; i < 6; i++ {
		<-done
	}

	if !<-set.MoveMemberTo(otherset, 2) {
		t.Error("Should be able to move 2")
	}

	if res := <-otherset.Size(); res != 3 {
		t.Error("There should now be 3 members in the other set")
	}

	if res := <-set.Size(); res != 1 {
		t.Error("There should now only be 1 member in the base set")
	}

	if <-set.MoveMemberTo(otherset, 3) {
		t.Error("Should not be able to move 3 (as it doesn't exist in the base set)")
	}

	if !<-set.MoveMemberTo(otherset, 1) {
		t.Error("Should be able to move 1 (even though it is already in the other set)")
	}

	if res := <-set.Size(); res != 0 {
		t.Error("There should now be no more members in the base set")
	}
}
