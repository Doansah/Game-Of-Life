// Pong Game for SSD1306 OLED (128x64) using TinyGO
// Complete implementation with proper collision detection
package main

// Note: When porting to ESP32, add these imports:
// import (
// 	"image/color"
// 	"machine"
// 	"time"
// 	"tinygo.org/x/drivers/ssd1306"
// )

const (
	SCREEN_WIDTH  = 128
	SCREEN_HEIGHT = 64

	BALL_SIZE     = 3  // 3x3 pixel ball
	PADDLE_WIDTH  = 3  // Thin paddle
	PADDLE_HEIGHT = 12 // Paddle height

	BALL_SPEED   = 1 // Pixels per frame
	PADDLE_SPEED = 2 // Paddle movement speed

	WINNING_SCORE = 5
)

// ═══════════════════════════════════════════════════════════════
// GAME STRUCTURES
// ═══════════════════════════════════════════════════════════════

type Ball struct {
	x  int16 // Current X position (center)
	y  int16 // Current Y position (center)
	dx int16 // X velocity (-1 or 1)
	dy int16 // Y velocity (-1 or 1)
}

type Paddle struct {
	x      int16 // X position (left edge)
	y      int16 // Y position (top edge)
	width  int16
	height int16
}

type GameState struct {
	ball        Ball
	player1     Paddle // Left paddle
	player2     Paddle // Right paddle
	score1      int
	score2      int
	gameRunning bool
	winner      int // 0 = none, 1 = player1, 2 = player2
}

// ═══════════════════════════════════════════════════════════════
// COLLISION DETECTION - THE KEY LOGIC!
// ═══════════════════════════════════════════════════════════════

// detectPaddleCollision checks if ball collides with a paddle
// Returns: (collision bool, hitPosition float)
// hitPosition: 0.0 (top) to 1.0 (bottom) - used for angle calculation
func detectPaddleCollision(ball *Ball, paddle *Paddle) (bool, float32) {
	// Calculate ball bounds
	ballLeft := ball.x - BALL_SIZE/2
	ballRight := ball.x + BALL_SIZE/2
	ballTop := ball.y - BALL_SIZE/2
	ballBottom := ball.y + BALL_SIZE/2

	// Calculate paddle bounds
	paddleLeft := paddle.x
	paddleRight := paddle.x + paddle.width
	paddleTop := paddle.y
	paddleBottom := paddle.y + paddle.height

	// Rectangle overlap test (AABB - Axis-Aligned Bounding Box)
	collision := ballRight >= paddleLeft &&
		ballLeft <= paddleRight &&
		ballBottom >= paddleTop &&
		ballTop <= paddleBottom

	if !collision {
		return false, 0.0
	}

	// Calculate where on the paddle the ball hit (0.0 to 1.0)
	hitY := float32(ball.y - paddleTop)
	hitPosition := hitY / float32(paddle.height)

	return true, hitPosition
}

// detectWallCollision checks if ball hits top/bottom walls
func detectWallCollision(ball *Ball) bool {
	ballTop := ball.y - BALL_SIZE/2
	ballBottom := ball.y + BALL_SIZE/2

	return ballTop <= 0 || ballBottom >= SCREEN_HEIGHT
}

// detectGoal checks if ball went past paddles (someone scored)
// Returns: (scored bool, scoringPlayer int)
func detectGoal(ball *Ball) (bool, int) {
	if ball.x < 0 {
		return true, 2 // Player 2 scored
	}
	if ball.x > SCREEN_WIDTH {
		return true, 1 // Player 1 scored
	}
	return false, 0
}

// ═══════════════════════════════════════════════════════════════
// GAME LOGIC
// ═══════════════════════════════════════════════════════════════

func newGame() *GameState {
	return &GameState{
		ball: Ball{
			x:  SCREEN_WIDTH / 2,
			y:  SCREEN_HEIGHT / 2,
			dx: 1, // Start moving right
			dy: 1, // Start moving down
		},
		player1: Paddle{
			x:      5, // Left side
			y:      SCREEN_HEIGHT/2 - PADDLE_HEIGHT/2,
			width:  PADDLE_WIDTH,
			height: PADDLE_HEIGHT,
		},
		player2: Paddle{
			x:      SCREEN_WIDTH - 5 - PADDLE_WIDTH, // Right side
			y:      SCREEN_HEIGHT/2 - PADDLE_HEIGHT/2,
			width:  PADDLE_WIDTH,
			height: PADDLE_HEIGHT,
		},
		score1:      0,
		score2:      0,
		gameRunning: true,
		winner:      0,
	}
}

func (g *GameState) resetBall() {
	g.ball.x = SCREEN_WIDTH / 2
	g.ball.y = SCREEN_HEIGHT / 2
	// Alternate direction
	if g.ball.dx > 0 {
		g.ball.dx = -1
	} else {
		g.ball.dx = 1
	}
	g.ball.dy = 1
}

func (g *GameState) updateBall() {
	// Move ball
	g.ball.x += g.ball.dx * BALL_SPEED
	g.ball.y += g.ball.dy * BALL_SPEED

	// Check wall collision (top/bottom)
	if detectWallCollision(&g.ball) {
		g.ball.dy = -g.ball.dy // Bounce vertically
		// Clamp position
		if g.ball.y < BALL_SIZE/2 {
			g.ball.y = BALL_SIZE / 2
		}
		if g.ball.y > SCREEN_HEIGHT-BALL_SIZE/2 {
			g.ball.y = SCREEN_HEIGHT - BALL_SIZE/2
		}
	}

	// Check paddle collision
	// Player 1 (left paddle)
	if collision, hitPos := detectPaddleCollision(&g.ball, &g.player1); collision && g.ball.dx < 0 {
		g.ball.dx = -g.ball.dx // Bounce horizontally

		// Add angle based on where ball hit paddle
		// Hit near top: ball goes up more
		// Hit near bottom: ball goes down more
		if hitPos < 0.33 {
			g.ball.dy = -1 // Hit top third
		} else if hitPos > 0.66 {
			g.ball.dy = 1 // Hit bottom third
		}
		// Middle third keeps current dy

		// Move ball out of paddle to prevent stuck
		g.ball.x = g.player1.x + g.player1.width + BALL_SIZE/2
	}

	// Player 2 (right paddle)
	if collision, hitPos := detectPaddleCollision(&g.ball, &g.player2); collision && g.ball.dx > 0 {
		g.ball.dx = -g.ball.dx // Bounce horizontally

		// Add angle based on hit position
		if hitPos < 0.33 {
			g.ball.dy = -1
		} else if hitPos > 0.66 {
			g.ball.dy = 1
		}

		// Move ball out of paddle
		g.ball.x = g.player2.x - BALL_SIZE/2
	}

	// Check goals
	if scored, player := detectGoal(&g.ball); scored {
		if player == 1 {
			g.score1++
		} else {
			g.score2++
		}

		// Check for winner
		if g.score1 >= WINNING_SCORE {
			g.winner = 1
			g.gameRunning = false
		} else if g.score2 >= WINNING_SCORE {
			g.winner = 2
			g.gameRunning = false
		} else {
			g.resetBall()
		}
	}
}

func (g *GameState) movePaddle(paddle *Paddle, direction int16) {
	paddle.y += direction * PADDLE_SPEED

	// Clamp to screen bounds
	if paddle.y < 0 {
		paddle.y = 0
	}
	if paddle.y > SCREEN_HEIGHT-paddle.height {
		paddle.y = SCREEN_HEIGHT - paddle.height
	}
}

// ═══════════════════════════════════════════════════════════════
// RENDERING (Placeholder - replace with actual display code)
// ═══════════════════════════════════════════════════════════════

// When porting to ESP32 with SSD1306, uncomment and use this:
/*
func (g *GameState) draw(display *ssd1306.Device) {
	display.ClearBuffer()

	// Draw center line (dashed)
	for y := int16(0); y < SCREEN_HEIGHT; y += 4 {
		display.SetPixel(SCREEN_WIDTH/2, y, color.RGBA{255, 255, 255, 255})
		display.SetPixel(SCREEN_WIDTH/2, y+1, color.RGBA{255, 255, 255, 255})
	}

	// Draw ball
	for dy := int16(-BALL_SIZE / 2); dy <= BALL_SIZE/2; dy++ {
		for dx := int16(-BALL_SIZE / 2); dx <= BALL_SIZE/2; dx++ {
			display.SetPixel(g.ball.x+dx, g.ball.y+dy, color.RGBA{255, 255, 255, 255})
		}
	}

	// Draw player 1 paddle
	for dy := int16(0); dy < g.player1.height; dy++ {
		for dx := int16(0); dx < g.player1.width; dx++ {
			display.SetPixel(g.player1.x+dx, g.player1.y+dy, color.RGBA{255, 255, 255, 255})
		}
	}

	// Draw player 2 paddle
	for dy := int16(0); dy < g.player2.height; dy++ {
		for dx := int16(0); dx < g.player2.width; dx++ {
			display.SetPixel(g.player2.x+dx, g.player2.y+dy, color.RGBA{255, 255, 255, 255})
		}
	}

	// Draw scores (simple dots)
	// Player 1 score (left side)
	for i := 0; i < g.score1; i++ {
		display.SetPixel(10+int16(i*5), 2, color.RGBA{255, 255, 255, 255})
		display.SetPixel(10+int16(i*5), 3, color.RGBA{255, 255, 255, 255})
	}

	// Player 2 score (right side)
	for i := 0; i < g.score2; i++ {
		display.SetPixel(SCREEN_WIDTH-15-int16(i*5), 2, color.RGBA{255, 255, 255, 255})
		display.SetPixel(SCREEN_WIDTH-15-int16(i*5), 3, color.RGBA{255, 255, 255, 255})
	}

	display.Display()
}
*/

// ═══════════════════════════════════════════════════════════════
// DEMO / EXPLANATION - Run this to understand the collision logic
// ═══════════════════════════════════════════════════════════════

func main() {
	// This is a structural example
	// To run on ESP32, you'd need to:
	// 1. Initialize I2C and display (like in tinygo_ssd1306_version.go)
	// 2. Add button inputs for paddle control
	// 3. Flash to ESP32

	println("Pong Game Structure - Collision Detection Demo")
	println("")
	println("═══════════════════════════════════════════════")
	println("COLLISION DETECTION EXPLAINED:")
	println("═══════════════════════════════════════════════")
	println("")
	println("1. BALL-PADDLE COLLISION (AABB Method):")
	println("   - Check if rectangles overlap on both axes")
	println("   - Ball bounds: (x±size/2, y±size/2)")
	println("   - Paddle bounds: (x, y) to (x+width, y+height)")
	println("   - Overlap = collision detected!")
	println("")
	println("2. COLLISION RESPONSE:")
	println("   - Reverse dx (horizontal velocity)")
	println("   - Adjust dy based on hit position (top/mid/bottom)")
	println("   - Push ball out of paddle to prevent 'stuck'")
	println("")
	println("3. WALL COLLISION:")
	println("   - Check if ball.y ≤ 0 or ≥ SCREEN_HEIGHT")
	println("   - Reverse dy (vertical velocity)")
	println("")
	println("4. GOAL DETECTION:")
	println("   - Ball.x < 0 = Player 2 scores")
	println("   - Ball.x > SCREEN_WIDTH = Player 1 scores")
	println("")
	println("═══════════════════════════════════════════════")
	println("")

	// Example collision test
	game := newGame()

	// Simulate ball moving toward player 1
	game.ball.x = 10
	game.ball.y = 32
	game.ball.dx = -1

	println("TEST CASE:")
	println("Ball position: (10, 32)")
	println("Player 1 paddle: x=5, y=26, width=3, height=12")

	collision, hitPos := detectPaddleCollision(&game.ball, &game.player1)

	if collision {
		println("✅ COLLISION DETECTED!")
		println("Hit position:", hitPos, "(0.0=top, 1.0=bottom)")
	} else {
		println("❌ No collision")
	}

	// AI opponent example
	println("")
	println("BONUS: Simple AI for Player 2:")
	println("  if ball.y > paddle.y + height/2:")
	println("    paddle.y += PADDLE_SPEED  // Move down")
	println("  else if ball.y < paddle.y + height/2:")
	println("    paddle.y -= PADDLE_SPEED  // Move up")
}
