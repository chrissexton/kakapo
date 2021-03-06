package lisp

import "fmt"

type primitive_t struct {
	name string
	f func(*scope, []sexpr) sexpr
}

func primitive(name string, f func(*scope, []sexpr) sexpr) primitive_t {
	return primitive_t{name, f}
}

// (go expr)
// Runs expr in the background.
func primitiveGo(sc *scope, ss []sexpr) sexpr {
	go func() {
		builtinEval(sc, ss)
	}()
	return Nil
}


// (if cond expr1 expr2)
func primitiveIf(sc *scope, ss []sexpr) sexpr {
	if len(ss) < 2 || len(ss) > 3 {
		panic("Invalid number of arguments to primitive if")
	}
	cond := ss[0]
	cv := eval(sc, cond)
	if IsTrue(cv) {
		return eval(sc, ss[1])
	} else if len(ss) == 3 {
		return eval(sc, ss[2])
	}
	return Nil
}

// (for cond expr)
func primitiveFor(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 2 {
		panic("Invalid number of arguments")
	}
	cond := ss[0]
	expr := ss[1]
	val := Nil
	cv := eval(sc, cond)
	for cv != nil {
		val = eval(sc, expr)
		cv = eval(sc, cond)
	}
	return val
}

// (lambda (arg1 ...) expr)
func primitiveLambda(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 2 {
		panic("Invalid number of arguments")
	}
	expr := ss[1]
	evalScopeParent := newScope(sc)
	var args_ = ss[0]
	// TODO type check the args list
	return function(func(callScope *scope, ss []sexpr) sexpr {
		args := args_
		evalScope := newScope(evalScopeParent)
		// Match args with ss
		aC, ok := args.(cons)
		for args != nil {
			if len(ss) == 0 {
				panic("Invalid number of arguments")
			}
			if !ok {
				// turn ss back into a cons
				val := unflatten(ss)
				s, k := args.(sym)
				if !k {
					panic("Invalid parameter specification")
				}
				evalScope.define(s, val)
				goto done
			}
			arg := aC.car
			val := ss[0]
			s, k := arg.(sym)
			if !k {
				panic("Invalid parameter specification")
			}
			evalScope.define(s, val)

			ss = ss[1:]
			args = aC.cdr
			aC, ok = args.(cons)
		}
		if len(ss) > 0 {
			panic("Invalid number of arguments")
		}
	done:
		return eval(evalScope, expr)
	})
}

// (let ((sym1 val1) ...) expr1 ...)
func primitiveLet(sc *scope, ss []sexpr) sexpr {
	if len(ss) < 1 {
		panic("Invalid number of arguments")
	}
	evalScope := newScope(sc)
	bindings := flatten(ss[0])
	for _, b := range bindings {
		bs := flatten(b)
		if len(bs) != 2 {
			panic("Invalid binding")
		}
		s, ok := bs[0].(sym)
		if !ok {
			panic("Invalid binding")
		}
		val := eval(sc, bs[1])
		evalScope.define(s, val)
	}

	prog := ss[1:]
	last := Nil
	for _, l := range prog {
		last = eval(evalScope, l)
	}
	return last
}

// (defmacro f (arg1 arg2 ...) body)
func primitiveDefmacro(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 3 {
		msg := fmt.Sprintf(
			"Invalid number of arguments to defmacro.  " +
			"Expected 3, got %d", len(ss))
		panic(msg)
	}
	idSym, ok := ss[0].(sym)
	if !ok {
		msg := fmt.Sprint("Expected a symbol as first argument to" +
			"defmacro, got", ss[0])
		panic(msg)
	}
	argNames := unpackSymList(ss[1])
	bodyScope := newScope(sc)
	for _, argName := range(argNames) {
		bodyScope.define(argName, argName)
	}
	body := eval(bodyScope, ss[2])
	m := macro{idSym, argNames, body}
	sc.defineHigh(idSym, m)
	return Nil
}

// (macroexpand-1 '(macroName arg1 arg2 ...))
// Runs the specified macro on the arguments, returning the resulting AST
// but not evaluating it.
func primitiveMacroexpand1(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 1 {
		msg := fmt.Sprint("Expected one argument to macroexpand-1,",
			"got", len(ss))
		panic(msg)
	}
	ss2 := flatten(eval(sc, ss[0]))
	first := eval(sc, ss2[0])
	switch m := first.(type) {
	case macro:
		return m.expand(ss2[1:len(ss2)])
	}
	msg := fmt.Sprintf("In macroexpand-1, expected a macro, got %s",
		asString(ss2[0]))
	panic(msg)
}

// unpackSymList converts an expression containing a list of symbols to a slice
// of syms.
func unpackSymList(e sexpr) []sym {
	symExprs := flatten(e)
	slice := make([]sym, len(symExprs))
	for i, e2 := range(symExprs) {
		s, ok := e2.(sym)
		if !ok {
			msg := fmt.Sprintf("Expected a symbol, got %s",
				asString(e2));
			panic(msg)
		}
		slice[i] = s
	}
	return slice
}

// (define keyword expression)
func primitiveDefine(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 2 {
		panic("Invalid number of arguments")
	}
	idSym, ok := ss[0].(sym)
	if !ok {
		panic("Invalid argument")
	}
	val := eval(sc, ss[1])
	sc.defineHigh(idSym, val)
	return Nil
}

// (quote expr)
func primitiveQuote(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 1 {
		panic("Invalid number of arguments")
	}
	return ss[0]
}

// (begin expr1 ...)
//
// This could be implemented (in large part, at least) as an ordinary function
// taking variable arguments; however, in the interest of clarity of behaviour,
// it is not.
func primitiveBegin(sc *scope, ss []sexpr) sexpr {
	last := Nil
	for _, l := range ss {
		last = eval(sc, l)
	}
	return last
}
