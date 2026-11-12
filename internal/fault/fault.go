package fault

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

func Flip(value float32, bit int) float32 {
	if bit < 0 || bit > 31 {
		panic("bit index out of range")
	}
	return math.Float32frombits(math.Float32bits(value) ^ (uint32(1) << uint(bit)))
}

func Class(bit int) string {
	switch {
	case bit == 31:
		return "sign"
	case bit >= 23:
		return "exponent"
	default:
		return "mantissa"
	}
}

func ChooseBit(r *rand.Rand, choice string) (int, error) {
	switch choice {
	case "sign":
		return 31, nil
	case "exponent":
		return 23 + r.Intn(8), nil
	case "mantissa":
		return r.Intn(23), nil
	case "all":
		return r.Intn(32), nil
	default:
		bit, err := strconv.Atoi(choice)
		if err != nil || bit < 0 || bit > 31 {
			return 0, fmt.Errorf("bit must be all, sign, exponent, mantissa, or 0..31")
		}
		return bit, nil
	}
}

func Inject(weights []float32, r *rand.Rand, bit int) int {
	index := r.Intn(len(weights))
	weights[index] = Flip(weights[index], bit)
	return index
}
