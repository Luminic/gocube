package main

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

func cc(pos CornerPiece) bool {
	switch pos {
	case 0, 3, 5, 6:
		return true
	default:
		return false
	}
}

func (c *Cube) getCornerColor(pos CornerPiece, offset CornerRot) Colors {
	piece := c.corner_piece[pos]
	if cc(piece) == cc(pos) {
		return corner_colors[piece][(offset-c.corner_rot[piece])%3]
	} else {
		return corner_colors[piece][(-offset+c.corner_rot[piece])%3]
	}
}

func (c *Cube) String() string {
	c.U()
	rstring := ""
	rstring += color_codes[c.getCornerColor(0, 0)] + "  @"
	rstring += color_codes[c.getCornerColor(1, 0)] + "@\n"
	rstring += color_codes[c.getCornerColor(2, 0)] + "  @"
	rstring += color_codes[c.getCornerColor(3, 0)] + "@\n"

	rstring += color_codes[c.getCornerColor(0, 2)] + "@"
	rstring += color_codes[c.getCornerColor(2, 2)] + "@"
	rstring += color_codes[c.getCornerColor(2, 1)] + "@"
	rstring += color_codes[c.getCornerColor(3, 1)] + "@"
	rstring += color_codes[c.getCornerColor(3, 2)] + "@"
	rstring += color_codes[c.getCornerColor(1, 2)] + "@"
	rstring += color_codes[c.getCornerColor(1, 1)] + "@"
	rstring += color_codes[c.getCornerColor(0, 1)] + "@\n"

	rstring += color_codes[c.getCornerColor(4, 2)] + "@"
	rstring += color_codes[c.getCornerColor(6, 2)] + "@"
	rstring += color_codes[c.getCornerColor(6, 1)] + "@"
	rstring += color_codes[c.getCornerColor(7, 1)] + "@"
	rstring += color_codes[c.getCornerColor(7, 2)] + "@"
	rstring += color_codes[c.getCornerColor(5, 2)] + "@"
	rstring += color_codes[c.getCornerColor(5, 1)] + "@"
	rstring += color_codes[c.getCornerColor(4, 1)] + "@\n"

	rstring += color_codes[c.getCornerColor(6, 0)] + "  @"
	rstring += color_codes[c.getCornerColor(7, 0)] + "@\n"
	rstring += color_codes[c.getCornerColor(4, 0)] + "  @"
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

func (c *Cube) moveCorner(clist [4]CornerPiece) {
	cornerP := c.corner_piece[clist[3]]
	c.corner_piece[clist[3]] = c.corner_piece[clist[2]]
	c.corner_piece[clist[2]] = c.corner_piece[clist[1]]
	c.corner_piece[clist[1]] = c.corner_piece[clist[0]]
	c.corner_piece[clist[0]] = cornerP
}
func (c *Cube) rotCorner(clist [4]CornerPiece, r1 CornerRot, r2 CornerRot) {
	for i := range 4 {
		piece := c.corner_piece[clist[i]]
		switch c.corner_rot[piece] {
		case r1:
			c.corner_rot[piece] = r2
		case r2:
			c.corner_rot[piece] = r1
		}
	}
}

func (c *Cube) moveEdge(elist [4]EdgePiece) {
	edgeP := c.edge_piece[elist[3]]
	c.edge_piece[elist[3]] = c.edge_piece[elist[2]]
	c.edge_piece[elist[2]] = c.edge_piece[elist[1]]
	c.edge_piece[elist[1]] = c.edge_piece[elist[0]]
	c.edge_piece[elist[0]] = edgeP
}
func (c *Cube) rotEdge(clist [4]EdgePiece) {
	for i := range 4 {
		piece := c.edge_piece[clist[i]]
		c.edge_rot[piece] = !c.edge_rot[piece]
	}
}

// Rotations
func (c *Cube) U() {
	//Corners 0->1->3->2
	cval := [4]CornerPiece{0, 1, 3, 2}
	c.moveCorner(cval)
	//FB<->LR
	c.rotCorner(cval, WY_updown, WY_leftright)
	//Edges 0->2->3->1
	c.moveEdge([4]EdgePiece{0, 2, 3, 1})
	//No change in rot
}

func (c *Cube) D() {
	//Corners 4->5->7->6
	cval := [4]CornerPiece{4, 5, 7, 6}
	c.moveCorner(cval)
	//FB<->LR
	c.rotCorner(cval, WY_updown, WY_leftright)
	//Edges 8->10->11->9
	c.moveEdge([4]EdgePiece{8, 10, 11, 9})
	//No change in rot
}

func (c *Cube) L() {
	//Corners 0->2->6->4
	cval := [4]CornerPiece{0, 2, 6, 4}
	c.moveCorner(cval)
	//UD<->FB
	c.rotCorner(cval, WY_updown, WY_frontback)
	//Edges 1->6->9->4
	c.moveEdge([4]EdgePiece{1, 6, 9, 4})
	//No change in rot
}

func (c *Cube) R() {
	//Corners 1->5->7->3
	cval := [4]CornerPiece{1, 5, 7, 3}
	c.moveCorner(cval)
	//UD<->FB
	c.rotCorner(cval, WY_updown, WY_frontback)
	//Edges 2->5->10->7
	c.moveEdge([4]EdgePiece{2, 5, 10, 7})
	//No change in rot
}

func (c *Cube) F() {
	//Corners 2->3->7->6
	cval := [4]CornerPiece{2, 3, 7, 6}
	c.moveCorner(cval)
	//UD<->LR
	c.rotCorner(cval, WY_updown, WY_leftright)
	//Edges 3->7->11->6
	eval := [4]EdgePiece{3, 7, 11, 6}
	c.moveEdge(eval)
	//Swap the 4 egdes
	c.rotEdge(eval)
}

func (c *Cube) B() {
	//Corners 0->4->5->1
	cval := [4]CornerPiece{0, 4, 5, 1}
	c.moveCorner(cval)
	//UD<->LR
	c.rotCorner(cval, WY_updown, WY_leftright)
	//Edges 0->4->8->5
	eval := [4]EdgePiece{0, 4, 8, 5}
	c.moveEdge(eval)
	//Swap the 4 egdes
	c.rotEdge(eval)
}
