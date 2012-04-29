; List manipulation

(S' "Basic")
(F' nil)
(T' 'nil)

(S' "len")
(T' (= 0 (len '())))
(T' (= 1 (len '(a))))
(T' (= 2 (len '(a b))))

(S' "map")
(T' (equal? '() (map identity '())))
(T' (equal? '(1) (map identity '(1))))
(T' (equal? '(1 2) (map identity '(1 2))))

