package guard

import (
	"math"
	"orbital-gpu-lab/internal/fault"
	"testing"
)

func TestChecksumDetectsEverySingleBit(t *testing.T) {
	w := []float32{1, -2, 3, 4, 5}
	g := New(2)
	g.Update(w)
	for i := range w {
		for bit := 0; bit < 32; bit++ {
			w[i] = fault.Flip(w[i], bit)
			if g.Verify(w) {
				t.Fatalf("missed weight %d bit %d", i, bit)
			}
			w[i] = fault.Flip(w[i], bit)
		}
	}
	if !g.Verify(w) {
		t.Fatal("restored weights failed verification")
	}
}

func TestCheckpointRestoresBits(t *testing.T) {
	w := []float32{1, math.Float32frombits(0x7fc01234), -3}
	c := Save(w, 7)
	w[0], w[1], w[2] = 0, 0, 0
	c.Restore(w)
	if c.Step != 7 || math.Float32bits(w[0]) != 0x3f800000 || math.Float32bits(w[1]) != 0x7fc01234 || math.Float32bits(w[2]) != 0xc0400000 {
		t.Fatal("checkpoint did not restore exact bits")
	}
}
