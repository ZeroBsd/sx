// SPDX-License-Identifier: 0BSD
package sx_test

import (
	"reflect"
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestStack(t *testing.T) {
	var s = sx.NewStack[int]()
	sx.ThrowIf(s.Length() != 0)
	s.Push(1)
	sx.ThrowIf(s.Length() != 1 || s.Peek().Value() != 1)
	s.Push(2)
	sx.ThrowIf(s.Length() != 2 || s.Peek().Value() != 2)
	s.Push(4)
	sx.ThrowIf(s.Length() != 3 || s.Peek().Value() != 4)
	s.Pop()
	sx.ThrowIf(s.Length() != 2 || s.Peek().Value() != 2)
	s.Pop()
	sx.ThrowIf(s.Length() != 1 || s.Peek().Value() != 1)
	s.Pop()
	sx.ThrowIf(s.Length() != 0)
}

func TestPeekEmptyStack(t *testing.T) {
	var s = sx.NewStack[int]()
	var i = s.Peek()
	sx.ThrowIf(i.Ok())
}

func TestArrayFrom(t *testing.T) {
	var a = sx.NewArrayFrom(2, 3, 4, 5, 6, 7)
	sx.ThrowIf(a.Length() != 6)
	sx.ThrowIf(!a.Has(0))
	sx.ThrowIf(!a.Has(5))
	sx.ThrowIf(a.Has(6))
	var sliceAll = a.SubSlice()
	sx.ThrowIf(!reflect.DeepEqual(sliceAll, []int{2, 3, 4, 5, 6, 7}))
	var slice2 = a.SubSlice(2)
	sx.ThrowIf(!reflect.DeepEqual(slice2, []int{4, 5, 6, 7}))
	var slice24 = a.SubSlice(2, 4)
	sx.ThrowIf(!reflect.DeepEqual(slice24, []int{4, 5}))
}

func TestArrayPushArray(t *testing.T) {
	var a = sx.NewArray[int]()
	a.PushArray(sx.NewArrayFrom(2, 3, 4))
	sx.ThrowIf(!reflect.DeepEqual(a.SubSlice(), []int{2, 3, 4}))
	a.PushArray(sx.NewArrayFrom(5, 6, 7))
	sx.ThrowIf(a.Length() != 6)
	sx.ThrowIf(!reflect.DeepEqual(a.SubSlice(), []int{2, 3, 4, 5, 6, 7}))
	sx.ThrowIf(!a.Has(2) || a.Get(2).Value() != 4)
	sx.ThrowIf(!a.Has(5) || a.Get(5).Value() != 7)
}

func TestArrayPushPopGet(t *testing.T) {
	var a = sx.NewArrayFrom(2)
	sx.ThrowIf(!reflect.DeepEqual(a.SubSlice(), []int{2}))
	a.Push(3)
	sx.ThrowIf(!reflect.DeepEqual(a.SubSlice(), []int{2, 3}))
	sx.ThrowIf(a.Get(1).Value() != 3)
	a.PushArray(sx.NewArrayFrom(4, 5))
	sx.ThrowIf(!reflect.DeepEqual(a.SubSlice(), []int{2, 3, 4, 5}))
	sx.ThrowIf(a.Get(2).Value() != 4)
	sx.ThrowIf(a.Get(3).Value() != 5)
	sx.ThrowIf(a.Length() != 4)
	a.Pop()
	sx.ThrowIf(!reflect.DeepEqual(a.SubSlice(), []int{2, 3, 4}))
	sx.ThrowIf(a.Length() != 3)
	sx.ThrowIf(a.Get(2).Value() != 4)
}

func TestArrayGetFailure(t *testing.T) {
	var a = sx.NewArrayFrom(3)
	var _ = a.Get(0).Value()
	var value = a.Get(1).ValueOr(99)
	sx.ThrowIf(value != 99)
}

func TestArrayPut(t *testing.T) {
	var a = sx.NewArrayFrom(2, 3, 4)
	sx.ThrowIf(!a.Put(1, 5))
	sx.ThrowIf(a.Put(3, 5))
}

func TestArrayDrop(t *testing.T) {
	var a = sx.NewArrayFrom(2, 3, 4, 5)
	sx.ThrowIf(!reflect.DeepEqual(a.SubSlice(), []int{2, 3, 4, 5}))
	sx.ThrowIf(a.Drop(4))
	sx.ThrowIf(!reflect.DeepEqual(a.SubSlice(), []int{2, 3, 4, 5}))
	sx.ThrowIf(!a.Drop(3))
	sx.ThrowIf(!reflect.DeepEqual(a.SubSlice(), []int{2, 3, 4}))
	sx.ThrowIf(!a.Drop(1))
	sx.ThrowIf(!reflect.DeepEqual(a.SubSlice(), []int{2, 4}))
	sx.ThrowIf(a.Drop(2))
	sx.ThrowIf(!reflect.DeepEqual(a.SubSlice(), []int{2, 4}))
	sx.ThrowIf(!a.Drop(0))
	sx.ThrowIf(!reflect.DeepEqual(a.SubSlice(), []int{4}))
}

func TestCompact(t *testing.T) {
	var a = sx.NewArrayFrom(2, 3, 4, 5, 6, 7)
	a.Push(8, 9, 10)
	var normal = cap(a.SubSlice())
	a.Pop()
	a.Pop()
	a.Pop()
	a.Pop()
	a.Compact()
	var compacted = cap(a.SubSlice())
	sx.ThrowIf(compacted >= normal)
}

func TestArrayIterationWithDeletion(t *testing.T) {
	{
		var arr = sx.NewArrayFrom("a", "b", "c", "d")
		for it := arr.NewIterator(); it.Ok(); it.Next() {
			if it.Key() == 1 {
				arr.Drop(2) //Drop "c" when at "b"
			}
		}
		if arr.Length() != 3 || arr.Get(0).Value() != "a" || arr.Get(1).Value() != "b" || arr.Get(2).Value() != "d" {
			t.FailNow()
		}
	}

}

func TestSliceIterator(t *testing.T) {
	var s = []string{"a", "b"}
	var it = sx.NewSliceIterator(s)
	if !it.Ok() || it.Key() != 0 || it.Value() != "a" {
		t.FailNow()
	}
	it.Next()
	if !it.Ok() || it.Key() != 1 || it.Value() != "b" {
		t.FailNow()
	}
	it.SetValue("c")
	if !it.Ok() || it.Value() != "c" {
		t.FailNow()
	}
	if len(s) != 2 {
		t.FailNow()
	}
}

var benchmarkResult int

func BenchmarkSliceIterator(b *testing.B) {
	b.StopTimer()
	var a = sx.NewArray[int](5000)
	for i := 0; i < 5000; i++ {
		a.Push(i)
	}
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		for it := a.NewIterator(); it.Ok(); it.Next() {
			benchmarkResult = it.Value()
		}
	}
}
