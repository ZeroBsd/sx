// SPDX-License-Identifier: 0BSD
package sx

var _ Container = hashMapImpl[string, int]{}
var _ Map[string, int] = NewHashMap[string, int]()
var _ Map[string, int] = hashMapImpl[string, int]{}
var _ Iterable[string, int] = hashMapImpl[string, int]{}
var _ Iterator[string, int] = hashMapImpl[string, int]{}.NewIterator()

func NewMap[K comparable, V any]() Map[K, V]                   { return NewHashMap[K, V]() }
func NewMapFrom[K comparable, V any](values map[K]V) Map[K, V] { return NewHashMapFrom(values) }

func NewHashMap[K comparable, V any]() Map[K, V] {
	var m hashMapImpl[K, V]
	m.Map = make(map[K]V)
	return &m
}

func NewHashMapFrom[K comparable, V any](values map[K]V) Map[K, V] {
	var m = NewHashMap[K, V]()
	for key, value := range values {
		m.Put(key, value)
	}
	return m
}

func NewSet[K comparable]() Map[K, struct{}] { return NewSetFrom[K]() }
func NewSetFrom[K comparable](values ...K) Map[K, struct{}] {
	var hs = NewHashMap[K, struct{}]()
	for _, v := range values {
		hs.Put(v, struct{}{})
	}
	return hs
}

type hashMapImpl[K comparable, V any] struct {
	Map map[K]V // reference type; needs to be public for reflection access (go maps can be iterated via reflection)
}

func (m hashMapImpl[K, V]) Length() int {
	return len(m.Map)
}

func (m hashMapImpl[K, V]) IsEmpty() bool {
	return m.Length() == 0
}

func (m hashMapImpl[K, V]) Has(key K) bool {
	var _, ok = m.Map[key]
	return ok
}

func (m hashMapImpl[K, V]) Get(key K) Result[V] {
	var value, ok = m.Map[key]
	if !ok {
		return NewResultError[V]("Fatal error: Key does not exist")
	}
	return NewResultFrom(value)
}

func (m hashMapImpl[K, V]) Put(key K, value V) bool {
	m.Map[key] = value
	return true
}

func (m hashMapImpl[K, V]) Drop(key K) bool {
	if m.Has(key) {
		delete(m.Map, key)
		return true
	}
	return false
}

func (m hashMapImpl[K, V]) NewIterator() Iterator[K, V] {
	return NewMapIterator(m.Map)
}

type mapIterator[K comparable, V any] struct {
	Map         map[K]V
	keyIterator Iterator[int, K]
}

func NewMapIterator[K comparable, V any](m map[K]V) Iterator[K, V] {
	var keys = NewArray[K](len(m))
	for k := range m {
		keys.Push(k)
	}
	var it = &mapIterator[K, V]{
		Map:         m,
		keyIterator: keys.NewIterator(),
	}
	return it
}

func (it mapIterator[K, V]) Ok() bool {
	for kit := it.keyIterator; kit.Ok(); kit.Next() {
		var _, ok = it.Map[kit.Value()]
		if ok {
			return ok
		}
	}
	return false
}

func (it mapIterator[K, V]) Key() K {
	return it.keyIterator.Value()
}

func (it mapIterator[K, V]) Value() V {
	return it.Map[it.keyIterator.Value()]
}

func (it *mapIterator[K, V]) Next() {
	it.keyIterator.Next()
}
