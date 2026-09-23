package main

import "fmt"

func main() {
	c := newCube()
	//c.scramble(3)
	c.U()
	fmt.Println(c.String())
	n := newNode(c)
	sol := n.bfs()
	fmt.Println(sol)
	c.manyMoves(sol)
	fmt.Println(c.String())
}
