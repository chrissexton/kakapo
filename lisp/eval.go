package lisp

// eval evaluates an s-expression, including syntax transformations (macros).
func eval(sc *scope, e sexpr) sexpr {
	switch e := e.(type) {
	case cons: // a function, primitive or macro to evaluate
		cons := e
		car := eval(sc, cons.car)
		cdr := cons.cdr
		args := flatten(cdr)
		switch f := car.(type) {
		case function:
			// Evaluate all arguments
			for i, a := range args {
				args[i] = eval(sc, a)
			}
			return f(sc, args)

		case primitive_t:
			// Run without first evaluating the arguments.
			return f.f(sc, args)

		case macro:
			// Expand the macro invocation and then evaluate the result.
			builtinPrint(sc, []sexpr{f})
			return eval(sc, f.expand(args))

		default:
			msg := ("Attempted application on something other " +
				"than a function, primitive or macro")
			panic(msg)
		}
	case sym:
		return sc.lookup(e)
	}
	return e
}

func apply(sc *scope, e sexpr, ss []sexpr) sexpr {
	f, ok := e.(function)
	if !ok {
		panic("Attempted application on non-function")
	}
	return f(sc, ss)
}

func unflatten(ss []sexpr) sexpr {
	c := sexpr(nil)
	for i := len(ss) - 1; i >= 0; i-- {
		c = cons{ss[i], c}
	}
	return c
}

func flatten(s sexpr) (ss []sexpr) {
	_, ok := s.(cons)
	for ok {
		ss = append(ss, s.(cons).car)
		s = s.(cons).cdr
		_, ok = s.(cons)
	}
	if s != nil {
		panic("List isn't flat")
	}
	return 
}
