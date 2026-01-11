// Pong Game for ESP32 with SSD1306 OLED (128x64)
// Two-player game with button controls
package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
)

const (
	SCREEN_WIDTH  = 128
	SCREEN_HEIGHT = 64

	BALL_SIZE     = 3
	PADDLE_WIDTH  = 3
	PADDLE_HEIGHT = 12

	INITIAL_BALL_SPEED = 5.0 // Starting speed (pixels per frame)
	SPEED_INCREMENT    = 0.4 // Speed increase per paddle hit
	JUMP_SPEED         = 8   // How fast paddle rises when button pressed
	FALL_SPEED         = 2   // How fast paddle falls when button released (gravity)

	WINNING_SCORE = 5
)

// ═══════════════════════════════════════════════════════════════
// GAME STRUCTURES
// ═══════════════════════════════════════════════════════════════

type Ball struct {
	x     int16   // Current X position (center)
	y     int16   // Current Y position (center)
	dx    int16   // X velocity (-1 or 1)
	dy    int16   // Y velocity (-1 or 1)
	speed float32 // Current speed (increases with each hit)
}

type Paddle struct {
	x      int16 // X position (left edge)
	y      int16 // Y position (top edge)
	width  int16
	height int16
}

type WindParticle struct {
	x    int16 // X position
	y    int16 // Y position
	life int   // Lifetime counter (0 = dead)
}

type GameState struct {
	ball           Ball
	player1        Paddle // Left paddle
	player2        Paddle // Right paddle (AI or player 2)
	score1         int
	score2         int
	gameRunning    bool
	winner         int            // 0 = none, 1 = player1, 2 = player2
	aiEnabled      bool           // AI mode vs 2-player
	windParticles  []WindParticle // Visual effect when jumping
	collisionCount int
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
			x:     SCREEN_WIDTH / 2,
			y:     SCREEN_HEIGHT / 2,
			dx:    1,
			dy:    1,
			speed: INITIAL_BALL_SPEED,
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
		score1:        0,
		score2:        0,
		gameRunning:   true,
		winner:        0,
		aiEnabled:     aiEnabled,
		windParticles: make([]WindParticle, 0, 10), // Pre-allocate for efficiency
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
	g.ball.speed = INITIAL_BALL_SPEED // Reset speed when ball resets (after scoring)
}

func (g *GameState) updateBall() {
	// Move ball using its current speed
	g.ball.x += int16(float32(g.ball.dx) * g.ball.speed)
	g.ball.y += int16(float32(g.ball.dy) * g.ball.speed)

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
		g.ball.speed += SPEED_INCREMENT // Increase speed on hit!

		if g.ball.speed > 6.5 { // cap increment speed at 6
			g.ball.speed = 6.5
		}

		g.collisionCount += 1

		println("[PONG] P1 hit! Ball speed:", g.ball.speed)
		println("collision count: ", g.collisionCount)
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
		g.ball.speed += SPEED_INCREMENT // Increase speed on hit!
		println("[PONG] P2 hit! Ball speed:", g.ball.speed)
		g.collisionCount += 1
		println("collision count: ", g.collisionCount)

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
// WIND PARTICLE EFFECTS (visual feedback when jumping)
// ═══════════════════════════════════════════════════════════════

// Spawn wind particles at the bottom of the paddle when jumping
func (g *GameState) spawnWindParticles() {
	// Create 4 particles at the bottom of player1 paddle with slight spread
	paddleBottom := g.player1.y + g.player1.height
	paddleCenterX := g.player1.x + g.player1.width/2

	// Spawn 4 particles with horizontal spread
	offsets := []int16{-1, 0, 1, 0} // Left, center, right, center again
	for i, offset := range offsets {
		particle := WindParticle{
			x:    paddleCenterX + offset,    // Spread around paddle center
			y:    paddleBottom + int16(i%2), // Slightly staggered vertically
			life: 3,                         // Lasts for 15 frames (longer trail)
		}
		g.windParticles = append(g.windParticles, particle)
	}
}

// Update wind particles (move down and fade)
func (g *GameState) updateWindParticles() {
	// Update each particle
	for i := len(g.windParticles) - 1; i >= 0; i-- {
		g.windParticles[i].life--
		g.windParticles[i].y++ // Fall downward

		// Remove dead particles
		if g.windParticles[i].life <= 0 || g.windParticles[i].y >= SCREEN_HEIGHT {
			// Remove particle by swapping with last and truncating
			g.windParticles[i] = g.windParticles[len(g.windParticles)-1]
			g.windParticles = g.windParticles[:len(g.windParticles)-1]
		}
	}
}

// ═══════════════════════════════════════════════════════════════
// RENDERING
// ═══════════════════════════════════════════════════════════════

func (g *GameState) draw(display *ssd1306.Device) {
	display.ClearBuffer()

	// Draw center line (dashed)
	for y := int16(0); y < SCREEN_HEIGHT; y += 5 {
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

	// Draw scores (using numbers instead of dots)
	white := color.RGBA{255, 255, 255, 255}

	// Draw P1 score (left side)
	scoreText := string(rune('0' + g.score1))
	tinyfont.WriteLine(display, &freesans.BoldOblique9pt7b, 15, 12, scoreText, white)

	// Draw P2 score (right side)
	scoreText = string(rune('0' + g.score2))
	tinyfont.WriteLine(display, &freesans.Bold9pt7b, SCREEN_WIDTH-22, 12, scoreText, white)

	// Draw wind particles (small dots that trail behind paddle when jumping)
	for _, particle := range g.windParticles {
		// Draw single pixel for each particle (simple and fast)
		if particle.x >= 0 && particle.x < SCREEN_WIDTH && particle.y >= 0 && particle.y < SCREEN_HEIGHT {
			display.SetPixel(particle.x, particle.y, color.RGBA{255, 255, 255, 255})
		}
	}

	display.Display()
}

func showWinner(display *ssd1306.Device, winner int) {
	display.ClearBuffer()

	white := color.RGBA{255, 255, 255, 255}

	if winner == 1 {
		tinyfont.WriteLine(display, &freesans.Bold12pt7b, 10, 35, "P1 WINS!", white)
	} else {
		tinyfont.WriteLine(display, &freesans.Bold12pt7b, 10, 35, "P2 WINS!", white)
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

	// Configure buttons
	buttonP1 := machine.GPIO18 // Player 1 button
	buttonP1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	buttonP2 := machine.GPIO19 // Player 2 button
	buttonP2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

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
	white := color.RGBA{255, 255, 255, 255}
	tinyfont.WriteLine(display, &freesans.Bold18pt7b, 25, 30, "PONG", white)
	tinyfont.WriteLine(display, &freesans.Regular9pt7b, 15, 50, "Press to Start", white)
	display.Display()
	time.Sleep(2 * time.Second)

	// Start game with AI enabled (single player mode)
	game := newGame(true)
	println("[PONG] Game started - AI mode")

	frameCount := 0
	lastP2ButtonState := false // Track if P2 is actively playing

	// Main game loop
	for {
		// Check buttons - Flappy Bird style controls for both players
		buttonP1Pressed := !buttonP1.Get()
		buttonP2Pressed := !buttonP2.Get()

		// Update P2 playing state (if button pressed recently, P2 is playing)
		if buttonP2Pressed {
			if !lastP2ButtonState {
				println("[PONG] Player 2 joined! (2-player mode)")
			}
			lastP2ButtonState = true
		}

		// Update game state
		if game.gameRunning {
			// PLAYER 1 - Flappy Bird mechanics: paddle always falls, button makes it rise
			if buttonP1Pressed {
				// Button pressed - JUMP UP (fast!)
				game.movePaddle(&game.player1, -1, JUMP_SPEED)
				game.spawnWindParticles() // Spawn particles on jump
			} else {
				// Button not pressed - FALL DOWN slowly (gravity)
				game.movePaddle(&game.player1, 1, FALL_SPEED)
			}

			// PLAYER 2 - Can be controlled by button OR AI
			if lastP2ButtonState {
				// Human player 2 controls (Flappy Bird style)
				if buttonP2Pressed {
					// Button pressed - JUMP UP (fast!)
					game.movePaddle(&game.player2, -1, JUMP_SPEED)
					game.spawnWindParticles() // Spawn particles on jump
				} else {
					// Button not pressed - FALL DOWN slowly (gravity)
					game.movePaddle(&game.player2, 1, FALL_SPEED)
				}
			} else {
				// AI controls player 2
				game.updateAI()
			}

			game.updateBall()
			game.updateWindParticles() // Update wind particles

			// Draw game
			game.draw(display)

			frameCount++
			if frameCount%90 == 0 { // changed from 60
				println("[PONG] Score:", game.score1, "-", game.score2)
			}
		} else {
			// Game over
			showWinner(display, game.winner)
			time.Sleep(3 * time.Second)

			// Reset for new game
			game = newGame(true)
			lastP2ButtonState = false // Reset to AI mode for new game
			println("[PONG] New game started")
		}

		// Frame delay (adjust for game speed)
		time.Sleep(50 * time.Millisecond) // 20 FPS
	}
}
