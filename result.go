// SPDX-License-Identifier: 0BSD
package sx

import "errors"

var _ Container = &Result[int]{}
var _ Container = NewResultFrom(1)
var _ Iterable[int, int] = &Result[int]{}
var _ Iterable[int, int] = NewResultFrom(1)

func NewResultFrom[V any](value V) Result[V] {
	var opt = Result[V]{
		data: value,
		err:  nil,
	}
	return opt
}

func NewResultError[V any](errorMessage string) Result[V] {
	var opt = Result[V]{
		err: errors.New(errorMessage),
	}
	return opt
}

type Result[T any] struct {
	data T
	err  error
}

func (r Result[T]) Ok() bool {
	return r.err == nil
}

func (r Result[T]) IsEmpty() bool {
	return !r.Ok()
}

func (r Result[T]) Length() int {
	if r.IsEmpty() {
		return 0
	}
	return 1
}

func (r Result[T]) Value() T {
	if r.IsEmpty() {
		Throw(r.err.Error())
	}
	return r.data
}

func (r Result[T]) ValueOr(defaultValue T) T {
	if r.IsEmpty() {
		return defaultValue
	}
	return r.data
}

func (r Result[T]) ValueOrDefault() (value T) {
	return r.ValueOr(value)
}

func (r Result[T]) Error() string {
	if r.Ok() {
		Throw("Fatal error: accessing empty error value from result")
	}
	return ""
}

func (r Result[T]) NewIterator() Iterator[int, T] {
	if r.IsEmpty() {
		return NewArray[T]().NewIterator()
	} else {
		return NewArrayFrom(r.data).NewIterator()
	}
}
