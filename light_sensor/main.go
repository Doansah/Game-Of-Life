package main

import (
	"fmt"
	"machine"
	"time"
)

// What should this program do?
// 1) Turn on and Off my Govee Lightbulbs (via HTTP Request)
// 2) Utilitze Motion Sensor Data, to smartly turn on and off my lights when i enter a room
// 3)
//
// Development Process:
// Firstly Print light sensor data to the terminal...

func main() {
	// Define the motion sensor pin (GPIO23)
	msOutputPin := machine.GPIO23

	// Configure it as INPUT (we're reading from the sensor)
	msOutputPin.Configure(machine.PinConfig{Mode: machine.PinInput})

	fmt.Println("Motion Sensor Ready... ")

	// Counter to track motion detection events
	motionCounter := 0

	for {
		if msOutputPin.Get() {
			motionCounter++
			fmt.Println("Can Detect!!!!")
			fmt.Printf("Motion Detected - Count: %d\n", motionCounter)
		}
		time.Sleep(2000 * time.Millisecond)
	}
}
