# Conway's Game of Life - Go Implementation

A Go implementation of Conway's Game of Life running on a 128x64 pixel grid, with terminal-based visualization.

## Features

- **128x64 Grid**: Matches your target OLED screen dimensions
- **Wrapping Edges**: Cells wrap around the edges for infinite grid simulation
- **Multiple Patterns**: Choose from several classic Game of Life patterns
- **Real-time Display**: Animated terminal visualization with generation counter
- **Compact View**: Efficient display that fits in most terminal windows

## How to Run

```bash
go run gpt_version.go
```

When prompted, select a starting pattern:
1. **Random** - 30% of cells randomly initialized as alive
2. **Glider** - A small pattern that moves diagonally across the grid
3. **Blinker** - A simple oscillator that alternates between two states
4. **Toad** - A period-2 oscillator
5. **Pulsar** - A larger period-3 oscillator

## Conway's Game of Life Rules

From `conwayGOLRules.txt`:

- **Underpopulation**: Live cell with < 2 live neighbors dies
- **Survival**: Live cell with 2 or 3 live neighbors lives
- **Overpopulation**: Live cell with > 3 live neighbors dies
- **Reproduction**: Dead cell with exactly 3 live neighbors becomes live

## Display

The program uses a compact display mode that:
- Samples every 2nd column and 2nd row (64x32 displayed from 128x64 grid)
- Uses `█` for live cells and `·` for dead cells
- Updates 10 times per second (100ms per frame)
- Shows generation count and live cell count

## Code Structure

### Main Components

- **Grid**: Represents the game board with a 2D boolean array
- **NewGrid()**: Creates a random initial state
- **NewGridWithPattern()**: Creates predefined patterns (glider, blinker, etc.)
- **CountNeighbors()**: Counts live neighbors with edge wrapping
- **Next()**: Computes the next generation following Game of Life rules
- **DisplayCompact()**: Renders the grid to terminal

### Wrapping Implementation

The grid uses modulo arithmetic to wrap around edges:
```go
nx := (x + dx + Width) % Width
ny := (y + dy + Height) % Height
```

This creates a toroidal topology where cells on opposite edges are neighbors.

## Adapting for OLED Screen

The code is designed to be easily adapted for an OLED display:

1. The grid dimensions (128x64) match common OLED screens
2. The `cells` array can be directly mapped to pixel data
3. Replace `DisplayCompact()` with your OLED driver code
4. The boolean array can be converted to monochrome bitmap data

Example OLED adaptation:
```go
func (g *Grid) SendToOLED(display *oled.Display) {
    for y := 0; y < Height; y++ {
        for x := 0; x < Width; x++ {
            display.SetPixel(x, y, g.cells[y][x])
        }
    }
    display.Show()
}
```

## Performance

- Grid: 128x64 = 8,192 cells
- Each generation checks 8 neighbors per cell
- Updates run smoothly at 10 FPS in terminal mode
- Minimal memory footprint (~8KB for grid)

## Next Steps

1. **OLED Integration**: Replace terminal display with OLED driver
2. **Buttons**: Add GPIO buttons to select patterns or pause/resume
3. **Color**: If using color OLED, add cell age coloring
4. **Save/Load**: Implement pattern saving and loading
5. **Speed Control**: Add variable speed control

## Terminal Controls

- Press `Ctrl+C` to stop the simulation

## Requirements

- Go 1.16 or higher
- Terminal with ANSI escape code support (most modern terminals)

## License

Free to use and modify for your projects!
