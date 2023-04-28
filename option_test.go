// SPDX-License-Identifier: 0BSD
package sx_test

import (
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestValidOptional(t *testing.T) {
	var opt = sx.NewOptionalFrom(1)
	if opt.IsEmpty() {
		t.FailNow()
	}
	if opt.Length() != 1 {
		t.FailNow()
	}
	if opt.Value() != 1 {
		t.FailNow()
	}
	if opt.ValueOrDefault() != 1 {
		t.FailNow()
	}
	if opt.ValueOr(42) != 1 {
		t.FailNow()
	}
	var i int
	for it := opt.NewIterator(); it.Ok(); it.Next() {
		i = it.Value()
	}
	if i != 1 {
		t.FailNow()
	}
}

func TestInvalidOptional(t *testing.T) {
	var opt = sx.NewOptional[int]()
	if !opt.IsEmpty() {
		t.FailNow()
	}
	if opt.Length() != 0 {
		t.FailNow()
	}
	if opt.ValueOrDefault() != 0 {
		t.FailNow()
	}
	if opt.ValueOr(42) != 42 {
		t.FailNow()
	}
	for it := opt.NewIterator(); it.Ok(); it.Next() {
		t.FailNow()
		break
	}
}

func TestFailingInvalidOptional(t *testing.T) {
	defer sx.Catch(func(err error) {})
	var opt = sx.NewOptional[int]()
	if !opt.IsEmpty() {
		t.FailNow()
	}
	if opt.Length() != 0 {
		t.FailNow()
	}
	if opt.ValueOrDefault() != 0 {
		t.FailNow()
	}
	if opt.ValueOr(42) != 42 {
		t.FailNow()
	}
	opt.Value() //this must throw
	t.FailNow()
}
