package redis

import (
	"testing"
	"time"
)

func TestLists(t *testing.T) {
	r, err := New(DefaultConfiguration())
	if err != nil {
		t.Fatal("Can't load redis")
	}

	list := r.List("Test_List")

	list.Delete()

	if <-list.LeftPushIfExists("A") != 0 {
		t.Error("LPUSHX - Shouldn't left push when doesn't exist yet")
	}
	if <-list.RightPushIfExists("B") != 0 {
		t.Error("RPUSHX - Shouldn't right push when doesn't exist yet")
	}
	if _, ok := <-list.LeftPop(); ok {
		t.Error("Shouldn't have anything to pop yet")
	}
	if _, ok := <-list.RightPop(); ok {
		t.Error("SHouldn't have anything to pop yet")
	}

	if res := <-list.LeftPush("C"); res != 1 { //C
		t.Error("LPUSH - Length should be at 1, not", res)
	}
	if res := <-list.RightPush("D"); res != 2 { //CD
		t.Error("RPUSH - Length should be at 2, not", res)
	}
	if res := <-list.LeftPushIfExists("E"); res != 3 { //ECD
		t.Error("LPUSHX - Length should be at 3, not", res)
	}
	if res := <-list.RightPushIfExists("F"); res != 4 { //ECDF
		t.Error("RPUSHX - Length should be at 4, not", res)
	}

	if res := <-list.Length(); res != 4 {
		t.Error("LLEN - Length should still be at 4, not", res)
	}

	if res := <-list.Index(1); res != "C" {
		t.Error("LINDEX - Index 1 should be C, not", res)
	}
	if _, ok := <-list.Index(5); ok {
		t.Error("LINDEX - There should be no index 5")
	}

	if res := <-list.LeftPop(); res != "E" { //CDF
		t.Error("LPOP - E should still be on the left, not", res)
	}
	if res := <-list.RightPop(); res != "F" { //CD
		t.Error("RPOP - F should still be on the right, not", res)
	}

	if res := <-list.InsertBefore("D", "G"); res != 3 { //CGD
		t.Error("LINSERT - Length should be at 3, not", res)
	}
	if res := <-list.InsertBefore("A", "H"); res != -1 {
		t.Error("LINSERT - should not insert when pivot not found, got", res)
	}
	if res := <-list.InsertAfter("G", "I"); res != 4 { //CGID
		t.Error("LINSERT - Length should be at 4, not", res)
	}
	if res := <-list.InsertAfter("B", "J"); res != -1 {
		t.Error("LINSERT - should not insert when pivot not found, got", res)
	}

	<-list.Set(3, "K") //CGIK
	if res := <-list.Index(3); res != "K" {
		t.Error("LSET - should have set to K, got", res)
	}

	/*//	Currently this throws an error, TODO: make this work (not sure if can easily)
	if _,ok := <-list.Set(10,"L"); ok {
		t.Error("Should not work")
	}
	*/

	<-list.LeftPush("G")  //GCGIK
	<-list.RightPush("C") //GCGIKC
	<-list.RightPush("I") //GCGIKCI

	if res := <-list.Remove("G"); res != 2 { //CIKCI
		t.Error("LREM (0) - should have removed 2 items")
	}
	if res := <-list.RemoveNFromLeft(1, "C"); res != 1 { //IKCI
		t.Error("LREM (+) - should have only removed 1 item")
	}
	if res := <-list.Index(2); res != "C" {
		t.Error("LREM (+) - should have removed the leftmost C - got", res)
	}
	if res := <-list.RemoveNFromRight(1, "I"); res != 1 { //IKC
		t.Error("LREM (-) - should have only removed 1 item")
	}
	if res := <-list.Index(-3); res != "I" {
		t.Error("LREM (-) - should have removed the rightmost I - got", res)
	}

	if res := <-list.GetFromRange(0, -1); len(res) != 3 || res[0] != "I" || res[1] != "K" || res[2] != "C" {
		t.Error("LRANGE - should have [I K C], instead, have", res)
	}

	<-list.TrimToRange(1, -2)

	if res := <-list.GetFromRange(0, -1); len(res) != 1 || res[0] != "K" {
		t.Error("LRANGE - should have [K], instead, have", res)
	}

	if res := <-list.BlockUntilLeftPop(); res != "K" {
		t.Error("BLPOP - should have gotten K instead of", res)
	}
	print("Working")
	if _, ok := <-list.BlockUntilLeftPopWithTimeout(1); ok {
		t.Error("BLPOP - should not have anything to get")
	}
	print(".")

	time.AfterFunc(1*time.Second, func() {
		list.LeftPush("M")
	})
	select {
	case res := <-list.BlockUntilLeftPop():
		if res != "M" {
			t.Error("BLPOP - Should get M, not", res)
		}
	case <-time.After(2 * time.Second):
		t.Error("BLPOP - didn't receive anything after 2 seconds")
	}
	print(".")

	time.AfterFunc(1*time.Second, func() {
		list.LeftPush("N")
	})
	if res, ok := <-list.BlockUntilLeftPopWithTimeout(2); !ok || res != "N" {
		if !ok {
			t.Error("BLPOP - didn't receive anything after 2 seconds")
		} else if res != "N" {
			t.Error("BLPOP - should get N, not", res)
		}
	}
	print(".")

	time.AfterFunc(1*time.Second, func() {
		list.LeftPush("O")
	})
	select {
	case res := <-list.BlockUntilRightPop():
		if res != "O" {
			t.Error("BRPOP - Should get O, not", res)
		}
	case <-time.After(2 * time.Second):
		t.Error("BRPOP - didn't receive anything after 2 seconds")
	}
	print(".")

	time.AfterFunc(1*time.Second, func() {
		list.LeftPush("P")
	})
	if res, ok := <-list.BlockUntilRightPopWithTimeout(2); !ok || res != "P" {
		if !ok {
			t.Error("BRPOP - didn't receive anything after 2 seconds")
		} else if res != "P" {
			t.Error("BRPOP - should get P, not", res)
		}
	}
	print(".")

	otherlist := r.List("Other_Test_List")
	<-otherlist.Delete()

	if _, ok := <-list.MoveLastItemToList(otherlist); ok {
		t.Error("RPOPLPUSH - Should not have any items to move")
	}
	<-list.LeftPush("Q")
	if res := <-list.MoveLastItemToList(otherlist); res != "Q" {
		t.Error("RPOPLPUSH - Should have moved Q to new list")
	}
	if <-otherlist.Index(-1) != "Q" {
		t.Error("RPOPLPUSH - Never received Q in new list")
	}

	if _, ok := <-list.BlockUntilMoveLastItemToListWithTimeout(otherlist, 1); ok {
		t.Error("BRPOPLPUSH - Should not have any items to move")
	}
	print(".")

	time.AfterFunc(1*time.Second, func() {
		list.LeftPush("R")
	})
	select {
	case res := <-list.BlockUntilMoveLastItemToList(otherlist):
		if res != "R" {
			t.Error("BRPOPLPUSH - Should get R, not", res)
		}
		if <-otherlist.Index(-2) != "R" {
			t.Error("RPOPLPUSH - Never received R in new list")
		}
	case <-time.After(2 * time.Second):
		t.Error("BRPOPLPUSH - didn't receive anything after 2 seconds")
	}
	print(".")

	time.AfterFunc(1*time.Second, func() {
		list.LeftPush("S")
	})
	if res, ok := <-list.BlockUntilMoveLastItemToListWithTimeout(otherlist, 2); !ok || res != "S" {
		if !ok {
			t.Error("BRPOPLPUSH - didn't receive anything after 2 seconds")
		} else {
			if res != "S" {
				t.Error("BRPOPLPUSH - should get S, not", res)
			}
			if <-otherlist.Index(-3) != "S" {
				t.Error("RPOPLPUSH - Never received S in new list")
			}
		}
	}
	print(".\n")
}
