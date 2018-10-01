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
	gravity = 0.1
)

type game struct {
	ball *ball
	slime *slime
	width int32
	height int32
}

type slime struct {
	x float64
	y float64
	radius float64
	velocityX float64
	velocityY float64
}

type ball struct {
	x float64
	y float64
	velocityX float64
	velocityY float64
	radius float64
}

func newGame (w, h int32) *game {
	return &game{
		ball: &ball{x:350, y: 200, radius: 20},
		slime: &slime{x:50, y: 600, radius: 50},
		width: w,
		height: h,

	}
}

func (g *game) handleLeftTouch(state uint8) {
	if state == 1 {
		g.slime.velocityX = -5
	} else {
		g.slime.velocityX = 0
	}
}

func (g *game) handleRightTouch(state uint8) {
	if state == 1 {
		g.slime.velocityX = 5
	} else {
		g.slime.velocityX = 0
	}
}

func (g *game) tick() {
	count++
	g.ball.velocityY += gravity
	g.ball.y += g.ball.velocityY
	g.ball.x += g.ball.velocityX


	g.slime.x += g.slime.velocityX

	g.slime.touch(g.ball)

	g.checkWalls()

	if int32(g.ball.y + g.ball.radius) >= 600 {
		g.ball.velocityY *= -1
	}
}

func (g *game) checkWalls() {
	if g.ball.x <= 0{
		g.ball.x = 0
		g.ball.velocityX *= -1
	} else if int32(g.ball.x) >= g.width {
		g.ball.x = float64(g.width)
		g.ball.velocityX *= - 1
	}

	if g.ball.y <= 0 {
		g.ball.y = 0
		g.ball.velocityY *= -1
	}
}

func (b *ball) touchWalls() {

}

func (s *slime) touch(b *ball) {
	HEIGHT := float64(600)
	ballY := HEIGHT - b.y
	slimeY := HEIGHT - s.y

	dx := 2 * (b.x - s.x)
	dy := ballY - slimeY
	dist := math.Sqrt(float64(dx * dx + dy * dy))

	ballVelocityY := -1 * b.velocityY
	slimeVelocityY := -1 * s.velocityY

	dVelocityX := b.velocityX  - s.velocityX
	dVelocityY := ballVelocityY - slimeVelocityY


	MAX_VELOCITY_Y:=float64(8)
	MAX_VELOCITY_X:=float64(10)

	if dy > 0 && dist < float64(b.radius + s.radius) {
//		fmt.Printf("oldBallX %v\n", b.x)
//		fmt.Printf("oldBallY %v\n", ballY)


		b.x = s.x + (s.radius + b.radius) / 2 * (dx / dist)
		ballY = slimeY + s.radius + b.radius * (dy/dist)

//		fmt.Printf("newBallX %v\n", b.x)
//		fmt.Printf("newBallY %v\n", ballY)

		b.y = HEIGHT - ballY

		smth := (dx * dVelocityX + dy * dVelocityY) / dist
		if smth <= 0 {
			b.velocityX = s.velocityX -2 * dx * smth / dist
			ballVelocityY = s.velocityY - 2 * dy * smth / dist

			if math.Abs(b.velocityX) >= MAX_VELOCITY_X {
				if b.velocityX > 0 {
					b.velocityX = MAX_VELOCITY_X
				} else {
					b.velocityX = -1 * MAX_VELOCITY_X
				}
			}
			if math.Abs(ballVelocityY) >= MAX_VELOCITY_Y {
				if ballVelocityY > 0 {
					ballVelocityY = MAX_VELOCITY_Y
				} else {
					ballVelocityY = -1 * MAX_VELOCITY_Y
				}
			}
			fmt.Printf("b.velocityX is now %v\n", b.velocityX)
			fmt.Printf("b.velocityY is now %v\n", ballVelocityY)

			b.velocityY = -1 * ballVelocityY
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