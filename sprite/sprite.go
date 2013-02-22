package sprite

import (
	"image"
	"time"
)

type Sprite struct {
	Image            image.Image
	Animations       map[string]*Animation
	Width, Height    int
	CurrentAnimation *Animation
	CurrentFrame     int
	FrameTime        time.Duration
}

type Animation struct {
	Frames []int
	Fps    int
	Next   *Animation
}
