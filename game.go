package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
	"github.com/veandco/go-sdl2/gfx"
	"math"
	"fmt"
	"github.com/veandco/go-sdl2/ttf"
	"github.com/veandco/go-sdl2/img"
)

var count = 0

const (
	gravity = 0.1
	initialSlime1X = 200
	initialSlime2X = 500
)

type gameMode int

const (
	playing gameMode = iota
	pointScored
	matchEnded
	intro
)

type opponentMode int

const (
	opponent2p = iota
	opponentAI
)

type color int

const (
	textWhite = iota
	textBlack
)

type aiState int

const (
	initialAiState = iota
	enemyCloseToNet
)

type game struct {
	ball *ball
	slime1 *slime
	slime2 *slime
	ai *ai
	net *net
	width int32
	height int32

	mode gameMode
	oppMode opponentMode

	slime1Pts, slime2Pts int
	lastPtPlayer1 bool

	// point scored
	endpointTimer time.Time

	backgroundTex *sdl.Texture
	font *ttf.Font
	smallFont *ttf.Font
}

type slime struct {
	x float64
	y float64
	radius float64
	velocityX float64
	velocityY float64
	color sdl.Color
	tex *sdl.Texture
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

type ai struct {
	movement float64
	state aiState
}

func newGame (r *sdl.Renderer, w, h int32) (*game, error) {
	bg, err := img.LoadTexture(r, "image/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	slime1, err := img.LoadTexture(r, "image/greenslimegopher.png")
	if err != nil {
		return nil, fmt.Errorf("could not load slime image: %v", err)
	}

	slime2, err := img.LoadTexture(r, "image/redslimegopher.png")
	if err != nil {
		return nil, fmt.Errorf("could not load slime image: %v", err)
	}


	f, err := ttf.OpenFont("fonts/OpenSans-Regular.ttf", 200)
	if err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	smallF, err := ttf.OpenFont("fonts/OpenSans-Regular.ttf", 12)
	if err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	var netHeight int32 = 40
	var netWidth int32 = 10



	return &game{
		ball: &ball{x:200, y: 250, radius: 10},
		slime1: &slime{
			x: initialSlime1X,
			y: float64(h),
			radius: 50,
			color: sdl.Color{255, 0, 0,255},
			tex: slime1,
		},
		slime2: &slime{
			x:550,
			y: float64(h),
			radius: 50,
			color: sdl.Color{0, 255, 0, 255},
			tex: slime2,
		},
		ai: &ai{
			movement: 0,
			state: enemyCloseToNet,
		},
		width: w,
		height: h,
		net: &net{x: w/2 - netWidth/2, y: h-netHeight, w: netWidth, h: netHeight},
		slime1Pts: 0,
		slime2Pts: 0,
		mode: intro,
		backgroundTex: bg,
		font: f,
		smallFont: smallF,
	}, nil
}

func (g *game) handlePlayer1LeftTouch(state uint8) {
	if state == 1 {
		g.slime1.velocityX = -5
	} else {
		g.slime1.velocityX = 0
	}
}

func (g *game) handlePlayer1RightTouch(state uint8) {
	if state == 1 {
		g.slime1.velocityX = 5
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
	if g.oppMode != opponent2p {
		return
	}
	if state == 1 {
		g.slime2.velocityX = -5
	} else {
		g.slime2.velocityX = 0
	}
}

func (g *game) handlePlayer2RightTouch(state uint8) {
	if g.oppMode != opponent2p {
		return
	}
	if state == 1 {
		g.slime2.velocityX = 5
	} else {
		g.slime2.velocityX = 0
	}
}

func (g *game) handlePlayer2UpTouch(state uint8) {
	if g.oppMode != opponent2p {
		return
	}
	if state == 1 {
		g.slime2.velocityY = -3
	}
}

func (g *game) handleSpaceBar(state uint8) {
	fmt.Printf("space bar pressed")
	if g.mode == matchEnded {
		fmt.Printf("space bar match endedd")
		g.resetPoint()
		g.slime2Pts = 0
		g.slime1Pts = 0
		g.mode = playing
	}
}

func (g *game) tick() {

	if g.oppMode == opponentAI {
		g.updateAI()
	}

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

	g.checkBallFloor()
}

func (g *game) checkBallFloor() {
	if int32(g.ball.y - g.ball.radius) >= g.height {
		if g.ball.x < float64(g.net.x) {
			g.slime2Pts++
			g.lastPtPlayer1 = false
		} else {
			g.slime1Pts++
			g.lastPtPlayer1 = true
		}

		if g.slime1Pts == 7 || g.slime2Pts == 7 {
			g.mode = matchEnded
			return
		}
		g.mode = pointScored
		g.endpointTimer = time.Now()
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

	if dy > 0 && dist < float64(b.radius + s.radius) && dist > 5 {
		b.x = s.x + (s.radius + b.radius) / 2 * (dx / dist)
		ballY = slimeY + s.radius + b.radius * (dy/dist)

		b.y = HEIGHT - ballY

		smth := (dx * dVelocityX + dy * dVelocityY) / dist
		if smth <= 0 {
			b.velocityX += s.velocityX -2 * dx * smth / dist
			ballVelocityY += s.velocityY - 2 * dy * smth / dist

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
//  ArcColor call might be useful for debugging to render physics arc more directly
	r.Copy(s.tex, nil, &sdl.Rect{int32(s.x-s.radius), int32(s.y-s.radius), int32(s.radius*2), int32(s.radius)})
}

func (g *game) renderBackground(r *sdl.Renderer) {
	r.Copy(g.backgroundTex, nil, nil)
}

func (g *game) renderScore(r *sdl.Renderer, x int32, points int) {
	amountToWin := 7
	curX := x
	xDiff := int32(40)

	for i := 0; i < points; i++ {
		gfx.FilledCircleColor(r, curX, 50, 10, sdl.Color{0, 0, 0, 255})
		curX += xDiff
	}
	for i := 0; i < int(amountToWin - points); i++ {
		gfx.CircleColor(r, curX, 50, 10, sdl.Color{0, 0, 0, 255})
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
	g.renderBackground(r)

	g.net.render(r)

	g.slime1.render(r)
	g.slime2.render(r)
	g.ball.render(r)

	g.renderScore(r, 25, g.slime1Pts)
	g.renderScore(r, 475, g.slime2Pts)
}

func (g *game) gameLoop(r  *sdl.Renderer) {
	if g.mode == pointScored {
		g.render(r)
		g.renderPoint(r)
		r.Present()
		if time.Since(g.endpointTimer) > 1 * time.Second {
			g.resetPoint()
			g.mode = playing
		}
	} else if g.mode == matchEnded {
		g.renderMatchOver(r)
		r.Present()
	} else if g.mode == playing {
		g.tick()
		g.render(r)
		r.Present()
	} else if g.mode == intro {
		g.renderIntro(r)
		r.Present()
	}
}

func (g *game) start(r *sdl.Renderer) {
	lastUpdate:=time.Now()
	go func() {
		for {
			diff := time.Since(lastUpdate)
			if diff > 10*time.Millisecond {
				g.gameLoop(r)
				lastUpdate = time.Now()
			}
		}
	}()
}

func (g *game) resetPoint() {
	g.ball.y = 250
	g.ball.x = 200

	g.ball.velocityY = 0
	g.ball.velocityX = 0

	g.slime1.velocityX = 0
	g.slime1.velocityY = 0
	g.slime1.x = initialSlime1X
	g.slime1.y = float64(g.height)

	g.slime2.velocityX = 0
	g.slime2.velocityY = 0
	g.slime2.x = initialSlime2X
	g.slime2.velocityY = float64(g.height)
}

func (g *game) writeText(r *sdl.Renderer, text string, rect *sdl.Rect, tc color, big bool) error {
	var c sdl.Color
	if tc == textBlack {
		c = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	} else if tc == textWhite {
		c = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	}
	var s *sdl.Surface
	var err error
	if (big) {
		s, err = g.font.RenderUTF8Solid(text, c)
	} else {
		s, err = g.smallFont.RenderUTF8Solid(text, c)
	}
	if err != nil {
		return fmt.Errorf("Could not create surface ")
	}
	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("Could not create texture ")
	}
	defer t.Destroy()

	if err := r.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("could not copy texture: ")
	}
	return nil
}

func (g *game) renderPoint(r *sdl.Renderer) error {
	var err error
	r.SetDrawColor(0, 0, 0, 255)
	if g.lastPtPlayer1 {
		err = g.writeText(r, "Player 1 scored!", &sdl.Rect{200, 100, 300, 200}, textBlack, true)
	} else {
		err = g.writeText(r, "Player 2 scored!", &sdl.Rect{200, 100, 300, 200}, textBlack, true)
	}
	return err
}

func (g *game) renderMatchOver(r *sdl.Renderer) error {
	r.SetDrawColor(255, 255, 255, 255)
	r.Clear()
	g.renderBackground(r)

	if g.slime1Pts == 7 {
		err := g.writeText(r, "Player 1 Won!", &sdl.Rect{200, 100, 300, 200}, textBlack, true)
		if err != nil {
			return err
		}
	} else {
		err := g.writeText(r, "Player 2 Won!", &sdl.Rect{200, 100, 300, 200}, textBlack, true)
		if err != nil {
			return err
		}
	}
	g.writeText(r,"Press spacebar to play again.", &sdl.Rect{200, 200, 300, 200}, textBlack, true)
	return nil
}



func (g *game) renderIntro(r *sdl.Renderer) error {
	r.SetDrawColor(255, 255, 255, 255)
	r.Clear()
	g.renderBackground(r)

	r.SetDrawColor(0, 0, 0, 255)

	rect := &sdl.Rect{int32(100), int32(100), int32(150), int32((100))}
	r.DrawRect(rect)
	r.FillRect(rect)
	g.writeText(r,"Play Vs AI!", &sdl.Rect{100, 100, 150, 50}, textWhite, false)

	rect2 := &sdl.Rect{int32(500), int32(100), int32(150), int32((100))}
	r.DrawRect(rect2)
	r.FillRect(rect2)
	g.writeText(r,"Play 2 Player!", &sdl.Rect{500, 100, 150, 50}, textWhite, false)

	g.writeText(r,"Welcome to Gopher Volleyball!", &sdl.Rect{200, 250, 300, 100}, textBlack, true)
	return nil
}


func (g *game) handleMouseUp(x int32, y int32) {
	if g.mode == intro {
		if x > 100 && x < (100 + 150) && y > 100 && y < (100 + 100) {
			g.oppMode = opponentAI
			g.mode = playing
		}
		if x > 500 && x < (500 + 150) && y > 100 && y < (100 + 100) {
			g.oppMode = opponent2p
			g.mode = playing
		}
	}
}


func (ball *ball) render(r *sdl.Renderer) {
	gfx.FilledCircleColor(r, int32(ball.x), int32(ball.y), 10, sdl.Color{255, 255, 255, 255})
}