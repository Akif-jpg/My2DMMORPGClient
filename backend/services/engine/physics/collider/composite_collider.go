package collider

import (
	"fmt"
	"math"

	"github.com/Akif-jpg/MyHobieMMORPGGame/services/engine/physics/geometry"
)

// CompositeCollider - Multi-part collider for complex entities (bosses, vehicles)
type CompositeCollider struct {
	*Collider
	BodyParts map[string]*BodyPart
}

// BodyPart represents a single part of a composite collider
type BodyPart struct {
	LocalShape    geometry.Shape   // Shape relative to parent
	LocalOffset   geometry.Vector2 // Offset from parent center
	LocalRotation float64          // Rotation relative to parent
}

// NewCompositeCollider creates a new composite collider
func NewCompositeCollider(transform geometry.Vector2, rotation float64) *CompositeCollider {
	return &CompositeCollider{
		Collider: &Collider{
			Transform: transform,
			Rotation:  rotation,
			Enabled:   true,
			ShapeList: []geometry.Shape{}, // Empty, will be computed
		},
		BodyParts: make(map[string]*BodyPart),
	}
}

// AddBodyPart adds a new body part to the composite collider
func (cc *CompositeCollider) AddBodyPart(name string, shape geometry.Shape, offset geometry.Vector2, rotation float64) {
	cc.BodyParts[name] = &BodyPart{
		LocalShape:    shape,
		LocalOffset:   offset,
		LocalRotation: rotation,
	}
}

// RemoveBodyPart removes a body part by name
func (cc *CompositeCollider) RemoveBodyPart(name string) {
	delete(cc.BodyParts, name)
}

// GetBodyPart returns a body part by name
func (cc *CompositeCollider) GetBodyPart(name string) (*BodyPart, bool) {
	part, exists := cc.BodyParts[name]
	return part, exists
}

// UpdateBodyPart updates an existing body part's properties
func (cc *CompositeCollider) UpdateBodyPart(name string, offset *geometry.Vector2, rotation *float64) error {
	part, exists := cc.BodyParts[name]
	if !exists {
		return fmt.Errorf("body part '%s' not found", name)
	}

	if offset != nil {
		part.LocalOffset = *offset
	}
	if rotation != nil {
		part.LocalRotation = *rotation
	}

	return nil
}

// GetWorldSpaceShapes returns all body parts transformed to world space
func (cc *CompositeCollider) GetWorldSpaceShapes() []geometry.Shape {
	var allShapes []geometry.Shape

	parentRotation := cc.Rotation
	parentPos := cc.Transform

	cos := math.Cos(parentRotation)
	sin := math.Sin(parentRotation)

	for _, part := range cc.BodyParts {
		// Step 1: Rotate local offset around origin by parent rotation
		rotatedOffsetX := part.LocalOffset.X*cos - part.LocalOffset.Y*sin
		rotatedOffsetY := part.LocalOffset.X*sin + part.LocalOffset.Y*cos

		// Step 2: Get shape's local center
		localCenter := part.LocalShape.GetCenter()

		// Step 3: Calculate total rotation (parent + local)
		totalRotation := parentRotation + part.LocalRotation

		// Step 4: Rotate shape center by total rotation
		totalCos := math.Cos(totalRotation)
		totalSin := math.Sin(totalRotation)

		rotatedCenterX := localCenter.X*totalCos - localCenter.Y*totalSin
		rotatedCenterY := localCenter.X*totalSin + localCenter.Y*totalCos

		// Step 5: Create world space shape
		worldShape := part.LocalShape
		worldShape.SetCenter(geometry.Point(geometry.Vector2{
			X: parentPos.X + rotatedOffsetX + rotatedCenterX,
			Y: parentPos.Y + rotatedOffsetY + rotatedCenterY,
		}))

		// Apply rotation if shape supports it
		if rotatable, ok := worldShape.(interface{ SetRotation(float64) }); ok {
			rotatable.SetRotation(totalRotation)
		}

		allShapes = append(allShapes, worldShape)
	}

	return allShapes
}

// GetShapes returns world space shapes (override parent method)
func (cc *CompositeCollider) GetShapes() []geometry.Shape {
	return cc.GetWorldSpaceShapes()
}

// GetBounds returns the AABB bounding box containing all body parts
func (cc *CompositeCollider) GetBounds() geometry.Rectangle {
	shapes := cc.GetWorldSpaceShapes()
	if len(shapes) == 0 {
		return geometry.Rectangle{
			Center: geometry.Point(cc.Transform),
			Width:  0,
			Height: 0,
		}
	}

	// Calculate AABB from all shapes
	firstBounds := shapes[0].GetBounds()
	minX := firstBounds.MinX
	minY := firstBounds.MinY
	maxX := firstBounds.MaxX
	maxY := firstBounds.MaxY

	for i := 1; i < len(shapes); i++ {
		bounds := shapes[i].GetBounds()
		minX = math.Min(minX, bounds.MinX)
		minY = math.Min(minY, bounds.MinY)
		maxX = math.Max(maxX, bounds.MaxX)
		maxY = math.Max(maxY, bounds.MaxY)
	}

	return geometry.Rectangle{
		Center: *geometry.NewPoint(minX, minY),
		Width:  maxX - minX,
		Height: maxY - minY,
	}
}
