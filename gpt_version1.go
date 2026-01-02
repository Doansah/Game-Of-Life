// This GO application will setup conways game of life on a screen, 128x64 pixels
package main

import (
	"fmt"
	"math/rand"
	"time"
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

	// Initialize with random cells
	// rand.Float32 generates a number between 0-1,
	// if random number is less than 0.3 true..
	// essentially returns true 30% of the time
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			g.cells[y][x] = rand.Float32() < 0.3
		}
	}
	return g
}

// Auto-generated Patterns
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
			// MODULAR ARITHMETIC
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

			// Conway's Game of Life rules
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

// Display renders the grid to the terminal
func (g *Grid) Display() {
	// Clear screen and move cursor to top-left
	fmt.Print("\033[H\033[2J")

	// Print top border
	fmt.Print("┌")
	for i := 0; i < Width; i++ {
		fmt.Print("─")
	}
	fmt.Println("┐")

	// Print grid (sample every 2 columns to fit on screen better)
	for y := 0; y < Height; y++ {
		fmt.Print("│")
		for x := 0; x < Width; x++ {
			if g.cells[y][x] {
				fmt.Print("█") // Live cell
			} else {
				fmt.Print(" ") // Dead cell
			}
		}
		fmt.Println("│")
	}

	// Print bottom border
	fmt.Print("└")
	for i := 0; i < Width; i++ {
		fmt.Print("─")
	}
	fmt.Println("┘")
}

// DisplayCompact renders a compact version of the grid
func (g *Grid) DisplayCompact() {
	// Clear screen and move cursor to top-left
	fmt.Print("\033[H\033[2J")

	fmt.Println("Conway's Game of Life - 128x64 Grid")
	fmt.Println("====================================")

	// Sample every 4th row and 2nd column for compact display
	for y := 0; y < Height; y += 2 {
		for x := 0; x < Width; x += 2 {
			if g.cells[y][x] {
				fmt.Print("█")
			} else {
				fmt.Print("·")
			}
		}
		fmt.Println()
	}
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

func main() {
	fmt.Println("Conway's Game of Life - Go Implementation")
	fmt.Println("=========================================")
	fmt.Println("\nChoose a starting pattern:")
	fmt.Println("1. Random")
	fmt.Println("2. Glider")
	fmt.Println("3. Blinker")
	fmt.Println("4. Toad")
	fmt.Println("5. Pulsar")
	fmt.Print("\nEnter choice (1-5): ")

	var choice int
	fmt.Scanln(&choice)

	var grid *Grid
	switch choice {
	case 2:
		grid = NewGridWithPattern("glider")
	case 3:
		grid = NewGridWithPattern("blinker")
	case 4:
		grid = NewGridWithPattern("toad")
	case 5:
		grid = NewGridWithPattern("pulsar")
	default:
		grid = NewGrid()
	}

	fmt.Println("\nStarting simulation... Press Ctrl+C to stop.")
	time.Sleep(2 * time.Second)

	generation := 0

	// Run the game loop
	for {
		// Display the current generation
		grid.DisplayCompact()
		fmt.Printf("\nGeneration: %d | Live Cells: %d\n", generation, grid.CountLiveCells())

		// Compute next generation
		grid = grid.Next()
		generation++

		// Wait before next frame
		time.Sleep(100 * time.Millisecond)
	}
}
