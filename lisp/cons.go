package lisp

func builtinCons(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 2 {
		panic("Invalid number of arguments")
	}
	return cons{ss[0], ss[1]}
}

func builtinCar(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 1 {
		panic("Invalid number of arguments")
	}
	if _, ok := ss[0].(cons); !ok {
		panic("Invalid argument")
	}
	return ss[0].(cons).car
}

func builtinCdr(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 1 {
		panic("Invalid number of arguments")
	}
	if _, ok := ss[0].(cons); !ok {
		panic("Invalid argument")
	}
	return ss[0].(cons).cdr
}

func builtinList(sc *scope, ss []sexpr) sexpr {
	if len(ss) == 0 {
		return Nil
	}
	return cons{ss[0], builtinList(sc, ss[1:len(ss)])}
}

