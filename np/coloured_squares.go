package main

import "fmt"

func main() {
	fmt.Println("hello")
	sq := square{}
	sq.setup("wybr")
	fmt.Println(sq.colours)
}

type board struct {
	n         int
	locations [][]*square
	edge      string
}

func (b *board) outside(x int, y int) string {
	xy := []int{-1, 0, 1, 0, -1} // top-right-bottom-left
	outside := ""
	for i := 1; i < len(xy); i++ {
		nx, ny := x+xy[i-1], y+xy[i]
		// new coordinates are off board so needs to match edge
		if nx < 0 || nx >= b.n || ny < 0 || ny >= b.n {
			outside = outside + b.edge
		} else { // on board so get value
			// if not yet set can be anything
			if b.locations[nx][ny] != nil {
				outside += "*"
			} else {
				l := b.locations[nx][ny]
				l.colours[l.rotation]
			}
		}
	}
	return outside
}

// string are always in the order top-right-bottom-left
type square struct {
	rotation int
	colours  []string
}

func (sq *square) setup(s string) {
	sq.colours = []string{s}
	// add 3 rotations
	for i := 0; i < 3; i++ {
		prev := sq.colours[i]
		cur := prev[1:] + string(prev[0])
		sq.colours = append(sq.colours, cur)
	}
}
