package main

import "fmt"

func main() {
	c := newCube()
	c.true_scramble(7)
	fmt.Println(c.String())
	n := newNode(c)
	sol := n.IDA()
	fmt.Println(sol)
	c.manyMoves(sol)
	fmt.Println(c.String())
}
