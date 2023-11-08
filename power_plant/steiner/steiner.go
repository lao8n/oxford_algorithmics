package main

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/triangulation"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func prims(powerPlants []loc) (float64, []edge) {
	// build visited map
	visited := make(map[loc]bool, len(powerPlants))
	// build adjacency list of edges
	adjList := make(map[loc][]edge, len(powerPlants))

	// initialize visited map and adjacency list
	for i, powerPlantI := range powerPlants {
		visited[powerPlantI] = false
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

	// initailize heap with point 0 edges
	pointZeroEdges := adjList[powerPlants[0]]
	h := make(Heap, len(pointZeroEdges))
	copy(h, pointZeroEdges)
	heap.Init(&h)
	visited[powerPlants[0]] = true
	visitedCount := 1
	rollingCost := 0.0

	// initialize data
	edges := make([]edge, 0)

	// loop until all power plants are visited
	for h.Len() > 0 && visitedCount < len(powerPlants) {
		// get neighbour edge with lowest cost
		minEdge := heap.Pop(&h).(edge)

		// we are going from current tree to new vertex
		newVertex := minEdge.to
		if visited[minEdge.to] {
			continue
		}
		// process edge
		visited[newVertex] = true
		edges = append(edges, minEdge)
		visitedCount++
		rollingCost += minEdge.cost

		// add neighbours of new edge to heap
		for _, neighbourEdge := range adjList[newVertex] {
			if !visited[neighbourEdge.to] {
				heap.Push(&h, neighbourEdge)
			}
		}
	}
	return rollingCost, edges
}

func calculateEdgeCost(powerPlants []loc, i, j int) float64 {
	sumSquares := math.Pow(powerPlants[i].x-powerPlants[j].x, 2) + math.Pow(powerPlants[i].y-powerPlants[j].y, 2)
	return math.Pow(sumSquares, 0.5)
}

func nonTerminal(powerPlants []loc, w int, h int) []loc {
	// delauney triangulation
	powerPlantMap := make(map[loc]bool, len(powerPlants))
	// Triangulate(points normgeom.NormPointGroup, w, h int) []geom.Triangle
	var points normgeom.NormPointGroup
	for _, powerPlant := range powerPlants {
		points = append(points, normgeom.NormPoint{powerPlant.x, powerPlant.y})
		powerPlantMap[powerPlant] = true
	}

	// add non-terminal nodes
	nonTerminalNodes := make([]loc, 0)
	for _, triangle := range triangulation.Triangulate(points, w, h) {
		nonTerminal := centroid(triangle)
		if !powerPlantMap[nonTerminal] {
			nonTerminalNodes = append(nonTerminalNodes, nonTerminal)
		}
	}
	return nonTerminalNodes
}

func centroid(triangle geom.Triangle) loc {
	centroid := loc{
		x: float64(triangle.Points[0].X+triangle.Points[1].X+triangle.Points[2].X) / 3.0,
		y: float64(triangle.Points[0].Y+triangle.Points[1].Y+triangle.Points[2].Y) / 3.0,
	}
	return centroid
}

func plotCosts(powerPlants []loc, edges []edge) {
	p := plot.New()
	for _, edge := range edges {
		pts := make(plotter.XYs, 2)
		pts[0].X = edge.from.x
		pts[0].Y = edge.from.y
		pts[1].X = edge.to.x
		pts[1].Y = edge.to.y
		line, err := plotter.NewLine(pts)
		if err != nil {
			fmt.Println(err)
		}
		p.Add(line)
	}
	pts := make(plotter.XYs, len(powerPlants))
	for i, powerPlant := range powerPlants {
		pts[i].X = powerPlant.x
		pts[i].Y = powerPlant.y
	}
	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		fmt.Println(err)
	}
	p.Add(scatter)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "line_graph.png"); err != nil {
		fmt.Println(err)
	}
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

type Heap []edge

func (h *Heap) Push(x interface{}) { *h = append(*h, x.(edge)) }
func (h *Heap) Pop() interface{} {
	previous, n, popped := *h, h.Len(), edge{}
	*h, popped = previous[:n-1], previous[n-1]
	return popped
}
func (h Heap) Len() int           { return len(h) }
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h Heap) Less(i, j int) bool { return h[i].cost < h[j].cost }

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
	minCost, minEdges := prims(powerPlants)
	finalNodes := powerPlants
	for _, nT := range nonTerminal(powerPlants, 1, 1) {
		cost, edges := prims(append(finalNodes, nT))
		if cost < minCost {
			fmt.Println(minCost)
			minCost = cost
			minEdges = edges
			finalNodes = append(finalNodes, nT)
		}
	}
	plotCosts(finalNodes, minEdges)
	fmt.Println(minCost, minEdges)
}
