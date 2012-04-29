package lisp

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func load(s string) {
	EvalFrom(strings.NewReader(s))
}

// TODO: Rename this ExecFrom
func EvalFrom(ior io.Reader) {
	// TODO parse and eval in separate goroutines

	r := bufio.NewReader(ior)
	e, err := parse(r)
	for err == nil {
		eval(global, e)
		e, err = parse(r)
	}
}

func EvalStr(s string) sexpr {
	r := bufio.NewReader(strings.NewReader(s))
	e, err := parse(r)
	if err != nil {
		panic(fmt.Sprint("Failed to evaluate", s))
	}
	return eval(global, e)
}

