package main

import (
	"math/rand"
)

func initRcube() {
	// Corners
	createTransitionCorner()
	// Edge group 1
	createTransitionEdges1()
	// Edge group 2
	createTransitionEdges2()
}

var RCUBE_BUFFER struct {
	pcorner [][]uint16 // 40320
	pedges1 [][]uint32 // 665280
	pedges2 [][]uint32 // 665280
	ccorner [][]uint16 // 2187
	cedges1 [][]uint8  // 64
	cedges2 [][]uint8  // 64
}

// Corners
func createTransitionCorner() {
	for i := range uint(40320) {
		cube := unrankCorner(i)
		for m := range 18 {
			c := cube
			c.move(Move(m))
			rank := c.getCornerRanking()
			pcorner := uint16(rank / 2187)
			ccorner := uint16(rank % 2187)
			RCUBE_BUFFER.pcorner[i][m] = pcorner
			RCUBE_BUFFER.ccorner[i][m] = ccorner
		}
	}
}

// Edge group 1
func unrankEdges1(rank uint) Cube {
	edges_pos := [12]EdgePiece{}
	edges_rot := [12]EdgeRot{}
	rank_positions := rank / 64
	rank_orientations := rank % 64
	// add other edges or something other than zero for re-ranking
	for i := range 12 {
		edges_pos[i] = 11
	}
	avalible_locations := make([]int, 12)
	for i := range 12 {
		avalible_locations[i] = i
	}
	for i := range 6 {
		current_value := 12 - i
		val := rank_positions % uint(current_value)
		edges_pos[avalible_locations[val]] = EdgePiece(i)
		avalible_locations = append(avalible_locations[:val], avalible_locations[val+1:]...)
		rank_positions /= uint(current_value)
	}
	//fix pieces after
	a := 6
	for i := range 12 {
		if edges_pos[i] == 11 {
			edges_pos[i] = EdgePiece(a)
			a += 1
		}
	}
	//Orientations
	for x := range 6 {
		rot := rank_orientations % 2
		edges_rot[x] = EdgeRot(rot == 1)
		rank_orientations /= 2
	}

	//Initialize cube, corners dont matter here
	corner_pos := [8]CornerPiece{}
	corner_rot := [8]CornerRot{}
	return Cube{corner_piece: corner_pos, corner_rot: corner_rot, edge_piece: edges_pos, edge_rot: edges_rot}
}

func createTransitionEdges1() {
	for i := range uint(665280) {
		cube := unrankEdges1(i)
		for m := range 18 {
			c := cube
			c.move(Move(m))
			rank := c.getEdge1Ranking()
			pedges1 := uint16(rank / 64)
			cedges1 := uint16(rank % 64)
			RCUBE_BUFFER.pcorner[i][m] = pedges1
			RCUBE_BUFFER.ccorner[i][m] = cedges1
		}
	}
}

// Edge group 2
func unrankEdges2(rank uint) Cube {
	edges_pos := [12]EdgePiece{}
	edges_rot := [12]EdgeRot{}
	rank_positions := rank / 64
	rank_orientations := rank % 64
	// add other edges or something other than zero for re-ranking
	avalible_locations := make([]int, 12)
	for i := range 12 {
		avalible_locations[i] = 11 - i
	}
	for i := range 6 {
		current_value := 12 - i
		val := rank_positions % uint(current_value)
		edges_pos[avalible_locations[val]] = EdgePiece(11 - i)
		avalible_locations = append(avalible_locations[:val], avalible_locations[val+1:]...)
		rank_positions /= uint(current_value)
	}
	//fix pieces after
	a := 0
	for i := range 12 {
		if edges_pos[i] == 0 {
			edges_pos[i] = EdgePiece(a)
			a += 1
		}
	}
	//Orientations
	for x := range 6 {
		rot := rank_orientations % 2
		edges_rot[x+6] = EdgeRot(rot == 1)
		rank_orientations /= 2
	}

	//Initialize cube, corners dont matter here
	corner_pos := [8]CornerPiece{}
	corner_rot := [8]CornerRot{}
	return Cube{corner_piece: corner_pos, corner_rot: corner_rot, edge_piece: edges_pos, edge_rot: edges_rot}
}

func createTransitionEdges2() {
	for i := range uint(665280) {
		cube := unrankEdges2(i)
		for m := range 18 {
			c := cube
			c.move(Move(m))
			rank := c.getEdge2Ranking()
			pedges2 := uint16(rank / 64)
			cedges2 := uint16(rank % 64)
			RCUBE_BUFFER.pcorner[i][m] = pedges2
			RCUBE_BUFFER.ccorner[i][m] = cedges2
		}
	}
}

// Cube
type RCube struct {
	pcorner uint16 // 8! -> 10
	pedges1 uint32 // 12!/6! -> 20
	pedges2 uint32 // 12!/6! -> 20
	ccorner uint16 // 3^7 -> 12
	cedges1 uint8  // 2^6 -> 6
	cedges2 uint8  // 2^6 -> 6
}

// Solved Cube
func (c *Cube) toRCube() RCube {
	corner := c.getCornerRanking()
	edges1 := c.getEdge1Ranking()
	edges2 := c.getEdge2Ranking()
	pcorner := uint16(corner / 2187)
	ccorner := uint16(corner % 2187)
	pedges1 := uint32(edges1 / 64)
	cedges1 := uint8(edges1 % 64)
	pedges2 := uint32(edges2 / 64)
	cedges2 := uint8(edges2 % 64)
	rcube := RCube{pcorner, pedges1, pedges2, ccorner, cedges1, cedges2}
	return rcube
}

func (c *RCube) manyMoves(moves []Move) {
	for _, move := range moves {
		c.move(move)
	}
}

func (c *RCube) move(m Move) {
	c.pcorner = RCUBE_BUFFER.pcorner[c.pcorner][m]
	c.pedges1 = RCUBE_BUFFER.pedges1[c.pedges1][m]
	c.pedges2 = RCUBE_BUFFER.pedges2[c.pedges2][m]
	c.ccorner = RCUBE_BUFFER.ccorner[c.ccorner][m]
	c.cedges1 = RCUBE_BUFFER.cedges1[c.cedges1][m]
	c.cedges2 = RCUBE_BUFFER.cedges2[c.cedges2][m]
}

func (c *RCube) scramble(count uint) []Move {
	moves := make([]Move, count)
	for x := range count {
		move := Move(rand.Intn(18))
		moves[x] = move
		c.move(move)
	}
	return moves
}
