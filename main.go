package main

import "fmt"

func main() {
	c := newCube()
	for range 5 {
		c.U()
		c.R()
		c.UP()
		c.RP()
	}
	fmt.Println(c.String())
}
