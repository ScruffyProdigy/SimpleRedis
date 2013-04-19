package redis

import (
	"testing"
	"time"
)

func TestIntLists(t *testing.T) {
	r := GetRedis(t)
	defer r.Close()

	list := r.IntList("Test_IntList")

	list.Delete()

	if <-list.LeftPushIfExists(1) != 0 {
		t.Error("LPUSHX - Shouldn't left push when doesn't exist yet")
	}
	if <-list.RightPushIfExists(2) != 0 {
		t.Error("RPUSHX - Shouldn't right push when doesn't exist yet")
	}
	if _, ok := <-list.LeftPop(); ok {
		t.Error("Shouldn't have anything to pop yet")
	}
	if _, ok := <-list.RightPop(); ok {
		t.Error("SHouldn't have anything to pop yet")
	}

	if res := <-list.LeftPush(3); res != 1 { //C
		t.Error("LPUSH - Length should be at 1, not", res)
	}
	if res := <-list.RightPush(4); res != 2 { //CD
		t.Error("RPUSH - Length should be at 2, not", res)
	}
	if res := <-list.LeftPushIfExists(5); res != 3 { //ECD
		t.Error("LPUSHX - Length should be at 3, not", res)
	}
	if res := <-list.RightPushIfExists(6); res != 4 { //ECDF
		t.Error("RPUSHX - Length should be at 4, not", res)
	}

	if res := <-list.Length(); res != 4 {
		t.Error("LLEN - Length should still be at 4, not", res)
	}

	if res := <-list.Index(1); res != 3 {
		t.Error("LINDEX - Index 1 should be 3, not", res)
	}
	if _, ok := <-list.Index(5); ok {
		t.Error("LINDEX - There should be no index 5")
	}

	if res := <-list.LeftPop(); res != 5 { //CDF
		t.Error("LPOP - 5 should still be on the left, not", res)
	}
	if res := <-list.RightPop(); res != 6 { //CD
		t.Error("RPOP - 6 should still be on the right, not", res)
	}

	if res := <-list.InsertBefore(4, 7); res != 3 { //CGD
		t.Error("LINSERT - Length should be at 3, not", res)
	}
	if res := <-list.InsertBefore(1, 8); res != -1 {
		t.Error("LINSERT - should not insert when pivot not found, got", res)
	}
	if res := <-list.InsertAfter(7, 9); res != 4 { //CGID
		t.Error("LINSERT - Length should be at 4, not", res)
	}
	if res := <-list.InsertAfter(2, 10); res != -1 {
		t.Error("LINSERT - should not insert when pivot not found, got", res)
	}

	<-list.Set(3, 11) //CGIK
	if res := <-list.Index(3); res != 11 {
		t.Error("LSET - should have set to 11, got", res)
	}

	/*//	Currently this throws an error, TODO: make this work (not sure if can easily)
	if _,ok := <-list.Set(10,12); ok {
		t.Error("Should not work")
	}
	*/

	<-list.LeftPush(7)  //GCGIK
	<-list.RightPush(3) //GCGIKC
	<-list.RightPush(9) //GCGIKCI

	if res := <-list.Remove(7); res != 2 { //CIKCI
		t.Error("LREM (0) - should have removed 2 items")
	}
	if res := <-list.RemoveNFromLeft(1, 3); res != 1 { //IKCI
		t.Error("LREM (+) - should have only removed 1 item")
	}
	if res := <-list.Index(2); res != 3 {
		t.Error("LREM (+) - should have removed the leftmost 3 - got", res)
	}
	if res := <-list.RemoveNFromRight(1, 9); res != 1 { //IKC
		t.Error("LREM (-) - should have only removed 1 item")
	}
	if res := <-list.Index(-3); res != 9 {
		t.Error("LREM (-) - should have removed the rightmost 9 - got", res)
	}

	if res := <-list.GetFromRange(0, -1); len(res) != 3 || res[0] != 9 || res[1] != 11 || res[2] != 3 {
		t.Error("LRANGE - should have [9 11 3], instead, have", res)
	}

	<-list.TrimToRange(1, -2)

	if res := <-list.GetFromRange(0, -1); len(res) != 1 || res[0] != 11 {
		t.Error("LRANGE - should have [11], instead, have", res)
	}

	if res := <-list.BlockUntilLeftPop(); res != 11 {
		t.Error("BLPOP - should have gotten 11 instead of", res)
	}
	print("Working")
	if _, ok := <-list.BlockUntilLeftPopWithTimeout(1); ok {
		t.Error("BLPOP - should not have anything to get")
	}
	print(".")

	time.AfterFunc(1*time.Second, func() {
		list.LeftPush(13)
	})
	select {
	case res := <-list.BlockUntilLeftPop():
		if res != 13 {
			t.Error("BLPOP - Should get 13, not", res)
		}
	case <-time.After(2 * time.Second):
		t.Error("BLPOP - didn't receive anything after 2 seconds")
	}
	print(".")

	time.AfterFunc(1*time.Second, func() {
		list.LeftPush(14)
	})
	if res, ok := <-list.BlockUntilLeftPopWithTimeout(2); !ok || res != 14 {
		if !ok {
			t.Error("BLPOP - didn't receive anything after 2 seconds")
		} else if res != 14 {
			t.Error("BLPOP - should get 14, not", res)
		}
	}
	print(".")

	time.AfterFunc(1*time.Second, func() {
		list.LeftPush(15)
	})
	select {
	case res := <-list.BlockUntilRightPop():
		if res != 15 {
			t.Error("BRPOP - Should get 15, not", res)
		}
	case <-time.After(2 * time.Second):
		t.Error("BRPOP - didn't receive anything after 2 seconds")
	}
	print(".")

	time.AfterFunc(1*time.Second, func() {
		list.LeftPush(16)
	})
	if res, ok := <-list.BlockUntilRightPopWithTimeout(2); !ok || res != 16 {
		if !ok {
			t.Error("BRPOP - didn't receive anything after 2 seconds")
		} else if res != 16 {
			t.Error("BRPOP - should get 16, not", res)
		}
	}
	print(".")

	otherlist := r.IntList("Other_Test_IntList")
	<-otherlist.Delete()

	if _, ok := <-list.MoveLastItemToList(otherlist); ok {
		t.Error("RPOPLPUSH - Should not have any items to move")
	}
	<-list.LeftPush(17)
	if res := <-list.MoveLastItemToList(otherlist); res != 17 {
		t.Error("RPOPLPUSH - Should have moved 17 to new list")
	}
	if <-otherlist.Index(-1) != 17 {
		t.Error("RPOPLPUSH - Never received 17 in new list")
	}

	if _, ok := <-list.BlockUntilMoveLastItemToListWithTimeout(otherlist, 1); ok {
		t.Error("BRPOPLPUSH - Should not have any items to move")
	}
	print(".")

	time.AfterFunc(1*time.Second, func() {
		list.LeftPush(18)
	})
	select {
	case res := <-list.BlockUntilMoveLastItemToList(otherlist):
		if res != 18 {
			t.Error("BRPOPLPUSH - Should get 18, not", res)
		}
	case <-time.After(2 * time.Second):
		t.Error("BRPOPLPUSH - didn't receive anything after 2 seconds")
	}
	print(".")

	time.AfterFunc(1*time.Second, func() {
		list.LeftPush(19)
	})
	if res, ok := <-list.BlockUntilMoveLastItemToListWithTimeout(otherlist, 2); !ok || res != 19 {
		if !ok {
			t.Error("BRPOPLPUSH - didn't receive anything after 2 seconds")
		} else if res != 19 {
			t.Error("BRPOPLPUSH - should get 19, not", res)
		}
	}
	print(".\n")
}
