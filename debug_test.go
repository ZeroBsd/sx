// SPDX-License-Identifier: 0BSD
package sx_test

import (
	"errors"
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestDebug(t *testing.T) {
	sx.ThrowIf(sx.DebugIsActive())
	sx.Debug("Hello Debug!")
	sx.DebugAny(errors.New("Hello Debug!"))
	sx.ThrowIf(sx.DebugIsActive())
	sx.DebugSetOn()
	sx.Debug("")
	sx.DebugAny("")
	sx.ThrowIf(!sx.DebugIsActive())
	sx.DebugPushOn()
	sx.ThrowIf(!sx.DebugIsActive())
	sx.DebugPop()
	sx.ThrowIf(!sx.DebugIsActive())
	sx.DebugPushOff()
	sx.ThrowIf(sx.DebugIsActive())
	sx.DebugPop()
	sx.ThrowIf(!sx.DebugIsActive())
	sx.DebugSetOff()
	sx.ThrowIf(sx.DebugIsActive())
	sx.ThrowIf(sx.DebugIsActive())
}
