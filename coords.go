package main

func ScreenSpaceToGridWorldSpace(p Vec2D) (x int, y int) {
	return int(GRID_WORLD_DIMENSION * (p.X / float64(WINDOW_WIDTH))),
		int(GRID_WORLD_DIMENSION * (1.0 - p.Y/float64(WINDOW_HEIGHT)))
}

func ScreenSpaceToGridCellSpace(p Vec2D) (x int, y int) {
	return int(GRID_CELL_DIMENSION * (p.X / float64(WINDOW_WIDTH))),
		int(GRID_CELL_DIMENSION * (1.0 - p.Y/float64(WINDOW_HEIGHT)))
}

func GridWorldSpaceToScreenSpace(p Vec2D) (x int, y int) {
	return int(WINDOW_WIDTH * (p.X / float64(GRID_WORLD_DIMENSION))),
		int(WINDOW_HEIGHT * (1.0 - p.Y/float64(GRID_WORLD_DIMENSION)))
}

func GridCellSpaceToGridWorldSpace(p Position) Vec2D {
	return Vec2D{
		float64(p.X*GRIDCELL_WORLD_W + GRIDCELL_WORLD_W/2),
		float64(p.Y*GRIDCELL_WORLD_H + GRIDCELL_WORLD_H/2)}
}

func GridWorldSpaceToGridCellSpace(p Vec2D) Position {
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
