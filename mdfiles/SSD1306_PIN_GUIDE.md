# SSD1306 OLED Pin Configuration Guide

## Display Module Pinout

Your 0.96" SSD1306 OLED typically has 4 pins:

```
┌─────────────────┐
│  128x64 OLED    │
│                 │
└─────────────────┘
  │  │  │  │
  │  │  │  └── SCL (I2C Clock)
  │  │  └───── SDA (I2C Data)
  │  └──────── GND (Ground)
  └─────────── VCC (Power: 3.3V or 5V)
```

## I2C Address Detection

Most SSD1306 displays use address **0x3C** (60 decimal)
Some displays use address **0x3D** (61 decimal)

If your display doesn't work, try changing the address in the code:

```go
display.Configure(ssd1306.Config{
    Address: 0x3D,  // Change from 0x3C to 0x3D
    Width:   128,
    Height:  64,
})
```

## Microcontroller Pin Mappings

### Raspberry Pi Pico (RP2040) - **RECOMMENDED**

```
Pico          Function    OLED
─────────────────────────────────
3.3V (Pin 36) → Power  →  VCC
GND (Pin 38)  → Ground →  GND
GP0 (Pin 1)   → SDA    →  SDA
GP1 (Pin 2)   → SCL    →  SCL
```

Flash command:
```bash
tinygo flash -target=pico tinygo_ssd1306_version.go
```

### Arduino Uno

```
Arduino       Function    OLED
─────────────────────────────────
5V            → Power  →  VCC
GND           → Ground →  GND
A4            → SDA    →  SDA
A5            → SCL    →  SCL
```

Flash command:
```bash
tinygo flash -target=arduino tinygo_ssd1306_version.go
```

### Arduino Nano

```
Nano          Function    OLED
─────────────────────────────────
5V            → Power  →  VCC
GND           → Ground →  GND
A4            → SDA    →  SDA
A5            → SCL    →  SCL
```

Flash command:
```bash
tinygo flash -target=arduino-nano tinygo_ssd1306_version.go
```

### ESP32 DevKit

```
ESP32         Function    OLED
─────────────────────────────────
3.3V          → Power  →  VCC
GND           → Ground →  GND
GPIO21        → SDA    →  SDA
GPIO22        → SCL    →  SCL
```

Flash command:
```bash
tinygo flash -target=esp32-coreboard tinygo_ssd1306_version.go
```

### Arduino Due

```
Due           Function    OLED
─────────────────────────────────
3.3V          → Power  →  VCC
GND           → Ground →  GND
SDA1 (Pin 20) → SDA    →  SDA
SCL1 (Pin 21) → SCL    →  SCL
```

Flash command:
```bash
tinygo flash -target=arduino-due tinygo_ssd1306_version.go
```

### STM32 Blue Pill

```
Blue Pill     Function    OLED
─────────────────────────────────
3.3V          → Power  →  VCC
GND           → Ground →  GND
PB7           → SDA    →  SDA
PB6           → SCL    →  SCL
```

Flash command:
```bash
tinygo flash -target=bluepill tinygo_ssd1306_version.go
```

## Custom Pin Configuration

If your board uses different I2C pins, modify the code:

```go
// Define custom pins
const (
    sdaPin = machine.GPIO0  // Change to your SDA pin
    sclPin = machine.GPIO1  // Change to your SCL pin
)

// Configure I2C with custom pins
machine.I2C0.Configure(machine.I2CConfig{
    Frequency: machine.TWI_FREQ_400KHZ,
    SDA:       sdaPin,
    SCL:       sclPin,
})
```

## Voltage Considerations

⚠️ **IMPORTANT**: Check your OLED module's voltage rating!

- **5V Tolerant**: Most modules with voltage regulator (usually have 4 pins)
  - Safe to use 5V from Arduino Uno/Nano
  
- **3.3V Only**: Some modules (usually with 7 pins)
  - **Must** use 3.3V from Pico, ESP32, or Due
  - Using 5V will **damage** the display!

**When in doubt**: Use 3.3V - it's safe for all modules

## Troubleshooting

### Display shows nothing
1. Check power (VCC/GND) connections
2. Try different I2C address (0x3C ↔ 0x3D)
3. Verify SDA/SCL are not swapped
4. Check voltage (3.3V vs 5V)

### Display shows garbled pixels
1. Check I2C frequency - try 100kHz instead of 400kHz
2. Use shorter wires (< 6 inches / 15 cm)
3. Add pull-up resistors (4.7kΩ on SDA and SCL)

### Can't flash code to board
1. Install correct USB drivers for your board
2. Check USB cable (must support data, not just power)
3. Press BOOTSEL button (Pico) or reset button during flash

### Display works but simulation is slow
- Normal for Arduino Uno (8MHz CPU)
- Increase `time.Sleep()` duration for smoother animation
- Or upgrade to faster board (Pico, ESP32)

## Testing I2C Connection

To scan for I2C devices (find the address):

```go
package main

import (
    "machine"
    "time"
)

func main() {
    machine.I2C0.Configure(machine.I2CConfig{
        Frequency: machine.TWI_FREQ_100KHZ,
    })
    
    println("Scanning I2C bus...")
    
    for addr := uint8(0); addr < 128; addr++ {
        buf := []byte{0}
        err := machine.I2C0.Tx(uint16(addr), buf, nil)
        if err == nil {
            println("Device found at address:", addr, "(0x", addr, ")")
        }
        time.Sleep(10 * time.Millisecond)
    }
    
    println("Scan complete")
}
```

## Pull-up Resistors

I2C requires pull-up resistors on SDA and SCL lines:

- Most breakout boards have built-in pull-ups ✅
- If not, add 4.7kΩ resistors between:
  - SDA → VCC
  - SCL → VCC

## Wire Length

- **Keep wires short**: < 6 inches (15 cm) preferred
- Longer wires may need lower I2C frequency
- Use twisted pair or shielded cable for > 12 inches

## Multiple I2C Devices

You can connect multiple I2C devices on the same bus:
- Each must have a unique address
- Add pull-ups only once (not per device)
- Maximum ~127 devices (address range)
