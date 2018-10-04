package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
	"github.com/veandco/go-sdl2/gfx"
	"math"
)

var count = 0

const (
	gravity = 0.1
)

type game struct {
	ball *ball
	slime1 *slime
	slime2 *slime
	net *net
	width int32
	height int32

	slime1Pts, slime2Pts int
}

type slime struct {
	x float64
	y float64
	radius float64
	velocityX float64
	velocityY float64
	color sdl.Color
}

type ball struct {
	x float64
	y float64
	velocityX float64
	velocityY float64
	radius float64
}

type net struct {
	x int32
	y int32
	h int32
	w int32
}

func newGame (w, h int32) *game {
	var netHeight int32 = 40
	var netWidth int32 = 10
	return &game{
		ball: &ball{x:200, y: 300, radius: 10},
		slime1: &slime{x:200, y: 600, radius: 50, color: sdl.Color{255, 0, 0,255}},
		slime2: &slime{x:700, y: 600, radius: 50, color: sdl.Color{0, 255, 0, 255}},
		width: w,
		height: h,
		net: &net{x: w/2 - netWidth/2, y: h-netHeight, w: netWidth, h: netHeight},
		slime1Pts: 3,
		slime2Pts: 2,
	}
}

func (g *game) handlePlayer1LeftTouch(state uint8) {
	if state == 1 {
		g.slime1.velocityX = -6
	} else {
		g.slime1.velocityX = 0
	}
}

func (g *game) handlePlayer1RightTouch(state uint8) {
	if state == 1 {
		g.slime1.velocityX = 6
	} else {
		g.slime1.velocityX = 0
	}
}

func (g *game) handlePlayer1UpTouch(state uint8) {
	if state == 1 {
		g.slime1.velocityY = -3
	}
}

func (g *game) handlePlayer2LeftTouch(state uint8) {
	if state == 1 {
		g.slime2.velocityX = -6
	} else {
		g.slime2.velocityX = 0
	}
}

func (g *game) handlePlayer2RightTouch(state uint8) {
	if state == 1 {
		g.slime2.velocityX = 6
	} else {
		g.slime2.velocityX = 0
	}
}

func (g *game) handlePlayer2UpTouch(state uint8) {
	if state == 1 {
		g.slime2.velocityY = -3
	}
}

func (g *game) tick() {
	count++
	g.ball.velocityY += gravity

	g.slime1.velocityY += gravity
	g.slime2.velocityY += gravity

	g.ball.y += g.ball.velocityY
	g.ball.x += g.ball.velocityX

	g.slime1.x += g.slime1.velocityX
	g.slime1.y += g.slime1.velocityY
	if g.slime1.y >= float64(g.height) {
		g.slime1.velocityY = 0
		g.slime1.y = float64(g.height)
	}


	g.slime2.x += g.slime2.velocityX
	g.slime2.y += g.slime2.velocityY
	if g.slime2.y >= float64(g.height) {
		g.slime2.velocityY = 0
		g.slime2.y = float64(g.height)
	}

	// Slime1 collide with walls
	if g.slime1.x - g.slime1.radius < 0 {
		g.slime1.x = g.slime1.radius
	} else if (g.slime1.x + g.slime1.radius) >= float64(g.net.x) {
		g.slime1.x = float64(g.net.x) - g.slime1.radius
	}

	// Slime 2 collide with walls
	if g.slime2.x + g.slime2.radius >= float64(g.width) {
		g.slime2.x = float64(g.width) - g.slime2.radius
	} else if (g.slime2.x - g.slime2.radius) <= float64(g.net.x + g.net.w) {
		g.slime2.x = float64(g.net.x + g.net.w) + g.slime2.radius
	}

	g.slime1.touch(g.ball)
	g.slime2.touch(g.ball)

	g.checkWallsBall()
	g.checkNetBall()

	if int32(g.ball.y + g.ball.radius) >= 600 {
		g.ball.velocityY *= -1
	}
}

func (g *game) checkNetBall() {
	ballY := int32(g.ball.y) + int32(g.ball.radius)
	ballX := int32(g.ball.x)
	topOfNet := g.net.y

	if ballY > topOfNet {
		if ballX+int32(g.ball.radius) >= g.net.x && ballX - int32(g.ball.radius) < g.net.x + g.net.w {
			if g.ball.velocityY > 0 && ballY < (topOfNet + 10){
				g.ball.velocityY *= -1
				g.ball.y = float64(g.net.y) - g.ball.radius
			} else if ballX + int32(g.ball.radius) < g.net.x + g.net.w/2  {
				g.ball.x = float64(g.net.x) - g.ball.radius
				if g.ball.velocityX >= 0 {
					g.ball.velocityX *= -1
				}
			} else {
				g.ball.x = float64(g.net.x + g.net.w) + g.ball.radius
				if g.ball.velocityX < 0 {
					g.ball.velocityX *= -1
				}
			}
		}
	}
}

func (g *game) checkWallsBall() {
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


	MAX_VELOCITY_Y:=float64(6)
	MAX_VELOCITY_X:=float64(8)

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
			b.velocityY = -1 * ballVelocityY
		}

	}
}

func (s *slime) render(r *sdl.Renderer) {
	gfx.ArcColor(r, int32(s.x), int32(s.y), int32(s.radius), 180, 360, s.color)
}

func (g *game) renderScore(r *sdl.Renderer, x int32, points int) {
	amountToWin := 6
	curX := x
	xDiff := int32(40)

	for i := 0; i < points; i++ {
	gfx.FilledCircleColor(r, curX, 50, 10, sdl.Color{150, 255, 255, 255})
		curX += xDiff
	}
	for i := 0; i < int(amountToWin - points); i++ {
		gfx.CircleColor(r, curX, 50, 10, sdl.Color{150, 255, 255, 255})
		curX += xDiff
	}

}

func (n *net) render(r *sdl.Renderer) {
	rect := &sdl.Rect{int32(n.x), n.y, n.w, n.h}
	r.SetDrawColor(255, 255, 255, 255)
	r.DrawRect(rect)
}

func (g *game) render(r *sdl.Renderer) {
	r.SetDrawColor(0, 0, 0, 255)
	r.Clear()

	g.net.render(r)

	g.slime1.render(r)
	g.slime2.render(r)
	g.ball.render(r)

	g.renderScore(r, 25, g.slime1Pts)
	g.renderScore(r, 525, g.slime2Pts)

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