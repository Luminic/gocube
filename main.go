package main

import "fmt"

func main() {
	cube := newCube()
	cube.scramble(5)
	rank := cube.getCornerRanking()
	other := unrankCorner(rank)
	fmt.Println(cube == other)
	fmt.Println(cube.String())
	fmt.Println(other.String())
	fmt.Println(rank == other.getCornerRanking())
	// err := initHeuristics()
	// if err != nil {
	// 	fmt.Println("Could not create file!")
	// }
	// c := newCube()
	// c.true_scramble(7)
	// fmt.Println(c.String())
	// n := newNode(c)
	// sol := n.IDA()
	// fmt.Println(sol)
	// c.manyMoves(sol)
	// fmt.Println(c.String())
}
