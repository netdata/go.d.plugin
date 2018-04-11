package raw

import "testing"

func TestValidChartType(t *testing.T) {

	for _, v := range []string{Line, Area, Stacked} {
		if !ValidChartType(v) {
			t.Fatalf("function returned false for correct chart type")
		}
	}

	for _, v := range []string{"", "this", "is", "wrong"} {
		if ValidChartType(v) {
			t.Fatalf("function returned true for incorrect chart type")
		}
	}

}

