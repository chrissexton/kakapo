package lisp

import "testing"

func TestTrue(t *testing.T) {
	trueVals := []string{"1", "2", "(not nil)", "(not (not (not nil)))",
			  "true", "(not false)"}
	for _, v := range(trueVals) {
		if !IsTrue(EvalStr(v)) {
			t.Error(v, "did not evaluate to truth.")
		}
	}
}

func testFalse(t *testing.T) {
	// falseVals := []string{"false", "(not true)"}
}

