(begin
  ; need help ?
  (doc println)
  (dump-scope)

  ; look under the hood
  (inspect hello)
  (println (type +))
  (println (tokenize "(println \"hello\")"))

  ; calling variadic Go functions
  (println "Hello" "from" "Parens!")

  ; binding values
  (label π 3.1412)

  ; let starts a new scope
  (let
    (label π 3)
    (printf "integer part of pi is %f\n" π))

  ; value of π should now be reset to original
  (printf "but real value of pi is %f\n" π)

  ; let's use some cool looking characters since
  ; parens supports UTF-8
  (label ∂ 0)
  (label ∑ +)
  (label ≠ (lambda (a b) (not (== a b))))

  (println "1 + 2 =" (∑ 1 2))
  (println "false and nil are not same =" (≠ false nil))

  ; defining a lambda
  (label square (lambda (a) (* a a)))

  ; calling a lambda, obviously
  (printf "square of 2 is = %f\n" (square 2))

  ; we need to do some math obviously
  (printf "complex math answer %f\n" (* 1 (- 2 (+ 1 (/ 3 3))))))


