package shapes

import (
	"math"
)

type Rectangle struct {
	Center Point
	Width  float64
	Height float64
}

func NewRectangle(center Point, width, height float64) *Rectangle {
	return &Rectangle{
		Center: center,
		Width:  width,
		Height: height,
	}
}

func (r *Rectangle) GetType() string {
	return "Rectangle"
}

func (r *Rectangle) GetCenter() Point {
	return r.Center
}

func (r *Rectangle) GetWidth() float64 {
	return r.Width
}

func (r *Rectangle) GetHeight() float64 {
	return r.Height
}

func (r *Rectangle) IntersectsRectangle(other *Rectangle) bool {
	return !(r.Center.X+r.Width/2 <= other.Center.X-other.Width/2 ||
		r.Center.X-r.Width/2 >= other.Center.X+other.Width/2 ||
		r.Center.Y+r.Height/2 <= other.Center.Y-other.Height/2 ||
		r.Center.Y-r.Height/2 >= other.Center.Y+other.Height/2)
}

func (r *Rectangle) IntersectsCircle(other *Circle) bool {
	// Find the closest point on the rectangle to the circle's center
	distanceX := math.Abs(other.Center.X - r.Center.X)
	distanceY := math.Abs(other.Center.Y - r.Center.Y)

	// If circle center is too far away, no intersection
	if distanceX > (r.Width/2 + other.Radius) {
		return false
	}
	if distanceY > (r.Height/2 + other.Radius) {
		return false
	}

	// If circle center is close enough to the rectangle's center, there's intersection
	if distanceX <= r.Width/2 {
		return true
	}
	if distanceY <= r.Height/2 {
		return true
	}

	// Check corner distance
	cornerDistanceSquared := math.Pow(distanceX-r.Width/2, 2) + math.Pow(distanceY-r.Height/2, 2)
	return cornerDistanceSquared <= other.Radius*other.Radius
}

func (r *Rectangle) IntersectsLine(line *Line) bool {
	if r.ContainsPoint(&line.Start) && !r.ContainsPoint(&line.End) {
		return true
	}

	if r.ContainsPoint(&line.End) && !r.ContainsPoint(&line.Start) {
		return true
	}

	if r.ContainsPoint(&line.Start) && r.ContainsPoint(&line.End) {
		return true
	}

	edge1 := &Line{Start: Point{X: r.Center.X - r.Width/2, Y: r.Center.Y - r.Height/2}, End: Point{X: r.Center.X + r.Width/2, Y: r.Center.Y - r.Height/2}}
	edge2 := &Line{Start: Point{X: r.Center.X + r.Width/2, Y: r.Center.Y - r.Height/2}, End: Point{X: r.Center.X + r.Width/2, Y: r.Center.Y + r.Height/2}}
	edge3 := &Line{Start: Point{X: r.Center.X + r.Width/2, Y: r.Center.Y + r.Height/2}, End: Point{X: r.Center.X - r.Width/2, Y: r.Center.Y + r.Height/2}}
	edge4 := &Line{Start: Point{X: r.Center.X - r.Width/2, Y: r.Center.Y + r.Height/2}, End: Point{X: r.Center.X - r.Width/2, Y: r.Center.Y - r.Height/2}}

	if edge1.IntersectsLine(line) || edge2.IntersectsLine(line) || edge3.IntersectsLine(line) || edge4.IntersectsLine(line) {
		return true
	}

	return false
}

func (r *Rectangle) IntersectsPoint(point *Point) bool {
	return r.ContainsPoint(point)
}

func (r *Rectangle) ContainsRectangle(other *Rectangle) bool {
	return r.Center.X+r.Width/2 >= other.Center.X+other.Width/2 &&
		r.Center.X-r.Width/2 <= other.Center.X-other.Width/2 &&
		r.Center.Y+r.Height/2 >= other.Center.Y+other.Height/2 &&
		r.Center.Y-r.Height/2 <= other.Center.Y-other.Height/2
}

func (r *Rectangle) ContainsCircle(other *Circle) bool {
	// For a rectangle to contain a circle, all points on the circle must be inside the rectangle
	// This means the circle's center must be at least 'radius' distance away from all edges

	distanceX := math.Abs(r.Center.X - other.Center.X)
	distanceY := math.Abs(r.Center.Y - other.Center.Y)

	// Check if circle fits within rectangle bounds
	return (distanceX+other.Radius) <= r.Width/2 && (distanceY+other.Radius) <= r.Height/2
}

func (r *Rectangle) ContainsLine(line *Line) bool {
	return r.ContainsPoint(&line.Start) && r.ContainsPoint(&line.End)
}

func (r *Rectangle) ContainsPoint(point *Point) bool {
	distanceX := math.Abs(r.Center.X - point.X)
	distanceY := math.Abs(r.Center.Y - point.Y)

	// Point is inside if it's within both width and height bounds
	return distanceX <= r.Width/2 && distanceY <= r.Height/2
}
