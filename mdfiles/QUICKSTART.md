# Quick Start Guide - Game of Life on SSD1306

## 🚀 Quick Flash (One Command)

```bash
# Raspberry Pi Pico (RECOMMENDED)
tinygo flash -target=pico tinygo_ssd1306_version.go

# Arduino Uno
tinygo flash -target=arduino tinygo_ssd1306_version.go

# Arduino Nano
tinygo flash -target=arduino-nano tinygo_ssd1306_version.go

# ESP32
tinygo flash -target=esp32-coreboard tinygo_ssd1306_version.go
```

## 🔌 Wiring (4 wires)

```
OLED → Microcontroller
─────────────────────
VCC  → 3.3V or 5V
GND  → GND
SDA  → SDA pin
SCL  → SCL pin
```

**Most common I2C pins:**
- **Pico**: SDA=GP0, SCL=GP1
- **Arduino Uno/Nano**: SDA=A4, SCL=A5
- **ESP32**: SDA=GPIO21, SCL=GPIO22

## ⚙️ Configuration Options

Edit `tinygo_ssd1306_version.go`:

### Change Pattern (line ~243)
```go
// Options: "glider", "blinker", "toad", "pulsar", "lightweight_spaceship", "random"
grid := NewGridWithPattern("glider")
```

### Change I2C Address (line ~233)
```go
Address: 0x3C,  // Try 0x3D if display doesn't work
```

### Change Speed (line ~257)
```go
time.Sleep(100 * time.Millisecond)  // 100ms = 10 FPS
time.Sleep(50 * time.Millisecond)   // 50ms = 20 FPS
time.Sleep(200 * time.Millisecond)  // 200ms = 5 FPS
```

### Change I2C Speed (line ~227)
```go
Frequency: machine.TWI_FREQ_400KHZ,  // Fast (400 kHz)
Frequency: machine.TWI_FREQ_100KHZ,  // Standard (100 kHz) - use if issues
```

## 🎨 Available Patterns

| Pattern | Description | Size | Behavior |
|---------|-------------|------|----------|
| `glider` | Small moving pattern | 5 cells | Moves diagonally |
| `blinker` | Simple oscillator | 3 cells | Flips horizontal/vertical |
| `toad` | Medium oscillator | 6 cells | Period-2 oscillation |
| `pulsar` | Large oscillator | 48 cells | Period-3 oscillation |
| `lightweight_spaceship` | Spaceship | 9 cells | Moves horizontally |
| `random` | Random cells | ~30% filled | Chaotic evolution |

## 🛠️ Troubleshooting

| Problem | Solution |
|---------|----------|
| Display blank | • Check wiring (VCC, GND, SDA, SCL)<br>• Try I2C address 0x3D<br>• Check voltage (3.3V vs 5V) |
| Garbled pixels | • Lower I2C frequency to 100kHz<br>• Use shorter wires<br>• Add 4.7kΩ pull-up resistors |
| Can't flash | • Install USB drivers<br>• Check USB cable (data capable)<br>• Pico: Hold BOOTSEL while plugging in |
| Too slow | • Normal for Arduino Uno<br>• Increase sleep duration<br>• Upgrade to Pico/ESP32 |
| Code won't compile | • Check TinyGO version (need 0.25+)<br>• Run: `tinygo version` |

## 📊 Performance Guide

| Board | CPU | FPS | Status |
|-------|-----|-----|--------|
| **Raspberry Pi Pico** | RP2040 133MHz | 15-20 | ⭐ Best |
| **ESP32** | Dual-core 240MHz | 20+ | ⭐ Excellent |
| Arduino Due | ARM 84MHz | 15-20 | ✅ Good |
| Arduino Nano | ATmega328 16MHz | 8-10 | ⚠️ OK |
| Arduino Uno | ATmega328 8MHz | 3-5 | ⚠️ Slow |

## 🔍 Debug Serial Output

To see generation count and stats over serial:

```bash
tinygo monitor
```

Uncomment this line in code (line ~254):
```go
println("Generation:", generation, "Live Cells:", grid.CountLiveCells())
```

## 📦 Dependencies

TinyGO automatically fetches:
```go
import "tinygo.org/x/drivers/ssd1306"
```

No manual installation needed!

## 🎯 Recommended Setup

**Best experience:**
- **Board**: Raspberry Pi Pico ($4)
- **Display**: 0.96" SSD1306 OLED ($3-5)
- **Wiring**: 4 jumper wires
- **Pattern**: Start with "glider"
- **Speed**: 100ms (10 FPS)

**Total cost**: ~$10

## 📖 More Help

- **Full README**: `README.md`
- **Wiring Guide**: `SSD1306_PIN_GUIDE.md`
- **Flash Script**: `./flash.sh` (interactive)
- **TinyGO Docs**: https://tinygo.org/docs/

## 🎮 Usage

1. **Wire** OLED to microcontroller (4 wires)
2. **Flash** code: `tinygo flash -target=pico tinygo_ssd1306_version.go`
3. **Watch** Conway's Game of Life run!
4. **Modify** patterns in code and reflash

---

**That's it!** The simulation runs forever (or until you unplug it). Enjoy watching the cellular automaton evolve! 🎉
