package main

import "fmt"

func main() {
	fmt.Println("hello")
}

type board struct {
	rows []string
	cols []string
}

// want to make the comparison operation really easy
// where BRWY == BRWY for inside and outside
type square struct {
	rotation int
	row      []string
	col      []string
}

// rbyw - red blue yellow white
func (sq *square) setup(s string) {
	// each type of rotation
	sq.row = []string{
		string(s[0] + s[1]),
	}

}
