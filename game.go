package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"time"
)

const (
	MODE_PLACING_START = iota
	MODE_PLACING_END   = iota
)

type Game struct {
	grid *Grid

	mode  int
	start *Position
	end   *Position

	// path computer
	pc *PathComputer

	r *sdl.Renderer
	f *ttf.Font
}

func NewGame(r *sdl.Renderer, f *ttf.Font) *Game {
	grid := NewGrid(r)
	return &Game{
		r:    r,
		f:    f,
		grid: grid,
		pc:   NewPathComputer(grid),
	}
}

func (g *Game) HandleKeyEvents(e sdl.Event) {
	switch e.(type) {
	case *sdl.KeyboardEvent:
		ke := e.(*sdl.KeyboardEvent)
		if ke.Type == sdl.KEYDOWN {
			if ke.Keysym.Sym == sdl.K_g {
				fmt.Println("pressed G")
			}
		}
	}
}

func (g *Game) HandleMouseButtonEvents(me *sdl.MouseButtonEvent) {
	p := MouseButtonEventToVec2D(me)
	if me.Type == sdl.MOUSEBUTTONDOWN {
		pos := g.grid.ToGridSpace(p)
		if g.mode == MODE_PLACING_START {
			g.grid.ClearGrid()
			g.start = &pos
			g.grid.SetStart(*g.start)
			g.end = nil
		} else {
			g.end = &pos
			g.grid.SetEnd(*g.end)
		}
		if g.start != nil && g.end != nil {
			g.grid.path = g.grid.path[:0]
			path := g.pc.AstarPath(*g.start, *g.end)
			for i, _ := range path {
				if i != len(path)-1 {
					g.grid.path = append(g.grid.path,
						PositionPair{path[i], path[i+1]})
				}
			}
		}
		g.mode = (g.mode + 1) % 2
		g.grid.UpdateTexture()
	}
}

func (g *Game) HandleQuit(e sdl.Event) bool {
	switch e.(type) {
	case *sdl.QuitEvent:
		return false
	case *sdl.KeyboardEvent:
		ke := e.(*sdl.KeyboardEvent)
		if ke.Keysym.Sym == sdl.K_ESCAPE ||
			ke.Keysym.Sym == sdl.K_q {
			return true
		}
	}
	return false
}

func (g *Game) HandleEvents() bool {
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		if g.HandleQuit(e) {
			return false
		}
		switch e.(type) {
		case *sdl.KeyboardEvent:
			g.HandleKeyEvents(e)
		case *sdl.MouseButtonEvent:
			g.HandleMouseButtonEvents(e.(*sdl.MouseButtonEvent))
		}
	}
	return true
}

func (g *Game) gameloop() int {

	fpsTicker := time.NewTicker(time.Millisecond * (1000 / FPS))

gameloop:
	for {
		// try to draw
		select {
		case _ = <-fpsTicker.C:
			sdl.Do(func() {
				g.r.Clear()
				g.r.Copy(g.grid.st, nil, nil)
				g.r.Present()
			})
		default:
		}

		// handle input
		if !g.HandleEvents() {
			break gameloop
		}

		// sleep
		sdl.Delay(1000 / (2 * FPS))
	}
	return 0
}
