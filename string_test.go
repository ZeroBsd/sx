package sx_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestPrint(t *testing.T) {
	sx.Print("Hello World!")
	sx.PrintLn("Hello World!")
	sx.PrintAny("Hello World!")
	sx.PrintAny(errors.New("Failed World!"))
	sx.PrintAnyLn("Hello World!")
	sx.PrintAnyLn(errors.New("Failed World!"))
	sx.PrintStdErr("Hello stderr!")
	sx.PrintStdErrLn("Hello stderr!")
	sx.PrintAnyStdErr("Hello stderr!")
	sx.PrintAnyStdErr(errors.New("Hello stderr!"))
	sx.PrintAnyStdErrLn("Hello stderr!")
	sx.PrintAnyStdErrLn(errors.New("Hello stderr!"))
	sx.PrintLn("")
	sx.PrintAny('\n')
}

func TestStringConcats(t *testing.T) {
	if sx.Str() != "" {
		t.FailNow()
	}
	if sx.Str("a", "b", "c") != "abc" {
		t.FailNow()
	}
	if sx.Str("a", "b", 'c') == "abc" {
		t.FailNow()
	}
	if sx.Str("a", "b", 'c', 42) != fmt.Sprint("a")+fmt.Sprint("b")+fmt.Sprint('c')+fmt.Sprint(42) {
		t.FailNow()
	}
	if sx.StrCat("a", "b", "c") != "abc" {
		t.FailNow()
	}
	if sx.StrJoin(",", "a", "b", "c") != "a,b,c" {
		t.FailNow()
	}
	if sx.StrJoin("-") != "" {
		t.Fail()
	}
}

func TestStringConversions(t *testing.T) {
	if sx.String2Bool("true").Value() != true {
		t.FailNow()
	}
	if sx.String2Bool("false").Value() != false {
		t.FailNow()
	}
	if sx.String2Bool("truthy").ValueOrInit() != false {
		t.FailNow()
	}
	if sx.String2Int("42").Value() != 42 {
		t.FailNow()
	}
	if sx.String2Int("A38").ValueOrInit() != 0 {
		t.FailNow()
	}
	if sx.String2Float("3.1415").Value() != 3.1415 {
		t.FailNow()
	}
	if sx.String2Float("3.14.15").ValueOrInit() != 0 {
		t.FailNow()
	}
}

func TestStringBuilder(t *testing.T) {
	var s = sx.NewStringBuilder()
	var r2d2 = s.WriteAny("R", 2, string('D'), uint64(2)).String()
	if r2d2 != "R2D2" {
		t.FailNow()
	}
	s = sx.NewStringBuilder()
	var c3po = s.WriteStrings("C", sx.Str(3), "P", string('O')).String()
	if c3po != "C3PO" {
		t.FailNow()
	}
}
