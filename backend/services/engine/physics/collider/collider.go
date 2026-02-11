package collider

import (
	"math"

	"github.com/Akif-jpg/MyHobieMMORPGGame/services/engine/physics/geometry"
)

type Collider struct {
	ShapeList []geometry.Shape
	Transform geometry.Vector2
	Rotation  float64
	LayerMask bitmask
	MatchMask bitmask

	IsTrigger bool
	Enabled   bool
	Tag       string

	EntityID string
	UserData any
}

func (c *Collider) CanCollideWith(other *Collider) bool {
	if !c.Enabled || !other.Enabled {
		return false
	}
	return c.MatchMask.CanMatch(other.LayerMask)
}

func (c *Collider) GetShapes() []geometry.Shape {
	return c.ShapeList
}

func (c *Collider) SetTransform(transform geometry.Vector2) {
	c.Transform = transform
}

func (c *Collider) SetRotation(rotation float64) {
	c.Rotation = rotation
}

func (c *Collider) SetUserData(userData any) {
	c.UserData = userData
}

func (c *Collider) SetTag(tag string) {
	c.Tag = tag
}

func (c *Collider) SetEntityID(entityID string) {
	c.EntityID = entityID
}

// GetRotatedShapes rotates shapes around origin (0,0)
func (c *Collider) GetRotatedShapes() []geometry.Shape {
	rotatedShapes := make([]geometry.Shape, len(c.ShapeList))

	for i, shape := range c.ShapeList {
		rotatedShapes[i] = shape
		center := shape.GetCenter()

		// Rotate point around origin
		cos := math.Cos(c.Rotation)
		sin := math.Sin(c.Rotation)

		rotatedX := center.X*cos - center.Y*sin
		rotatedY := center.X*sin + center.Y*cos

		rotatedShapes[i].SetCenter(geometry.Point(geometry.Vector2{X: rotatedX, Y: rotatedY}))

		// If shape has rotation property (for rectangles, polygons)
		if rotatable, ok := rotatedShapes[i].(interface{ SetRotation(float64) }); ok {
			rotatable.SetRotation(c.Rotation)
		}
	}

	return rotatedShapes
}

// GetWorldSpaceShapes returns fully transformed shapes (rotation + translation)
func (c *Collider) GetWorldSpaceShapes() []geometry.Shape {
	worldShapes := make([]geometry.Shape, len(c.ShapeList))

	for i, shape := range c.ShapeList {
		center := shape.GetCenter()

		// First rotate around origin
		cos := math.Cos(c.Rotation)
		sin := math.Sin(c.Rotation)

		rotatedX := center.X*cos - center.Y*sin
		rotatedY := center.X*sin + center.Y*cos

		// Then translate
		worldShapes[i] = shape
		worldShapes[i].SetCenter(geometry.Point(geometry.Vector2{
			X: rotatedX + c.Transform.X,
			Y: rotatedY + c.Transform.Y,
		}))

		// Apply rotation to shape itself if supported
		if rotatable, ok := worldShapes[i].(interface{ SetRotation(float64) }); ok {
			rotatable.SetRotation(c.Rotation)
		}
	}

	return worldShapes
}

// AddTransform applies relative transformation
func (c *Collider) AddTransform(delta geometry.Vector2) {
	c.Transform.X += delta.X
	c.Transform.Y += delta.Y
}

// AddRotation applies relative rotation
func (c *Collider) AddRotation(deltaRotation float64) {
	c.Rotation += deltaRotation
}
