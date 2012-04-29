package lisp

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

// Circumvent lame initialization loop detection. An explicit init() allows
// builtinDefine et al to reference global.
func init() {
	globalData := map[sym]sexpr{
		// Misc. primitives (primitives.go)
		"if":     primitive("if", primitiveIf),
		"for":    primitive("for", primitiveFor),
		"lambda": primitive("lambda", primitiveLambda),
		"let":    primitive("let", primitiveLet),
		"define": primitive("define", primitiveDefine),
		"quote":  primitive("quote", primitiveQuote),
		"begin":  primitive("begin", primitiveBegin),

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
		"list": function(builtinList),
		"list?": function(builtinIsList),

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
		"go": primitive("go", primitiveGo),
		"<-": function(builtinLeftArrow),

		// Macros
		"defmacro": primitive("defmacro", primitiveDefmacro),
		"macroexpand-1": primitive("macroexpand-1", primitiveMacroexpand1),
	}

	global = &scope{globalData, nil}

	// Now interpret init_lisp
	load(init_lisp)
}

// (list? expr)
//
// Tells whether the given expression is a list.
func builtinIsList(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 1 {
		msg := fmt.Sprint("Expected exactly one arg to list?, got",
			len(ss))
		panic(msg)
	}
	return isList(ss[0])
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
	return builtinEq(sc, ss);
}

// (= ...)
//
// Returns true if and only if all arguments are equal.
func builtinEq(sc *scope, ss []sexpr) sexpr {
	if len(ss) == 0 {
		return true
	}
	for _, e := range(ss[1:len(ss)]) {
		if !reflect.DeepEqual(ss[0], e) {
			return Nil
		}
	}
	return true
}

// (/= ...)
//
// Returns true if not all arguments are equal.
func builtinNe(sc *scope, ss []sexpr) sexpr {
	r := builtinEq(sc, ss)
	return builtinNot(sc, []sexpr{r})
}

