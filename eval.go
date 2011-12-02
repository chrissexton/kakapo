package main

import (
	"fmt"
)

type scope map[string]sexpr

func (v cons) String() string {
	return "<cons>"
}

func doEval(c chan sexpr) {
	for e := range c {
		v := eval(e)
		switch v.(type) {
		case cons:
			fmt.Printf("<cons>\n")
		case sym:
			fmt.Printf("<sym : %s>\n", string(v.(sym)))
		case float64:
			fmt.Printf("%f\n", v.(float64))
		case string:
			fmt.Printf("\"%s\"\n", v.(string))
		default:
			fmt.Printf("nil\n")
		}
	}
}

func isFunction(s sexpr) bool {
	_, ok := s.(func ([]sexpr) sexpr)
	return ok
}

func isSyntax(s sexpr) bool {
	return false
}

// Perform appropriate syntax transformations on the given s-expression. Note
// that some s-expressions that 'should' involve syntax transformations, such
// as (if cond x y) and (lambda ...), don't - they just aren't evaluated as
// normal functions. (TODO make user-defined transformations more flexible to
// add symmetry.)
func transform(e sexpr) sexpr {
	return e // TODO
}

// Evaluates an s-expression, excluding syntax transformations (macros).
func eval(e sexpr) sexpr {
	// TODO match against syntax declarations
	switch e.(type) {
	case cons: // a function to evaluate
		cons := e.(cons)
		car := eval(cons.car)
		cdr := cons.cdr
		if !isFunction(car) {
			panic("Attempted application on non-function")
		}
		args := flatten(cdr)
		f := car.(func([]sexpr) sexpr)
		// This is a function - evaluate all arguments
		for i, a := range args {
			args[i] = eval(a)
		}
		return f(args)
	case sym:
		return lookup(string(e.(sym)))
	case float64:
		return e
	case string:
		return e
	}
	return Nil
}

func flatten(s sexpr) (ss []sexpr) {
	_, ok := s.(cons)
	for ok {
		ss = append(ss, s.(cons).car)
		s = s.(cons).cdr
		_, ok = s.(cons)
	}
	// TODO what if s isn't nil now?
	return
}

// Performs lookup of symbols for values.
func lookup(sym string) sexpr {
	// TODO attempt lookup in reflect

	v, ok := global[sym]
	if ok {
		return v
	}
	panic("undefined")
}
