package main

import "fmt"

func main() {
	err := initHeuristics()
	if err != nil {
		fmt.Println("Could not create file!")
	}
	var i uint
	for {
		fmt.Print("Enter max turns: ")
		fmt.Scan(&i)
		c := newCube()
		c.true_scramble(i)
		fmt.Println(c.String())
		n := newNode(c)
		sol := n.IDA()
		fmt.Println(sol)
		c.manyMoves(sol)
		fmt.Println(c.String())
	}
}
