package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
	"github.com/veandco/go-sdl2/gfx"
	"math"
	"fmt"
)

var count = 0

const (
	gravity = 0.05
)

type game struct {
	ball *ball
	slime *slime
}

type slime struct {
	x float32
	y float32
	radius float32
	direction int32
}

type ball struct {
	x float32
	y float32
	speed float32
	radius float32
}

func newGame () *game {
	return &game{
		ball: &ball{x:350, y: 200, radius: 20},
		slime: &slime{x:50, y: 600, radius: 50},
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
	count++
	g.ball.speed += gravity
	g.ball.y += g.ball.speed

	if g.slime.direction == 1 {
		g.slime.x += 3
	} else if g.slime.direction == 2 {
		g.slime.x -= 3
	}

	g.slime.touch(g.ball)

	if int32(g.ball.y + g.ball.radius) >= 600 {
		g.ball.speed *= -1
	}
}

func (s *slime) touch(b *ball) {
	FUDGE := 0
	dx := 2 * (b.x - s.x)
	dy := b.y - s.y
	dist := math.Sqrt(float64(dx * dx + dy * dy))

	if count % 10 == 0 {
		//fmt.Printf("b.x is %v s.x is %v dx is %v\n", b.x, s.x, dx)
	}
	if dy < 0 && dist < float64(b.radius + s.radius) && dist > float64(FUDGE) {
		fmt.Printf("Collission!")
		if b.speed > 0 {
			b.speed *= -1
		}
	}

}

func (s *slime) render(r *sdl.Renderer) {
	gfx.ArcColor(r, int32(s.x), int32(s.y), int32(s.radius), 180, 360, sdl.Color{255, 0, 0, 255})
}

func (g *game) render(r *sdl.Renderer) {
	r.SetDrawColor(0, 0, 0, 255)
	r.Clear()
	g.slime.render(r)
	g.ball.render(r)

	r.Present()
}

func (g *game) start(r *sdl.Renderer) {
	lastUpdate:=time.Now()
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
	gfx.ArcColor(r, int32(ball.x), int32(ball.y), int32(ball.radius), 1, 360, sdl.Color{255, 255, 255, 255})
}