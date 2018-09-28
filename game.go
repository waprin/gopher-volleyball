package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"time"
	"github.com/veandco/go-sdl2/gfx"
)

var lastUpdate time.Time

const (
	gravity = 0.05
)

type game struct {
	leftX int32
	ball *ball
}


type ball struct {
	y float32
	speed float32
	radius int32
}

func newGame () *game {
	return &game{
		ball: &ball{y: 200, radius: 20},
		leftX: 50,
	}
}

func (g *game) handleLeftTouch() {
	fmt.Printf("moving left x to ", g.leftX + 10)
    g.leftX += 10
}

func (g *game) tick() {
	g.ball.speed += gravity
	g.ball.y += g.ball.speed

	if int32(g.ball.y) + g.ball.radius >= 600 {
		g.ball.speed *= -1
	}

}

func (g *game) render(r *sdl.Renderer) {
	r.SetDrawColor(0, 0, 0, 255)
	r.Clear()
	r.SetDrawColor(255, 100, 0, 255)
	rect := sdl.Rect{g.leftX, 400, 250, 250}
	r.FillRect(&rect)

	g.ball.render(r)

	r.Present()
}

func (g *game) start(r *sdl.Renderer) {
	lastUpdate=time.Now()
	go func() {
		for {
			diff := time.Since(lastUpdate)
			if diff > 10*time.Millisecond {
				lastUpdate = time.Now()
				g.tick()
				g.render(r)
			}
		}
	}()
}


func (ball *ball) render(r *sdl.Renderer) {
	gfx.ArcColor(r, 350, int32(ball.y), int32(ball.radius), 1, 360, sdl.Color{255, 255, 255, 255})
}