package main

import (
	"fmt"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func bruteForce(powerPlants []loc) (int, loc, loc, []point) {
	// get most east and most west x values
	mostWest, mostEast := powerPlants[0].x, powerPlants[0].x
	for _, powerPlant := range powerPlants {
		if powerPlant.x < mostWest {
			mostWest = powerPlant.x
		}
		if powerPlant.x > mostEast {
			mostEast = powerPlant.x
		}
	}

	// try every possible transmission line and calculate the costs
	lowestCost, bestNorth, bestSouth := math.MaxInt, loc{}, loc{}
	data := make([]point, 0)
	for x := mostWest; x <= mostEast; x++ {
		cost, north, south := calculateCost(powerPlants, x)
		data = append(data, point{x, cost})
		if cost < lowestCost {
			lowestCost, bestNorth, bestSouth = cost, north, south
		}
	}
	return lowestCost, bestNorth, bestSouth, data
}

func calculateCost(powerPlants []loc, transmissionLine int) (int, loc, loc) {
	horizontalCost, verticalCost := 0, 0
	mostNorth, mostSouth := powerPlants[0].y, powerPlants[0].y
	for _, powerPlant := range powerPlants {
		// calculate north-south transmission line length
		if powerPlant.y > mostNorth {
			mostNorth = powerPlant.y
		}
		if powerPlant.y < mostSouth {
			mostSouth = powerPlant.y
		}
		// calculate east-west transmission line length
		horizontalCost += abs(powerPlant.x - transmissionLine)
	}
	verticalCost = mostNorth - mostSouth
	return horizontalCost + verticalCost, loc{transmissionLine, mostNorth}, loc{transmissionLine, mostSouth}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type loc struct {
	x int
	y int
}

type point struct {
	x    int
	cost int
}

func plotCosts(data []point) {
	p := plot.New()
	p.X.Label.Text = "Transmission Line"
	p.Y.Label.Text = "Cost"
	pts := make(plotter.XYs, len(data))
	for i, pt := range data {
		pts[i].X = float64(pt.x)
		pts[i].Y = float64(pt.cost)
	}
	line, err := plotter.NewLine(pts)
	if err != nil {
		fmt.Println(err)
	}
	p.Add(line)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "line_graph.png"); err != nil {
		fmt.Println(err)
	}
}

func main() {
	powerPlants := []loc{
		{12, 22},
		{16, 38},
		{18, 30},
		{23, 23},
		{22, 35},
		{36, 26},
		{32, 36},
		{40, 35},
	}
	cost, north, south, data := bruteForce(powerPlants)
	fmt.Println(data)
	fmt.Println(cost, north, south)
	plotCosts(data)
}
