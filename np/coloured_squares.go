package main

import "fmt"

func main() {
	fmt.Println("hello")
	sq := square{}
	sq.setup("wybr")
	fmt.Println(sq.colours)
}

type board struct {
	locations [][]*square
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
