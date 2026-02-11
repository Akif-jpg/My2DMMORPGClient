package geometry

import (
	"math"
)

type Vector2 struct {
	X, Y float64
}

func NewVector2(x, y float64) *Vector2 {
	return &Vector2{X: x, Y: y}
}

func (v *Vector2) Add(other *Vector2) *Vector2 {
	return &Vector2{X: v.X + other.X, Y: v.Y + other.Y}
}

func (v *Vector2) Subtract(other *Vector2) *Vector2 {
	return &Vector2{X: v.X - other.X, Y: v.Y - other.Y}
}

func (v *Vector2) Multiply(scalar float64) *Vector2 {
	return &Vector2{X: v.X * scalar, Y: v.Y * scalar}
}

func (v *Vector2) Divide(scalar float64) *Vector2 {
	return &Vector2{X: v.X / scalar, Y: v.Y / scalar}
}

func (v *Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector2) Normalize() *Vector2 {
	length := v.Length()
	if length == 0 {
		return &Vector2{X: 0, Y: 0}
	}
	return &Vector2{X: v.X / length, Y: v.Y / length}
}

func (v *Vector2) Rotate(angle float64) *Vector2 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return &Vector2{X: v.X*cos - v.Y*sin, Y: v.X*sin + v.Y*cos}
}

func (v *Vector2) Transform(vector Vector2) *Vector2 {
	return &Vector2{X: v.X + vector.X, Y: v.Y + vector.Y}
}

func (v *Vector2) ApplyTransformToPoint(point Point) Point {
	return Point{X: point.X + v.X, Y: point.Y + v.Y}
}
