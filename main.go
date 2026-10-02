package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	err := initHeuristics()
	if err != nil {
		fmt.Println("Could not create file!")
	}
	c := newCube()
	for x := range 18 {
		t := c
		t.move(Move(x))
		fmt.Println(t.String())
	}
	// var i uint
	// for {
	// 	fmt.Print("Enter max turns: ")
	// 	fmt.Scan(&i)
	// 	c := newCube()
	// 	c.true_scramble(i)
	// 	fmt.Println(c.String())
	// 	n := newNode(c)
	// 	sol := n.IDA()
	// 	fmt.Println(sol)
	// 	c.manyMoves(sol)
	// 	fmt.Println(c.String())
	// }
}

func analyze() {
	//Create file
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatalf("could not create CPU profile: %v", err)
	}
	defer f.Close()
	// Start analysis
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatalf("could not start CPU profile: %v", err)
	}
	defer pprof.StopCPUProfile()
	c := newCube()
	c.true_scramble(17)
	n := newNode(c)
	n.IDA()
	// Run with go tool pprof -http=:8080 cpu.pprof
}
