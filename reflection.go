// SPDX-License-Identifier: 0BSD
package sx

import (
	"reflect"
	"runtime"
)

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
