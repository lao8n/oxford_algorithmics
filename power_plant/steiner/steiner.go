package main

import (
	"math"
)

func dijkstras(powerPlants []loc) {
	// build adjacency list of edges
	adjList := make(map[loc][]edge, len(powerPlants))

	// initialize adjacency list
	for i, powerPlantI := range powerPlants {
		for j := i + 1; j < len(powerPlants); j++ {
			powerPlantJ := powerPlants[j]
			// from power plant I to power plant J
			fromIToJ := edge{powerPlantI, powerPlantJ, calculateEdgeCost(powerPlants, i, j)}
			adjList[powerPlantI] = append(adjList[powerPlantI], fromIToJ)
			// from power plant J to power plant I
			fromJToI := edge{powerPlantJ, powerPlantI, calculateEdgeCost(powerPlants, i, j)}
			adjList[powerPlantJ] = append(adjList[powerPlantJ], fromJToI)
		}
	}

	// build dp table
	dist := make(map[loc]map[loc]float64, len(powerPlants))
	for _, powerPlant := range powerPlants {
		dist[powerPlant] = make(map[loc]float64, len(powerPlants))

	}

}

func calculateEdgeCost(powerPlants []loc, i, j int) float64 {
	sumSquares := math.Pow(powerPlants[i].x-powerPlants[j].x, 2) + math.Pow(powerPlants[i].y-powerPlants[j].y, 2)
	return math.Pow(sumSquares, 0.5)
}

type loc struct {
	x float64
	y float64
}

type edge struct {
	from loc
	to   loc
	cost float64
}
