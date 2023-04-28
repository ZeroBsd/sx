# sx - Simple Extensions for the Go Programming Language
![](https://img.shields.io/badge/License-0BSD-brightgreen "License: 0BSD") ![](https://img.shields.io/badge/StaticCheck-0%20Warnings-brightgreen "StaticCheck: 0 Warnings") [![Go Report Card](https://goreportcard.com/badge/github.com/ZeroBsd/sx)](https://goreportcard.com/report/github.com/ZeroBsd/sx) ![](https://img.shields.io/badge/TestCoverage-100%25-brightgreen "TestCoverage: 100%")

---
License
===

Everything in this repository is provided under the [__0BSD__](https://github.com/ZeroBsd/sx/blob/main/LICENSE) License (see: [spdx.org](https://spdx.org/licenses/0BSD.html)), a public domain equivalent license\
(when in doubt, please contact your lawyer)


---
$~$


About
===
>Unfortunately, the Go standard library does not provide generic types.\
So I wrote this library to make my life easier - and maybe yours, too.\
I hope you like it 🙂

$~$

__Features:__
* Standard Containers (Array, HashMap, Optional)
* Iterators (extensible)
* Result Type (a bit 'rusticious' - similar to std::expected)
* Exceptions (based on panic/recover)
* ... and some other small helpers

$~$

__Quality:__
* high quality library with 100% test coverage
* dependency free and non-invasive
* minimalistic, but 'all-in-one' solution
* written in pure Go
* _0BSD_ licence, so you don't have to give credit or carry around copyright notices
* No legacy stuff - you need an up-to-date Go version (1.20+)


$~$

---
$~$

Usage / Examples
===

$~$

### Basic Usage
```go
// tldr; just use the library as a module (legacy non-module usage is not supported)
// 1. include the header for the library
// 2. update/check your go.mod file
// 3. run "go mod tidy" to update your go.sum file
// 4. check your go.sum file for the correct version, just in case

// do your package thing
package main

// import the library
import (
	"github.com/ZeroBsd/sx"
)

// enjoy
func Main() {
	sx.PrintLn("Hello, ", "World")
}
```

### Arrays / Maps
```go
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
```

### Error Handling and Exceptions
```go
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
```
