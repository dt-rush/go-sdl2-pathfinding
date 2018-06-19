package main

// G:			reference to the grid of cells we're pathing over
// OH			NodeHeap ("Open Heap") used to pop off the nodes with the lowest
//					F during search (open nodes)
// N: 			incremented each time we calculate (used to avoid having to
//					clear values in various arrays)
//
// WhichList:	2d array shadowing used to keep track of which list, open or
//					closed, the node is on
// From:		link to the prior node in path
// G:	 		path cost
// H:	 		heuristic
// F:			G + H
// HeapIX:		keeps track of the heap index of the element at this position
type PathComputer struct {
	Grid *Grid
	OH   *NodeHeap
	N    int
	// these 2D arrays store info about each node
	WhichList [][]int
	From      [][]Position
	G         [][]int
	H         [][]int
	F         [][]int
	HeapIX    [][]int
}

func NewPathComputer(grid *Grid) *PathComputer {

	// make 2D array rows
	// NOTE: in array-speak, the "rows" are columns. It's just nicer to put
	// X as the first coordinate instead of Y
	whichList := make([][]int, grid.W)
	from := make([][]Position, grid.W)
	g := make([][]int, grid.W)
	h := make([][]int, grid.W)
	f := make([][]int, grid.W)
	heapIX := make([][]int, grid.W)
	// make 2D array columns
	for x := 0; x < grid.W; x++ {
		whichList[x] = make([]int, grid.H)
		from[x] = make([]Position, grid.H)
		g[x] = make([]int, grid.H)
		h[x] = make([]int, grid.H)
		f[x] = make([]int, grid.H)
		heapIX[x] = make([]int, grid.H)
	}
	// make node heap
	pc := &PathComputer{
		Grid:      grid,
		N:         0,
		WhichList: whichList,
		From:      from,
		G:         g,
		H:         h,
		F:         f,
		HeapIX:    heapIX,
	}
	oh := NewNodeHeap(pc)
	pc.OH = oh
	return pc
}
