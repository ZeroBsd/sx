// SPDX-License-Identifier: 0BSD
package sx

import (
	"fmt"
	"os"
)

var debugStack = NewArrayFrom(false)

func DebugIsActive() bool { return debugStack.Peek().ValueOr(false) == true }

func DebugSetOn()   { debugStack.Pop(); debugStack.Push(true) }
func DebugSetOff()  { debugStack.Pop(); debugStack.Push(false) }
func DebugPushOn()  { debugStack.Push(true) }
func DebugPushOff() { debugStack.Push(false) }
func DebugPop() {
	if debugStack.Length() > 1 {
		debugStack.Pop()
	}
}

// Print debug messages
// will not print by default, please call "SetDebugOn" first
func Debug(values ...string) {
	if DebugIsActive() {
		fmt.Fprintln(os.Stderr, StrCat(values...))
	}
}

// Print debug messages
// will not print by default, please call "SetDebugOn" first
func DebugAny(values ...any) {
	if DebugIsActive() {
		PrintAnyStdErr(values...)
	}
}
