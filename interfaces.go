// SPDX-License-Identifier: 0BSD
package sx

type Container interface {
	Length() int
	IsEmpty() bool
}

type Iterable[K any, V any] interface {
	NewIterator() Iterator[K, V]
}

type Iterator[K any, V any] interface {
	Ok() bool // Checks if iterator is valid and can return a key or a value. For empty containers, this is false after iterator creation.
	Key() K   // Returns the key of the current item
	Value() V // Returns the value of the current item
	Next()    // Finds the next item. If it reached the end, Ok() will return false after this call
}

type Set[K comparable, V any] interface {
	Container
	Iterable[K, V]
	Has(K) bool    // Checks if key is present
	Put(K, V) bool // Sets value for a key
	Drop(K) bool   // Removes a key and its value
}

type Map[K comparable, V any] interface {
	Set[K, V]
	Get(K) Result[V] // Returns Value for the key (if present)
}

type Stack[V any] interface {
	Container
	Push(value ...V)             // Inserts all values at the end
	PushArray(otherVec Array[V]) // Inserts all values of an sx.Array at the end
	Pop() (value Result[V])      // Removes and returns last value (if present)
	Peek() (value Result[V])     // Returns last value (if present)
}

type Array[V any] interface {
	Iterable[int, V]
	Stack[V]
	Map[int, V]
	Compact()                             // Copies Array and removes excessive memory
	SubSlice(fromIndexToIndex ...int) []V // Returns a go slice
}
