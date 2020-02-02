package main

import (
	"fmt"
	"math"
)

func (g *game) updateAI() {

	/*
	p1ToP2Wall := g.width - int32(g.slime1.x)
	p2ToP2Wall := g.width - int32(g.slime2.x)

	ballToP2Wall := g.width - int32(g.slime2.x)
	ballToP1Wall := g.width - int32(g.slime1.x)

	// currently AI must always be slime2

	meToEnemyWall := int32(g.slime2.x)
	enemyToTheirWall := int32(g.slime1.x)

    ballToEnemyWall := g.ball.x
	ballVToEnmyWall := g.ball.velocityX
    */


    // if the ball is far away from enemy's wall, or it's velocity is 0

    // when velocity is 0 on the right of the net, the AI is serving

    // aiIsServing := g.ball.velocityX == 0 && g.ball.x > g.net.x


	// g.slime2.velocityX = 3

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
	}
}


func (g *game) calculateXWhenBallBelow(yLimit float64) float64 {
  frameCount := countFramesUntilBelow(g.ball.y, g.ball.velocityY, yLimit)
  movingBallX := g.ball.x
  for i := 0; i < frameCount; i++ {
  	movingBallX += g.ball.velocityX
  }
  return math.Min(math.Max(movingBallX, g.width), 0)
}

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