package lisp

import (
	"fmt"
	"io"
	"os"
)

// Circumvent lame initialization loop detection. An explicit init() allows
// builtinDefine et al to reference global.
func init() {
	globalData := map[sym]sexpr{
		// Misc. primitives (primitives.go)
		"if":     primitive(primitiveIf),
		"for":    primitive(primitiveFor),
		"lambda": primitive(primitiveLambda),
		"let":    primitive(primitiveLet),
		"define": primitive(primitiveDefine),
		"quote":  primitive(primitiveQuote),
		"begin":  primitive(primitiveBegin),

		// Nil
		"nil": Nil,

		// Misc
		"read":  function(builtinRead),
		"eval":  function(builtinEval),
		"apply": function(builtinApply),
		"print": function(builtinPrint),

		// Cons manipulation (cons.go)
		"cons": function(builtinCons),
		"car":  function(builtinCar),
		"cdr":  function(builtinCdr),

		// Basic stuff
		"equal?": function(builtinEqual),

		// Boolean arithmetic (bool.go)
		"not": function(builtinNot),

		// Arithmetic (math.go)
		"=":  function(builtinEq),
		"/=": function(builtinNe),
		"<":  function(builtinLt),
		">":  function(builtinGt),
		"<=": function(builtinLe),
		">=": function(builtinGe),
		"+":  function(builtinAdd),
		"-":  function(builtinSub),
		"*":  function(builtinMul),
		"/":  function(builtinDiv),
		"%":  function(builtinMod),

		// Go runtime (compat.go)
		"import": function(builtinImport),

		// Panics (panic.go)
		"recover": function(builtinRecover),
		"panic":   function(builtinPanic),

		// Concurrency
		"make-chan": function(builtinMakeChan),
		"go": primitive(primitiveGo),
		"<-": function(builtinLeftArrow),
	}

	global = &scope{globalData, nil}

	// Now interpret init_lisp
	load(init_lisp)
}

// (make-chan)
//
// Makes a channel for sexprs.
func builtinMakeChan(sc *scope, ss []sexpr) sexpr {
	return make(chan sexpr)
}

// (<- c)      Gets the next value from a channel c.
// (<- c val)  Sends val to a channel c.
func builtinLeftArrow(sc *scope, ss []sexpr) sexpr {
	if len(ss) == 1 {
		// (<- channel) returns the next value from the channel
		switch c := ss[0].(type) {
		default:
			panic(fmt.Sprintf("Expected a channel. Got type %T", c))
		case chan sexpr:
			return <-c
		}
	} else if len(ss) == 2 {
		// (<- channel value) 
		val := ss[1]
		switch c := ss[0].(type) {
		default:
			panic(fmt.Sprintf("Expected a channel. Got type %T", c))
		case chan sexpr:
			c <- val
			return Nil
		}
	}
	msg := fmt.Sprint("Unexpected arguments to <-.",
			"Wanted (<- c) or (<- c val), got (<-", ss,
			")")
	panic(msg)
}

// (read)
//
// Reads one s-expression from standard input.
func builtinRead(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 0 {
		panic("Invalid number of arguments")
	}
	v, err := parse(GetRuneScanner(os.Stdin))
	if err != nil && err != io.EOF {
		panic(err)
	} else if err == io.EOF {
		panic(sym("eof"))
	}
	return v
}

// (eval expr)
//
// Evaluates an s-expression.
func builtinEval(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 1 {
		panic("Invalid number of arguments")
	}
	return eval(sc, ss[0]) // TODO custom scope
}

// (print expr)
//
// Prints an s-expression.
func builtinPrint(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 1 {
		panic("Invalid number of arguments")
	}
	fmt.Printf("%s\n", asString(ss[0]))
	return Nil
}

// (apply func '(arg1 ...))
func builtinApply(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 2 {
		panic("Invalid number of arguments")
	}
	return apply(sc, ss[0], flatten(ss[1]))
}

// (equal? arg1 ...)
func builtinEqual(sc *scope, ss []sexpr) sexpr {
	if len(ss) == 0 {
		return 1.0
	}
	r := true
	for _, s := range ss[1:] {
		r = r && (s == ss[0])
	}
	if r {
		return 1.0
	}
	return Nil
}
