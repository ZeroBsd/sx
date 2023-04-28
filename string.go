// SPDX-License-Identifier: 0BSD
package sx

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Print(values ...string) {
	fmt.Fprint(os.Stdout, StrCat(values...))
}

func PrintLn(values ...string) {
	fmt.Fprintln(os.Stdout, StrCat(values...))
}

func PrintStdErr(values ...string) {
	fmt.Fprint(os.Stderr, StrCat(values...))
}

func PrintStdErrLn(values ...string) {
	fmt.Fprintln(os.Stderr, StrCat(values...))
}

func PrintAny(values ...any) {
	fmt.Fprint(os.Stdout, Str(values...))
}

func PrintAnyLn(values ...any) {
	fmt.Fprintln(os.Stdout, Str(values...))
}

func PrintAnyStdErr(values ...any) {
	fmt.Fprint(os.Stderr, Str(values...))
}

func PrintAnyStdErrLn(values ...any) {
	fmt.Fprintln(os.Stderr, Str(values...))
}

func String2Int(s string) Result[int64] {
	var val, e = strconv.ParseInt(s, 10, 64)
	if e != nil {
		NewResultError[int64]("Fatal error: cannot convert string to int64")
	}
	return NewResultFrom(val)
}

func String2Float(s string) Result[float64] {
	var val, e = strconv.ParseFloat(s, 64)
	if e != nil {
		NewResultError[int64]("Fatal error: cannot convert string to float64")
	}
	return NewResultFrom(val)
}

func String2Bool(s string) Result[bool] {
	switch s {
	default:
		return NewResultError[bool]("Fatal error: cannot convert string to bool")
	case "true":
		return NewResultFrom(true)
	case "false":
		return NewResultFrom(false)
	}
}

// Concatenates strings without adding any extra (white-)spaces
func StrCat(stringsToJoin ...string) string {
	return StrJoin("", stringsToJoin...)
}

// Joins strings using a separator, but without adding any extra (white-)spaces
func StrJoin(separator string, stringsToJoin ...string) string {
	return strings.Join(stringsToJoin, separator)
}

// Concatenates values without adding any extra (white-)spaces
func Str(values ...any) string {
	switch len(values) {
	case 0:
		return ""
	case 1:
		return fmt.Sprint(values[0])
	default:
		return NewStringBuilder().WriteAny(values...).String()
	}
}

type StringBuilder struct {
	strings.Builder
}

func NewStringBuilder() *StringBuilder {
	return &StringBuilder{strings.Builder{}}
}

func (sb *StringBuilder) WriteStrings(values ...string) *StringBuilder {
	sb.WriteString(StrCat(values...))
	return sb
}

func (sb *StringBuilder) WriteAny(values ...any) *StringBuilder {
	for _, v := range values {
		sb.WriteString(fmt.Sprint(v))
	}
	return sb
}
