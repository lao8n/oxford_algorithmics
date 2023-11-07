package main

import (
	"fmt"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func binarySearch(powerPlants []loc) (float64, loc, loc, []point, []point) {
	// get most east and most west x values
	mostWestInt, mostEastInt := powerPlants[0].x, powerPlants[0].x
	for _, powerPlant := range powerPlants {
		if powerPlant.x < mostWestInt {
			mostWestInt = powerPlant.x
		}
		if powerPlant.x > mostEastInt {
			mostEastInt = powerPlant.x
		}
	}

	// adapted binary search
	mostWest, mostEast := float64(mostWestInt), float64(mostEastInt)
	eps := 1e-1
	data1 := make([]point, 0)
	data2 := make([]point, 0)
	for mostEast-mostWest > eps {
		// search space: [mostWest - mid1 - mid2 - mostEast]
		mid1 := mostWest + (mostEast-mostWest)/3
		cost1, _, _ := calculateCost(powerPlants, mid1)
		mid2 := mostEast - (mostEast-mostWest)/3
		cost2, _, _ := calculateCost(powerPlants, mid2)
		data1 = append(data1, point{mid1, cost1})
		data2 = append(data2, point{mid2, cost2})
		// pick best search space
		if cost1 < cost2 {
			mostEast = mid2
			data1 = append(data1, point{mid1, cost1})
		} else {
			mostWest = mid1
			data1 = append(data1, point{mid2, cost2})
		}
	}
	lowestCost, bestNorth, bestSouth := calculateCost(powerPlants, mostWest)
	return lowestCost, bestNorth, bestSouth, data1, data1
}

func calculateCost(powerPlants []loc, transmissionLine float64) (float64, loc, loc) {
	horizontalCost, verticalCost := 0.0, 0.0
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
		horizontalCost += math.Abs(powerPlant.x - transmissionLine)
	}
	verticalCost = mostNorth - mostSouth
	return horizontalCost + verticalCost, loc{transmissionLine, mostNorth}, loc{transmissionLine, mostSouth}
}

type loc struct {
	x float64
	y float64
}

type point struct {
	x    float64
	cost float64
}

func plotCosts(data1 []point, data2 []point) {
	p := plot.New()
	p.X.Label.Text = "Transmission Line"
	p.Y.Label.Text = "Cost"
	pts1 := make(plotter.XYs, len(data1))
	for i, pt := range data1 {
		pts1[i].X = pt.x
		pts1[i].Y = pt.cost
	}
	line1, err := plotter.NewScatter(pts1)
	if err != nil {
		fmt.Println(err)
	}
	p.Add(line1)
	fmt.Println(line1)
	pts2 := make(plotter.XYs, len(data2))
	for i, pt := range data2 {
		pts2[i].X = pt.x
		pts2[i].Y = pt.cost
	}
	line2, err := plotter.NewScatter(pts2)
	if err != nil {
		fmt.Println(err)
	}
	p.Add(line2)
	// fmt.Println(line2)
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
	cost, north, south, data1, data2 := binarySearch(powerPlants)
	fmt.Println(data1, data2)
	fmt.Println(cost, north, south)
	plotCosts(data1, data2)
}
