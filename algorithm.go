// SPDX-License-Identifier: 0BSD
package sx

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

var debug = false

// Print debug messages
// will not print by default, please call "SetDebugOn" first
func Debug(value any, moreValues ...any) {
	if !debug {
		return
	}
	fmt.Fprintln(os.Stderr, Str(value, moreValues...))
}
func SetDebugOn()  { debug = true }
func SetDebugOff() { debug = false }

// Throws an error by panic-ing
func Throw(errorMessages ...string) {
	var err error = errors.New(strings.Join(errorMessages, " "))
	panic(err)
}

// Throws an error by panic-ing iff the condition is true
func ThrowIf(condition bool, errorMessages ...string) {
	if condition {
		Throw(errorMessages...)
	}
}

// Catches errors thrown by "Throw" or "ThrowIf", recovers the panic
// Needs to be called before throwing errors
func Catch(handler ...func(error)) {
	if x := recover(); x != nil {
		var err = TypeCast[error](x)
		if len(handler) > 0 && !err.IsEmpty() {
			for _, h := range handler {
				h(err.Value())
			}
		} else {
			log.Fatalln("Fatal Error:", x)
			os.Exit(1)
		}
	}
}

// Typecast with Result, allowed to fail
func TypeCast[T any](it any) Result[T] {
	var r, ok = it.(T)
	if ok {
		return NewResultFrom(r)
	}
	return NewResultError[T]("Cannot typecast from '" + reflect.TypeOf(it).Name() + "' to '" + reflect.TypeOf(r).Name() + "'")
}

// Create new Array with TypeCast-ed items
// Best effort, only adds items that can be casted to the new array
func TypeCastArray[TO any, FROM any](fromArray Array[FROM]) Array[TO] {
	var arr = NewFixedSizeArray[TO](fromArray.Length())
	for it := arr.NewIterator(); it.Ok(); it.Next() {
		var val = TypeCast[TO](it.Value())
		if !val.IsEmpty() {
			arr.Push(val.Value())
		}
	}
	return arr
}

func FindFirstWhere[K comparable, V any](m Map[K, V], condition func(key K, value V) bool) Optional[Pair[K, V]] {
	for it := m.NewIterator(); it.Ok(); it.Next() {
		if condition(it.Key(), it.Value()) {
			return NewOptionalFrom(Pair[K, V]{it.Key(), it.Value()})
		}
	}
	return NewOptional[Pair[K, V]]()
}
func FindAll[K comparable, V any](m Map[K, V], condition func(key K, value V) bool) Array[Pair[K, V]] {
	var result = NewArray[Pair[K, V]]()
	for it := m.NewIterator(); it.Ok(); it.Next() {
		if condition(it.Key(), it.Value()) {
			result.Push(Pair[K, V]{Key: it.Key(), Value: it.Value()})
		}
	}
	return result
}
func ContainsValue[K comparable, V comparable](m Map[K, V], value V) bool {
	return !FindFirstWhere(m, func(_ K, v V) bool { return v == value }).IsEmpty()
}
func RemoveAll[V any](arr Array[V], condition func(value V) bool) Array[V] {
	var result = NewArray[V](arr.Length())
	for it := arr.NewIterator(); it.Ok(); it.Next() {
		if !condition(it.Value()) {
			result.Push(it.Value())
		}
	}
	return result
}

func String2Int(s string) Result[int64] {
	var val, e = strconv.ParseInt(s, 10, 64)
	if e != nil {
		NewResultError[int64]("Fatal error: cannot convert string to int64")
	}
	return NewResultFrom(val)
}

func String2Float(s string) Result[float64] {
	var val, e = strconv.ParseFloat(s, 64)
	if e != nil {
		NewResultError[int64]("Fatal error: cannot convert string to float64")
	}
	return NewResultFrom(val)
}

func String2Bool(s string) Result[bool] {
	switch s {
	default:
		return NewResultError[bool]("Fatal error: cannot convert string to bool")
	case "true":
		return NewResultFrom(true)
	case "false":
		return NewResultFrom(false)
	}
}

// Concatenates values without adding any whitespaces
func Str(value any, moreValues ...any) string {
	var result = fmt.Sprint(value)
	if len(moreValues) > 0 {
		var sb = strings.Builder{}
		sb.WriteString(result)
		for _, v := range moreValues {
			sb.WriteString(fmt.Sprint(v))
		}
		result = sb.String()
	}
	return result
}

// Shows function name of the caller of the current function
// returns empty string when in "main"
func ReflectFunctionName(optFuncLevel ...int) string {
	frameLevel := 1
	if len(optFuncLevel) > 0 {
		frameLevel = optFuncLevel[0] + 1
	}
	pc, _, _, ok := runtime.Caller(frameLevel)
	if ok {
		fp := runtime.FuncForPC(pc)
		if fp != nil {
			return fp.Name()
		}
	}
	return ""
}
