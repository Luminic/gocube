package main

import "fmt"

func main() {
	err := initHeuristics()
	if err != nil {
		fmt.Println("Could not create file!")
	}
	// c.true_scramble(7)
	// fmt.Println(c.String())
	// n := newNode(c)
	// sol := n.IDA()
	// fmt.Println(sol)
	// c.manyMoves(sol)
	// fmt.Println(c.String())
}
