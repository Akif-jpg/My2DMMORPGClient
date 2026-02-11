package geometry

import "math"

type Point struct {
	X float64
	Y float64
}

func (p *Point) DistanceTo(point *Point) float64 {
	powX := (p.X - point.X) * (p.X - point.X)
	powY := (p.Y - point.Y) * (p.Y - point.Y)
	return math.Sqrt(powX + powY)
}

func NewPoint(x float64, y float64) *Point {
	return &Point{X: x, Y: y}
}
