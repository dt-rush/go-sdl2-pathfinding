// D* Algorithm, as described in Stentz, 1993
package main

import (
	"fmt"
)

// "Tags" which cells can have
// New: never visited
// Open: path cost being computed and neighbors being checked
// Closed: already computed
const (
	NEW    = iota
	OPEN   = iota
	CLOSED = iota
)

// Used to associate data with a grid cell during computation
type Node struct {
	Pos    Position // position in grid
	From   *Node    // "backpointer"
	T      int      // NEW, OPEN, CLOSED
	P      int      // "previous value" function (H before insert to OPEN)
	H      int      // "path cost" estimate
	K      int      // "key function" value (min of P and H)
	HeapIX int      // index in heap array
}

// set's the node's K = min (P, H)
func (n *Node) ComputeK() {
	if n.P < n.H {
		n.K = n.P
	} else {
		n.K = n.H
	}
}

// Prints in format (k)[x, y]
func (n *Node) String() string {
	return fmt.Sprintf("(%d)%v", n.K, n.Pos)
}

type FieldDStarPathComputer struct {
	// basic members needed for computing a path
	Grid      *Grid
	OH        *NodeHeap
	Nodes     [][]Node
	startNode *Node
	endNode   *Node
}

func NewFieldDStarPathComputer(grid *Grid) *FieldDStarPathComputer {
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
	return &FieldDStarPathComputer{
		Grid:  grid,
		Nodes: nodes,
		OH:    NewNodeHeap(),
	}
}

// used to clear all state
func (c *FieldDStarPathComputer) Clear() {
	c.startNode = nil
	c.endNode = nil
	c.OH.Clear()
	for x := 0; x < c.Grid.W; x++ {
		for y := 0; y < c.Grid.H; y++ {
			c.Nodes[x][y] = Node{Pos: Position{x, y}}
		}
	}
}

// Initialize the path computer, calculting the path data for the various
// grid cells, which can later be updated
func (c *FieldDStarPathComputer) FieldDStarPathInit(
	start Position, end Position) bool {

	// clear state
	c.Clear()
	// Add end cell to OPEN heap
	c.endNode = &c.Nodes[end.X][end.Y]
	*c.endNode = Node{
		Pos:  c.endNode.Pos,
		From: nil,
		H:    0,
	}
	c.OH.Add(c.endNode)
	// ProcessState() is repeatedly called until start is removed from the
	// OPEN list (ie. T(start) == CLOSED), or a value of -1 is
	// returned, at which point either the path has been constructed,
	// or does not exist, respectively
	c.startNode = &c.Nodes[start.X][start.Y]
	kmin := 0
	for c.startNode.T != CLOSED && kmin != -1 {
		kmin = c.ProcessState()
	}
	return kmin != -1
}

// Used by FieldDStarPathInit
func (c *FieldDStarPathComputer) ProcessState() (kmin int) {
	// pop min element from OPEN heap (minimizes by K, which is min of H and P)
	cur, err := c.MinState()
	if err != nil {
		return -1
	}
	kOld := cur.K
	cur.T = CLOSED
	if cur == c.startNode {
		return kOld
	}
	// reduce H(cur) by lowest-cost neighbor if possible
	for _, ix := range deltas {
		nbrPos, err := c.Grid.NbrOf(cur.Pos, ix)
		if err != nil {
			continue
		}
		nbr := &c.Nodes[nbrPos.X][nbrPos.Y]
		pathCost := cur.H + c.C(cur, nbr)
		if nbr.T == CLOSED &&
			nbr.H <= kOld &&
			cur.H > pathCost {

			cur.From = nbr
			cur.H = pathCost
		}
	}
	// process each neighbor
	for _, ix := range deltas {
		nbrPos, err := c.Grid.NbrOf(cur.Pos, ix)
		if err != nil {
			continue
		}
		nbr := &c.Nodes[nbrPos.X][nbrPos.Y]
		pathCost := cur.H + c.C(cur, nbr)

		if nbr.T == NEW {
			nbr.From = cur
			H := pathCost
			nbr.H = H
			nbr.P = H
			c.Insert(nbr)
		} else {
			// propagate cost change along backpointer
			if nbr.From == cur && nbr.H != pathCost {
				if nbr.T == OPEN {
					if nbr.H < nbr.P {
						nbr.P = nbr.H
					}
					nbr.H = pathCost
				} else {
					nbr.H = pathCost
					nbr.P = pathCost
				}
				c.Insert(nbr)
			} else {
				// reduce cost of neighbor if possible
				if nbr.From != cur && nbr.H > pathCost {
					if cur.P >= cur.H {
						nbr.From = cur
						nbr.H = pathCost
						if nbr.T == CLOSED {
							nbr.P = nbr.H
						}
						c.Insert(nbr)
					} else {
						cur.P = cur.H
						c.Insert(cur)
					}
				} else {
					// set up cost reduction by neighbor if possible
					reversePathCost := nbr.H + c.C(nbr, cur)
					if nbr.From != cur &&
						cur.H > reversePathCost &&
						nbr.T == CLOSED &&
						nbr.H > kOld {

						nbr.P = nbr.H
						c.Insert(nbr)
					}
				}
			}
		}
	}
	return c.GetKMin()
}

// returns the node with min K from the OPEN heap
func (c *FieldDStarPathComputer) MinState() (*Node, error) {
	return c.OH.Pop()
}

// inserts or reinserts an element to the OPEN heap at the appropriate position
func (c *FieldDStarPathComputer) Insert(n *Node) {
	n.ComputeK()
	if n.T == OPEN {
		c.OH.Modified(n)
	} else {
		c.OH.Add(n)
	}
}

// returns the current minimum K value on the OPEN heap
func (c *FieldDStarPathComputer) GetKMin() int {
	if len(c.OH.Arr) < 2 {
		return -1
	} else {
		return c.OH.Arr[1].K
	}
}

// cost of traversing p1 -> p2
func (c *FieldDStarPathComputer) C(n1 *Node, n2 *Node) int {
	dx := n1.Pos.X - n2.Pos.X
	dy := n1.Pos.Y - n2.Pos.Y
	if dx*dy != 0 {
		return 14
	} else {
		return 10
	}
}
