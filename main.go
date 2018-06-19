package main

import (
	"github.com/fatih/color"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"log"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	color.NoColor = false
}

func InitSDL() (*sdl.Renderer, *ttf.Font) {
	// init SDL
	sdl.Init(sdl.INIT_EVERYTHING)
	// init SDL TTF
	err := ttf.Init()
	if err != nil {
		panic(err)
	}
	r, rendererError := GetRenderer()
	if rendererError != 0 {
		log.Fatalf("failed to build renderer (reason %d)\n", rendererError)
	}
	f, err := GetFont()
	if err != nil {
		log.Fatalf("couldn't load font: %v\n", err)
	}
	return r, f
}

func main() {

	var exitcode int
	sdl.Main(func() {
		r, f := InitSDL()
		g := NewGame(r, f)
		exitcode = g.gameloop()
	})
	os.Exit(exitcode)

}
