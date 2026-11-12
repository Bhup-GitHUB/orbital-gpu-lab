package nn

import (
	"math"
	"math/rand"

	"orbital-gpu-lab/internal/data"
)

const Hidden = 8
const WeightCount = Hidden*3 + Hidden + 1

type Network struct{ Weights []float32 }

func New(seed int64) *Network {
	r := rand.New(rand.NewSource(seed))
	w := make([]float32, WeightCount)
	for i := 0; i < Hidden*3; i++ {
		w[i] = float32(r.NormFloat64() * 0.3)
	}
	for i := Hidden * 3; i < Hidden*4; i++ {
		w[i] = float32(r.NormFloat64() * 0.3)
	}
	return &Network{Weights: w}
}

func finite(x float32) bool { return !math.IsNaN(float64(x)) && !math.IsInf(float64(x), 0) }

func (n *Network) forward(s data.Sample) ([Hidden]float32, float32, bool) {
	var h [Hidden]float32
	ok := true
	for i := range h {
		z := n.Weights[3*i]*s.X + n.Weights[3*i+1]*s.Y + n.Weights[3*i+2]
		if !finite(z) {
			ok = false
		}
		h[i] = float32(math.Tanh(float64(z)))
	}
	z := n.Weights[Hidden*4]
	for i := range h {
		z += h[i] * n.Weights[Hidden*3+i]
	}
	if !finite(z) {
		ok = false
	}
	return h, z, ok
}

func (n *Network) Accuracy(samples []data.Sample) (float64, bool) {
	correct := 0
	invalid := false
	for _, s := range samples {
		_, z, ok := n.forward(s)
		if !ok {
			invalid = true
			continue
		}
		if (z >= 0) == (s.Label == 1) {
			correct++
		}
	}
	return float64(correct) / float64(len(samples)), invalid
}

func (n *Network) Step(samples []data.Sample, rate float32) {
	grad := make([]float32, len(n.Weights))
	for _, s := range samples {
		h, z, _ := n.forward(s)
		p := float32(1 / (1 + math.Exp(-float64(z))))
		d := p - float32(s.Label)
		grad[Hidden*4] += d
		for i := range h {
			grad[Hidden*3+i] += d * h[i]
			dh := d * n.Weights[Hidden*3+i] * (1 - h[i]*h[i])
			grad[3*i] += dh * s.X
			grad[3*i+1] += dh * s.Y
			grad[3*i+2] += dh
		}
	}
	scale := rate / float32(len(samples))
	for i := range n.Weights {
		n.Weights[i] -= scale * grad[i]
	}
}
