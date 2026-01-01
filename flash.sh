#!/bin/bash

# Conway's Game of Life - TinyGO Build Script
# This script helps you flash the Game of Life to your microcontroller

echo "╔═══════════════════════════════════════════════════╗"
echo "║   Conway's Game of Life - SSD1306 OLED Builder   ║"
echo "╚═══════════════════════════════════════════════════╝"
echo ""

# Check if TinyGO is installed
if ! command -v tinygo &> /dev/null; then
    echo "❌ TinyGO is not installed!"
    echo ""
    echo "Install TinyGO first:"
    echo "  macOS:   brew tap tinygo-org/tools && brew install tinygo"
    echo "  Linux:   See https://tinygo.org/getting-started/install/"
    echo "  Windows: Download from https://github.com/tinygo-org/tinygo/releases"
    exit 1
fi

echo "✅ TinyGO found: $(tinygo version)"
echo ""

# Select target board
echo "Select your microcontroller:"
echo ""
echo "  1) Raspberry Pi Pico (RP2040) ⭐ RECOMMENDED"
echo "  2) Arduino Uno"
echo "  3) Arduino Nano"
echo "  4) Arduino Due"
echo "  5) ESP32"
echo "  6) STM32 Blue Pill"
echo "  7) Custom target"
echo ""
read -p "Enter choice (1-7): " choice

case $choice in
    1)
        TARGET="pico"
        BOARD="Raspberry Pi Pico"
        ;;
    2)
        TARGET="arduino"
        BOARD="Arduino Uno"
        ;;
    3)
        TARGET="arduino-nano"
        BOARD="Arduino Nano"
        ;;
    4)
        TARGET="arduino-due"
        BOARD="Arduino Due"
        ;;
    5)
        TARGET="esp32-coreboard"
        BOARD="ESP32"
        ;;
    6)
        TARGET="bluepill"
        BOARD="STM32 Blue Pill"
        ;;
    7)
        read -p "Enter custom target name: " TARGET
        BOARD="$TARGET"
        ;;
    *)
        echo "❌ Invalid choice"
        exit 1
        ;;
esac

echo ""
echo "📋 Selected board: $BOARD"
echo "🎯 Target: $TARGET"
echo ""

# Select pattern
echo "Select starting pattern:"
echo ""
echo "  1) Glider (small moving pattern)"
echo "  2) Blinker (oscillator)"
echo "  3) Toad (oscillator)"
echo "  4) Pulsar (large oscillator)"
echo "  5) Lightweight Spaceship (spaceship)"
echo "  6) Random"
echo ""
read -p "Enter choice (1-6): " pattern_choice

case $pattern_choice in
    1) PATTERN="glider" ;;
    2) PATTERN="blinker" ;;
    3) PATTERN="toad" ;;
    4) PATTERN="pulsar" ;;
    5) PATTERN="lightweight_spaceship" ;;
    6) PATTERN="random" ;;
    *)
        echo "Invalid choice, using glider"
        PATTERN="glider"
        ;;
esac

echo ""
echo "🎨 Pattern: $PATTERN"
echo ""

# Check if we need to modify the pattern in the code
if [ "$PATTERN" != "glider" ]; then
    echo "📝 Note: Edit tinygo_ssd1306_version.go line ~243 to change pattern"
    echo "   Current: grid := NewGridWithPattern(\"glider\")"
    echo "   Change to: grid := NewGridWithPattern(\"$PATTERN\")"
    echo ""
    read -p "Press Enter to continue or Ctrl+C to edit file first..."
fi

# Build and flash
echo ""
echo "🔨 Building and flashing to $BOARD..."
echo ""
echo "Command: tinygo flash -target=$TARGET tinygo_ssd1306_version.go"
echo ""

# For Pico, show BOOTSEL instructions
if [ "$TARGET" = "pico" ]; then
    echo "⚠️  IMPORTANT for Raspberry Pi Pico:"
    echo "   1. Hold BOOTSEL button"
    echo "   2. Connect USB cable"
    echo "   3. Release BOOTSEL"
    echo "   4. Pico will appear as USB drive"
    echo ""
    read -p "Press Enter when Pico is in BOOTSEL mode..."
    echo ""
fi

# Flash the board
tinygo flash -target=$TARGET tinygo_ssd1306_version.go

# Check result
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Successfully flashed to $BOARD!"
    echo ""
    echo "📺 Your Game of Life should now be running on the OLED display!"
    echo ""
    echo "🔧 Wiring reminder:"
    echo "   VCC → 3.3V or 5V"
    echo "   GND → GND"
    echo "   SDA → SDA pin"
    echo "   SCL → SCL pin"
    echo ""
    echo "📖 See SSD1306_PIN_GUIDE.md for detailed wiring"
    echo ""
    
    # Ask if user wants to monitor serial output
    read -p "Monitor serial output? (y/n): " monitor_choice
    if [ "$monitor_choice" = "y" ] || [ "$monitor_choice" = "Y" ]; then
        echo ""
        echo "Starting serial monitor (Ctrl+C to exit)..."
        tinygo monitor
    fi
else
    echo ""
    echo "❌ Flash failed!"
    echo ""
    echo "Common issues:"
    echo "  • Check USB connection"
    echo "  • Install USB drivers for your board"
    echo "  • For Pico: Enter BOOTSEL mode"
    echo "  • Try different USB cable (must support data)"
    echo "  • Close other programs using the serial port"
    echo ""
    exit 1
fi
