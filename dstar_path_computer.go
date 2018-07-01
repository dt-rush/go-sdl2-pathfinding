package main

const (
	NEW    = iota
	OPEN   = iota
	CLOSED = iota
)

// G:			reference to the grid of cells we're pathing over
// OH			NodeHeap ("Open Heap") used to pop off the nodes with the lowest
//					F during search (open nodes)
// N: 			incremented each time we calculate (used to avoid having to
//					clear values in various arrays)
type DstarPathComputer struct {
	Grid *Grid
	OH   *NodeHeap

	start Position
	end   Position

	// these 2D arrays store info about each node
	From [][]Position // "backpointer"

	T [][]int // NEW, OPEN, CLOSED
	P [][]int // "previous value" function (H before insert to OPEN)
	H [][]int // "path cost" estimate
	K [][]int // "key function" value (min of P and H)

	HeapIX [][]int // index in heap array
}

func NewDstarPathComputer(grid *Grid) *DstarPathComputer {

	// make 2D array rows
	// NOTE: in array-speak, the "rows" are columns. It's just nicer to put
	// X as the first coordinate instead of Y
	from := make([][]Position, grid.W)
	t := make([][]int, grid.W)
	p := make([][]int, grid.W)
	h := make([][]int, grid.W)
	k := make([][]int, grid.W)
	heapIX := make([][]int, grid.W)
	// make 2D array columns
	for x := 0; x < grid.W; x++ {
		from[x] = make([]Position, grid.H)
		t[x] = make([]int, grid.H)
		p[x] = make([]int, grid.H)
		h[x] = make([]int, grid.H)
		k[x] = make([]int, grid.H)
		heapIX[x] = make([]int, grid.H)
	}
	// make node heap
	pc := &DstarPathComputer{
		Grid:   grid,
		From:   from,
		T:      t,
		P:      p,
		H:      h,
		K:      k,
		HeapIX: heapIX,
	}
	oh := NewNodeHeap(pc)
	pc.OH = oh
	return pc
}

func (pc *DstarPathComputer) Clear() {
	pc.start = NOWHERE
	pc.end = NOWHERE
	pc.OH.Clear()
	for x := 0; x < pc.Grid.W; x++ {
		for y := 0; y < pc.Grid.H; y++ {
			pc.From[x][y] = NOWHERE
			pc.T[x][y] = NEW
			pc.P[x][y] = 0
			pc.H[x][y] = 0
			pc.K[x][y] = 0
			pc.HeapIX[x][y] = 0
		}
	}
}
