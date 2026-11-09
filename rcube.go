package main

import (
	"fmt"
	"math/rand"
	"os"
	"unsafe"
)

type UInts interface {
	uint8 | uint16 | uint32
}

func getFile[T UInts](filename string, buffer *[]T, create func() []T) {
	data, err := os.ReadFile(filename)
	if err != nil {
		items := create()

		var byteSlice []byte
		if len(items) > 0 {
			// Zero-copy cast []T to []byte for writing
			byteSlice = unsafe.Slice((*byte)(unsafe.Pointer(&items[0])), len(items)*int(unsafe.Sizeof(items[0])))
		}

		if err := os.WriteFile(filename, byteSlice, 0644); err != nil {
			fmt.Println("Failed to write to: ", filename)
		}
		*buffer = items
	} else {
		var zero T
		elemSize := int(unsafe.Sizeof(zero))

		if len(data) == 0 {
			*buffer = []T{}
		}

		// Zero-copy reinterpret the file's byte buffer as []T
		*buffer = unsafe.Slice((*T)(unsafe.Pointer(&data[0])), len(data)/elemSize)
	}
}

func initRcube() {
	getFile("data/corner_t_p.bin", &RCUBE_BUFFER.pcorner, createTransitionCornerP)
	getFile("data/edges1_t_p.bin", &RCUBE_BUFFER.pedges1, createTransitionEdges1P)
	getFile("data/edges2_t_p.bin", &RCUBE_BUFFER.pedges2, createTransitionEdges2P)

	getFile("data/corner_t_m.bin", &RCUBE_BUFFER.swapcorners, createTransitionCornersRank)
	getFile("data/corner_t_c.bin", &RCUBE_BUFFER.ccorner, createTransitionCornerC)
	getFile("data/edges1_t_m.bin", &RCUBE_BUFFER.swapedges1, createTransitionEdges1Mask)
	getFile("data/edges2_t_m.bin", &RCUBE_BUFFER.swapedges2, createTransitionEdges2Mask)
}

func (c *RCube) getHeuristic() uint8 {
	max := c.getCornerHeuristic()
	other := c.getEdge1Heuristic()
	if max < other {
		max = other
	}
	other = c.getEdge2Heuristic()
	if max < other {
		max = other
	}
	return max
}

func (c *RCube) getCornerHeuristic() uint8 {
	rank := uint(c.pcorner)*2187 + uint(c.ccorner)
	value := GLOBAL_BUFFER.corner[rank/2]
	if rank%2 == 1 {
		value >>= 4
	} else {
		value &= 0xF
	}
	return value
}
func (c *RCube) getEdge1Heuristic() uint8 {
	rank := uint(c.pedges1)*64 + uint(c.cedges1)
	value := GLOBAL_BUFFER.edges1[rank/2]
	if rank%2 == 1 {
		value >>= 4
	} else {
		value &= 0xF
	}
	return value
}
func (c *RCube) getEdge2Heuristic() uint8 {
	rank := uint(c.pedges2)*64 + uint(c.cedges2)
	value := GLOBAL_BUFFER.edges2[rank/2]
	if rank%2 == 1 {
		value >>= 4
	} else {
		value &= 0xF
	}
	return value
}

var RCUBE_BUFFER struct {
	pcorner     []uint16 // 40320 * 18
	pedges1     []uint32 // 665280 * 18
	pedges2     []uint32 // 665280 * 18
	swapcorners []uint8  // 40320 * 18
	ccorner     []uint16 // 256 * 2187 * 18
	swapedges1  []uint8  // 665280 * 18
	swapedges2  []uint8  // 665280 * 18
}

// This function pre-calculates how one permutation moves to the next
func createTransitionCornerP() []uint16 {
	pbuffer := make([]uint16, 18*40320)
	for i := range uint(40320) {
		cube := unrankCorner(i * 2187)
		for m := range 18 {
			c := cube
			c.move(Move(m))
			rank := c.getCornerRanking()
			pcorner := uint16(rank / 2187)
			pbuffer[18*i+uint(m)] = pcorner
		}
	}
	return pbuffer
}

func nchoosek(n uint, k uint) uint {
	if n < k {
		return 0
	}
	o := n - k
	if o > k {
		o, k = k, o
	}
	p := uint(1)
	for n > k {
		p *= n
		n -= 1
	}
	for o > 1 {
		p /= o
		o -= 1
	}
	return p
}

func swapCornersRank(a *Cube, b *Cube) uint8 {
	val := uint8(1)
	exists := false
	left := 4
	//remap pieces
	var piece_positions [8]uint8
	for i, piece := range a.corner_piece {
		piece_positions[piece] = uint8(i)
	}
	for x := range 8 {
		if CornerPiece(x) != b.corner_piece[piece_positions[x]] {
			val += uint8(nchoosek(uint(7-x), uint(left)))
			left -= 1
			exists = true
		}
	}
	if !exists {
		return 0
	}
	return val
}

// This function pre-calculates which corners get their orientation changed depending on permutation and move
func createTransitionCornersRank() []uint8 {
	buffer := make([]uint8, 18*40320)
	for i := range uint(40320) {
		cube := unrankCorner(i * 2187)
		for m := range 18 {
			c := cube
			c.move(Move(m))
			val := swapCornersRank(&cube, &c)
			buffer[18*i+uint(m)] = val
		}
	}
	return buffer
}

// This function pre-calculates what a new cubes orienation is depending on current orientation, which corners move, and what move
func createTransitionCornerC() []uint16 {
	buffer := make([]uint16, 18*71*2187)
	for j := range uint(40320) {
		for i := range uint(2187) {
			cube := unrankCorner(j*2187 + i)
			for m := range 18 {
				c := cube
				c.move(Move(m))
				val := swapCornersRank(&cube, &c)
				buffer[18*71*i+18*uint(val)+uint(m)] = uint16(c.getCornerRanking() % 2187)
			}
		}
		if j%1000 == 0 {
			fmt.Println(j)
		}
	}
	return buffer
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

func createTransitionEdges1P() []uint32 {
	pbuffer := make([]uint32, 18*665280)
	for i := range uint(665280) {
		cube := unrankEdges1(i * 64)
		for m := range 18 {
			c := cube
			c.move(Move(m))
			rank := c.getEdge1Ranking()
			pedges1 := uint32(rank / 64)
			pbuffer[18*i+uint(m)] = pedges1
		}
	}
	return pbuffer
}
func createTransitionEdges1Mask() []uint8 {
	buffer := make([]uint8, 18*665280)
	for i := range uint(665280) {
		cube := unrankEdges1(i * 64)
		for m := range 18 {
			c := cube
			c.move(Move(m))
			//remap pieces
			var piece_positions [12]uint8
			for i, piece := range cube.edge_piece {
				piece_positions[piece] = uint8(i)
			}
			val := uint8(0)
			for x := range 6 {
				if (EdgePiece(x) != c.edge_piece[piece_positions[x]]) && (cube.edge_rot[x] != c.edge_rot[x]) {
					val |= 1 << x
				}
			}
			buffer[18*i+uint(m)] = val
		}
	}
	return buffer
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

func createTransitionEdges2P() []uint32 {
	pbuffer := make([]uint32, 18*665280)
	for i := range uint(665280) {
		cube := unrankEdges2(i * 64)
		for m := range 18 {
			c := cube
			c.move(Move(m))
			rank := c.getEdge2Ranking()
			pedges2 := uint32(rank / 64)
			pbuffer[18*i+uint(m)] = pedges2
		}
	}
	return pbuffer
}
func createTransitionEdges2Mask() []uint8 {
	buffer := make([]uint8, 18*665280)
	for i := range uint(665280) {
		cube := unrankEdges2(i * 64)
		for m := range 18 {
			c := cube
			c.move(Move(m))
			//remap pieces
			var piece_positions [12]uint8
			for i, piece := range cube.edge_piece {
				piece_positions[piece] = uint8(i)
			}
			val := uint8(0)
			for x := range 6 {
				if (EdgePiece(x+6) != c.edge_piece[piece_positions[x+6]]) && (cube.edge_rot[x+6] != c.edge_rot[x+6]) {
					val |= 1 << x
				}
			}
			buffer[18*i+uint(m)] = val
		}
	}
	return buffer
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

func (c *RCube) toCube() Cube {
	cc := unrankCorner(uint(c.pcorner)*2187 + uint(c.ccorner))
	ce1 := unrankEdges1(uint(c.pedges1)*64 + uint(c.cedges1))
	ce2 := unrankEdges2(uint(c.pedges2)*64 + uint(c.cedges2))
	for i := range 12 {
		cc.edge_rot[i] = ce1.edge_rot[i] || ce2.edge_rot[i]
		if ce1.edge_piece[i] < 6 {
			cc.edge_piece[i] = ce1.edge_piece[i]
		} else {
			cc.edge_piece[i] = ce2.edge_piece[i]
		}
	}
	return cc
}

func (c *RCube) manyMoves(moves []Move) {
	for _, move := range moves {
		c.move(move)
	}
}

func (c *RCube) move(m Move) {
	c.ccorner = RCUBE_BUFFER.ccorner[71*18*uint(c.ccorner)+18*uint(RCUBE_BUFFER.swapcorners[18*uint(c.pcorner)+uint(m)])+uint(m)]
	c.cedges1 ^= RCUBE_BUFFER.swapedges1[18*uint(c.pedges1)+uint(m)]
	c.cedges2 ^= RCUBE_BUFFER.swapedges2[18*uint(c.pedges2)+uint(m)]
	c.pcorner = RCUBE_BUFFER.pcorner[18*uint(c.pcorner)+uint(m)]
	c.pedges1 = RCUBE_BUFFER.pedges1[18*uint(c.pedges1)+uint(m)]
	c.pedges2 = RCUBE_BUFFER.pedges2[18*uint(c.pedges2)+uint(m)]
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
