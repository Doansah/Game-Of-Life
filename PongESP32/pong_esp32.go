// Pong Game for ESP32 with SSD1306 OLED (128x64)
// Two-player game with button controls
package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ssd1306"
)

const (
	SCREEN_WIDTH  = 128
	SCREEN_HEIGHT = 64

	BALL_SIZE     = 3  // 3x3 pixel ball
	PADDLE_WIDTH  = 3  // Thin paddle
	PADDLE_HEIGHT = 12 // Paddle height

	BALL_SPEED = 3 // Pixels per frame (was 1, now 3x faster!)
	JUMP_SPEED = 8 // How fast paddle rises when button pressed
	FALL_SPEED = 2 // How fast paddle falls when button released (gravity)

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
	player2     Paddle // Right paddle (AI or player 2)
	score1      int
	score2      int
	gameRunning bool
	winner      int  // 0 = none, 1 = player1, 2 = player2
	aiEnabled   bool // AI mode vs 2-player
}

// ═══════════════════════════════════════════════════════════════
// COLLISION DETECTION
// ═══════════════════════════════════════════════════════════════

func detectPaddleCollision(ball *Ball, paddle *Paddle) (bool, float32) {
	ballLeft := ball.x - BALL_SIZE/2
	ballRight := ball.x + BALL_SIZE/2
	ballTop := ball.y - BALL_SIZE/2
	ballBottom := ball.y + BALL_SIZE/2

	paddleLeft := paddle.x
	paddleRight := paddle.x + paddle.width
	paddleTop := paddle.y
	paddleBottom := paddle.y + paddle.height

	collision := ballRight >= paddleLeft &&
		ballLeft <= paddleRight &&
		ballBottom >= paddleTop &&
		ballTop <= paddleBottom

	if !collision {
		return false, 0.0
	}

	hitY := float32(ball.y - paddleTop)
	hitPosition := hitY / float32(paddle.height)

	return true, hitPosition
}

func detectWallCollision(ball *Ball) bool {
	ballTop := ball.y - BALL_SIZE/2
	ballBottom := ball.y + BALL_SIZE/2
	return ballTop <= 0 || ballBottom >= SCREEN_HEIGHT
}

func detectGoal(ball *Ball) (bool, int) {
	if ball.x < 0 {
		return true, 2
	}
	if ball.x > SCREEN_WIDTH {
		return true, 1
	}
	return false, 0
}

// ═══════════════════════════════════════════════════════════════
// GAME LOGIC
// ═══════════════════════════════════════════════════════════════

func newGame(aiEnabled bool) *GameState {
	return &GameState{
		ball: Ball{
			x:  SCREEN_WIDTH / 2,
			y:  SCREEN_HEIGHT / 2,
			dx: 1,
			dy: 1,
		},
		player1: Paddle{
			x:      5,
			y:      SCREEN_HEIGHT/2 - PADDLE_HEIGHT/2,
			width:  PADDLE_WIDTH,
			height: PADDLE_HEIGHT,
		},
		player2: Paddle{
			x:      SCREEN_WIDTH - 5 - PADDLE_WIDTH,
			y:      SCREEN_HEIGHT/2 - PADDLE_HEIGHT/2,
			width:  PADDLE_WIDTH,
			height: PADDLE_HEIGHT,
		},
		score1:      0,
		score2:      0,
		gameRunning: true,
		winner:      0,
		aiEnabled:   aiEnabled,
	}
}

func (g *GameState) resetBall() {
	g.ball.x = SCREEN_WIDTH / 2
	g.ball.y = SCREEN_HEIGHT / 2
	if g.ball.dx > 0 {
		g.ball.dx = -1
	} else {
		g.ball.dx = 1
	}
	g.ball.dy = 1
}

func (g *GameState) updateBall() {
	g.ball.x += g.ball.dx * BALL_SPEED
	g.ball.y += g.ball.dy * BALL_SPEED

	if detectWallCollision(&g.ball) {
		g.ball.dy = -g.ball.dy
		if g.ball.y < BALL_SIZE/2 {
			g.ball.y = BALL_SIZE / 2
		}
		if g.ball.y > SCREEN_HEIGHT-BALL_SIZE/2 {
			g.ball.y = SCREEN_HEIGHT - BALL_SIZE/2
		}
	}

	// Player 1 paddle collision
	if collision, hitPos := detectPaddleCollision(&g.ball, &g.player1); collision && g.ball.dx < 0 {
		g.ball.dx = -g.ball.dx
		if hitPos < 0.33 {
			g.ball.dy = -1
		} else if hitPos > 0.66 {
			g.ball.dy = 1
		}
		g.ball.x = g.player1.x + g.player1.width + BALL_SIZE/2
	}

	// Player 2 paddle collision
	if collision, hitPos := detectPaddleCollision(&g.ball, &g.player2); collision && g.ball.dx > 0 {
		g.ball.dx = -g.ball.dx
		if hitPos < 0.33 {
			g.ball.dy = -1
		} else if hitPos > 0.66 {
			g.ball.dy = 1
		}
		g.ball.x = g.player2.x - BALL_SIZE/2
	}

	// Check goals
	if scored, player := detectGoal(&g.ball); scored {
		if player == 1 {
			g.score1++
		} else {
			g.score2++
		}

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

func (g *GameState) movePaddle(paddle *Paddle, direction int16, speed int16) {
	paddle.y += direction * speed
	if paddle.y < 0 {
		paddle.y = 0
	}
	if paddle.y > SCREEN_HEIGHT-paddle.height {
		paddle.y = SCREEN_HEIGHT - paddle.height
	}
}

// Simple AI that tracks ball
func (g *GameState) updateAI() {
	if !g.aiEnabled {
		return
	}

	paddleCenter := g.player2.y + g.player2.height/2

	// Add some delay/imperfection to make AI beatable
	if g.ball.y > paddleCenter+2 {
		g.movePaddle(&g.player2, 1, 2) // AI uses moderate speed
	} else if g.ball.y < paddleCenter-2 {
		g.movePaddle(&g.player2, -1, 2) // AI uses moderate speed
	}
}

// ═══════════════════════════════════════════════════════════════
// RENDERING
// ═══════════════════════════════════════════════════════════════

func (g *GameState) draw(display *ssd1306.Device) {
	display.ClearBuffer()

	// Draw center line (dashed)
	for y := int16(0); y < SCREEN_HEIGHT; y += 4 {
		display.SetPixel(SCREEN_WIDTH/2, y, color.RGBA{255, 255, 255, 255})
	}

	// Draw ball
	for dy := int16(-BALL_SIZE / 2); dy <= BALL_SIZE/2; dy++ {
		for dx := int16(-BALL_SIZE / 2); dx <= BALL_SIZE/2; dx++ {
			x := g.ball.x + dx
			y := g.ball.y + dy
			if x >= 0 && x < SCREEN_WIDTH && y >= 0 && y < SCREEN_HEIGHT {
				display.SetPixel(x, y, color.RGBA{255, 255, 255, 255})
			}
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

	// Draw scores (dots)
	for i := 0; i < g.score1 && i < 5; i++ {
		x := int16(10 + i*6)
		display.SetPixel(x, 2, color.RGBA{255, 255, 255, 255})
		display.SetPixel(x+1, 2, color.RGBA{255, 255, 255, 255})
		display.SetPixel(x, 3, color.RGBA{255, 255, 255, 255})
		display.SetPixel(x+1, 3, color.RGBA{255, 255, 255, 255})
	}

	for i := 0; i < g.score2 && i < 5; i++ {
		x := SCREEN_WIDTH - int16(16+i*6)
		display.SetPixel(x, 2, color.RGBA{255, 255, 255, 255})
		display.SetPixel(x+1, 2, color.RGBA{255, 255, 255, 255})
		display.SetPixel(x, 3, color.RGBA{255, 255, 255, 255})
		display.SetPixel(x+1, 3, color.RGBA{255, 255, 255, 255})
	}

	display.Display()
}

// Draw simple text for winner screen
func drawText(display *ssd1306.Device, text string, x, y int16) {
	// Very simple 3x5 font patterns
	fontMap := map[rune][][]bool{
		'P': {
			{true, true, true},
			{true, false, true},
			{true, true, true},
			{true, false, false},
			{true, false, false},
		},
		'1': {
			{false, true, false},
			{true, true, false},
			{false, true, false},
			{false, true, false},
			{true, true, true},
		},
		'2': {
			{true, true, true},
			{false, false, true},
			{true, true, true},
			{true, false, false},
			{true, true, true},
		},
		'W': {
			{true, false, true},
			{true, false, true},
			{true, false, true},
			{true, true, true},
			{true, false, true},
		},
		'I': {
			{true, true, true},
			{false, true, false},
			{false, true, false},
			{false, true, false},
			{true, true, true},
		},
		'N': {
			{true, false, true},
			{true, true, true},
			{true, true, true},
			{true, false, true},
			{true, false, true},
		},
		'S': {
			{true, true, true},
			{true, false, false},
			{true, true, true},
			{false, false, true},
			{true, true, true},
		},
	}

	offsetX := int16(0)
	for _, char := range text {
		pattern, exists := fontMap[char]
		if !exists {
			offsetX += 4
			continue
		}

		for row := 0; row < 5; row++ {
			for col := 0; col < 3; col++ {
				if pattern[row][col] {
					display.SetPixel(x+offsetX+int16(col), y+int16(row), color.RGBA{255, 255, 255, 255})
				}
			}
		}
		offsetX += 4
	}
}

func showWinner(display *ssd1306.Device, winner int) {
	display.ClearBuffer()

	if winner == 1 {
		drawText(display, "P1 WINS", 40, 28)
	} else {
		drawText(display, "P2 WINS", 40, 28)
	}

	display.Display()
}

// ═══════════════════════════════════════════════════════════════
// MAIN - ESP32 VERSION
// ═══════════════════════════════════════════════════════════════

func main() {
	println("[PONG] Initializing...")

	// Configure I2C
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
		SDA:       machine.GPIO21,
		SCL:       machine.GPIO22,
	})

	// Configure button (GPIO18)
	button := machine.GPIO18
	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// Initialize display
	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})
	display.ClearDisplay()

	println("[PONG] Display initialized")

	// Show splash screen
	display.ClearBuffer()
	drawText(display, "PONG", 48, 20)
	display.Display()
	time.Sleep(2 * time.Second)

	// Start game with AI enabled (single player mode)
	game := newGame(true)
	println("[PONG] Game started - AI mode")

	frameCount := 0

	// Main game loop
	for {
		// Check button - controls player 1 paddle (Flappy Bird style)
		buttonPressed := !button.Get()

		// Update game state
		if game.gameRunning {
			// Flappy Bird mechanics: paddle always falls, button makes it rise
			if buttonPressed {
				// Button pressed - JUMP UP (fast!)
				game.movePaddle(&game.player1, -1, JUMP_SPEED)
			} else {
				// Button not pressed - FALL DOWN slowly (gravity)
				game.movePaddle(&game.player1, 1, FALL_SPEED)
			}
			game.updateBall()
			game.updateAI() // AI controls player 2

			// Draw game
			game.draw(display)

			frameCount++
			if frameCount%60 == 0 {
				println("[PONG] Score:", game.score1, "-", game.score2)
			}
		} else {
			// Game over
			showWinner(display, game.winner)
			time.Sleep(3 * time.Second)

			// Reset for new game
			game = newGame(true)
			println("[PONG] New game started")
		}

		// Frame delay (adjust for game speed)
		time.Sleep(50 * time.Millisecond) // 20 FPS
	}
}
