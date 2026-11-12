package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"

	"orbital-gpu-lab/internal/data"
	"orbital-gpu-lab/internal/fault"
	"orbital-gpu-lab/internal/guard"
	"orbital-gpu-lab/internal/nn"
)

type result struct {
	mode                                             string
	accuracy                                         float64
	invalid                                          bool
	injected, detected, rollbacks, wasted, completed int
	overhead                                         float64
}

func run(mode string, seed int64, steps, interval, blockSize int, rate float64, bitChoice string) result {
	train, test := data.Blobs(seed+1, 400), data.Blobs(seed+2, 200)
	net := nn.New(seed)
	r := rand.New(rand.NewSource(seed + 3))
	g := guard.New(blockSize)
	checkpoint := guard.Save(net.Weights, 0)
	out := result{mode: mode}
	start := time.Now()
	var checksumTime time.Duration
	update := func() { t := time.Now(); g.Update(net.Weights); checksumTime += time.Since(t) }
	if mode != "none" {
		update()
	}
	step := 0
	attempts := 0
	for step < steps {
		if attempts >= steps*100+1000 {
			fmt.Fprintf(os.Stderr, "%s: stopped after %d attempts; fault rate prevents progress\n", mode, attempts)
			break
		}
		attempts++
		if mode != "none" {
			t := time.Now()
			valid := g.Verify(net.Weights)
			checksumTime += time.Since(t)
			if !valid {
				out.detected++
				if mode == "detect" {
					fmt.Fprintf(os.Stderr, "detect: checksum mismatch before step %d; stopping\n", step+1)
					break
				}
				out.rollbacks++
				out.wasted += step - checkpoint.Step
				fmt.Fprintf(os.Stderr, "recover: checksum mismatch before step %d; rolling back to step %d\n", step+1, checkpoint.Step)
				checkpoint.Restore(net.Weights)
				step = checkpoint.Step
				update()
			}
		}
		net.Step(train, 0.1)
		step++
		if mode != "none" {
			update()
		}
		if mode == "recover" && step%interval == 0 {
			checkpoint = guard.Save(net.Weights, step)
		}
		if step < steps {
			flips := int(rate / 1000)
			if r.Float64() < math.Mod(rate, 1000)/1000 {
				flips++
			}
			for i := 0; i < flips; i++ {
				bit, _ := fault.ChooseBit(r, bitChoice)
				fault.Inject(net.Weights, r, bit)
				out.injected++
			}
		}
	}
	elapsed := time.Since(start)
	out.completed = step
	out.accuracy, out.invalid = net.Accuracy(test)
	if base := elapsed - checksumTime; base > 0 {
		out.overhead = float64(checksumTime) / float64(base) * 100
	}
	return out
}

func main() {
	seed := flag.Int64("seed", 1, "random seed")
	mode := flag.String("mode", "all", "all, none, detect, or recover")
	steps := flag.Int("steps", 300, "training steps")
	interval := flag.Int("checkpoint", 25, "checkpoint interval in steps")
	block := flag.Int("block", 8, "weights per checksum block")
	rate := flag.Float64("rate", 10, "expected flips per 1000 step boundaries")
	bit := flag.String("bit", "all", "all, sign, exponent, mantissa, or bit index 0..31")
	flag.Parse()
	if *steps <= 0 || *interval <= 0 || *block <= 0 || math.IsNaN(*rate) || math.IsInf(*rate, 0) || *rate < 0 || *rate > 1000 {
		log.Fatal("steps, checkpoint, and block must be positive; rate must be between 0 and 1000")
	}
	if _, err := fault.ChooseBit(rand.New(rand.NewSource(*seed)), *bit); err != nil {
		log.Fatal(err)
	}
	if *mode != "all" && *mode != "none" && *mode != "detect" && *mode != "recover" {
		log.Fatal("mode must be all, none, detect, or recover")
	}
	modes := []string{"none", "detect", "recover"}
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "Mode\tAccuracy\tNaN/Inf\tSteps\tInjected\tDetected\tRollbacks\tWasted steps\tChecksum overhead")
	for _, m := range modes {
		if *mode != "all" && *mode != m {
			continue
		}
		v := run(m, *seed, *steps, *interval, *block, *rate, *bit)
		fmt.Fprintf(w, "%s\t%.2f%%\t%t\t%d\t%d\t%d\t%d\t%d\t%.2f%%\n", v.mode, v.accuracy*100, v.invalid, v.completed, v.injected, v.detected, v.rollbacks, v.wasted, v.overhead)
	}
	w.Flush()
}
