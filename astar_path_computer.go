package main

import (
	"fmt"
)

// Used to associate data with a grid cell during computation
type Node struct {
	Pos       Position // position in grid
	From      *Node    // pointer to next
	WhichList int      // == N means OPEN, N + 1 means CLOSED
	G         int      // path cost
	H         int      // heuristic
	F         int      // path cost + heuristic
	HeapIX    int      // index in heap array
}

// Prints in format (k)[x, y]
func (n *Node) String() string {
	return fmt.Sprintf("(%d)%v", n.F, n.Pos)
}

type AStarPathComputer struct {
	Grid      *Grid
	OH        *NodeHeap
	N         int
	Nodes     [][]Node
	startNode *Node
	endNode   *Node
}

func NewAStarPathComputer(grid *Grid) *AStarPathComputer {
	// build nodes grid
	// NOTE: in array-speak, the "rows" are columns. It's just nicer to put
	// X as the first coordinate instead of Y
	nodes := make([][]Node, grid.W)
	for x := 0; x < grid.W; x++ {
		nodes[x] = make([]Node, grid.H)
		for y := 0; y < grid.H; y++ {
			nodes[x][y] = Node{Pos: Position{x, y}}
		}
	}
	// make node heap
	c := &AStarPathComputer{
		Grid:  grid,
		N:     0,
		Nodes: nodes,
		OH:    NewNodeHeap(),
	}
	return c
}

func (c *AStarPathComputer) AStarPath(start Position, end Position) (path []Position) {
	// clear the heap which contains leftover nodes from the last calculation
	c.OH.Clear()
	// increment N so WhichList works properly
	c.N += 2

	// set start and end nodes
	c.startNode = &c.Nodes[start.X][start.Y]
	c.endNode = &c.Nodes[end.X][end.Y]
	// add first node to OPEN heap
	*c.startNode = Node{
		Pos:       c.startNode.Pos,
		From:      nil,
		WhichList: c.N,
		G:         0,
		H:         ManhattanDistance(start, end),
	}
	c.OH.Add(c.startNode)

	// while open heap has elements...
	for c.OH.Len() > 0 {
		// pop from open heap and set as closed (whichlist == c.N + 1)
		cur, err := c.OH.Pop()
		// if err, we have exhausted all squares on open heap and found no path
		// return empty list
		if err != nil {
			return []Position{}
		}
		// set popped node to CLOSED
		cur.WhichList = c.N + 1
		// if the current cell is the end, we're here. build the return list
		if cur == c.endNode {
			path = make([]Position, 0)
			for cur != nil {
				path = append(path, cur.Pos)
				cur = cur.From
			}
			// return the path to the user
			return path
		}
		// else, we have yet to complete the path. So:
		// for each neighbor
		for _, delta := range deltas {
			nbrPos, dist, err := c.Grid.NbrOf(cur.Pos, delta)
			if err != nil {
				continue
			}
			nbr := &c.Nodes[nbrPos.X][nbrPos.Y]
			// compute g, h for the neighbor
			g := cur.G + dist
			h := ManhattanDistance(nbr.Pos, end)
			// don't consider this neighbor if the neighbor is in the closed
			// list *and* our g is greater or equal to its g score (we already
			// have a better way to get to it)
			isClosed := nbr.WhichList == c.N+1
			if isClosed && g >= nbr.G {
				continue
			}
			// if not on open heap, add it with "From" == cur
			isOpen := nbr.WhichList == c.N
			if !isOpen {
				// set From, G, H
				nbr.From = cur
				nbr.G = g
				nbr.H = h
				// set whichlist == OPEN
				nbr.WhichList = c.N
				// push to open heap
				c.OH.Add(nbr)
			} else {
				// if it *is* on the open heap already, check to see if
				// this is a better path to that square
				// on -> "open node"
				gAlready := nbr.G
				if g < gAlready {
					// if the open node could be reached better by
					// this path, set the g to the new lower g, set the
					// "From" reference to cur and fix up the heap because
					// we've changed the value of one of its elements
					nbr.From = cur
					// compute new F after setting new G and add to OPEN heap
					nbr.G = g
					nbr.F = nbr.G + nbr.H
					c.OH.Modified(nbr)
				}
			}
		}
	}
	return path
}

// Manhattan distance (times 10, since 10 = 1 orthogonal, 14 = 1 diagonal)
func ManhattanDistance(p1 Position, p2 Position) int {
	dx := p1.X - p2.X
	if dx < 0 {
		dx *= -1
	}
	dy := p1.Y - p2.Y
	if dy < 0 {
		dy *= -1
	}
	return 10 * (dx + dy)
}
