// SPDX-License-Identifier: 0BSD
package sx

type Iterable[K any, V any] interface {
	NewIterator() Iterator[K, V]
}

type Iterator[K any, V any] interface {
	Ok() bool
	Key() K
	Value() V
	Next()
}

type Map[K comparable, V any] interface {
	Container
	Iterable[K, V]
	Has(K) bool
	Get(K) Result[V]
	Put(K, V) bool
	Drop(K) bool
}

type Stack[V any] interface {
	Container
	Push(value ...V)
	PushAll(otherVec Array[V])
	Pop() (value Result[V])
	Peek() (value Result[V])
}

type Array[V any] interface {
	Iterable[int, V]
	Stack[V]
	Map[int, V]
	Compact()
	Slice(fromIndexToIndex ...int) []V
}

func NewArray[V any](capacity ...int) Array[V] {
	if len(capacity) > 0 {
		return &array[V]{slice: make([]V, 0, int(capacity[0]))}
	}
	return &array[V]{slice: make([]V, 0)}
}

func NewStack[V any]() Stack[V] {
	return NewArray[V]()
}

func NewFixedSizeArray[V any](numberOfItems int) Array[V] {
	if numberOfItems > 0 {
		return &array[V]{slice: make([]V, numberOfItems)}
	}
	return &array[V]{slice: make([]V, 0)}
}

func NewArrayFrom[V any](items ...V) Array[V] {
	var a = NewArray[V](len(items))
	a.Push(items...)
	return a
}
