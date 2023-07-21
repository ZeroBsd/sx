// SPDX-License-Identifier: 0BSD
package sx

var _ Container = &Optional[int]{}
var _ Container = NewOptional[int]()
var _ Container = NewOptionalFrom(1)
var _ Iterable[int, int] = &Optional[int]{}
var _ Iterable[int, int] = NewOptional[int]()
var _ Iterable[int, int] = NewOptionalFrom(1)

func NewOptional[V any]() Optional[V] {
	var opt = Optional[V]{
		valid: false,
	}
	return opt
}

func NewOptionalFrom[V any](value V) Optional[V] {
	var opt = Optional[V]{
		data:  value,
		valid: true,
	}
	return opt
}

type Optional[T any] struct {
	data  T
	valid bool
}

func (opt Optional[T]) IsEmpty() bool {
	return !opt.valid
}

func (opt Optional[T]) Length() int {
	if opt.IsEmpty() {
		return 0
	}
	return 1
}

func (opt Optional[T]) Value() T {
	if opt.IsEmpty() {
		Throw("Fatal error: accessing empty optional value")
	}
	return opt.data
}

func (opt Optional[T]) ValueOr(defaultValue T) T {
	if !opt.IsEmpty() {
		return opt.data
	}
	return defaultValue
}

func (opt Optional[T]) ValueOrDefault() (value T) {
	return opt.ValueOr(value)
}

func (opt Optional[T]) NewIterator() Iterator[int, T] {
	if !opt.IsEmpty() {
		return NewArrayFrom(opt.data).NewIterator()
	}
	return NewArray[T]().NewIterator()
}
