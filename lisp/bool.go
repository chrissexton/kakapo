package lisp

// (not b)
func builtinNot(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 1 {
		panic("Invalid number of arguments")
	}
	return !IsTrue(ss[0])
}

// nil and false are the only false values.
// Everything else is considered true.
func IsTrue(s sexpr) bool {
	if s == nil {
		return false
	}
	switch s := s.(type) {
	case bool:
		return s
	}
	return true
}

