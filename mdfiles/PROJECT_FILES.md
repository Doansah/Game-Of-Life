# Project Files Overview

## Main Files

### 1. `tinygo_ssd1306_version.go` ⭐ **USE THIS FOR HARDWARE**
- **Purpose**: Production code for SSD1306 OLED display
- **Language**: TinyGO (subset of Go for microcontrollers)
- **Output**: Physical 128x64 OLED display via I2C
- **Flash**: `tinygo flash -target=pico tinygo_ssd1306_version.go`
- **Features**:
  - Direct pixel control of OLED
  - I2C communication
  - Hardware-optimized
  - Multiple pattern support
  - Runs on battery power

### 2. `gpt_version.go` - Terminal Testing Version
- **Purpose**: Development and pattern testing on computer
- **Language**: Standard Go
- **Output**: Terminal with ASCII art
- **Run**: `go run gpt_version.go`
- **Features**:
  - Interactive pattern selection
  - Generation counter
  - Live cell statistics
  - Compact display mode
  - Easy pattern debugging

## Documentation Files

### `README.md` - Complete Documentation
Comprehensive guide covering:
- Both versions (terminal & hardware)
- Installation instructions
- Pattern descriptions
- Code structure
- Performance metrics
- Troubleshooting

### `QUICKSTART.md` - Quick Reference
One-page guide with:
- Fast flash commands
- Wiring diagram
- Configuration options
- Troubleshooting table
- Performance comparison

### `SSD1306_PIN_GUIDE.md` - Hardware Reference
Detailed wiring guide:
- Pin mappings for all boards
- I2C address detection
- Voltage considerations
- Pull-up resistor info
- Multiple device setup

### `conwayGOLRules.txt` - Game Rules
Original specification of Conway's Game of Life rules

## Utility Files

### `flash.sh` - Interactive Flash Script
Bash script that:
- Checks TinyGO installation
- Guides board selection
- Handles pattern selection
- Provides wiring reminders
- Monitors serial output

Usage: `./flash.sh`

### `go.mod` - Go Module File
Defines the Go module for the terminal version

## Directory Structure

```
go1stproject/
├── tinygo_ssd1306_version.go  ← Hardware version (MAIN)
├── gpt_version.go              ← Terminal test version
├── screen_display.go           ← Old prototype (ignore)
├── README.md                   ← Full documentation
├── QUICKSTART.md               ← Quick reference
├── SSD1306_PIN_GUIDE.md        ← Wiring guide
├── PROJECT_FILES.md            ← This file
├── conwayGOLRules.txt          ← Game rules
├── flash.sh                    ← Flash helper script
├── go.mod                      ← Go module
└── GameOfLife/                 ← (folder)
```

## Development Workflow

### Recommended Process:

1. **Test Patterns on Computer**
   ```bash
   go run gpt_version.go
   ```
   - Select different patterns
   - Observe behavior
   - Verify rules are working
   - Check for interesting patterns

2. **Modify Pattern (Optional)**
   Edit `tinygo_ssd1306_version.go` line ~243:
   ```go
   grid := NewGridWithPattern("glider")  // Change pattern here
   ```

3. **Flash to Hardware**
   ```bash
   # Using script (recommended for beginners)
   ./flash.sh
   
   # Or direct command
   tinygo flash -target=pico tinygo_ssd1306_version.go
   ```

4. **Test on OLED**
   - Connect power
   - Watch simulation run
   - Adjust speed if needed

5. **Iterate**
   - Modify code
   - Reflash
   - Test again

## File Size Comparison

| File | Lines | Size | Purpose |
|------|-------|------|---------|
| `tinygo_ssd1306_version.go` | ~270 | 8KB | Production hardware code |
| `gpt_version.go` | ~280 | 9KB | Development/testing |
| `README.md` | ~200 | 10KB | Full documentation |
| `QUICKSTART.md` | ~150 | 5KB | Quick reference |
| `SSD1306_PIN_GUIDE.md` | ~200 | 8KB | Wiring guide |

## Which File Should I Use?

### For Testing Patterns:
→ **`gpt_version.go`**
- Run on your computer
- See patterns in terminal
- Quick iteration

### For Physical OLED Display:
→ **`tinygo_ssd1306_version.go`**
- Flash to microcontroller
- Real hardware display
- Production ready

### For Learning Wiring:
→ **`SSD1306_PIN_GUIDE.md`**
- Pin mappings for your board
- Voltage information
- Troubleshooting hardware

### For Quick Commands:
→ **`QUICKSTART.md`**
- Fast reference
- Common commands
- Configuration options

### For Complete Understanding:
→ **`README.md`**
- Full project documentation
- Code structure explanation
- Performance metrics

## Code Differences

### Terminal Version (`gpt_version.go`)
```go
// Uses standard Go
import "fmt"

// Display function
func (g *Grid) DisplayCompact() {
    fmt.Print("\033[H\033[2J")  // ANSI codes
    // ... ASCII art rendering
}

// Can use float32, full stdlib
rand.Float32() < 0.3
```

### Hardware Version (`tinygo_ssd1306_version.go`)
```go
// Uses TinyGO + hardware drivers
import "machine"
import "tinygo.org/x/drivers/ssd1306"

// Display function
func (g *Grid) DrawToOLED(display *ssd1306.Device) {
    display.ClearBuffer()
    display.SetPixel(x, y, true)
    display.Display()
}

// Must use TinyGO-compatible functions
rand.Intn(100) < 30
```

## Next Steps

1. ✅ **You're here** - Understanding the project files
2. 🔌 **Wire up** your OLED display (see `SSD1306_PIN_GUIDE.md`)
3. 🚀 **Flash** the code (use `./flash.sh` or manual command)
4. 🎮 **Watch** Conway's Game of Life run!
5. 🎨 **Customize** patterns and settings
6. 📦 **Build** an enclosure (optional)

## Getting Help

- **Hardware issues?** → Check `SSD1306_PIN_GUIDE.md`
- **Quick command?** → Check `QUICKSTART.md`
- **Understanding code?** → Check `README.md`
- **Pattern ideas?** → Visit https://conwaylife.com/wiki/
- **TinyGO issues?** → Visit https://tinygo.org/docs/

## Contributing

This is your project! Feel free to:
- Add new patterns
- Optimize performance
- Add button controls
- Create custom displays
- Build an enclosure
- Share your build!

---

**Start here**: `QUICKSTART.md` for fastest path to running hardware!
