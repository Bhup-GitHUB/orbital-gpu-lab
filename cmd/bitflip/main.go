package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"orbital-gpu-lab/internal/data"
	"orbital-gpu-lab/internal/fault"
	"orbital-gpu-lab/internal/nn"
	"orbital-gpu-lab/internal/report"
)

func main() {
	seed := flag.Int64("seed", 1, "random seed")
	trials := flag.Int("trials", 200, "trials per bit")
	bitChoice := flag.String("bit", "all", "all, sign, exponent, mantissa, or bit index 0..31")
	csvPath := flag.String("csv", "bitflip.csv", "output CSV path")
	flag.Parse()
	if *trials <= 0 {
		log.Fatal("trials must be positive")
	}
	r := rand.New(rand.NewSource(*seed + 3))
	if _, err := fault.ChooseBit(r, *bitChoice); err != nil {
		log.Fatal(err)
	}
	train, test := data.Blobs(*seed+1, 400), data.Blobs(*seed+2, 200)
	net := nn.New(*seed)
	for i := 0; i < 300; i++ {
		net.Step(train, 0.1)
	}
	baseline, _ := net.Accuracy(test)
	fmt.Printf("Baseline accuracy: %.2f%% (seed %d)\n", baseline*100, *seed)
	value := net.Weights[0]
	fmt.Printf("Example weight: %g; bit 0: %g; bit 30: %g\n", value, fault.Flip(value, 0), fault.Flip(value, 30))
	rows := make([]report.BitResult, 0, 32)
	for bit := 0; bit < 32; bit++ {
		selectedBit, parseErr := strconv.Atoi(*bitChoice)
		if *bitChoice != "all" && *bitChoice != fault.Class(bit) && (parseErr != nil || selectedBit != bit) {
			continue
		}
		row := report.BitResult{Bit: bit, Class: fault.Class(bit)}
		for i := 0; i < *trials; i++ {
			index := r.Intn(len(net.Weights))
			before := net.Weights[index]
			net.Weights[index] = fault.Flip(before, bit)
			accuracy, invalid := net.Accuracy(test)
			net.Weights[index] = before
			row.Accuracy += accuracy
			if baseline-accuracy > 0.10 {
				row.Drops++
			}
			if invalid {
				row.Invalid++
			}
		}
		row.Accuracy /= float64(*trials)
		row.Drops /= float64(*trials)
		row.Invalid /= float64(*trials)
		rows = append(rows, row)
	}
	report.BitTable(os.Stdout, rows)
	f, err := os.Create(*csvPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := report.BitCSV(f, rows); err != nil {
		f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("CSV: %s\n", *csvPath)
}
