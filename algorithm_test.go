// SPDX-License-Identifier: 0BSD
package sx_test

import (
	"errors"
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestTypeCastWithInterface(t *testing.T) {
	var ti testInterface = testImpl{}
	ti.makeInterfaceCall()
	var target = sx.TypeCast[testImpl](ti)
	if !target.Ok() {
		t.FailNow()

	}
	target.Value().makeInterfaceCall()
}

func TestFailingTypeCastWithInterface(t *testing.T) {
	var someStruct = struct{}{}
	var timpl = sx.TypeCast[testImpl](someStruct)
	if timpl.Ok() {
		t.FailNow()
	}
}

func TestTypeCastWithBasicType(t *testing.T) {
	var iany any = 3
	var iok = sx.TypeCast[int](iany)
	if !iok.Ok() {
		t.FailNow()
	}
	var i32 int32 = 42
	var i64 = sx.TypeCast[int64](i32)
	if i64.Ok() {
		t.FailNow()
	}
	var str = "some string"
	var i = sx.TypeCast[int](str)
	if i.Ok() {
		t.FailNow()
	}
}

func TestTypeCastWithError(t *testing.T) {
	var ti testInterface = testImpl{}
	var i2 = sx.TypeCast[error](ti)
	_ = i2

	var str = "some string"
	var i = sx.TypeCast[testInterface](str)
	if i.Ok() {
		t.FailNow()
	}
}

type testInterface interface {
	makeInterfaceCall()
}
type testImpl struct{}

func (t testImpl) makeInterfaceCall() {
	return
}

func TestArrayTypeCast(t *testing.T) {
	var concreteArray = sx.NewArrayFrom(testImpl{}, testImpl{})
	var abstractArray = sx.NewArray[any]()
	for it := concreteArray.NewIterator(); it.Ok(); it.Next() {
		abstractArray.Push(it.Value())
	}
	var concreteArray2 = sx.TypeCastArray[testImpl](abstractArray)
	for it := concreteArray2.NewIterator(); it.Ok(); it.Next() {
		it.Value().makeInterfaceCall()
	}
}

func TestFirstFindWhere(t *testing.T) {
	var arr = sx.NewArrayFrom(1, 2, 3, 2, 4)
	{
		var r = sx.FindFirstWhere[int, int](arr, func(key int, value int) bool {
			return value == 2
		})
		if r.IsEmpty() || r.Value().Key != 1 || r.Value().Value != 2 {
			t.FailNow()
		}
	}
	{
		var r = sx.FindFirstWhere[int, int](arr, func(key int, value int) bool {
			return value == 5
		})
		if r.Ok() {
			t.Log("Error optional should be empty, but was [key: ", r.Value().Key, "; value: ", r.Value().Value, "]")
			t.FailNow()
		}
	}
}

func TestContainsValue(t *testing.T) {
	var arr = sx.NewArrayFrom(1, 2, 3, 2, 4)
	if sx.ContainsValue[int, int](arr, 5) || !sx.ContainsValue[int, int](arr, 1) || !sx.ContainsValue[int, int](arr, 2) || !sx.ContainsValue[int, int](arr, 3) || !sx.ContainsValue[int, int](arr, 4) {
		t.FailNow()
	}
}

func TestFirstFindAll(t *testing.T) {
	var arr = sx.NewArrayFrom(1, 2, 3, 2, 4)
	{
		var r = sx.FindAll[int, int](arr, func(key int, value int) bool {
			return value == 3
		})
		if r.Length() != 1 || r.Get(0).Value().Key != 2 || r.Get(0).Value().Value != 3 {
			t.FailNow()
		}
	}
	{
		var r = sx.FindAll[int, int](arr, func(key int, value int) bool {
			return value == 2
		})
		if r.Length() != 2 || r.Get(0).Value().Key != 1 || r.Get(0).Value().Value != 2 || r.Get(1).Value().Key != 3 || r.Get(1).Value().Value != 2 {
			t.FailNow()
		}
	}
	{
		var r = sx.FindAll[int, int](arr, func(key int, value int) bool {
			return value == 5
		})
		if r.Length() > 0 {
			t.FailNow()
		}
	}
}

func TestRemoveAll(t *testing.T) {
	{
		var arr = sx.NewArrayFrom(1, 2, 3, 2, 4)
		var arr2 = sx.RemoveAll(arr, func(x int) bool {
			return x == 3
		})
		if arr2.Length() != arr.Length()-1 || sx.ContainsValue[int, int](arr2, 3) {
			t.Fail()
		}
	}
	{
		var arr = sx.NewArrayFrom(1, 2, 3, 2, 4)
		var arr2 = sx.RemoveAll(arr, func(x int) bool {
			return x == 2
		})
		if arr2.Length() != arr.Length()-2 || sx.ContainsValue[int, int](arr2, 2) {
			t.Fail()
		}
	}
	{
		var arr = sx.NewArrayFrom(1, 2, 3, 2, 4)
		var arr2 = sx.RemoveAll(arr, func(x int) bool {
			return x == 5
		})
		if arr2.Length() != arr.Length() {
			t.Fail()
		}
	}
}

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
