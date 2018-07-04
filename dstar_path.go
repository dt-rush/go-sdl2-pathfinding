// D* Algorithm, as described in Stentz, 1993
package main

// Initialize the path computer, calculting the path data for the various
// grid cells, which can later be updated
func (pc *DStarPathComputer) DStarPathInit(
	start Position, end Position) bool {

	// clear and set initial state
	pc.Clear()
	pc.start = start
	pc.end = end
	// Add end cell to OPEN list
	pc.H[end.X][end.Y] = 0
	pc.From[end.X][end.Y] = NOWHERE
	pc.OH.Add(end)
	// ProcessState() is repeatedly called until start is removed from the
	// OPEN list (ie. T(start) == CLOSED), or a value of -1 is
	// returned, at which point either the path has been constructed,
	// or does not exist, respectively
	kmin := 0
	for pc.T[start.X][start.Y] != CLOSED &&
		kmin != -1 {

		kmin = pc.ProcessState()
	}
	return kmin != -1
}

// Used by DStarPathInit
func (pc *DStarPathComputer) ProcessState() (kmin int) {
	cur, err := pc.MinState()
	if err != nil {
		return -1
	}
	kold := pc.K[cur.X][cur.Y]
	pc.T[cur.X][cur.Y] = CLOSED
	if cur == pc.start {
		return kold
	}
	// reduce H(cur) by lowest-cost neighbor if possible
	for _, ix := range ixs {
		nbr := Position{
			cur.X + ix[0],
			cur.Y + ix[1],
		}
		if !pc.Grid.InGrid(nbr) ||
			pc.Grid.Cells[nbr.X][nbr.Y] == OBSTACLE ||
			(ix[0]*ix[1] != 0 &&
				pc.Grid.Cells[cur.X][nbr.Y] == OBSTACLE ||
				pc.Grid.Cells[nbr.X][cur.Y] == OBSTACLE) {
			continue
		}
		pathCost := pc.H[cur.X][cur.Y] + pc.C(cur, nbr)
		if pc.T[nbr.X][nbr.Y] == CLOSED &&
			pc.H[nbr.X][nbr.Y] <= kold &&
			pc.H[cur.X][cur.Y] > pathCost {

			pc.From[cur.X][cur.Y] = nbr
			pc.H[cur.X][cur.Y] = pathCost
		}
	}
	// process each neighbor
	for _, ix := range ixs {
		nbr := Position{
			cur.X + ix[0],
			cur.Y + ix[1],
		}
		if !pc.Grid.InGrid(nbr) ||
			pc.Grid.Cells[nbr.X][nbr.Y] == OBSTACLE ||
			(ix[0]*ix[1] != 0 &&
				pc.Grid.Cells[cur.X][nbr.Y] == OBSTACLE ||
				pc.Grid.Cells[nbr.X][cur.Y] == OBSTACLE) {
			continue
		}
		pathCost := pc.H[cur.X][cur.Y] + pc.C(cur, nbr)

		if pc.T[nbr.X][nbr.Y] == NEW {
			pc.From[nbr.X][nbr.Y] = cur
			H := pathCost
			pc.H[nbr.X][nbr.Y] = H
			pc.P[nbr.X][nbr.Y] = H
			pc.Insert(nbr)
		} else {
			// propagate cost change along backpointer
			if pc.From[nbr.X][nbr.Y] == cur &&
				pc.H[nbr.X][nbr.Y] != pc.H[cur.X][cur.Y]+pc.C(cur, nbr) {
				if pc.T[nbr.X][nbr.Y] == OPEN {
					if pc.H[nbr.X][nbr.Y] < pc.P[nbr.X][nbr.Y] {
						pc.P[nbr.X][nbr.Y] = pc.H[nbr.X][nbr.Y]
					}
					pc.H[nbr.X][nbr.Y] = pathCost
				} else {
					pc.H[nbr.X][nbr.Y] = pathCost
					pc.P[nbr.X][nbr.Y] = pathCost
				}
				pc.Insert(nbr)
			} else {
				// reduce cost of neighbor if possible
				if pc.From[nbr.X][nbr.Y] != cur &&
					pc.H[nbr.X][nbr.Y] > pathCost {

					if pc.P[cur.X][cur.Y] >= pc.H[cur.X][cur.Y] {
						pc.From[nbr.X][nbr.Y] = cur
						pc.H[nbr.X][nbr.Y] = pathCost
						if pc.T[nbr.X][nbr.Y] == CLOSED {
							pc.P[nbr.X][nbr.Y] = pc.H[nbr.X][nbr.Y]
						}
						pc.Insert(nbr)
					} else {
						pc.P[cur.X][cur.Y] = pc.H[cur.X][cur.Y]
						pc.Insert(cur)
					}
				} else {
					// set up cost reduction by neighbor if possible
					reversePathCost := pc.H[nbr.X][nbr.Y] + pc.C(nbr, cur)
					if pc.From[nbr.X][nbr.Y] != cur &&
						pc.H[cur.X][cur.Y] > reversePathCost &&
						pc.T[nbr.X][nbr.Y] == CLOSED &&
						pc.H[nbr.X][nbr.Y] > kold {

						pc.P[nbr.X][nbr.Y] = pc.H[nbr.X][nbr.Y]
						pc.Insert(nbr)
					}
				}
			}
		}
	}
	return pc.GetKMin()
}

func (pc *DStarPathComputer) MinState() (Position, error) {
	return pc.OH.Pop()
}

func (pc *DStarPathComputer) Insert(s Position) {
	wasOpenAlready := pc.T[s.X][s.Y] == OPEN
	if wasOpenAlready {
		pc.OH.Modify(pc.HeapIX[s.X][s.Y], pc.H[s.X][s.Y])
	} else {
		pc.T[s.X][s.Y] = OPEN
		pc.OH.Add(s)
	}
}

func (pc *DStarPathComputer) GetKMin() int {
	if len(pc.OH.Arr) < 2 {
		return -1
	} else {
		return pc.K[pc.OH.Arr[1].X][pc.OH.Arr[1].Y]
	}
}

// cost of traversing p2 -> p1
func (pc *DStarPathComputer) C(p1, p2 Position) int {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	if dx*dy != 0 {
		return 14
	} else {
		return 10
	}
}
