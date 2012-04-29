// This is the entry point for the REPL (Read Eval Print Loop).

package main

import (
	"flag"
	"fmt"
	. "kakapo/lisp"
	"os"
	"strings"
)

const VERSION = `0.5`

var (
	version = flag.Bool("V", false, "Display version information and exit")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("Kakapo %s\n", VERSION)
		return
	}

	// Expose imports
	for name, pkg := range _go_imports {
		ExposeImport(name, pkg)
	}

	// Expose globals
	ExposeGlobal("-interpreter", "Kakapo")
	ExposeGlobal("-interpreter-version", VERSION)

	args := flag.Args()
	// Define functions and macros from core.lisp
	EvalFrom(strings.NewReader(core))
	if len(args) == 0 {
		// Start the read-eval-print loop (repl.lisp)
		EvalFrom(strings.NewReader(repl))
	} else {
		for _, path := range args {
			if path == "-" {
				EvalFrom(os.Stdin)
			} else {
				file, err := os.Open(path)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				EvalFrom(file)
			}
		}
	}
}
