// Conway's Game of Life for SSD1306 OLED (128x64) using TinyGO
// Compatible with 0.96" SSD1306 OLED Display via I2C
package main

import (
	"image/color"
	"machine"
	"math/rand"
	"time"

	"tinygo.org/x/drivers/ssd1306"
)

const (
	Width  = 128
	Height = 64
)

// Grid represents the game board
type Grid struct {
	cells [Height][Width]bool
}

// NewGrid creates a new grid with random initial state
func NewGrid() *Grid {
	g := &Grid{}
	rand.Seed(time.Now().UnixNano())
	
	// Initialize with random cells (about 30% alive)
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			g.cells[y][x] = rand.Intn(100) < 30
		}
	}
	return g
}

// NewGridWithPattern creates a grid with a specific pattern
func NewGridWithPattern(pattern string) *Grid {
	g := &Grid{}
	
	switch pattern {
	case "glider":
		// Place a glider in the center
		cx, cy := Width/2, Height/2
		g.cells[cy][cx+1] = true
		g.cells[cy+1][cx+2] = true
		g.cells[cy+2][cx] = true
		g.cells[cy+2][cx+1] = true
		g.cells[cy+2][cx+2] = true
		
	case "blinker":
		// Place a blinker in the center
		cx, cy := Width/2, Height/2
		g.cells[cy][cx-1] = true
		g.cells[cy][cx] = true
		g.cells[cy][cx+1] = true
		
	case "toad":
		// Place a toad oscillator
		cx, cy := Width/2, Height/2
		g.cells[cy][cx] = true
		g.cells[cy][cx+1] = true
		g.cells[cy][cx+2] = true
		g.cells[cy+1][cx-1] = true
		g.cells[cy+1][cx] = true
		g.cells[cy+1][cx+1] = true
		
	case "pulsar":
		// Place a pulsar pattern
		cx, cy := Width/2, Height/2
		// Top half
		for i := 0; i < 3; i++ {
			g.cells[cy-6][cx-4+i] = true
			g.cells[cy-6][cx+2+i] = true
			g.cells[cy-1][cx-4+i] = true
			g.cells[cy-1][cx+2+i] = true
		}
		// Bottom half (mirror)
		for i := 0; i < 3; i++ {
			g.cells[cy+1][cx-4+i] = true
			g.cells[cy+1][cx+2+i] = true
			g.cells[cy+6][cx-4+i] = true
			g.cells[cy+6][cx+2+i] = true
		}
		// Left side
		for i := 0; i < 3; i++ {
			g.cells[cy-4+i][cx-6] = true
			g.cells[cy+2+i][cx-6] = true
			g.cells[cy-4+i][cx-1] = true
			g.cells[cy+2+i][cx-1] = true
		}
		// Right side
		for i := 0; i < 3; i++ {
			g.cells[cy-4+i][cx+1] = true
			g.cells[cy+2+i][cx+1] = true
			g.cells[cy-4+i][cx+6] = true
			g.cells[cy+2+i][cx+6] = true
		}
		
	case "lightweight_spaceship":
		// LWSS - moves horizontally
		cx, cy := Width/2, Height/2
		g.cells[cy][cx+1] = true
		g.cells[cy][cx+4] = true
		g.cells[cy+1][cx] = true
		g.cells[cy+2][cx] = true
		g.cells[cy+2][cx+4] = true
		g.cells[cy+3][cx] = true
		g.cells[cy+3][cx+1] = true
		g.cells[cy+3][cx+2] = true
		g.cells[cy+3][cx+3] = true
		
	case "gosper_glider_gun":
		// Famous pattern that continuously produces gliders
		// Left square
		g.cells[20][24] = true
		g.cells[20][25] = true
		g.cells[21][24] = true
		g.cells[21][25] = true
		
		// Left part
		g.cells[20][34] = true
		g.cells[21][34] = true
		g.cells[22][34] = true
		g.cells[19][35] = true
		g.cells[23][35] = true
		g.cells[18][36] = true
		g.cells[24][36] = true
		g.cells[18][37] = true
		g.cells[24][37] = true
		g.cells[21][38] = true
		g.cells[19][39] = true
		g.cells[23][39] = true
		g.cells[20][40] = true
		g.cells[21][40] = true
		g.cells[22][40] = true
		g.cells[21][41] = true
		
		// Right part
		g.cells[18][44] = true
		g.cells[19][44] = true
		g.cells[20][44] = true
		g.cells[18][45] = true
		g.cells[19][45] = true
		g.cells[20][45] = true
		g.cells[17][46] = true
		g.cells[21][46] = true
		g.cells[16][48] = true
		g.cells[17][48] = true
		g.cells[21][48] = true
		g.cells[22][48] = true
		
		// Right square
		g.cells[18][58] = true
		g.cells[19][58] = true
		g.cells[18][59] = true
		g.cells[19][59] = true
		
	case "explosion":
		// Creates chaotic explosions across the screen
		cx, cy := Width/2, Height/2
		// Multiple R-pentominos (famous for chaotic behavior)
		for i := 0; i < 3; i++ {
			ox, oy := cx-40+i*40, cy-10+i*10
			g.cells[oy][ox+1] = true
			g.cells[oy][ox+2] = true
			g.cells[oy+1][ox] = true
			g.cells[oy+1][ox+1] = true
			g.cells[oy+2][ox+1] = true
		}
		
	case "traffic_lights":
		// Multiple oscillators creating a light show
		for y := 10; y < Height-10; y += 15 {
			for x := 10; x < Width-10; x += 20 {
				// Blinker
				g.cells[y][x] = true
				g.cells[y][x+1] = true
				g.cells[y][x+2] = true
			}
		}
		for y := 18; y < Height-10; y += 15 {
			for x := 15; x < Width-10; x += 20 {
				// Toad
				g.cells[y][x] = true
				g.cells[y][x+1] = true
				g.cells[y][x+2] = true
				g.cells[y+1][x-1] = true
				g.cells[y+1][x] = true
				g.cells[y+1][x+1] = true
			}
		}
		
	case "acorn":
		// Small pattern that evolves for 5000+ generations
		cx, cy := Width/2, Height/2
		g.cells[cy][cx+1] = true
		g.cells[cy+1][cx+3] = true
		g.cells[cy+2][cx] = true
		g.cells[cy+2][cx+1] = true
		g.cells[cy+2][cx+4] = true
		g.cells[cy+2][cx+5] = true
		g.cells[cy+2][cx+6] = true
		
	case "fireworks":
		// Multiple gliders shooting in all directions
		cx, cy := Width/2, Height/2
		// Center explosion
		for i := 0; i < 8; i++ {
			angle := i * 45
			offsetX, offsetY := 0, 0
			switch angle {
			case 0: offsetX, offsetY = 15, 0
			case 45: offsetX, offsetY = 10, -10
			case 90: offsetX, offsetY = 0, -15
			case 135: offsetX, offsetY = -10, -10
			case 180: offsetX, offsetY = -15, 0
			case 225: offsetX, offsetY = -10, 10
			case 270: offsetX, offsetY = 0, 15
			case 315: offsetX, offsetY = 10, 10
			}
			x, y := cx+offsetX, cy+offsetY
			g.cells[y][x+1] = true
			g.cells[y+1][x+2] = true
			g.cells[y+2][x] = true
			g.cells[y+2][x+1] = true
			g.cells[y+2][x+2] = true
		}
		
	case "spaceship_fleet":
		// Multiple spaceships moving together
		for i := 0; i < 4; i++ {
			cx, cy := 20+i*25, 10+i*10
			// LWSS
			g.cells[cy][cx+1] = true
			g.cells[cy][cx+4] = true
			g.cells[cy+1][cx] = true
			g.cells[cy+2][cx] = true
			g.cells[cy+2][cx+4] = true
			g.cells[cy+3][cx] = true
			g.cells[cy+3][cx+1] = true
			g.cells[cy+3][cx+2] = true
			g.cells[cy+3][cx+3] = true
		}
		
	case "dense_chaos":
		// 50% density random - maximum chaos!
		for y := 0; y < Height; y++ {
			for x := 0; x < Width; x++ {
				g.cells[y][x] = rand.Intn(100) < 50
			}
		}
		
	default:
		// Random initialization
		return NewGrid()
	}
	
	return g
}

// CountNeighbors counts the live neighbors of a cell at (x, y)
func (g *Grid) CountNeighbors(x, y int) int {
	count := 0
	
	// Check all 8 neighbors with wrapping
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue // Skip the cell itself
			}
			
			// Wrap around the edges
			nx := (x + dx + Width) % Width
			ny := (y + dy + Height) % Height
			
			if g.cells[ny][nx] {
				count++
			}
		}
	}
	
	return count
}

// Next computes the next generation of the grid
func (g *Grid) Next() *Grid {
	next := &Grid{}
	
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			neighbors := g.CountNeighbors(x, y)
			alive := g.cells[y][x]
			
			// Apply Conway's Game of Life rules
			if alive {
				// Cell is alive
				if neighbors < 2 {
					// Underpopulation
					next.cells[y][x] = false
				} else if neighbors == 2 || neighbors == 3 {
					// Survival
					next.cells[y][x] = true
				} else {
					// Overpopulation
					next.cells[y][x] = false
				}
			} else {
				// Cell is dead
				if neighbors == 3 {
					// Reproduction
					next.cells[y][x] = true
				}
			}
		}
	}
	
	return next
}

// DrawToOLED renders the grid directly to the SSD1306 OLED display
func (g *Grid) DrawToOLED(display *ssd1306.Device) {
	// Clear the display buffer
	display.ClearBuffer()
	
	// Set each pixel based on cell state
	for y := int16(0); y < Height; y++ {
		for x := int16(0); x < Width; x++ {
			if g.cells[y][x] {
				display.SetPixel(x, y, color.RGBA{255, 255, 255, 255}) // White pixel
			}
		}
	}
	
	// Send buffer to display
	display.Display()
}

// CountLiveCells returns the number of live cells
func (g *Grid) CountLiveCells() int {
	count := 0
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			if g.cells[y][x] {
				count++
			}
		}
	}
	return count
}

// DrawText draws a simple 5x7 character at position (x, y)
func DrawText(display *ssd1306.Device, text string, x, y int16) {
	// Simple 3x5 font for basic characters
	for i, char := range text {
		drawChar(display, char, x+int16(i*4), y)
	}
}

// drawChar draws a single character using a tiny font
func drawChar(display *ssd1306.Device, char rune, x, y int16) {
	// Tiny 3x5 font patterns
	var pattern []byte
	switch char {
	case 'A': pattern = []byte{0x0E, 0x11, 0x1F, 0x11, 0x11}
	case 'B': pattern = []byte{0x1E, 0x11, 0x1E, 0x11, 0x1E}
	case 'C': pattern = []byte{0x0E, 0x11, 0x10, 0x11, 0x0E}
	case 'D': pattern = []byte{0x1E, 0x11, 0x11, 0x11, 0x1E}
	case 'E': pattern = []byte{0x1F, 0x10, 0x1E, 0x10, 0x1F}
	case 'F': pattern = []byte{0x1F, 0x10, 0x1E, 0x10, 0x10}
	case 'G': pattern = []byte{0x0E, 0x10, 0x17, 0x11, 0x0E}
	case 'H': pattern = []byte{0x11, 0x11, 0x1F, 0x11, 0x11}
	case 'I': pattern = []byte{0x0E, 0x04, 0x04, 0x04, 0x0E}
	case 'L': pattern = []byte{0x10, 0x10, 0x10, 0x10, 0x1F}
	case 'M': pattern = []byte{0x11, 0x1B, 0x15, 0x11, 0x11}
	case 'N': pattern = []byte{0x11, 0x19, 0x15, 0x13, 0x11}
	case 'O': pattern = []byte{0x0E, 0x11, 0x11, 0x11, 0x0E}
	case 'P': pattern = []byte{0x1E, 0x11, 0x1E, 0x10, 0x10}
	case 'R': pattern = []byte{0x1E, 0x11, 0x1E, 0x14, 0x12}
	case 'S': pattern = []byte{0x0E, 0x10, 0x0E, 0x01, 0x0E}
	case 'T': pattern = []byte{0x1F, 0x04, 0x04, 0x04, 0x04}
	case 'U': pattern = []byte{0x11, 0x11, 0x11, 0x11, 0x0E}
	case 'W': pattern = []byte{0x11, 0x11, 0x15, 0x1B, 0x11}
	case 'X': pattern = []byte{0x11, 0x0A, 0x04, 0x0A, 0x11}
	case 'Y': pattern = []byte{0x11, 0x0A, 0x04, 0x04, 0x04}
	case '0': pattern = []byte{0x0E, 0x13, 0x15, 0x19, 0x0E}
	case '1': pattern = []byte{0x04, 0x0C, 0x04, 0x04, 0x0E}
	case '2': pattern = []byte{0x0E, 0x11, 0x02, 0x04, 0x1F}
	case '3': pattern = []byte{0x1F, 0x02, 0x0E, 0x01, 0x1E}
	case '4': pattern = []byte{0x11, 0x11, 0x1F, 0x01, 0x01}
	case '5': pattern = []byte{0x1F, 0x10, 0x1E, 0x01, 0x1E}
	case '6': pattern = []byte{0x0E, 0x10, 0x1E, 0x11, 0x0E}
	case '7': pattern = []byte{0x1F, 0x01, 0x02, 0x04, 0x04}
	case '8': pattern = []byte{0x0E, 0x11, 0x0E, 0x11, 0x0E}
	case '9': pattern = []byte{0x0E, 0x11, 0x0F, 0x01, 0x0E}
	case '>': pattern = []byte{0x08, 0x04, 0x02, 0x04, 0x08}
	case ' ': pattern = []byte{0x00, 0x00, 0x00, 0x00, 0x00}
	case '-': pattern = []byte{0x00, 0x00, 0x0E, 0x00, 0x00}
	default: pattern = []byte{0x00, 0x00, 0x00, 0x00, 0x00}
	}
	
	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			if pattern[row]&(1<<uint(4-col)) != 0 {
				display.SetPixel(x+int16(col), y+int16(row), color.RGBA{255, 255, 255, 255})
			}
		}
	}
}

// ShowMenu displays the pattern selection menu
func ShowMenu(display *ssd1306.Device, patterns []string, selected int) {
	display.ClearBuffer()
	
	// Title
	DrawText(display, "GAME OF LIFE", 20, 2)
	DrawText(display, "SELECT PATTERN", 10, 10)
	
	// Show 5 items at a time
	startIdx := selected - 2
	if startIdx < 0 {
		startIdx = 0
	}
	if startIdx > len(patterns)-5 {
		startIdx = len(patterns) - 5
	}
	if startIdx < 0 {
		startIdx = 0
	}
	
	for i := 0; i < 5 && startIdx+i < len(patterns); i++ {
		y := int16(22 + i*8)
		idx := startIdx + i
		
		// Draw selection indicator
		if idx == selected {
			DrawText(display, ">", 2, y)
		}
		
		// Draw pattern name (truncate if needed)
		name := patterns[idx]
		if len(name) > 18 {
			name = name[:18]
		}
		DrawText(display, name, 10, y)
	}
	
	display.Display()
}

func main() {
	// Configure I2C pins to match your wiring
	// Your wiring: SDA=GPIO21, SCL=GPIO22 (standard ESP32)
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ, // 400kHz I2C speed
		SDA:       machine.GPIO21,          // Your SDA pin
		SCL:       machine.GPIO22,          // Your SCL pin
	})

	// Configure external button on GPIO18 for mode switching
	button := machine.GPIO18
	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// Initialize SSD1306 display
	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})

	display.ClearDisplay()

	// Available patterns - visually striking ones!
	patterns := []string{
		"RANDOM",
		"DENSE CHAOS",
		"EXPLOSION",
		"FIREWORKS",
		"TRAFFIC LIGHTS",
		"GLIDER GUN",
		"SPACESHIP FLEET",
		"ACORN",
		"PULSAR",
		"SPACESHIP",
		"GLIDER",
		"TOAD",
	}
	
	patternKeys := []string{
		"random",
		"dense_chaos",
		"explosion",
		"fireworks",
		"traffic_lights",
		"gosper_glider_gun",
		"spaceship_fleet",
		"acorn",
		"pulsar",
		"lightweight_spaceship",
		"glider",
		"toad",
	}
	
	// MENU MODE
	menuMode := true
	selectedPattern := 0
	lastButtonState := true
	buttonPressTime := time.Now()
	
	for menuMode {
		// Show menu
		ShowMenu(display, patterns, selectedPattern)
		
		// Check button
		buttonPressed := !button.Get()
		
		if buttonPressed && lastButtonState && time.Since(buttonPressTime) > 200*time.Millisecond {
			// Short press = next pattern
			selectedPattern = (selectedPattern + 1) % len(patterns)
			buttonPressTime = time.Now()
		} else if buttonPressed && time.Since(buttonPressTime) > 1000*time.Millisecond {
			// Long press (1 second) = select and start
			menuMode = false
			time.Sleep(300 * time.Millisecond) // Debounce
		}
		
		lastButtonState = buttonPressed
		time.Sleep(50 * time.Millisecond)
	}
	
	// GAME MODE
	grid := NewGridWithPattern(patternKeys[selectedPattern])
	generation := 0
	lastButtonState = true
	
	// Main game loop
	for {
		// Check button (pressed = LOW on pullup)
		buttonPressed := !button.Get()
		if buttonPressed && lastButtonState {
			// Button pressed - cycle to next pattern
			selectedPattern = (selectedPattern + 1) % len(patterns)
			grid = NewGridWithPattern(patternKeys[selectedPattern])
			generation = 0
			time.Sleep(200 * time.Millisecond) // Debounce
		}
		lastButtonState = buttonPressed
		
		// Draw current generation
		grid.DrawToOLED(display)
		
		// Compute next generation
		grid = grid.Next()
		generation++
		
		// Delay between frames
		time.Sleep(100 * time.Millisecond)
	}
}
