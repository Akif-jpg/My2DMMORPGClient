package collision

import (
	"github.com/Akif-jpg/MyHobieMMORPGGame/services/engine/physics/collider"
	"github.com/Akif-jpg/MyHobieMMORPGGame/services/engine/physics/geometry"
)

type CollisionBody struct {
	transform geometry.Vector2
	Radius    float64
	collider  collider.Collider
}

func NewCollisionBody(transform geometry.Vector2, radius float64, collider collider.Collider) *CollisionBody {
	return &CollisionBody{
		transform: transform,
		Radius:    radius,
		collider:  collider,
	}
}
