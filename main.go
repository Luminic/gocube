package main

import (
	"fmt"
)

func main() {
	err := initHeuristics()
	initRcube()
	if err != nil {
		fmt.Println("Could not create file!")
	}
	var i uint
	for {
		fmt.Print("Enter num: ")
		fmt.Scan(&i)
		cc := newCube()
		c := cc.toRCube()
		seq := c.scramble(i)
		cc.manyMoves(seq)
		fmt.Println(cc.String())
		n := newNode(c)
		fmt.Println(english(seq))
		sol := n.IDA()
		fmt.Println(english(sol))
		c.manyMoves(sol)
		cc.manyMoves(sol)
		fmt.Println(cc.String())
	}
}

// func analyze() {
// 	//Create file
// 	f, err := os.Create("cpu.pprof")
// 	if err != nil {
// 		log.Fatalf("could not create CPU profile: %v", err)
// 	}
// 	defer f.Close()
// 	// Start analysis
// 	if err := pprof.StartCPUProfile(f); err != nil {
// 		log.Fatalf("could not start CPU profile: %v", err)
// 	}
// 	defer pprof.StopCPUProfile()
// 	c := newCube()
// 	c.scramble(17)
// 	n := newNode(c)
// 	n.IDA()
// 	// Run with go tool pprof -http=:8080 cpu.pprof
// }
