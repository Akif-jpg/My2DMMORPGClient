package shapes

import (
	"math"
)

type Circle struct {
	Center Point
	Radius float64
}

func (c *Circle) GetType() string {
	return "Circle"
}

func (c *Circle) GetCenter() Point {
	return c.Center
}

func (c *Circle) GetRadius() float64 {
	return c.Radius
}

func NewCircle(center Point, radius float64) Circle {
	return Circle{Center: center, Radius: radius}
}

func (c *Circle) IntersectsCircle(other *Circle) bool {
	distance := c.Center.DistanceTo(&other.Center)
	return distance <= c.Radius+other.Radius
}

func (c *Circle) IntersectsRectangle(other *Rectangle) bool {
	distanceX := math.Abs(c.Center.X - other.Center.X)
	distanceY := math.Abs(c.Center.Y - other.Center.Y)

	if distanceX > c.Radius+other.Width/2 || distanceY > c.Radius+other.Height/2 {
		return false
	}

	if distanceX <= other.Width/2 || distanceY <= other.Height/2 {
		return true
	}

	cornerDistance := math.Pow(distanceX-other.Width/2, 2) + math.Pow(distanceY-other.Height/2, 2)
	return cornerDistance <= c.Radius*c.Radius
}

func (c *Circle) ContainsPoint(point *Point) bool {
	distance := c.Center.DistanceTo(point)
	return distance <= c.Radius
}

func (c *Circle) ContainsCircle(other *Circle) bool {
	distance := c.Center.DistanceTo(&other.Center)
	return distance <= c.Radius-other.Radius
}

func (c *Circle) ContainsRectangle(other *Rectangle) bool {
	distanceX := math.Abs(c.Center.X - other.Center.X)
	distanceY := math.Abs(c.Center.Y - other.Center.Y)

	if distanceX > c.Radius+other.Width/2 || distanceY > c.Radius+other.Height/2 {
		return false
	}

	if distanceX <= other.Width/2 || distanceY <= other.Height/2 {
		return true
	}

	cornerDistance := math.Pow(distanceX-other.Width/2, 2) + math.Pow(distanceY-other.Height/2, 2)
	return cornerDistance <= c.Radius*c.Radius
}

// IntersectsLine checks if the circle intersects with a line segment.
// It calculates the closest point on the line to the circle's center
// and checks if that distance is within the radius.
func (c *Circle) IntersectsLine(other *Line) bool {
	// Vector from line start to circle center
	dx := c.Center.X - other.Start.X
	dy := c.Center.Y - other.Start.Y

	// Vector from line start to line end
	lx := other.End.X - other.Start.X
	ly := other.End.Y - other.Start.Y

	// Length squared of the line segment
	lineLength := lx*lx + ly*ly

	// Handle degenerate case where line is actually a point
	if lineLength == 0 {
		return c.ContainsPoint(&other.Start)
	}

	// Calculate projection parameter t (0 <= t <= 1 means point is on segment)
	t := (dx*lx + dy*ly) / lineLength

	// Clamp t to [0, 1] to stay on the line segment
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	// Find the closest point on the line segment
	closestX := other.Start.X + t*lx
	closestY := other.Start.Y + t*ly

	// Calculate distance from circle center to closest point
	distX := c.Center.X - closestX
	distY := c.Center.Y - closestY
	distanceSquared := distX*distX + distY*distY

	return distanceSquared <= c.Radius*c.Radius
}

// IntersectsPoint checks if a point intersects (touches) the circle.
// A point intersects if it's on or inside the circle.
func (c *Circle) IntersectsPoint(point *Point) bool {
	return c.ContainsPoint(point)
}

// ContainsLine checks if the circle completely contains a line segment.
// Both endpoints must be within the circle.
func (c *Circle) ContainsLine(line *Line) bool {
	return c.ContainsPoint(&line.Start) && c.ContainsPoint(&line.End)
}
