package main

import (
	"math/rand"
)

const (
	EMPTY    = 0
	OBSTACLE = iota
	START    = iota
	END      = iota
)

// generates random terrain of grid cells
func MakeTerrain(w int, h int) [][]int {
	t := make([][]int, w)
	for x := 0; x < w; x++ {
		t[x] = make([]int, h)
		for y := 0; y < h; y++ {
			t[x][y] = EMPTY
		}
	}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if rand.Float64() < DENSITY {
				t[x][y] = OBSTACLE
				for _, ix := range deltas {
					xx := x + ix[0]
					yy := y + ix[1]
					if xx < 0 || xx > w-1 ||
						yy < 0 || yy > h-1 {
						continue
					}
					t[xx][yy] = OBSTACLE
				}
			}
		}
	}
	return t
}
