Kakapo is an experimental lisp interpreter in Go. It aims to provide maximum
compatibility between the interpreted lisp code and any present go code. For
usage, look at repl.lsp (the most complex program written with kakapo so far)
and kakapo.go (which showcases the entire api of the backing library).

# Installing
	1. Make sure kakapo is in the src/ dir of a directory in your GOPATH.
	2. ./configure
	3. go install

# Running
	$ kakapo
	Welcome to Kakapo
	kakapo> (print "Hello, 世界") 
	"Hello, 世界"
	nil

