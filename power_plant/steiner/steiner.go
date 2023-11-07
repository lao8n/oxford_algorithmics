package main

import (
	"container/heap"
	"fmt"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func dijkstras(powerPlants []loc) map[loc]map[loc]float64 {
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
	shortestPaths := make(map[loc]map[loc]float64, len(powerPlants))
	for _, powerPlantI := range powerPlants {
		shortestPaths[powerPlantI] = make(map[loc]float64, len(powerPlants))
		for _, powerPlantJ := range powerPlants {
			// same power plant so 0 distance
			if powerPlantI.x == powerPlantJ.x && powerPlantI.y == powerPlantJ.y {
				shortestPaths[powerPlantI][powerPlantJ] = 0
			} else {
				shortestPaths[powerPlantI][powerPlantJ] = math.MaxFloat64
			}
		}
	}

	// dijkstra's algorithm
	for _, powerPlant := range powerPlants {
		queue := []loc{powerPlant}
		visited := make(map[loc]bool, len(powerPlants))
		visited[powerPlant] = true
		for len(queue) > 0 {
			popped := queue[0]
			queue = queue[1:]
			for _, edge := range adjList[popped] {
				if !visited[edge.to] {
					shortestPaths[powerPlant][edge.to] = shortestPaths[powerPlant][edge.from] + edge.cost
					visited[edge.to] = true
					queue = append(queue, edge.to)
				}
			}
		}
	}
	return shortestPaths
}

func metricClosure(powerPlants []loc, shortestPaths map[loc]map[loc]float64) map[loc][]edge {
	adjList := make(map[loc][]edge, len(powerPlants))

	for from, m := range shortestPaths {
		for to, cost := range m {
			// exclude edges from power plant to itself
			if cost != 0.0 {
				// only need to go one way as shortestPaths has all combinations
				fromTo := edge{from, to, cost}
				adjList[from] = append(adjList[from], fromTo)
			}
		}
	}
	return adjList
}

func prims(powerPlants []loc, adjList map[loc][]edge) (float64, []edge) {
	// build visited map
	visited := make(map[loc]bool, len(powerPlants))

	// initialize heap with point 0 edges
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

type Heap []edge

func (h *Heap) Push(x interface{}) { *h = append(*h, x.(edge)) }
func (h *Heap) Pop() interface{} {
	previous, n := *h, h.Len()
	*h = previous[:n-1]
	return previous[n-1]
}
func (h Heap) Len() int           { return len(h) }
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h Heap) Less(i, j int) bool { return h[i].cost < h[j].cost }

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

// func nonTerminalNodes(southWest loc, northEast loc, numPoints float64) []loc {
// 	x := (northEast.x - southWest.x) / numPoints
// 	y := (northEast.y - southWest.y) / numPoints
// 	nonTerminal := make([]loc, 0)
// 	for i := southWest.x; i <= northEast.x; i += x {
// 		for j := southWest.y; j <= northEast.y; j += y {
// 			nonTerminal = append(nonTerminal, loc{i, j})
// 		}
// 	}
// 	return nonTerminal
// }

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
	// nonTerminal := nonTerminalNodes(loc{10, 20}, loc{40, 40}, 20)
	// allNodes := append(powerPlants, nonTerminal...)
	shortestPaths := dijkstras(powerPlants)

	adjList := metricClosure(powerPlants, shortestPaths)
	// fmt.Println(adjList)

	cost, edges := prims(powerPlants, adjList)

	// edges = steinerTree(powerPlants, edges, nonTerminal)

	plotCosts(powerPlants, edges)
	fmt.Println(cost, edges)
}
