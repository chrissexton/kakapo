# Generates a file with keywords for completion in the ka wrapper script.

grep -o '"[^"]*":' lisp/builtins.go | sed 's/[":]//g'

