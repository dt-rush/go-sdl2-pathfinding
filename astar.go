package main

func (pc *AstarPathComputer) AstarPath(start Position, end Position) (path []Position) {
	// clear the heap which contains leftover nodes from the last calculation
	pc.OH.Clear()
	// increment N (easier than clearing arrays)
	// we increment by 2 since we use WhichList == pc.N for OPEN and
	// WhichList == pc.N + 1 for CLOSED
	pc.N += 2

	// add first node to open heap (whichlist == pc.N)
	pc.WhichList[start.X][start.Y] = pc.N
	// store a special value for the "From" of the first node
	pc.From[start.X][start.Y] = NOWHERE
	pc.G[start.X][start.Y] = 0
	pc.H[start.X][start.Y] = ManhattanDistance(start, end)
	pc.WhichList[start.X][start.Y] = pc.N
	pc.OH.Add(start)
	// while open heap has elements...
	for pc.OH.Len() > 0 {
		// pop from open heap and set as closed (whichlist == pc.N + 1)
		cur, err := pc.OH.Pop()
		pc.WhichList[cur.X][cur.Y] = pc.N + 1
		// if err, we have exhausted all squares on open heap and found no path
		// return empty list
		if err != nil {
			return []Position{}
		}
		// if the current cell is the end, we're here. build the return list
		if cur.X == end.X && cur.Y == end.Y {
			path = make([]Position, 0)
			for cur != NOWHERE {
				path = append(path, Position{cur.X, cur.Y})
				cur = pc.From[cur.X][cur.Y]
			}
			// return the path to the user
			return path
		}
		// else, we have yet to complete the path. So:
		// for each neighbor
		for _, ix := range ixs {
			// get the coordinates of the cell we will check the cost to
			// by applying an offset to cur's coordinates
			x := cur.X + ix[0]
			y := cur.Y + ix[1]
			// continue loop to next neighbor early if not in grid
			inGrid := x >= 0 && x < pc.Grid.W && y >= 0 && y < pc.Grid.H
			if !inGrid || pc.Grid.Cells[x][y] == OBSTACLE {
				continue
			}
			// continue to consider next neighbor if diagonal is flanked by
			// an obstacle
			isDiagonal := ix[0]*ix[1] != 0
			if isDiagonal && (pc.Grid.CellHasObstacle(cur.X+ix[0], cur.Y) ||
				pc.Grid.CellHasObstacle(cur.X, cur.Y+ix[1])) {
				continue
			}
			// dist is an integer expression of the distance from
			// cur to the neighbor cell we're looking at here.
			// if either x or y offset is 0, we're moving straight,
			// so put 10. Otherwise we're moving diagonal, so put 14
			// (these are 1 and sqrt(2), but made into integers for speed)
			var dist int
			if ix[0]*ix[1] == 0 {
				dist = 10
			} else {
				dist = 14
			}
			// compute g, h, f for the current cell
			g := pc.G[cur.X][cur.Y] + dist
			h := ManhattanDistance(Position{x, y}, end)
			// don't consider this neighbor if the neighbor is in the closed
			// list *and* our g is greater or equal to its g score (we already
			// have a better way to get to it)
			closed := pc.WhichList[x][y] == pc.N+1
			if closed && g >= pc.G[x][y] {
				continue
			}

			// if not on open heap, add it with "From" == cur
			open := pc.WhichList[x][y] == pc.N
			if !open {
				// set From, G, H
				pc.From[x][y] = cur
				pc.G[x][y] = g
				pc.H[x][y] = h
				// set whichlist == OPEN
				pc.WhichList[x][y] = pc.N
				// push to open heap
				pc.OH.Add(Position{x, y})
			} else {
				// if it *is* on the open heap already, check to see if
				// this is a better path to that square
				// on -> "open node"
				gAlready := pc.G[x][y]
				if g < gAlready {
					// if the open node could be reached better by
					// this path, set the g to the new lower g, set the
					// "From" reference to cur and fix up the heap because
					// we've changed the value of one of its elements
					pc.From[x][y] = cur
					pc.OH.Modify(pc.HeapIX[x][y], g)
				}
			}
		}
	}
	return path
}
