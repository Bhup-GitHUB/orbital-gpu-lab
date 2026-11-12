package radiator

import (
	"fmt"
	"math"
)

const Sigma = 5.670374419e-8

func Area(watts, celsius, emissivity float64, sides int) (float64, error) {
	if math.IsNaN(watts) || math.IsInf(watts, 0) || watts < 0 || math.IsNaN(celsius) || math.IsInf(celsius, 0) || celsius <= -273.15 || math.IsNaN(emissivity) || math.IsInf(emissivity, 0) || emissivity <= 0 || emissivity > 1 || (sides != 1 && sides != 2) {
		return 0, fmt.Errorf("heat must be nonnegative, temperature above absolute zero, emissivity in (0,1], and sides 1 or 2")
	}
	area := watts / (emissivity * Sigma * float64(sides) * math.Pow(celsius+273.15, 4))
	if math.IsInf(area, 0) || math.IsNaN(area) {
		return 0, fmt.Errorf("inputs exceed numerical range")
	}
	return area, nil
}
