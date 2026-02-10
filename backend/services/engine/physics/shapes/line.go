package shapes

type Line struct {
	Start Point
	End   Point
}

func NewLine(start, end Point) Line {
	return Line{
		Start: start,
		End:   end,
	}
}

func (l *Line) Length() float64 {
	return l.Start.DistanceTo(&l.End)
}

func (l *Line) IntersectsLine(other *Line) bool {
	p1, p2 := l.Start, l.End
	p3, p4 := other.Start, other.End

	o1 := orientation(p1, p2, p3)
	o2 := orientation(p1, p2, p4)
	o3 := orientation(p3, p4, p1)
	o4 := orientation(p3, p4, p2)

	if o1*o2 < 0 && o3*o4 < 0 {
		return true
	}

	if o1 == 0 && onSegment(p1, p2, p3) {
		return true
	}
	if o2 == 0 && onSegment(p1, p2, p4) {
		return true
	}
	if o3 == 0 && onSegment(p3, p4, p1) {
		return true
	}
	if o4 == 0 && onSegment(p3, p4, p2) {
		return true
	}

	return false
}

// Helper functions for line intersection detection
func orientation(a, b, c Point) float64 {
	return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
}

func onSegment(a, b, c Point) bool {
	return c.X >= min(a.X, b.X) &&
		c.X <= max(a.X, b.X) &&
		c.Y >= min(a.Y, b.Y) &&
		c.Y <= max(a.Y, b.Y)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// GetType returns the type name of the shape.
func (l *Line) GetType() string {
	return "Line"
}

// GetCenter returns the midpoint of the line segment.
func (l *Line) GetCenter() Point {
	return Point{
		X: (l.Start.X + l.End.X) / 2,
		Y: (l.Start.Y + l.End.Y) / 2,
	}
}

// IntersectsRectangle checks if the line intersects with a rectangle.
// This delegates to Rectangle.IntersectsLine for the actual implementation.
func (l *Line) IntersectsRectangle(other *Rectangle) bool {
	return other.IntersectsLine(l)
}

// IntersectsCircle checks if the line intersects with a circle.
// This delegates to Circle.IntersectsLine for the actual implementation.
func (l *Line) IntersectsCircle(other *Circle) bool {
	return other.IntersectsLine(l)
}

// IntersectsPoint checks if a point lies on the line segment.
func (l *Line) IntersectsPoint(point *Point) bool {
	return l.ContainsPoint(point)
}

// ContainsRectangle always returns false because a line cannot contain a 2D shape.
func (l *Line) ContainsRectangle(other *Rectangle) bool {
	return false
}

// ContainsCircle always returns false because a line cannot contain a 2D shape.
func (l *Line) ContainsCircle(other *Circle) bool {
	return false
}

// ContainsLine checks if the other line is completely contained within this line.
// For a line to contain another, they must be collinear and the other line's
// endpoints must lie on this line segment.
func (l *Line) ContainsLine(other *Line) bool {
	// Check if both endpoints of the other line are on this line
	return l.ContainsPoint(&other.Start) && l.ContainsPoint(&other.End)
}

// ContainsPoint checks if a point lies exactly on the line segment.
// Uses the cross product to check collinearity and bounding box to check if on segment.
func (l *Line) ContainsPoint(point *Point) bool {
	// Check if point is within the bounding box of the line
	if point.X < min(l.Start.X, l.End.X) || point.X > max(l.Start.X, l.End.X) ||
		point.Y < min(l.Start.Y, l.End.Y) || point.Y > max(l.Start.Y, l.End.Y) {
		return false
	}

	// Check if the point is collinear with the line using cross product
	// Cross product should be zero (or very close to zero) if collinear
	crossProduct := (point.Y-l.Start.Y)*(l.End.X-l.Start.X) - (point.X-l.Start.X)*(l.End.Y-l.Start.Y)

	// Use a small epsilon for floating point comparison
	const epsilon = 1e-10
	return crossProduct*crossProduct < epsilon
}
