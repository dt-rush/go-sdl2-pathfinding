// Responds to input, draws the scene, handles quit event
package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"time"
)

// control input modes
const (
	MODE_PLACING_START = iota
	MODE_PLACING_END   = iota
)

type Game struct {
	grid      *Grid
	mode      int
	fdpc      *FieldDStarPathComputer
	fpsTicker *time.Ticker
	r         *sdl.Renderer
	f         *ttf.Font
}

func NewGame(r *sdl.Renderer, f *ttf.Font) *Game {
	grid := NewGrid(r)
	return &Game{
		grid:      grid,
		fdpc:      NewFieldDStarPathComputer(grid),
		fpsTicker: time.NewTicker(time.Millisecond * (1000 / FPS)),
		r:         r,
		f:         f,
	}
}

func (g *Game) gameloop() int {

gameloop:
	for {
		// draw if fps ticket ticked
		select {
		case _ = <-g.fpsTicker.C:
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

// handle quit, key, and mouse events
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

// check for a quit event or key escape/q
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

// handle keyboard input
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

// handle mouse input
func (g *Game) HandleMouseButtonEvents(me *sdl.MouseButtonEvent) {
	p := MouseButtonEventToGridCellPosition(me)
	if me.Type != sdl.MOUSEBUTTONDOWN {
		return
	}
	// place either start or end
	if g.mode == MODE_PLACING_START {
		// if placing start, clear any prior grid data
		g.grid.Clear()
		g.grid.SetStart(p)
	} else {
		g.grid.SetEnd(p)
	}
	// mode is toggled between start/end whenever a click event is processed
	g.mode = (g.mode + 1) % 2
	// if g.start and g.end are defined, compute the path
	if g.grid.start != nil && g.grid.end != nil {
		g.grid.path = g.grid.path[:0]
		pathExists := g.fdpc.FieldDStarPathInit(*g.grid.start, *g.grid.end)
		if pathExists {
			path := make([]Position, 0)
			cur := g.fdpc.startNode
			for cur != nil {
				path = append(path, cur.Pos)
				cur = cur.From
			}
			for i, _ := range path {
				if i != len(path)-1 {
					g.grid.path = append(g.grid.path,
						PositionPair{path[i], path[i+1]})
				}
			}
		}
	}
	g.grid.UpdateTexture()
}
