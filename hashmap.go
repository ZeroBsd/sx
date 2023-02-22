// SPDX-License-Identifier: 0BSD
package sx

var _ Container = hashMap[string, int]{}
var _ Map[string, int] = NewHashMap[string, int]()
var _ Map[string, int] = hashMap[string, int]{}
var _ Iterable[string, int] = hashMap[string, int]{}
var _ Iterator[string, int] = hashMap[string, int]{}.NewIterator()

func NewMap[K comparable, V any]() Map[K, V]                   { return NewHashMap[K, V]() }
func NewMapFrom[K comparable, V any](values map[K]V) Map[K, V] { return NewHashMapFrom(values) }

func NewHashMap[K comparable, V any]() Map[K, V] {
	var m hashMap[K, V]
	m.data = make(map[K]V)
	return &m
}

func NewHashMapFrom[K comparable, V any](values map[K]V) Map[K, V] {
	var m = NewHashMap[K, V]()
	for key, value := range values {
		m.Put(key, value)
	}
	return m
}

func NewSetFrom[K comparable](values ...K) Map[K, bool] {
	var hs = NewHashMap[K, bool]()
	for _, v := range values {
		hs.Put(v, true)
	}
	return hs
}

type hashMap[K comparable, V any] struct {
	data map[K]V
}

func (m hashMap[K, V]) Length() int {
	return len(m.data)
}

func (m hashMap[K, V]) IsEmpty() bool {
	return m.Length() == 0
}

func (m hashMap[K, V]) Has(key K) bool {
	var _, ok = m.data[key]
	return ok
}

func (m hashMap[K, V]) Get(key K) Result[V] {
	var value, ok = m.data[key]
	if !ok {
		return NewResultError[V]("Fatal error: Key does not exist")
	}
	return NewResultFrom(value)
}

func (m hashMap[K, V]) GetOrDefault(key K) (value V) {
	if m.Has(key) {
		value = m.data[key]
	}
	return value
}

func (m hashMap[K, V]) Put(key K, value V) bool {
	m.data[key] = value
	return true
}

func (m hashMap[K, V]) Drop(key K) bool {
	if m.Has(key) {
		delete(m.data, key)
		return true
	}
	return false
}

func (m hashMap[K, V]) NewIterator() Iterator[K, V] {
	var it = &HashMapIterator[K, V]{}
	if m.Length() == 0 {
		return it
	}
	it.mp = &m
	it.keys = NewFixedSizeArray[K](len(m.data))
	it.currentKeyIndex = 0
	var i = 0
	for k := range m.data {
		it.keys.Put(i, k)
		i++
	}
	it.current.Key = it.keys.Get(0).ValueOrDefault()
	it.current.Value = it.mp.GetOrDefault(it.current.Key)
	return it
}

type HashMapIterator[K comparable, V any] struct {
	mp              *hashMap[K, V]
	keys            Array[K]
	currentKeyIndex int
	current         Pair[K, V]
}

func (it HashMapIterator[K, V]) Ok() bool {
	return it.currentKeyIndex < it.keys.Length() && it.mp.Has(it.current.Key)
}

func (it HashMapIterator[K, V]) Key() K {
	return it.current.Key
}
func (it HashMapIterator[K, V]) Value() V {
	return it.current.Value
}

func (it *HashMapIterator[K, V]) Next() {
	it.currentKeyIndex++
	var nextKey = it.keys.Get(it.currentKeyIndex)
	// if the nextKey is not valid, it has been removed during iteration
	// this is supported - so we try and find the next valid key, until our key-array runs out of items
	for !nextKey.IsEmpty() && !it.mp.Has(nextKey.Value()) && it.currentKeyIndex < it.keys.Length() {
		it.currentKeyIndex++
		nextKey = it.keys.Get(it.currentKeyIndex)
	}
	it.current.Key = nextKey.ValueOrDefault()
	it.current.Value = it.mp.GetOrDefault(it.current.Key)
}
