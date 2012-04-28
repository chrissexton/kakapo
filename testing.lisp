; Testing rig library

(let ((-section "UNDEFINED-SECTION"))
  (define S'
    (lambda (s)
      (define -section s)))

  ;; T' tests its argument. If the argument evaluates to true then nothing
  ;; happens. Otherwise an exception is thrown.
  (define T'
    (lambda (b)
      (if b nil
        (panic -section))))

  (define F'
    (lambda (b)
      (T' (not b)))))
