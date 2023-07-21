// SPDX-License-Identifier: 0BSD
package sx_test

import (
	"errors"
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestResultFrom(t *testing.T) {
	var r = sx.NewResultFrom(2)
	var i int = r.Value()
	sx.ThrowIf(i != 2)
}

func TestErrorResultWithDefault(t *testing.T) {
	var r = sx.NewResultError[int]("failed")
	sx.ThrowIf(r.Error() != "failed")
	var i int = r.ValueOr(10)
	sx.ThrowIf(i != 10)
}

func TestErrorResultWithInit(t *testing.T) {
	var r = sx.NewResultError[int]("failed")
	sx.ThrowIf(r.Error() != "failed")
	var i int = r.ValueOrInit()
	sx.ThrowIf(i != 0)
}

func TestErrorResult(t *testing.T) {
	defer sx.Catch(func(err error) {})
	var r = sx.NewResultError[int]("failed")
	sx.ThrowIf(r.Error() != "failed")
	var _ int = r.Value()
	t.FailNow()
}

func TestErrorResultFromError(t *testing.T) {
	defer sx.Catch(func(err error) {})
	var r = sx.NewResultFromError[int](errors.New("failed"))
	sx.ThrowIf(r.Error() != "failed")
	var _ int = r.Value()
	t.FailNow()
}

func TestAccessingErrorOfValidResult(t *testing.T) {
	defer sx.Catch(func(err error) {})
	var r = sx.NewResultFrom(2)
	var _ = r.Error()
	t.FailNow()
}

func TestValueOrThrow(t *testing.T) {
	var r = sx.NewResultFrom(2)
	r.ValueOrThrow("failed")
}

func TestFailedValueOrThrow(t *testing.T) {
	defer sx.Catch(func(err error) {
		if err.Error() != "failed" {
			t.FailNow()
		}
	})
	var r = sx.NewResultError[int]("error")
	r.ValueOrThrow("failed")
	t.FailNow()
}

func TestValueisOptional(t *testing.T) {
	{
		var r = sx.NewResultFrom(2)
		var opt = r.ValueIsOptional()
		if opt.IsEmpty() {
			t.FailNow()
		}
	}
	{
		var r = sx.NewResultError[int]("error")
		var opt = r.ValueIsOptional()
		if !opt.IsEmpty() {
			t.FailNow()
		}
	}
}
