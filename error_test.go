// SPDX-License-Identifier: 0BSD
package sx_test

import (
	"errors"
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestError(t *testing.T) {
	defer sx.Catch(func(err error) {
		if err.Error() != "this will be catched" {
			t.FailNow()
		}
	})
	sx.Throw("this will be catched")
}

func TestThrowIfTrue(t *testing.T) {
	defer sx.Catch(func(err error) {
		if err.Error() != "this will be thrown and catched" {
			t.FailNow()
		}
	})
	sx.ThrowIf(2 == 2, "this will be thrown and catched")
}

func TestThrowIfFalse(t *testing.T) {
	sx.ThrowIf(1 == 2, "This will never be thrown")
}

func TestThrowIfErrorTrue(t *testing.T) {
	var err = errors.New("some error")
	defer sx.Catch(func(err error) {
		if err.Error() != "some error" {
			t.FailNow()
		}
	})
	sx.ThrowIfError(err)
}

func TestThrowIfErrorFalse(t *testing.T) {
	var err error
	sx.ThrowIfError(err)
}

func TestCatch(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			var i = sx.TypeCast[int](x)
			if i.Value() != 2 {
				t.FailNow()
			}
		}
	}()
	defer sx.Catch(func(err error) {
		t.FailNow()
	})
	panic(2)
}
