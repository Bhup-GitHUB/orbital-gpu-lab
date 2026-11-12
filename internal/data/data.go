package data

import "math/rand"

type Sample struct {
	X, Y  float32
	Label int
}

func Blobs(seed int64, count int) []Sample {
	r := rand.New(rand.NewSource(seed))
	samples := make([]Sample, count)
	for i := range samples {
		label := i % 2
		center := float64(-1.4)
		if label == 1 {
			center = 1.4
		}
		samples[i] = Sample{float32(center + r.NormFloat64()*0.55), float32(center + r.NormFloat64()*0.55), label}
	}
	return samples
}
