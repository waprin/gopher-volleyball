package main

import (
	"math"
	"math/rand"
)

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func (g *game) randomJump(percent float64) {
	// do a jump `percent` percent of the time

	// if already in the air, don't jump again
	if g.slime2.y != float64(g.height) {
		return
	}
	if rand.Float64() < percent {
		g.slime2.velocityY = -5
	}
}

func (g *game) moveAIToNet() {
	g.slime2.velocityX = -2
}

func (g *game) moveAIToWall() {
	g.slime2.velocityX = 2
}

func (g *game) makeAIJump() {
	g.slime2.velocityY = -5
}


// calculateXWhenBallBelow returns the X value of the ball
// when it will be below the given yLimit y value
func (g *game) calculateXWhenBallBelow(yLimit float64) float64 {
	frameCount := countFramesUntilBelow(g.ball.y, g.ball.velocityY, yLimit)
	movingBallX := g.ball.x
	for i := 0; i < frameCount; i++ {
		movingBallX += g.ball.velocityX
	}

	//  fmt.Printf("got x when ball below %v\n", movingBallX)
	localMax := math.Max(movingBallX, 0)
	return math.Min(localMax, float64(g.width))
}

func (g *game) aiIsServing() bool {
	return g.ball.velocityX == 0 && g.ball.x == initialSlime2X
}

func (g *game) aiServe() {
	if g.ai.state == notServing {
		if rand.Float64() < .5 {
			g.ai.state = serve1
		} else {
			g.ai.state = serve2
		}
	}
	if g.ai.state == serve1 {
		if g.ball.y > 200 && g.ball.velocityY > 3 {
			g.moveAIToWall()
			g.makeAIJump()
		}
	} else if g.ai.state == serve2 {
		if g.ball.velocityY < -3 && g.slime2.x < 660 {
			g.slime2.velocityX = 2
		}
		if g.slime2.x >= 660 {
			g.slime2.velocityX = 0
		}

		if g.slime2.velocityY == 3  && g.slime2.x != 600 {
			g.slime2.velocityY = -5
		}

		if g.ball.velocityY >  3 && g.slime2.y != 0 && g.slime2.x >= 660 {
			g.slime2.velocityX = -2
		}
	}
}



func (g *game) updateAI() {
	xWhenBallBelow475 := g.calculateXWhenBallBelow(475)
	var closeEnough float64
	closeEnough = 20 // how close the slime should move to where the ball will be

	if g.ball.x < float64(g.width/ 2) {
		g.ai.state = notServing
	}

	if g.aiIsServing() {
		g.aiServe()
		return
	}

	if xWhenBallBelow475 < float64(g.width) / 2 { // if balls' trajectory is towards opponents side, just get centered
	    // handle return logic here

		differenceFromCenter := int32(g.slime2.x) - 600
		if float64(Abs(differenceFromCenter)) < closeEnough {
			g.slime2.velocityX = 0
		} else if differenceFromCenter >= 20 {
			//console.log('moveToNet 2');
			g.slime2.velocityX = -2
		} else {
			g.slime2.velocityX = 2
		}
		return
	}

	if math.Abs(g.slime2.x -xWhenBallBelow475) <=  closeEnough || g.slime2.x == 700 && g.ball.x > 700 {
		// if the AI is close to where the ball will be, do some random jumps
		if (g.slime2.x >= 600 && g.ball.x > 530 ||
			g.slime2.x <= 580 && g.ball.x < 530) && math.Abs(g.ball.x - g.slime2.x) < 100 {
			g.randomJump(100)
		} else if math.Pow(g.ball.x  - g.slime2.x, 2) * 2 + math.Pow(g.ball.y - g.slime2.y, 2) < 28900 &&

			g.ball.x != g.slime2.x {
			g.randomJump(40)
		} else if math.Pow((-1 * g.ball.velocityX), 2) + math.Pow(g.ball.velocityY, 2) < 20 &&
			g.ball.x - g.slime2.x < 30 &&
			g.ball.x != g.slime2.x {
		    g.randomJump(40)
	    } else if (math.Abs(g.ball.x - g.slime2.x) < 150 && g.ball.y < 600 && g.ball.y > 400) {
	    	g.randomJump(100)
		}
	}

	// handles moving the slime left and right to intercept the ball
	if math.Abs(g.slime2.x -xWhenBallBelow475) <  closeEnough {
		g.slime2.velocityX = 0
	} else if xWhenBallBelow475+ closeEnough >= g.slime2.x {
		// move to net
		g.slime2.velocityX = 2
	} else if xWhenBallBelow475- closeEnough <= g.slime2.x {
		g.slime2.velocityX = -2
	} else{
		panic("Logic error in ai")
	}

}

// countFramesUntilBelow returns the number of frames it will take
// for y to be below the limit
// keep in mind that gravity is positive since the origin is the top left of the screen
func countFramesUntilBelow(y float64, vy float64, limit float64) int {
	count := 0
	for {
		vy += gravity
		y += vy
		if y > limit {
			return count
		}
		count++
  }
}