package lisp

import (
	"fmt"
	"strings"
)

type cons struct {
	car sexpr
	cdr sexpr
}

type sexpr interface{}
type atom interface{}
type sym string
type function func(*scope, []sexpr) sexpr

type native interface{}

func (v cons) String() string {
	return fmt.Sprintf("(%s . %s)", asString(v.car), asString(v.cdr))
}

func asString(v sexpr) string {
	switch v := v.(type) {
	case cons:
		if isList(v) {
			items := flatten(v)
			strs := make([]string, len(items))
			for i, x := range(items) {
				strs[i] = asString(x)
			}
			return "(" + strings.Join(strs, " ") + ")"
		}
		return v.String()
	case sym:
		return string(v)
	case float64:
		return fmt.Sprintf("%G", v)
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case function:
		return "<func>"
	case primitive_t:
		return fmt.Sprintf("<primitive: %s>", v.name)
	case macro:
		return fmt.Sprintf("<macro: %s>", v.name)
	case nil:
		return "nil"
	case native:
		return "<native>"
	}
	return "<unknown>"
}

func isFunction(s sexpr) bool {
	_, ok := s.(function)
	return ok
}

func isPrimitive(s sexpr) bool {
	_, ok := s.(primitive_t)
	return ok
}

func isList(s sexpr) bool {
	if s == nil {
		return true
	}
	switch c := s.(type) {
	case cons:
		return isList(c.cdr)
	}
	return false;
}

