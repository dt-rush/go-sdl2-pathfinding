package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func MouseButtonEventToGridCellPosition(me *sdl.MouseButtonEvent) Position {
	x, y := ScreenSpaceToGridCellSpace(Vec2D{float64(me.X), float64(me.Y)})
	return Position{x, y}
}

func MouseMotionEventToGridCellPosition(me *sdl.MouseMotionEvent) Position {
	x, y := ScreenSpaceToGridCellSpace(Vec2D{float64(me.X), float64(me.Y)})
	return Position{x, y}
}
