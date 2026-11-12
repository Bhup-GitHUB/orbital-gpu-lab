package guard

import (
	"encoding/binary"
	"hash/crc32"
	"math"
)

type Guard struct {
	blockSize int
	sums      []uint32
}

func New(blockSize int) *Guard {
	if blockSize <= 0 {
		panic("block size must be positive")
	}
	return &Guard{blockSize: blockSize}
}

func (g *Guard) sumsFor(weights []float32) []uint32 {
	sums := make([]uint32, 0, (len(weights)+g.blockSize-1)/g.blockSize)
	buf := make([]byte, 4*g.blockSize)
	for start := 0; start < len(weights); start += g.blockSize {
		end := start + g.blockSize
		if end > len(weights) {
			end = len(weights)
		}
		for i, w := range weights[start:end] {
			binary.LittleEndian.PutUint32(buf[4*i:], math.Float32bits(w))
		}
		sums = append(sums, crc32.ChecksumIEEE(buf[:4*(end-start)]))
	}
	return sums
}

func (g *Guard) Update(weights []float32) { g.sums = g.sumsFor(weights) }

func (g *Guard) Verify(weights []float32) bool {
	sums := g.sumsFor(weights)
	if len(sums) != len(g.sums) {
		return false
	}
	for i := range sums {
		if sums[i] != g.sums[i] {
			return false
		}
	}
	return true
}

type Checkpoint struct {
	weights []float32
	Step    int
}

func Save(weights []float32, step int) Checkpoint {
	return Checkpoint{weights: append([]float32(nil), weights...), Step: step}
}

func (c Checkpoint) Restore(weights []float32) {
	if len(weights) != len(c.weights) {
		panic("checkpoint weight count mismatch")
	}
	copy(weights, c.weights)
}
