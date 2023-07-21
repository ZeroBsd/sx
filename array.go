// SPDX-License-Identifier: 0BSD
package sx

var _ Container = &arrayImpl[int]{}
var _ Array[int] = &arrayImpl[int]{}
var _ Array[int] = NewArray[int]()
var _ Map[int, int] = &arrayImpl[int]{}
var _ Iterable[int, int] = &arrayImpl[int]{}
var _ Iterator[int, int] = arrayImpl[int]{}.NewIterator()
var _ Iterator[int, string] = &SliceIterator[int, string]{}

type arrayImpl[V any] []V

func NewArray[V any](capacity ...int) Array[V] {
	var v arrayImpl[V]
	if len(capacity) > 0 {
		v = arrayImpl[V](make([]V, 0, int(capacity[0])))
	} else {
		v = arrayImpl[V](make([]V, 0))
	}
	return &v
}

func NewStack[V any](capacity ...int) Stack[V] {
	return NewArray[V](capacity...)
}

func NewArrayFrom[V any](items ...V) Array[V] {
	var a = NewArray[V](len(items))
	a.Push(items...)
	return a
}

func (v arrayImpl[V]) Length() int        { return len(v) }
func (v arrayImpl[V]) IsEmpty() bool      { return v.Length() == 0 }
func (v arrayImpl[V]) Has(index int) bool { return index < v.Length() }
func (v *arrayImpl[V]) Push(value ...V)   { (*v) = append(*v, value...) }

func (v *arrayImpl[V]) PushArray(arr Array[V]) {
	if cap(*v) == 0 {
		(*v) = arrayImpl[V](make([]V, 0, v.Length()))
	}
	(*v) = append(*v, arr.SubSlice()...)
}

func (v arrayImpl[V]) Peek() Result[V] {
	if v.IsEmpty() {
		return NewResultError[V]("Fatal error: trying to peek empty stack")
	}
	return NewResultFrom(v[v.Length()-1])
}

func (v *arrayImpl[V]) Pop() Result[V] {
	var temp = v.Peek()
	if temp.Ok() {
		(*v) = (*v)[:v.Length()-1]
	}
	return temp
}

func (v arrayImpl[V]) Get(key int) Result[V] {
	if !v.Has(key) {
		return NewResultError[V]("Fatal error: Index out of bounds")
	}
	return NewResultFrom(v[key])
}

func (v arrayImpl[V]) Put(key int, value V) bool {
	if !v.Has(key) {
		return false
	}
	v[key] = value
	return true
}

func (v *arrayImpl[V]) Drop(key int) bool {
	if key == 0 {
		(*v) = (*v)[1:]
	} else if key == v.Length()-1 {
		(*v) = (*v)[:v.Length()-1]
	} else if key < v.Length()-1 {
		(*v) = append((*v)[0:int(key)], (*v)[int(key+1):]...)
	} else {
		return false
	}
	return true
}

func (v *arrayImpl[V]) Compact() {
	var newVec arrayImpl[V]
	newVec.Push((*v)...)
	(*v) = newVec
}

func (v arrayImpl[V]) SubSlice(fromIndexToIndex ...int) []V {
	switch len(fromIndexToIndex) {
	default:
		return v
	case 1:
		return v[fromIndexToIndex[0]:]
	case 2:
		return v[fromIndexToIndex[0]:fromIndexToIndex[1]]
	}
}

func (v arrayImpl[V]) NewIterator() Iterator[int, V] {
	var it = NewSliceIterator(v)
	return &it
}

type SliceIterator[K int, V any] struct {
	slice []V
	index int
}

func NewSliceIterator[V any](slice []V) SliceIterator[int, V] {
	var index = 0
	return SliceIterator[int, V]{
		slice: slice,
		index: index,
	}
}

func (it SliceIterator[int, V]) Ok() bool {
	return it.index < len(it.slice)
}

func (it SliceIterator[int, V]) Key() int {
	return int(it.index)
}

func (it SliceIterator[int, V]) Value() V {
	return it.slice[it.index]
}

func (it SliceIterator[int, V]) SetValue(v V) {
	it.slice[it.index] = v
}

func (it *SliceIterator[int, V]) Next() {
	it.index++
}
