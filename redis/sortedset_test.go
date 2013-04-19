package redis

import (
	"testing"
)

func TestSortedSets(t *testing.T) {
	r := GetRedis(t)
	defer r.Close()

	ss := r.SortedSet("Test_SortedSetTest")
	ss.Delete()

	if _, ok := <-ss.ScoreOf("A"); ok {
		t.Error("Should not be any score to get")
	}
	if _, ok := <-ss.IndexOf("A"); ok {
		t.Error("Should not be any rank to get")
	}
	if res := <-ss.Size(); res != 0 {
		t.Error("Should have a size of 0, not", res)
	}

	a := ss.Add("A", 1)
	bb := ss.IncrementBy("B", 3)
	c := ss.Add("C", 2)

	<-a
	<-bb
	<-c

	if res := <-ss.ScoreOf("C"); res != 2 {
		t.Error("C should have a score of 2, not", res)
	}
	if res := <-ss.IndexOf("C"); res != 1 {
		t.Error("C should have an index of 1, not", res)
	}
	if res := <-ss.ReverseIndexOf("C"); res != 1 {
		t.Error("C should have a reverse index of 1, not", res)
	}

	if res := <-ss.Size(); res != 3 {
		t.Error("Should have a size of 3, not", res)
	}

	<-ss.IncrementBy("C", 2)

	if res := <-ss.ScoreOf("C"); res != 4 {
		t.Error("C should now have a score of 4, not", res)
	}
	if res := <-ss.IndexOf("C"); res != 2 {
		t.Error("C should now have an index of 2, not", res)
	}
	if res := <-ss.ReverseIndexOf("C"); res != 0 {
		t.Error("C should now have a reverse index of 0, not", res)
	}

	<-ss.Remove("B")

	if res := <-ss.ScoreOf("C"); res != 4 {
		t.Error("Score of C should be unaffected; should still be 4, not", res)
	}
	if res := <-ss.IndexOf("C"); res != 1 {
		t.Error("C should now have an index of 1, not", res)
	}
	if res := <-ss.ReverseIndexOf("C"); res != 0 {
		t.Error("C should still have a reverse index of 0, not", res)
	}

	//	"A" 1
	a = ss.Add("B", 9)
	//	"C" 4
	b := ss.Add("D", 6)
	c = ss.Add("E", 8)
	d := ss.Add("F", 3)
	e := ss.Add("G", 10)
	f := ss.Add("H", 2)
	g := ss.Add("I", 7)
	h := ss.Add("J", 5)

	<-a
	<-b
	<-c
	<-d
	<-e
	<-f
	<-g
	<-h

	done := make(chan bool)

	go func() {
		//get top three scores
		res := <-ss.ReverseIndexedBetween(0, 2)
		if len(res) != 3 {
			t.Error("Top 3 should have 3 members, not", len(res))
		} else {
			if res[0] != "G" {
				t.Error("Top score should be G, not", res[0])
			}
			if res[1] != "B" {
				t.Error("2nd place should be B, not", res[1])
			}
			if res[2] != "E" {
				t.Error("3rd place should be E, not", res[2])
			}
		}

		done <- true
	}()

	go func() {
		//get bottom three scores
		res := <-ss.IndexedBetween(0, 2)
		if len(res) != 3 {
			t.Error("Bottom 3 should have 3 members, not", len(res))
		} else {
			if res[0] != "A" {
				t.Error("Bottom score should be A, not", res[0])
			}
			if res[1] != "H" {
				t.Error("2nd worst should be H, not", res[1])
			}
			if res[2] != "F" {
				t.Error("3rd worst should be F, not", res[2])
			}
		}

		done <- true
	}()

	go func() {
		//get all scores at or above 7
		base := ss.Scores().AboveOrEqualTo(7)

		count := <-base.Count()
		if count != 4 {
			t.Error("Count should be 4, not", count)
		}

		get := <-base.Get()
		if len(get) != 4 || get[3] != "G" || get[2] != "B" || get[1] != "E" || get[0] != "I" {
			t.Error("Scores should be [I E B G], not", get)
		}

		getwithscores := <-base.GetWithScores()
		if len(getwithscores) != 4 || getwithscores["G"] != 10 || getwithscores["B"] != 9 || getwithscores["E"] != 8 || getwithscores["I"] != 7 {
			t.Error("Scores should be map[G:10 B:9 E:8 I:7], not", getwithscores)
		}

		done <- true
	}()

	go func() {
		//get all scores at or below 4
		base := ss.Scores().BelowOrEqualTo(4).Reversed()

		count := <-base.Count()
		if count != 4 {
			t.Error("Count should be 4, not", count)
		}

		get := <-base.Get()
		if len(get) != 4 || get[3] != "A" || get[2] != "H" || get[1] != "F" || get[0] != "C" {
			t.Error("Scores should be [C F H A], not", get)
		}

		getwithscores := <-base.GetWithScores()
		if len(getwithscores) != 4 || getwithscores["C"] != 4 || getwithscores["F"] != 3 || getwithscores["H"] != 2 || getwithscores["A"] != 1 {
			t.Error("Scores should be map[C:1 F:2 H:3 A:4], not", getwithscores)
		}

		done <- true
	}()

	//make sure Above/Below override AboveOrEqualTo/BelowOrEqualTo properly
	go func() {
		res := <-ss.Scores().AboveOrEqualTo(3).Above(3).Limit(0, 3).Get()
		if res[0] != "C" || res[1] != "J" || res[2] != "D" {
			t.Error("First 3 Above 3 should be [C J D], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().Above(3).AboveOrEqualTo(3).Limit(0, 3).GetWithScores()
		if res["C"] != 4 || res["J"] != 5 || res["D"] != 6 {
			t.Error("First 3 Above 3 should be map[C:4 J:5 D:6], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().BelowOrEqualTo(8).Below(8).Reversed().Limit(0, 3).Get()
		if res[0] != "I" || res[1] != "D" || res[2] != "J" {
			t.Error("First 3 below 8 should be [I D J], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().Below(8).BelowOrEqualTo(8).Reversed().Limit(0, 3).GetWithScores()
		if res["I"] != 7 || res["D"] != 6 || res["J"] != 5 {
			t.Error("First 3 Below 8 should be [I D J], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().AboveOrEqualTo(5).Above(3).Limit(0, 3).Get()
		if res[0] != "J" || res[1] != "D" || res[2] != "I" {
			t.Error("First 3 Above or Equal to 5 should be [J D I], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().Above(3).AboveOrEqualTo(5).Limit(0, 3).Get()
		if res[0] != "J" || res[1] != "D" || res[2] != "I" {
			t.Error("First 3 Above or Equal to 5 should be [J D I], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().BelowOrEqualTo(6).Below(8).Reversed().Limit(0, 3).Get()
		if res[0] != "D" || res[1] != "J" || res[2] != "C" {
			t.Error("First 3 below or equal to 6 should be [D J C], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().Below(8).BelowOrEqualTo(6).Reversed().Limit(0, 3).Get()
		if res[0] != "D" || res[1] != "J" || res[2] != "C" {
			t.Error("First 3 Below or equal to 6 should be [D J C], not", res)
		}
		done <- true
	}()

	for i := 0; i < 12; i++ {
		<-done
	}

	if res := <-ss.Scores().Above(4).Below(7).Remove(); res != 2 {
		t.Error("Should be removing 2 elements, not", res)
	}
	if res := <-ss.Scores().Get(); len(res) != 8 ||
		res[0] != "A" ||
		res[1] != "H" ||
		res[2] != "F" ||
		res[3] != "C" ||
		res[4] != "I" ||
		res[5] != "E" ||
		res[6] != "B" ||
		res[7] != "G" {
		t.Error("Should get [A H F C I E B G], not", res)
	}

	if res := <-ss.RemoveIndexedBetween(2, 5); res != 4 {
		t.Error("Should be removing 4 elements, not", res)
	}
	if res := <-ss.Scores().Get(); len(res) != 4 ||
		res[0] != "A" ||
		res[1] != "H" ||
		res[2] != "B" ||
		res[3] != "G" {
		t.Error("Should get [A H B G], not", res)
	}

	//Test Combos here
	otherss := r.SortedSet("Other_Test_SortedSet")

	otherss.Add("A", 5)
	otherss.Add("B", 6)
	otherss.Add("C", 3)
	otherss.Add("D", 4)

	go func() {
		resultset := r.SortedSet("InterSum_Test")
		if res := <-resultset.StoreIntersection().OfSet(ss).OfSet(otherss).UseCombinedScores(); res != 2 {
			t.Error("result set should have 2 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 2 || res["A"] != 6 || res["B"] != 15 {
			t.Error("Should end up with map[A:6 B:15], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedSet("InterMin_Test")
		if res := <-resultset.StoreIntersection().OfSet(ss).OfSet(otherss).UseLowerScore(); res != 2 {
			t.Error("result set should have 2 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 2 || res["A"] != 1 || res["B"] != 6 {
			t.Error("Should end up with map[A:1 B:6], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedSet("InterMax_Test")
		if res := <-resultset.StoreIntersection().OfSet(ss).OfSet(otherss).UseHigherScore(); res != 2 {
			t.Error("result set should have 2 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 2 || res["A"] != 5 || res["B"] != 9 {
			t.Error("Should end up with map[A:5 B:9], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedSet("InterWeight_Test")
		if res := <-resultset.StoreIntersection().OfSet(ss).OfWeightedSet(otherss, -1).UseCombinedScores(); res != 2 {
			t.Error("result set should have 2 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 2 || res["A"] != -4 || res["B"] != 3 {
			t.Error("Should end up with map[A:-4 B:3], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedSet("UnionSum_Test")
		if res := <-resultset.StoreUnion().OfSet(ss).OfSet(otherss).UseCombinedScores(); res != 6 {
			t.Error("result set should have 6 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 6 || res["A"] != 6 || res["B"] != 15 || res["C"] != 3 || res["D"] != 4 || res["G"] != 10 || res["H"] != 2 {
			t.Error("Should end up with map[A:6 B:15 C:3 D:4 G:10 H:2], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedSet("UnionMin_Test")
		if res := <-resultset.StoreUnion().OfSet(ss).OfSet(otherss).UseLowerScore(); res != 6 {
			t.Error("result set should have 6 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 6 || res["A"] != 1 || res["B"] != 6 || res["C"] != 3 || res["D"] != 4 || res["G"] != 10 || res["H"] != 2 {
			t.Error("Should end up with map[A:1 B:6 C:3 D:4 G:10 H:2], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedSet("UnionMax_Test")
		if res := <-resultset.StoreUnion().OfSet(ss).OfSet(otherss).UseHigherScore(); res != 6 {
			t.Error("result set should have 6 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 6 || res["A"] != 5 || res["B"] != 9 || res["C"] != 3 || res["D"] != 4 || res["G"] != 10 || res["H"] != 2 {
			t.Error("Should end up with map[A:5 B:9 C:3 D:4 G:10 H:2], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedSet("UnionWeight_Test")
		if res := <-resultset.StoreUnion().OfSet(ss).OfWeightedSet(otherss, -1).UseCombinedScores(); res != 6 {
			t.Error("result set should have 6 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 6 || res["A"] != -4 || res["B"] != 3 || res["C"] != -3 || res["D"] != -4 || res["G"] != 10 || res["H"] != 2 {
			t.Error("Should end up with map[A:6 B:15 C:-3 D:-4 G:10 H:2], not", res)
		}
		done <- true
	}()

	for i := 0; i < 8; i++ {
		<-done
	}

}
