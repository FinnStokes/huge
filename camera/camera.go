package camera

// A Camera indicates an in-world rectangle to be drawn to an on-screen rectangle
type Camera struct {
	World  Rectangle
	Screen Screen
}

// Zoom performs a relative adjustment to the size of the in-world rectangle.
// A relZoom > 1 causes the rectangle to shrink (zooming in) while a relZoom < 1
// causes it to grow (zooming out)
func (c *Camera) Zoom(relZoom float32) {
	ow, oh := c.World.Width, c.World.Height
	c.World.Width /= relZoom
	c.World.Height /= relZoom

	c.World.X -= (c.World.Width - ow) / 2.0
	c.World.Y -= (c.World.Height - oh) / 2.0
}

// Focus returns the coordinates of the centre of the in-world rectangle
func (c *Camera) Focus() (x, y float32) {
	x = c.World.X + c.World.Width/2.0
	y = c.World.Y + c.World.Height/2.0
	return
}

// SetFocus adjusts the coordinates centre of the in-world rectangle
func (c *Camera) SetFocus(x, y float32) {
	c.World.X = x - c.World.Width/2.0
	c.World.Y = y - c.World.Height/2.0
}

// A Rectangle represents a rectangular region
type Rectangle struct {
	X, Y, Width, Height float32
}

// Intersects returns true if r and r2 overlap at any point
func (r *Rectangle) Intersects(r2 *Rectangle) bool {
	if r.X < r2.X+r2.Width && r.X+r.Width > r2.X && r.Y < r2.Y+r2.Height && r.Y+r.Height > r2.Y {
		return true
	}
	return false
}

// Contains returns true if r2 lies wholly within r
func (r *Rectangle) Contains(r2 *Rectangle) bool {
	if r.X < r2.X && r.X+r.Width > r2.X+r2.Width && r.Y < r2.Y && r.Y+r.Height > r2.Y+r2.Height {
		return true
	}
	return false
}

// ContainsPoint returns true if the point (x,y) lies within r
func (r *Rectangle) ContainsPoint(x, y float32) bool {
	if r.X < x && r.X+r.Width > x && r.Y < y && r.Y+r.Height > y {
		return true
	}
	return false
}

type Screen struct {
	Width, Height int
}
