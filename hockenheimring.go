package main

import "fmt"

type pit struct {
	i int
	f int
	g int
}

func searchStartPit(pits []pit) int {
	pitStart, fuelInTank, totalNetFuel := 1, 0, 0
	for _, pit := range pits {
		netFuelI := pit.f - pit.g
		if fuelInTank+netFuelI < 0 {
			pitStart = pit.i + 1
			fuelInTank = 0
		} else {
			fuelInTank += netFuelI
		}
		totalNetFuel += netFuelI
	}
	if totalNetFuel < 0 {
		return -1
	}
	return pitStart
}

func printFuel(pits []pit, start int) {
	fuelInTank := 0
	for i := 0; i < len(pits); i++ {
		mI := (start + i - 1) % len(pits)
		pitI := mI + 1
		fmt.Printf("Pit %d: fuel in tank %d\n", pitI, fuelInTank)
		fuelInTank += pits[mI].f - pits[mI].g
	}
}

func main() {
	testPits := []pit{
		{i: 1, f: 9, g: 5},
		{i: 2, f: 1, g: 4},
		{i: 3, f: 9, g: 4},
		{i: 4, f: 5, g: 2},
		{i: 5, f: 7, g: 9},
		{i: 6, f: 3, g: 2},
		{i: 7, f: 2, g: 4},
		{i: 8, f: 6, g: 3},
		{i: 9, f: 1, g: 6},
		{i: 10, f: 2, g: 4},
		{i: 11, f: 7, g: 4},
		{i: 12, f: 0, g: 2},
		{i: 13, f: 1, g: 4},
	}
	start := searchStartPit(testPits)
	printFuel(testPits, start)
	testPits2 := []pit{
		{i: 1, f: 2, g: 4},
		{i: 2, f: 6, g: 3},
		{i: 3, f: 1, g: 6},
		{i: 4, f: 2, g: 4},
		{i: 5, f: 7, g: 4},
		{i: 6, f: 0, g: 2},
		{i: 7, f: 1, g: 4},
		{i: 8, f: 9, g: 5},
		{i: 9, f: 1, g: 4},
		{i: 10, f: 9, g: 4},
		{i: 11, f: 5, g: 2},
		{i: 12, f: 7, g: 9},
		{i: 13, f: 3, g: 2},
	}
	start2 := searchStartPit(testPits2)
	printFuel(testPits2, start2)
}
