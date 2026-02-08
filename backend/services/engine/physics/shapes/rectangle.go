package shapes

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

func (r *Rectangle) IntersectsRectangle(other Rectangle) bool {
	return !(r.Center.X+r.Width/2 <= other.Center.X-other.Width/2 ||
		r.Center.X-r.Width/2 >= other.Center.X+other.Width/2 ||
		r.Center.Y+r.Height/2 <= other.Center.Y-other.Height/2 ||
		r.Center.Y-r.Height/2 >= other.Center.Y+other.Height/2)
}

func (r *Rectangle) ContainsRectangle(other Rectangle) bool {
	return r.Center.X+r.Width/2 >= other.Center.X+other.Width/2 &&
		r.Center.X-r.Width/2 <= other.Center.X-other.Width/2 &&
		r.Center.Y+r.Height/2 >= other.Center.Y+other.Height/2 &&
		r.Center.X-r.Height/2 <= other.Center.Y-other.Height/2
}
