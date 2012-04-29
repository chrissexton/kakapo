(S' "Macros")

(defmacro square (x) (list '* x x))
(S' "Macros 1")
(T' (equal? '(* 9 9)
            (macroexpand-1 '(square 9))))
(S' "Macros 2")
(T' (= 9 (square 3)))

(defmacro or2 (a b) (list 'if a a b))
(S' "Macros 3")
(T' (= 1 (or2 1 nil)))
(S' "Macros 4")
(T' (= 1 (or2 1 (panic "macros are hosed"))))
(S' "Macros 5")
(T' (= nil (or2 nil nil)))

