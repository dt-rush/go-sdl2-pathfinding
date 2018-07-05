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

type Node struct {
	Pos    Position // position in grid
	From   *Node    // "backpointer"
	T      int      // NEW, OPEN, CLOSED
	P      int      // "previous value" function (H before insert to OPEN)
	H      int      // "path cost" estimate
	K      int      // "key function" value (min of P and H)
	HeapIX int      // index in heap array
}

// K = min (P, H)
func (n *Node) ComputeK() {
	if n.P < n.H {
		n.K = n.P
	} else {
		n.K = n.H
	}
}

// for pretty printing
func (n *Node) String() string {
	return fmt.Sprintf("%v(%d)", n.Pos, n.K)
}

type DStarPathComputer struct {
	// basic members needed for computing a path
	Grid  *Grid
	OH    *NodeHeap
	Nodes [][]Node
	start Position
	end   Position
}

func NewDStarPathComputer(grid *Grid) *DStarPathComputer {
	// make 2D array rows
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
	c := &DStarPathComputer{
		Grid:  grid,
		Nodes: nodes,
		OH:    NewNodeHeap(),
	}
	return c
}

// used to clear all state
func (c *DStarPathComputer) Clear() {
	c.start = NOWHERE
	c.end = NOWHERE
	c.OH.Clear()
	for x := 0; x < c.Grid.W; x++ {
		for y := 0; y < c.Grid.H; y++ {
			c.Nodes[x][y] = Node{Pos: Position{x, y}}
		}
	}
}

// Initialize the path computer, calculting the path data for the various
// grid cells, which can later be updated
func (c *DStarPathComputer) DStarPathInit(
	start Position, end Position) bool {

	// clear and set initial state
	c.Clear()
	c.start = start
	c.end = end
	// Add end cell to OPEN list
	endNode := &c.Nodes[end.X][end.Y]
	endNode.H = 0
	endNode.From = nil
	c.OH.Add(endNode)
	// ProcessState() is repeatedly called until start is removed from the
	// OPEN list (ie. T(start) == CLOSED), or a value of -1 is
	// returned, at which point either the path has been constructed,
	// or does not exist, respectively
	startNode := &c.Nodes[start.X][start.Y]
	kmin := 0
	for startNode.T != CLOSED && kmin != -1 {
		kmin = c.ProcessState()
	}
	return kmin != -1
}

// Used by DStarPathInit
func (c *DStarPathComputer) ProcessState() (kmin int) {
	// pop min element from OPEN heap (minimizes by K, which is min of H and P)
	cur, err := c.MinState()
	if err != nil {
		return -1
	}
	kOld := cur.K
	cur.T = CLOSED
	if cur == c.start {
		return kOld
	}
	// reduce H(cur) by lowest-cost neighbor if possible
	for _, ix := range ixs {
		nbr, err := c.Grid.NbrOf(cur, ix)
		if err != nil {
			continue
		}
		pathCost := cur.H + c.C(cur, nbr)
		if nbr.T == CLOSED &&
			nbr.H <= kOld &&
			cur.H > pathCost {

			cur.From = nbr
			cur.H = pathCost
		}
	}
	// process each neighbor
	for _, ix := range ixs {
		nbr := Position{
			cur.X + ix[0],
			cur.Y + ix[1],
		}
		if !c.Grid.InGrid(nbr) ||
			c.Grid.Cells[nbr.X][nbr.Y] == OBSTACLE ||
			(ix[0]*ix[1] != 0 &&
				c.Grid.Cells[cur.X][nbr.Y] == OBSTACLE ||
				c.Grid.Cells[nbr.X][cur.Y] == OBSTACLE) {
			continue
		}
		pathCost := cur.H + c.C(cur, nbr)

		if nbr.T == NEW {
			nbr.From = cur
			H := pathCost
			nbr.H = H
			nbr.P = H
			c.Insert(nbr)
		} else {
			// propagate cost change along backpointer
			if nbr.From == cur &&
				nbr.H != pathCost {
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
				if nbr.From != cur &&
					nbr.H > pathCost {

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

func (c *DStarPathComputer) MinState() (*Node, error) {
	return c.OH.Pop()
}

func (c *DStarPathComputer) Insert(n *Node) {
	n.ComputeK()
	if n.T == OPEN {
		c.OH.Modified(n)
	} else {
		c.OH.Add(n)
	}
}

func (c *DStarPathComputer) GetKMin() int {
	if len(c.OH.Arr) < 2 {
		return -1
	} else {
		return c.OH.Arr[1].K
	}
}

// cost of traversing p2 -> p1
func (c *DStarPathComputer) C(p1, p2 Position) int {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	if dx*dy != 0 {
		return 14
	} else {
		return 10
	}
}

func (c *DStarPathComputer) From(p Position) *Position {
	return &c.from[p.X][p.Y]
}
func (c *DStarPathComputer) T(p Position) *int {
	return &c.t[p.X][p.Y]
}
func (c *DStarPathComputer) P(p Position) *int {
	return &c.p[p.X][p.Y]
}
func (c *DStarPathComputer) H(p Position) *int {
	return &c.h[p.X][p.Y]
}
func (c *DStarPathComputer) K(p Position) *int {
	return &c.k[p.X][p.Y]
}
func (c *DStarPathComputer) HeapIX(p Position) *int {
	return &c.heapIX[p.X][p.Y]
}
