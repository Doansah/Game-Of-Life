
package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"log"
)

func main() {
	file, err := os.Open("/Users/dillonansah/Desktop/golang.png")
	if err != nil {
		log.Println("err with file opening")
	}
	defer file.Close()

	img, _, _ := image.Decode(file)
	b := img.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			
			// rgba returned as 0–65535; convert to 0–255
			fmt.Printf("(%d,%d,%d)\n", r/257, g/257, b/257)
		}
	}
}

