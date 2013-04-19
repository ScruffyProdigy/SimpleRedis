package redis

import (
	"testing"
)

func TestSets(t *testing.T) {
	r := GetRedis(t)
	defer r.Close()

	set := r.Set("Test_Set")
	<-set.Delete()

	if res := <-set.Size(); res != 0 {
		t.Error("Set should have 0 members to start out, not", res)
	}
	if res := <-set.Members(); len(res) != 0 {
		t.Error("Should get an empty slice from an empty set")
	}
	if <-set.IsMember("A") {
		t.Error("A should not be a member of the empty set")
	}

	if <-set.Remove("A") {
		t.Error("Should not be able to remove non-existant A")
	}
	if !<-set.Add("A") {
		t.Error("Should be able to add A")
	}
	if <-set.Add("A") {
		t.Error("A should already exist, should not be able to add")
	}

	if res := <-set.Size(); res != 1 {
		t.Error("Set should have 1 member")
	}
	if res := <-set.Members(); len(res) != 1 || res[0] != "A" {
		t.Error("Should get a slice with 1 item, which is an \"A\"")
	}
	if !<-set.IsMember("A") {
		t.Error("A should be a member of the set")
	}

	if !<-set.Remove("A") {
		t.Error("A should be removable now")
	}

	a := set.Add("A")
	b := set.Add("B")
	c := set.Add("C")
	<-a
	<-b
	<-c

	for i := 0; i < 100; i++ {
		if res := <-set.RandomMember(); res != "A" && res != "B" && res != "C" {
			t.Error("found incorrect member of set: ", res)
		}
	}

	for i := 0; i < 3; i++ {
		res := <-set.Pop()
		if res != "A" && res != "B" && res != "C" {
			t.Error("found incorrect member of set: ", res)
		}
		if <-set.IsMember(res) {
			t.Error(res, "should've been popped out of set")
		}
		if <-set.Size()+i != 2 {
			t.Error("There should only be", 2-i, "items left")
		}
	}

	otherset := r.Set("Other_Test_Set")
	<-otherset.Delete()

	a = set.Add("A")
	b = set.Add("B")
	c = otherset.Add("A")
	d := otherset.Add("C")
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
		if len(res) == 1 && res[0] != "A" {
			t.Error("A should be in the intersection of the sets")
		}

		done <- true
	}()

	go func() {
		res := <-set.Union(otherset)
		if len(res) != 3 {
			t.Error("There should be 3 items in the union, not", len(res), ":", res)
		}
		if len(res) == 3 && res[0] != "A" && res[1] != "A" && res[2] != "A" {
			t.Error("A should be in the union of the sets:", res)
		}
		if len(res) == 3 && res[0] != "B" && res[1] != "B" && res[2] != "B" {
			t.Error("B should be in the union of the sets:", res)
		}
		if len(res) == 3 && res[0] != "C" && res[1] != "C" && res[2] != "C" {
			t.Error("C should be in the union of the sets:", res)
		}

		done <- true
	}()

	go func() {
		res := <-set.Difference(otherset)
		if len(res) != 1 {
			t.Error("There should only be 1 item in the difference, not", len(res), ":", res)
		}
		if len(res) == 1 && res[0] != "B" {
			t.Error("B should be in the difference of the sets:", res)
		}

		done <- true
	}()

	go func() {
		inter := r.Set("Intersection_Set")

		if res := <-inter.StoreIntersectionOf(set, otherset); res != 1 {
			t.Error("There should only be 1 item in the intersection, not", res)
		}

		res := <-inter.Members()
		if len(res) != 1 {
			t.Error("Result slice should have 1 item, not", len(res))
		}
		if len(res) == 1 && res[0] != "A" {
			t.Error("A should be in the intersection of the sets:", res)
		}

		done <- true
	}()

	go func() {
		union := r.Set("Union_Set")

		if res := <-union.StoreUnionOf(set, otherset); res != 3 {
			t.Error("There should be 3 items in the union, not", res)
		}

		res := <-union.Members()
		if len(res) != 3 {
			t.Error("Result slice should have 3 members, not", len(res))
		}
		if len(res) == 3 && res[0] != "A" && res[1] != "A" && res[2] != "A" {
			t.Error("A should be in the union of the sets:", res)
		}
		if len(res) == 3 && res[0] != "B" && res[1] != "B" && res[2] != "B" {
			t.Error("B should be in the union of the sets:", res)
		}
		if len(res) == 3 && res[0] != "C" && res[1] != "C" && res[2] != "C" {
			t.Error("C should be in the union of the sets:", res)
		}

		done <- true
	}()

	go func() {
		diff := r.Set("Diff_Set")

		if res := <-diff.StoreDifferenceOf(set, otherset); res != 1 {
			t.Error("There should only be 1 item in the difference, not", res)
		}

		res := <-diff.Members()
		if len(res) != 1 {
			t.Error("Result slice should have 1 member, not", len(res))
		}
		if len(res) == 1 && res[0] != "B" {
			t.Error("B should be in the difference of the sets:", res)
		}

		done <- true
	}()

	for i := 0; i < 6; i++ {
		<-done
	}

	if !<-set.MoveMemberTo(otherset, "B") {
		t.Error("Should be able to move B")
	}

	if res := <-otherset.Size(); res != 3 {
		t.Error("There should now be 3 members in the other set")
	}

	if res := <-set.Size(); res != 1 {
		t.Error("There should now only be 1 member in the base set")
	}

	if <-set.MoveMemberTo(otherset, "C") {
		t.Error("Should not be able to move C (as it doesn't exist in the base set)")
	}

	if !<-set.MoveMemberTo(otherset, "A") {
		t.Error("Should be able to move A (even though it is already in the other set)")
	}

	if res := <-set.Size(); res != 0 {
		t.Error("There should now be no more members in the base set")
	}
}
