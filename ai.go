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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (g *game) randomJump40() {
	// 40% of the time will do a jump
	if rand.Float64() < .4 {
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





func (g *game) updateAI() {

	/*
	p1ToP2Wall := g.width - int32(g.slime1.x)
	p2ToP2Wall := g.width - int32(g.slime2.x)

	ballToP2Wall := g.width - int32(g.slime2.x)
	ballToP1Wall := g.width - int32(g.slime1.x)

	// currently AI must always be slime2


	enemyToTheirWall := int32(g.slime1.x)

    ballToEnemyWall := g.ball.x
	ballVToEnmyWall := g.ball.velocityX
    */
	//meToEnemyWall := int32(g.slime2.x)

	// currently AI is hardcode to player 2
	meToEnemyWall := g.slime2.x


    // if the ball is far away from enemy's wall, or it's velocity is 0

    // when velocity is 0 on the right of the net, the AI is serving

    // aiIsServing := g.ball.velocityX == 0 && g.ball.x > g.net.x


	// g.slime2.velocityX = 3
	randomSaltInt := 15 + rand.Float64() * 10
	randSalt := randomSaltInt
	xWhenBallBelow125 := g.calculateXWhenBallBelow(125)
/*
	if g.ai.state == enemyCloseToNet {
		// TODO: generalize these checks to not rely on magic numbers

		// if ball is low and moving down
		if g.ball.y < 300 && g.ball.velocityY < -3 {
			fmt.Printf("ball is low and moving left %v\n", g.ball.velocityY)
			// move to wall
			g.slime2.velocityX = 2
			// jump
		}
		//something := 23 +
		//s1 := rand.NewSource(time.Now().UnixNano())
		//r1 := rand.New(s1)
	}*/


	if xWhenBallBelow125 < float64(g.width) / 2 { // if balls' trajectory is towards opponents side, just get centered
	    // handle return logic here
		return
	}

	//
	if math.Abs(meToEnemyWall - xWhenBallBelow125) <  randSalt { // if AI is close to where ball will land
		// 30% of the time just stay there
		if rand.Float64() < .3 {
			return
		}

		if (meToEnemyWall >= 900 && g.ball.x > 830 ||
			meToEnemyWall <= 580 && g.ball.x < 530 && math.Abs(g.ball.x - meToEnemyWall) < 100) {
			g.randomJump40()
		} else if (math.Pow(g.ball.velocityX, 2) + math.Pow(g.ball.velocityY, 2) < 20 &&
			g.ball.x - meToEnemyWall < 30 &&
			g.ball.x != meToEnemyWall) {
			    g.randomJump40()
	    }
	//else if math.Pow(g.ball.x - meToEnemyWall, 2) * 2 + math.Pow() < 28900 {
	}

	/*

    // if close to where the ball will end up
    if (Math.abs(this.meToEnemyWall - xWhenBallBelow125) < something)
    {
      if (me.y != 0 || Math.random() < 0.3)
        return;

        if (
            (this.meToEnemyWall >= 900 && this.ballToEnemyWall > 830) ||
            (this.meToEnemyWall <= 580 && this.ballToEnemyWall < 530 && Math.abs(this.ballToEnemyWall - this.meToEnemyWall) < 100)
        ) {
            this.randomJump40Percent();
      } else if ((Math.pow(this.ballToEnemyWall - this.meToEnemyWall, 2) * 2 + Math.pow(ball.y - me.y, 2) < 28900) &&
            (this.ballToEnemyWall != this.meToEnemyWall)) {
            this.randomJump40Percent();
        } else if ((Math.pow(this.ballVXToEnemyWall, 2) + Math.pow(ball.velocityY, 2) < 20) &&
            (this.ballToEnemyWall - this.meToEnemyWall < 30) &&
            (this.ballToEnemyWall != this.meToEnemyWall)) {
            this.randomJump40Percent();
        } else if ((Math.abs(this.ballToEnemyWall - this.meToEnemyWall) < 150) &&
            (ball.y > 50) && (ball.y < 400) && (Math.random() < 0.5)) {
            this.randomJump40Percent();
        }
    }

	 */


	// handles moving the slime left and right to intercept the ball
	if g.ai.state == initialAiState {
		fmt.Printf("in tick %v %v %v\n", meToEnemyWall, xWhenBallBelow125, randSalt)
		if math.Abs(meToEnemyWall - xWhenBallBelow125) <  randSalt {
			g.slime2.velocityX = 0
		} else if xWhenBallBelow125 + randSalt >= meToEnemyWall {
			// move to net
			g.slime2.velocityX = 2
		} else if xWhenBallBelow125 - randSalt <= meToEnemyWall {
			fmt.Printf("got here when %v and %v", xWhenBallBelow125, meToEnemyWall)
			g.slime2.velocityX = -2
		} else{
			fmt.Printf("should never get here")
		}
	}

}

// countFramesUntilBelow returns the number of frames it will take
// for y to be less than limit given a velocity of vy
func countFramesUntilBelow(y float64, vy float64, limit float64) int {
	count := 0
	for {
		vy -= 1
		y += vy
		if y < limit {
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