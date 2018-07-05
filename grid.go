package main

import (
	"errors"
	"github.com/veandco/go-sdl2/sdl"
)

type Grid struct {
	W     int
	H     int
	Cells [][]int

	// path
	path []PositionPair
	// renderer reference
	r *sdl.Renderer
	// screen texture
	st *sdl.Texture
}

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

func (m *Grid) ClearGrid() {
	m.path = m.path[:0]
	for x := 0; x < GRID_CELL_DIMENSION; x++ {
		for y := 0; y < GRID_CELL_DIMENSION; y++ {
			if m.Cells[x][y] == START ||
				m.Cells[x][y] == PATH ||
				m.Cells[x][y] == END {
				m.Cells[x][y] = EMPTY
			}
		}
	}
}

func (m *Grid) InGrid(p Position) bool {
	return p.X >= 0 && p.X < GRID_CELL_DIMENSION &&
		p.Y >= 0 && p.Y < GRID_CELL_DIMENSION
}

func (m *Grid) IsObstacle(p Position) bool {
	return m.Cells[p.X][p.Y] == OBSTACLE
}

func (m *Grid) NbrOf(cur Position, ix [2]int) (Position, error) {
	nbr := Position{
		cur.X + ix[0],
		cur.Y + ix[1],
	}
	if !m.InGrid(nbr) ||
		m.IsObstacle(nbr) ||
		// excludes cells which cross an obstacle on the corner
		(ix[0]*ix[1] != 0 &&
			m.Cells[cur.X][nbr.Y] == OBSTACLE ||
			m.Cells[nbr.X][cur.Y] == OBSTACLE) {
		return NOWHERE, errors.New("invalid neighbor")
	} else {
		return nbr, nil
	}
}

func (m *Grid) SetStart(start Position) {
	m.Cells[start.X][start.Y] = START
}

func (m *Grid) SetEnd(end Position) {
	m.Cells[end.X][end.Y] = END
}

func (m *Grid) UpdateTexture() {
	m.r.SetRenderTarget(m.st)
	defer m.r.SetRenderTarget(nil)

	m.r.SetDrawColor(0, 0, 0, 0)
	m.r.Clear()

	m.DrawGrid()
	m.DrawPath()
}

func (m *Grid) DrawGrid() {
	for x := 0; x < GRID_CELL_DIMENSION; x++ {
		for y := 0; y < GRID_CELL_DIMENSION; y++ {
			var c sdl.Color
			kind := m.Cells[x][y]
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

			drawRect(m.r,
				Rect2D{
					float64(x * GRIDCELL_WORLD_W),
					float64(y * GRIDCELL_WORLD_H),
					GRIDCELL_WORLD_W,
					GRIDCELL_WORLD_H},
				c)
		}
	}
}

// draws the path to the screen texture
func (m *Grid) DrawPath() {
	for _, pp := range m.path {
		p1 := GridCellSpaceToGridWorldSpace(pp.p1)
		p2 := GridCellSpaceToGridWorldSpace(pp.p2)
		drawVector(m.r,
			p1,
			p2.Sub(p1),
			sdl.Color{R: 255, G: 255, B: 255})
	}
}
