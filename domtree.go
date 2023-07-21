// SPDX-License-Identifier: 0BSD
package sx

import (
	"reflect"
)

type DomVisitor[Context any] struct {
	Context     Context
	BeforeVisit func(context Context, domTreeNode any) (continueVisit bool)
	AfterVisit  func(context Context, domTreeNode any)
}

func NewDomVisitor[Context any](ctx Context) (visitor DomVisitor[Context]) {
	visitor.Context = ctx
	visitor.BeforeVisit = func(ctx Context, domTreeNode any) (continueVisit bool) { return true }
	visitor.AfterVisit = func(ctx Context, domTreeNode any) {}
	return visitor
}

func (visitor DomVisitor[Context]) VisitRecursively(domTreeNode any) {
	var continueVisit = visitor.BeforeVisit(visitor.Context, domTreeNode)
	if !continueVisit {
		return
	}
	var nodeType = reflect.TypeOf(domTreeNode)
	var nodeValue = reflect.ValueOf(domTreeNode)
	switch nodeType.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for i := 0; i < nodeValue.Len(); i++ {
			var iv = reflect.Value.Index(nodeValue, i)
			if iv.CanInterface() {
				visitor.VisitRecursively(iv.Interface())
			}
		}
	case reflect.Map:
		keys := nodeValue.MapKeys()
		for _, key := range keys {
			val := nodeValue.MapIndex(key)
			if val.CanInterface() {
				visitor.VisitRecursively(val.Interface())
			}
		}
	case reflect.Struct:
		for i := 0; i < nodeType.NumField(); i++ {
			var v = nodeValue.Field(i)
			var name = nodeType.Field(i).Name
			_ = name
			if v.CanInterface() {
				visitor.VisitRecursively(v.Interface())
			}
		}
	case reflect.Pointer:
		var dref = reflect.ValueOf(domTreeNode).Elem()
		if dref.CanInterface() {
			visitor.VisitRecursively(dref.Interface())
		}
	}
	visitor.AfterVisit(visitor.Context, domTreeNode)
}
