package main

type Node struct {
	cube       Cube
	prev_state *Node
	last_turn  Move
	depth      uint
}

func newNode(c Cube) Node {
	return Node{cube: c, prev_state: nil, last_turn: NO_MOVE, depth: 0}
}

func (n *Node) get_children() []Node {
	// Check if its the root node
	var nodes []Node
	if n.last_turn == NO_MOVE {
		nodes = make([]Node, 18)
		for x := range 18 {
			move := Move(x)
			c := n.cube
			c.move(move)
			nodes[x] = Node{cube: c, prev_state: n, last_turn: move, depth: n.depth + 1}
		}
	} else {
		v := int(n.last_turn) / 3
		// Check if its in the first half
		if v < 3 {
			nodes = make([]Node, 12)
			i := 0
			for x := range 18 {
				// Disallow turns on opposite side
				if (x/3 != v) && (x/3 != v+3) {
					move := Move(x)
					c := n.cube
					c.move(move)
					nodes[i] = Node{cube: c, prev_state: n, last_turn: move, depth: n.depth + 1}
					i++
				}
			}
		} else {
			nodes = make([]Node, 15)
			i := 0
			for x := range 18 {
				if x/3 != v {
					move := Move(x)
					c := n.cube
					c.move(move)
					nodes[i] = Node{cube: c, prev_state: n, last_turn: move, depth: n.depth + 1}
					i++
				}
			}
		}
	}
	return nodes
}

func (n *Node) bfs() []Move {
	target_state := newCube()
	queue := []Node{*n}
	for {
		node := queue[0]
		if node.cube == target_state {
			// We did it!
			sol := []Move{}
			for node.prev_state != nil {
				sol = append([]Move{node.last_turn}, sol...)
				node = *node.prev_state
			}
			return sol
		}
		queue = append(queue[1:], node.get_children()...)
	}
}

func (n *Node) IDA() []Move {
	target_state := newCube()
	// Do Depth first search with max depth and heuristic
	for max_depth := range uint(20) {
		node := n.IDA_rec(&target_state, max_depth)
		if node != nil {
			// Solution found!
			sol := []Move{}
			for node.prev_state != nil {
				sol = append([]Move{node.last_turn}, sol...)
				node = node.prev_state
			}
			return sol
		}
	}
	return nil
}

func (n *Node) IDA_rec(target *Cube, max_depth uint) *Node {
	heuristic := n.cube.getHeuristic()
	if max_depth < n.depth+uint(heuristic) {
		return nil
	} else if n.cube == *target {
		return n
	}
	for _, child := range n.get_children() {
		potential_solution := child.IDA_rec(target, max_depth)
		if potential_solution != nil {
			return potential_solution
		}
	}
	return nil
}
