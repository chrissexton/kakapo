(S' "Goroutines and channels")

;; Channels are truthy.
(T' (make-chan))

;; Run a goroutine that does nothing.
(F' (go nil))

;; Send a value over a channel from a goroutine.
(T'
  (let ((c (make-chan)))
    (go (<- c 3))
    (let ((x (<- c)))
      (equal? x 3))))

