package main

import (
	"errors"
	"github.com/veandco/go-sdl2/sdl"
)

// special value used for the "From" of the start node
var NOWHERE = Position{-1, -1}

// deltas = neighbor x, y offsets
//
//   Y ^
//     |
//     |    -1,  1     0,  1     1,  1
//     |
//     |    -1,  0     [cur]     1,  0
//     |
//     |    -1, -1     0, -1     1, -1
//     |
//     |
//      --------------------------------->
//                                       X
var deltas = [][2]int{
	[2]int{-1, 1},
	[2]int{0, 1},
	[2]int{1, 1},
	[2]int{-1, 0},
	[2]int{1, 0},
	[2]int{-1, -1},
	[2]int{0, -1},
	[2]int{1, -1},
}

type Grid struct {
	W     int
	H     int
	Cells [][]int
	start *Position
	end   *Position
	path  []PositionPair
	r     *sdl.Renderer
	st    *sdl.Texture
}

// Construct a new grid, along with its SDL texture, generating random terrain
func NewGrid(r *sdl.Renderer) *Grid {
	st, err := r.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET,
		WINDOW_WIDTH,
		WINDOW_HEIGHT)
	if err != nil {
		panic(err)
	}
	g := Grid{
		W:     GRID_CELL_DIMENSION,
		H:     GRID_CELL_DIMENSION,
		Cells: MakeTerrain(GRID_CELL_DIMENSION, GRID_CELL_DIMENSION),
		r:     r,
		st:    st,
	}
	g.UpdateTexture()
	return &g
}

// delete the saved path, start, and end data
func (g *Grid) Clear() {
	g.path = g.path[:0]
	if g.start != nil {
		g.Cells[g.start.X][g.start.Y] = EMPTY
		g.start = nil
	}
	if g.end != nil {
		g.Cells[g.end.X][g.end.Y] = EMPTY
		g.end = nil
	}
}

// tests if a position is in the grid bounds
func (g *Grid) InGrid(p Position) bool {
	return p.X >= 0 && p.X < GRID_CELL_DIMENSION &&
		p.Y >= 0 && p.Y < GRID_CELL_DIMENSION
}

// tests if a position contains an obstacle
func (g *Grid) IsObstacle(p Position) bool {
	return g.Cells[p.X][p.Y] == OBSTACLE
}

// returns the neighbor position given an offset 'delta' or error if not a valid
// neighbor (returns error on cross-corners)
func (g *Grid) NbrOf(cur Position, delta [2]int) (
	pos Position, err error) {
	nbr := Position{
		cur.X + delta[0],
		cur.Y + delta[1],
	}
	if !g.InGrid(nbr) ||
		g.IsObstacle(nbr) ||
		// excludes cells which cross an obstacle on the corner
		(delta[0]*delta[1] != 0 &&
			g.Cells[cur.X][nbr.Y] == OBSTACLE ||
			g.Cells[nbr.X][cur.Y] == OBSTACLE) {
		return NOWHERE, errors.New("invalid neighbor")
	} else {
		return nbr, nil
	}
}

func (g *Grid) SetStart(start Position) {
	g.start = &start
	g.Cells[start.X][start.Y] = START
}

func (g *Grid) SetEnd(end Position) {
	g.end = &end
	g.Cells[end.X][end.Y] = END
}

// redraw the texture according to current state
func (g *Grid) UpdateTexture() {
	g.r.SetRenderTarget(g.st)
	defer g.r.SetRenderTarget(nil)
	g.r.SetDrawColor(0, 0, 0, 0)
	g.r.Clear()
	g.DrawGrid()
	g.DrawPath()
}

// draw the grid cells (EMPTY, OBSTACLE, START, END) to `st`
func (g *Grid) DrawGrid() {
	for x := 0; x < GRID_CELL_DIMENSION; x++ {
		for y := 0; y < GRID_CELL_DIMENSION; y++ {
			var c sdl.Color
			kind := g.Cells[x][y]
			switch kind {
			case EMPTY:
				c = sdl.Color{R: 0, G: 0, B: 0}
			case OBSTACLE:
				c = sdl.Color{R: 255, G: 0, B: 0}
			case START:
				c = sdl.Color{R: 0, G: 255, B: 0}
			case END:
				c = sdl.Color{R: 0, G: 255, B: 255}
			}

			drawRect(g.r,
				Rect2D{
					float64(x * GRIDCELL_WORLD_W),
					float64(y * GRIDCELL_WORLD_H),
					GRIDCELL_WORLD_W,
					GRIDCELL_WORLD_H},
				c)
		}
	}
}

// draw the path to `st`
func (g *Grid) DrawPath() {
	for _, pp := range g.path {
		p1 := GridCellSpaceToGridWorldSpace(pp.p1)
		p2 := GridCellSpaceToGridWorldSpace(pp.p2)
		drawVector(g.r,
			p1,
			p2.Sub(p1),
			sdl.Color{R: 255, G: 255, B: 255})
	}
}
