package redis

import (
	"testing"
)

func TestHashes(t *testing.T) {
	r := GetRedis(t)
	defer r.Close()

	h := r.Hash("Test_Hash")
	h.Delete()

	//basic hash funcs
	if res := <-h.Size(); res != 0 {
		t.Error("An empty hash should have size 0, not", res)
	}

	if res := <-h.Get(); len(res) != 0 {
		t.Error("Getting from an empty hash should yield an empty map, not", res)
	}

	//hash string field funcs
	s := h.String("String")

	if _, ok := <-s.Get(); ok {
		t.Error("Should not have anything to get")
	}

	if !<-s.SetIfEmpty("A") {
		t.Error("Should be able to set if empty on an empty string")
	}

	if res, ok := <-s.Get(); !ok || res != "A" {
		t.Error("Should be A, not", res)
	}

	if <-s.SetIfEmpty("B") {
		t.Error("Should not be able to set if empty on a non-empty string")
	}

	if res := <-s.Get(); res != "A" {
		t.Error("Should be A, not", res)
	}

	if <-s.Set("C") {
		t.Error("Should overwrite")
	}

	if res := <-s.Get(); res != "C" {
		t.Error("Should be C, not", res)
	}

	if !<-s.Exists() {
		t.Error("Should exist")
	}

	if !<-s.Delete() {
		t.Error("Should be able to delete")
	}

	if _, ok := <-s.Get(); ok {
		t.Error("Should not have anything to get")
	}

	if <-s.Delete() {
		t.Error("Should already be deleted")
	}

	if <-s.Exists() {
		t.Error("Should not exist")
	}

	if !<-s.Set("D") {
		t.Error("Should have created a new value")
	}

	if res, ok := <-s.Get(); !ok || res != "D" {
		t.Error("Should be D, not", res)
	}

	//hash int field funcs
	i := h.Integer("Integer")

	if _, ok := <-i.Get(); ok {
		t.Error("Should not have anything to get")
	}

	if !<-i.SetIfEmpty(1) {
		t.Error("Should be able to set if empty on an empty string")
	}

	if res, ok := <-i.Get(); !ok || res != 1 {
		t.Error("Should be 1, not", res)
	}

	if <-i.SetIfEmpty(2) {
		t.Error("Should not be able to set if empty on a non-empty integer")
	}

	if res := <-i.Get(); res != 1 {
		t.Error("Should be 1, not", res)
	}

	if <-i.Set(3) {
		t.Error("Should overwrite")
	}

	if res := <-i.Get(); res != 3 {
		t.Error("Should be 3, not", res)
	}

	if !<-i.Exists() {
		t.Error("Should exist")
	}

	if !<-i.Delete() {
		t.Error("Should be able to delete")
	}

	if _, ok := <-i.Get(); ok {
		t.Error("Should not have anything to get")
	}

	if <-i.Delete() {
		t.Error("Should already be deleted")
	}

	if <-i.Exists() {
		t.Error("Should not exist")
	}

	if !<-i.Set(4) {
		t.Error("Should have created a new value")
	}

	if res, ok := <-i.Get(); !ok || res != 4 {
		t.Error("Should be 4, not", res)
	}

	if res := <-i.IncrementBy(3); res != 7 {
		t.Error("Should be 7, not", res)
	}

	if res := <-i.DecrementBy(2); res != 5 {
		t.Error("Should be 5, not", res)
	}

	//hash float field funcs
	f := h.Float("Float")

	if _, ok := <-f.Get(); ok {
		t.Error("Should not have anything to get")
	}

	if !<-f.SetIfEmpty(.1) {
		t.Error("Should be able to set if empty on an empty string")
	}

	if res, ok := <-f.Get(); !ok || res != .1 {
		t.Error("Should be 1, not", res)
	}

	if <-f.SetIfEmpty(.2) {
		t.Error("Should not be able to set if empty on a non-empty integer")
	}

	if res := <-f.Get(); res != .1 {
		t.Error("Should be .1, not", res)
	}

	if <-f.Set(.3) {
		t.Error("Should overwrite")
	}

	if res := <-f.Get(); res != .3 {
		t.Error("Should be .3, not", res)
	}

	if !<-f.Exists() {
		t.Error("Should exist")
	}

	if !<-f.Delete() {
		t.Error("Should be able to delete")
	}

	if _, ok := <-f.Get(); ok {
		t.Error("Should not have anything to get")
	}

	if <-f.Delete() {
		t.Error("Should already be deleted")
	}

	if <-f.Exists() {
		t.Error("Should not exist")
	}

	if !<-f.Set(.4) {
		t.Error("Should have created a new value")
	}

	if res, ok := <-f.Get(); !ok || res != .4 {
		t.Error("Should be 4, not", res)
	}

	if res := <-f.IncrementBy(.3); res != .7 {
		t.Error("Should be 7, not", res)
	}

	if res := <-f.DecrementBy(.2); res != .5 {
		t.Error("Should be 5, not", res)
	}

	//basic hash funcs again

	if res := <-h.Size(); res != 3 {
		t.Error("should have size 3, not", res)
	}

	if res := <-h.Get(); len(res) != 3 || res["String"] != "D" || res["Integer"] != "5" || res["Float"] != "0.5" {
		t.Error("Should have map[String:D Integer:5 Float:0.5], not", res)
	}

}
