package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"github.com/veandco/go-sdl2/gfx"
	"time"
)

var lastUpdate time.Time


type game struct {
	leftX int32
}

func (g *game) handleLeftTouch() {
	fmt.Printf("moving left x to ", g.leftX + 10)
  g.leftX += 10
}

func (g *game) render(r *sdl.Renderer) {


/*	rect := sdl.Rect{g.leftX, 400, 250, 250}

	r.SetDrawColor(0, 0, 0, 255)
	r.Clear()
	r.SetDrawColor(255, 100, 0, 255)
	r.FillRect(&rect)
*/
	gfx.ArcColor(r, 350, 200, 20, 1, 360, sdl.Color{255, 255, 255, 255})
	r.Present()

}