// SPDX-License-Identifier: 0BSD
package sx

var _ Container = &array[int]{}
var _ Array[int] = &array[int]{}
var _ Array[int] = NewArray[int]()
var _ Map[int, int] = &array[int]{}
var _ Iterable[int, int] = &array[int]{}
var _ Iterator[int, int] = array[int]{}.NewIterator()

// array is the main implementation for Array and Stack interfaces
type array[V any] struct {
	slice []V
}

func (v array[V]) Length() int        { return len(v.slice) }
func (v array[V]) IsEmpty() bool      { return v.Length() == 0 }
func (v array[V]) Has(index int) bool { return index < v.Length() }
func (v *array[V]) Push(value ...V)   { v.slice = append(v.slice, value...) }
func (v *array[V]) PushAll(arr Array[V]) {
	if cap(v.slice) == 0 {
		*v = array[V]{slice: make([]V, 0, arr.Length())}
	}
	v.Push(arr.Slice()...)
}

func (v array[V]) Peek() Result[V] {
	if v.IsEmpty() {
		return NewResultError[V]("Fatal error: trying to peek empty stack")
	}
	return NewResultFrom(v.slice[v.Length()-1])
}

func (v *array[V]) Pop() Result[V] {
	var temp = v.Peek()
	if temp.Ok() {
		v.slice = v.slice[:v.Length()-1]
	}
	return temp
}

func (v array[V]) Get(key int) Result[V] {
	if !v.Has(key) {
		return NewResultError[V]("Fatal error: Index out of bounds")
	}
	return NewResultFrom(v.slice[key])
}

func (v array[V]) GetOrDefault(key int) (value V) {
	if v.Has(key) {
		value = v.slice[key]
	}
	return value
}

func (v *array[V]) Put(key int, value V) bool {
	if !v.Has(key) {
		return false
	}
	v.slice[key] = value
	return true
}

func (v *array[V]) Drop(key int) bool {
	if key == 0 {
		v.slice = v.slice[1:]
	} else if key == v.Length()-1 {
		v.slice = v.slice[:v.Length()-1]
	} else if key < v.Length()-1 {
		v.slice = append(v.slice[0:int(key)], v.slice[int(key+1):]...)
	} else {
		return false
	}
	return true
}

func (v *array[V]) Compact() {
	var newVec array[V]
	newVec.Push(v.slice...)
	*v = newVec
}

func (v array[V]) Slice(fromIndexToIndex ...int) []V {
	switch len(fromIndexToIndex) {
	default:
		return v.slice
	case 1:
		return v.slice[fromIndexToIndex[0]:]
	case 2:
		return v.slice[fromIndexToIndex[0]:fromIndexToIndex[1]]
	}
}

func (v array[V]) NewIterator() Iterator[int, V] {
	var it = &ArrayIterator[int, V]{}
	it.vec = &v
	return it
}

type ArrayIterator[K int, V any] struct {
	vec   *array[V]
	index int
}

func (it ArrayIterator[int, V]) Ok() bool {
	return it.vec.Has(it.index)
}

func (it ArrayIterator[int, V]) Key() int {
	return int(it.index)
}

func (it ArrayIterator[int, V]) Value() V {
	return it.vec.GetOrDefault(it.index)
}

func (it *ArrayIterator[int, V]) Next() {
	it.index++
}
