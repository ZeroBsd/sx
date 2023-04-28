// SPDX-License-Identifier: 0BSD
package sx_test

import (
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestHashMap(t *testing.T) {
	var refMap = map[string]int{
		"a": 0x1010,
		"b": 0x1011,
	}
	var m = sx.NewMap[string, int]()
	if !m.IsEmpty() {
		t.Fail()
	}
	m.Put("a", 0x1010)
	m.Put("b", 0x1011)
	var fm = sx.NewMapFrom(map[string]int{
		"a": 0x1010,
		"b": 0x1011,
	})
	var hm = sx.NewHashMap[string, int]()
	hm.Put("a", 0x1010)
	hm.Put("b", 0x1011)
	var gm = sx.NewHashMapFrom(map[string]int{
		"a": 0x1010,
		"b": 0x1011,
	})
	if m.IsEmpty() {
		t.Fail()
	}
	if len(refMap) != m.Length() || len(refMap) != hm.Length() || len(refMap) != gm.Length() || len(refMap) != fm.Length() {
		t.Fail()
	}
	if m.Get("b").Value() != refMap["b"] {
		t.FailNow()
	}
	var it2 = m.NewIterator()
	_ = it2

	for it := m.NewIterator(); it.Ok(); it.Next() {
		if it.Value() != refMap[it.Key()] {
			t.FailNow()
		}
		if it.Value() != hm.Get(it.Key()).Value() {
			t.FailNow()
		}
		if it.Value() != gm.Get(it.Key()).Value() {
			t.FailNow()
		}
		if it.Value() != fm.Get(it.Key()).Value() {
			t.FailNow()
		}
	}
}

func TestMapSet(t *testing.T) {
	var set = sx.NewSetFrom("a", "b", "c")
	if set.Get("b").Value() != true {
		t.FailNow()
	}
	if set.Get("d").Ok() {
		t.FailNow()
	}
	if !set.Drop("a") {
		t.FailNow()
	}
	if set.Drop("e") {
		t.FailNow()
	}
}

func TestIterationWithDeletion(t *testing.T) {
	{
		var m = sx.NewMap[string, int]()
		m.Put("a", 1)
		m.Put("b", 2)
		m.Put("c", 3)
		m.Put("d", 4)
		var foundKeys = sx.NewArray[string]()
		for it := m.NewIterator(); it.Ok(); it.Next() {
			foundKeys.Push(it.Key())
			if it.Key() == "b" {
				m.Drop("c")
			}
		}
		if sx.ContainsValue[int, string](foundKeys, "c") || !sx.ContainsValue[int, string](foundKeys, "a") || !sx.ContainsValue[int, string](foundKeys, "b") || !sx.ContainsValue[int, string](foundKeys, "d") {
			t.FailNow()
		}
	}
}
