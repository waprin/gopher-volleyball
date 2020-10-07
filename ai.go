package main

import (
	"math"
	"math/rand"
	"fmt"
)

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func (g *game) randomJump(percent float64) {
	// 40% of the time will do a jump
	if g.slime2.y != float64(g.height) {
		return
	}
	if rand.Float64() < percent {
		g.slime2.velocityY = -5
	}
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

func (g *game) updateAI() {
    // if the ball is far away from enemy's wall, or it's velocity is 0

    // when velocity is 0 on the right of the net, the AI is serving

    // aiIsServing := g.ball.velocityX == 0 && g.ball.x > g.net.x


	// g.slime2.velocityX = 3
	randomSaltInt := 15 + rand.Float64() * 10
	randSalt := randomSaltInt
	xWhenBallBelow125 := g.calculateXWhenBallBelow(125)

	if g.ball.x < float64(g.width/ 2) {
		g.ai.state = initialAiState
	}

	if g.ai.state == initialAiState {
		// consider opponents position when deciding next move
		if g.slime1.x > 250 {
			g.ai.state = enemyCloseToNet
		} else if g.slime1.x < 200 {
			g.ai.state = enemyFarFromNet
		} else {
			g.ai.state = enemyInMiddle
		}

		// randomly do a different move then the usual one
		if rand.Float64() < .35 {
			g.ai.state = aiState(rand.Intn(4))
		}
	}

	if g.aiIsServing() {
		if g.ai.state == enemyCloseToNet {
			fmt.Printf("serving enemy close to net %v %v\n", g.ball.y, g.ball.velocityY)
			// if ball moving down
			if g.ball.y > 200 && g.ball.velocityY > 3 {
				g.slime2.velocityX = 2
				g.slime2.velocityY = -5
			}
		} else if g.ai.state == enemyFarFromNet {
			fmt.Printf("serving enemy far from net\n")
			if g.ball.y > 150 && g.ball.velocityY > 4 {
				g.slime2.velocityX = -2
				g.slime2.velocityY = -5
			}
		} else if g.ai.state == enemyInMiddle {
			fmt.Printf("serving enemy in middle %v %v\n", g.ball.velocityY, g.slime2.x)
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
				g.slime2.velocityX = -2;
			}
		}

		return
	}

	if xWhenBallBelow125 < float64(g.width) / 2 { // if balls' trajectory is towards opponents side, just get centered
	    // handle return logic here

		differenceFromCenter := int32(g.slime2.x) - 600
		if Abs(differenceFromCenter) < 20 {
			fmt.Printf("got centered\n")
			g.slime2.velocityX = 0
		} else if differenceFromCenter >= 20 {
			//console.log('moveToNet 2');
			fmt.Printf("centering to net\n")
			g.slime2.velocityX = -2
		} else {
			fmt.Printf("centering to wall\n")
			g.slime2.velocityX = 2
		}
		return
	}

	//
	if math.Abs(g.slime2.x - xWhenBallBelow125) <  randSalt { // if AI is close to where ball will land
		// 30% of the time just stay there
		if rand.Float64() < .3 {
			return
		}

		if g.slime2.x >= 900 && g.ball.x > 830 ||
			g.slime2.x <= 580 && g.ball.x < 530 && math.Abs(g.ball.x - g.slime2.x) < 100 {
			g.randomJump(40)
		} else if math.Pow(g.ball.x  - g.slime2.x, 2) * 2 + math.Pow(g.ball.y - g.slime2.y, 2) < 28900 &&
				g.ball.x != g.slime2.x {
			g.randomJump(40)
		} else if math.Pow((-1 * g.ball.velocityX), 2) + math.Pow(g.ball.velocityY, 2) < 20 &&
			g.ball.x - g.slime2.x < 30 &&
			g.ball.x != g.slime2.x {
			    g.randomJump(40)
	    } else if (math.Abs(g.ball.x - g.slime2.x) < 150 && g.ball.y > 50 && g.ball.y < 400  && rand.Float64() < .5) {
	    	g.randomJump(40)
		}
	//else if math.Pow(g.ball.x - meToEnemyWall, 2) * 2 + math.Pow() < 28900 {
	}

	// handles moving the slime left and right to intercept the ball
	if math.Abs(g.slime2.x - xWhenBallBelow125) <  randSalt {
		g.slime2.velocityX = 0
	} else if xWhenBallBelow125 + randSalt >= g.slime2.x {
		// move to net
		g.slime2.velocityX = 2
	} else if xWhenBallBelow125 - randSalt <= g.slime2.x {
		g.slime2.velocityX = -2
	} else{
		fmt.Printf("should never get here")
	}

}

// countFramesUntilBelow returns the number of frames it will take
// for y to be less than limit given a velocity of vy
func countFramesUntilBelow(y float64, vy float64, limit float64) int {
	count := 0
	for {
		vy += 1
		y += vy
		if y > limit {
			return count
		}
		count++
  }
}

/*
calculateXWhenBallBelow : function(yLimit) {
var frameCount = countFramesTillBelow(ball.y, ball.velocityY, yLimit);
var toEnemyWall         = this.ballToEnemyWall;
var velocityToEnemyWall = this.ballVXToEnemyWall;
for(var i = 0; i < frameCount; i++) {
toEnemyWall += velocityToEnemyWall;
if(toEnemyWall < 0) {
toEnemyWall = 0;
velocityToEnemyWall = -velocityToEnemyWall;
} else if(toEnemyWall > gameWidth) {
toEnemyWall = gameWidth;
velocityToEnemyWall = -velocityToEnemyWall;
}
}
return toEnemyWall;
}*/