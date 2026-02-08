package physics

import "github.com/Akif-jpg/MyHobieMMORPGGame/services/engine/physics/shapes"

// RectangleCollision checks if two rectangles are colliding.
func RectangleCollision(rect1, rect2 shapes.Rectangle) bool {
	return rect1.IntersectsRectangle(rect2) && rect1.ContainsRectangle(rect2) && rect2.ContainsRectangle(rect1)
}
