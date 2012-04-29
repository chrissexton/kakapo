package lisp

var init_lisp = `; Kakapo interpreter initialization file

(define len (lambda (ls)
  (if (equal? ls nil)
    0
    (+ 1 (len (cdr ls))))))

(define map
  (lambda (f ls)
    (if (equal? ls '())
      '()
      (cons (f (car ls))
            (map f (cdr ls))))))

(define identity
  (lambda (x) x))

`
