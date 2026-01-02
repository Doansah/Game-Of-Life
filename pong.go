package main

func main() {
	current = newFB(current) 
	next = newFB(next)
	
}

const (
	Height = 64
	Width  = 128
	BallArea = 8
)

type Framebuffer struct {
	name  string
	cells [Height][Width]bool
}

// Constructor for Framebuffer
func newFB(name string) Framebuffer {
	return Framebuffer{
		name,
		cells[Height][Width]
	},
}
const (
	BALLAREA = 9
)
type Ball struct {
	ballArea uint
	location_x uint
	location_y uint
	direction_x uint
	direction_y uint 
}

func newBall(ballArea, location, direction_y, direction_x uint) Ball {
	return Ball{
		ballArea,
		location_x,
		location_y
		direction_x
		direction_y
	}
}
//Ball and Upper Lower Collision
func detectBorderCollision(location_x, location_y uint) bool{
	
	if (location_x == 1 || location_x == SCREENHEIGHT -1 ) {
		return true
	}
	return false 

}


const (
PLENGTH = 8
PHEIGHT = 16
)

type Player struct {
	length uint
	height uint 
	location uint
}
// Ball and Player Collision

func detectPlayerCollision(p1Location, p2Location, ballx, bally uint) (bool Player) {
	// For players I need to account for the 7 pixel  positions below...
	

	
}

// Player Details



// Questions for Ball Movement
// How do the physics work (an angle vector?) (Movement Speed?) 
//
// Design Decissions 
// Ball Bouncing Mechanics:
// Bounces back on horizontal walls 
// When a Collision Occurs on Vertical Wall 
