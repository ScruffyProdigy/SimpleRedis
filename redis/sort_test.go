package redis

import (
	"testing"
)

func TestSorting(t *testing.T) {
	r := GetRedis(t)
	defer r.Close()

	//we have a lot to test here

	//need to test to make sure it works well with List's, Set's, and SortedSet's, and strings ints and floats for each of them
	//need to test alphabetically, numerically
	//need to test with and without limits
	//need to test with and without reverse
	//need to test with and without By
	//need to test with and without Get
	//need to test with and without storing

	//need to test each result type: regular returns, maybe returns, and storage returns for each of strings, ints, and floats

	//	backup
	str_to_str := func(index string) String {
		return r.String("string_" + index)
	}

	int_to_str := func(index int) String {
		return r.String("string_" + itoa(index))
	}

	flt_to_str := func(index float64) String {
		return r.String("string_" + ftoa(index))
	}

	str_to_int := func(index string) Integer {
		return r.Integer("integer_" + index)
	}

	int_to_int := func(index int) Integer {
		return r.Integer("integer_" + itoa(index))
	}

	flt_to_int := func(index float64) Integer {
		return r.Integer("integer_" + ftoa(index))
	}

	str_to_flt := func(index string) Float {
		return r.Float("float_" + index)
	}

	int_to_flt := func(index int) Float {
		return r.Float("float_" + itoa(index))
	}

	flt_to_flt := func(index float64) Float {
		return r.Float("float_" + ftoa(index))
	}

	// Setup
	str_to_str("A").Set("H")
	str_to_str("B").Set("J")
	str_to_str("C").Set("G")
	str_to_str("D").Set("I")
	str_to_str("E").Set("F")
	str_to_int("A").Set(10)
	str_to_int("B").Set(7)
	str_to_int("C").Set(9)
	str_to_int("D").Set(6)
	str_to_int("E").Set(8)
	str_to_flt("A").Set(0.9)
	str_to_flt("B").Set(0.7)
	str_to_flt("C").Set(1.0)
	str_to_flt("D").Set(0.8)
	str_to_flt("E").Set(0.6)
	int_to_str(1).Set("H")
	int_to_str(2).Set("J")
	int_to_str(3).Set("G")
	int_to_str(4).Set("I")
	int_to_str(5).Set("F")
	int_to_int(1).Set(10)
	int_to_int(2).Set(7)
	int_to_int(3).Set(9)
	int_to_int(4).Set(6)
	int_to_int(5).Set(8)
	int_to_flt(1).Set(0.9)
	int_to_flt(2).Set(0.7)
	int_to_flt(3).Set(1.0)
	int_to_flt(4).Set(0.8)
	int_to_flt(5).Set(0.6)
	flt_to_str(0.1).Set("H")
	flt_to_str(0.2).Set("J")
	flt_to_str(0.3).Set("G")
	flt_to_str(0.4).Set("I")
	flt_to_str(0.5).Set("F")
	flt_to_int(0.1).Set(10)
	flt_to_int(0.2).Set(7)
	flt_to_int(0.3).Set(9)
	flt_to_int(0.4).Set(6)
	flt_to_int(0.5).Set(8)
	flt_to_flt(0.1).Set(0.9)
	flt_to_flt(0.2).Set(0.7)
	flt_to_flt(0.3).Set(1.0)
	flt_to_flt(0.4).Set(0.8)
	flt_to_flt(0.5).Set(0.6)

	//	Storage
	StringStorage := r.List("Test_String_Storage")
	IntStorage := r.IntList("Test_Int_Storage")
	FloatStorage := r.FloatList("Test_Float_Storage")

	//	Lists

	str_list := r.List("Test_Sort_List")
	str_list.Delete()
	<-str_list.RightPush("C")
	<-str_list.RightPush("A")
	<-str_list.RightPush("D")
	<-str_list.RightPush("B")
	<-str_list.RightPush("E")

	if res := <-str_list.SortAlphabetically().Get(); len(res) != 5 || res[0] != "A" || res[1] != "B" || res[2] != "C" || res[3] != "D" || res[4] != "E" {
		t.Error("Should be [A B C D E], not", res)
	}

	if res := <-str_list.SortAlphabetically().Reverse().Get(); len(res) != 5 || res[0] != "E" || res[1] != "D" || res[2] != "C" || res[3] != "B" || res[4] != "A" {
		t.Error("Should be [E D C B A], not", res)
	}

	if res := <-str_list.SortAlphabetically().Limit(1, 3).Get(); len(res) != 3 || res[0] != "B" || res[1] != "C" || res[2] != "D" {
		t.Error("Should be [B C D], not", res)
	}

	if res := <-str_list.SortAlphabetically().Limit(0, 3).Reverse().Get(); len(res) != 3 || res[0] != "E" || res[1] != "D" || res[2] != "C" {
		t.Error("Should be [E D C], not", res)
	}

	if res := <-str_list.SortAlphabetically().By("string_*").Get(); len(res) != 5 || res[0] != "E" || res[1] != "C" || res[2] != "A" || res[3] != "D" || res[4] != "B" {
		t.Error("Should be [E C A D B], not", res)
	}

	if res := <-str_list.SortNumerically().By("integer_*").Get(); len(res) != 5 || res[0] != "D" || res[1] != "B" || res[2] != "E" || res[3] != "C" || res[4] != "A" {
		t.Error("Should be [D B E C A], not", res)
	}

	if res := <-str_list.SortNumerically().By("float_*").Get(); len(res) != 5 || res[0] != "E" || res[1] != "B" || res[2] != "D" || res[3] != "A" || res[4] != "C" {
		t.Error("Should be [E B D A C], not", res)
	}

	if res := <-str_list.SortAlphabetically().GetFrom("string_*"); len(res) != 5 || res[0] == nil || *res[0] != "H" || res[1] == nil || *res[1] != "J" || res[2] == nil || *res[2] != "G" || res[3] == nil || *res[3] != "I" || res[4] == nil || *res[4] != "F" {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = *sec
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [H J G I F], not [", result, "]")
	}

	if res := <-str_list.SortAlphabetically().GetIntsFrom("integer_*"); len(res) != 5 || res[0] == nil || *res[0] != 10 || res[1] == nil || *res[1] != 7 || res[2] == nil || *res[2] != 9 || res[3] == nil || *res[3] != 6 || res[4] == nil || *res[4] != 8 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = itoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [10 7 9 6 8], not [", result, "]")
	}

	if res := <-str_list.SortAlphabetically().GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.9 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 1.0 || res[3] == nil || *res[3] != 0.8 || res[4] == nil || *res[4] != 0.6 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.9 0.7 1.0 0.8 0.6], not [", result, "]")
	}

	if res := <-str_list.SortNumerically().By("integer_*").GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.8 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 0.6 || res[3] == nil || *res[3] != 1.0 || res[4] == nil || *res[4] != 0.9 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.8 0.7 0.6 1.0 0.8], not [", result, "]")
	}

	if res := <-str_list.SortAlphabetically().StoreStrings(StringStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-StringStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != "A" || end[1] != "B" || end[2] != "C" || end[3] != "D" || end[4] != "E" {
			t.Error("Result Should be [A B C D E], not", end)
		}
	}

	if res := <-str_list.SortAlphabetically().GetFromAndStoreIn("string_*", StringStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-StringStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != "H" || end[1] != "J" || end[2] != "G" || end[3] != "I" || end[4] != "F" {
			t.Error("Result Should be [E C A D B], not", end)
		}
	}

	if res := <-str_list.SortAlphabetically().GetIntsFromAndStoreIn("integer_*", IntStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-IntStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 10 || end[1] != 7 || end[2] != 9 || end[3] != 6 || end[4] != 8 {
			t.Error("Result Should be [10 7 9 6 8], not", end)
		}
	}

	if res := <-str_list.SortAlphabetically().GetFloatsFromAndStoreIn("float_*", FloatStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-FloatStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 0.9 || end[1] != 0.7 || end[2] != 1.0 || end[3] != 0.8 || end[4] != 0.6 {
			t.Error("Result Should be [0.9 0.7 1.0 0.8 0.6], not", end)
		}
	}

	<-str_list.RightPush("F")

	if res := <-str_list.SortAlphabetically().GetFrom("string_*"); len(res) != 6 || res[5] != nil {
		t.Error("New element should not be found in lookup")
	}

	int_list := r.IntList("Test_Sort_IntList")
	int_list.Delete()
	<-int_list.RightPush(3)
	<-int_list.RightPush(1)
	<-int_list.RightPush(4)
	<-int_list.RightPush(2)
	<-int_list.RightPush(5)

	if res := <-int_list.SortNumerically().GetInts(); len(res) != 5 || res[0] != 1 || res[1] != 2 || res[2] != 3 || res[3] != 4 || res[4] != 5 {
		t.Error("Should be [1 2 3 4 5], not", res)
	}

	if res := <-int_list.SortNumerically().Reverse().GetInts(); len(res) != 5 || res[0] != 5 || res[1] != 4 || res[2] != 3 || res[3] != 2 || res[4] != 1 {
		t.Error("Should be [E D C B A], not", res)
	}

	if res := <-int_list.SortNumerically().Limit(1, 3).GetInts(); len(res) != 3 || res[0] != 2 || res[1] != 3 || res[2] != 4 {
		t.Error("Should be [B C D], not", res)
	}

	if res := <-int_list.SortNumerically().Limit(0, 3).Reverse().GetInts(); len(res) != 3 || res[0] != 5 || res[1] != 4 || res[2] != 3 {
		t.Error("Should be [E D C], not", res)
	}

	if res := <-int_list.SortAlphabetically().By("string_*").GetInts(); len(res) != 5 || res[0] != 5 || res[1] != 3 || res[2] != 1 || res[3] != 4 || res[4] != 2 {
		t.Error("Should be [E C A D B], not", res)
	}

	if res := <-int_list.SortNumerically().By("integer_*").GetInts(); len(res) != 5 || res[0] != 4 || res[1] != 2 || res[2] != 5 || res[3] != 3 || res[4] != 1 {
		t.Error("Should be [D B E C A], not", res)
	}

	if res := <-int_list.SortNumerically().By("float_*").GetInts(); len(res) != 5 || res[0] != 5 || res[1] != 2 || res[2] != 4 || res[3] != 1 || res[4] != 3 {
		t.Error("Should be [E B D A C], not", res)
	}

	if res := <-int_list.SortNumerically().GetFrom("string_*"); len(res) != 5 || res[0] == nil || *res[0] != "H" || res[1] == nil || *res[1] != "J" || res[2] == nil || *res[2] != "G" || res[3] == nil || *res[3] != "I" || res[4] == nil || *res[4] != "F" {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = *sec
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [H J G I F], not [", result, "]")
	}

	if res := <-int_list.SortNumerically().GetIntsFrom("integer_*"); len(res) != 5 || res[0] == nil || *res[0] != 10 || res[1] == nil || *res[1] != 7 || res[2] == nil || *res[2] != 9 || res[3] == nil || *res[3] != 6 || res[4] == nil || *res[4] != 8 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = itoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [10 7 9 6 8], not [", result, "]")
	}

	if res := <-int_list.SortNumerically().GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.9 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 1.0 || res[3] == nil || *res[3] != 0.8 || res[4] == nil || *res[4] != 0.6 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.9 0.7 1.0 0.8 0.6], not [", result, "]")
	}

	if res := <-int_list.SortNumerically().By("integer_*").GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.8 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 0.6 || res[3] == nil || *res[3] != 1.0 || res[4] == nil || *res[4] != 0.9 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.8 0.7 0.6 1.0 0.8], not [", result, "]")
	}

	if res := <-int_list.SortNumerically().StoreInts(IntStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-IntStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 1 || end[1] != 2 || end[2] != 3 || end[3] != 4 || end[4] != 5 {
			t.Error("Result Should be [1 2 3 4 5], not", end)
		}
	}

	if res := <-int_list.SortNumerically().GetFromAndStoreIn("string_*", StringStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-StringStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != "H" || end[1] != "J" || end[2] != "G" || end[3] != "I" || end[4] != "F" {
			t.Error("Result Should be [E C A D B], not", end)
		}
	}

	if res := <-int_list.SortNumerically().GetIntsFromAndStoreIn("integer_*", IntStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-IntStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 10 || end[1] != 7 || end[2] != 9 || end[3] != 6 || end[4] != 8 {
			t.Error("Result Should be [10 7 9 6 8], not", end)
		}
	}

	if res := <-int_list.SortNumerically().GetFloatsFromAndStoreIn("float_*", FloatStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-FloatStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 0.9 || end[1] != 0.7 || end[2] != 1.0 || end[3] != 0.8 || end[4] != 0.6 {
			t.Error("Result Should be [0.9 0.7 1.0 0.8 0.6], not", end)
		}
	}

	<-int_list.RightPush(6)

	if res := <-int_list.SortNumerically().GetFrom("string_*"); len(res) != 6 || res[5] != nil {
		t.Error("New element should not be found in lookup")
	}

	float_list := r.FloatList("Test_Sort_FloatList")
	float_list.Delete()
	<-float_list.RightPush(0.3)
	<-float_list.RightPush(0.1)
	<-float_list.RightPush(0.4)
	<-float_list.RightPush(0.2)
	<-float_list.RightPush(0.5)

	if res := <-float_list.SortNumerically().GetFloats(); len(res) != 5 || res[0] != 0.1 || res[1] != 0.2 || res[2] != 0.3 || res[3] != 0.4 || res[4] != 0.5 {
		t.Error("Should be [0.1 0.2 0.3 0.4 0.5], not", res)
	}

	if res := <-float_list.SortNumerically().Reverse().GetFloats(); len(res) != 5 || res[0] != 0.5 || res[1] != 0.4 || res[2] != 0.3 || res[3] != 0.2 || res[4] != 0.1 {
		t.Error("Should be [E D C B A], not", res)
	}

	if res := <-float_list.SortNumerically().Limit(1, 3).GetFloats(); len(res) != 3 || res[0] != 0.2 || res[1] != 0.3 || res[2] != 0.4 {
		t.Error("Should be [B C D], not", res)
	}

	if res := <-float_list.SortNumerically().Limit(0, 3).Reverse().GetFloats(); len(res) != 3 || res[0] != 0.5 || res[1] != 0.4 || res[2] != 0.3 {
		t.Error("Should be [E D C], not", res)
	}

	if res := <-float_list.SortAlphabetically().By("string_*").GetFloats(); len(res) != 5 || res[0] != 0.5 || res[1] != 0.3 || res[2] != 0.1 || res[3] != 0.4 || res[4] != 0.2 {
		t.Error("Should be [E C A D B], not", res)
	}

	if res := <-float_list.SortNumerically().By("integer_*").GetFloats(); len(res) != 5 || res[0] != 0.4 || res[1] != 0.2 || res[2] != 0.5 || res[3] != 0.3 || res[4] != 0.1 {
		t.Error("Should be [D B E C A], not", res)
	}

	if res := <-float_list.SortNumerically().By("float_*").GetFloats(); len(res) != 5 || res[0] != 0.5 || res[1] != 0.2 || res[2] != 0.4 || res[3] != 0.1 || res[4] != 0.3 {
		t.Error("Should be [E B D A C], not", res)
	}

	if res := <-float_list.SortNumerically().GetFrom("string_*"); len(res) != 5 || res[0] == nil || *res[0] != "H" || res[1] == nil || *res[1] != "J" || res[2] == nil || *res[2] != "G" || res[3] == nil || *res[3] != "I" || res[4] == nil || *res[4] != "F" {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = *sec
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [H J G I F], not [", result, "]")
	}

	if res := <-float_list.SortNumerically().GetIntsFrom("integer_*"); len(res) != 5 || res[0] == nil || *res[0] != 10 || res[1] == nil || *res[1] != 7 || res[2] == nil || *res[2] != 9 || res[3] == nil || *res[3] != 6 || res[4] == nil || *res[4] != 8 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = itoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [10 7 9 6 8], not [", result, "]")
	}

	if res := <-float_list.SortNumerically().GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.9 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 1.0 || res[3] == nil || *res[3] != 0.8 || res[4] == nil || *res[4] != 0.6 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.9 0.7 1.0 0.8 0.6], not [", result, "]")
	}

	if res := <-float_list.SortNumerically().By("integer_*").GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.8 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 0.6 || res[3] == nil || *res[3] != 1.0 || res[4] == nil || *res[4] != 0.9 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.8 0.7 0.6 1.0 0.8], not [", result, "]")
	}

	if res := <-float_list.SortNumerically().StoreFloats(FloatStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-FloatStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 0.1 || end[1] != 0.2 || end[2] != 0.3 || end[3] != 0.4 || end[4] != 0.5 {
			t.Error("Result Should be [0.1 0.2 0.3 0.4 0.5], not", end)
		}
	}

	if res := <-float_list.SortNumerically().GetFromAndStoreIn("string_*", StringStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-StringStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != "H" || end[1] != "J" || end[2] != "G" || end[3] != "I" || end[4] != "F" {
			t.Error("Result Should be [E C A D B], not", end)
		}
	}

	if res := <-float_list.SortNumerically().GetIntsFromAndStoreIn("integer_*", IntStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-IntStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 10 || end[1] != 7 || end[2] != 9 || end[3] != 6 || end[4] != 8 {
			t.Error("Result Should be [10 7 9 6 8], not", end)
		}
	}

	if res := <-float_list.SortNumerically().GetFloatsFromAndStoreIn("float_*", FloatStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-FloatStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 0.9 || end[1] != 0.7 || end[2] != 1.0 || end[3] != 0.8 || end[4] != 0.6 {
			t.Error("Result Should be [0.9 0.7 1.0 0.8 0.6], not", end)
		}
	}

	<-float_list.RightPush(0.25)

	if res := <-float_list.SortNumerically().GetFrom("string_*"); len(res) != 6 || res[2] != nil {
		t.Error("New element should not be found in lookup")
	}

	//	Sets

	str_set := r.Set("Test_Sort_Set")
	str_set.Delete()
	<-str_set.Add("C")
	<-str_set.Add("A")
	<-str_set.Add("D")
	<-str_set.Add("B")
	<-str_set.Add("E")

	if res := <-str_set.SortAlphabetically().Get(); len(res) != 5 || res[0] != "A" || res[1] != "B" || res[2] != "C" || res[3] != "D" || res[4] != "E" {
		t.Error("Should be [A B C D E], not", res)
	}

	if res := <-str_set.SortAlphabetically().Reverse().Get(); len(res) != 5 || res[0] != "E" || res[1] != "D" || res[2] != "C" || res[3] != "B" || res[4] != "A" {
		t.Error("Should be [E D C B A], not", res)
	}

	if res := <-str_set.SortAlphabetically().Limit(1, 3).Get(); len(res) != 3 || res[0] != "B" || res[1] != "C" || res[2] != "D" {
		t.Error("Should be [B C D], not", res)
	}

	if res := <-str_set.SortAlphabetically().Limit(0, 3).Reverse().Get(); len(res) != 3 || res[0] != "E" || res[1] != "D" || res[2] != "C" {
		t.Error("Should be [E D C], not", res)
	}

	if res := <-str_set.SortAlphabetically().By("string_*").Get(); len(res) != 5 || res[0] != "E" || res[1] != "C" || res[2] != "A" || res[3] != "D" || res[4] != "B" {
		t.Error("Should be [E C A D B], not", res)
	}

	if res := <-str_set.SortNumerically().By("integer_*").Get(); len(res) != 5 || res[0] != "D" || res[1] != "B" || res[2] != "E" || res[3] != "C" || res[4] != "A" {
		t.Error("Should be [D B E C A], not", res)
	}

	if res := <-str_set.SortNumerically().By("float_*").Get(); len(res) != 5 || res[0] != "E" || res[1] != "B" || res[2] != "D" || res[3] != "A" || res[4] != "C" {
		t.Error("Should be [E B D A C], not", res)
	}

	if res := <-str_set.SortAlphabetically().GetFrom("string_*"); len(res) != 5 || res[0] == nil || *res[0] != "H" || res[1] == nil || *res[1] != "J" || res[2] == nil || *res[2] != "G" || res[3] == nil || *res[3] != "I" || res[4] == nil || *res[4] != "F" {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = *sec
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [H J G I F], not [", result, "]")
	}

	if res := <-str_set.SortAlphabetically().GetIntsFrom("integer_*"); len(res) != 5 || res[0] == nil || *res[0] != 10 || res[1] == nil || *res[1] != 7 || res[2] == nil || *res[2] != 9 || res[3] == nil || *res[3] != 6 || res[4] == nil || *res[4] != 8 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = itoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [10 7 9 6 8], not [", result, "]")
	}

	if res := <-str_set.SortAlphabetically().GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.9 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 1.0 || res[3] == nil || *res[3] != 0.8 || res[4] == nil || *res[4] != 0.6 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.9 0.7 1.0 0.8 0.6], not [", result, "]")
	}

	if res := <-str_set.SortNumerically().By("integer_*").GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.8 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 0.6 || res[3] == nil || *res[3] != 1.0 || res[4] == nil || *res[4] != 0.9 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.8 0.7 0.6 1.0 0.9], not [", result, "]")
	}

	if res := <-str_set.SortAlphabetically().StoreStrings(StringStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-StringStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != "A" || end[1] != "B" || end[2] != "C" || end[3] != "D" || end[4] != "E" {
			t.Error("Result Should be [A B C D E], not", end)
		}
	}

	if res := <-str_set.SortAlphabetically().GetFromAndStoreIn("string_*", StringStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-StringStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != "H" || end[1] != "J" || end[2] != "G" || end[3] != "I" || end[4] != "F" {
			t.Error("Result Should be [E C A D B], not", end)
		}
	}

	if res := <-str_set.SortAlphabetically().GetIntsFromAndStoreIn("integer_*", IntStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-IntStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 10 || end[1] != 7 || end[2] != 9 || end[3] != 6 || end[4] != 8 {
			t.Error("Result Should be [10 7 9 6 8], not", end)
		}
	}

	if res := <-str_set.SortAlphabetically().GetFloatsFromAndStoreIn("float_*", FloatStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-FloatStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 0.9 || end[1] != 0.7 || end[2] != 1.0 || end[3] != 0.8 || end[4] != 0.6 {
			t.Error("Result Should be [0.9 0.7 1.0 0.8 0.6], not", end)
		}
	}

	<-str_set.Add("F")

	if res := <-str_set.SortAlphabetically().GetFrom("string_*"); len(res) != 6 || res[5] != nil {
		t.Error("New element should not be found in lookup")
	}

	int_set := r.IntSet("Test_Sort_IntSet")
	int_set.Delete()
	<-int_set.Add(3)
	<-int_set.Add(1)
	<-int_set.Add(4)
	<-int_set.Add(2)
	<-int_set.Add(5)

	if res := <-int_set.SortNumerically().GetInts(); len(res) != 5 || res[0] != 1 || res[1] != 2 || res[2] != 3 || res[3] != 4 || res[4] != 5 {
		t.Error("Should be [1 2 3 4 5], not", res)
	}

	if res := <-int_set.SortNumerically().Reverse().GetInts(); len(res) != 5 || res[0] != 5 || res[1] != 4 || res[2] != 3 || res[3] != 2 || res[4] != 1 {
		t.Error("Should be [E D C B A], not", res)
	}

	if res := <-int_set.SortNumerically().Limit(1, 3).GetInts(); len(res) != 3 || res[0] != 2 || res[1] != 3 || res[2] != 4 {
		t.Error("Should be [B C D], not", res)
	}

	if res := <-int_set.SortNumerically().Limit(0, 3).Reverse().GetInts(); len(res) != 3 || res[0] != 5 || res[1] != 4 || res[2] != 3 {
		t.Error("Should be [E D C], not", res)
	}

	if res := <-int_set.SortAlphabetically().By("string_*").GetInts(); len(res) != 5 || res[0] != 5 || res[1] != 3 || res[2] != 1 || res[3] != 4 || res[4] != 2 {
		t.Error("Should be [E C A D B], not", res)
	}

	if res := <-int_set.SortNumerically().By("integer_*").GetInts(); len(res) != 5 || res[0] != 4 || res[1] != 2 || res[2] != 5 || res[3] != 3 || res[4] != 1 {
		t.Error("Should be [D B E C A], not", res)
	}

	if res := <-int_set.SortNumerically().By("float_*").GetInts(); len(res) != 5 || res[0] != 5 || res[1] != 2 || res[2] != 4 || res[3] != 1 || res[4] != 3 {
		t.Error("Should be [E B D A C], not", res)
	}

	if res := <-int_set.SortNumerically().GetFrom("string_*"); len(res) != 5 || res[0] == nil || *res[0] != "H" || res[1] == nil || *res[1] != "J" || res[2] == nil || *res[2] != "G" || res[3] == nil || *res[3] != "I" || res[4] == nil || *res[4] != "F" {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = *sec
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [H J G I F], not [", result, "]")
	}

	if res := <-int_set.SortNumerically().GetIntsFrom("integer_*"); len(res) != 5 || res[0] == nil || *res[0] != 10 || res[1] == nil || *res[1] != 7 || res[2] == nil || *res[2] != 9 || res[3] == nil || *res[3] != 6 || res[4] == nil || *res[4] != 8 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = itoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [10 7 9 6 8], not [", result, "]")
	}

	if res := <-int_set.SortNumerically().GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.9 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 1.0 || res[3] == nil || *res[3] != 0.8 || res[4] == nil || *res[4] != 0.6 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.9 0.7 1.0 0.8 0.6], not [", result, "]")
	}

	if res := <-int_set.SortNumerically().By("integer_*").GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.8 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 0.6 || res[3] == nil || *res[3] != 1.0 || res[4] == nil || *res[4] != 0.9 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.8 0.7 0.6 1.0 0.8], not [", result, "]")
	}

	if res := <-int_set.SortNumerically().StoreInts(IntStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-IntStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 1 || end[1] != 2 || end[2] != 3 || end[3] != 4 || end[4] != 5 {
			t.Error("Result Should be [1 2 3 4 5], not", end)
		}
	}

	if res := <-int_set.SortNumerically().GetFromAndStoreIn("string_*", StringStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-StringStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != "H" || end[1] != "J" || end[2] != "G" || end[3] != "I" || end[4] != "F" {
			t.Error("Result Should be [E C A D B], not", end)
		}
	}

	if res := <-int_set.SortNumerically().GetIntsFromAndStoreIn("integer_*", IntStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-IntStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 10 || end[1] != 7 || end[2] != 9 || end[3] != 6 || end[4] != 8 {
			t.Error("Result Should be [10 7 9 6 8], not", end)
		}
	}

	if res := <-int_set.SortNumerically().GetFloatsFromAndStoreIn("float_*", FloatStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-FloatStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 0.9 || end[1] != 0.7 || end[2] != 1.0 || end[3] != 0.8 || end[4] != 0.6 {
			t.Error("Result Should be [0.9 0.7 1.0 0.8 0.6], not", end)
		}
	}

	<-int_set.Add(6)

	if res := <-int_set.SortNumerically().GetFrom("string_*"); len(res) != 6 || res[5] != nil {
		t.Error("New element should not be found in lookup")
	}

	float_set := r.FloatSet("Test_Sort_FloatSet")
	float_set.Delete()
	<-float_set.Add(0.3)
	<-float_set.Add(0.1)
	<-float_set.Add(0.4)
	<-float_set.Add(0.2)
	<-float_set.Add(0.5)

	if res := <-float_set.SortNumerically().GetFloats(); len(res) != 5 || res[0] != 0.1 || res[1] != 0.2 || res[2] != 0.3 || res[3] != 0.4 || res[4] != 0.5 {
		t.Error("Should be [1 2 3 4 5], not", res)
	}

	if res := <-float_set.SortNumerically().GetFloats(); len(res) != 5 || res[0] != 0.1 || res[1] != 0.2 || res[2] != 0.3 || res[3] != 0.4 || res[4] != 0.5 {
		t.Error("Should be [0.1 0.2 0.3 0.4 0.5], not", res)
	}

	if res := <-float_set.SortNumerically().Reverse().GetFloats(); len(res) != 5 || res[0] != 0.5 || res[1] != 0.4 || res[2] != 0.3 || res[3] != 0.2 || res[4] != 0.1 {
		t.Error("Should be [E D C B A], not", res)
	}

	if res := <-float_set.SortNumerically().Limit(1, 3).GetFloats(); len(res) != 3 || res[0] != 0.2 || res[1] != 0.3 || res[2] != 0.4 {
		t.Error("Should be [B C D], not", res)
	}

	if res := <-float_set.SortNumerically().Limit(0, 3).Reverse().GetFloats(); len(res) != 3 || res[0] != 0.5 || res[1] != 0.4 || res[2] != 0.3 {
		t.Error("Should be [E D C], not", res)
	}

	if res := <-float_set.SortAlphabetically().By("string_*").GetFloats(); len(res) != 5 || res[0] != 0.5 || res[1] != 0.3 || res[2] != 0.1 || res[3] != 0.4 || res[4] != 0.2 {
		t.Error("Should be [E C A D B], not", res)
	}

	if res := <-float_set.SortNumerically().By("integer_*").GetFloats(); len(res) != 5 || res[0] != 0.4 || res[1] != 0.2 || res[2] != 0.5 || res[3] != 0.3 || res[4] != 0.1 {
		t.Error("Should be [D B E C A], not", res)
	}

	if res := <-float_set.SortNumerically().By("float_*").GetFloats(); len(res) != 5 || res[0] != 0.5 || res[1] != 0.2 || res[2] != 0.4 || res[3] != 0.1 || res[4] != 0.3 {
		t.Error("Should be [E B D A C], not", res)
	}

	if res := <-float_set.SortNumerically().GetFrom("string_*"); len(res) != 5 || res[0] == nil || *res[0] != "H" || res[1] == nil || *res[1] != "J" || res[2] == nil || *res[2] != "G" || res[3] == nil || *res[3] != "I" || res[4] == nil || *res[4] != "F" {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = *sec
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [H J G I F], not [", result, "]")
	}

	if res := <-float_set.SortNumerically().GetIntsFrom("integer_*"); len(res) != 5 || res[0] == nil || *res[0] != 10 || res[1] == nil || *res[1] != 7 || res[2] == nil || *res[2] != 9 || res[3] == nil || *res[3] != 6 || res[4] == nil || *res[4] != 8 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = itoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [10 7 9 6 8], not [", result, "]")
	}

	if res := <-float_set.SortNumerically().GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.9 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 1.0 || res[3] == nil || *res[3] != 0.8 || res[4] == nil || *res[4] != 0.6 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.9 0.7 1.0 0.8 0.6], not [", result, "]")
	}

	if res := <-float_set.SortNumerically().By("integer_*").GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.8 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 0.6 || res[3] == nil || *res[3] != 1.0 || res[4] == nil || *res[4] != 0.9 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.8 0.7 0.6 1.0 0.8], not [", result, "]")
	}

	if res := <-float_set.SortNumerically().StoreFloats(FloatStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-FloatStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 0.1 || end[1] != 0.2 || end[2] != 0.3 || end[3] != 0.4 || end[4] != 0.5 {
			t.Error("Result Should be [0.1 0.2 0.3 0.4 0.5], not", end)
		}
	}

	if res := <-float_set.SortNumerically().GetFromAndStoreIn("string_*", StringStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-StringStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != "H" || end[1] != "J" || end[2] != "G" || end[3] != "I" || end[4] != "F" {
			t.Error("Result Should be [E C A D B], not", end)
		}
	}

	if res := <-float_set.SortNumerically().GetIntsFromAndStoreIn("integer_*", IntStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-IntStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 10 || end[1] != 7 || end[2] != 9 || end[3] != 6 || end[4] != 8 {
			t.Error("Result Should be [10 7 9 6 8], not", end)
		}
	}

	if res := <-float_set.SortNumerically().GetFloatsFromAndStoreIn("float_*", FloatStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-FloatStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 0.9 || end[1] != 0.7 || end[2] != 1.0 || end[3] != 0.8 || end[4] != 0.6 {
			t.Error("Result Should be [0.9 0.7 1.0 0.8 0.6], not", end)
		}
	}

	<-float_set.Add(0.25)

	if res := <-float_set.SortNumerically().GetFrom("string_*"); len(res) != 6 || res[2] != nil {
		t.Error("New element should not be found in lookup")
	}

	//	SortedSets

	str_ss := r.SortedSet("Test_Sort_SortedSet")
	str_ss.Delete()
	<-str_ss.Add("C", 1)
	<-str_ss.Add("A", 2)
	<-str_ss.Add("D", 3)
	<-str_ss.Add("B", 4)
	<-str_ss.Add("E", 5)

	if res := <-str_ss.SortAlphabetically().Get(); len(res) != 5 || res[0] != "A" || res[1] != "B" || res[2] != "C" || res[3] != "D" || res[4] != "E" {
		t.Error("Should be [A B C D E], not", res)
	}

	if res := <-str_ss.SortAlphabetically().Reverse().Get(); len(res) != 5 || res[0] != "E" || res[1] != "D" || res[2] != "C" || res[3] != "B" || res[4] != "A" {
		t.Error("Should be [E D C B A], not", res)
	}

	if res := <-str_ss.SortAlphabetically().Limit(1, 3).Get(); len(res) != 3 || res[0] != "B" || res[1] != "C" || res[2] != "D" {
		t.Error("Should be [B C D], not", res)
	}

	if res := <-str_ss.SortAlphabetically().Limit(0, 3).Reverse().Get(); len(res) != 3 || res[0] != "E" || res[1] != "D" || res[2] != "C" {
		t.Error("Should be [E D C], not", res)
	}

	if res := <-str_ss.SortAlphabetically().By("string_*").Get(); len(res) != 5 || res[0] != "E" || res[1] != "C" || res[2] != "A" || res[3] != "D" || res[4] != "B" {
		t.Error("Should be [E C A D B], not", res)
	}

	if res := <-str_ss.SortNumerically().By("integer_*").Get(); len(res) != 5 || res[0] != "D" || res[1] != "B" || res[2] != "E" || res[3] != "C" || res[4] != "A" {
		t.Error("Should be [D B E C A], not", res)
	}

	if res := <-str_ss.SortNumerically().By("float_*").Get(); len(res) != 5 || res[0] != "E" || res[1] != "B" || res[2] != "D" || res[3] != "A" || res[4] != "C" {
		t.Error("Should be [E B D A C], not", res)
	}

	if res := <-str_ss.SortAlphabetically().GetFrom("string_*"); len(res) != 5 || res[0] == nil || *res[0] != "H" || res[1] == nil || *res[1] != "J" || res[2] == nil || *res[2] != "G" || res[3] == nil || *res[3] != "I" || res[4] == nil || *res[4] != "F" {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = *sec
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [H J G I F], not [", result, "]")
	}

	if res := <-str_ss.SortAlphabetically().GetIntsFrom("integer_*"); len(res) != 5 || res[0] == nil || *res[0] != 10 || res[1] == nil || *res[1] != 7 || res[2] == nil || *res[2] != 9 || res[3] == nil || *res[3] != 6 || res[4] == nil || *res[4] != 8 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = itoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [10 7 9 6 8], not [", result, "]")
	}

	if res := <-str_ss.SortAlphabetically().GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.9 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 1.0 || res[3] == nil || *res[3] != 0.8 || res[4] == nil || *res[4] != 0.6 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.9 0.7 1.0 0.8 0.6], not [", result, "]")
	}

	if res := <-str_ss.SortNumerically().By("integer_*").GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.8 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 0.6 || res[3] == nil || *res[3] != 1.0 || res[4] == nil || *res[4] != 0.9 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.8 0.7 0.6 1.0 0.8], not [", result, "]")
	}

	if res := <-str_ss.SortAlphabetically().StoreStrings(StringStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-StringStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != "A" || end[1] != "B" || end[2] != "C" || end[3] != "D" || end[4] != "E" {
			t.Error("Result Should be [A B C D E], not", end)
		}
	}

	if res := <-str_ss.SortAlphabetically().GetFromAndStoreIn("string_*", StringStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-StringStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != "H" || end[1] != "J" || end[2] != "G" || end[3] != "I" || end[4] != "F" {
			t.Error("Result Should be [E C A D B], not", end)
		}
	}

	if res := <-str_ss.SortAlphabetically().GetIntsFromAndStoreIn("integer_*", IntStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-IntStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 10 || end[1] != 7 || end[2] != 9 || end[3] != 6 || end[4] != 8 {
			t.Error("Result Should be [10 7 9 6 8], not", end)
		}
	}

	if res := <-str_ss.SortAlphabetically().GetFloatsFromAndStoreIn("float_*", FloatStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-FloatStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 0.9 || end[1] != 0.7 || end[2] != 1.0 || end[3] != 0.8 || end[4] != 0.6 {
			t.Error("Result Should be [0.9 0.7 1.0 0.8 0.6], not", end)
		}
	}

	<-str_ss.Add("F", 2.5)

	if res := <-str_ss.SortAlphabetically().GetFrom("string_*"); len(res) != 6 || res[5] != nil {
		t.Error("New element should not be found in lookup")
	}

	int_ss := r.SortedIntSet("Test_Sort_SortedIntSet")
	int_ss.Delete()
	<-int_ss.Add(3, 1)
	<-int_ss.Add(1, 2)
	<-int_ss.Add(4, 3)
	<-int_ss.Add(2, 4)
	<-int_ss.Add(5, 5)

	if res := <-int_ss.SortNumerically().GetInts(); len(res) != 5 || res[0] != 1 || res[1] != 2 || res[2] != 3 || res[3] != 4 || res[4] != 5 {
		t.Error("Should be [1 2 3 4 5], not", res)
	}

	if res := <-int_ss.SortNumerically().Reverse().GetInts(); len(res) != 5 || res[0] != 5 || res[1] != 4 || res[2] != 3 || res[3] != 2 || res[4] != 1 {
		t.Error("Should be [E D C B A], not", res)
	}

	if res := <-int_ss.SortNumerically().Limit(1, 3).GetInts(); len(res) != 3 || res[0] != 2 || res[1] != 3 || res[2] != 4 {
		t.Error("Should be [B C D], not", res)
	}

	if res := <-int_ss.SortNumerically().Limit(0, 3).Reverse().GetInts(); len(res) != 3 || res[0] != 5 || res[1] != 4 || res[2] != 3 {
		t.Error("Should be [E D C], not", res)
	}

	if res := <-int_ss.SortAlphabetically().By("string_*").GetInts(); len(res) != 5 || res[0] != 5 || res[1] != 3 || res[2] != 1 || res[3] != 4 || res[4] != 2 {
		t.Error("Should be [E C A D B], not", res)
	}

	if res := <-int_ss.SortNumerically().By("integer_*").GetInts(); len(res) != 5 || res[0] != 4 || res[1] != 2 || res[2] != 5 || res[3] != 3 || res[4] != 1 {
		t.Error("Should be [D B E C A], not", res)
	}

	if res := <-int_ss.SortNumerically().By("float_*").GetInts(); len(res) != 5 || res[0] != 5 || res[1] != 2 || res[2] != 4 || res[3] != 1 || res[4] != 3 {
		t.Error("Should be [E B D A C], not", res)
	}

	if res := <-int_ss.SortNumerically().GetFrom("string_*"); len(res) != 5 || res[0] == nil || *res[0] != "H" || res[1] == nil || *res[1] != "J" || res[2] == nil || *res[2] != "G" || res[3] == nil || *res[3] != "I" || res[4] == nil || *res[4] != "F" {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = *sec
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [H J G I F], not [", result, "]")
	}

	if res := <-int_ss.SortNumerically().GetIntsFrom("integer_*"); len(res) != 5 || res[0] == nil || *res[0] != 10 || res[1] == nil || *res[1] != 7 || res[2] == nil || *res[2] != 9 || res[3] == nil || *res[3] != 6 || res[4] == nil || *res[4] != 8 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = itoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [10 7 9 6 8], not [", result, "]")
	}

	if res := <-int_ss.SortNumerically().GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.9 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 1.0 || res[3] == nil || *res[3] != 0.8 || res[4] == nil || *res[4] != 0.6 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.9 0.7 1.0 0.8 0.6], not [", result, "]")
	}

	if res := <-int_ss.SortNumerically().By("integer_*").GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.8 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 0.6 || res[3] == nil || *res[3] != 1.0 || res[4] == nil || *res[4] != 0.9 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.8 0.7 0.6 1.0 0.8], not [", result, "]")
	}

	if res := <-int_ss.SortNumerically().StoreInts(IntStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-IntStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 1 || end[1] != 2 || end[2] != 3 || end[3] != 4 || end[4] != 5 {
			t.Error("Result Should be [1 2 3 4 5], not", end)
		}
	}

	if res := <-int_ss.SortNumerically().GetFromAndStoreIn("string_*", StringStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-StringStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != "H" || end[1] != "J" || end[2] != "G" || end[3] != "I" || end[4] != "F" {
			t.Error("Result Should be [E C A D B], not", end)
		}
	}

	if res := <-int_ss.SortNumerically().GetIntsFromAndStoreIn("integer_*", IntStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-IntStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 10 || end[1] != 7 || end[2] != 9 || end[3] != 6 || end[4] != 8 {
			t.Error("Result Should be [10 7 9 6 8], not", end)
		}
	}

	if res := <-int_ss.SortNumerically().GetFloatsFromAndStoreIn("float_*", FloatStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-FloatStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 0.9 || end[1] != 0.7 || end[2] != 1.0 || end[3] != 0.8 || end[4] != 0.6 {
			t.Error("Result Should be [0.9 0.7 1.0 0.8 0.6], not", end)
		}
	}

	<-int_ss.Add(6, 2.5)

	if res := <-int_ss.SortNumerically().GetFrom("string_*"); len(res) != 6 || res[5] != nil {
		t.Error("New element should not be found in lookup")
	}

	float_ss := r.SortedFloatSet("Test_Sort_SortedFloatSet")
	float_ss.Delete()
	<-float_ss.Add(0.3, 1)
	<-float_ss.Add(0.1, 2)
	<-float_ss.Add(0.4, 3)
	<-float_ss.Add(0.2, 4)
	<-float_ss.Add(0.5, 5)

	if res := <-float_ss.SortNumerically().GetFloats(); len(res) != 5 || res[0] != 0.1 || res[1] != 0.2 || res[2] != 0.3 || res[3] != 0.4 || res[4] != 0.5 {
		t.Error("Should be [1 2 3 4 5], not", res)
	}

	if res := <-float_ss.SortNumerically().GetFloats(); len(res) != 5 || res[0] != 0.1 || res[1] != 0.2 || res[2] != 0.3 || res[3] != 0.4 || res[4] != 0.5 {
		t.Error("Should be [1 2 3 4 5], not", res)
	}

	if res := <-float_ss.SortNumerically().GetFloats(); len(res) != 5 || res[0] != 0.1 || res[1] != 0.2 || res[2] != 0.3 || res[3] != 0.4 || res[4] != 0.5 {
		t.Error("Should be [0.1 0.2 0.3 0.4 0.5], not", res)
	}

	if res := <-float_ss.SortNumerically().Reverse().GetFloats(); len(res) != 5 || res[0] != 0.5 || res[1] != 0.4 || res[2] != 0.3 || res[3] != 0.2 || res[4] != 0.1 {
		t.Error("Should be [E D C B A], not", res)
	}

	if res := <-float_ss.SortNumerically().Limit(1, 3).GetFloats(); len(res) != 3 || res[0] != 0.2 || res[1] != 0.3 || res[2] != 0.4 {
		t.Error("Should be [B C D], not", res)
	}

	if res := <-float_ss.SortNumerically().Limit(0, 3).Reverse().GetFloats(); len(res) != 3 || res[0] != 0.5 || res[1] != 0.4 || res[2] != 0.3 {
		t.Error("Should be [E D C], not", res)
	}

	if res := <-float_ss.SortAlphabetically().By("string_*").GetFloats(); len(res) != 5 || res[0] != 0.5 || res[1] != 0.3 || res[2] != 0.1 || res[3] != 0.4 || res[4] != 0.2 {
		t.Error("Should be [E C A D B], not", res)
	}

	if res := <-float_ss.SortNumerically().By("integer_*").GetFloats(); len(res) != 5 || res[0] != 0.4 || res[1] != 0.2 || res[2] != 0.5 || res[3] != 0.3 || res[4] != 0.1 {
		t.Error("Should be [D B E C A], not", res)
	}

	if res := <-float_ss.SortNumerically().By("float_*").GetFloats(); len(res) != 5 || res[0] != 0.5 || res[1] != 0.2 || res[2] != 0.4 || res[3] != 0.1 || res[4] != 0.3 {
		t.Error("Should be [E B D A C], not", res)
	}

	if res := <-float_ss.SortNumerically().GetFrom("string_*"); len(res) != 5 || res[0] == nil || *res[0] != "H" || res[1] == nil || *res[1] != "J" || res[2] == nil || *res[2] != "G" || res[3] == nil || *res[3] != "I" || res[4] == nil || *res[4] != "F" {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = *sec
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [H J G I F], not [", result, "]")
	}

	if res := <-float_ss.SortNumerically().GetIntsFrom("integer_*"); len(res) != 5 || res[0] == nil || *res[0] != 10 || res[1] == nil || *res[1] != 7 || res[2] == nil || *res[2] != 9 || res[3] == nil || *res[3] != 6 || res[4] == nil || *res[4] != 8 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = itoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [10 7 9 6 8], not [", result, "]")
	}

	if res := <-float_ss.SortNumerically().GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.9 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 1.0 || res[3] == nil || *res[3] != 0.8 || res[4] == nil || *res[4] != 0.6 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.9 0.7 1.0 0.8 0.6], not [", result, "]")
	}

	if res := <-float_ss.SortNumerically().By("integer_*").GetFloatsFrom("float_*"); len(res) != 5 || res[0] == nil || *res[0] != 0.8 || res[1] == nil || *res[1] != 0.7 || res[2] == nil || *res[2] != 0.6 || res[3] == nil || *res[3] != 1.0 || res[4] == nil || *res[4] != 0.9 {
		var result [5]string
		for i, sec := range res {
			if sec != nil {
				result[i] = ftoa(*sec)
			} else {
				result[i] = "nil"
			}
		}
		t.Error("Should be [0.8 0.7 0.6 1.0 0.8], not [", result, "]")
	}

	if res := <-float_ss.SortNumerically().StoreFloats(FloatStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-FloatStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 0.1 || end[1] != 0.2 || end[2] != 0.3 || end[3] != 0.4 || end[4] != 0.5 {
			t.Error("Result Should be [0.1 0.2 0.3 0.4 0.5], not", end)
		}
	}

	if res := <-float_ss.SortNumerically().GetFromAndStoreIn("string_*", StringStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-StringStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != "H" || end[1] != "J" || end[2] != "G" || end[3] != "I" || end[4] != "F" {
			t.Error("Result Should be [E C A D B], not", end)
		}
	}

	if res := <-float_ss.SortNumerically().GetIntsFromAndStoreIn("integer_*", IntStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-IntStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 10 || end[1] != 7 || end[2] != 9 || end[3] != 6 || end[4] != 8 {
			t.Error("Result Should be [10 7 9 6 8], not", end)
		}
	}

	if res := <-float_ss.SortNumerically().GetFloatsFromAndStoreIn("float_*", FloatStorage); res != 5 {
		t.Error("Should store 5 elements")
	} else {
		end := <-FloatStorage.GetFromRange(0, -1)
		if len(end) != 5 || end[0] != 0.9 || end[1] != 0.7 || end[2] != 1.0 || end[3] != 0.8 || end[4] != 0.6 {
			t.Error("Result Should be [0.9 0.7 1.0 0.8 0.6], not", end)
		}
	}

	<-float_ss.Add(0.25, 0.25)

	if res := <-float_ss.SortNumerically().GetFrom("string_*"); len(res) != 6 || res[2] != nil {
		t.Error("New element should not be found in lookup")
	}
}
