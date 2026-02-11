package geometry

const (
	CircleType    = "circle"
	RectangleType = "rectangle"
)

type Bounds struct {
	MinX, MinY, MaxX, MaxY float64
}

func (b Bounds) Width() float64  { return b.MaxX - b.MinX }
func (b Bounds) Height() float64 { return b.MaxY - b.MinY }

func (b Bounds) Intersects(other Bounds) bool {
	return !(b.MaxX < other.MinX || b.MinX > other.MaxX ||
		b.MaxY < other.MinY || b.MinY > other.MaxY)
}

type Shape interface {
	GetType() string
	GetCenter() Point
	SetCenter(center Point)

	GetBounds() Bounds

	IntersectsRectangle(other *Rectangle) bool
	IntersectsCircle(other *Circle) bool
	IntersectsLine(other *Line) bool
	IntersectsPoint(point *Point) bool

	ContainsRectangle(other *Rectangle) bool
	ContainsCircle(other *Circle) bool
	ContainsLine(other *Line) bool
	ContainsPoint(point *Point) bool
}
