package lisp

import (
	"fmt"
)

type macro struct {
	name sym
	argNames []sym
	body sexpr
}

func (m macro) expand(args []sexpr) sexpr {
	if len(args) != len(m.argNames) {
		pattern := ("Wrong number of arguments to %s macro. " +
			"Expected %i args, got %i")
		msg := fmt.Sprintf(pattern, len(args), len(m.argNames))
		panic(msg)
	}
	newBody := m.body
	for i, argName := range(m.argNames) {
		arg := args[i]
		newBody = replaceSym(argName, arg, newBody)
	}
	return newBody
}

// replaceSym(s, val, e) replaces s with val in e.
func replaceSym(s sym, val sexpr, e sexpr) sexpr {
	if e == s {
		return val
	}
	switch e2 := e.(type) {
	case cons:
		newCar := replaceSym(s, val, e2.car)
		newCdr := replaceSym(s, val, e2.cdr)
		return cons{newCar, newCdr}
	}
	return e
}

