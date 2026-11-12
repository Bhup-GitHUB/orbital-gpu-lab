package main

import (
	"flag"
	"fmt"
	"log"
	"orbital-gpu-lab/internal/radiator"
)

func main() {
	watts := flag.Float64("watts", 700, "heat load in watts")
	temp := flag.Float64("temp", 27, "radiator temperature in Celsius")
	emissivity := flag.Float64("emissivity", 0.9, "radiator emissivity")
	sides := flag.Int("sides", 1, "radiating sides: 1 or 2")
	flag.Parse()
	area, err := radiator.Area(*watts, *temp, *emissivity, *sides)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Required one-sided panel area: %.3f m² (%.1f W, %.1f °C, emissivity %.2f, %d radiating side(s))\n", area, *watts, *temp, *emissivity, *sides)
	fmt.Println("Temperature (°C)  Panel area (m²)")
	for t := 0; t <= 80; t += 10 {
		a, _ := radiator.Area(*watts, float64(t), *emissivity, *sides)
		fmt.Printf("%16d  %15.3f\n", t, a)
	}
	fmt.Println("Simplification: ignores absorbed sunlight, Earth infrared, and heat transport from chip to radiator; assumes radiation into a 0 K sink.")
}
