# orbital-gpu-lab

Three small, CPU-only Go experiments inspired by the challenges of operating compute in orbit. They use the Go standard library and require Go 1.22 or newer; no GPU, external data, or network access is needed.

This is a **software simulation of bit flips**, not a model of real radiation rates, hardware failure behavior, or Starcloud's actual design. The fault rates and thermal assumptions are illustrative.

## Run

```sh
make all                    # run all three demonstrations
make bitflip                # writes bitflip.csv in the current directory
make protect
make radiator
make test

go run ./cmd/bitflip --seed 1 --trials 200 --bit all --csv results.csv
go run ./cmd/bitflip --bit exponent --trials 500
go run ./cmd/bitflip --bit 30
go run ./cmd/protect --mode all --seed 1 --steps 300 --rate 10 --checkpoint 25 --block 8
go run ./cmd/protect --mode recover --bit sign --rate 20
go run ./cmd/radiator --watts 700 --temp 27 --emissivity 0.9 --sides 2
```

Use `go run ./cmd/<name> --help` for all options. The random data, initialization, and fault choices are repeatable with the same seed and flags. Timing measurements vary between runs.

## Bit-flip lab

A two-layer float32 MLP learns a seeded two-class Gaussian dataset with full-batch gradient descent. Each trial flips exactly one bit in a uniformly chosen weight or bias, measures test accuracy, and restores the original bits. By default the command sweeps all 32 positions; `--bit` can select `sign`, `exponent`, `mantissa`, or an index from 0 to 31. Bit 0 is the lowest mantissa bit, bit 30 the highest exponent bit, and bit 31 the sign bit. The printed example applies each flip independently to the same weight.

Results include the mean test accuracy, the fraction of trials more than 10 percentage points below baseline, and the fraction in which any test sample produced a non-finite intermediate or output. Samples with non-finite values count as incorrect in accuracy. The CSV contains percentages with six decimal places. Example with `--seed 1 --trials 20` (selected rows):

```text
Baseline accuracy: 100.00% (seed 1)
Bit  Class     Mean accuracy  Drop >10 points  NaN/Inf
0    mantissa  100.00%        0.00%            0.00%
22   mantissa  100.00%        0.00%            0.00%
30   exponent   78.45%       35.00%           30.00%
31   sign      100.00%        0.00%            0.00%
```

## Checksums and rollback

CRC32 covers blocks of weight and bias bit patterns. Protected modes verify before each step and recompute after each legitimate step. Flips occur at step boundaries, after the update and any checkpoint, and before the next verification. The rate is the probability per boundary multiplied by 1,000, between 0 and 1,000; it is a simulation parameter, not a physical rate.

- `none` trains without checking for corruption.
- `detect` stops at the first checksum mismatch.
- `recover` restores the latest checkpoint, including the initial weights, and repeats lost steps until the requested number of completed steps is reached. `--checkpoint` sets the interval and `--block` sets the number of weights per CRC.

Each mode starts with the same model, data, and fault RNG seed; recovery consumes additional fault draws as steps are repeated, so later injections do not necessarily match across modes. Detected counts are mismatch events, not a count of individual flipped bits. Wasted steps are completed training steps discarded by rollback. A run stops early if repeated rollbacks prevent progress (after `100 * steps + 1000` attempted steps). An undetected corruption or a non-finite training value can still affect results; CRC is not a guarantee against all possible multi-bit patterns.

Checksum overhead is measured as time spent computing and verifying CRCs divided by the rest of that mode's training-loop time. It includes allocations in checksum code but excludes checkpoint copying from the numerator; it is an approximate laptop timing, not a controlled benchmark. Accuracy is evaluated at the end, including an early stop. Example (`--seed 1`, timings vary):

```text
Mode     Accuracy  NaN/Inf  Steps  Injected  Detected  Rollbacks  Wasted steps  Checksum overhead
none     0.00%     true     300    5         0         0          0             0.00%
detect   100.00%   false    20     1         1         0          0             0.60%
recover  100.00%   false    300    7         7         7          88            0.50%
```

## Radiator calculator

The calculator solves `P = emissivity × sigma × A × T⁴` for panel area, with temperature in kelvin and sigma = 5.670374419 × 10⁻⁸ W m⁻² K⁻⁴. The reported area is the physical area of one side of the panel: two equally effective radiating sides halve the required panel area. The sweep holds heat load, emissivity, and side count fixed.

It assumes a 0 K radiative sink and ignores absorbed sunlight, Earth infrared, and heat transport from chip to radiator. It is not a spacecraft thermal design.

```text
Required one-sided panel area: 1.690 m² (700.0 W, 27.0 °C, emissivity 0.90, 1 radiating side(s))
Temperature (°C)  Panel area (m²)
               0            2.464
              20            1.857
              40            1.426
              60            1.113
              80            0.882
```

## Background

[Starcloud reports](https://www.starcloud.com/starcloud-1) that Starcloud-1 carried an NVIDIA H100 into orbit and trained nanoGPT in space. The satellite launched in November 2025, as [reported by Data Center Dynamics](https://www.datacenterdynamics.com/en/news/starcloud-1-satellite-reaches-space-with-nvidia-h100-gpu-now-operating-in-orbit/). These experiments are independent educational examples and make no claim to reproduce its hardware or software.
