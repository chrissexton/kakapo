for x in repl core; do
    ./txt2go.sh $x < $x.lisp > $x.go
done

