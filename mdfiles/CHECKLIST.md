# Installation Checklist

Use this checklist to ensure everything is set up correctly for your Conway's Game of Life on SSD1306 OLED.

## ☐ Prerequisites

### Software
- [ ] **TinyGO installed** (0.25 or higher)
  ```bash
  tinygo version
  ```
  - macOS: `brew install tinygo`
  - Linux/Windows: See https://tinygo.org/getting-started/install/

- [ ] **Go installed** (1.16+) - Optional, for terminal testing
  ```bash
  go version
  ```

- [ ] **Git** (for downloading if needed)
  ```bash
  git --version
  ```

### Hardware
- [ ] **Microcontroller board** (choose one):
  - [ ] Raspberry Pi Pico (RP2040) ⭐ **RECOMMENDED**
  - [ ] Arduino Uno
  - [ ] Arduino Nano
  - [ ] ESP32
  - [ ] Arduino Due
  - [ ] Other TinyGO-compatible board

- [ ] **SSD1306 OLED Display** (0.96", 128x64, I2C)
  - [ ] 4 pins: VCC, GND, SDA, SCL
  - [ ] Voltage: 3.3V-5V (check your specific model)

- [ ] **USB Cable** (data-capable, not just charging)
  - [ ] Micro-USB for Arduino/ESP32
  - [ ] USB-C for newer Picos
  - [ ] USB-B for Arduino Uno

- [ ] **4 Jumper Wires** (male-to-female or male-to-male depending on board)
  - [ ] Red (VCC)
  - [ ] Black (GND)
  - [ ] Green or Blue (SCL)
  - [ ] Yellow or White (SDA)

## ☐ Hardware Setup

### Wiring
- [ ] **Identified correct pins** on your microcontroller
  - See `SSD1306_PIN_GUIDE.md` for your specific board
  
- [ ] **Connected VCC** (power)
  - OLED VCC → Board 3.3V or 5V
  - [ ] Verified voltage is correct for your OLED module
  
- [ ] **Connected GND** (ground)
  - OLED GND → Board GND
  
- [ ] **Connected SDA** (I2C data)
  - OLED SDA → Board SDA pin
  - Pico: GP0, Arduino: A4, ESP32: GPIO21
  
- [ ] **Connected SCL** (I2C clock)
  - OLED SCL → Board SCL pin
  - Pico: GP1, Arduino: A5, ESP32: GPIO22

### Physical Check
- [ ] **Wires secured** (not loose)
- [ ] **No shorts** (wires not touching each other)
- [ ] **OLED display clean** (no scratches on screen)
- [ ] **Board powered off** while wiring

## ☐ Software Setup

### Files Downloaded
- [ ] **Project files downloaded/cloned**
  ```bash
  cd "/Users/dillonansah/Desktop/Programming Projects/go1stproject"
  ```

- [ ] **Main file present**: `tinygo_ssd1306_version.go`
- [ ] **Documentation present**: README.md, QUICKSTART.md, etc.

### Configuration (Optional)
- [ ] **I2C Address set** (line ~233 in code)
  - [ ] Default: 0x3C
  - [ ] Alternative: 0x3D (if display doesn't work)

- [ ] **Pattern selected** (line ~243 in code)
  - [ ] glider ✓ (default)
  - [ ] blinker
  - [ ] toad
  - [ ] pulsar
  - [ ] lightweight_spaceship
  - [ ] random

- [ ] **Speed adjusted** (line ~257 in code)
  - [ ] 100ms = 10 FPS ✓ (default)
  - [ ] 50ms = 20 FPS (faster)
  - [ ] 200ms = 5 FPS (slower)

## ☐ First Flash

### Preparation
- [ ] **Board connected** via USB
- [ ] **USB cable verified** (data-capable)
- [ ] **Drivers installed** (if needed for your board)
  - Windows: May need CH340, FTDI, or CP210x drivers
  - macOS/Linux: Usually automatic

### For Raspberry Pi Pico Only
- [ ] **BOOTSEL mode entered**:
  1. [ ] Hold BOOTSEL button
  2. [ ] Plug in USB
  3. [ ] Release BOOTSEL
  4. [ ] Pico appears as USB drive named "RPI-RP2"

### Flash Command
- [ ] **Choose your method**:

  **Option A: Interactive Script** (Recommended for beginners)
  ```bash
  ./flash.sh
  ```
  - [ ] Selected correct board
  - [ ] Selected pattern
  - [ ] Flash completed successfully

  **Option B: Direct Command**
  ```bash
  # Replace "pico" with your board target
  tinygo flash -target=pico tinygo_ssd1306_version.go
  ```
  - [ ] Command executed without errors
  - [ ] Flash completed successfully

### Verification
- [ ] **Code flashed** (no error messages)
- [ ] **Board restarted** automatically
- [ ] **OLED display powered on**

## ☐ Testing

### Visual Check
- [ ] **OLED shows pattern** (not blank)
- [ ] **Cells visible** (white pixels on black background)
- [ ] **Animation running** (cells changing over time)
- [ ] **No flickering** (smooth transitions)

### Behavior Check
- [ ] **Pattern evolves** according to Game of Life rules
- [ ] **Cells wrap around** edges (toroidal topology)
- [ ] **Performance acceptable** (5-20 FPS depending on board)

### Serial Monitor (Optional)
- [ ] **Serial monitor opened**
  ```bash
  tinygo monitor
  ```
- [ ] **Debug output visible** (if uncommented in code)

## ☐ Troubleshooting

### If Display is Blank
- [ ] **Power connected** (checked with multimeter if available)
- [ ] **Correct voltage** (3.3V or 5V as appropriate)
- [ ] **I2C address changed** to 0x3D (if 0x3C doesn't work)
- [ ] **Wires reseated** (unplugged and plugged back in)
- [ ] **Display tested** with different board/code

### If Garbled Display
- [ ] **I2C frequency lowered** to 100kHz (in code)
- [ ] **Wires shortened** (under 6 inches)
- [ ] **Pull-up resistors added** (4.7kΩ) if needed

### If Can't Flash
- [ ] **TinyGO target correct** for your board
- [ ] **USB drivers installed** (Windows especially)
- [ ] **Other programs closed** (Arduino IDE, serial monitors)
- [ ] **Different USB port tried**
- [ ] **Different USB cable tried**

### If Pattern Not Evolving
- [ ] **Code flashed completely** (not interrupted)
- [ ] **Board not stuck** in bootloader mode
- [ ] **Power cycled** (unplug and replug)

## ☐ Optimization

### Performance
- [ ] **Frame rate satisfactory**
  - [ ] Adjusted sleep duration if too fast/slow
  - [ ] Considered faster board if too slow (Pico/ESP32)

### Pattern
- [ ] **Interesting behavior observed**
  - [ ] Tried different patterns
  - [ ] Found favorite pattern
  - [ ] Considered creating custom patterns

### Physical Setup
- [ ] **Enclosure planned** (optional)
  - [ ] Measurements taken
  - [ ] Design sketched
  - [ ] Materials gathered

## ☐ Next Steps

### Enhancements
- [ ] **Buttons added** for pattern selection
- [ ] **Speed control** via potentiometer
- [ ] **Pause button** implemented
- [ ] **Statistics display** on OLED
- [ ] **Battery power** for portability

### Documentation
- [ ] **Photos taken** of working setup
- [ ] **Notes made** on configuration
- [ ] **Customizations documented**

### Sharing
- [ ] **Project shared** with friends/community
- [ ] **Improvements contributed** back to project
- [ ] **Experience documented** for others

---

## Quick Status Check

Mark your current status:

- [ ] **✅ FULLY WORKING** - Display shows Game of Life animation
- [ ] **⚙️ IN PROGRESS** - Hardware connected, troubleshooting
- [ ] **📋 PLANNING** - Gathering parts and information
- [ ] **🎯 READY TO START** - Have everything, ready to wire

---

## Need Help?

**Check these files:**
1. `QUICKSTART.md` - Fast commands and quick fixes
2. `WIRING_DIAGRAMS.md` - Visual connection guides
3. `SSD1306_PIN_GUIDE.md` - Pin-specific information
4. `README.md` - Complete documentation

**Common issues:**
- Blank display → Check power and I2C address
- Garbled display → Lower I2C frequency, shorten wires
- Can't flash → Check USB cable and drivers
- Slow performance → Normal for Arduino Uno, upgrade to Pico

**Still stuck?**
- Double-check wiring against diagrams
- Try the I2C scanner code (in `SSD1306_PIN_GUIDE.md`)
- Test with simpler code first (blink LED)
- Verify board with other TinyGO examples

---

**When everything works, celebrate! 🎉 You've built a physical cellular automaton!**
