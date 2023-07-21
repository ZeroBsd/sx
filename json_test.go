// SPDX-License-Identifier: 0BSD
package sx_test

import (
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestToJson(t *testing.T) {
	if json := sx.ToJson("test").ValueOrInit(); json != `"test"` {
		t.FailNow()
	}
	if json := sx.ToJson(42).ValueOrInit(); json != `42` {
		t.FailNow()
	}
	if json := sx.ToJson(JsonTestStruct{X: 42, Y: "abc"}).ValueOrInit(); json != `{"X":42,"Y":"abc"}` {
		t.FailNow()
	}
	if json := sx.ToJsonPretty(JsonTestStruct{X: 42, Y: "abc"}, "  ").ValueOrInit(); json != "{\n  \"X\": 42,\n  \"Y\": \"abc\"\n}" {
		t.FailNow()
	}
	if json := sx.ToJson([]any{"a", "b", 42}).ValueOrInit(); json != `["a","b",42]` {
		t.FailNow()
	}
	if json := sx.ToJson(sx.NewArrayFrom("a", "b")).ValueOrInit(); json != `["a","b"]` {
		t.FailNow()
	}
	if json := sx.ToJson(map[string]int{"a": 42}).ValueOrInit(); json != `{"a":42}` {
		t.FailNow()
	}
	if json := sx.ToJson(sx.NewMapFrom(map[string]int{"a": 42})).ValueOrInit(); json != `{"a":42}` {
		t.FailNow()
	}

	// must fail because of recursion
	var recs = JsonRecursiveTestStruct{}
	recs.R = &recs
	if sx.ToJson(recs).Ok() {
		t.FailNow()
	}
	if sx.ToJsonPretty(recs, "  ").Ok() {
		t.FailNow()
	}
}

func TestFromJson(t *testing.T) {
	if obj := sx.FromJson[string](`"abc"`).ValueOrInit(); obj != "abc" {
		t.FailNow()
	}
	if obj := sx.FromJson[int](`42`).ValueOrInit(); obj != 42 {
		t.FailNow()
	}
	if obj := sx.FromJson[JsonTestStruct](`{"X":42,"Y":"abc"}`).ValueOrInit(); (obj != JsonTestStruct{X: 42, Y: "abc"}) {
		t.FailNow()
	}
	if obj := sx.FromJson[JsonTestStruct]("{\n  \"X\": 42,\n  \"Y\": \"abc\"\n}").ValueOrInit(); (obj != JsonTestStruct{X: 42, Y: "abc"}) {
		t.FailNow()
	}
	if obj := sx.FromJson[[]any](`["a","b",42]`).ValueOrInit(); len(obj) != 3 || obj[0].(string) != "a" || obj[1].(string) != "b" || obj[2].(float64) != 42 {
		t.FailNow()
	}

	var arrayTestInstance = sx.NewArrayFrom("a", "b")
	if obj := sx.FromJsonInterface(`["a","b"]`, &arrayTestInstance).ValueOrInit(); obj.Length() != 2 || obj.Get(0).ValueOrInit() != "a" || obj.Get(1).ValueOrInit() != "b" {
		t.FailNow()
	}

	var mapTestInstance = sx.NewMapFrom(map[string]int{})
	if obj := sx.FromJsonInterface(`{"a":42}`, &mapTestInstance).ValueOrInit(); obj.Length() != 1 || obj.Get("a").ValueOrInit() != 42 {
		t.FailNow()
	}

	// negative test, float64 is not a string
	if obj := sx.FromJson[string](`42`).ValueOrInit(); obj != "" {
		t.FailNow()
	}

	// negative test, 42 is not sx.Array[string]
	var arrayNegTestInstance = sx.NewArray[string](0)
	if obj := sx.FromJsonInterface(`42`, &arrayNegTestInstance).ValueOrInit(); obj != nil {
		t.FailNow()
	}
}

func TestJsonWithArray(t *testing.T) {
	var a = sx.NewArrayFrom(1, 2)
	var json = sx.ToJson(a)
	if json.ValueOrInit() != `[1,2]` {
		t.FailNow()
	}

	sx.FromJsonInterface[sx.Array[int]](`[3,4]`, &a)
	if a.Get(0).ValueOrInit() != 3 || a.Get(1).ValueOrInit() != 4 {
		t.FailNow()
	}
}

func TestJsonWithMap(t *testing.T) {
	var m = sx.NewHashMapFrom(map[string]int{"a": 1, "b": 2})
	var json = sx.ToJson(m)
	if json.ValueOrInit() != `{"a":1,"b":2}` {
		t.FailNow()
	}

	sx.FromJsonInterface[sx.Map[string, int]](`{"c":3,"d":4}`, &m)
	if m.Get("c").ValueOrInit() != 3 || m.Get("d").ValueOrInit() != 4 {
		t.FailNow()
	}
	// test the merging, a and b should still be valid
	if m.Get("a").ValueOrInit() != 1 || m.Get("b").ValueOrInit() != 2 {
		t.FailNow()
	}
}
