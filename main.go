package main

import "fmt"

func main() {
	c := newCube()
	c.true_scramble(5)
	fmt.Println(c.String())
	n := newNode(c)
	sol := n.bfs()
	fmt.Println(sol)
	c.manyMoves(sol)
	fmt.Println(c.String())
}
