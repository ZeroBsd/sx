// SPDX-License-Identifier: 0BSD
package sx

import (
	"reflect"
	"runtime"
)

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
	for it := fromArray.NewIterator(); it.Ok(); it.Next() {
		var val = TypeCast[TO](it.Value())
		if val.Ok() {
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

// Shows function name of the caller of the current function
//
// returns empty string when in "main"
//
// optFuncLevel 0 (default) means current Function name
//
// optFuncLevel 1 (default) means caller of current function
//
// optFuncLevel 2+ (default) means caller of caller of current function...
func ReflectFunctionName(optFuncLevel ...int) (functionName string) {
	frameLevel := 1
	if len(optFuncLevel) > 0 {
		frameLevel = optFuncLevel[0] + 1
	}
	if pc, _, _, ok := runtime.Caller(frameLevel); ok {
		if fp := runtime.FuncForPC(pc); fp != nil {
			functionName = fp.Name()
		}
	}
	return functionName
}

// Returns the reflect.Type for an instance of a type.
//
// Also works without an actual instance by passing 'nil'.
func ReflectType[T any](instances ...T) reflect.Type {
	var reflectedType reflect.Type
	if len(instances) != 1 || (any)(instances[0]) == nil {
		reflectedType = reflect.TypeOf((*T)(nil)).Elem()
	} else {
		reflectedType = reflect.TypeOf(instances[0])
	}
	return reflectedType
}

func ReflectTypeName[T any](instance ...T) string {
	var sb = NewStringBuilder()
	var reflType = ReflectType(instance...)
	reflectTypeNameImpl(sb, reflType)
	var typeName = sb.String()
	return typeName
}

func reflectTypeNameImpl(sb *StringBuilder, reflType reflect.Type) {
	var reflTypeName = reflType.Name()
	if reflTypeName != "" {
		sb.WriteString(reflTypeName)
		return
	}
	var kind = reflType.Kind()
	switch kind {
	case reflect.Pointer:
		sb.WriteString("*")
		reflectTypeNameImpl(sb, reflType.Elem())
	case reflect.Slice:
		sb.WriteString("[]")
		reflectTypeNameImpl(sb, reflType.Elem())
	case reflect.Array:
		sb.WriteString("[")
		sb.WriteAny(reflType.Len())
		sb.WriteString("]")
		reflectTypeNameImpl(sb, reflType.Elem())
	case reflect.Map:
		sb.WriteString("map[")
		reflectTypeNameImpl(sb, reflType.Key())
		sb.WriteString("]")
		reflectTypeNameImpl(sb, reflType.Elem())
	}
}
