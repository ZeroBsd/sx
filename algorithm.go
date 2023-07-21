// SPDX-License-Identifier: 0BSD
package sx

// Typecast with Result, allowed to fail
func TypeCast[T any](it any) Result[T] {
	var r, ok = it.(T)
	if ok {
		return NewResultFrom(r)
	}
	var sourceTypeName = ReflectType(it)
	var targetTypeName = ReflectType[T]()
	return NewResultError[T](Str("Cannot typecast from '", sourceTypeName, "' to '", targetTypeName, "'"))
}

// Create new Array with TypeCast-ed items
//
// Best effort, only adds items that can be casted to the new array
func TypeCastArray[TO any, FROM any](fromArray Array[FROM]) Array[TO] {
	var arr = NewArray[TO](fromArray.Length())
	MapValues[int, FROM, TO](fromArray.NewIterator(), func(key int, value FROM) Optional[TO] { return TypeCast[TO](value).ValueIsOptional() })
	return arr
}

// Maps the Iterator Values to a new type
//
// Can return
func MapValues[K any, V any, TO any](fromIterator Iterator[K, V], mapperFunc func(key K, value V) Optional[TO]) Array[TO] {
	var arr = NewArray[TO]()
	for ; fromIterator.Ok(); fromIterator.Next() {
		var val = mapperFunc(fromIterator.Key(), fromIterator.Value())
		if !val.IsEmpty() {
			arr.Push(val.Value())
		}
	}
	return arr
}

// Finds first entry in sx.Map or sx.Array where the condition is true and returns a key/value pair
func FindFirstWhere[K comparable, V any](m Map[K, V], condition func(key K, value V) bool) Optional[Pair[K, V]] {
	for it := m.NewIterator(); it.Ok(); it.Next() {
		if condition(it.Key(), it.Value()) {
			return NewOptionalFrom(Pair[K, V]{it.Key(), it.Value()})
		}
	}
	return NewOptional[Pair[K, V]]()
}

// Finds all entries in sx.Map or sx.Array where the condition is true and returns a list of key/value pairs
func FindAll[K comparable, V any](m Map[K, V], condition func(key K, value V) bool) Array[Pair[K, V]] {
	var result = NewArray[Pair[K, V]]()
	for it := m.NewIterator(); it.Ok(); it.Next() {
		if condition(it.Key(), it.Value()) {
			result.Push(Pair[K, V]{Key: it.Key(), Value: it.Value()})
		}
	}
	return result
}

// Checks if a sx.Map or sx.Array contains a certain value (for keys we can just use map.Has(key))
func ContainsValue[K comparable, V comparable](m Map[K, V], value V) bool {
	return !FindFirstWhere(m, func(_ K, v V) bool { return v == value }).IsEmpty()
}

// Copies sx.Array with all values removed where the condition is true and returns the copy
func RemoveAll[V any](arr Array[V], condition func(value V) bool) Array[V] {
	var result = NewArray[V](arr.Length())
	for it := arr.NewIterator(); it.Ok(); it.Next() {
		if !condition(it.Value()) {
			result.Push(it.Value())
		}
	}
	return result
}
