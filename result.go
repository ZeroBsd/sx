// SPDX-License-Identifier: 0BSD
package sx

import "errors"

func NewResultFrom[V any](value V) Result[V] {
	var opt = Result[V]{data: value, err: nil}
	return opt
}

func NewResultError[V any](errorMessageParts ...string) Result[V] {
	var opt = Result[V]{err: errors.New(StrCat(errorMessageParts...))}
	return opt
}

func NewResultFromError[V any](err error) Result[V] {
	var opt = Result[V]{err: err}
	return opt
}

type Result[T any] struct {
	data T
	err  error
}

func (r Result[T]) Ok() bool {
	return r.err == nil
}

func (r Result[T]) Value() T {
	if !r.Ok() {
		Throw(r.err.Error())
	}
	return r.data
}

func (r Result[T]) ValueOrThrow(errorMessage string) T {
	if !r.Ok() {
		ThrowIfError(errors.New(errorMessage))
	}
	return r.data
}

func (r Result[T]) ValueOr(defaultValue T) T {
	if !r.Ok() {
		return defaultValue
	}
	return r.data
}

func (r Result[T]) ValueOrInit() (value T) {
	return r.ValueOr(value)
}

func (r Result[T]) ValueIsOptional() Optional[T] {
	if !r.Ok() {
		return NewOptional[T]()
	}
	return NewOptionalFrom(r.Value())
}

func (r Result[T]) Error() string {
	if r.Ok() {
		Throw("Fatal error: accessing empty error value from result")
	}
	return r.err.Error()
}
