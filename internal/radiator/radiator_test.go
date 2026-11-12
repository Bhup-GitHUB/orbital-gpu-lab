package radiator

import (
	"math"
	"testing"
)

func TestArea(t *testing.T) {
	a, err := Area(700, 27, 0.9, 1)
	if err != nil || math.Abs(a-1.69) > 0.01 {
		t.Fatalf("area = %g, err = %v", a, err)
	}
	b, err := Area(700, 27, 0.9, 2)
	if err != nil || math.Abs(b-a/2) > 1e-12 {
		t.Fatalf("two-sided area = %g, err = %v", b, err)
	}
	for _, input := range []struct {
		w, t, e float64
		sides   int
	}{{-1, 20, 1, 1}, {1, -273.15, 1, 1}, {1, 20, 0, 1}, {1, 20, 1.1, 1}, {1, 20, 1, 3}} {
		if _, err := Area(input.w, input.t, input.e, input.sides); err == nil {
			t.Fatalf("accepted invalid input %+v", input)
		}
	}
}
