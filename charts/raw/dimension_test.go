package raw

import "testing"

func TestValidAlgorithm(t *testing.T) {
	for _, v := range []string{Absolute, Incremental, PercentOfAbsolute, PercentOfIncremental} {
		if !ValidAlgorithm(v) {
			t.Fatalf("function returned false for correct dimension algorithm")
		}
	}

	for _, v := range []string{"", "this", "is", "wrong"} {
		if ValidAlgorithm(v) {
			t.Fatalf("function returned true for incorrect dimension algorithm")
		}
	}

}
