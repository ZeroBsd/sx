// SPDX-License-Identifier: 0BSD
package sx_test

import (
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestPair(t *testing.T) {
	var p1 = sx.NewPair(3, 0xAA55)
	var p2 = sx.NewPair(3, 0b1010101001010101)
	sx.ThrowIf(p1 != p2)
	sx.ThrowIf(p1.Key != 3)
	sx.ThrowIf(p1.Value != p2.Value)
}
