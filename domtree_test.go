// SPDX-License-Identifier: 0BSD
package sx_test

import (
	"strings"
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestDomDefaultVisitor(t *testing.T) {
	var v = sx.NewDomVisitor[any](nil)
	var x = 42
	v.VisitRecursively(x)
	v.AfterVisit = func(context, domTreeNode any) {
		if domTreeNode.(int) != 42 {
			t.FailNow()
		}
	}
}

func NewIntVisitor() sx.DomVisitor[*sx.StringBuilder] {
	var sb = sx.NewStringBuilder()
	var visitor = sx.NewDomVisitor(sb)
	visitor.BeforeVisit = func(context *sx.StringBuilder, domTreeNode any) bool {
		context.WriteString("{")
		if asInt := sx.TypeCast[int](domTreeNode); asInt.Ok() {
			context.WriteAny(asInt.Value())
		}
		return true
	}
	visitor.AfterVisit = func(context *sx.StringBuilder, domTreeNode any) {
		context.WriteString("}")
	}
	return visitor
}

func TestDomTreeArray(t *testing.T) {
	var v = NewIntVisitor()
	var arr = [3]int{1, 2, 3}
	v.VisitRecursively(arr)
	if r := v.Context.String(); r != "{{1}{2}{3}}" {
		t.FailNow()
	}
}

func TestDomTreeSlice(t *testing.T) {
	var v = NewIntVisitor()
	var slice = []int{1, 2, 3}
	v.VisitRecursively(slice)
	if r := v.Context.String(); r != "{{1}{2}{3}}" {
		t.FailNow()
	}
}

func TestDomTreeMap(t *testing.T) {
	var v = NewIntVisitor()
	var m = map[string]int{"a": 42, "b": 43}
	v.VisitRecursively(m)
	// order is random, so we need to check both possibilities
	if r := v.Context.String(); r != "{{42}{43}}" && r != "{{43}{42}}" {
		t.FailNow()
	}
}

func TestDomTreeStruct(t *testing.T) {
	var v = NewIntVisitor()
	var s = struct {
		X int
		Y int
	}{X: 42, Y: 43}
	v.VisitRecursively(s)
	if r := v.Context.String(); r != "{{42}{43}}" {
		t.FailNow()
	}
}

func TestDomTreePtr(t *testing.T) {
	var v = NewIntVisitor()
	var s = &struct {
		X int
		Y int
	}{X: 42, Y: 43}
	v.VisitRecursively(s)
	if r := v.Context.String(); r != "{{{42}{43}}}" {
		t.FailNow()
	}
}

func TestDomTreeSxArray(t *testing.T) {
	var v = NewIntVisitor()
	var s = sx.NewArrayFrom(42, 43)
	v.VisitRecursively(s)
	if r := v.Context.String(); r != "{{{42}{43}}}" {
		t.FailNow()
	}
}

func TestDomTreeSxMap(t *testing.T) {
	var v = NewIntVisitor()
	var s = sx.NewMapFrom(map[string]int{"a": 42, "b": 43})
	v.VisitRecursively(s)
	if r := v.Context.String(); r != "{{{{42}{43}}}}" && r != "{{{{43}{42}}}}" {
		t.FailNow()
	}
}

func TestDomTreeAst(t *testing.T) {
	var v = NewIntVisitor()
	v.BeforeVisit = func(context *sx.StringBuilder, domTreeNode any) bool {
		context.WriteString("{")
		switch node := domTreeNode.(type) {
		case int:
			context.WriteAny(domTreeNode)
		case string:
			if node == "XXX" {
				context.WriteAny("ZZZ")
				context.WriteAny("}") // after visit will not be executed
				return false
			} else {
				context.WriteAny(domTreeNode)
			}
		}
		return true
	}
	var ast = struct {
		Expression any
	}{
		Expression: struct {
			Operator []string
			LHS      any
			RHS      any
		}{
			Operator: []string{">", "="},
			LHS: struct {
				Identifier string
			}{
				Identifier: "a",
			},
			RHS: struct {
				Operator string
				Value1   any
				Value2   any
				Value3   any
				Value4   any
				ValueX   any
			}{
				Operator: "+",
				Value1: &struct {
					Value int
				}{
					Value: 42,
				},
				Value2: struct {
					Value int
				}{
					Value: 43,
				},
				Value3: sx.NewArrayFrom(44, 45),
				Value4: sx.NewHashMapFrom(map[string]string{"a": "aa", "b": "bb"}),
				ValueX: "XXX",
			},
		},
	}
	v.VisitRecursively(ast)
	var r = v.Context.String()
	if r != "{{{{>}{=}}{{a}}{{+}{{{42}}}{{43}}{{{44}{45}}}{{{{aa}{bb}}}}{ZZZ}}}}" {
		t.FailNow()
	}
	if openCurlyCount, closeCurlyCount := strings.Count(r, "{"), strings.Count(r, "}"); openCurlyCount != closeCurlyCount {
		t.FailNow()
	}
}
