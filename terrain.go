package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
)

const (
	EMPTY    = 0
	OBSTACLE = iota
	START    = iota
	PATH     = iota
	END      = iota
)

var bgRed = color.New(color.BgRed).SprintFunc()
var bgBlack = color.New(color.BgBlack).SprintFunc()
var bgGreen = color.New(color.BgGreen).SprintFunc()
var bgWhite = color.New(color.BgWhite).SprintFunc()
var bgCyan = color.New(color.BgCyan).SprintFunc()

func CellString(val int) (rep string) {
	switch val {
	case EMPTY:
		rep = bgBlack("  ")
	case OBSTACLE:
		rep = bgRed("  ")
	case START:
		rep = bgGreen("  ")
	case PATH:
		rep = bgWhite("  ")
	case END:
		rep = bgCyan("  ")
	}
	return rep
}

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
				for _, ix := range ixs {
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

func PrintGrid(g Grid) {
	for y := g.H - 1; y >= 0; y-- {
		for x := 0; x < g.W; x++ {
			fmt.Printf("%s", CellString(g.Cells[x][y]))
		}
		fmt.Println()
	}
	fmt.Println()
}

func PrintGridWithPath(
	g Grid, start Position, end Position, path []Position) {

	inPath := make(map[Position]bool)
	for _, p := range path {
		inPath[p] = true
	}
	for y := g.H - 1; y >= 0; y-- {
		for x := 0; x < g.W; x++ {
			p := Position{x, y}
			if p == start {
				fmt.Printf("%s", CellString(START))
			} else if p == end {
				fmt.Printf("%s", CellString(END))
			} else if _, ok := inPath[Position{x, y}]; ok {
				fmt.Printf("%s", CellString(PATH))
			} else {
				fmt.Printf("%s", CellString(g.Cells[x][y]))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
