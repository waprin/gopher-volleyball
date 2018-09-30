package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
	"github.com/veandco/go-sdl2/gfx"
)

var lastUpdate time.Time

const (
	gravity = 0.05
)

type game struct {
	ball *ball
	slime *slime
}

type slime struct {
	center int32
	radius int32
	direction int32
}

type ball struct {
	y float32
	speed float32
	radius int32
}

func newGame () *game {
	return &game{
		ball: &ball{y: 200, radius: 20},
		slime: &slime{center:50, radius: 50},
	}
}

func (g *game) handleLeftTouch(state uint8) {
	if state == 1 {
		g.slime.direction = 2
	} else {
		g.slime.direction = 0
	}
}

func (g *game) handleRightTouch(state uint8) {
	if state == 1 {
		g.slime.direction = 1
	} else {
		g.slime.direction = 0
	}
}

func (g *game) tick() {
	g.ball.speed += gravity
	g.ball.y += g.ball.speed

	if g.slime.direction == 1 {
		g.slime.center += 2
	} else if g.slime.direction == 2 {
		g.slime.center -= 2
	}

	g.slime.touch(g.ball)

	if int32(g.ball.y) + g.ball.radius >= 600 {
		g.ball.speed *= -1
	}


}

func (s *slime) touch(_ *ball) {

}

func (s *slime) render(r *sdl.Renderer) {
	gfx.ArcColor(r, s.center, 600, int32(s.radius), 180, 360, sdl.Color{255, 0, 0, 255})
}

func (g *game) render(r *sdl.Renderer) {
	r.SetDrawColor(0, 0, 0, 255)
	r.Clear()
	g.slime.render(r)
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