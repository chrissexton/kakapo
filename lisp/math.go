package lisp

func builtinAdd(sc *scope, ss []sexpr) sexpr {
	// add all numeric arguments
	r := 0.
	for _, s := range ss {
		n, ok := s.(float64)
		if !ok {
			panic("Invalid argument")
		}
		r += n
	}
	return r
}

func builtinSub(sc *scope, ss []sexpr) sexpr {
	if len(ss) == 0 {
		return 0
	}
	s := ss[0]
	r, ok := s.(float64)
	if !ok {
		panic("Invalid argument")
	}
	if len(ss) == 1 {
		return -r
	}
	for _, s := range ss[1:] {
		n, ok := s.(float64)
		if !ok {
			panic("Invalid argument")
		}
		r -= n
	}
	return r
}

func builtinMul(sc *scope, ss []sexpr) sexpr {
	// add all numeric arguments
	r := 1.
	for _, s := range ss {
		n, ok := s.(float64)
		if !ok {
			panic("Invalid argument")
		}
		r *= n
	}
	return r
}

func builtinDiv(sc *scope, ss []sexpr) sexpr {
	if len(ss) == 0 {
		return 0
	}
	s := ss[0]
	r, ok := s.(float64)
	if !ok {
		panic("Invalid argument")
	}
	if len(ss) == 1 {
		return 1 / r
	}
	for _, s := range ss[1:] {
		n, ok := s.(float64)
		if !ok {
			panic("Invalid argument")
		}
		r /= n
	}
	return r
}

func builtinMod(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 2 {
		panic("Invalid number of arguments")
	}
	a, ok1 := ss[0].(float64)
	b, ok2 := ss[1].(float64)
	if !ok1 || !ok2 {
		panic("Invalid argument")
	}
	return int(a) % int(b) // TODO fixme to work with floats
}

func builtinGt(sc *scope, ss []sexpr) sexpr {
	if len(ss) == 0 {
		return 1.0
	}
	r := true
	f, ok := ss[0].(float64)
	if !ok {
		panic("Invalid argument")
	}
	for _, s := range ss[1:] {
		n, ok := s.(float64)
		if !ok {
			panic("Invalid argument")
		}
		r = r && (f > n)
		f = n
	}
	if r {
		return 1.0
	}
	return Nil
}

func builtinLt(sc *scope, ss []sexpr) sexpr {
	if len(ss) == 0 {
		return 1.0
	}
	r := true
	f, ok := ss[0].(float64)
	if !ok {
		panic("Invalid argument")
	}
	for _, s := range ss[1:] {
		n, ok := s.(float64)
		if !ok {
			panic("Invalid argument")
		}
		r = r && (f < n)
		f = n
	}
	if r {
		return 1.0
	}
	return Nil
}

func builtinGe(sc *scope, ss []sexpr) sexpr {
	if len(ss) == 0 {
		return 1.0
	}
	r := true
	f, ok := ss[0].(float64)
	if !ok {
		panic("Invalid argument")
	}
	for _, s := range ss[1:] {
		n, ok := s.(float64)
		if !ok {
			panic("Invalid argument")
		}
		r = r && (f >= n)
		f = n
	}
	if r {
		return 1.0
	}
	return Nil
}

func builtinLe(sc *scope, ss []sexpr) sexpr {
	if len(ss) == 0 {
		return 1.0
	}
	r := true
	f, ok := ss[0].(float64)
	if !ok {
		panic("Invalid argument")
	}
	for _, s := range ss[1:] {
		n, ok := s.(float64)
		if !ok {
			panic("Invalid argument")
		}
		r = r && (f <= n)
		f = n
	}
	if r {
		return 1.0
	}
	return Nil
}
