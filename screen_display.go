package main

import ( "fmt"
					"math"
					"maps"
)


func main() {
	fmt.Println("hello world")

} 

struct RuleSet {
	haslegs bool
	jumpheight int
	position int

}



// FrameBuffer Struct
type FrameBuffer struct {
	pixelArray [128][64] pixel // where 128 is xcord and 64 is ycord
	isDisplayed bool
	isProcessing bool 
}


func (fb FrameBuffer) update FrameBuffer {
	
}


// Return a array of the indicies of all possible neighbors
func (Pixel) findNeighborLocations   {
	neighbors := make(map[])
	// 1)compute all of all 8 neighbors indexes
	
	pixel.xcord
	// 2) ensure that the indexes 'wrap around screen' modules? 
	xRemainder = math.Mod(pixel.xcord, 127)
	yRemainder =  math.Mod(pixel.ycord, 63 )
	
	if xRemainder == 0 {
		
	}

	// all possible combinations 
	// (a,b) = (a+1, b) (a-1, b) (a, b+1) (a, b-1)  | (a+1, b-1) (a+1, b+1) (a -1, b- 1 )
	// Create Neighbors
	//
	var x = pixel.xcord
	
	neighbors := [][]int {
	{x+1, y}, //  right middle
	{x-1, y }, // left middle
	{x, y+1}, // top middle
	{x, y-1}, // bottom middle
	{x+1, y+1}, // top right
	{x+1, y-1}, // bottom right
	{x-1, y+1}, // top left
	{x -1, y-1} // bottom left
	}


}

func (pixels []Pixel) countNeighbors uint8 {
	count := len(pixels)
	return count
}

func (neighbors [][]int ) validateNeighborIndicies [][]int {
	for neighbors.len() -1;
	// apply wrapping to ensure the indicies are correct 
	
}

/*
* Given our array Size 128 x 64 = 8192 = 2^(13)
* we can represent pixel location in a 14 byte array
*/

type Pixel struct {
	isLit bool 
	xcord uint8
	ycord unit8

}

func (Pixel pixel) validateLocation bool {
	// ensure x & y cordinate are within range of 128 x 64
	isValid = false
	if pixel.xcord > 128 && pixel.ycord = 64 {
		isValid = true
	}
	return isValid

} 


// Anonymous Structs
job := struct {
	title string

	salary int
}

// When do we want to deal with pointers
// A pointer is a value that points to memory address

jobPtr := &jobPtr
jobPrt.salary = 100000
 

gdbr 1




