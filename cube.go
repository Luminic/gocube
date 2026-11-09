package main

import (
	"math/rand"
)

// Red front, white top

// Corner values
type CornerPiece uint8

type Colors int8

const (
	White Colors = iota
	Yellow
	Red
	Orange
	Green
	Blue
)

var color_codes = [6]string{
	"\033[37m",
	"\033[33m",
	"\033[31m",
	"\033[38;2;255;165;0m",
	"\033[32m",
	"\033[34m",
}

var corner_colors = [8][3]Colors{
	{White, Orange, Green},
	{White, Orange, Blue},
	{White, Red, Green},
	{White, Red, Blue},
	{Yellow, Orange, Green},
	{Yellow, Orange, Blue},
	{Yellow, Red, Green},
	{Yellow, Red, Blue},
}

var edge_colors = [12][2]Colors{
	{White, Orange},
	{White, Green},
	{White, Blue},
	{White, Red},
	{Orange, Green},
	{Orange, Blue},
	{Red, Green},
	{Red, Blue},
	{Yellow, Orange},
	{Yellow, Green},
	{Yellow, Blue},
	{Yellow, Red},
}

func is_corner_cc(pos CornerPiece) bool {
	switch pos {
	case 0, 3, 5, 6:
		return true
	default:
		return false
	}
}

func (c *Cube) getCornerColor(pos CornerPiece, offset CornerRot) Colors {
	piece := c.corner_piece[pos]
	if is_corner_cc(piece) == is_corner_cc(pos) {
		num := ((int(offset)-int(c.corner_rot[piece]))%3 + 3) % 3
		return corner_colors[piece][num]
	} else {
		num := ((-int(offset)+int(c.corner_rot[piece]))%3 + 3) % 3
		return corner_colors[piece][num]
	}
}

func (c *Cube) getEdgeColor(pos EdgePiece, rot EdgeRot) Colors {
	piece := c.edge_piece[pos]
	if rot != c.edge_rot[piece] {
		return edge_colors[piece][0]
	} else {
		return edge_colors[piece][1]
	}
}

func (c *Cube) String() string {
	rstring := ""
	rstring += color_codes[c.getCornerColor(0, 0)] + "   @"
	rstring += color_codes[c.getEdgeColor(0, false)] + "@"
	rstring += color_codes[c.getCornerColor(1, 0)] + "@\n"
	rstring += color_codes[c.getEdgeColor(1, false)] + "   @"
	rstring += color_codes[White] + "@"
	rstring += color_codes[c.getEdgeColor(2, false)] + "@\n"
	rstring += color_codes[c.getCornerColor(2, 0)] + "   @"
	rstring += color_codes[c.getEdgeColor(3, false)] + "@"
	rstring += color_codes[c.getCornerColor(3, 0)] + "@\n"

	rstring += color_codes[c.getCornerColor(0, 2)] + "@"
	rstring += color_codes[c.getEdgeColor(1, true)] + "@"
	rstring += color_codes[c.getCornerColor(2, 2)] + "@"
	rstring += color_codes[c.getCornerColor(2, 1)] + "@"
	rstring += color_codes[c.getEdgeColor(3, true)] + "@"
	rstring += color_codes[c.getCornerColor(3, 1)] + "@"
	rstring += color_codes[c.getCornerColor(3, 2)] + "@"
	rstring += color_codes[c.getEdgeColor(2, true)] + "@"
	rstring += color_codes[c.getCornerColor(1, 2)] + "@"
	rstring += color_codes[c.getCornerColor(1, 1)] + "@"
	rstring += color_codes[c.getEdgeColor(0, true)] + "@"
	rstring += color_codes[c.getCornerColor(0, 1)] + "@\n"

	rstring += color_codes[c.getEdgeColor(4, true)] + "@"
	rstring += color_codes[Green] + "@"
	rstring += color_codes[c.getEdgeColor(6, true)] + "@"
	rstring += color_codes[c.getEdgeColor(6, false)] + "@"
	rstring += color_codes[Red] + "@"
	rstring += color_codes[c.getEdgeColor(7, false)] + "@"
	rstring += color_codes[c.getEdgeColor(7, true)] + "@"
	rstring += color_codes[Blue] + "@"
	rstring += color_codes[c.getEdgeColor(5, true)] + "@"
	rstring += color_codes[c.getEdgeColor(5, false)] + "@"
	rstring += color_codes[Orange] + "@"
	rstring += color_codes[c.getEdgeColor(4, false)] + "@\n"

	rstring += color_codes[c.getCornerColor(4, 2)] + "@"
	rstring += color_codes[c.getEdgeColor(9, true)] + "@"
	rstring += color_codes[c.getCornerColor(6, 2)] + "@"
	rstring += color_codes[c.getCornerColor(6, 1)] + "@"
	rstring += color_codes[c.getEdgeColor(11, true)] + "@"
	rstring += color_codes[c.getCornerColor(7, 1)] + "@"
	rstring += color_codes[c.getCornerColor(7, 2)] + "@"
	rstring += color_codes[c.getEdgeColor(10, true)] + "@"
	rstring += color_codes[c.getCornerColor(5, 2)] + "@"
	rstring += color_codes[c.getCornerColor(5, 1)] + "@"
	rstring += color_codes[c.getEdgeColor(8, true)] + "@"
	rstring += color_codes[c.getCornerColor(4, 1)] + "@\n"

	rstring += color_codes[c.getCornerColor(6, 0)] + "   @"
	rstring += color_codes[c.getEdgeColor(11, false)] + "@"
	rstring += color_codes[c.getCornerColor(7, 0)] + "@\n"
	rstring += color_codes[c.getEdgeColor(9, false)] + "   @"
	rstring += color_codes[Yellow] + "@"
	rstring += color_codes[c.getEdgeColor(10, false)] + "@\n"
	rstring += color_codes[c.getCornerColor(4, 0)] + "   @"
	rstring += color_codes[c.getEdgeColor(8, false)] + "@"
	rstring += color_codes[c.getCornerColor(5, 0)] + "@\n"

	rstring += "\033[0m"
	return rstring
}

const (
	WOG CornerPiece = iota
	WOB
	WRG
	WRB
	YOG
	YOB
	YRG
	YRB
)

type CornerRot uint8

const (
	WY_updown CornerRot = iota
	WY_frontback
	WY_leftright
)

// Edge values
type EdgePiece uint8

const (
	WO EdgePiece = iota
	WG
	WB
	WR
	OG
	OB
	RG
	RB
	YO
	YG
	YB
	YR
)

type EdgeRot bool

const WYthenRO_updown_or_RGRBOGOB = true
const WYthenRO_leftright_or_WRWOYRYO = false

// Cube
type Cube struct {
	corner_piece [8]CornerPiece
	corner_rot   [8]CornerRot
	edge_piece   [12]EdgePiece
	edge_rot     [12]EdgeRot
}

// Solved Cube
func newCube() Cube {
	corner_pos := [8]CornerPiece{WOG, WOB, WRG, WRB, YOG, YOB, YRG, YRB}
	corner_rot := [8]CornerRot{WY_updown, WY_updown, WY_updown, WY_updown, WY_updown, WY_updown, WY_updown, WY_updown}
	edge_pos := [12]EdgePiece{WO, WG, WB, WR, OG, OB, RG, RB, YO, YG, YB, YR}
	edge_rot := [12]EdgeRot{WYthenRO_updown_or_RGRBOGOB, WYthenRO_updown_or_RGRBOGOB, WYthenRO_updown_or_RGRBOGOB, WYthenRO_updown_or_RGRBOGOB, WYthenRO_updown_or_RGRBOGOB, WYthenRO_updown_or_RGRBOGOB, WYthenRO_updown_or_RGRBOGOB, WYthenRO_updown_or_RGRBOGOB, WYthenRO_updown_or_RGRBOGOB, WYthenRO_updown_or_RGRBOGOB, WYthenRO_updown_or_RGRBOGOB, WYthenRO_updown_or_RGRBOGOB}
	cube := Cube{corner_piece: corner_pos, corner_rot: corner_rot, edge_piece: edge_pos, edge_rot: edge_rot}
	return cube
}

type Move uint8

const (
	U Move = iota
	UP
	U2
	F
	FP
	F2
	L
	LP
	L2
	D
	DP
	D2
	B
	BP
	B2
	R
	RP
	R2
	NO_MOVE
)

func (c *Cube) manyMoves(moves []Move) {
	for _, move := range moves {
		c.move(move)
	}
}

func (c *Cube) moveCorner(clist [4]CornerPiece, reverse bool) {
	if reverse {
		cornerP := c.corner_piece[clist[0]]
		c.corner_piece[clist[0]] = c.corner_piece[clist[1]]
		c.corner_piece[clist[1]] = c.corner_piece[clist[2]]
		c.corner_piece[clist[2]] = c.corner_piece[clist[3]]
		c.corner_piece[clist[3]] = cornerP
	} else {
		cornerP := c.corner_piece[clist[3]]
		c.corner_piece[clist[3]] = c.corner_piece[clist[2]]
		c.corner_piece[clist[2]] = c.corner_piece[clist[1]]
		c.corner_piece[clist[1]] = c.corner_piece[clist[0]]
		c.corner_piece[clist[0]] = cornerP
	}
}
func (c *Cube) rotCorner(clist [4]CornerPiece, r CornerRot) {
	for i := range 4 {
		r1, r2 := (r+1)%3, (r+2)%3
		piece := c.corner_piece[clist[i]]
		switch c.corner_rot[piece] {
		case r1:
			c.corner_rot[piece] = r2
		case r2:
			c.corner_rot[piece] = r1
		}
	}
}

func (c *Cube) moveEdge(elist [4]EdgePiece, reverse bool) {
	if reverse {
		edgeP := c.edge_piece[elist[0]]
		c.edge_piece[elist[0]] = c.edge_piece[elist[1]]
		c.edge_piece[elist[1]] = c.edge_piece[elist[2]]
		c.edge_piece[elist[2]] = c.edge_piece[elist[3]]
		c.edge_piece[elist[3]] = edgeP
	} else {
		edgeP := c.edge_piece[elist[3]]
		c.edge_piece[elist[3]] = c.edge_piece[elist[2]]
		c.edge_piece[elist[2]] = c.edge_piece[elist[1]]
		c.edge_piece[elist[1]] = c.edge_piece[elist[0]]
		c.edge_piece[elist[0]] = edgeP

	}
}
func (c *Cube) rotEdge(clist [4]EdgePiece) {
	for i := range 4 {
		piece := c.edge_piece[clist[i]]
		c.edge_rot[piece] = !c.edge_rot[piece]
	}
}

func (c *Cube) moveCorner2(clist [4]CornerPiece) {
	cornerP := c.corner_piece[clist[3]]
	c.corner_piece[clist[3]] = c.corner_piece[clist[1]]
	c.corner_piece[clist[1]] = cornerP
	cornerP = c.corner_piece[clist[2]]
	c.corner_piece[clist[2]] = c.corner_piece[clist[0]]
	c.corner_piece[clist[0]] = cornerP
}

func (c *Cube) moveEdge2(elist [4]EdgePiece) {
	edgeP := c.edge_piece[elist[3]]
	c.edge_piece[elist[3]] = c.edge_piece[elist[1]]
	c.edge_piece[elist[1]] = edgeP
	edgeP = c.edge_piece[elist[2]]
	c.edge_piece[elist[2]] = c.edge_piece[elist[0]]
	c.edge_piece[elist[0]] = edgeP
}

// Corners
var cvals = [...][4]CornerPiece{
	{0, 1, 3, 2},
	{2, 3, 7, 6},
	{0, 2, 6, 4},
	{4, 6, 7, 5},
	{0, 4, 5, 1},
	{1, 5, 7, 3},
}

// Edges
var evals = [...][4]EdgePiece{
	{0, 2, 3, 1},
	{3, 7, 11, 6},
	{1, 6, 9, 4},
	{8, 9, 11, 10},
	{0, 4, 8, 5},
	{2, 5, 10, 7},
}

func (c *Cube) move(m Move) {
	v := m / 3
	t := m % 3

	cval := cvals[v]
	eval := evals[v]

	if t == 2 {
		//no rot corner
		c.moveCorner2(cval)
		c.moveEdge2(eval)
	} else {
		c.moveCorner(cval, t == 1)
		c.rotCorner(cval, CornerRot(v%3))
		c.moveEdge(eval, t == 1)

		// F and R
		if v == 1 || v == 4 {
			c.rotEdge(eval)
		}
	}
}

func (c *Cube) scramble(count uint) []Move {
	moves := make([]Move, count)
	for x := range count {
		move := Move(rand.Intn(18))
		moves[x] = move
		c.move(move)
	}
	return moves
}

func english(moves []Move) string {
	rep := [...]string{"U", "F", "L", "D", "B", "R"}
	other := [...]string{"", "P", "2"}
	str := ""
	for _, m := range moves {
		str += rep[m/3] + other[m%3] + " "
	}
	return str
}
