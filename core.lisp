;; Core Lisp functions and macros

(define len (lambda (ls)
  (if (equal? ls nil)
    0
    (+ 1 (len (cdr ls))))))

