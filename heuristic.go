package main

import (
	"fmt"
	"os"
)

func initHeuristics() error {
	// Corners
	corner, err := os.ReadFile("corner_data.bin")
	if err != nil {
		corner = createHeuristicCorner()
		err = os.WriteFile("corner_data.bin", corner, 0644)
	}
	GLOBAL_BUFFER.corner = corner
	// Edge group 1
	edges1, err := os.ReadFile("edges1_data.bin")
	if err != nil {
		edges1 = createHeuristicEdge1()
		err = os.WriteFile("edges1_data.bin", edges1, 0644)
	}
	GLOBAL_BUFFER.edges1 = edges1
	// Edge group 2
	edges2, err := os.ReadFile("edges2_data.bin")
	if err != nil {
		edges2 = createHeuristicEdge2()
		err = os.WriteFile("edges2_data.bin", edges2, 0644)
	}
	GLOBAL_BUFFER.edges2 = edges2
	return err
}

var GLOBAL_BUFFER struct {
	corner []uint8
	edges1 []uint8
	edges2 []uint8
}

func (c *Cube) getHeuristic() uint8 {
	heuristics := [3]uint8{c.getCornerHeuristic(), c.getEdge1Heuristic(), c.getEdge2Heuristic()}
	max := uint8(0)
	for _, h := range heuristics {
		if max < h {
			max = h
		}
	}
	return max
}

// CORNER SECTION
func (c *Cube) getCornerHeuristic() uint8 {
	rank := c.getCornerRanking()
	value := GLOBAL_BUFFER.corner[rank/2]
	if rank%2 == 1 {
		value >>= 4
	} else {
		value &= 0xF
	}
	return value
}

func unrankCorner(rank uint) Cube {
	corner_pos := [8]CornerPiece{}
	corner_rot := [8]CornerRot{}
	rank_positions := rank / 2187
	rank_orientations := rank % 2187
	// positions
	avalible_pieces := make([]CornerPiece, 8)
	for i := range 8 {
		avalible_pieces[i] = CornerPiece(i)
	}
	for i := range 8 {
		pieces_left := 8 - i
		val := rank_positions % uint(pieces_left)
		corner_pos[i] = avalible_pieces[val]
		avalible_pieces = append(avalible_pieces[:val], avalible_pieces[val+1:]...)
		rank_positions /= uint(pieces_left)
	}
	//Orientations
	last := 0
	for x := range 7 {
		rot := rank_orientations % 3
		corner_rot[x] = CornerRot(rot)
		last = (last + int(rot)) % 3
		rank_orientations /= 3
	}
	corner_rot[7] = CornerRot(last)

	//Initialize cube, edges dont matter here
	edge_pos := [12]EdgePiece{}
	edge_rot := [12]EdgeRot{}
	return Cube{corner_piece: corner_pos, corner_rot: corner_rot, edge_piece: edge_pos, edge_rot: edge_rot}
}

func (c *Cube) getCornerRanking() uint {
	// Rank the positions
	rank_positions := uint(0)
	m := uint(1)
	// For each corner pos
	for corner_index, corner_pos := range c.corner_piece {
		count_lesser := uint(0)
		for _, other_corner := range c.corner_piece[corner_index+1:] {
			// Count the number of corners that are less than it
			if other_corner < corner_pos {
				count_lesser++
			}
		}
		// update rank
		rank_positions += m * count_lesser
		// Increase weight to be max of the current ranked corners
		m *= uint(8 - corner_index)
	}
	// Rank the orientations
	rank_orientations := uint(0)
	m = uint(1)
	for _, corner_rot := range c.corner_rot[:7] {
		rank_orientations += m * uint(corner_rot)
		m *= 3
	}
	return rank_positions*2187 + rank_orientations
}

func createHeuristicCorner() []uint8 {
	//8*7*6*5*4*3*2*1 * 3^7
	max_possible := 88179840
	// can fit 2 values into one uint8
	buffer := make([]uint8, max_possible/2)
	written := make([]bool, max_possible)
	// Keep track of current and next layer
	now := make([]bool, max_possible)
	next := make([]bool, max_possible)
	/*
		This function will start from the solved state
		and track how far away a state is from the end state
	*/
	now[0] = true
	written[0] = true
	i := 0
	depth := 0
	for i < max_possible {
		for rcube, search := range now {
			// If we want to search it
			if search {
				if rcube%2 == 0 {
					buffer[rcube/2] |= uint8(depth)
				} else {
					buffer[rcube/2] |= uint8(depth << 4)
				}
				if i%100000 == 0 {
					fmt.Println(i)
				}
				i++
				// Get our cube
				cube := unrankCorner(uint(rcube))
				for x := range 18 {
					move := Move(x)
					c := cube
					c.move(move)
					// add if not visited
					c_rank := c.getCornerRanking()
					if !written[c_rank] {
						written[c_rank] = true
						next[c_rank] = true
					}
				}
			}
		}
		depth++
		now = next
		next = make([]bool, max_possible)
	}
	return buffer
}

// EDGE1 SECTION
func (c *Cube) getEdge1Heuristic() uint8 {
	rank := c.getEdge1Ranking()
	value := GLOBAL_BUFFER.edges1[rank/2]
	if rank%2 == 1 {
		value >>= 4
	} else {
		value &= 0xF
	}
	return value
}

func (c *Cube) getEdge1Ranking() uint {
	// Rank the positions
	rank_positions := uint(0)
	m := uint(1)
	// For each edge pos
	for edge_index, edge_val := range c.edge_piece[:6] {
		count_lesser := uint(0)
		for _, other_edge := range c.edge_piece[edge_index+1:] {
			// Count the number of edges that are less than it
			if other_edge < edge_val {
				count_lesser++
			}
		}
		// update rank
		rank_positions += m * count_lesser
		// Increase weight to be max of the current ranked corners
		m *= uint(12 - edge_index)
	}
	// Rank the orientations
	rank_orientations := uint(0)
	m = uint(1)
	for _, edge_rot := range c.edge_rot[:6] {
		if edge_rot {
			rank_orientations += m
		}
		m *= 2
	}
	return rank_positions*64 + rank_orientations
}

func createHeuristicEdge1() []uint8 {
	//12*11*10*9*8*7 * 2^6
	max_possible := 42577920
	// can fit 2 values into one uint8
	buffer := make([]uint8, max_possible/2)
	written := make([]bool, max_possible)
	// Keep track of current and next layer
	now := make([]*Cube, max_possible)
	next := make([]*Cube, max_possible)
	/*
		This function will start from the solved state
		and track how far away a state is from the end state
	*/
	c := newCube()
	r := c.getEdge1Ranking()
	now[r] = &c
	written[r] = true
	i := 0
	depth := 0
	for i < max_possible {
		for rcube, cube := range now {
			// If we want to search it
			if cube != nil {
				if rcube%2 == 0 {
					buffer[rcube/2] |= uint8(depth)
				} else {
					buffer[rcube/2] |= uint8(depth << 4)
				}
				if i%100000 == 0 {
					fmt.Println(i)
				}
				i++
				for x := range 18 {
					move := Move(x)
					c := *cube
					c.move(move)
					// add if not visited
					c_rank := c.getEdge1Ranking()
					if !written[c_rank] {
						written[c_rank] = true
						next[c_rank] = &c
					}
				}
			}
		}
		depth++
		now = next
		next = make([]*Cube, max_possible)
	}
	return buffer
}

// EDGE2 SECTION
func (c *Cube) getEdge2Heuristic() uint8 {
	rank := c.getEdge2Ranking()
	value := GLOBAL_BUFFER.edges2[rank/2]
	if rank%2 == 1 {
		value >>= 4
	} else {
		value &= 0xF
	}
	return value
}

func (c *Cube) getEdge2Ranking() uint {
	// Rank the positions
	rank_positions := uint(0)
	m := uint(1)
	// For each edge pos
	edge_index := 11
	for edge_index >= 6 {
		edge_val := c.edge_piece[edge_index]
		count_greater := uint(0)
		for _, other_edge := range c.edge_piece[:edge_index] {
			// Count the number of edges that are greater than it
			if other_edge > edge_val {
				count_greater++
			}
		}
		// update rank
		rank_positions += m * count_greater
		// Increase weight to be max of the current ranked corners
		m *= uint(edge_index + 1)
		edge_index--
	}
	// Rank the orientations
	rank_orientations := uint(0)
	m = uint(1)
	for _, edge_rot := range c.edge_rot[6:] {
		if edge_rot {
			rank_orientations += m
		}
		m *= 2
	}
	return rank_positions*64 + rank_orientations
}

func createHeuristicEdge2() []uint8 {
	//12*11*10*9*8*7 * 2^6
	max_possible := 42577920
	// can fit 2 values into one uint8
	buffer := make([]uint8, max_possible/2)
	written := make([]bool, max_possible)
	// Keep track of current and next layer
	now := make([]*Cube, max_possible)
	next := make([]*Cube, max_possible)
	/*
		This function will start from the solved state
		and track how far away a state is from the end state
	*/
	c := newCube()
	r := c.getEdge2Ranking()
	now[r] = &c
	written[r] = true
	i := 0
	depth := 0
	for i < max_possible {
		for rcube, cube := range now {
			// If we want to search it
			if cube != nil {
				if rcube%2 == 0 {
					buffer[rcube/2] |= uint8(depth)
				} else {
					buffer[rcube/2] |= uint8(depth << 4)
				}
				if i%100000 == 0 {
					fmt.Println(i)
				}
				i++
				for x := range 18 {
					move := Move(x)
					c := *cube
					c.move(move)
					// add if not visited
					c_rank := c.getEdge2Ranking()
					if !written[c_rank] {
						written[c_rank] = true
						next[c_rank] = &c
					}
				}
			}
		}
		depth++
		now = next
		next = make([]*Cube, max_possible)
	}
	return buffer
}
