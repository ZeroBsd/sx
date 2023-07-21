// SPDX-License-Identifier: 0BSD
package sx_test

import (
	"errors"
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestReflectFunctionName(t *testing.T) {
	if reflectedCallerFuncName := subFunc(); reflectedCallerFuncName != "github.com/ZeroBsd/sx_test.TestReflectFunctionName" {
		t.FailNow()
	}
	if reflectedFunctionName := sx.ReflectFunctionName(-1); reflectedFunctionName != "github.com/ZeroBsd/sx.ReflectFunctionName" {
		t.FailNow()
	}
	var lambdaFunc = func() string {
		return sx.ReflectFunctionName(1)
	}
	if reflectedCallerNameFromLambda := lambdaFunc(); reflectedCallerNameFromLambda != "github.com/ZeroBsd/sx_test.TestReflectFunctionName" {
		t.FailNow()
	}
}
func subFunc() string {
	return sx.ReflectFunctionName(1)
}

func TestReflectTypeName(t *testing.T) {
	{
		var instance int = 42
		var typeName1 = sx.ReflectType(instance).Name()
		var typeName2 = sx.ReflectTypeName(instance)
		if typeName1 != "int" || typeName2 != typeName1 {
			t.FailNow()
		}
	}
	{
		var temp int = 42
		var instance *int = &temp
		var typeName1 = sx.ReflectType(instance).Name()
		var typeName2 = sx.ReflectTypeName(instance)
		if typeName1 != "" || typeName2 != "*int" {
			t.FailNow()
		}
	}
	{
		var instance [2]int = [...]int{42, 4711}
		var typeName1 = sx.ReflectType(instance).Name()
		var typeName2 = sx.ReflectTypeName(instance)
		if typeName1 != "" || typeName2 != "[2]int" {
			t.FailNow()
		}
	}
	{
		var instance []int = []int{42}
		var typeName1 = sx.ReflectType(instance).Name()
		var typeName2 = sx.ReflectTypeName(instance)
		if typeName1 != "" || typeName2 != "[]int" {
			t.FailNow()
		}
	}
	{
		var instance map[int]string = map[int]string{42: "some string"}
		var typeName1 = sx.ReflectType(instance).Name()
		var typeName2 = sx.ReflectTypeName(instance)
		if typeName1 != "" || typeName2 != "map[int]string" {
			t.FailNow()
		}
	}
	{
		var temp = "some string"
		var instance map[int][]*string = map[int][]*string{42: {&temp}}
		var typeName1 = sx.ReflectType(instance).Name()
		var typeName2 = sx.ReflectTypeName(instance)
		if typeName1 != "" || typeName2 != "map[int][]*string" {
			t.FailNow()
		}
	}
	{
		var instance error = errors.New("some error")
		var typeName1 = sx.ReflectType(instance).Name()
		var typeName2 = sx.ReflectTypeName(instance)
		if typeName1 != "" || typeName2 != "*errorString" {
			// this might fail when the internal implementation in Go changes!
			t.FailNow()
		}
	}
}

func TestReflectType(t *testing.T) {
	if reflectedTypeName := sx.ReflectType[int]().Name(); reflectedTypeName != "int" {
		t.FailNow()
	}

	if reflectedTypeName := sx.ReflectType[error](nil).Name(); reflectedTypeName != "error" {
		t.FailNow()
	}

	if reflectedTypeName := sx.ReflectType(errors.New("some error")).Name(); reflectedTypeName != "" {
		// there is no name because the type is a pointer (more precisely: *errorString)
		t.FailNow()
	}

	if reflectedTypeName := sx.ReflectType((error)(MyError{}), errors.New("some error")).Name(); reflectedTypeName != "error" {
		// there is a name because we passeed multiple arguments, all of them are of interface 'error'
		t.FailNow()
	}

	if reflectedTypeName := reflectWithinGenericFuncWithBasicType(42); reflectedTypeName != "int" {
		t.FailNow()
	}

	if reflectedTypeName := reflectWithinGenericFuncWithIterfaceType(42); reflectedTypeName != "int" {
		t.FailNow()
	}

	var mye = MyError{}
	if reflectedTypeName := sx.ReflectType(mye).Name(); reflectedTypeName != "MyError" {
		t.FailNow()
	}
	var err error = MyError{}
	if reflectedTypeName := sx.ReflectType(err).Name(); reflectedTypeName != "MyError" {
		t.FailNow()
	}

	err = nil
	if reflectedTypeName := sx.ReflectType(err).Name(); reflectedTypeName != "error" {
		t.FailNow()
	}
}

type MyError struct{}

func (e MyError) Error() string {
	return "MyError"
}
func reflectWithinGenericFuncWithBasicType[T int](x T) string {
	var name = sx.ReflectType[T]().Name()
	return name
}
func reflectWithinGenericFuncWithIterfaceType[T any](x T) string {
	var name = sx.ReflectType[T]().Name()
	return name
}
