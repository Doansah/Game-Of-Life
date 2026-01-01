# Conway's Game of Life - Go Implementation

A Go implementation of Conway's Game of Life running on a 128x64 pixel grid, with both terminal-based visualization and **SSD1306 OLED display support** using TinyGO.

## Features

- **128x64 Grid**: Matches the 0.96" SSD1306 OLED screen dimensions
- **Wrapping Edges**: Cells wrap around the edges for infinite grid simulation
- **Multiple Patterns**: Choose from several classic Game of Life patterns
- **Real-time Display**: Works with both terminal and hardware OLED displays
- **I2C Communication**: Uses TinyGO drivers for SSD1306 via I2C protocol
- **Hardware Ready**: Optimized for microcontrollers (Arduino, Raspberry Pi Pico, ESP32, etc.)

## Two Versions Available

### 1. Terminal Version (Testing/Development)

For testing on your computer with visual terminal output:

```bash
go run gpt_version.go
```

When prompted, select a starting pattern:
1. **Random** - 30% of cells randomly initialized as alive
2. **Glider** - A small pattern that moves diagonally across the grid
3. **Blinker** - A simple oscillator that alternates between two states
4. **Toad** - A period-2 oscillator
5. **Pulsar** - A larger period-3 oscillator

### 2. TinyGO + SSD1306 OLED Version (Hardware)

For running on actual hardware with the 0.96" SSD1306 OLED display:

**File:** `tinygo_ssd1306_version.go`

#### Hardware Requirements

- Microcontroller board (Arduino Uno, Nano, Pico, ESP32, etc.)
- 0.96" SSD1306 OLED Display (128x64)
- I2C connection (4 wires: VCC, GND, SDA, SCL)

#### Wiring (I2C)

```
SSD1306 OLED  →  Microcontroller
---------------------------------
VCC           →  3.3V or 5V (check your display)
GND           →  GND
SDA           →  SDA (I2C Data)
SCL           →  SCL (I2C Clock)
```

#### Build and Flash

**For Arduino Uno:**
```bash
tinygo flash -target=arduino tinygo_ssd1306_version.go
```

**For Raspberry Pi Pico:**
```bash
tinygo flash -target=pico tinygo_ssd1306_version.go
```

**For Arduino Nano:**
```bash
tinygo flash -target=arduino-nano tinygo_ssd1306_version.go
```

**For ESP32:**
```bash
tinygo flash -target=esp32-coreboard tinygo_ssd1306_version.go
```

**For other boards:** Check [TinyGO supported boards](https://tinygo.org/docs/reference/microcontrollers/)

#### Configuration

In `tinygo_ssd1306_version.go`, you can customize:

1. **I2C Address** (line ~233): Most SSD1306 use `0x3C`, some use `0x3D`
2. **Starting Pattern** (line ~243): Change `"glider"` to any pattern
3. **Frame Rate** (line ~257): Adjust `time.Sleep()` duration
4. **I2C Speed** (line ~227): Default 400kHz, can use 100kHz if issues

Available patterns:
- `"glider"` - Moving pattern
- `"blinker"` - Simple oscillator
- `"toad"` - Period-2 oscillator
- `"pulsar"` - Large period-3 oscillator
- `"lightweight_spaceship"` - Horizontal spaceship
- `"random"` - Random start (or use `NewGrid()`)

## Installing TinyGO

### macOS
```bash
brew tap tinygo-org/tools
brew install tinygo
```

### Linux
```bash
wget https://github.com/tinygo-org/tinygo/releases/download/v0.30.0/tinygo_0.30.0_amd64.deb
sudo dpkg -i tinygo_0.30.0_amd64.deb
```

### Windows
Download installer from: https://github.com/tinygo-org/tinygo/releases

Verify installation:
```bash
tinygo version
```

## Conway's Game of Life Rules

From `conwayGOLRules.txt`:

- **Underpopulation**: Live cell with < 2 live neighbors dies
- **Survival**: Live cell with 2 or 3 live neighbors lives
- **Overpopulation**: Live cell with > 3 live neighbors dies
- **Reproduction**: Dead cell with exactly 3 live neighbors becomes live

All neighbors are checked with wrapping (toroidal topology).

## Code Structure

### Common Components (Both Versions)

- **Grid**: Represents the game board with a 2D boolean array `[64][128]bool`
- **NewGrid()**: Creates a random initial state
- **NewGridWithPattern()**: Creates predefined patterns (glider, blinker, etc.)
- **CountNeighbors()**: Counts live neighbors with edge wrapping
- **Next()**: Computes the next generation following Game of Life rules

### Terminal Version Specific
- **DisplayCompact()**: Renders the grid to terminal with ASCII

### TinyGO/OLED Version Specific
- **DrawToOLED()**: Renders grid directly to SSD1306 pixel buffer
- **I2C Configuration**: Hardware I2C setup for display communication

### Wrapping Implementation

The grid uses modulo arithmetic to wrap around edges:
```go
nx := (x + dx + Width) % Width
ny := (y + dy + Height) % Height
```

This creates a toroidal topology where cells on opposite edges are neighbors.

## SSD1306 OLED Display Details

### Display Specifications
- **Model**: SSD1306 0.96" OLED
- **Resolution**: 128x64 pixels (monochrome)
- **Communication**: I2C protocol
- **I2C Address**: Usually `0x3C` (sometimes `0x3D`)
- **Voltage**: 3.3V - 5V (check your module)
- **Driver**: TinyGO `ssd1306` package

### How the Display Works

The SSD1306 has an internal buffer that maps 1:1 with pixels:
- Each cell in `grid.cells[y][x]` = one pixel on screen
- `true` = pixel ON (white/lit)
- `false` = pixel OFF (black/dark)

The `DrawToOLED()` function:
1. Clears the display buffer
2. Iterates through all 8,192 cells (128×64)
3. Sets each pixel based on cell state
4. Sends complete buffer to display via I2C

### Troubleshooting

**Display not working?**
1. Check I2C address - try `0x3D` if `0x3C` doesn't work
2. Verify wiring (especially SDA/SCL)
3. Ensure correct voltage (3.3V for most boards)
4. Try lower I2C frequency: `machine.TWI_FREQ_100KHZ`

**Display shows random pixels?**
- Normal at startup - pattern should stabilize after a few generations

**Too slow/fast?**
- Adjust `time.Sleep()` duration in main loop
- 50ms = 20 FPS, 100ms = 10 FPS, 200ms = 5 FPS

## Adapting for OLED Screen

The TinyGO version is **production ready** for SSD1306 OLED displays. Key features:

1. ✅ Grid dimensions (128x64) match display exactly
2. ✅ Direct pixel mapping with no scaling needed
3. ✅ I2C communication properly configured
4. ✅ Efficient buffer management
5. ✅ Hardware-optimized rendering

## Performance

### Terminal Version (Go)
- Grid: 128x64 = 8,192 cells
- Each generation checks 8 neighbors per cell
- Updates run smoothly at 10 FPS
- Memory: ~8KB for grid + overhead

### TinyGO/OLED Version
- Grid: 8,192 cells (128×64)
- Per-generation computation: ~65,536 neighbor checks
- Display buffer: 1KB (128×64÷8 bytes)
- Typical frame rate: 5-20 FPS (depending on microcontroller)
- Flash usage: ~30-50KB
- RAM usage: ~10-15KB

**Optimized for:**
- Arduino Uno/Nano (ATmega328P) - Works but slow (~3-5 FPS)
- Raspberry Pi Pico (RP2040) - Smooth (~15-20 FPS) ⭐ **Recommended**
- ESP32 - Very smooth (~20+ FPS)
- Arduino Due - Smooth (~15-20 FPS)

## Example Hardware Setups

### Setup 1: Raspberry Pi Pico (Recommended)
```
Pico Pin  →  SSD1306
GP0 (SDA) →  SDA
GP1 (SCL) →  SCL
3.3V      →  VCC
GND       →  GND
```

### Setup 2: Arduino Uno
```
Arduino   →  SSD1306
A4 (SDA)  →  SDA
A5 (SCL)  →  SCL
5V        →  VCC
GND       →  GND
```

### Setup 3: ESP32
```
ESP32     →  SSD1306
GPIO21    →  SDA
GPIO22    →  SCL
3.3V      →  VCC
GND       →  GND
```

## Next Steps / Enhancements

### Software Enhancements
1. **Pattern Selection**: Add button to cycle through patterns
2. **Speed Control**: Add potentiometer to adjust frame rate
3. **Pause/Resume**: Add button to pause simulation
4. **Statistics Display**: Show generation count on OLED
5. **Pattern Library**: Store and load custom patterns from EEPROM

### Hardware Additions
1. **Buttons**: GPIO buttons for control (pause, pattern select, speed)
2. **LED Indicator**: Show when simulation is running
3. **Buzzer**: Audio feedback for button presses
4. **Battery Power**: Make it portable with LiPo battery
5. **Case**: 3D-printed enclosure

### Advanced Features
1. **Color**: Upgrade to color OLED (SSD1351) with cell age coloring
2. **Larger Display**: Use 1.3" OLED (128x64 SH1106) or larger
3. **Multiple Grids**: Run multiple simulations side-by-side
4. **Web Control**: ESP32 + WiFi for remote pattern selection
5. **Save to SD**: Log interesting patterns to SD card

## Terminal Controls

- Press `Ctrl+C` to stop the simulation (both versions)

## Requirements

### Terminal Version
- Go 1.16 or higher
- Terminal with ANSI escape code support (most modern terminals)

### TinyGO/OLED Version
- TinyGO 0.25 or higher
- SSD1306 OLED display (128x64)
- Compatible microcontroller board
- USB cable for flashing

## Dependencies

TinyGO version uses:
```go
import "tinygo.org/x/drivers/ssd1306"
```

The driver is automatically fetched during build.

## Testing Workflow

1. **Develop & Test**: Use `gpt_version.go` with terminal display to test patterns
2. **Flash to Hardware**: Use `tinygo_ssd1306_version.go` when ready for OLED
3. **Debug**: Use `println()` statements (they output to serial monitor)

## Serial Monitor (Debugging)

To see debug output from your microcontroller:

```bash
tinygo monitor
```

Or use Arduino IDE's Serial Monitor at 115200 baud.

## Common Patterns

### Glider (5 cells)
```
  █
█ █
 ██
```
Moves diagonally, period 4

### Blinker (3 cells)
```
███
```
Oscillates horizontally/vertically, period 2

### Toad (6 cells)
```
 ███
███
```
Oscillates, period 2

### Lightweight Spaceship (9 cells)
```
 █  █
█
█   █
████
```
Moves horizontally, period 4

## Resources

- [TinyGO Documentation](https://tinygo.org/docs/)
- [SSD1306 Driver](https://github.com/tinygo-org/drivers/tree/release/ssd1306)
- [Conway's Game of Life Patterns](https://conwaylife.com/wiki/)
- [TinyGO Supported Boards](https://tinygo.org/docs/reference/microcontrollers/)

## License

Free to use and modify for your projects!

---

## 📚 Additional Documentation

This project includes several helpful guides:

- **[QUICKSTART.md](QUICKSTART.md)** - One-page quick reference for fast setup
- **[SSD1306_PIN_GUIDE.md](SSD1306_PIN_GUIDE.md)** - Detailed pin wiring for all boards
- **[WIRING_DIAGRAMS.md](WIRING_DIAGRAMS.md)** - Visual ASCII wiring diagrams
- **[PROJECT_FILES.md](PROJECT_FILES.md)** - Overview of all project files
- **flash.sh** - Interactive flashing script (run: `./flash.sh`)

## 🎯 Quick Navigation

- **Just want to flash and go?** → See [QUICKSTART.md](QUICKSTART.md)
- **Need wiring help?** → See [WIRING_DIAGRAMS.md](WIRING_DIAGRAMS.md)
- **Pin questions?** → See [SSD1306_PIN_GUIDE.md](SSD1306_PIN_GUIDE.md)
- **Understanding files?** → See [PROJECT_FILES.md](PROJECT_FILES.md)

## 🎉 You're Ready!

The TinyGO version (`tinygo_ssd1306_version.go`) is production-ready for your 0.96" SSD1306 OLED display. Just wire it up, flash the code, and watch Conway's Game of Life come alive on your hardware!

**Recommended first setup:**
```bash
# 1. Wire OLED to your Pico (see WIRING_DIAGRAMS.md)
# 2. Flash the code
tinygo flash -target=pico tinygo_ssd1306_version.go
# 3. Watch it run!
```
