package redis

import (
	"testing"
)

func TestSortedFloatSets(t *testing.T) {
	r := GetRedis(t)
	defer r.Close()

	ss := r.SortedFloatSet("Test_SortedFloatSet")
	ss.Delete()

	if _, ok := <-ss.ScoreOf(0.1); ok {
		t.Error("Should not be any score to get")
	}
	if _, ok := <-ss.IndexOf(0.1); ok {
		t.Error("Should not be any rank to get")
	}
	if res := <-ss.Size(); res != 0 {
		t.Error("Should have a size of 0, not", res)
	}

	a := ss.Add(0.1, 1)
	bb := ss.IncrementBy(0.2, 3)
	c := ss.Add(0.3, 2)

	<-a
	<-bb
	<-c

	if res := <-ss.ScoreOf(0.3); res != 2 {
		t.Error("0.3 should have a score of 2, not", res)
	}
	if res := <-ss.IndexOf(0.3); res != 1 {
		t.Error("0.3 should have an index of 1, not", res)
	}
	if res := <-ss.ReverseIndexOf(0.3); res != 1 {
		t.Error("0.3 should have a reverse index of 1, not", res)
	}

	if res := <-ss.Size(); res != 3 {
		t.Error("Should have a size of 3, not", res)
	}

	<-ss.IncrementBy(0.3, 2)

	if res := <-ss.ScoreOf(0.3); res != 4 {
		t.Error("0.3 should now have a score of 4, not", res)
	}
	if res := <-ss.IndexOf(0.3); res != 2 {
		t.Error("0.3 should now have an index of 2, not", res)
	}
	if res := <-ss.ReverseIndexOf(0.3); res != 0 {
		t.Error("0.3 should now have a reverse index of 0, not", res)
	}

	<-ss.Remove(0.2)

	if res := <-ss.ScoreOf(0.3); res != 4 {
		t.Error("Score of 0.3 should be unaffected; should still be 4, not", res)
	}
	if res := <-ss.IndexOf(0.3); res != 1 {
		t.Error("0.3 should now have an index of 1, not", res)
	}
	if res := <-ss.ReverseIndexOf(0.3); res != 0 {
		t.Error("0.3 should still have a reverse index of 0, not", res)
	}

	//	0.1 1
	a = ss.Add(0.2, 9)
	//	0.3 4
	b := ss.Add(0.4, 6)
	c = ss.Add(0.5, 8)
	d := ss.Add(0.6, 3)
	e := ss.Add(0.7, 10)
	f := ss.Add(0.8, 2)
	g := ss.Add(0.9, 7)
	h := ss.Add(1.0, 5)

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
			if res[0] != 0.7 {
				t.Error("Top score should be 0.7, not", res[0])
			}
			if res[1] != 0.2 {
				t.Error("2nd place should be 0.2, not", res[1])
			}
			if res[2] != 0.5 {
				t.Error("3rd place should be 0.5, not", res[2])
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
			if res[0] != 0.1 {
				t.Error("Bottom score should be 0.1, not", res[0])
			}
			if res[1] != 0.8 {
				t.Error("2nd worst should be 0.8, not", res[1])
			}
			if res[2] != 0.6 {
				t.Error("3rd worst should be 0.6, not", res[2])
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
		if len(get) != 4 || get[3] != 0.7 || get[2] != 0.2 || get[1] != 0.5 || get[0] != 0.9 {
			t.Error("Scores should be [0.9 0.5 0.2 0.7], not", get)
		}

		getwithscores := <-base.GetWithScores()
		if len(getwithscores) != 4 || getwithscores[0.7] != 10 || getwithscores[0.2] != 9 || getwithscores[0.5] != 8 || getwithscores[0.9] != 7 {
			t.Error("Scores should be map[0.7:10 0.2:9 0.5:8 0.9:7], not", getwithscores)
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
		if len(get) != 4 || get[3] != 0.1 || get[2] != 0.8 || get[1] != 0.6 || get[0] != 0.3 {
			t.Error("Scores should be [0.3 0.6 0.8 0.1], not", get)
		}

		getwithscores := <-base.GetWithScores()
		if len(getwithscores) != 4 || getwithscores[0.3] != 4 || getwithscores[0.6] != 3 || getwithscores[0.8] != 2 || getwithscores[0.1] != 1 {
			t.Error("Scores should be map[0.3:1 0.6:2 0.8:3 0.1:4], not", getwithscores)
		}

		done <- true
	}()

	//make sure Above/Below override AboveOrEqualTo/BelowOrEqualTo properly
	go func() {
		res := <-ss.Scores().AboveOrEqualTo(3).Above(3).Limit(0, 3).Get()
		if res[0] != 0.3 || res[1] != 1.0 || res[2] != 0.4 {
			t.Error("First 3 Above 3 should be [0.3 1.0 0.4], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().Above(3).AboveOrEqualTo(3).Limit(0, 3).GetWithScores()
		if res[0.3] != 4 || res[1.0] != 5 || res[0.4] != 6 {
			t.Error("First 3 Above 3 should be map[0.3:4 1.0:5 0.4:6], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().BelowOrEqualTo(8).Below(8).Reversed().Limit(0, 3).Get()
		if res[0] != 0.9 || res[1] != 0.4 || res[2] != 1.0 {
			t.Error("First 3 below 8 should be [0.9 0.4 1.0], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().Below(8).BelowOrEqualTo(8).Reversed().Limit(0, 3).GetWithScores()
		if res[0.9] != 7 || res[0.4] != 6 || res[1.0] != 5 {
			t.Error("First 3 Below 8 should be [0.9 0.4 1.0], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().AboveOrEqualTo(5).Above(3).Limit(0, 3).Get()
		if res[0] != 1.0 || res[1] != 0.4 || res[2] != 0.9 {
			t.Error("First 3 Above or Equal to 5 should be [1.0 0.4 0.9], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().Above(3).AboveOrEqualTo(5).Limit(0, 3).Get()
		if res[0] != 1.0 || res[1] != 0.4 || res[2] != 0.9 {
			t.Error("First 3 Above or Equal to 5 should be [1.0 0.4 0.9], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().BelowOrEqualTo(6).Below(8).Reversed().Limit(0, 3).Get()
		if res[0] != 0.4 || res[1] != 1.0 || res[2] != 0.3 {
			t.Error("First 3 below or equal to 6 should be [0.4 1.0 0.3], not", res)
		}
		done <- true
	}()

	go func() {
		res := <-ss.Scores().Below(8).BelowOrEqualTo(6).Reversed().Limit(0, 3).Get()
		if res[0] != 0.4 || res[1] != 1.0 || res[2] != 0.3 {
			t.Error("First 3 Below or equal to 6 should be [0.4 1.0 0.3], not", res)
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
		res[0] != 0.1 ||
		res[1] != 0.8 ||
		res[2] != 0.6 ||
		res[3] != 0.3 ||
		res[4] != 0.9 ||
		res[5] != 0.5 ||
		res[6] != 0.2 ||
		res[7] != 0.7 {
		t.Error("Should get [0.1 0.8 0.6 0.3 0.9 0.5 0.2 0.7], not", res)
	}

	if res := <-ss.RemoveIndexedBetween(2, 5); res != 4 {
		t.Error("Should be removing 4 elements, not", res)
	}
	if res := <-ss.Scores().Get(); len(res) != 4 ||
		res[0] != 0.1 ||
		res[1] != 0.8 ||
		res[2] != 0.2 ||
		res[3] != 0.7 {
		t.Error("Should get [0.1 0.8 0.2 0.7], not", res)
	}

	//Test Combos here
	otherss := r.SortedFloatSet("Other_Test_SortedFloatSet")

	otherss.Add(0.1, 5)
	otherss.Add(0.2, 6)
	otherss.Add(0.3, 3)
	otherss.Add(0.4, 4)

	go func() {
		resultset := r.SortedFloatSet("InterSum_FloatTest")
		if res := <-resultset.StoreIntersection().OfSet(ss).OfSet(otherss).UseCombinedScores(); res != 2 {
			t.Error("result set should have 2 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 2 || res[0.1] != 6 || res[0.2] != 15 {
			t.Error("Should end up with map[0.1:6 0.2:15], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedFloatSet("InterMin_FloatTest")
		if res := <-resultset.StoreIntersection().OfSet(ss).OfSet(otherss).UseLowerScore(); res != 2 {
			t.Error("result set should have 2 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 2 || res[0.1] != 1 || res[0.2] != 6 {
			t.Error("Should end up with map[0.1:1 0.2:6], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedFloatSet("InterMax_FloatTest")
		if res := <-resultset.StoreIntersection().OfSet(ss).OfSet(otherss).UseHigherScore(); res != 2 {
			t.Error("result set should have 2 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 2 || res[0.1] != 5 || res[0.2] != 9 {
			t.Error("Should end up with map[0.1:5 0.2:9], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedFloatSet("InterWeight_FloatTest")
		if res := <-resultset.StoreIntersection().OfSet(ss).OfWeightedSet(otherss, -1).UseCombinedScores(); res != 2 {
			t.Error("result set should have 2 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 2 || res[0.1] != -4 || res[0.2] != 3 {
			t.Error("Should end up with map[0.1:-4 0.2:3], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedFloatSet("UnionSum_FloatTest")
		if res := <-resultset.StoreUnion().OfSet(ss).OfSet(otherss).UseCombinedScores(); res != 6 {
			t.Error("result set should have 6 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 6 || res[0.1] != 6 || res[0.2] != 15 || res[0.3] != 3 || res[0.4] != 4 || res[0.7] != 10 || res[0.8] != 2 {
			t.Error("Should end up with map[0.1:6 0.2:15 0.3:3 0.4:4 0.7:10 0.8:2], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedFloatSet("UnionMin_FloatTest")
		if res := <-resultset.StoreUnion().OfSet(ss).OfSet(otherss).UseLowerScore(); res != 6 {
			t.Error("result set should have 6 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 6 || res[0.1] != 1 || res[0.2] != 6 || res[0.3] != 3 || res[0.4] != 4 || res[0.7] != 10 || res[0.8] != 2 {
			t.Error("Should end up with map[0.1:1 0.2:6 0.3:3 0.4:4 0.7:10 0.8:2], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedFloatSet("UnionMax_FloatTest")
		if res := <-resultset.StoreUnion().OfSet(ss).OfSet(otherss).UseHigherScore(); res != 6 {
			t.Error("result set should have 6 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 6 || res[0.1] != 5 || res[0.2] != 9 || res[0.3] != 3 || res[0.4] != 4 || res[0.7] != 10 || res[0.8] != 2 {
			t.Error("Should end up with map[0.1:5 0.2:9 0.3:3 0.4:4 0.7:10 0.8:2], not", res)
		}
		done <- true
	}()

	go func() {
		resultset := r.SortedFloatSet("UnionWeight_FloatTest")
		if res := <-resultset.StoreUnion().OfSet(ss).OfWeightedSet(otherss, -1).UseCombinedScores(); res != 6 {
			t.Error("result set should have 6 items, not", res)
		}
		if res := <-resultset.Scores().GetWithScores(); len(res) != 6 || res[0.1] != -4 || res[0.2] != 3 || res[0.3] != -3 || res[0.4] != -4 || res[0.7] != 10 || res[0.8] != 2 {
			t.Error("Should end up with map[0.1:6 0.2:15 0.3:-3 0.4:-4 0.7:10 0.8:2], not", res)
		}
		done <- true
	}()

	for i := 0; i < 8; i++ {
		<-done
	}
}
