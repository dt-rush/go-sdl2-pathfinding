package main

// special value used for the "From" of the start node
var NOWHERE = Position{-1, -1}

// neighbor x, y offsets
//
//                                       X
//      --------------------------------->
//     |
//     |    -1,  1     0,  1     1,  1
//     |
//     |    -1,  0               1,  0
//     |
//     |    -1, -1     0, -1     1, -1
//     |
//  Y  v
//
//
var ixs = [][2]int{
	[2]int{-1, 1},
	[2]int{0, 1},
	[2]int{1, 1},
	[2]int{-1, 0},
	[2]int{1, 0},
	[2]int{-1, -1},
	[2]int{0, -1},
	[2]int{1, -1},
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
