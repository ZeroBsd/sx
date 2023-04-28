// SPDX-License-Identifier: 0BSD
package sx_test

import (
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestExamples(t *testing.T) {
	base()
	arraysMaps()
	if playThrowAndCatch(true).Ok() {
		t.FailNow()
	}
	if x := playThrowAndCatch(false).Value(); x != 42 {
		t.FailNow()
	}
}

func base() {
	sx.PrintLn("Hello, ", "World")
}

func arraysMaps() {
	var array1 = sx.NewArray[int]()       // create dynamic array using a generic value type
	var array2 = sx.NewArrayFrom(1, 2, 3) // same as above, but infer the value type from the values
	array1.PushAll(array2)                // add all items from array2 to array1
	array1.Push(array2.Slice()...)        // same
	array1.Push(4)                        // add another value

	// iterate through all key/value pairs. This works for every array and map
	for it := array1.NewIterator(); it.Ok(); it.Next() {
		sx.PrintAnyLn("Key: '", it.Key(), "' ; Value: '", it.Value(), "'")
	}

	// arrays are maps - they satisfy the 'Map' interface
	var myMap sx.Map[int, int] = array1

	// replace the entry at index 0
	// since this is actually an array, the index must exist
	// for real maps, this restriction is not present
	if myMap.Has(0) {
		myMap.Put(0, 42)
	}
	sx.ThrowIf(myMap.Get(0).Value() != 42, "Oh no, the value is not 42!")

	// the iterator is part of the 'Map' interface
	// when iterating over a real map, the order is undefined. For arrays, iteration is ordered
	for it := myMap.NewIterator(); it.Ok(); it.Next() {
		sx.PrintAnyLn("Key: '", it.Key(), "' ; Value: '", it.Value(), "'")
	}
}

// this function returns a result which contains either the value (of type int) or an error
func playThrowAndCatch(throws bool) (result sx.Result[int]) {

	// provide an error handler in case something goes wrong
	// !!! Important: don't forget the 'defer' !!!
	defer sx.Catch(func(err error) {

		// we print the exception message here
		// this may be ok, since we do not re-throw (log-and-throw-antipattern)
		sx.PrintStdErrLn(err.Error())

		// we transform the exception (panic) to a valid result (that contains the error)
		// be aware that we cannot use 'return' here, so we make use of the named return value
		result = sx.NewResultFromError[int](err)
	})

	// throws, if the input variable (bool) was true. The text will be encapsulated in an error
	sx.ThrowIf(throws, "This is throws as an error if throws is true. Internally, this panics with an error")

	// if the above didn't throw, we return a valid number. In this case, the type argument is inferred
	return sx.NewResultFrom(42)
}
