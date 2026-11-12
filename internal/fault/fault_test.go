package fault

import (
	"math"
	"testing"
)

func TestFlipTwice(t *testing.T) {
	for _, bits := range []uint32{0, 0x80000000, 0x3f800000, 0x7fc01234} {
		for bit := 0; bit < 32; bit++ {
			value := math.Float32frombits(bits)
			flipped := Flip(value, bit)
			if got := math.Float32bits(flipped); got != bits^(uint32(1)<<uint(bit)) {
				t.Fatalf("bit %d: got %08x", bit, got)
			}
			if got := math.Float32bits(Flip(flipped, bit)); got != bits {
				t.Fatalf("bit %d: got %08x, want %08x", bit, got, bits)
			}
		}
	}
}
