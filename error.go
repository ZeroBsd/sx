// SPDX-License-Identifier: 0BSD
package sx

import (
	"errors"
	"strings"
)

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

func ThrowIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// Catches errors thrown by "Throw" or "ThrowIf", recovers the panic
// Needs to be called before throwing errors
// Only handles errors. If you throw/panic something else, you need to provide a custom recover function
func Catch(handler ...func(error)) {
	if x := recover(); x != nil {
		var err = TypeCast[error](x)
		if len(handler) > 0 && err.Ok() {
			for _, h := range handler {
				h(err.Value())
			}
		} else {
			panic(x)
		}
	}
}
