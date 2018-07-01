package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
)

func GetFont() (*ttf.Font, error) {
	return ttf.OpenFont("test.ttf", FONTSZ)
}

func GetRenderer() (*sdl.Renderer, int) {
	var err error
	var window *sdl.Window
	var renderer *sdl.Renderer

	window, err = sdl.CreateWindow(
		"GRIDMAP TEST",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return nil, 1
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create renderer: %s\n", err)
		return nil, 2
	}

	renderer.Clear()
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	return renderer, 0
}
