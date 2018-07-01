package main

import (
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

func (m *Grid) SetPath(path []Position) {
	for _, p := range path {
		val := m.Cells[p.X][p.Y]
		if val != START && val != END {
			m.Cells[p.X][p.Y] = PATH
		}
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

func (m *Grid) DrawPath() {
	for _, pp := range m.path {
		p1 := m.ToWorldSpace(pp.p1)
		p2 := m.ToWorldSpace(pp.p2)
		drawVector(m.r,
			p1,
			p2.Sub(p1),
			sdl.Color{R: 255, G: 255, B: 255})
	}
}

func (m *Grid) InGrid(p Position) bool {
	return p.X >= 0 && p.X < GRID_CELL_DIMENSION &&
		p.Y >= 0 && p.Y < GRID_CELL_DIMENSION
}

func (m *Grid) CellHasObstacle(x int, y int) bool {
	return m.Cells[x][y] == OBSTACLE
}

func (m *Grid) ToWorldSpace(p Position) Vec2D {
	return Vec2D{
		float64(p.X*GRIDCELL_WORLD_W + GRIDCELL_WORLD_W/2),
		float64(p.Y*GRIDCELL_WORLD_H + GRIDCELL_WORLD_H/2)}
}

func (m *Grid) ToGridSpace(p Vec2D) Position {
	x := int(p.X / GRIDCELL_WORLD_W)
	y := int(p.Y / GRIDCELL_WORLD_H)
	if x > GRID_CELL_DIMENSION-1 {
		x = GRID_CELL_DIMENSION - 1
	}
	if x < 0 {
		x = 0
	}
	if y > GRID_CELL_DIMENSION-1 {
		y = GRID_CELL_DIMENSION - 1
	}
	if y < 0 {
		y = 0
	}
	return Position{x, y}
}
